package validators

import (
	"regexp"
	"strings"
)

func IsValidHexColor(color string) bool {
	if color == "" {
		return false
	}

	// Remove # if present
	color = strings.TrimPrefix(color, "#")

	// Check if it's a valid 6-digit hex color
	validHex := regexp.MustCompile(`^[0-9A-Fa-f]{6}$`)
	return validHex.MatchString(color)
}
