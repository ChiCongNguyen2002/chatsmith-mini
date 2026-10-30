// Package telemetry installs OpenTelemetry metrics and traces that export to Cloud Monitoring and Cloud Trace.
package telemetry

import (
	"context"
	"errors"

	"chatsmith-mini/constants"
	"chatsmith-mini/pkg/failure"

	mexporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric"
	texporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func Setup(_ context.Context, projectID string) (func(context.Context) error, error) {
	otel.SetTextMapPropagator(propagation.TraceContext{})
	if projectID == "" {
		return func(context.Context) error { return nil }, nil
	}

	res, err := resource.Merge(resource.Default(), resource.NewSchemaless(attribute.String("service.name", constants.ServiceName)))
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}
	metricExporter, err := mexporter.New(mexporter.WithProjectID(projectID))
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}
	traceExporter, err := texporter.New(texporter.WithProjectID(projectID))
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}

	meter := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
	)
	tracer := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(traceExporter),
	)
	otel.SetMeterProvider(meter)
	otel.SetTracerProvider(tracer)

	return func(ctx context.Context) error {
		return errors.Join(tracer.Shutdown(ctx), meter.Shutdown(ctx))
	}, nil
}
