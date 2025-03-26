package handlers

import (
	"bytes"
	"image"
	"image/jpeg"
	"net/http"

	"github.com/spossner/img-sizer/internal/config"
	"github.com/spossner/img-sizer/internal/handlers/helpers"
	"github.com/spossner/img-sizer/internal/processing"
	"github.com/spossner/img-sizer/internal/storage"
	"github.com/spossner/img-sizer/internal/utils"
	"github.com/spossner/img-sizer/internal/validators"

	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
)

func GetImageSizerHandler(cfg *config.Config, s3Client *storage.S3Client, paramsParser ParamsParser) fiber.Handler {
	return func(c *fiber.Ctx) error {
		params := paramsParser(c, cfg)
		if !validators.IsAllowedDimension(cfg, params.Width, params.Height) {
			cfg.Logger.Error("invalid dimensions", "width", params.Width, "height", params.Height)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid dimensions",
			})
		}

		if err := validators.ValidateOutputDimensions(cfg, params.Width, params.Height); err != nil {
			cfg.Logger.Error("requested output dimensions exceed limit", "width", params.Width, "height", params.Height)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "requested output dimensions exceed limit",
			})
		}

		sourceURL := c.Query("src")
		if sourceURL == "" {
			cfg.Logger.Error("source URL is required")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "source URL is required",
			})
		}

		bucket, key, err := helpers.ParseS3Url(cfg, sourceURL)
		if err != nil {
			cfg.Logger.Error("invalid source URL", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid source URL",
			})
		}

		var img image.Image
		if bucket != "" { // load from s3 if bucket is configured
			img, err = helpers.LoadImageFromS3(c.Context(), s3Client, bucket, key)
			if err != nil {
				cfg.Logger.Error("error loading image from S3", "bucket", bucket, "key", key, "error", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "error processing image",
				})
			}
		} else {
			cfg.Logger.Warn("unmapped source URL - loading from URL", "url", sourceURL)
			img, err = helpers.LoadImageFromURL(c.Context(), sourceURL)
			if err != nil {
				cfg.Logger.Error("error loading image from URL", "url", sourceURL, "error", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "error processing image",
				})
			}
		}

		if err = validators.ValidateInputDimensions(cfg, img.Bounds().Size().X, img.Bounds().Size().Y); err != nil {
			cfg.Logger.Error("image dimensions exceed limit", "bucket", bucket, "key", key, "width", img.Bounds().Size().X, "height", img.Bounds().Size().Y)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "image dimensions exceed limit",
			})
		}

		if params.Crop.Dx() > 0 && params.Crop.Dy() > 0 {
			if err = validators.ValidateCropZone(cfg, img.Bounds().Size().X, img.Bounds().Size().Y, params.Crop); err != nil {
				cfg.Logger.Error("invalid crop zone", "bounds", img.Bounds().Size(), "crop", params.Crop)
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "invalid crop zone",
				})
			}
			// Crop image
			img = imaging.Crop(img, params.Crop)
		}

		// Resize image
		img = processing.ResizeImage(img, params.Width, params.Height)

		// Fill background
		img, err = processing.FillBackground(img, params.BgColor)
		if err != nil {
			cfg.Logger.Error("invalid background color", "bgColor", params.BgColor, "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid background color",
			})
		}

		// Create a buffer to store the encoded image
		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: params.Quality}); err != nil {
			cfg.Logger.Error("error encoding image", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "error processing image",
			})
		}

		// Response header with ETag
		etag := utils.CalculateETag(buf.Bytes(), params.String())
		helpers.SetResponseHeaders(c, etag)

		// Check if client has matching ETag
		if match := c.Get("If-None-Match"); match == etag {
			return c.Status(http.StatusNotModified).Send(nil)
		}
		// Send the processed image
		return c.Send(buf.Bytes())
	}
}
