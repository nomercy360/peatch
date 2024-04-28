package main

import (
	"context"
	"errors"
	"github.com/caarlos0/env/v11"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/peatch-io/peatch/docs"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/handler"
	storage "github.com/peatch-io/peatch/internal/s3"
	"github.com/peatch-io/peatch/internal/service"
	"github.com/peatch-io/peatch/internal/terrors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type config struct {
	DatabaseURL string `env:"DATABASE_URL,required"`
	Server      ServerConfig
	BotToken    string `env:"BOT_TOKEN,required"`
	AWS         AWSConfig
}

type ServerConfig struct {
	Port string `env:"SERVER_PORT" envDefault:"8080"`
	Host string `env:"SERVER_HOST" envDefault:"localhost"`
}

type AWSConfig struct {
	AccessKey string `env:"AWS_ACCESS_KEY_ID,required"`
	SecretKey string `env:"AWS_SECRET_ACCESS_KEY,required"`
	Bucket    string `env:"AWS_BUCKET,required"`
	Endpoint  string `env:"AWS_ENDPOINT,required"`
}

// @title Peatch API
// @version 1.0
// @description This is a sample server ClanPlatform server.

// @host localhost:8080
// @BasePath /
func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse config: %v\n", err)
	}

	pg, err := db.New(cfg.DatabaseURL)

	if err != nil {
		log.Fatalf("Failed to initialize database: %v\n", err)
	}

	defer pg.Close()

	e := echo.New()

	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		var (
			code = http.StatusInternalServerError
			msg  interface{}
		)

		var he *echo.HTTPError
		var terror *terrors.Error
		if errors.As(err, &he) {
			code = he.Code
			msg = he.Message
		} else if errors.As(err, &terror) {
			code = terror.Code
			msg = terror.Message
		} else {
			msg = err.Error()
		}

		if _, ok := msg.(string); ok {
			msg = map[string]interface{}{"error": msg}
		}

		if !c.Response().Committed {
			if c.Request().Method == http.MethodHead {
				err = c.NoContent(code)
			} else {
				err = c.JSON(code, msg)
			}

			if err != nil {
				e.Logger.Error(err)
			}
		}
	}

	s3Client, err := storage.NewS3Client(
		cfg.AWS.AccessKey, cfg.AWS.SecretKey, cfg.AWS.Endpoint, cfg.AWS.Bucket)

	if err != nil {
		log.Fatalf("Failed to initialize AWS S3 client: %v\n", err)
	}

	svc := service.New(pg, s3Client, service.Config{BotToken: cfg.BotToken})

	h := handler.New(svc)

	h.RegisterRoutes(e)

	//e.GET("/swagger/*", echoSwagger.WrapHandler)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := e.Start(cfg.Server.Host + ":" + cfg.Server.Port); err != nil {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
