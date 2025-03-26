package handlers

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSizerParamsString(t *testing.T) {
	tests := []struct {
		name     string
		params   SizerParams
		expected string
	}{
		{
			name: "basic parameters",
			params: SizerParams{
				Width:   100,
				Height:  200,
				Quality: 80,
				BgColor: "000000",
				Density: 1.0,
				Scale:   1.0,
				Crop:    image.Rectangle{},
			},
			expected: "100x200-q80-bg000000-d1.00-s1.00-c(0,0)-(0,0)",
		},
		{
			name: "with crop rectangle",
			params: SizerParams{
				Width:   800,
				Height:  600,
				Quality: 90,
				BgColor: "FFFFFF",
				Density: 2.0,
				Scale:   0.5,
				Crop:    image.Rect(10, 20, 810, 620),
			},
			expected: "800x600-q90-bgFFFFFF-d2.00-s0.50-c(10,20)-(810,620)",
		},
		{
			name: "zero values",
			params: SizerParams{
				Width:   0,
				Height:  0,
				Quality: 0,
				BgColor: "",
				Density: 0.0,
				Scale:   0.0,
				Crop:    image.Rectangle{},
			},
			expected: "0x0-q0-bg-d0.00-s0.00-c(0,0)-(0,0)",
		},
		{
			name: "negative values",
			params: SizerParams{
				Width:   -100,
				Height:  -200,
				Quality: 70,
				BgColor: "FF0000",
				Density: 1.5,
				Scale:   0.75,
				Crop:    image.Rect(-10, -20, 90, 180),
			},
			expected: "-100x-200-q70-bgFF0000-d1.50-s0.75-c(-10,-20)-(90,180)",
		},
		{
			name: "decimal values",
			params: SizerParams{
				Width:   1024,
				Height:  768,
				Quality: 85,
				BgColor: "808080",
				Density: 1.25,
				Scale:   0.333,
				Crop:    image.Rect(0, 0, 1024, 768),
			},
			expected: "1024x768-q85-bg808080-d1.25-s0.33-c(0,0)-(1024,768)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.params.String()
			assert.Equal(t, tt.expected, result, "String() should match expected")
		})
	}
}

func TestSizerParamsEquality(t *testing.T) {
	params1 := SizerParams{
		Width:   100,
		Height:  200,
		Quality: 80,
		BgColor: "000000",
		Density: 1.0,
		Scale:   1.0,
		Crop:    image.Rect(0, 0, 100, 200),
	}

	params2 := SizerParams{
		Width:   100,
		Height:  200,
		Quality: 80,
		BgColor: "000000",
		Density: 1.0,
		Scale:   1.0,
		Crop:    image.Rect(0, 0, 100, 200),
	}

	assert.Equal(t, params1.String(), params2.String(), "Identical SizerParams should have identical string representations")

	// Test with different values
	params3 := SizerParams{
		Width:   200,
		Height:  100,
		Quality: 80,
		BgColor: "000000",
		Density: 1.0,
		Scale:   1.0,
		Crop:    image.Rect(0, 0, 200, 100),
	}

	assert.NotEqual(t, params1.String(), params3.String(), "Different SizerParams should have different string representations")
}
