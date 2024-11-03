package main

import (
	"context"
	"github.com/gbooksoft/backend/internal/database"
	"github.com/gbooksoft/backend/internal/routing"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	e := gin.Default()

	_, err = database.SetUpConnection(ctx)

	if err != nil {
		log.Fatal().Err(err).Msg("can not connect to database")
	}

	routing.Add(e)

	srv := &http.Server{
		Addr:    "localhost:8080",
		Handler: e,
	}

	go func() {
		log.Debug().Msg("Server started: http://localhost:8080")
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server shutdown")
		}
	}()

	<-ctx.Done()

	stop()

	log.Info().Msg("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("server forced to shutdown")
	}

	log.Info().Msg("Server exiting")
}
