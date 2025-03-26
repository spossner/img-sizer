package helpers

import (
	"context"
	"image"
	"net/http"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
)

func LoadImageFromURL(ctx context.Context, url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, ErrLoadingImage
	}
	defer resp.Body.Close()

	img, err := imaging.Decode(resp.Body)
	if err != nil {
		return nil, ErrProcessingImage
	}
	return img, nil
}

func SetResponseHeaders(c *fiber.Ctx, etag string) {
	c.Set("Content-Type", "image/jpeg")
	c.Set("ETag", etag)
	c.Set("Cache-Control", "public, max-age=2592000, immutable")
	c.Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
}
