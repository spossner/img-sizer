package handlers

import (
	"image"
	"img-sizer/internal/config"
	"img-sizer/internal/storage"

	"github.com/gofiber/fiber/v2"
)

func cropParamsParser(c *fiber.Ctx, cfg *config.Config) SizerParams {
	// Get dimensions from query parameters
	w := c.QueryInt("width", 0)
	h := c.QueryInt("height", 0)
	x := c.QueryInt("x", 0)
	y := c.QueryInt("y", 0)

	quality := c.QueryInt("quality", cfg.Jpeg.Quality)
	bgColor := c.Query("background", cfg.Jpeg.Background)

	// Get density parameter for final dimensions
	density := c.QueryFloat("density", 1.0)
	if density <= 0 {
		density = 1.0
	}

	width := int(float64(w) * density)
	height := int(float64(h) * density)

	// Get scale parameter for cropping
	scale := c.QueryFloat("scale", 1.0)
	if scale <= 0 {
		scale = 1.0
	}

	scaledWidth := int(float64(w) / scale)
	scaledHeight := int(float64(h) / scale)
	scaledX := int(float64(x) / scale)
	scaledY := int(float64(y) / scale)

	crop := image.Rect(scaledX, scaledY, scaledX+scaledWidth, scaledY+scaledHeight)

	return SizerParams{
		Width:   width,
		Height:  height,
		Quality: quality,
		BgColor: bgColor,
		Density: density,
		Scale:   scale,
		Crop:    crop,
	}
}

func GetCropHandler(cfg *config.Config, s3Client *storage.S3Client) fiber.Handler {
	return GetImageSizerHandler(cfg, s3Client, cropParamsParser)
}
