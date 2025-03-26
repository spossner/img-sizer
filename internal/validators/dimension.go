package validators

import (
	"fmt"
	"image"
	"img-sizer/internal/config"
)

// IsAllowedDimension checks if the given width and height match any of the allowed dimensions
func IsAllowedDimension(cfg *config.Config, width, height int) bool {
	if cfg.AllowAllDimensions {
		cfg.Logger.Warn("unknown dimension", "width", width, "height", height)
		return true
	}
	for _, dim := range cfg.AllowedDimensions {
		if (width == 0 || dim.Width == width) && (height == 0 || dim.Height == height) {
			return true
		}
	}
	return false
}

func ValidateInputDimensions(cfg *config.Config, width, height int) error {
	if width > cfg.MaxInputDimension || height > cfg.MaxInputDimension {
		return fmt.Errorf("image dimensions too large")
	}
	return nil
}

func ValidateOutputDimensions(cfg *config.Config, width, height int) error {
	if width > cfg.MaxOutputDimension || height > cfg.MaxOutputDimension {
		return fmt.Errorf("output dimensions too large")
	}
	return nil
}

func ValidateCropZone(cfg *config.Config, width, height int, crop image.Rectangle) error {
	if crop.Min.X < 0 || crop.Min.Y < 0 || crop.Max.X > width || crop.Max.Y > height || crop.Empty() {
		return fmt.Errorf("invalid crop zone")
	}
	return nil
}
