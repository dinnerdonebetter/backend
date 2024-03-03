package config

import (
	"github.com/dinnerdonebetter/backend/internal/observability/logging"

	"github.com/google/wire"
)

var (
	ProvidersLogConfig = wire.NewSet(
		ProvideLogger,
	)
)

func ProvideLogger(cfg *Config) logging.Logger {
	return cfg.ProvideLogger()
}
