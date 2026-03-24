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

		// Ensure Server is constructed (triggers Router + NLQueryService chain)
		fx.Invoke(func(*di.Server) {}),
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
