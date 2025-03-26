package processing

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertColorEqual(t *testing.T, expected, actual color.Color, msgAndArgs ...interface{}) bool {
	// convert to RGBA
	expectedR, expectedG, expectedB, expectedA := expected.RGBA()
	actualR, actualG, actualB, actualA := actual.RGBA()

	// check if the colors are equal (using inDelta to allow for small differences)
	return assert.InDelta(t, float64(expectedR), float64(actualR), 1, msgAndArgs...) &&
		assert.InDelta(t, float64(expectedG), float64(actualG), 1, msgAndArgs...) &&
		assert.InDelta(t, float64(expectedB), float64(actualB), 1, msgAndArgs...) &&
		assert.InDelta(t, float64(expectedA), float64(actualA), 1, msgAndArgs...)
}

func TestFillBackground(t *testing.T) {
	// create test image with red triangle in the top right half
	original := createTestImage(100, 100, color.RGBA{R: 255, G: 0, B: 0, A: 255})

	tests := []struct {
		name          string
		bgColor       string
		expectedLeft  color.Color
		expectedRight color.Color
		shouldErr     bool
	}{
		{
			name:          "black background",
			bgColor:       Black,
			expectedLeft:  color.RGBA{R: 0, G: 0, B: 0, A: 0},
			expectedRight: color.RGBA{R: 255, G: 0, B: 0, A: 255},
		},
		{
			name:          "white background",
			bgColor:       "FFFFFF",
			expectedLeft:  color.RGBA{R: 255, G: 255, B: 255, A: 255},
			expectedRight: color.RGBA{R: 255, G: 0, B: 0, A: 255},
		},
		{
			name:          "red background",
			bgColor:       "FF0000",
			expectedLeft:  color.RGBA{R: 255, G: 0, B: 0, A: 255},
			expectedRight: color.RGBA{R: 255, G: 0, B: 0, A: 255},
		},
		{
			name:          "green background",
			bgColor:       "00FF00",
			expectedLeft:  color.RGBA{R: 0, G: 255, B: 0, A: 255},
			expectedRight: color.RGBA{R: 255, G: 0, B: 0, A: 255},
		},
		{
			name:          "blue background",
			bgColor:       "0000FF",
			expectedLeft:  color.RGBA{R: 0, G: 0, B: 255, A: 255},
			expectedRight: color.RGBA{R: 255, G: 0, B: 0, A: 255},
		},
		{
			name:      "invalid hex color",
			bgColor:   "invalid",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Fill background
			result, err := FillBackground(original, tt.bgColor)
			if tt.shouldErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Check dimensions
			assert.Equal(t, original.Bounds(), result.Bounds(), "image bounds should not change")

			// Check color at upper right area
			upperRightX := result.Bounds().Dx() * 2 / 3
			upperRightY := result.Bounds().Dy() / 3
			actualColor := result.At(upperRightX, upperRightY)
			assertColorEqual(t, actualColor, tt.expectedRight, "overlap color should match expected")

			// Check color at lower left area
			lowerLeftX := result.Bounds().Dx() / 3
			lowerLeftY := result.Bounds().Dy() * 2 / 3
			actualColor = result.At(lowerLeftX, lowerLeftY)
			assertColorEqual(t, actualColor, tt.expectedLeft, "background color should match expected")
		})
	}
}

func TestFillBackgroundWithTransparency(t *testing.T) {
	// Create a fully transparent image
	transparent := createTestImage(100, 100, color.RGBA{R: 0, G: 0, B: 0, A: 0})

	tests := []struct {
		name     string
		bgColor  string
		expected color.Color
	}{
		{
			name:     "transparent with white background",
			bgColor:  "FFFFFF",
			expected: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
		{
			name:     "transparent with red background",
			bgColor:  "FF0000",
			expected: color.RGBA{R: 255, G: 0, B: 0, A: 255},
		},
		{
			name:     "transparent with black background",
			bgColor:  Black,
			expected: color.RGBA{R: 0, G: 0, B: 0, A: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FillBackground(transparent, tt.bgColor)
			assert.NoError(t, err)

			// Check dimensions
			assert.Equal(t, transparent.Bounds(), result.Bounds(), "image bounds should not change")

			// Check color at upper right area
			upperRightX := result.Bounds().Dx() * 2 / 3
			upperRightY := result.Bounds().Dy() / 3
			actualColor := result.At(upperRightX, upperRightY)
			assertColorEqual(t, actualColor, tt.expected, "color in upper right area should match expected")

			// Check color at lower left area
			lowerLeftX := result.Bounds().Dx() / 3
			lowerLeftY := result.Bounds().Dy() * 2 / 3
			actualColor = result.At(lowerLeftX, lowerLeftY)
			assertColorEqual(t, actualColor, tt.expected, "color in lower left area should match expected")
		})
	}
}

func TestFillBackgroundWithPartialTransparency(t *testing.T) {
	imgColor := color.RGBA{R: 255, G: 0, B: 0, A: 255}

	// Create an image with red triangle in the top right half but a transparent center
	img := createTestImage(100, 100, imgColor)
	for y := 40; y < 60; y++ {
		for x := 40; x < 60; x++ {
			img.Set(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 0})
		}
	}

	tests := []struct {
		name     string
		bgColor  string
		expected struct {
			center color.Color
		}
	}{
		{
			name:    "partially transparent with white background",
			bgColor: "FFFFFF",
			expected: struct {
				center color.Color
			}{
				center: color.RGBA{R: 255, G: 255, B: 255, A: 255},
			},
		},
		{
			name:    "partially transparent with red background",
			bgColor: "FF0000",
			expected: struct {
				center color.Color
			}{
				center: color.RGBA{R: 255, G: 0, B: 0, A: 255},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FillBackground(img, tt.bgColor)
			assert.NoError(t, err)

			// Check dimensions
			assert.Equal(t, img.Bounds(), result.Bounds(), "image bounds should not change")

			// Check color at edge (should be original red)
			edgeColor := result.At(50, 10)
			assertColorEqual(t, edgeColor, imgColor, "top edge color should match expected")

			// Check color at center (should be background color)
			centerColor := result.At(50, 50)
			assertColorEqual(t, centerColor, tt.expected.center, "center color should match expected")

			// Check color at left edge (should be background color due to transparency of original test image)
			centerColor = result.At(10, 50)
			assertColorEqual(t, centerColor, tt.expected.center, "left edge color should match expected")
		})
	}
}
