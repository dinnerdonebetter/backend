// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package server

import (
	"context"
	"github.com/prixfixeco/backend/internal/authentication"
	"github.com/prixfixeco/backend/internal/config"
	"github.com/prixfixeco/backend/internal/database"
	config2 "github.com/prixfixeco/backend/internal/database/config"
	"github.com/prixfixeco/backend/internal/email"
	"github.com/prixfixeco/backend/internal/encoding"
	"github.com/prixfixeco/backend/internal/features/recipeanalysis"
	config3 "github.com/prixfixeco/backend/internal/messagequeue/config"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/metrics"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/internal/random"
	"github.com/prixfixeco/backend/internal/routing/chi"
	"github.com/prixfixeco/backend/internal/server"
	"github.com/prixfixeco/backend/internal/services/admin"
	"github.com/prixfixeco/backend/internal/services/apiclients"
	authentication2 "github.com/prixfixeco/backend/internal/services/authentication"
	"github.com/prixfixeco/backend/internal/services/householdinvitations"
	"github.com/prixfixeco/backend/internal/services/households"
	"github.com/prixfixeco/backend/internal/services/mealplanevents"
	"github.com/prixfixeco/backend/internal/services/mealplangrocerylistitems"
	"github.com/prixfixeco/backend/internal/services/mealplanoptions"
	"github.com/prixfixeco/backend/internal/services/mealplanoptionvotes"
	"github.com/prixfixeco/backend/internal/services/mealplans"
	"github.com/prixfixeco/backend/internal/services/mealplantasks"
	"github.com/prixfixeco/backend/internal/services/meals"
	"github.com/prixfixeco/backend/internal/services/recipepreptasks"
	"github.com/prixfixeco/backend/internal/services/recipes"
	recipestepingredients2 "github.com/prixfixeco/backend/internal/services/recipestepcompletionconditions"
	"github.com/prixfixeco/backend/internal/services/recipestepingredients"
	"github.com/prixfixeco/backend/internal/services/recipestepinstruments"
	"github.com/prixfixeco/backend/internal/services/recipestepproducts"
	"github.com/prixfixeco/backend/internal/services/recipesteps"
	"github.com/prixfixeco/backend/internal/services/recipestepvessels"
	"github.com/prixfixeco/backend/internal/services/users"
	"github.com/prixfixeco/backend/internal/services/validingredientmeasurementunits"
	"github.com/prixfixeco/backend/internal/services/validingredientpreparations"
	"github.com/prixfixeco/backend/internal/services/validingredients"
	"github.com/prixfixeco/backend/internal/services/validingredientstateingredients"
	"github.com/prixfixeco/backend/internal/services/validingredientstates"
	"github.com/prixfixeco/backend/internal/services/validinstruments"
	"github.com/prixfixeco/backend/internal/services/validmeasurementconversions"
	"github.com/prixfixeco/backend/internal/services/validmeasurementunits"
	"github.com/prixfixeco/backend/internal/services/validpreparationinstruments"
	"github.com/prixfixeco/backend/internal/services/validpreparations"
	"github.com/prixfixeco/backend/internal/services/webhooks"
	"github.com/prixfixeco/backend/internal/uploads/images"
)

// Injectors from build.go:

// Build builds a server.
func Build(ctx context.Context, logger logging.Logger, cfg *config.InstanceConfig, tracerProvider tracing.TracerProvider, unitCounterProvider metrics.UnitCounterProvider, metricsHandler metrics.Handler, dataManager database.DataManager, emailer email.Emailer) (*server.HTTPServer, error) {
	serverConfig := cfg.Server
	servicesConfigurations := &cfg.Services
	authenticationConfig := &servicesConfigurations.Auth
	authenticator := authentication.ProvideArgon2Authenticator(logger, tracerProvider)
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
	configConfig := &cfg.Events
	publisherProvider, err := config3.ProvidePublisherProvider(logger, tracerProvider, configConfig)
	if err != nil {
		return nil, err
	}
	generator := random.NewGenerator(logger, tracerProvider)
	authService, err := authentication2.ProvideService(logger, authenticationConfig, authenticator, userDataManager, apiClientDataManager, householdUserMembershipDataManager, sessionManager, serverEncoderDecoder, tracerProvider, publisherProvider, generator, emailer)
	if err != nil {
		return nil, err
	}
	usersConfig := &servicesConfigurations.Users
	householdDataManager := database.ProvideHouseholdDataManager(dataManager)
	householdInvitationDataManager := database.ProvideHouseholdInvitationDataManager(dataManager)
	mediaUploadProcessor := images.NewImageUploadProcessor(logger, tracerProvider)
	routeParamManager := chi.NewRouteParamManager()
	passwordResetTokenDataManager := database.ProvidePasswordResetTokenDataManager(dataManager)
	userDataService, err := users.ProvideUsersService(ctx, usersConfig, authenticationConfig, logger, userDataManager, householdDataManager, householdInvitationDataManager, householdUserMembershipDataManager, authenticator, serverEncoderDecoder, mediaUploadProcessor, routeParamManager, tracerProvider, publisherProvider, generator, passwordResetTokenDataManager, emailer)
	if err != nil {
		return nil, err
	}
	householdsConfig := servicesConfigurations.Households
	householdDataService, err := households.ProvideService(logger, householdsConfig, householdDataManager, householdInvitationDataManager, householdUserMembershipDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	householdinvitationsConfig := &servicesConfigurations.HouseholdInvitations
	householdInvitationDataService, err := householdinvitations.ProvideHouseholdInvitationsService(logger, householdinvitationsConfig, userDataManager, householdInvitationDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider, emailer, generator)
	if err != nil {
		return nil, err
	}
	apiclientsConfig := apiclients.ProvideConfig(authenticationConfig)
	apiClientDataService := apiclients.ProvideAPIClientsService(logger, apiClientDataManager, userDataManager, authenticator, serverEncoderDecoder, routeParamManager, apiclientsConfig, tracerProvider, generator)
	validinstrumentsConfig := &servicesConfigurations.ValidInstruments
	validInstrumentDataManager := database.ProvideValidInstrumentDataManager(dataManager)
	validInstrumentDataService, err := validinstruments.ProvideService(logger, validinstrumentsConfig, validInstrumentDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	validingredientsConfig := &servicesConfigurations.ValidIngredients
	validIngredientDataManager := database.ProvideValidIngredientDataManager(dataManager)
	validIngredientDataService, err := validingredients.ProvideService(logger, validingredientsConfig, validIngredientDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
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
	validIngredientPreparationDataService, err := validingredientpreparations.ProvideService(logger, validingredientpreparationsConfig, validIngredientPreparationDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	mealsConfig := &servicesConfigurations.Meals
	mealDataManager := database.ProvideMealDataManager(dataManager)
	mealDataService, err := meals.ProvideService(logger, mealsConfig, mealDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	recipesConfig := &servicesConfigurations.Recipes
	recipeDataManager := database.ProvideRecipeDataManager(dataManager)
	recipeMediaDataManager := database.ProvideRecipeMediaDataManager(dataManager)
	recipeAnalyzer := recipeanalysis.NewRecipeAnalyzer(logger, tracerProvider)
	recipeDataService, err := recipes.ProvideService(ctx, logger, recipesConfig, recipeDataManager, recipeMediaDataManager, recipeAnalyzer, serverEncoderDecoder, routeParamManager, publisherProvider, mediaUploadProcessor, tracerProvider)
	if err != nil {
		return nil, err
	}
	recipestepsConfig := &servicesConfigurations.RecipeSteps
	recipeStepDataManager := database.ProvideRecipeStepDataManager(dataManager)
	recipeStepDataService, err := recipesteps.ProvideService(ctx, logger, recipestepsConfig, recipeStepDataManager, recipeMediaDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider, mediaUploadProcessor)
	if err != nil {
		return nil, err
	}
	recipestepproductsConfig := &servicesConfigurations.RecipeStepProducts
	recipeStepProductDataManager := database.ProvideRecipeStepProductDataManager(dataManager)
	recipeStepProductDataService, err := recipestepproducts.ProvideService(logger, recipestepproductsConfig, recipeStepProductDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	recipestepinstrumentsConfig := &servicesConfigurations.RecipeStepInstruments
	recipeStepInstrumentDataManager := database.ProvideRecipeStepInstrumentDataManager(dataManager)
	recipeStepInstrumentDataService, err := recipestepinstruments.ProvideService(logger, recipestepinstrumentsConfig, recipeStepInstrumentDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	recipestepingredientsConfig := &servicesConfigurations.RecipeStepIngredients
	recipeStepIngredientDataManager := database.ProvideRecipeStepIngredientDataManager(dataManager)
	recipeStepIngredientDataService, err := recipestepingredients.ProvideService(logger, recipestepingredientsConfig, recipeStepIngredientDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	mealplansConfig := &servicesConfigurations.MealPlans
	mealPlanDataManager := database.ProvideMealPlanDataManager(dataManager)
	mealPlanDataService, err := mealplans.ProvideService(logger, mealplansConfig, mealPlanDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	mealplanoptionsConfig := &servicesConfigurations.MealPlanOptions
	mealPlanOptionDataManager := database.ProvideMealPlanOptionDataManager(dataManager)
	mealPlanOptionDataService, err := mealplanoptions.ProvideService(logger, mealplanoptionsConfig, mealPlanOptionDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	mealplanoptionvotesConfig := &servicesConfigurations.MealPlanOptionVotes
	mealPlanOptionVoteDataService, err := mealplanoptionvotes.ProvideService(logger, mealplanoptionvotesConfig, dataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	validmeasurementunitsConfig := &servicesConfigurations.ValidMeasurementUnits
	validMeasurementUnitDataManager := database.ProvideValidMeasurementUnitDataManager(dataManager)
	validMeasurementUnitDataService, err := validmeasurementunits.ProvideService(logger, validmeasurementunitsConfig, validMeasurementUnitDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	validingredientstatesConfig := &servicesConfigurations.ValidIngredientStates
	validIngredientStateDataManager := database.ProvideValidIngredientStateDataManager(dataManager)
	validIngredientStateDataService, err := validingredientstates.ProvideService(ctx, logger, validingredientstatesConfig, validIngredientStateDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	validpreparationinstrumentsConfig := &servicesConfigurations.ValidPreparationInstruments
	validPreparationInstrumentDataManager := database.ProvideValidPreparationInstrumentDataManager(dataManager)
	validPreparationInstrumentDataService, err := validpreparationinstruments.ProvideService(logger, validpreparationinstrumentsConfig, validPreparationInstrumentDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	validingredientmeasurementunitsConfig := &servicesConfigurations.ValidInstrumentMeasurementUnits
	validIngredientMeasurementUnitDataManager := database.ProvideValidIngredientMeasurementUnitDataManager(dataManager)
	validIngredientMeasurementUnitDataService, err := validingredientmeasurementunits.ProvideService(logger, validingredientmeasurementunitsConfig, validIngredientMeasurementUnitDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	mealplaneventsConfig := &servicesConfigurations.MealPlanEvents
	mealPlanEventDataManager := database.ProvideMealPlanEventDataManager(dataManager)
	mealPlanEventDataService, err := mealplanevents.ProvideService(logger, mealplaneventsConfig, mealPlanEventDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	mealplantasksConfig := &servicesConfigurations.MealPlanTasks
	mealPlanTaskDataManager := database.ProvideMealPlanTaskDataManager(dataManager)
	mealPlanTaskDataService, err := mealplantasks.ProvideService(logger, mealplantasksConfig, mealPlanTaskDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	recipepreptasksConfig := &servicesConfigurations.RecipePrepTasks
	recipePrepTaskDataManager := database.ProvideRecipePrepTaskDataManager(dataManager)
	recipePrepTaskDataService, err := recipepreptasks.ProvideService(logger, recipepreptasksConfig, recipePrepTaskDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	mealplangrocerylistitemsConfig := &servicesConfigurations.MealPlanGroceryListItems
	mealPlanGroceryListItemDataManager := database.ProvideMealPlanGroceryListItemDataManager(dataManager)
	mealPlanGroceryListItemDataService, err := mealplangrocerylistitems.ProvideService(logger, mealplangrocerylistitemsConfig, mealPlanGroceryListItemDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	validmeasurementconversionsConfig := &servicesConfigurations.ValidMeasurementConversions
	validMeasurementConversionDataManager := database.ProvideValidMeasurementConversionDataManager(dataManager)
	validMeasurementConversionDataService, err := validmeasurementconversions.ProvideService(ctx, logger, validmeasurementconversionsConfig, validMeasurementConversionDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	config4 := &servicesConfigurations.RecipeStepCompletionConditions
	recipeStepCompletionConditionDataManager := database.ProvideRecipeStepCompletionConditionDataManager(dataManager)
	recipeStepCompletionConditionDataService, err := recipestepingredients2.ProvideService(logger, config4, recipeStepCompletionConditionDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	validingredientstateingredientsConfig := &servicesConfigurations.ValidIngredientStateIngredients
	validIngredientStateIngredientDataManager := database.ProvideValidIngredientStateIngredientDataManager(dataManager)
	validIngredientStateIngredientDataService, err := validingredientstateingredients.ProvideService(logger, validingredientstateingredientsConfig, validIngredientStateIngredientDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	recipestepvesselsConfig := &servicesConfigurations.RecipeStepVessels
	recipeStepVesselDataManager := database.ProvideRecipeStepVesselDataManager(dataManager)
	recipeStepVesselDataService, err := recipestepvessels.ProvideService(logger, recipestepvesselsConfig, recipeStepVesselDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
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
	httpServer, err := server.ProvideHTTPServer(ctx, serverConfig, authService, userDataService, householdDataService, householdInvitationDataService, apiClientDataService, validInstrumentDataService, validIngredientDataService, validPreparationDataService, validIngredientPreparationDataService, mealDataService, recipeDataService, recipeStepDataService, recipeStepProductDataService, recipeStepInstrumentDataService, recipeStepIngredientDataService, mealPlanDataService, mealPlanOptionDataService, mealPlanOptionVoteDataService, validMeasurementUnitDataService, validIngredientStateDataService, validPreparationInstrumentDataService, validIngredientMeasurementUnitDataService, mealPlanEventDataService, mealPlanTaskDataService, recipePrepTaskDataService, mealPlanGroceryListItemDataService, validMeasurementConversionDataService, recipeStepCompletionConditionDataService, validIngredientStateIngredientDataService, recipeStepVesselDataService, webhookDataService, adminService, logger, serverEncoderDecoder, router, tracerProvider, metricsHandler)
	if err != nil {
		return nil, err
	}
	return httpServer, nil
}
