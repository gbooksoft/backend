package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gbooksoft/backend/internal/database"
	"github.com/gbooksoft/backend/internal/domain/email/queue"
	"github.com/gbooksoft/backend/internal/domain/email/worker"
	"github.com/gbooksoft/backend/internal/routing"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/wneessen/go-mail"
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

	e.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("WEB_FRONT_FULL_DOMAIN")},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	db, err := database.SetUpConnection(ctx)

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
		if listenErr := srv.ListenAndServe(); listenErr != nil && !errors.Is(http.ErrServerClosed, listenErr) {
			log.Fatal().Err(listenErr).Msg("server shutdown")
		}
	}()

	client, err := mail.NewClient(
		"localhost",
		mail.WithPort(1025),
		mail.WithSMTPAuth(mail.SMTPAuthNoAuth),
		mail.WithTimeout(10*time.Second),
		mail.WithTLSPolicy(mail.NoTLS),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Ошибка создания SMTP клиента")
	}

	go func(db *sql.DB, client *mail.Client) {
		emailWorker := worker.NewEmailWorker(queue.NewEmailQueue(db), client)

		// Запускаем воркера
		emailWorker.StartSending()
	}(db, client)

	<-ctx.Done()

	stop()

	log.Info().Msg("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if shutdownErr := srv.Shutdown(ctx); shutdownErr != nil {
		log.Fatal().Err(shutdownErr).Msg("server forced to shutdown")
	}

	log.Info().Msg("Server exiting")
}
