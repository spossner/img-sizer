package helpers

import (
	"img-sizer/internal/config"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseS3Url(t *testing.T) {
	// Test cases with different configurations
	testCases := []struct {
		name           string
		configContent  string
		urlToParse     string
		expectedBucket string
		expectedKey    string
		expectError    bool
	}{
		{
			name: "exact match with bucket",
			configContent: `{
				"allowed_sources": [
					{
						"pattern": "static.nebenan.de",
						"bucket": "nebenande"
					}
				]
			}`,
			urlToParse:     "https://static.nebenan.de/images/test.jpg",
			expectedBucket: "nebenande",
			expectedKey:    "images/test.jpg",
			expectError:    false,
		},
		{
			name: "wildcard match with bucket",
			configContent: `{
				"allowed_sources": [
					{
						"pattern": "*.nebenan.de",
						"bucket": "nebenande"
					}
				]
			}`,
			urlToParse:     "https://cdn.nebenan.de/images/test.jpg",
			expectedBucket: "nebenande",
			expectedKey:    "images/test.jpg",
			expectError:    false,
		},
		{
			name: "exact match without bucket",
			configContent: `{
				"allowed_sources": [
					{
						"pattern": "static.nebenan.de"
					}
				]
			}`,
			urlToParse:     "https://static.nebenan.de/images/test.jpg",
			expectedBucket: "",
			expectedKey:    "images/test.jpg",
			expectError:    false,
		},
		{
			name: "first matching pattern wins",
			configContent: `{
				"allowed_sources": [
					{
						"pattern": "static.nebenan.de",
						"bucket": "nebenande"
					},
					{
						"pattern": "*.nebenan.de",
						"bucket": "other-bucket"
					}
				]
			}`,
			urlToParse:     "https://static.nebenan.de/images/test.jpg",
			expectedBucket: "nebenande",
			expectedKey:    "images/test.jpg",
			expectError:    false,
		},
		{
			name: "no matching pattern",
			configContent: `{
				"allowed_sources": [
					{
						"pattern": "static.nebenan.de"
					}
				]
			}`,
			urlToParse:  "https://other-domain.com/images/test.jpg",
			expectError: true,
		},
		{
			name: "invalid URL",
			configContent: `{
				"allowed_sources": [
					{
						"pattern": "static.nebenan.de"
					}
				]
			}`,
			urlToParse:  "not-a-url",
			expectError: true,
		},
		{
			name: "empty URL",
			configContent: `{
				"allowed_sources": [
					{
						"pattern": "static.nebenan.de"
					}
				]
			}`,
			urlToParse:  "",
			expectError: true,
		},
		{
			name: "S3 URL with bucket",
			configContent: `{
				"allowed_sources": [
					{
						"pattern": "possnerde.s3.eu-central-1.amazonaws.com",
						"bucket": "possnerde"
					}
				]
			}`,
			urlToParse:     "https://possnerde.s3.eu-central-1.amazonaws.com/images/test.jpg",
			expectedBucket: "possnerde",
			expectedKey:    "images/test.jpg",
			expectError:    false,
		},
		{
			name: "S3 URL without bucket",
			configContent: `{
				"allowed_sources": [
					{
						"pattern": "possnerde.s3.eu-central-1.amazonaws.com"
					}
				]
			}`,
			urlToParse:     "https://possnerde.s3.eu-central-1.amazonaws.com/images/test.jpg",
			expectedBucket: "",
			expectedKey:    "images/test.jpg",
			expectError:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create temporary config file
			tmpfile, err := os.CreateTemp("", "config-*.json")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpfile.Name())

			// Write test config
			if _, err := tmpfile.Write([]byte(tc.configContent)); err != nil {
				t.Fatal(err)
			}

			// Set environment variable to use our test config
			os.Setenv("CONFIG_PATH", tmpfile.Name())

			// Load config
			cfg, err := config.Load(slog.Default())
			if err != nil {
				t.Fatal(err)
			}

			// Parse URL
			bucket, key, err := ParseS3Url(cfg, tc.urlToParse)

			// Check error
			if tc.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Check results
			assert.Equal(t, tc.expectedBucket, bucket)
			assert.Equal(t, tc.expectedKey, key)
		})
	}
}
