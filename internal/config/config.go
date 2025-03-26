package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spossner/img-sizer/internal/utils"
)

type Dimension struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type RateLimit struct {
	MaxRequests int           `json:"max_requests"`
	Window      time.Duration `json:"window"`
}

func (r *RateLimit) UnmarshalJSON(data []byte) error {
	type Alias RateLimit
	aux := &struct {
		Window string `json:"window"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	duration, err := time.ParseDuration(aux.Window)
	if err != nil {
		return fmt.Errorf("invalid duration format: %v", err)
	}
	r.Window = duration
	return nil
}

type SourceConfig struct {
	Pattern *regexp.Regexp `json:"pattern"`
	Matcher *regexp.Regexp `json:"matcher,omitempty"`
	Bucket  string         `json:"bucket,omitempty"`
}

func (s *SourceConfig) UnmarshalJSON(data []byte) error {
	type Alias SourceConfig
	aux := &struct {
		Pattern string `json:"pattern"`
		Matcher string `json:"matcher,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return fmt.Errorf("error unmarshalling source config: %v", err)
	}

	if aux.Pattern != "" {
		// Convert wildcard pattern to regex
		pattern := strings.ReplaceAll(aux.Pattern, ".", "\\.")
		pattern = strings.ReplaceAll(pattern, "*", "[a-zA-Z0-9-]+")
		pattern = "^" + pattern + ".*$"

		re, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("error compiling pattern %s: %v", aux.Pattern, err)
		}
		s.Pattern = re
	}
	if aux.Matcher != "" {
		re, err := regexp.Compile(aux.Matcher)
		if err != nil {
			return fmt.Errorf("error compiling matcher %s: %v", aux.Matcher, err)
		}
		s.Matcher = re
	}
	return nil
}

type Jpeg struct {
	Background string `json:"background"`
	Quality    int    `json:"quality"`
}

type Config struct {
	AllowedSources     []SourceConfig `json:"allowed_sources"`
	AllowedDimensions  []Dimension    `json:"allowed_dimensions"`
	AllowAllDimensions bool           `json:"allow_all_dimensions"`
	MaxInputDimension  int            `json:"max_input_dimension"`
	MaxOutputDimension int            `json:"max_output_dimension"`
	RateLimit          RateLimit      `json:"rate_limit"`
	Jpeg               Jpeg           `json:"jpeg"`
	Logger             *slog.Logger   `json:"-"`
}

func Load(logger *slog.Logger) (*Config, error) {
	// Get environment from ENV or default to "local"
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}

	// Determine config path based on environment
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		projectRoot := utils.GetProjectRoot(logger)
		configPath = filepath.Join(projectRoot, "config", strings.ToLower(env)+".json")
	}

	// Load config
	logger.Info("using config file", "file", configPath)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file %s: %v", configPath, err)
	}

	// Parse config
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error decoding config: %v", err)
	}

	// Set logger
	config.Logger = logger

	// Set default rate limit if not configured
	if config.RateLimit.MaxRequests == 0 {
		config.RateLimit.MaxRequests = 300
	}
	if config.RateLimit.Window == 0 {
		config.RateLimit.Window = 1 * time.Minute
	}

	for i, source := range config.AllowedSources {
		logger.Info("allowed source", "index", i, "pattern", source.Pattern, "bucket", source.Bucket, "matcher", source.Matcher)
	}

	return &config, nil
}
