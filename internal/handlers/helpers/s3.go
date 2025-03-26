package helpers

import (
	"context"
	"image"
	"net/url"
	"strings"

	"github.com/spossner/img-sizer/internal/config"
	"github.com/spossner/img-sizer/internal/storage"

	"github.com/disintegration/imaging"
)

// ParseS3Url parses the S3 URL and returns the bucket name and key or an error if the URL is invalid or not configured
func ParseS3Url(cfg *config.Config, sourceURL string) (string, string, error) {
	parsedURL, err := url.Parse(sourceURL)
	if err != nil {
		return "", "", ErrInvalidURL
	}

	// The key is the path without the leading slash
	key := strings.TrimPrefix(parsedURL.Path, "/")

	// Find the bucket name from pattern
	for _, source := range cfg.AllowedSources {
		if source.Pattern.MatchString(parsedURL.Host) {
			if source.Matcher != nil {
				matches := source.Matcher.FindStringSubmatch(sourceURL)
				if len(matches) == 3 {
					return matches[1], matches[2], nil
				}
				return "", "", ErrInvalidURL
			}
			return source.Bucket, key, nil
		}
	}

	return "", "", ErrURLNotAllowed
}

func LoadImageFromS3(ctx context.Context, s3Client *storage.S3Client, bucket, key string) (image.Image, error) {
	// Download image from S3
	reader, err := s3Client.GetObject(ctx, bucket, key)
	if err != nil {
		return nil, ErrLoadingImage
	}
	defer reader.Close()

	// Decode image
	img, err := imaging.Decode(reader)
	if err != nil {
		return nil, ErrProcessingImage
	}
	return img, nil
}
