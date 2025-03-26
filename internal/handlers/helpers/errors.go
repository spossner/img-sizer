package helpers

import "errors"

var (
	ErrInvalidURL      = errors.New("invalid URL")
	ErrLoadingImage    = errors.New("error loading image")
	ErrURLNotAllowed   = errors.New("URL not allowed")
	ErrProcessingImage = errors.New("error processing image")
)
