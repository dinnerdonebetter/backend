package config

import (
	"github.com/google/wire"
)

var (
	// ServiceConfigProviders represents this package's offering to the dependency injector.
	ServiceConfigProviders = wire.NewSet(
		wire.FieldsOf(
			new(*InstanceConfig),
			"Observability",
			"Email",
			"Analytics",
			"FeatureFlags",
			"Encoding",
			"Routing",
			"Database",
			"Meta",
			"Events",
			"Search",
			"Server",
			"Services",
		),
		wire.FieldsOf(
			new(*ServicesConfig),
			"AuditLogEntries",
			"Auth",
			"Households",
			"HouseholdInvitations",
			"ServiceSettings",
			"ServiceSettingConfigurations",
			"Users",
			"UserNotifications",
			"Webhooks",
			"Workers",
			"ValidMeasurementUnits",
			"ValidInstruments",
			"ValidIngredients",
			"ValidPreparations",
			"MealPlanEvents",
			"MealPlanOptionVotes",
			"ValidIngredientPreparations",
			"ValidPreparationInstruments",
			"ValidInstrumentMeasurementUnits",
			"Meals",
			"Recipes",
			"RecipeSteps",
			"RecipeStepProducts",
			"RecipeStepInstruments",
			"RecipeStepIngredients",
			"MealPlans",
			"MealPlanOptions",
			"MealPlanTasks",
			"RecipePrepTasks",
			"MealPlanGroceryListItems",
			"ValidMeasurementUnitConversions",
			"ValidIngredientStates",
			"ValidIngredientGroups",
			"RecipeStepCompletionConditions",
			"ValidIngredientStateIngredients",
			"RecipeStepVessels",
			"UserIngredientPreferences",
			"RecipeRatings",
			"HouseholdInstrumentOwnerships",
			"ValidVessels",
			"ValidPreparationVessels",
		),
	)
)
