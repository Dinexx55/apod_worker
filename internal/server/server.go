package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"nasa-apod-app/internal/config"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	server *http.Server
}

func NewServer(config config.ServerConfig, handler http.Handler) *Server {
	return &Server{
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", config.Host, config.Port),
			Handler: handler,
		},
	}
}

func (s *Server) Run() error {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully.")
	return nil
}
