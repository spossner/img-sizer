package utils

import (
	"strconv"
)

// ParseInt parses a string to an integer, returning the default value if parsing fails
func ParseInt(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return val
}

// ParseFloat parses a string to a float64, returning the default value if parsing fails
func ParseFloat(s string, defaultValue float64) float64 {
	if s == "" {
		return defaultValue
	}
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultValue
	}
	return val
}
