package validators

import (
	"image"
	"log/slog"
	"testing"

	"github.com/spossner/img-sizer/internal/config"
)

func TestIsAllowedDimension(t *testing.T) {
	cfg := &config.Config{
		AllowedDimensions: []config.Dimension{
			{Width: 100, Height: 100},
			{Width: 200, Height: 200},
			{Width: 300, Height: 300},
			{Width: 400, Height: 300},
			{Width: 800, Height: 600},
			{Width: 1024, Height: 768},
		},
		AllowAllDimensions: false,
		Logger:             slog.Default(),
	}

	tests := []struct {
		name     string
		width    int
		height   int
		expected bool
	}{
		{
			name:     "exact match",
			width:    100,
			height:   100,
			expected: true,
		},
		{
			name:     "exact match with different dimensions",
			width:    800,
			height:   600,
			expected: true,
		},
		{
			name:     "not in allowed dimensions",
			width:    150,
			height:   150,
			expected: false,
		},
		{
			name:     "zero width",
			width:    0,
			height:   100,
			expected: true,
		},
		{
			name:     "zero height",
			width:    100,
			height:   0,
			expected: true,
		},
		{
			name:     "both zero",
			width:    0,
			height:   0,
			expected: true,
		},
	}

	// Test with AllowAllDimensions enabled
	cfgAllAllowed := &config.Config{
		AllowedDimensions:  cfg.AllowedDimensions,
		AllowAllDimensions: true,
		Logger:             slog.Default(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsAllowedDimension(cfg, tt.width, tt.height)
			if result != tt.expected {
				t.Errorf("IsAllowedDimension(%d, %d) = %v; want %v", tt.width, tt.height, result, tt.expected)
			}

			// With AllowAllDimensions enabled, all dimensions should be allowed
			resultAllAllowed := IsAllowedDimension(cfgAllAllowed, tt.width, tt.height)
			if !resultAllAllowed {
				t.Errorf("IsAllowedDimension with AllowAllDimensions=true(%d, %d) = %v; want true", tt.width, tt.height, resultAllAllowed)
			}
		})
	}
}

func TestValidateInputDimensions(t *testing.T) {
	cfg := &config.Config{
		MaxInputDimension: 5000,
	}

	tests := []struct {
		name    string
		width   int
		height  int
		wantErr bool
	}{
		{
			name:    "valid dimensions",
			width:   1000,
			height:  1000,
			wantErr: false,
		},
		{
			name:    "width too large",
			width:   6000,
			height:  1000,
			wantErr: true,
		},
		{
			name:    "height too large",
			width:   1000,
			height:  6000,
			wantErr: true,
		},
		{
			name:    "both too large",
			width:   6000,
			height:  6000,
			wantErr: true,
		},
		{
			name:    "zero dimensions",
			width:   0,
			height:  0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateInputDimensions(cfg, tt.width, tt.height)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateInputDimensions(%d, %d) error = %v, wantErr %v", tt.width, tt.height, err, tt.wantErr)
			}
		})
	}
}

func TestValidateOutputDimensions(t *testing.T) {
	cfg := &config.Config{
		MaxOutputDimension: 2000,
	}

	tests := []struct {
		name    string
		width   int
		height  int
		wantErr bool
	}{
		{
			name:    "valid dimensions",
			width:   1000,
			height:  1000,
			wantErr: false,
		},
		{
			name:    "width too large",
			width:   3000,
			height:  1000,
			wantErr: true,
		},
		{
			name:    "height too large",
			width:   1000,
			height:  3000,
			wantErr: true,
		},
		{
			name:    "both too large",
			width:   3000,
			height:  3000,
			wantErr: true,
		},
		{
			name:    "zero dimensions",
			width:   0,
			height:  0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOutputDimensions(cfg, tt.width, tt.height)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateOutputDimensions(%d, %d) error = %v, wantErr %v", tt.width, tt.height, err, tt.wantErr)
			}
		})
	}
}

func TestValidateCropZone(t *testing.T) {
	tests := []struct {
		name    string
		width   int
		height  int
		crop    image.Rectangle
		wantErr bool
	}{
		{
			name:    "valid crop zone",
			width:   1000,
			height:  1000,
			crop:    image.Rect(100, 100, 500, 500),
			wantErr: false,
		},
		{
			name:    "negative x",
			width:   1000,
			height:  1000,
			crop:    image.Rect(-100, 100, 500, 500),
			wantErr: true,
		},
		{
			name:    "negative y",
			width:   1000,
			height:  1000,
			crop:    image.Rect(100, -100, 500, 500),
			wantErr: true,
		},
		{
			name:    "x exceeds width",
			width:   1000,
			height:  1000,
			crop:    image.Rect(100, 100, 1100, 500),
			wantErr: true,
		},
		{
			name:    "y exceeds height",
			width:   1000,
			height:  1000,
			crop:    image.Rect(100, 100, 500, 1100),
			wantErr: true,
		},
		{
			name:    "empty rectangle",
			width:   1000,
			height:  1000,
			crop:    image.Rect(100, 100, 100, 100),
			wantErr: true,
		},
		{
			name:    "zero dimensions",
			width:   0,
			height:  0,
			crop:    image.Rect(0, 0, 0, 0),
			wantErr: true,
		},
	}

	cfg := &config.Config{} // No specific config needed for crop validation

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCropZone(cfg, tt.width, tt.height, tt.crop)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCropZone(%d, %d, %v) error = %v, wantErr %v", tt.width, tt.height, tt.crop, err, tt.wantErr)
			}
		})
	}
}
