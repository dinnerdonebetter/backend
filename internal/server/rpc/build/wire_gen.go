// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package build

import (
	"context"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	config2 "github.com/dinnerdonebetter/backend/internal/observability/logging/config"
	config3 "github.com/dinnerdonebetter/backend/internal/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/server/rpc"
)

// Injectors from build.go:

// Build builds a server.
func Build(ctx context.Context, cfg *config.InstanceConfig) (*rpc.Server, error) {
	observabilityConfig := &cfg.Observability
	configConfig := &observabilityConfig.Tracing
	config4 := &observabilityConfig.Logging
	logger, err := config2.ProvideLogger(config4)
	if err != nil {
		return nil, err
	}
	tracerProvider, err := config3.ProvideTracerProvider(ctx, configConfig, logger)
	if err != nil {
		return nil, err
	}
	config5 := &cfg.Database
	dataManager, err := postgres.ProvideDatabaseClient(ctx, logger, config5, tracerProvider)
	if err != nil {
		return nil, err
	}
	validIngredientDataManager := database.ProvideValidIngredientDataManager(dataManager)
	server, err := rpc.ProvideRPCServer(tracerProvider, logger, validIngredientDataManager)
	if err != nil {
		return nil, err
	}
	return server, nil
}
