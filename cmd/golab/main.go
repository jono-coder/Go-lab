package main

import (
	"Go-lab/config"
	"Go-lab/internal/client"
	"Go-lab/internal/contact"
	myMiddleware "Go-lab/internal/middleware"
	"Go-lab/internal/security"
	"Go-lab/internal/utils"
	"Go-lab/internal/utils/httpconst"
	"log/slog"
	"strconv"

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
		slog.Error("failed to start application: %v", err)
		panic(err)
	}

	// listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if err != nil {
				slog.Error("Error", err)
				cancel()
			}
		}
	}()

	go func() {
		<-quit
		slog.Info("Signal received, shutting down...")
		cancel() // unblock context-aware goroutines
	}()

	// block until context is done
	<-ctx.Done()

	defer destroy()
}

func initialise(ctx context.Context) (*http.Server, error) {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config: %v", err)
		panic(err)
	}
	slog.Info("Running in env: %s", "env", cfg.App.Env)

	dbUtils = utils.NewDbUtils(&cfg.DB)
	err = dbUtils.Init()
	if err != nil {
		slog.Error("Error", "err", err)
		panic(err)
	}

	oauthConfig := security.NewOAuthConfig(ctx, cfg.App.BaseUrl)

	serviceRegistry = utils.NewServiceRegistry()
	////////// plumbing //////////

	////////// client //////////
	clientRepo = client.NewRepo(dbUtils)
	clientApi, err := client.NewAPI(oauthConfig)
	if err != nil {
		slog.Error("Error", "err", err)
		panic(err)
	}
	clientService = client.NewService(dbUtils, clientRepo, clientApi)
	serviceRegistry.Register(clientService)
	////////// client //////////

	////////// contact //////////
	contactRepo := contact.NewRepo(dbUtils)
	contactService := contact.NewService(dbUtils, contactRepo)
	serviceRegistry.Register(contactService)
	////////// contact //////////

	serviceRegistry.StartAll()

	////////// router //////////
	router := chi.NewRouter()
	router.Use(middleware.Compress(5, httpconst.CONTENT_TYPE_JSON, httpconst.CONTENT_TYPE_XML,
		httpconst.CONTENT_TYPE_HTML, httpconst.CONTENT_TYPE_TEXT))
	router.Use(middleware.Logger)
	router.Use(myMiddleware.SecureHandler)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Throttle(int(cfg.App.Throttle)))
	router.Use(middleware.Timeout(cfg.App.TimeoutInSeconds))

	router.Get(cfg.App.Root, func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Go is listening"))
		if err != nil {
			slog.Error("Error", "err", err)
		}
	})

	clientHandler := client.NewHandler(ctx, clientService, cfg.App)
	router.Route(cfg.App.Root+"/client", func(r chi.Router) {
		r.With(utils.Medium).Get("/", clientHandler.List)
		r.With(utils.Low).Get("/{id}", clientHandler.Get)
		r.With(middleware.NoCache).Post("/", clientHandler.Create)
		r.With(middleware.NoCache).Put("/{id}", clientHandler.Update)
		r.With(middleware.NoCache).Delete("/{id}", clientHandler.Delete)
	})

	oauthHandler := security.NewHandler(ctx, cfg.App)
	router.Route(cfg.App.Root+"/security", func(r chi.Router) {
		r.Post("/oauth/token", oauthHandler.Auth)
	})
	////////// router //////////

	return &http.Server{
		Addr:    ":" + strconv.FormatInt(int64(cfg.App.Port), 10),
		Handler: router,
	}, nil
}

func destroy() {
	serviceRegistry.StopAll()
	clientRepo.Close()
	dbUtils.Close()
}

func findById(ctx context.Context, id int) {
	slog.Info("findById(%d)...", "id", id)

	entity, err := clientService.FindById(ctx, id)
	if err != nil {
		slog.Error("Error", "err", err)
	} else {
		slog.Info("Entity", "entity", entity)
	}

	slog.Info("findById(%d).", "id", id)
}

func doBusinessStuff(ctx context.Context) {
	slog.Info("doBusinessStuff()...")

	err := clientService.DoBusinessStuff(ctx)
	if err != nil {
		slog.Error("Error", "err", err)
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
						slog.Error("Task exiting early due to cancellation")
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
			slog.Error("Submit aborted: ", "err", err)
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

	c, err := clientService.GetById(1)
	if err != nil {
		slog.Error("Error", "err", err)
	} else {
		slog.Info("clients retrieved", "clients", c)
	}

	slog.Info("test1Api()!")
}

func test2Api() {
	c, err := clientService.GetAll()
	if err != nil {
		slog.Error(err.Error(), "err", err)
	} else {
		slog.Info("clients retrieved", "clients", c)
	}
}
