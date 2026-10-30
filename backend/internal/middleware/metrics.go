package middleware

import (
	"time"

	"chatsmith-mini/pkg/failure"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	httpMeterName         = "chatsmith/http"
	metricRequestDuration = "http.server.request.duration"
	attrHTTPRoute         = "http.route"
	attrHTTPStatus        = "http.response.status_code"
	unmatchedRoute        = "unmatched"
)

var requestDurationBuckets = []float64{0.05, 0.1, 0.25, 0.5, 1, 2, 3, 5, 10, 30, 60}

func Metrics() (gin.HandlerFunc, error) {
	duration, err := otel.Meter(httpMeterName).Float64Histogram(metricRequestDuration,
		metric.WithUnit("s"), metric.WithExplicitBucketBoundaries(requestDurationBuckets...))
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}

	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()

		route := ctx.FullPath()
		if route == "" {
			route = unmatchedRoute
		}
		duration.Record(ctx.Request.Context(), time.Since(start).Seconds(), metric.WithAttributes(
			attribute.String(attrHTTPRoute, route),
			attribute.Int(attrHTTPStatus, ctx.Writer.Status()),
		))
	}, nil
}
