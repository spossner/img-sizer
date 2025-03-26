package handlers

import (
	"img-sizer/internal/config"
	"img-sizer/internal/storage"

	"github.com/gofiber/fiber/v2"
)

func resizeParamsParser(c *fiber.Ctx, cfg *config.Config) SizerParams {
	// Get dimensions from query parameters
	width := c.QueryInt("width", 0)
	height := c.QueryInt("height", 0)
	quality := c.QueryInt("quality", cfg.Jpeg.Quality)
	bgColor := c.Query("background", cfg.Jpeg.Background)

	// Get density parameter
	density := c.QueryFloat("density", 1.0)
	density = c.QueryFloat("scale", density)
	if density <= 0 {
		density = 1.0
	}

	// Calculate final dimensions with density
	finalWidth := int(float64(width) * density)
	finalHeight := int(float64(height) * density)

	return SizerParams{
		Width:   finalWidth,
		Height:  finalHeight,
		Quality: quality,
		BgColor: bgColor,
		Density: density,
	}
}

func GetResizeHandler(cfg *config.Config, s3Client *storage.S3Client) fiber.Handler {
	return GetImageSizerHandler(cfg, s3Client, resizeParamsParser)
}
