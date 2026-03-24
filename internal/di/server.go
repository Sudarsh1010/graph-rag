package di

import (
	"context"
	"fmt"
	go_http "net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Server struct {
	server *go_http.Server
	logger *zap.Logger
}

type ServerParams struct {
	fx.In

	Config *Config
	Logger *zap.Logger
	// Router    *chi.Mux
	Lifecycle fx.Lifecycle
}

func NewServer(p ServerParams) (*Server, error) {
	srv := &go_http.Server{
		Addr: fmt.Sprintf(":%s", p.Config.Port),
		// Handler: p.Router,
	}

	server := &Server{
		server: srv,
		logger: p.Logger,
	}

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			p.Logger.Info(
				"Starting HTTP server",
				zap.String("port", p.Config.Port),
				zap.String("env", p.Config.Env),
			)
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
	if err != nil && err != go_http.ErrServerClosed {
		return fmt.Errorf("server failed: %w", err)
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
