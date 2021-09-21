package skeleton

import (
	"context"

	jaegerLog "backend-pipeline-security/pkg/log"

	"github.com/opentracing/opentracing-go"
)

// ISkeletonSvc is an interface to Skeleton Service
// Masukkan function dari service ke dalam interface ini
type ISkeletonSvc interface {
	GetSkeleton(ctx context.Context) error
}

type (
	// Handler ...
	Handler struct {
		skeletonSvc ISkeletonSvc
		tracer      opentracing.Tracer
		logger      jaegerLog.Factory
	}
)

// New for bridging product handler initialization
func New(is ISkeletonSvc, tracer opentracing.Tracer, logger jaegerLog.Factory) *Handler {
	return &Handler{
		skeletonSvc: is,
		tracer:      tracer,
		logger:      logger,
	}
}
