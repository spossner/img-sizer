package handlers

import (
	"image"
	"testing"

	"github.com/spossner/img-sizer/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestResizeParamsParser(t *testing.T) {
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
			name:  "basic resize parameters",
			query: "width=100&height=100",
			expectedParams: SizerParams{
				Width:   100,
				Height:  100,
				Quality: 70,
				BgColor: "000000",
				Density: 1.0,
				Crop:    image.Rectangle{},
			},
		},
		{
			name:  "with scale parameter",
			query: "width=100&height=100&scale=0.5",
			expectedParams: SizerParams{
				Width:   50,
				Height:  50,
				Quality: 70,
				BgColor: "000000",
				Density: 0.5,
				Crop:    image.Rectangle{},
			},
		},
		{
			name:  "with density parameter",
			query: "width=100&height=100&density=2.0",
			expectedParams: SizerParams{
				Width:   200,
				Height:  200,
				Quality: 70,
				BgColor: "000000",
				Density: 2.0,
				Crop:    image.Rectangle{},
			},
		},
		{
			name:  "with quality parameter",
			query: "width=100&height=100&quality=85",
			expectedParams: SizerParams{
				Width:   100,
				Height:  100,
				Quality: 85,
				BgColor: "000000",
				Density: 1.0,
				Crop:    image.Rectangle{},
			},
		},
		{
			name:  "with background color",
			query: "width=100&height=100&background=FF0000",
			expectedParams: SizerParams{
				Width:   100,
				Height:  100,
				Quality: 70,
				BgColor: "FF0000",
				Density: 1.0,
				Crop:    image.Rectangle{},
			},
		},
		{
			name:  "with all parameters",
			query: "width=100&height=100&scale=0.5&density=2.0&quality=85&background=FF0000",
			expectedParams: SizerParams{
				Width:   50,
				Height:  50,
				Quality: 85,
				BgColor: "FF0000",
				Density: 0.5,
				Crop:    image.Rectangle{},
			},
		},
		{
			name:  "zero dimensions",
			query: "width=0&height=0",
			expectedParams: SizerParams{
				Width:   0,
				Height:  0,
				Quality: 70,
				BgColor: "000000",
				Density: 1.0,
				Crop:    image.Rectangle{},
			},
		},
		{
			name:  "missing dimensions",
			query: "",
			expectedParams: SizerParams{
				Width:   0,
				Height:  0,
				Quality: 70,
				BgColor: "000000",
				Density: 1.0,
				Crop:    image.Rectangle{},
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
			params := resizeParamsParser(ctx, cfg)

			// Check if parameters match expected values
			assert.Equal(t, tt.expectedParams.Width, params.Width, "Width mismatch")
			assert.Equal(t, tt.expectedParams.Height, params.Height, "Height mismatch")
			assert.Equal(t, tt.expectedParams.Quality, params.Quality, "Quality mismatch")
			assert.Equal(t, tt.expectedParams.BgColor, params.BgColor, "BgColor mismatch")
			assert.Equal(t, tt.expectedParams.Density, params.Density, "Density mismatch")
			assert.Equal(t, tt.expectedParams.Scale, params.Scale, "Scale mismatch")
			assert.Equal(t, tt.expectedParams.Crop, params.Crop, "Crop mismatch")

			// Release the context
			app.ReleaseCtx(ctx)
		})
	}
}

func TestResizeParamsParserWithAllowAllDimensions(t *testing.T) {
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
				Crop:    image.Rectangle{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
			ctx.Request().SetRequestURI("?" + tt.query)

			params := resizeParamsParser(ctx, cfg)

			assert.Equal(t, tt.expectedParams.Width, params.Width, "Width mismatch")
			assert.Equal(t, tt.expectedParams.Height, params.Height, "Height mismatch")
			assert.Equal(t, tt.expectedParams.Quality, params.Quality, "Quality mismatch")
			assert.Equal(t, tt.expectedParams.BgColor, params.BgColor, "BgColor mismatch")
			assert.Equal(t, tt.expectedParams.Density, params.Density, "Density mismatch")
			assert.Equal(t, tt.expectedParams.Scale, params.Scale, "Scale mismatch")
			assert.Equal(t, tt.expectedParams.Crop, params.Crop, "Crop mismatch")

			app.ReleaseCtx(ctx)
		})
	}
}
