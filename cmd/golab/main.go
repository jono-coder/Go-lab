package main

import (
	"Go-lab/config"
	"Go-lab/internal/client"
	"Go-lab/internal/contact"
	myMiddleware "Go-lab/internal/middleware"
	"Go-lab/internal/player"
	"Go-lab/internal/security"
	"Go-lab/internal/utils"
	"Go-lab/internal/utils/dbutils"
	"Go-lab/internal/utils/httpconst"
	"Go-lab/internal/utils/session"
	"Go-lab/internal/utils/session/session_db"
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	dbUtils         *dbutils.DbUtils
	clientRepo      *client.Repo
	clientService   *client.Service
	serviceRegistry *utils.ServiceRegistry
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	////////// plumbing //////////

	backgroundCtx := context.Background()

	server, err := initialise(backgroundCtx)
	if err != nil {
		slog.Error("failed to start application", "error", err)
		panic(err)
	}

	// listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := server.ListenAndServe(); err != nil {
			slog.Error("failed to listenAndServe", "error", err)
			cancel()
		}
	}()

	go func() {
		<-quit
		slog.Info("signal received, shutting down...")
		cancel() // unblock context-aware goroutines
	}()

	// block until context is done
	<-ctx.Done()

	destroy()
}

func initialise(ctx context.Context) (*http.Server, error) {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		panic(err)
	}
	slog.Info("Running in env", slog.String("env", cfg.App.Env))

	dbUtils = dbutils.NewDbUtils(&cfg.DB)
	/*if err = dbUtils.Init(); err != nil {
		slog.Error("db init", "error", err)
		panic(err)
	}*/

	// now load the database if scripts are present
	loader := dbutils.NewDbLoader(ctx, dbUtils)
	if cfg.App.IsDev() {
		err = loader.Load(ctx, "db_scripts.dev")
	} else {
		err = loader.Load(ctx, "db_scripts.prod")
	}
	if err != nil {
		slog.Error("db load", "error", err)
		panic(err)
	}

	oauthConfig := security.NewOAuthConfig(ctx, cfg.App.BaseUrl)

	serviceRegistry = utils.NewServiceRegistry()
	////////// plumbing //////////

	////////// client //////////
	clientRepo = client.NewRepo(dbUtils)
	clientApi, err := client.NewAPI(oauthConfig)
	if err != nil {
		slog.Error("client.NewAPI", "error", err)
		panic(err)
	}
	clientService = client.NewService(dbUtils, clientRepo, clientApi)
	serviceRegistry.Register(clientService)
	////////// client //////////

	////////// player //////////
	playerRepo := player.NewRepo(dbUtils)
	playerApi, err := player.NewAPI(oauthConfig)
	if err != nil {
		slog.Error("client.NewAPI", "error", err)
		panic(err)
	}
	playerService := player.NewService(dbUtils, playerRepo, playerApi)
	////////// player //////////

	////////// contact //////////
	contactRepo := contact.NewRepo(dbUtils)
	contactService := contact.NewService(dbUtils, contactRepo)
	serviceRegistry.Register(contactService)
	////////// contact //////////

	serviceRegistry.StartAll()

	////////// router //////////
	router := chi.NewRouter()
	router.Use(middleware.Compress(5, httpconst.ContentTypeJson, httpconst.ContentTypeXml,
		httpconst.ContentTypeHtml, httpconst.ContentTypeText))
	router.Use(middleware.Logger)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(myMiddleware.SecureHandler)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Throttle(int(cfg.App.Throttle)))
	router.Use(middleware.Timeout(cfg.App.TimeoutInSeconds))

	/*router.Get(cfg.App.Root, func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Go is listening")); err != nil {
			slog.Error("write to response.writer", "error", err)
		}
	})*/

	// static files will be mounted after API routes to avoid shadowing them

	clientHandler := client.NewHandler(ctx, clientService, cfg.App)
	router.With(middleware.NoCache).Route(cfg.App.Root+"/client", func(r chi.Router) {
		r.With(utils.Low).Get("/", clientHandler.List)
		r.Get("/{id}", clientHandler.Get)
		r.Post("/", clientHandler.Create)
		r.Put("/{id}", clientHandler.Update)
		r.Delete("/{id}", clientHandler.Delete)
	})
	router.With(middleware.NoCache).Route(cfg.App.Root+"/session", func(r chi.Router) {
		r.Get("/currentUserId", func(w http.ResponseWriter, r *http.Request) {
			ctx := session.ContextWithUserID(r.Context(), 1001)

			err := dbUtils.WithTransaction(ctx, func(tx *sql.Tx) error {
				if userId, err := session_db.GetUserIdFromDB(ctx, tx); err != nil {
					slog.Error("session.GetUserIdFromDB", "error", err)
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("error"))
				} else {
					w.Header().Set(httpconst.HeaderContentType, httpconst.ContentTypeJson)
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(strconv.Itoa(*userId)))
				}
				return err
			})
			if err != nil {
				return
			}
		})
	})
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetReqID(r.Context())
		pong := fmt.Sprintf("pong - request id: %s; IP=%s", requestID, r.RemoteAddr)
		w.Write([]byte(pong))
	})

	playerHandler := player.NewHandler(ctx, playerService, cfg.App)
	router.Route(cfg.App.Root+"/player", func(r chi.Router) {
		r.With(middleware.NoCache).Get("/", playerHandler.List)
		r.With(middleware.NoCache).Get("/{id}", playerHandler.Get)
		r.With(middleware.NoCache).Get("/resource/{resource_id}", playerHandler.GetResource)
		r.Put("/{id}", playerHandler.Checkin)
	})

	oauthHandler := security.NewHandler(ctx, cfg.App)
	router.Route(cfg.App.Root+"/security", func(r chi.Router) {
		r.Post("/oauth/token", oauthHandler.Auth)
	})
	////////// router //////////

	
	fileServer := http.FileServer(http.Dir("./web"))

	router.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := "./web" + r.URL.Path

		if _, err := os.Stat(path); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}

		http.ServeFile(w, r, "./web/index.html")
	}))

	handlerWithCors := cors(router)

	port := int(cfg.App.Port)
	slog.Info("starting server on port", slog.Int("port", port))
	return &http.Server{
		Addr:              ":" + strconv.Itoa(port),
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
		Handler:           handlerWithCors,
	}, nil
}

////////// CORS //////////

var allowedOrigins = []string{
    "http://localhost:5173", // Vite default
}
var allowedMethods = "GET,POST,PUT,PATCH,DELETE,OPTIONS"
var allowedHeaders = "Content-Type, Authorization, X-Requested-With"
var allowCredentials = true

func cors(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")

        if origin != "" && isAllowedOrigin(origin) {
            w.Header().Set("Vary", "Origin")
            w.Header().Set("Access-Control-Allow-Origin", origin)
            w.Header().Set("Access-Control-Allow-Methods", allowedMethods)
            w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
            if allowCredentials {
                w.Header().Set("Access-Control-Allow-Credentials", "true")
            }
        }

        // ALWAYS handle OPTIONS after headers are set
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func isAllowedOrigin(origin string) bool {
    for _, o := range allowedOrigins {
        if strings.EqualFold(o, origin) {
            return true
        }
    }
    return false
}

func destroy() {
	serviceRegistry.StopAll()
	clientRepo.Close()
	dbUtils.Close()
}

func findById(ctx context.Context, id uint) {
	slog.Info("findById...", slog.Int("id", int(id)))

	if entity, err := clientService.FindById(ctx, id); err != nil {
		slog.Error("clientService.FindById", "error", err)
	} else {
		slog.Info("clientService.FindById", "entity", entity)
	}

	slog.Info("clientService.FindById", slog.Int("id", int(id)))
}

func doBusinessStuff(ctx context.Context) {
	slog.Info("doBusinessStuff()...")

	if err := clientService.DoBusinessStuff(ctx); err != nil {
		slog.Error("clientService.DoBusinessStuff", "error", err)
	}

	slog.Info("doBusinessStuff().")
}

func doSomeStuffInTheWorkerPool(ctx context.Context) {
	slog.Info("doSomeStuffInTheWorkerPool()...")

	workerPool := utils.NewWorkerPool(10, 100)

	for i := 1; i <= 10_000; i++ {
		i := i
		err := workerPool.Submit(func(ctx context.Context) error {
			// CPU-bound work
			sum := 0
			for n := 0; n < 1_000_000; n++ {
				sum += n % 3

				if n%10_000 == 0 {
					if err := ctx.Err(); err != nil {
						slog.Error("task exiting early due to cancellation")
						return err
					}
				}
			}

			// Simulate a single failure
			if i == 7 {
				return fmt.Errorf("task %d error", i)
			}

			return nil
		})

		// Stop submission if the pool has been cancelled
		if err != nil {
			slog.Error("submit aborted", "error", err)
			break
		}
	}

	slog.Info("after for loop")

	start := time.Now()
	err := workerPool.Wait()
	duration := time.Since(start)

	if err != nil {
		fmt.Println("Pool finished with error:", err)
	} else {
		fmt.Println("All tasks completed successfully")
	}

	fmt.Printf("Total execution time: %s\n", duration)

	slog.Info("doSomeStuffInTheWorkerPool().")
}

func test1Api() {
	slog.Info("test1Api()...")

	if c, err := clientService.GetById(1); err != nil {
		slog.Error("clientService.GetById", "error", err)
	} else {
		slog.Info("clients retrieved", "clients", c)
	}

	slog.Info("test1Api()!")
}

func test2Api() {
	if c, err := clientService.GetAll(); err != nil {
		slog.Error(err.Error(), "error", err)
	} else {
		slog.Info("clients retrieved", "clients", c)
	}
}
