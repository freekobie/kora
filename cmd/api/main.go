package main


import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/freekobie/kora/handler"
	"github.com/freekobie/kora/mail"
	"github.com/freekobie/kora/postgres"
	"github.com/freekobie/kora/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

//	@title			Backend for Kora file management
//	@version		1.0
//	@description	This is the backend API for the Kora file management application.
//	@contact.name	API Support
//	@contact.url	https://github.com/freekobie/kora
//	@contact.email	support@kora.local
func main() {

	_ = godotenv.Load()

	setupLogging()

	cfg := loadConfig()

	db, err := pgxpool.New(context.Background(), cfg.PostgresURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	mailer := mail.NewMailer(cfg.MailConfig)
	userService := service.NewUserService(postgres.NewUserStore(db), mailer)

	handler := handler.NewHandler(userService)

	app := newApplication(handler, cfg.ServerAddress)

	// Graceful shutdown setup
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	serverErr := make(chan error, 1)
	go func() {
		slog.Info("Starting server")
		serverErr <- app.start()
	}()

	select {
	case err := <-serverErr:
		if err != nil {
			panic(err)
		}
	case sig := <-stop:
		slog.Info("Shutting down server", "signal", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := app.shutdown(ctx); err != nil {
			slog.Error("Graceful shutdown failed", "error", err)
		}
	}
}