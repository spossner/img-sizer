package processing

import (
	"image"

	"github.com/disintegration/imaging"
)

func ResizeImage(img image.Image, width, height int) image.Image {
	// Resize image
	if width > 0 && height > 0 {
		img = imaging.Fill(img, width, height, imaging.Center, imaging.Lanczos)
	} else if width > 0 || height > 0 {
		img = imaging.Resize(img, width, height, imaging.Lanczos)
	}

	return img
}
