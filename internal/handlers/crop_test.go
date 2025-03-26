package handlers

import (
	"image"
	"testing"

	"github.com/spossner/img-sizer/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func TestCropParamsParser(t *testing.T) {
	// Create test config
	cfg := &config.Config{
		MaxInputDimension:  5000,
		MaxOutputDimension: 2000,
		AllowedDimensions: []config.Dimension{
			{Width: 100, Height: 100},
			{Width: 200, Height: 200},
			{Width: 300, Height: 300},
			{Width: 400, Height: 300},
			{Width: 800, Height: 600},
			{Width: 1024, Height: 768},
		},
		AllowAllDimensions: false,
		Jpeg: config.Jpeg{
			Quality:    70,
			Background: "000000",
		},
	}

	tests := []struct {
		name           string
		query          string
		expectedParams SizerParams
	}{
		{
			name:  "basic crop parameters",
			query: "width=100&height=100&x=10&y=10",
			expectedParams: SizerParams{
				Width:   100,
				Height:  100,
				Quality: 70,
				BgColor: "000000",
				Density: 1.0,
				Scale:   1.0,
				Crop:    image.Rect(10, 10, 110, 110),
			},
		},
		{
			name:  "with scale parameter",
			query: "width=100&height=100&x=10&y=10&scale=0.5",
			expectedParams: SizerParams{
				Width:   100,
				Height:  100,
				Quality: 70,
				BgColor: "000000",
				Density: 1.0,
				Scale:   0.5,
				Crop:    image.Rect(20, 20, 220, 220),
			},
		},
		{
			name:  "with density parameter",
			query: "width=100&height=100&x=10&y=10&density=2.0",
			expectedParams: SizerParams{
				Width:   200,
				Height:  200,
				Quality: 70,
				BgColor: "000000",
				Density: 2.0,
				Scale:   1.0,
				Crop:    image.Rect(10, 10, 110, 110),
			},
		},
		{
			name:  "with quality parameter",
			query: "width=100&height=100&x=10&y=10&quality=85",
			expectedParams: SizerParams{
				Width:   100,
				Height:  100,
				Quality: 85,
				BgColor: "000000",
				Density: 1.0,
				Scale:   1.0,
				Crop:    image.Rect(10, 10, 110, 110),
			},
		},
		{
			name:  "with background color",
			query: "width=100&height=100&x=10&y=10&background=FF0000",
			expectedParams: SizerParams{
				Width:   100,
				Height:  100,
				Quality: 70,
				BgColor: "FF0000",
				Density: 1.0,
				Scale:   1.0,
				Crop:    image.Rect(10, 10, 110, 110),
			},
		},
		{
			name:  "with all parameters",
			query: "width=100&height=100&x=10&y=10&scale=0.5&density=2.0&quality=85&background=FF0000",
			expectedParams: SizerParams{
				Width:   200,
				Height:  200,
				Quality: 85,
				BgColor: "FF0000",
				Density: 2.0,
				Scale:   0.5,
				Crop:    image.Rect(20, 20, 220, 220),
			},
		},
		{
			name:  "zero dimensions",
			query: "width=0&height=0&x=0&y=0",
			expectedParams: SizerParams{
				Width:   0,
				Height:  0,
				Quality: 70,
				BgColor: "000000",
				Density: 1.0,
				Scale:   1.0,
				Crop:    image.Rect(0, 0, 0, 0),
			},
		},
		{
			name:  "missing dimensions",
			query: "x=10&y=10",
			expectedParams: SizerParams{
				Width:   0,
				Height:  0,
				Quality: 70,
				BgColor: "000000",
				Density: 1.0,
				Scale:   1.0,
				Crop:    image.Rect(10, 10, 10, 10),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new Fiber app
			app := fiber.New()

			// Create a test context
			ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
			ctx.Request().SetRequestURI("?" + tt.query)

			// Parse parameters
			params := cropParamsParser(ctx, cfg)

			// Check if parameters match expected values
			if params.Width != tt.expectedParams.Width {
				t.Errorf("Width = %v, want %v", params.Width, tt.expectedParams.Width)
			}
			if params.Height != tt.expectedParams.Height {
				t.Errorf("Height = %v, want %v", params.Height, tt.expectedParams.Height)
			}
			if params.Quality != tt.expectedParams.Quality {
				t.Errorf("Quality = %v, want %v", params.Quality, tt.expectedParams.Quality)
			}
			if params.BgColor != tt.expectedParams.BgColor {
				t.Errorf("BgColor = %v, want %v", params.BgColor, tt.expectedParams.BgColor)
			}
			if params.Density != tt.expectedParams.Density {
				t.Errorf("Density = %v, want %v", params.Density, tt.expectedParams.Density)
			}
			if params.Scale != tt.expectedParams.Scale {
				t.Errorf("Scale = %v, want %v", params.Scale, tt.expectedParams.Scale)
			}
			if params.Crop != tt.expectedParams.Crop {
				t.Errorf("Crop = %v, want %v", params.Crop, tt.expectedParams.Crop)
			}

			// Release the context
			app.ReleaseCtx(ctx)
		})
	}
}

func TestCropParamsParserWithAllowAllDimensions(t *testing.T) {
	// Create test config with AllowAllDimensions enabled
	cfg := &config.Config{
		MaxInputDimension:  5000,
		MaxOutputDimension: 2000,
		AllowAllDimensions: true,
		Jpeg: config.Jpeg{
			Quality:    70,
			Background: "000000",
		},
	}

	tests := []struct {
		name           string
		query          string
		expectedParams SizerParams
	}{
		{
			name:  "any dimension allowed",
			query: "width=150&height=250&x=10&y=10",
			expectedParams: SizerParams{
				Width:   150,
				Height:  250,
				Quality: 70,
				BgColor: "000000",
				Density: 1.0,
				Scale:   1.0,
				Crop:    image.Rect(10, 10, 160, 260),
			},
		},
		{
			name:  "dimensions within max limits",
			query: "width=3000&height=2000&x=100&y=100",
			expectedParams: SizerParams{
				Width:   3000,
				Height:  2000,
				Quality: 70,
				BgColor: "000000",
				Density: 1.0,
				Scale:   1.0,
				Crop:    image.Rect(100, 100, 3100, 2100),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
			ctx.Request().SetRequestURI("?" + tt.query)

			params := cropParamsParser(ctx, cfg)

			if params.Width != tt.expectedParams.Width {
				t.Errorf("Width = %v, want %v", params.Width, tt.expectedParams.Width)
			}
			if params.Height != tt.expectedParams.Height {
				t.Errorf("Height = %v, want %v", params.Height, tt.expectedParams.Height)
			}
			if params.Quality != tt.expectedParams.Quality {
				t.Errorf("Quality = %v, want %v", params.Quality, tt.expectedParams.Quality)
			}
			if params.BgColor != tt.expectedParams.BgColor {
				t.Errorf("BgColor = %v, want %v", params.BgColor, tt.expectedParams.BgColor)
			}
			if params.Density != tt.expectedParams.Density {
				t.Errorf("Density = %v, want %v", params.Density, tt.expectedParams.Density)
			}
			if params.Scale != tt.expectedParams.Scale {
				t.Errorf("Scale = %v, want %v", params.Scale, tt.expectedParams.Scale)
			}
			if params.Crop != tt.expectedParams.Crop {
				t.Errorf("Crop = %v, want %v", params.Crop, tt.expectedParams.Crop)
			}

			app.ReleaseCtx(ctx)
		})
	}
}
