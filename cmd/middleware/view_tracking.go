package middleware

import (
	"fmt"
	"sora_landing_be/pkg/observability"
	"time"

	"github.com/gin-gonic/gin"
)

type ViewTrackingMiddleware struct {
	metrics *observability.Metrics
	logger  *observability.Logger
	tracer  *observability.Tracer
}

func NewViewTrackingMiddleware(metrics *observability.Metrics, logger *observability.Logger, tracer *observability.Tracer) *ViewTrackingMiddleware {
	return &ViewTrackingMiddleware{
		metrics: metrics,
		logger:  logger,
		tracer:  tracer,
	}
}

func (m *ViewTrackingMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Start trace for view tracking
		ctx, trace := m.tracer.StartTrace(c.Request.Context(), "article_view_tracking")
		c.Request = c.Request.WithContext(ctx)

		// Add initial trace tags
		m.tracer.AddTag(ctx, "path", c.Request.URL.Path)
		m.tracer.AddTag(ctx, "method", c.Request.Method)

		// Process request
		c.Next()

		// Record metrics
		duration := time.Since(startTime)

		// Track views by path
		m.metrics.IncrementCounter("article_views_total")
		m.metrics.IncrementCounter("article_views_by_path." + c.Param("slug"))

		// Track response time
		m.metrics.RecordDuration("article_view_processing_time", duration)

		// Track concurrent views
		m.metrics.SetGauge("concurrent_article_views", float64(c.Writer.Size()))

		// Complete trace
		m.tracer.AddTag(ctx, "status", fmt.Sprint(c.Writer.Status()))
		m.tracer.AddTag(ctx, "duration_ms", fmt.Sprintf("%d", duration.Milliseconds()))
		m.tracer.EndTrace(ctx, trace)

		// Log view with context
		m.logger.WithContext(ctx).Info("Article view processed", observability.LogFields{
			"path":        c.Request.URL.Path,
			"status":      c.Writer.Status(),
			"duration_ms": duration.Milliseconds(),
			"user_agent":  c.Request.UserAgent(),
			"referer":     c.Request.Referer(),
		})
	}
}
