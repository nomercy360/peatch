package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	_ "github.com/peatch-io/peatch/docs"
	"github.com/peatch-io/peatch/internal/config"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/handler"
	"github.com/peatch-io/peatch/internal/middleware"
	"github.com/peatch-io/peatch/internal/s3"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func gracefulShutdown(e *echo.Echo, logr *slog.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logr.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logr.Error("Error during server shutdown", "error", err)
	}
	logr.Info("Server gracefully stopped")
}

// @title Peatch API
// @version 1.0
// @description API Documentation for the Api Dating Project

// @host api.peatch.io
// @schemes https
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	storage, err := db.ConnectDB(cfg.DBURL, cfg.DBName)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	logr := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	e := echo.New()

	e.HideBanner = true
	e.HidePort = true
	e.IPExtractor = func(req *http.Request) string {
		if cfIP := req.Header.Get("CF-Connecting-IP"); cfIP != "" {
			return cfIP
		}
		return echo.ExtractIPFromXFFHeader()(req)
	}

	middleware.Setup(e, logr)

	hConfig := handler.Config{
		JWTSecret:        cfg.JWTSecret,
		TelegramBotToken: cfg.TelegramBotToken,
		AssetsURL:        cfg.AssetsURL,
	}

	s3Client, err := s3.NewClient(
		cfg.AWSConfig.AccessKey,
		cfg.AWSConfig.SecretKey,
		cfg.AWSConfig.Endpoint,
		cfg.AWSConfig.Bucket,
	)

	h := handler.New(storage, hConfig, s3Client, logr)

	h.SetupRoutes(e)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	go gracefulShutdown(e, logr)

	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	logr.Info("Starting server", "address", address)
	if err := e.Start(address); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logr.Error("Error starting server", "error", err)
	}
}
