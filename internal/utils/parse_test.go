package utils

import (
	"testing"
)

func TestParseInt(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		defaultValue  int
		expectedValue int
	}{
		{
			name:          "valid positive integer",
			input:         "123",
			defaultValue:  0,
			expectedValue: 123,
		},
		{
			name:          "valid negative integer",
			input:         "-456",
			defaultValue:  0,
			expectedValue: -456,
		},
		{
			name:          "empty string",
			input:         "",
			defaultValue:  42,
			expectedValue: 42,
		},
		{
			name:          "invalid integer",
			input:         "abc",
			defaultValue:  42,
			expectedValue: 42,
		},
		{
			name:          "float string",
			input:         "123.45",
			defaultValue:  42,
			expectedValue: 42,
		},
		{
			name:          "zero value",
			input:         "0",
			defaultValue:  42,
			expectedValue: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseInt(tt.input, tt.defaultValue)
			if result != tt.expectedValue {
				t.Errorf("ParseInt(%q, %d) = %d; want %d", tt.input, tt.defaultValue, result, tt.expectedValue)
			}
		})
	}
}

func TestParseFloat(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		defaultValue  float64
		expectedValue float64
	}{
		{
			name:          "valid positive float",
			input:         "123.45",
			defaultValue:  0.0,
			expectedValue: 123.45,
		},
		{
			name:          "valid negative float",
			input:         "-456.78",
			defaultValue:  0.0,
			expectedValue: -456.78,
		},
		{
			name:          "valid integer as float",
			input:         "123",
			defaultValue:  0.0,
			expectedValue: 123.0,
		},
		{
			name:          "empty string",
			input:         "",
			defaultValue:  42.5,
			expectedValue: 42.5,
		},
		{
			name:          "invalid float",
			input:         "abc",
			defaultValue:  42.5,
			expectedValue: 42.5,
		},
		{
			name:          "zero value",
			input:         "0",
			defaultValue:  42.5,
			expectedValue: 0.0,
		},
		{
			name:          "scientific notation",
			input:         "1.23e-4",
			defaultValue:  0.0,
			expectedValue: 0.000123,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseFloat(tt.input, tt.defaultValue)
			if result != tt.expectedValue {
				t.Errorf("ParseFloat(%q, %f) = %f; want %f", tt.input, tt.defaultValue, result, tt.expectedValue)
			}
		})
	}
}
