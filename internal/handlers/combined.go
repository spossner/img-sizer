package handlers

import (
	"image"
	"img-sizer/internal/config"
	"img-sizer/internal/storage"

	"github.com/gofiber/fiber/v2"
)

func combinedParamsParser(c *fiber.Ctx, cfg *config.Config) SizerParams {
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

	// Get dimensions from query parameters
	cropWidth := c.QueryInt("crop[width]", 0)
	cropHeight := c.QueryInt("crop[height]", 0)
	cropX := c.QueryInt("crop[x]", 0)
	cropY := c.QueryInt("crop[y]", 0)

	// Get scale parameter for cropping
	scale := c.QueryFloat("crop[scale]", 1.0)
	if scale <= 0 {
		scale = 1.0
	}

	scaledWidth := int(float64(cropWidth) / scale)
	scaledHeight := int(float64(cropHeight) / scale)
	scaledX := int(float64(cropX) / scale)
	scaledY := int(float64(cropY) / scale)

	crop := image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight)

	return SizerParams{
		Width:   finalWidth,
		Height:  finalHeight,
		Quality: quality,
		BgColor: bgColor,
		Density: density,
		Scale:   scale,
		Crop:    crop,
	}
}

func GetCombinedHandler(cfg *config.Config, s3Client *storage.S3Client) fiber.Handler {
	return GetImageSizerHandler(cfg, s3Client, combinedParamsParser)
}
