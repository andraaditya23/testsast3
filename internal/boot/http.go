package boot

import (
	"backend-pipeline-security/docs"
	"backend-pipeline-security/internal/data/auth"
	"backend-pipeline-security/pkg/httpclient"
	"backend-pipeline-security/pkg/tracing"
	"log"
	"net/http"

	"backend-pipeline-security/internal/config"
	jaegerLog "backend-pipeline-security/pkg/log"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	skeletonData "backend-pipeline-security/internal/data/skeleton"
	skeletonServer "backend-pipeline-security/internal/delivery/http"
	skeletonHandler "backend-pipeline-security/internal/delivery/http/skeleton"
	skeletonService "backend-pipeline-security/internal/service/skeleton"
)

// HTTP will load configuration, do dependency injection and then start the HTTP server
func HTTP() error {
	err := config.Init()
	if err != nil {
		log.Fatalf("[CONFIG] Failed to initialize config: %v", err)
	}
	cfg := config.Get()
	// Open MySQL DB Connection
	db, err := sqlx.Open("mysql", cfg.Database.Master)
	if err != nil {
		log.Fatalf("[DB] Failed to initialize database connection: %v", err)
	}

	//
	docs.SwaggerInfo.Host = cfg.Swagger.Host
	docs.SwaggerInfo.Schemes = cfg.Swagger.Schemes

	// Set logger used for jaeger
	logger, _ := zap.NewDevelopment(
		zap.AddStacktrace(zapcore.FatalLevel),
		zap.AddCallerSkip(1),
	)
	zapLogger := logger.With(zap.String("service", "skeleton"))
	zlogger := jaegerLog.NewFactory(zapLogger)

	// Set tracer for service
	tracer, closer := tracing.Init("skeleton", zlogger)
	defer closer.Close()

	httpc := httpclient.NewClient(tracer)
	ad := auth.New(httpc, cfg.API.Auth)

	// Diganti dengan domain yang anda buat
	sd := skeletonData.New(db, tracer, zlogger)
	ss := skeletonService.New(sd, ad, tracer, zlogger)
	sh := skeletonHandler.New(ss, tracer, zlogger)

	s := skeletonServer.Server{
		Skeleton: sh,
	}

	if err := s.Serve(cfg.Server.Port); err != http.ErrServerClosed {
		return err
	}

	return nil
}
