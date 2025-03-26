package processing

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

// createTestImage creates a test image with the specified dimensions and color
// the image will have top right triangle filled with the color and the lower left triangle unfilled (transparent)
func createTestImage(width, height int, c color.Color) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := y; x < width; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}

func TestResizeImage(t *testing.T) {
	// Create a test image with red triangle in the top right half
	original := createTestImage(100, 100, color.RGBA{R: 255, G: 0, B: 0, A: 255})

	tests := []struct {
		name           string
		width          int
		height         int
		expectedWidth  int
		expectedHeight int
	}{
		{
			name:           "resize to smaller dimensions",
			width:          50,
			height:         50,
			expectedWidth:  50,
			expectedHeight: 50,
		},
		{
			name:           "resize to larger dimensions",
			width:          200,
			height:         200,
			expectedWidth:  200,
			expectedHeight: 200,
		},
		{
			name:           "resize width only",
			width:          150,
			height:         0,
			expectedWidth:  150,
			expectedHeight: 150,
		},
		{
			name:           "resize height only",
			width:          0,
			height:         150,
			expectedWidth:  150,
			expectedHeight: 150,
		},
		{
			name:           "no resize",
			width:          0,
			height:         0,
			expectedWidth:  100,
			expectedHeight: 100,
		},
		{
			name:           "maintain aspect ratio",
			width:          200,
			height:         0,
			expectedWidth:  200,
			expectedHeight: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Resize the image
			resized := ResizeImage(original, tt.width, tt.height)

			// Check dimensions
			bounds := resized.Bounds()
			assert.Equal(t, tt.expectedWidth, bounds.Dx(), "width should match expected")
			assert.Equal(t, tt.expectedHeight, bounds.Dy(), "height should match expected")
		})
	}
}

func TestResizeImageWithDifferentAspectRatios(t *testing.T) {
	tests := []struct {
		name           string
		originalWidth  int
		originalHeight int
		targetWidth    int
		targetHeight   int
		expectedWidth  int
		expectedHeight int
	}{
		{
			name:           "landscape to portrait",
			originalWidth:  200,
			originalHeight: 100,
			targetWidth:    100,
			targetHeight:   200,
			expectedWidth:  100,
			expectedHeight: 200,
		},
		{
			name:           "portrait to landscape",
			originalWidth:  100,
			originalHeight: 200,
			targetWidth:    200,
			targetHeight:   100,
			expectedWidth:  200,
			expectedHeight: 100,
		},
		{
			name:           "square to wide",
			originalWidth:  100,
			originalHeight: 100,
			targetWidth:    200,
			targetHeight:   100,
			expectedWidth:  200,
			expectedHeight: 100,
		},
		{
			name:           "wide to square",
			originalWidth:  200,
			originalHeight: 100,
			targetWidth:    100,
			targetHeight:   100,
			expectedWidth:  100,
			expectedHeight: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test image with triangle in the top right half
			original := createTestImage(tt.originalWidth, tt.originalHeight, color.RGBA{R: 255, G: 0, B: 0, A: 255})

			// Resize the image
			resized := ResizeImage(original, tt.targetWidth, tt.targetHeight)

			// Check dimensions
			bounds := resized.Bounds()
			assert.Equal(t, tt.expectedWidth, bounds.Dx(), "width should match expected")
			assert.Equal(t, tt.expectedHeight, bounds.Dy(), "height should match expected")
		})
	}
}
