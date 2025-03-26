package processing

import (
	"encoding/hex"
	"errors"
	"image"
	"image/color"
	"img-sizer/internal/validators"

	"github.com/disintegration/imaging"
)

const (
	Black = "000000"
)

func FillBackground(img image.Image, bgColor string) (image.Image, error) {
	if bgColor == Black {
		return img, nil
	}

	if !validators.IsValidHexColor(bgColor) {
		return nil, errors.New("invalid hex color")
	}

	hexBytes, err := hex.DecodeString(bgColor)
	if err != nil {
		return nil, err
	}

	bgImg := imaging.New(img.Bounds().Dx(), img.Bounds().Dy(), color.RGBA{
		R: hexBytes[0],
		G: hexBytes[1],
		B: hexBytes[2],
		A: 255,
	})

	// Composite original image over background
	return imaging.OverlayCenter(bgImg, img, 1.0), nil
}
