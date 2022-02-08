package config

import (
	"bytes"
	"context"
	"encoding/json"
	"os"

	"github.com/prixfixeco/api_server/internal/database"

	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
)

// GetConfigFromCloudSecretManager fetches and InstanceConfig from AWS SSM Parameter Store.
func GetConfigFromCloudSecretManager(ctx context.Context) (*InstanceConfig, error) {
	logger := zerolog.NewZerologLogger()
	logger.Debug("setting up secret manager client")

	var cfg *InstanceConfig
	configFilepath := os.Getenv("CONFIGURATION_FILEPATH")

	configBytes, configReadErr := os.ReadFile(configFilepath)
	if configReadErr != nil {
		return nil, configReadErr
	}

	if encodeErr := json.NewDecoder(bytes.NewReader(configBytes)).Decode(&cfg); encodeErr != nil || cfg == nil {
		return nil, encodeErr
	}

	// fetch supplementary data from env vars
	cfg.Database.ConnectionDetails = database.ConnectionDetails(os.Getenv("PRIXFIXE_DATABASE_URL"))

	cfg.Services.Auth.Cookies.HashKey = os.Getenv("PRIXFIXE_COOKIE_HASH_KEY")
	cfg.Services.Auth.Cookies.BlockKey = os.Getenv("PRIXFIXE_COOKIE_BLOCK_KEY")
	cfg.Services.Auth.PASETO.LocalModeKey = []byte(os.Getenv("PRIXFIXE_PASETO_LOCAL_KEY"))

	dataChangesTopicName := os.Getenv("PRIXFIXE_DATA_CHANGES_TOPIC")

	cfg.Email.APIToken = os.Getenv("PRIXFIXE_SENDGRID_API_TOKEN")
	cfg.CustomerData.APIToken = os.Getenv("PRIXFIXE_SEGMENT_API_TOKEN")

	cfg.Services.ValidInstruments.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidIngredients.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidPreparations.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidIngredientPreparations.DataChangesTopicName = dataChangesTopicName

	cfg.Services.Recipes.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipeSteps.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipeStepProducts.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipeStepInstruments.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipeStepIngredients.DataChangesTopicName = dataChangesTopicName

	cfg.Services.Meals.DataChangesTopicName = dataChangesTopicName
	cfg.Services.MealPlans.DataChangesTopicName = dataChangesTopicName
	cfg.Services.MealPlanOptions.DataChangesTopicName = dataChangesTopicName
	cfg.Services.MealPlanOptionVotes.DataChangesTopicName = dataChangesTopicName

	cfg.Services.Households.DataChangesTopicName = dataChangesTopicName
	cfg.Services.HouseholdInvitations.DataChangesTopicName = dataChangesTopicName
	cfg.Services.Webhooks.DataChangesTopicName = dataChangesTopicName

	if err := cfg.ValidateWithContext(ctx); err != nil {
		return nil, err
	}

	return cfg, nil
}