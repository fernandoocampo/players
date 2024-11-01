package grpc

import (
	"context"

	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

const scopeName = "players-otl"

func makeTracerUnaryInterceptor(version, appName string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		serverInfo *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		tracer := otel.GetTracerProvider().Tracer(
			scopeName,
			oteltrace.WithInstrumentationVersion(version),
		)

		spanName := serverInfo.FullMethod

		opts := []oteltrace.SpanStartOption{
			oteltrace.WithAttributes(semconv.ServiceName(appName)),
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		}

		ctx, span := tracer.Start(ctx, spanName, opts...)
		defer span.End()

		return handler(ctx, req)
	}
}
