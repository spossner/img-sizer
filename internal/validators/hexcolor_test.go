package validators

import (
	"testing"
)

func TestIsValidHexColor(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid 6-digit hex with #",
			input:    "#FF0000",
			expected: true,
		},
		{
			name:     "valid 6-digit hex without #",
			input:    "00FF00",
			expected: true,
		},
		{
			name:     "valid 6-digit hex with lowercase",
			input:    "#0000ff",
			expected: true,
		},
		{
			name:     "valid 6-digit hex with mixed case",
			input:    "#FF00fF",
			expected: true,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "invalid hex with wrong length",
			input:    "#FF00",
			expected: false,
		},
		{
			name:     "invalid hex with wrong length without #",
			input:    "FF00",
			expected: false,
		},
		{
			name:     "invalid hex with non-hex characters",
			input:    "#FF00GG",
			expected: false,
		},
		{
			name:     "invalid hex with non-hex characters without #",
			input:    "FF00GG",
			expected: false,
		},
		{
			name:     "invalid hex with special characters",
			input:    "#FF00@0",
			expected: false,
		},
		{
			name:     "invalid hex with spaces",
			input:    "#FF 000",
			expected: false,
		},
		{
			name:     "invalid hex with multiple #",
			input:    "##FF0000",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidHexColor(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidHexColor(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}
