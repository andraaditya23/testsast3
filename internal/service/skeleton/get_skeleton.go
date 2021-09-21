package skeleton

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

// GetSkeleton ...
func (s Service) GetSkeleton(ctx context.Context) error {
	// Check if have span on context
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := s.tracer.StartSpan("GetSkeleton", opentracing.ChildOf(span.Context()))
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}

	// checkPermission
	// s.checkPermission(ctx, "")

	return nil
}
