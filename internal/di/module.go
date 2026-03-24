package di

import "go.uber.org/fx"

var Module = fx.Options(
	// Core infrastructure
	fx.Provide(NewConfig),
	fx.Provide(NewLogger),
	fx.Provide(NewDatabase),

	// HTTP Layer
	fx.Provide(NewServer),
)
