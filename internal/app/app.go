package app

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"nasa-apod-app/internal/config"
	"nasa-apod-app/internal/handler"
	"nasa-apod-app/internal/migration"
	"nasa-apod-app/internal/repository"
	"nasa-apod-app/internal/server"
	"nasa-apod-app/internal/service"
	"net/http"
)

type App struct {
	server *server.Server
}

func NewApp(config *config.Config, logger *zap.Logger) (*App, error) {
	migrator := migration.NewMigration()

	apodImagesRepository, err := repository.InitDB(config.DatabaseConfig, migrator, logger)
	if err != nil {
		logger.With(
			zap.String("place", "main"),
			zap.Error(err),
		).Panic("Failed to establish database connection")
	}

	apodImagesService := service.NewApodImagesService(logger, apodImagesRepository)
	apodWorker := service.NewAPODWorker(apodImagesService, config.NasaApiKey, config.WorkerConfig, logger)
	apodImagesHandler := handler.NewApodImagesHandler(apodImagesService, logger)

	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodOptions, http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodPatch},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Set-Cookie", "User-Agent", "Origin"},
		AllowOriginFunc: func(origin string) bool {
			return true
		},
	})

	mux := mux.NewRouter()
	apodImagesHandler.Init(mux)
	handler := c.Handler(mux)

	httpServer := server.NewServer(config.ServerConfig, handler)

	go apodWorker.Start()

	return &App{
		server: httpServer,
	}, nil
}

func (app *App) Run() error {
	return app.server.Run()
}
