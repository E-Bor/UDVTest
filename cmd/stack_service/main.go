package main

import (
	"StackService/internal/config"
	"StackService/internal/http_server"
	"StackService/internal/http_server/handlers"
	"StackService/internal/services"
	"StackService/internal/storage/postgresql"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	logger := initLogger(cfg.Env)
	logger.Info("starting service")

	db := postgresql.InitDB(logger, cfg)

	stack := services.Stack{DB: db}
	err := stack.LoadStackIfExist()

	if err != nil {
		logger.Error("error getting stored messages", err)
	}

	stackHandler := &handlers.StackHandler{Stack: &stack, Logger: logger}

	router := http_server.NewStackServiceRouter(stackHandler)

	err = http.ListenAndServe(cfg.HTTPServer.Address, router)

	if err != nil {
		logger.Error("error starting http server", err)
	}
}

func initLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)
	}
	return log
}
