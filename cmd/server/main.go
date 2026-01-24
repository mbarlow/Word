package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/mbarlow/word/internal/api/handler"
	"github.com/mbarlow/word/internal/config"
	"github.com/mbarlow/word/internal/repository"
)

func main() {
	cfg := config.Load()
	setupLogger(cfg)

	log.Info().
		Str("env", cfg.Env).
		Int("port", cfg.Port).
		Msg("starting word-service")

	db, err := repository.NewDB(cfg.DBPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(zerologMiddleware())
	e.Use(echoprometheus.NewMiddleware("word"))

	// Health and metrics
	e.GET("/health", healthHandler)
	e.GET("/metrics", echoprometheus.NewHandler())

	// API routes
	h := handler.New(db)
	v1 := e.Group("/v1")
	v1.GET("/works", h.ListWorks)
	v1.GET("/books", h.ListBooks)
	v1.GET("/text/:work/:book/:chapter", h.GetChapter)
	v1.GET("/verse/:work/:book/:chapter/:verse", h.GetVerse)
	v1.GET("/compare", h.Compare)
	v1.GET("/search", h.Search)
	v1.GET("/random-verse", h.GetRandomVerse)

	// Start server
	go func() {
		addr := cfg.Address()
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	log.Info().Str("addr", cfg.Address()).Msg("server started")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("server shutdown error")
	}
}

func setupLogger(cfg *config.Config) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	if cfg.LogFormat == "console" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

func zerologMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			req := c.Request()
			res := c.Response()

			log.Info().
				Str("method", req.Method).
				Str("uri", req.RequestURI).
				Int("status", res.Status).
				Dur("latency", time.Since(start)).
				Str("request_id", c.Response().Header().Get(echo.HeaderXRequestID)).
				Msg("request")

			return err
		}
	}
}

func healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
