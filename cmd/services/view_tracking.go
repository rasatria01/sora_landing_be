package services

import (
	"context"
	"sora_landing_be/cmd/domain"
	"sora_landing_be/cmd/repository"
	"sora_landing_be/pkg/observability"
	"time"
)

type ViewTrackingService struct {
	repo    repository.BlogRepository
	metrics *observability.Metrics
	logger  *observability.Logger
	tracer  *observability.Tracer
}

func NewViewTrackingService(
	repo repository.BlogRepository,
	metrics *observability.Metrics,
	logger *observability.Logger,
	tracer *observability.Tracer,
) *ViewTrackingService {
	return &ViewTrackingService{
		repo:    repo,
		metrics: metrics,
		logger:  logger,
		tracer:  tracer,
	}
}

func (s *ViewTrackingService) TrackView(ctx context.Context, article *domain.BlogArtikel) error {
	// Start tracing
	ctx, trace := s.tracer.StartTrace(ctx, "track_article_view")
	defer s.tracer.EndTrace(ctx, trace)

	startTime := time.Now()

	// Add context to trace
	s.tracer.AddTag(ctx, "article_id", article.ID)
	s.tracer.AddTag(ctx, "article_slug", article.Slug)

	// Increment view count
	err := s.repo.IncrementViews(ctx, article.ID)
	if err != nil {
		s.metrics.IncrementCounter("article_view_errors")
		s.logger.Error("Failed to increment view count", err, observability.LogFields{
			"article_id": article.ID,
			"slug":       article.Slug,
		})
		return err
	}

	// Record metrics
	s.metrics.IncrementCounter("article_views_total")
	s.metrics.IncrementCounter("article_views_by_id." + article.ID)
	s.metrics.IncrementCounter("category_views." + article.CategoryID)
	s.metrics.RecordDuration("view_processing_time", time.Since(startTime))

	// Update view rate gauge
	s.updateViewRateGauge(article.ID)

	// Log success
	s.logger.Info("Article view tracked", observability.LogFields{
		"article_id":  article.ID,
		"slug":        article.Slug,
		"category_id": article.CategoryID,
		"duration_ms": time.Since(startTime).Milliseconds(),
	})

	return nil
}

// updateViewRateGauge calculates and updates the views per minute gauge
func (s *ViewTrackingService) updateViewRateGauge(articleID string) {
	// This would typically involve calculating a rolling average of views
	// For now, we'll use a simple gauge update
	s.metrics.SetGauge("article_views_per_minute."+articleID, 1.0)
}
