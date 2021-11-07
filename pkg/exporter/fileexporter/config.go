package fileexporter

import (
	"errors"

	"go.opentelemetry.io/collector/config"
)

// Config defines configuration for exporter.
type Config struct {
	config.ExporterSettings `mapstructure:",squash"`

	// Path is the file path where this exporter will export data to.
	Path string `mapstructure:"path"`

	// Append indicates whether to append to a file or just truncate it and start
	// from scratch.
	Append bool `mapstructure:"append"`
}

var (
	errPathNotConfigured = errors.New("path not configured")
)

func (c *Config) Validate() error {
	if c.Path == "" {
		return errPathNotConfigured
	}
	return nil
}
