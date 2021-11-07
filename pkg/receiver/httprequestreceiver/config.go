package httprequestreceiver

import (
	"go.opentelemetry.io/collector/config"
)

const (
	DefaultAddress string = "localhost:9999"
)

type Config struct {
	*config.ReceiverSettings `mapstructure:"-"`

	// Address on which to listen on.
	// Default localhost:9999
	Address string `mapstructure:"address"`
}

func (c *Config) Validate() error {
	// Perhaps validate the configured address?
	return nil
}
