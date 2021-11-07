package httprequestreceiver

import (
	"context"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/model/pdata"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.uber.org/zap"
)

const (
	metricName        = "request_content_length"
	metricDescription = "Describes the request content length"
	metricUnit        = "byte(s)"

	tagNameMethod     = "method"
	tagNameRequestURI = "request_uri"

	internalChannelSize = 1024
)

type receiver struct {
	consumer consumer.Metrics
	logger   *zap.Logger
	config   *Config

	chQuit chan struct{}
}

// Ensure this receiver adheres to required interface.
var _ component.MetricsReceiver = (*receiver)(nil)

func newReceiver(
	nextConsumer consumer.Metrics,
	logger *zap.Logger,
	config *Config,
) *receiver {
	return &receiver{
		consumer: nextConsumer,
		logger:   logger,
		config:   config,
		chQuit:   make(chan struct{}),
	}
}

// Start tells the receiver to start.
func (r *receiver) Start(ctx context.Context, host component.Host) error {
	type metricData struct {
		ContentLength int64
		Method        string
		RequestURI    string
	}
	ch := make(chan metricData, internalChannelSize)

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case md := <-ch:
				// Create pdata.Metrics
				metrics := pdata.NewMetrics()
				resMetricsSlice := metrics.ResourceMetrics()
				// Append new pdata.ResourceMetrics
				resMetrics := resMetricsSlice.AppendEmpty()

				// Get the resource attributes to manipulate them
				resAttrs := resMetrics.Resource().Attributes()
				resAttrs.InsertString(string(semconv.HostNameKey), hostname)

				// Append new pdata.InstrumentationLibraryMetrics.
				ilms := resMetrics.InstrumentationLibraryMetrics().AppendEmpty()

				// Create new pdata.Metric and fill describe it.
				metric := ilms.Metrics().AppendEmpty()
				metric.SetName(metricName)
				metric.SetDataType(pdata.MetricDataTypeGauge)
				metric.SetDescription(metricDescription)
				metric.SetUnit(metricUnit)

				// Add a data point and set its value and tags.
				dp := metric.Gauge().DataPoints().AppendEmpty()
				dp.SetTimestamp(pdata.NewTimestampFromTime(time.Now()))
				dp.SetIntVal(md.ContentLength)
				dp.Attributes().InitFromMap(map[string]pdata.AttributeValue{
					tagNameMethod:     pdata.NewAttributeValueString(md.Method),
					tagNameRequestURI: pdata.NewAttributeValueString(md.RequestURI),
				})

				// Forward it to the next consumer
				err := r.consumer.ConsumeMetrics(context.Background(), metrics)
				if err != nil {
					r.logger.Error("Failed sending metrics", zap.Error(err))
				}

			case <-r.chQuit:
				r.logger.Info("Closing the worker goroutine...")
				return
			}
		}
	}()

	go func() {
		http.ListenAndServe(r.config.Address, http.HandlerFunc(
			func(w http.ResponseWriter, req *http.Request) {
				ch <- metricData{
					ContentLength: req.ContentLength,
					Method:        req.Method,
					RequestURI:    req.RequestURI,
				}
			}))
	}()

	r.logger.Info("Receiver started", zap.String("address", r.config.Address))

	return nil
}

// Shutdown tells the receiver to shutdown.
func (r *receiver) Shutdown(context.Context) error {
	r.logger.Info("Receiver shutdown")
	close(r.chQuit)

	return nil
}
