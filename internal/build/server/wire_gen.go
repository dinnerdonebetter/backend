// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package server

import (
	"context"
	"github.com/prixfixeco/api_server/internal/authentication"
	"github.com/prixfixeco/api_server/internal/config"
	config3 "github.com/prixfixeco/api_server/internal/customerdata/config"
	"github.com/prixfixeco/api_server/internal/database"
	config2 "github.com/prixfixeco/api_server/internal/database/config"
	"github.com/prixfixeco/api_server/internal/database/queriers/postgres"
	"github.com/prixfixeco/api_server/internal/encoding"
	config4 "github.com/prixfixeco/api_server/internal/messagequeue/config"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/metrics"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/routing/chi"
	"github.com/prixfixeco/api_server/internal/server"
	"github.com/prixfixeco/api_server/internal/services/admin"
	"github.com/prixfixeco/api_server/internal/services/apiclients"
	authentication2 "github.com/prixfixeco/api_server/internal/services/authentication"
	"github.com/prixfixeco/api_server/internal/services/householdinvitations"
	"github.com/prixfixeco/api_server/internal/services/households"
	"github.com/prixfixeco/api_server/internal/services/mealplanoptions"
	"github.com/prixfixeco/api_server/internal/services/mealplanoptionvotes"
	"github.com/prixfixeco/api_server/internal/services/mealplans"
	"github.com/prixfixeco/api_server/internal/services/meals"
	"github.com/prixfixeco/api_server/internal/services/recipes"
	"github.com/prixfixeco/api_server/internal/services/recipestepingredients"
	"github.com/prixfixeco/api_server/internal/services/recipestepinstruments"
	"github.com/prixfixeco/api_server/internal/services/recipestepproducts"
	"github.com/prixfixeco/api_server/internal/services/recipesteps"
	"github.com/prixfixeco/api_server/internal/services/users"
	"github.com/prixfixeco/api_server/internal/services/validingredientpreparations"
	"github.com/prixfixeco/api_server/internal/services/validingredients"
	"github.com/prixfixeco/api_server/internal/services/validinstruments"
	"github.com/prixfixeco/api_server/internal/services/validpreparations"
	"github.com/prixfixeco/api_server/internal/services/webhooks"
	"github.com/prixfixeco/api_server/internal/services/websockets"
	"github.com/prixfixeco/api_server/internal/storage"
	"github.com/prixfixeco/api_server/internal/uploads"
	"github.com/prixfixeco/api_server/internal/uploads/images"
)

// Injectors from build.go:

// Build builds a server.
func Build(ctx context.Context, logger logging.Logger, cfg *config.InstanceConfig, tracerProvider tracing.TracerProvider, unitCounterProvider metrics.UnitCounterProvider, metricsHandler metrics.Handler) (*server.HTTPServer, error) {
	serverConfig := cfg.Server
	servicesConfigurations := &cfg.Services
	authenticationConfig := &servicesConfigurations.Auth
	authenticator := authentication.ProvideArgon2Authenticator(logger, tracerProvider)
	configConfig := &cfg.Database
	dataManager, err := postgres.ProvideDatabaseClient(ctx, logger, configConfig, tracerProvider)
	if err != nil {
		return nil, err
	}
	userDataManager := database.ProvideUserDataManager(dataManager)
	apiClientDataManager := database.ProvideAPIClientDataManager(dataManager)
	householdUserMembershipDataManager := database.ProvideHouseholdUserMembershipDataManager(dataManager)
	cookieConfig := authenticationConfig.Cookies
	sessionManager, err := config2.ProvideSessionManager(cookieConfig, dataManager)
	if err != nil {
		return nil, err
	}
	encodingConfig := cfg.Encoding
	contentType := encoding.ProvideContentType(encodingConfig)
	serverEncoderDecoder := encoding.ProvideServerEncoderDecoder(logger, tracerProvider, contentType)
	config5 := &cfg.CustomerData
	collector, err := config3.ProvideCollector(config5, logger)
	if err != nil {
		return nil, err
	}
	authService, err := authentication2.ProvideService(logger, authenticationConfig, authenticator, userDataManager, apiClientDataManager, householdUserMembershipDataManager, sessionManager, serverEncoderDecoder, collector, tracerProvider)
	if err != nil {
		return nil, err
	}
	householdDataManager := database.ProvideHouseholdDataManager(dataManager)
	imageUploadProcessor := images.NewImageUploadProcessor(logger, tracerProvider)
	uploadsConfig := &cfg.Uploads
	storageConfig := &uploadsConfig.Storage
	routeParamManager := chi.NewRouteParamManager()
	uploader, err := storage.NewUploadManager(ctx, logger, tracerProvider, storageConfig, routeParamManager)
	if err != nil {
		return nil, err
	}
	uploadManager := uploads.ProvideUploadManager(uploader)
	userDataService := users.ProvideUsersService(authenticationConfig, logger, userDataManager, householdDataManager, authenticator, serverEncoderDecoder, unitCounterProvider, imageUploadProcessor, uploadManager, routeParamManager, collector, tracerProvider)
	householdsConfig := servicesConfigurations.Households
	householdInvitationDataManager := database.ProvideHouseholdInvitationDataManager(dataManager)
	config6 := &cfg.Events
	publisherProvider, err := config4.ProvidePublisherProvider(logger, tracerProvider, config6)
	if err != nil {
		return nil, err
	}
	householdDataService, err := households.ProvideService(logger, householdsConfig, householdDataManager, householdInvitationDataManager, householdUserMembershipDataManager, serverEncoderDecoder, unitCounterProvider, routeParamManager, publisherProvider, collector, tracerProvider)
	if err != nil {
		return nil, err
	}
	householdinvitationsConfig := &servicesConfigurations.HouseholdInvitations
	householdInvitationDataService, err := householdinvitations.ProvideHouseholdInvitationsService(logger, householdinvitationsConfig, userDataManager, householdInvitationDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, collector, tracerProvider)
	if err != nil {
		return nil, err
	}
	apiclientsConfig := apiclients.ProvideConfig(authenticationConfig)
	apiClientDataService := apiclients.ProvideAPIClientsService(logger, apiClientDataManager, userDataManager, authenticator, serverEncoderDecoder, unitCounterProvider, routeParamManager, apiclientsConfig, collector, tracerProvider)
	websocketsConfig := servicesConfigurations.Websockets
	consumerProvider, err := config4.ProvideConsumerProvider(logger, tracerProvider, config6)
	if err != nil {
		return nil, err
	}
	websocketDataService, err := websockets.ProvideService(ctx, authenticationConfig, websocketsConfig, logger, serverEncoderDecoder, consumerProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	validinstrumentsConfig := &servicesConfigurations.ValidInstruments
	validInstrumentDataManager := database.ProvideValidInstrumentDataManager(dataManager)
	validInstrumentDataService, err := validinstruments.ProvideService(ctx, logger, validinstrumentsConfig, validInstrumentDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	validingredientsConfig := &servicesConfigurations.ValidIngredients
	validIngredientDataManager := database.ProvideValidIngredientDataManager(dataManager)
	validIngredientDataService, err := validingredients.ProvideService(ctx, logger, validingredientsConfig, validIngredientDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	validpreparationsConfig := &servicesConfigurations.ValidPreparations
	validPreparationDataManager := database.ProvideValidPreparationDataManager(dataManager)
	validPreparationDataService, err := validpreparations.ProvideService(ctx, logger, validpreparationsConfig, validPreparationDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	validingredientpreparationsConfig := &servicesConfigurations.ValidIngredientPreparations
	validIngredientPreparationDataManager := database.ProvideValidIngredientPreparationDataManager(dataManager)
	validIngredientPreparationDataService, err := validingredientpreparations.ProvideService(ctx, logger, validingredientpreparationsConfig, validIngredientPreparationDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	mealsConfig := &servicesConfigurations.Meals
	mealDataManager := database.ProvideMealDataManager(dataManager)
	mealDataService, err := meals.ProvideService(ctx, logger, mealsConfig, mealDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, collector, tracerProvider)
	if err != nil {
		return nil, err
	}
	recipesConfig := &servicesConfigurations.Recipes
	recipeDataManager := database.ProvideRecipeDataManager(dataManager)
	recipeDataService, err := recipes.ProvideService(ctx, logger, recipesConfig, recipeDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, collector, tracerProvider)
	if err != nil {
		return nil, err
	}
	recipestepsConfig := &servicesConfigurations.RecipeSteps
	recipeStepDataManager := database.ProvideRecipeStepDataManager(dataManager)
	recipeStepDataService, err := recipesteps.ProvideService(ctx, logger, recipestepsConfig, recipeStepDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	recipestepinstrumentsConfig := &servicesConfigurations.RecipeStepInstruments
	recipeStepInstrumentDataManager := database.ProvideRecipeStepInstrumentDataManager(dataManager)
	recipeStepInstrumentDataService, err := recipestepinstruments.ProvideService(ctx, logger, recipestepinstrumentsConfig, recipeStepInstrumentDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	recipestepingredientsConfig := &servicesConfigurations.RecipeStepIngredients
	recipeStepIngredientDataManager := database.ProvideRecipeStepIngredientDataManager(dataManager)
	recipeStepIngredientDataService, err := recipestepingredients.ProvideService(ctx, logger, recipestepingredientsConfig, recipeStepIngredientDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	recipestepproductsConfig := &servicesConfigurations.RecipeStepProducts
	recipeStepProductDataManager := database.ProvideRecipeStepProductDataManager(dataManager)
	recipeStepProductDataService, err := recipestepproducts.ProvideService(ctx, logger, recipestepproductsConfig, recipeStepProductDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	mealplansConfig := &servicesConfigurations.MealPlans
	mealPlanDataManager := database.ProvideMealPlanDataManager(dataManager)
	mealPlanDataService, err := mealplans.ProvideService(ctx, logger, mealplansConfig, mealPlanDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, collector, tracerProvider)
	if err != nil {
		return nil, err
	}
	mealplanoptionsConfig := &servicesConfigurations.MealPlanOptions
	mealPlanOptionDataManager := database.ProvideMealPlanOptionDataManager(dataManager)
	mealPlanOptionDataService, err := mealplanoptions.ProvideService(ctx, logger, mealplanoptionsConfig, mealPlanOptionDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	mealplanoptionvotesConfig := &servicesConfigurations.MealPlanOptionVotes
	mealPlanOptionVoteDataService, err := mealplanoptionvotes.ProvideService(ctx, logger, mealplanoptionvotesConfig, dataManager, serverEncoderDecoder, routeParamManager, publisherProvider, collector, tracerProvider)
	if err != nil {
		return nil, err
	}
	webhooksConfig := &servicesConfigurations.Webhooks
	webhookDataManager := database.ProvideWebhookDataManager(dataManager)
	webhookDataService, err := webhooks.ProvideWebhooksService(logger, webhooksConfig, webhookDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	adminUserDataManager := database.ProvideAdminUserDataManager(dataManager)
	adminService := admin.ProvideService(logger, authenticationConfig, authenticator, adminUserDataManager, sessionManager, serverEncoderDecoder, routeParamManager, tracerProvider)
	routingConfig := &cfg.Routing
	router := chi.NewRouter(logger, tracerProvider, routingConfig)
	httpServer, err := server.ProvideHTTPServer(ctx, serverConfig, authService, userDataService, householdDataService, householdInvitationDataService, apiClientDataService, websocketDataService, validInstrumentDataService, validIngredientDataService, validPreparationDataService, validIngredientPreparationDataService, mealDataService, recipeDataService, recipeStepDataService, recipeStepInstrumentDataService, recipeStepIngredientDataService, recipeStepProductDataService, mealPlanDataService, mealPlanOptionDataService, mealPlanOptionVoteDataService, webhookDataService, adminService, logger, serverEncoderDecoder, router, tracerProvider, metricsHandler)
	if err != nil {
		return nil, err
	}
	return httpServer, nil
}
