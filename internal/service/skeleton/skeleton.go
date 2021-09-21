package skeleton

import (
	"backend-pipeline-security/internal/entity"
	"backend-pipeline-security/internal/entity/auth"
	jaegerLog "backend-pipeline-security/pkg/log"
	"context"
	"errors"

	"github.com/opentracing/opentracing-go"
)

// Data ...
// Masukkan function dari package data ke dalam interface ini
type Data interface {
}

// AuthData ...
type AuthData interface {
	CheckAuth(ctx context.Context, _token, code string) (auth.Auth, error)
}

// Service ...
// Tambahkan variable sesuai banyak data layer yang dibutuhkan
type Service struct {
	data     Data
	authData AuthData
	tracer   opentracing.Tracer
	logger   jaegerLog.Factory
}

// New ...
// Tambahkan parameter sesuai banyak data layer yang dibutuhkan
func New(data Data, authData AuthData, tracer opentracing.Tracer, logger jaegerLog.Factory) Service {
	// Assign variable dari parameter ke object
	return Service{
		data:     data,
		authData: authData,
		tracer:   tracer,
		logger:   logger,
	}
}

func (s Service) checkPermission(ctx context.Context, _permissions ...string) error {
	claims := ctx.Value(entity.ContextKey("claims"))
	if claims != nil {
		actions := claims.(entity.ContextValue).Get("permissions").(map[string]interface{})
		for _, action := range actions {
			permissions := action.([]interface{})
			for _, permission := range permissions {
				for _, _permission := range _permissions {
					if permission.(string) == _permission {
						return nil
					}
				}
			}
		}
	}
	return errors.New("401 unauthorized")
}
