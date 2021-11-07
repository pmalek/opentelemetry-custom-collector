package httprequestreceiver

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver/receiverhelper"
)

const (
	typeStr    = "httprequest"
	versionStr = "v0.0.1"
)

// NewFactory creates a factory for httprequestreceiver.
func NewFactory() component.ReceiverFactory {
	return receiverhelper.NewFactory(
		typeStr,
		createDefaultConfig,
		receiverhelper.WithMetrics(createMetricsReceiver),
	)
}

func createDefaultConfig() config.Receiver {
	rs := config.NewReceiverSettings(config.NewComponentID(typeStr))

	return &Config{
		ReceiverSettings: &rs,
		Address:          DefaultAddress,
	}
}

// createMetricsReceiver creates a metrics receiver based on provided config.
func createMetricsReceiver(
	ctx context.Context,
	params component.ReceiverCreateSettings,
	cfg config.Receiver,
	nextConsumer consumer.Metrics,
) (component.MetricsReceiver, error) {
	rConfig, ok := cfg.(*Config)
	if !ok {
		return nil, fmt.Errorf("failed reading httprequestreceiver config from config.Receiver")
	}

	return newReceiver(nextConsumer, params.Logger, rConfig), nil
}
