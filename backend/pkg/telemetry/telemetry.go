// Package telemetry installs the OpenTelemetry meter provider that exports metrics to Cloud Monitoring.
package telemetry

import (
	"context"

	"chatsmith-mini/pkg/failure"

	mexporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

const serviceName = "chatsmith-mini"

func Setup(_ context.Context, projectID string) (func(context.Context) error, error) {
	if projectID == "" {
		return func(context.Context) error { return nil }, nil
	}

	exporter, err := mexporter.New(mexporter.WithProjectID(projectID))
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}
	res, err := resource.Merge(resource.Default(), resource.NewSchemaless(attribute.String("service.name", serviceName)))
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
	)
	otel.SetMeterProvider(provider)

	return provider.Shutdown, nil
}
