package main

import (
	"Go-lab/config"
	"Go-lab/internal/client"
	"Go-lab/internal/contact"
	myMiddleware "Go-lab/internal/middleware"
	"Go-lab/internal/player"
	"Go-lab/internal/security"
	"Go-lab/internal/utils"
	"Go-lab/internal/utils/httpconst"
	"log/slog"
	"path/filepath"
	"strconv"
	"strings"

	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	dbUtils       *utils.DbUtils
	clientRepo    *client.Repo
	clientService *client.Service
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

	dbUtils = utils.NewDbUtils(&cfg.DB)
	if err = dbUtils.Init(); err != nil {
		slog.Error("db init", "error", err)
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
	router.Route(cfg.App.Root+"/client", func(r chi.Router) {
		r.With(utils.Medium).Get("/", clientHandler.List)
		r.With(utils.Low).Get("/{id}", clientHandler.Get)
		r.With(middleware.NoCache).Post("/", clientHandler.Create)
		r.With(middleware.NoCache).Put("/{id}", clientHandler.Update)
		r.With(middleware.NoCache).Delete("/{id}", clientHandler.Delete)
	})

	playerHandler := player.NewHandler(ctx, playerService, cfg.App)
	router.Route(cfg.App.Root+"/player", func(r chi.Router) {
		r.With(utils.Medium).Get("/", playerHandler.List)
		r.With(middleware.NoCache).Get("/{id}", playerHandler.Get)
		r.With(utils.Low).Get("/resource/{resource_id}", playerHandler.GetResource)
		r.Put("/{id}", playerHandler.Checkin)
	})

	oauthHandler := security.NewHandler(ctx, cfg.App)
	router.Route(cfg.App.Root+"/security", func(r chi.Router) {
		r.Post("/oauth/token", oauthHandler.Auth)
	})
	////////// router //////////

	// serve static files from the project's `web` folder
	static := http.FileServer(http.Dir("./web"))

	// serve absolute asset paths (e.g. /assets/...) so built files referencing /assets work
	assetsFS := http.FileServer(http.Dir("./web/assets"))
	router.Handle("/assets/*", http.StripPrefix("/assets/", assetsFS))

	// favicon
	router.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/favicon.ico")
	})

	// serve the app root (SPA). This also allows /lab/... paths to resolve inside the web folder.
	if cfg.App.Root == "/" {
		router.Mount("/", static)
	} else {
		router.Mount(cfg.App.Root, http.StripPrefix(cfg.App.Root, static))
		// Optional: redirect plain root to app root so navigating to http://localhost:PORT/ goes to the SPA
		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, cfg.App.Root, http.StatusFound)
		})

		// SPA fallback: for any path under the app root that doesn't match a file, serve index.html
		router.Get(cfg.App.Root+"/*", func(w http.ResponseWriter, r *http.Request) {
			// compute filesystem path
			rel := strings.TrimPrefix(r.URL.Path, cfg.App.Root)
			rel = strings.TrimPrefix(rel, "/")
			fsPath := filepath.Join("./web", rel)

			if _, err := os.Stat(fsPath); err == nil {
				// file exists, let the mounted static handler serve it by redirecting
				http.ServeFile(w, r, fsPath)
				return
			}

			// otherwise serve index.html so the SPA can handle routing
			http.ServeFile(w, r, "./web/index.html")
		})
	}

	port := int(cfg.App.Port)
	slog.Info("starting server on port", slog.Int("port", port))
	return &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: router,
	}, nil
}

func destroy() {
	serviceRegistry.StopAll()
	clientRepo.Close()
	dbUtils.Close()
}

func findById(ctx context.Context, id int) {
	slog.Info("findById...", slog.Int("id", id))

	if entity, err := clientService.FindById(ctx, id); err != nil {
		slog.Error("clientService.FindById", "error", err)
	} else {
		slog.Info("clientService.FindById", "entity", entity)
	}

	slog.Info("clientService.FindById", slog.Int("id", id))
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
