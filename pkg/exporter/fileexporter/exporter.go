package fileexporter

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/kr/pretty"
	"go.opentelemetry.io/collector/model/pdata"
)

type exporter struct {
	config *Config

	file io.Writer
}

func newExporter(config *Config) (*exporter, error) {
	flags := os.O_CREATE | os.O_RDWR | os.O_TRUNC
	if config.Append {
		flags |= os.O_APPEND
	}

	file, err := os.OpenFile(config.Path, flags, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return &exporter{
		file:   file,
		config: config,
	}, nil
}

func (exp *exporter) pushMetricsData(ctx context.Context, md pdata.Metrics) error {
	_, err := fmt.Fprintf(exp.file, "%# v\n", pretty.Formatter(md))
	return err
}
