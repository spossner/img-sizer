package handlers

import (
	"fmt"
	"image"
	"img-sizer/internal/config"

	"github.com/gofiber/fiber/v2"
)

type SizerParams struct {
	Width   int
	Height  int
	Quality int
	BgColor string
	Density float64
	Scale   float64
	Crop    image.Rectangle
}

func (p SizerParams) String() string {
	return fmt.Sprintf("%dx%d-q%d-bg%s-d%.2f-s%.2f-c%v", p.Width, p.Height, p.Quality, p.BgColor, p.Density, p.Scale, p.Crop)
}

type ParamsParser func(c *fiber.Ctx, cfg *config.Config) SizerParams
