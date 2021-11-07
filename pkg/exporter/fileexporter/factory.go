package fileexporter

import (
	"context"
	"errors"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	// The value of "type" key in configuration.
	typeStr = "file"
)

func NewFactory() component.ExporterFactory {
	return exporterhelper.NewFactory(
		typeStr,
		createDefaultConfig,
		exporterhelper.WithMetrics(createMetricsExporter),
	)
}

func createDefaultConfig() config.Exporter {
	return &Config{
		ExporterSettings: config.NewExporterSettings(config.NewComponentID(typeStr)),
	}
}

func createMetricsExporter(
	_ context.Context,
	params component.ExporterCreateSettings,
	cfg config.Exporter,
) (component.MetricsExporter, error) {

	expConfig, ok := cfg.(*Config)
	if !ok {
		return nil, errors.New("failed to convert ot config to fileexporter config")
	}

	exp, err := newExporter(expConfig)
	if err != nil {
		return nil, err
	}

	return exporterhelper.NewMetricsExporter(
		cfg,
		params,
		exp.pushMetricsData,
	)
}
