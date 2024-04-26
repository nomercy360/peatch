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
	DBConnString string `env:"DB_CONN_STRING,required"`
	Server       ServerConfig
	BotToken     string `env:"BOT_TOKEN,required"`
}

type ServerConfig struct {
	Port string `env:"PORT" envDefault:"8080"`
	Host string `env:"HOST" envDefault:"localhost"`
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

	pg, err := db.New(cfg.DBConnString)

	if err != nil {
		log.Fatalf("Failed to initialize database: %v\n", err)
	}

	defer pg.Close()

	e := echo.New()

	e.Use(middleware.Recover())

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

	svc := service.New(pg, service.Config{BotToken: cfg.BotToken})

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
