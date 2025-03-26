package utils

import (
	"testing"
)

func TestCalculateETag(t *testing.T) {
	tests := []struct {
		name     string
		img      []byte
		params   string
		expected string
	}{
		{
			name:     "empty image and params",
			img:      []byte{},
			params:   "",
			expected: "\"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\"",
		},
		{
			name:     "simple image with no params",
			img:      []byte("test image"),
			params:   "",
			expected: "\"1187327c6d0f0b0b19b33ab211a549023aa9a41f359c6d0a827d7bd99f8d5994\"",
		},
		{
			name:     "empty image with params",
			img:      []byte{},
			params:   "width=100&height=200",
			expected: "\"7fa23bffc0744388748ccf7a4cf32aa03b9585c7ddad6642e3009280e6699f16\"",
		},
		{
			name:     "image with params",
			img:      []byte("test image"),
			params:   "width=100&height=200",
			expected: "\"b069e63657dd034c0ee70cabbd83e3ddfc17bb4f4dacbd5738a952d82382b252\"",
		},
		{
			name:     "same image different params",
			img:      []byte("test image"),
			params:   "width=200&height=100",
			expected: "\"79d5c5e301e66c9923f28c05f87954aaffbb699139aff495282b6b9c3c047400\"",
		},
		{
			name:     "different image same params",
			img:      []byte("different image"),
			params:   "width=100&height=200",
			expected: "\"edba1ca32a6b7e087b05f003d9e330474f2311b1cba028bbf83de72e6fb235ab\"",
		},
		{
			name:     "binary image data",
			img:      []byte{0x00, 0xFF, 0x42, 0xDE, 0xAD, 0xBE, 0xEF},
			params:   "width=100&height=200",
			expected: "\"2453e4a8bb8e70384df47d51b8c4bc319b1788e2746c486dc6415e8b8365c8f0\"",
		},
		{
			name:     "special characters in params",
			img:      []byte("test image"),
			params:   "width=100&height=200&quality=90&background=#FF0000",
			expected: "\"cdca577f33a4d4477ef5710efa4ddf9cf307d9b19cee0c5a54d4f1f3a3fc5594\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateETag(tt.img, tt.params)
			if result != tt.expected {
				t.Errorf("CalculateETag(%v, %q) = %q; want %q", tt.img, tt.params, result, tt.expected)
			}
		})
	}
}

// Test that ETags are consistent for the same input
func TestCalculateETagConsistency(t *testing.T) {
	img := []byte("test image")
	params := "width=100&height=200"

	// Calculate ETag multiple times
	etag1 := CalculateETag(img, params)
	etag2 := CalculateETag(img, params)
	etag3 := CalculateETag(img, params)

	// All ETags should be the same
	if etag1 != etag2 || etag2 != etag3 {
		t.Errorf("ETags are not consistent: %q, %q, %q", etag1, etag2, etag3)
	}
}

// Test that ETags are different for different inputs
func TestCalculateETagUniqueness(t *testing.T) {
	// Test with different images
	etag1 := CalculateETag([]byte("image1"), "width=100&height=200")
	etag2 := CalculateETag([]byte("image2"), "width=100&height=200")
	if etag1 == etag2 {
		t.Error("ETags should be different for different images")
	}

	// Test with different params
	etag3 := CalculateETag([]byte("image1"), "width=200&height=100")
	if etag1 == etag3 {
		t.Error("ETags should be different for different params")
	}
}
