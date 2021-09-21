package http

import (
	"net/http"

	"backend-pipeline-security/pkg/grace"

	"github.com/rs/cors"
)

// SkeletonHandler ...
type SkeletonHandler interface {
	GetSkeleton(w http.ResponseWriter, r *http.Request)
}

// Server ...
type Server struct {
	server   *http.Server
	Skeleton SkeletonHandler
}

// Serve is serving HTTP gracefully on port x ...
func (s *Server) Serve(port string) error {
	handler := cors.AllowAll().Handler(s.Handler())
	return grace.Serve(port, handler)
}
