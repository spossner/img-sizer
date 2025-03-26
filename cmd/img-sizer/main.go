package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spossner/img-sizer/internal/config"
	"github.com/spossner/img-sizer/internal/handlers"
	"github.com/spossner/img-sizer/internal/storage"
	"github.com/spossner/img-sizer/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// Create loggger
	var logLevel slog.Level
	_ = logLevel.UnmarshalText([]byte(os.Getenv("LOG_LEVEL")))
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))

	projectRoot := utils.GetProjectRoot(logger)
	if projectRoot != "" {
		envFile := fmt.Sprintf("%s/.env", projectRoot)
		logger.Info("using .env file", "file", envFile)
		if err := godotenv.Load(envFile); err != nil {
			logger.Error("error loading .env file", "file", envFile, "error", err)
		}
	}

	// Load configuration
	cfg, err := config.Load(logger)
	if err != nil {
		logger.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize S3 client
	s3Client, err := storage.NewS3Client(context.Background())
	if err != nil {
		logger.Error("failed to initialize S3 client", "error", err)
		os.Exit(1)
	}

	// Create Fiber app with optimized config
	app := fiber.New(fiber.Config{
		AppName:                 "Image Sizer v2.0",
		EnableTrustedProxyCheck: true,
		ProxyHeader:             fiber.HeaderXForwardedFor,
		GETOnly:                 true,
		ServerHeader:            "Image-Sizer",
	})

	// Add middleware
	app.Use(recover.New()) // Recover from panics
	app.Use(fiberLogger.New(fiberLogger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} | ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))
	app.Use(limiter.New(limiter.Config{
		Max:        cfg.RateLimit.MaxRequests,
		Expiration: cfg.RateLimit.Window,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("X-Forwarded-For") + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "rate limit exceeded",
			})
		},
	}))

	// Add health routes
	app.Get("/ping", handlers.GetPingHandler())
	app.Get("/healthz", handlers.GetHealthHandler(cfg, s3Client))
	app.Get("/readyz", handlers.GetReadinessHandler(cfg, s3Client))

	// Add sizer routes
	app.Get("/v2/resize.jpg", handlers.GetCombinedHandler(cfg, s3Client))
	app.Get("/resize.jpg", handlers.GetResizeHandler(cfg, s3Client))
	app.Get("/crop.jpg", handlers.GetCropHandler(cfg, s3Client))

	// Create channel to listen for errors coming from the server
	serverErrors := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		port := utils.ParseInt(os.Getenv("PORT"), 8080)
		addr := fmt.Sprintf(":%d", port)
		logger.Info("starting server", "addr", addr, "rate_limit", cfg.RateLimit.MaxRequests, "window", cfg.RateLimit.Window)
		serverErrors <- app.Listen(addr)
	}()

	// Create channel to listen for an interrupt or terminate signal from the OS
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking main and waiting for shutdown
	select {
	case err := <-serverErrors:
		logger.Error("error starting server", "error", err)
		os.Exit(1)

	case <-shutdown:
		logger.Info("initiating shutdown...")

		// Create shutdown context with 10 second timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Shutdown the server
		if err := app.ShutdownWithContext(ctx); err != nil {
			logger.Error("error during server shutdown", "error", err)
			os.Exit(1)
		}

		logger.Info("server stopped")
	}
}
