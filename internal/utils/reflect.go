package utils

import (
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
)

func GetProjectRoot(logger *slog.Logger) string {
	pc, fpath, line, ok := runtime.Caller(1)
	if !ok {
		logger.Error("failed to get caller info", "pc", pc, "fpath", fpath, "line", line)
		return ""
	}

	dir := filepath.Dir(fpath)

	// Walk up the directory tree until we find go.mod
	for {
		// Check for go.mod in current directory
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		// Get parent directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// We've reached the root without finding go.mod
			logger.Error("failed to find project root")
			return ""
		}
		dir = parent
	}
}
