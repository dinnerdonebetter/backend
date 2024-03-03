package config

import (
	"context"
	"fmt"

	auditlogentriesservice "github.com/dinnerdonebetter/backend/internal/services/auditlogentries"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	householdinstrumentownershipsservice "github.com/dinnerdonebetter/backend/internal/services/householdinstrumentownerships"
	householdinvitationsservice "github.com/dinnerdonebetter/backend/internal/services/householdinvitations"
	householdsservice "github.com/dinnerdonebetter/backend/internal/services/households"
	mealplaneventsservice "github.com/dinnerdonebetter/backend/internal/services/mealplanevents"
	"github.com/dinnerdonebetter/backend/internal/services/mealplangrocerylistitems"
	mealplanoptionsservice "github.com/dinnerdonebetter/backend/internal/services/mealplanoptions"
	mealplanoptionvotesservice "github.com/dinnerdonebetter/backend/internal/services/mealplanoptionvotes"
	mealplansservice "github.com/dinnerdonebetter/backend/internal/services/mealplans"
	"github.com/dinnerdonebetter/backend/internal/services/mealplantasks"
	mealsservice "github.com/dinnerdonebetter/backend/internal/services/meals"
	oauth2clientsservice "github.com/dinnerdonebetter/backend/internal/services/oauth2clients"
	"github.com/dinnerdonebetter/backend/internal/services/recipepreptasks"
	reciperatingsservice "github.com/dinnerdonebetter/backend/internal/services/reciperatings"
	recipesservice "github.com/dinnerdonebetter/backend/internal/services/recipes"
	recipestepcompletionconditionsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepcompletionconditions"
	recipestepingredientsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepingredients"
	recipestepinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepinstruments"
	recipestepproductsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepproducts"
	recipestepsservice "github.com/dinnerdonebetter/backend/internal/services/recipesteps"
	recipestepvesselsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepvessels"
	servicesettingconfigurationsservice "github.com/dinnerdonebetter/backend/internal/services/servicesettingconfigurations"
	servicesettingsservice "github.com/dinnerdonebetter/backend/internal/services/servicesettings"
	useringredientpreferencesservice "github.com/dinnerdonebetter/backend/internal/services/useringredientpreferences"
	usernotificationsservice "github.com/dinnerdonebetter/backend/internal/services/usernotifications"
	usersservice "github.com/dinnerdonebetter/backend/internal/services/users"
	validingredientgroupsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientgroups"
	validingredientmeasurementunitsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientmeasurementunits"
	validingredientpreparationsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientpreparations"
	validingredientsservice "github.com/dinnerdonebetter/backend/internal/services/validingredients"
	"github.com/dinnerdonebetter/backend/internal/services/validingredientstateingredients"
	"github.com/dinnerdonebetter/backend/internal/services/validingredientstates"
	validinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/validinstruments"
	"github.com/dinnerdonebetter/backend/internal/services/validmeasurementunitconversions"
	validmeasurementunitsservice "github.com/dinnerdonebetter/backend/internal/services/validmeasurementunits"
	validpreparationinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/validpreparationinstruments"
	validpreparationsservice "github.com/dinnerdonebetter/backend/internal/services/validpreparations"
	"github.com/dinnerdonebetter/backend/internal/services/validpreparationvessels"
	validvesselsservice "github.com/dinnerdonebetter/backend/internal/services/validvessels"
	webhooksservice "github.com/dinnerdonebetter/backend/internal/services/webhooks"
	workersservice "github.com/dinnerdonebetter/backend/internal/services/workers"

	"github.com/hashicorp/go-multierror"
)

type (
	// ServicesConfig collects the various service configurations.
	ServicesConfig struct {
		_                               struct{}                                      `json:"-"`
		AuditLogEntries                 auditlogentriesservice.Config                 `json:"auditLogEntries"                 toml:"audit_log_entries,omitempty"`
		RecipeStepProducts              recipestepproductsservice.Config              `json:"recipeStepProducts"              toml:"recipe_step_products,omitempty"`
		ValidInstrumentMeasurementUnits validingredientmeasurementunitsservice.Config `json:"validInstrumentMeasurementUnits" toml:"valid_ingredient_measurement_units,omitempty"`
		RecipeRatings                   reciperatingsservice.Config                   `json:"recipeRatings"                   toml:"recipe_ratings,omitempty"`
		MealPlanGroceryListItems        mealplangrocerylistitems.Config               `json:"mealPlanGroceryListItems"        toml:"meal_plan_grocery_list_items,omitempty"`
		ValidMeasurementUnitConversions validmeasurementunitconversions.Config        `json:"validMeasurementUnitConversions" toml:"valid_measurement_conversions,omitempty"`
		ServiceSettingConfigurations    servicesettingconfigurationsservice.Config    `json:"serviceSettingConfigurations"    toml:"service_setting_configurations,omitempty"`
		ServiceSettings                 servicesettingsservice.Config                 `json:"serviceSettings"                 toml:"service_settings,omitempty"`
		ValidIngredientStateIngredients validingredientstateingredients.Config        `json:"validIngredientStateIngredients" toml:"valid_ingredient_state_ingredients,omitempty"`
		RecipeStepInstruments           recipestepinstrumentsservice.Config           `json:"recipeStepInstruments"           toml:"recipe_step_instruments,omitempty"`
		RecipeStepIngredients           recipestepingredientsservice.Config           `json:"recipeStepIngredients"           toml:"recipe_step_ingredients,omitempty"`
		HouseholdInstrumentOwnerships   householdinstrumentownershipsservice.Config   `json:"householdInstrumentOwnerships"   toml:"household_instrument_ownerships,omitempty"`
		RecipePrepTasks                 recipepreptasks.Config                        `json:"recipePrepTasks"                 toml:"recipe_prep_tasks,omitempty"`
		MealPlanEvents                  mealplaneventsservice.Config                  `json:"mealPlanEvents"                  toml:"meal_plan_events,omitempty"`
		UserIngredientPreferences       useringredientpreferencesservice.Config       `json:"userIngredientPreferences"       toml:"user_ingredient_preferences,omitempty"`
		Households                      householdsservice.Config                      `json:"households"                      toml:"households,omitempty"`
		MealPlans                       mealplansservice.Config                       `json:"mealPlans"                       toml:"meal_plans,omitempty"`
		RecipeStepVessels               recipestepvesselsservice.Config               `json:"recipeStepVessels"               toml:"recipe_step_vessels,omitempty"`
		ValidIngredientPreparations     validingredientpreparationsservice.Config     `json:"validIngredientPreparations"     toml:"valid_ingredient_preparations,omitempty"`
		MealPlanTasks                   mealplantasks.Config                          `json:"mealPlanTasks"                   toml:"meal_plan_tasks,omitempty"`
		MealPlanOptionVotes             mealplanoptionvotesservice.Config             `json:"mealPlanOptionVotes"             toml:"meal_plan_option_votes,omitempty"`
		ValidPreparationInstruments     validpreparationinstrumentsservice.Config     `json:"validPreparationInstruments"     toml:"valid_preparation_instruments,omitempty"`
		RecipeStepCompletionConditions  recipestepcompletionconditionsservice.Config  `json:"recipeStepCompletionConditions"  toml:"recipe_step_completion_conditions,omitempty"`
		ValidIngredientGroups           validingredientgroupsservice.Config           `json:"validIngredientGroups"           toml:"valid_ingredient_groups,omitempty"`
		ValidPreparationVessels         validpreparationvessels.Config                `json:"validPreparationVessels"         toml:"valid_preparation_vessels,omitempty"`
		Workers                         workersservice.Config                         `json:"workers"                         toml:"workers,omitempty"`
		UserNotifications               usernotificationsservice.Config               `json:"userNotifications"               toml:"user_notifications,omitempty"`
		MealPlanOptions                 mealplanoptionsservice.Config                 `json:"mealPlanOptions"                 toml:"meal_plan_options,omitempty"`
		Users                           usersservice.Config                           `json:"users"                           toml:"users,omitempty"`
		RecipeSteps                     recipestepsservice.Config                     `json:"recipeSteps"                     toml:"recipe_steps,omitempty"`
		ValidPreparations               validpreparationsservice.Config               `json:"validPreparations"               toml:"valid_preparations,omitempty"`
		ValidIngredients                validingredientsservice.Config                `json:"validIngredients"                toml:"valid_ingredients,omitempty"`
		ValidMeasurementUnits           validmeasurementunitsservice.Config           `json:"validMeasurementUnits"           toml:"valid_measurement_units,omitempty"`
		Meals                           mealsservice.Config                           `json:"meals"                           toml:"meals,omitempty"`
		OAuth2Clients                   oauth2clientsservice.Config                   `json:"oauth2Clients"                   toml:"oauth2_clients,omitempty"`
		ValidIngredientStates           validingredientstates.Config                  `json:"validIngredientStates"           toml:"valid_ingredient_states,omitempty"`
		Webhooks                        webhooksservice.Config                        `json:"webhooks"                        toml:"webhooks,omitempty"`
		ValidInstruments                validinstrumentsservice.Config                `json:"validInstruments"                toml:"valid_instruments,omitempty"`
		ValidVessels                    validvesselsservice.Config                    `json:"validVessels"                    toml:"auth,omitempty"`
		HouseholdInvitations            householdinvitationsservice.Config            `json:"householdInvitations"            toml:"household_invitations,omitempty"`
		Recipes                         recipesservice.Config                         `json:"recipes"                         toml:"recipes,omitempty"`
		Auth                            authservice.Config                            `json:"auth"                            toml:"auth,omitempty"`
	}
)

// ValidateWithContext validates a InstanceConfig struct.
func (cfg *ServicesConfig) ValidateWithContext(ctx context.Context) error {
	var result *multierror.Error

	validatorsToRun := map[string]func(context.Context) error{
		"Auth":                            cfg.Auth.ValidateWithContext,
		"Users":                           cfg.Users.ValidateWithContext,
		"Webhooks":                        cfg.Webhooks.ValidateWithContext,
		"ValidInstruments":                cfg.ValidInstruments.ValidateWithContext,
		"ValidIngredients":                cfg.ValidIngredients.ValidateWithContext,
		"ValidIngredientGroups":           cfg.ValidIngredientGroups.ValidateWithContext,
		"ValidPreparations":               cfg.ValidPreparations.ValidateWithContext,
		"ValidMeasurementUnits":           cfg.ValidMeasurementUnits.ValidateWithContext,
		"ValidIngredientPreparations":     cfg.ValidIngredientPreparations.ValidateWithContext,
		"ValidIngredientStateIngredients": cfg.ValidIngredientStateIngredients.ValidateWithContext,
		"ValidPreparationInstruments":     cfg.ValidPreparationInstruments.ValidateWithContext,
		"ValidInstrumentMeasurementUnits": cfg.ValidInstrumentMeasurementUnits.ValidateWithContext,
		"Recipes":                         cfg.Recipes.ValidateWithContext,
		"RecipeSteps":                     cfg.RecipeSteps.ValidateWithContext,
		"RecipeStepInstruments":           cfg.RecipeStepInstruments.ValidateWithContext,
		"RecipeStepVessels":               cfg.RecipeStepVessels.ValidateWithContext,
		"RecipeStepIngredients":           cfg.RecipeStepIngredients.ValidateWithContext,
		"RecipeStepCompletionConditions":  cfg.RecipeStepCompletionConditions.ValidateWithContext,
		"MealPlans":                       cfg.MealPlans.ValidateWithContext,
		"MealPlanEvents":                  cfg.MealPlanEvents.ValidateWithContext,
		"MealPlanOptions":                 cfg.MealPlanOptions.ValidateWithContext,
		"MealPlanOptionVotes":             cfg.MealPlanOptionVotes.ValidateWithContext,
		"RecipePrepTasks":                 cfg.RecipePrepTasks.ValidateWithContext,
		"MealPlanGroceryListItems":        cfg.MealPlanGroceryListItems.ValidateWithContext,
		"ValidMeasurementUnitConversions": cfg.ValidMeasurementUnitConversions.ValidateWithContext,
		"ValidIngredientStates":           cfg.ValidIngredientStates.ValidateWithContext,
		"ServiceSettings":                 cfg.ServiceSettings.ValidateWithContext,
		"ServiceSettingConfigurations":    cfg.ServiceSettingConfigurations.ValidateWithContext,
		"UserIngredientPreferences":       cfg.UserIngredientPreferences.ValidateWithContext,
		"ValidVessels":                    cfg.ValidVessels.ValidateWithContext,
		"ValidPreparationVessels":         cfg.ValidPreparationVessels.ValidateWithContext,
		"Workers":                         cfg.Workers.ValidateWithContext,
		"UserNotifications":               cfg.UserNotifications.ValidateWithContext,
		"AuditLogEntries":                 cfg.AuditLogEntries.ValidateWithContext,
	}

	for name, validator := range validatorsToRun {
		if err := validator(ctx); err != nil {
			result = multierror.Append(fmt.Errorf("error validating %s config: %w", name, err), result)
		}
	}

	return result.ErrorOrNil()
}
