package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sudarsh1010/graph-rag/internal/di"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main() {
	app := fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),

		// Core modules
		di.Module,

		// Lifecycle hooks
		fx.Invoke(func(
			lifecycle fx.Lifecycle,
			server *di.Server,
			logger *zap.Logger,
		) {
			lifecycle.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					go func() {
						if err := server.Start(); err != nil {
							logger.Error("server failed", zap.Error(err))
						}
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return server.Stop(ctx)
				},
			})
		}),
	)

	if err := app.Start(context.Background()); err != nil {
		log.Fatal(err)
	}

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	if err := app.Stop(context.Background()); err != nil {
		log.Fatal(err)
	}
}
