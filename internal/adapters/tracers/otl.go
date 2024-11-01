package tracers

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type TracerSetup struct {
	Exporter    *otlptrace.Exporter
	ServiceName string
}

type ExporterSetup struct {
	TraceCollectorUrl string
}

type TracerService struct {
	tp        *sdktrace.TracerProvider
	CloseFunc func(context.Context) error
}

// NewTracerService creates a new tracer service
func NewTracerService(setup *TracerSetup) *TracerService {
	sdktraceResource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(setup.ServiceName),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(
			setup.Exporter,
			sdktrace.WithMaxExportBatchSize(sdktrace.DefaultMaxExportBatchSize),
			sdktrace.WithBatchTimeout(sdktrace.DefaultScheduleDelay*time.Millisecond),
			sdktrace.WithMaxExportBatchSize(sdktrace.DefaultMaxExportBatchSize),
		),
		sdktrace.WithResource(sdktraceResource),
	)

	otel.SetTracerProvider(tp)

	return &TracerService{
		tp:        tp,
		CloseFunc: tp.Shutdown,
	}
}

// CreateOTLPExporter creates a new OTL exporter and starts it
func CreateOTLPExporter(ctx context.Context, setup *ExporterSetup) (*otlptrace.Exporter, error) {
	exporter, err := otlptrace.New(ctx, newOTLPTraceClient(setup))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to the trace collector: %w", err)
	}

	return exporter, nil
}

func (t *TracerService) TraceProvider() *sdktrace.TracerProvider {
	return t.tp
}

func (t *TracerService) Close() error {
	return t.CloseFunc(context.TODO())
}

func newOTLPTraceClient(setup *ExporterSetup) otlptrace.Client {
	return otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(setup.TraceCollectorUrl),
	)
}
