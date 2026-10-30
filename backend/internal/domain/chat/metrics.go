package chat

import (
	"context"
	"time"

	"chatsmith-mini/pkg/failure"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type chatMetrics struct {
	outcome  metric.Int64Counter
	duration metric.Float64Histogram
}

func newChatMetrics() (*chatMetrics, error) {
	meter := otel.Meter(chatMeterName)

	outcome, err := meter.Int64Counter(metricChatOutcome)
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}
	duration, err := meter.Float64Histogram(metricChatDuration,
		metric.WithUnit("s"), metric.WithExplicitBucketBoundaries(chatDurationBuckets...))
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}

	return &chatMetrics{outcome: outcome, duration: duration}, nil
}

func (m *chatMetrics) reject(ctx context.Context, model string, code failure.Code) {
	m.outcome.Add(ctx, 1, metric.WithAttributes(
		attribute.String(attrStatus, statusRejected),
		attribute.String(attrReason, string(code)),
		attribute.String(attrModel, model),
	))
}

func (m *chatMetrics) record(ctx context.Context, model string, status MessageStatus, reason failure.Code, latency time.Duration) {
	attrs := metric.WithAttributes(
		attribute.String(attrStatus, string(status)),
		attribute.String(attrReason, string(reason)),
		attribute.String(attrModel, model),
	)
	m.outcome.Add(ctx, 1, attrs)
	m.duration.Record(ctx, latency.Seconds(), attrs)
}
