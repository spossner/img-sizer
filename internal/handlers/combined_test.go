package handlers

import (
	"image"
	"testing"

	"img-sizer/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func TestCombinedParamsParser(t *testing.T) {
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
			name:  "valid dimensions",
			query: "width=800&height=600",
			expectedParams: SizerParams{
				Width:   800,
				Height:  600,
				Quality: 70,
				BgColor: "000000",
				Density: 1.0,
				Scale:   1.0,
				Crop:    image.Rectangle{},
			},
		},
		{
			name:  "with crop parameters",
			query: "width=800&height=600&crop[x]=100&crop[y]=100&crop[width]=400&crop[height]=300&crop[scale]=1.0",
			expectedParams: SizerParams{
				Width:   800,
				Height:  600,
				Quality: 70,
				BgColor: "000000",
				Density: 1.0,
				Scale:   1.0,
				Crop:    image.Rect(100, 100, 500, 400),
			},
		},
		{
			name:  "with quality parameter",
			query: "width=800&height=600&quality=85",
			expectedParams: SizerParams{
				Width:   800,
				Height:  600,
				Quality: 85,
				BgColor: "000000",
				Density: 1.0,
				Scale:   1.0,
				Crop:    image.Rectangle{},
			},
		},
		{
			name:  "with background color",
			query: "width=800&height=600&background=FF0000",
			expectedParams: SizerParams{
				Width:   800,
				Height:  600,
				Quality: 70,
				BgColor: "FF0000",
				Density: 1.0,
				Scale:   1.0,
				Crop:    image.Rectangle{},
			},
		},
		{
			name:  "with density parameter",
			query: "width=800&height=600&density=2.0",
			expectedParams: SizerParams{
				Width:   1600,
				Height:  1200,
				Quality: 70,
				BgColor: "000000",
				Density: 2.0,
				Scale:   1.0,
				Crop:    image.Rectangle{},
			},
		},
		{
			name:  "with regular scale parameter",
			query: "width=800&height=600&scale=2.0",
			expectedParams: SizerParams{
				Width:   1600,
				Height:  1200,
				Quality: 70,
				BgColor: "000000",
				Density: 2.0,
				Scale:   1.0,
				Crop:    image.Rectangle{},
			},
		},
		{
			name:  "with all parameters",
			query: "width=800&height=600&crop[x]=100&crop[y]=100&crop[width]=400&crop[height]=300&crop[scale]=1.0&quality=85&background=FF0000&density=2.0",
			expectedParams: SizerParams{
				Width:   1600,
				Height:  1200,
				Quality: 85,
				BgColor: "FF0000",
				Density: 2.0,
				Scale:   1.0,
				Crop:    image.Rect(100, 100, 500, 400),
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
			params := combinedParamsParser(ctx, cfg)

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

func TestCombinedParamsParserWithAllowAllDimensions(t *testing.T) {
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
			query: "width=150&height=250",
			expectedParams: SizerParams{
				Width:   150,
				Height:  250,
				Quality: 70,
				BgColor: "000000",
				Density: 1.0,
				Scale:   1.0,
				Crop:    image.Rectangle{},
			},
		},
		{
			name:  "dimensions within max limits",
			query: "width=3000&height=2000",
			expectedParams: SizerParams{
				Width:   3000,
				Height:  2000,
				Quality: 70,
				BgColor: "000000",
				Density: 1.0,
				Scale:   1.0,
				Crop:    image.Rectangle{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
			ctx.Request().SetRequestURI("?" + tt.query)

			params := combinedParamsParser(ctx, cfg)

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
