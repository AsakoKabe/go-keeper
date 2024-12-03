package rest

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	slogchi "github.com/samber/slog-chi"
	"go-keeper/config/server"
	"go-keeper/internal/db/connection"
	"go-keeper/internal/db/repository"
	"go-keeper/internal/http/errs"
	"go-keeper/internal/http/rest/handlers"
	"go-keeper/internal/keeper/data"
	"go-keeper/internal/keeper/user"
	"go-keeper/pkg/hashing"
	"go-keeper/pkg/middleware"
	"golang.org/x/crypto/acme/autocert"

	"github.com/go-chi/chi/v5"
)

// App Приложение
type App struct {
	httpServer   *http.Server
	dbPool       *sql.DB
	repositories *repository.Repositories
	userService  *user.Service
	dataService  *data.Service
}

// NewApp Конструктор для App
func NewApp(cfg *server.Config) (*App, error) {
	if cfg.DatabaseDSN == "" {
		return &App{}, nil
	}
	pool, err := connection.NewDBPool(cfg.DatabaseDSN)
	if err != nil {
		slog.Error("error to create db pool", slog.String("err", err.Error()))
		return nil, errs.ErrCreateDBPoll
	}

	pgRepositories, err := repository.NewPostgresRepositories(pool)
	if err != nil {
		slog.Error("error to create pg repository", slog.String("err", err.Error()))
		return nil, errs.ErrCreateServices
	}

	return &App{
		dbPool:       pool,
		repositories: pgRepositories,
		userService:  user.NewService(pgRepositories.UserRepository),
		dataService: data.NewService(
			pgRepositories.DataRepository, hashing.NewCrypter(cfg.SecretHash),
		),
	}, nil
}

// Run Запуск приложения
func (a *App) Run(cfg *server.Config) error {

	router := chi.NewRouter()
	router.Use(slogchi.New(slog.New(slog.NewTextHandler(os.Stdout, nil))))
	var err error

	a.registerHandlers(router)

	if err != nil {
		return errs.ErrRegisterEndpoints
	}

	manager := &autocert.Manager{
		// директория для хранения сертификатов
		Cache: autocert.DirCache("cache-dir"),
		// функция, принимающая Terms of Service издателя сертификатов
		Prompt: autocert.AcceptTOS,
		// перечень доменов, для которых будут поддерживаться сертификаты
		HostPolicy: autocert.HostWhitelist("mysite.ru", "www.mysite.ru"),
	}

	a.httpServer = &http.Server{
		Addr:           cfg.Addr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		TLSConfig:      manager.TLSConfig(),
	}

	go func() {
		if cfg.EnableHTTPS {
			slog.Info("run server with HTTPS")
			err = http.ListenAndServeTLS(
				cfg.Addr,
				cfg.CertFile,
				cfg.KeyFile,
				router,
			)
		} else {
			slog.Info("run server with HTTP")
			err = http.ListenAndServe(
				cfg.Addr,
				router,
			)
		}
		if err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()
	go func() {
		<-ctx.Done()
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			slog.Error("forcing exit")
		}
	}()
	if err = a.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func (a *App) registerHandlers(router *chi.Mux) {
	pingHandler := handlers.NewPingHandler(a.repositories.PingRepository)
	router.Get("/ping", pingHandler.HealthDB)

	userHandler := handlers.NewUserHandler(a.userService)
	router.Post("/user/register", userHandler.Register)
	router.Post("/user/auth", userHandler.Authorize)

	dataHandler := handlers.NewDataHandler(a.dataService)
	router.Route(
		"/user/data", func(r chi.Router) {
			r.Use(middleware.Auth)
			r.Post("/{type}", dataHandler.Add)
			r.Get("/{type}", dataHandler.GetAllData)
			r.Get("/{type}/{dataID}", dataHandler.GeByID)
			r.Delete("/{type}/{dataID}", dataHandler.DeleteByID)
			r.Put("/{type}/{dataID}", dataHandler.Update)
		},
	)

}

// Stop Завершение работы приложения
func (a *App) Stop() {
	slog.Info("goroutines stopped")

	if a.dbPool == nil {
		return
	}
	err := a.dbPool.Close()
	if err != nil {
		return
	}
	slog.Info("db connection closed")
}
