package handlers

import (
	"fmt"
	"img-sizer/internal/config"
	"img-sizer/internal/storage"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ServerStatus struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	S3Status  string `json:"s3_status"`
}

func getStatus(s3Client *storage.S3Client, c *fiber.Ctx, ok, notOk string) (ServerStatus, error) {
	status := ServerStatus{
		Status:    ok,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	// Check S3 connectivity
	err := s3Client.CheckHealth(c.Context())
	if err != nil {
		status.Status = notOk
		status.S3Status = fmt.Sprintf("error: %v", err)
	} else {
		status.S3Status = "connected"
	}

	return status, err
}

func GetPingHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	}
}

func GetHealthHandler(cfg *config.Config, s3Client *storage.S3Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		status, _ := getStatus(s3Client, c, "healthy", "unhealthy")
		return c.JSON(status)
	}
}

func GetReadinessHandler(cfg *config.Config, s3Client *storage.S3Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		status, err := getStatus(s3Client, c, "ready", "not_ready")
		if err != nil {
			c.Status(fiber.StatusServiceUnavailable)
		}
		return c.JSON(status)
	}
}
