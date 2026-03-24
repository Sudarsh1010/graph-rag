package di

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Server struct {
	server *http.Server
	logger *zap.Logger
}

type ServerParams struct {
	fx.In

	Config    *Config
	Logger    *zap.Logger
	Router    http.Handler
	Lifecycle fx.Lifecycle
}

func NewServer(p ServerParams) (*Server, error) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", p.Config.Port),
		Handler: p.Router,
	}

	server := &Server{
		server: srv,
		logger: p.Logger,
	}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			p.Logger.Info(
				"Starting HTTP server",
				zap.String("port", p.Config.Port),
				zap.String("env", p.Config.Env),
			)
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					p.Logger.Error("HTTP server error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			p.Logger.Info("Shutting down HTTP server...")
			return srv.Shutdown(ctx)
		},
	})

	return server, nil
}

func (s *Server) Start() error {
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server failed: %w", err)
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
