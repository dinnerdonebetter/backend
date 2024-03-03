package authorization

import (
	"github.com/mikespook/gorbac/v2"
)

var _ gorbac.Permission = (*Permission)(nil)

type (
	// Permission is a simple string alias.
	Permission string
)

const (
	// CycleCookieSecretPermission is a service admin permission.
	CycleCookieSecretPermission Permission = "update.cookie_secret"
	// UpdateUserStatusPermission is a service admin permission.
	UpdateUserStatusPermission Permission = "update.user_status"
	// ReadUserPermission is a service admin permission.
	ReadUserPermission Permission = "read.user"
	// SearchUserPermission is a service admin permission.
	SearchUserPermission Permission = "search.user"

	// UpdateHouseholdPermission is a household admin permission.
	UpdateHouseholdPermission Permission = "update.household"
	// ArchiveHouseholdPermission is a household admin permission.
	ArchiveHouseholdPermission Permission = "archive.household"
	// InviteUserToHouseholdPermission is a household admin permission.
	InviteUserToHouseholdPermission Permission = "household.add.member"
	// ModifyMemberPermissionsForHouseholdPermission is a household admin permission.
	ModifyMemberPermissionsForHouseholdPermission Permission = "household.membership.modify"
	// RemoveMemberHouseholdPermission is a household admin permission.
	RemoveMemberHouseholdPermission Permission = "remove_member.household"
	// TransferHouseholdPermission is a household admin permission.
	TransferHouseholdPermission Permission = "transfer.household"

	// CreateWebhooksPermission is a household admin permission.
	CreateWebhooksPermission Permission = "create.webhooks"
	// ReadWebhooksPermission is a household admin permission.
	ReadWebhooksPermission Permission = "read.webhooks"
	// UpdateWebhooksPermission is a household admin permission.
	UpdateWebhooksPermission Permission = "update.webhooks"
	// ArchiveWebhooksPermission is a household admin permission.
	ArchiveWebhooksPermission Permission = "archive.webhooks"
	// CreateWebhookTriggerEventsPermission is a household admin permission.
	CreateWebhookTriggerEventsPermission Permission = "create.webhook_trigger_events"
	// ArchiveWebhookTriggerEventsPermission is a household admin permission.
	ArchiveWebhookTriggerEventsPermission Permission = "archive.webhook_trigger_events"

	// ReadAuditLogEntriesPermission is a service permission.
	ReadAuditLogEntriesPermission Permission = "read.audit_log_entries"

	// CreateValidInstrumentsPermission is a household user permission.
	CreateValidInstrumentsPermission Permission = "create.valid_instruments"
	// ReadValidInstrumentsPermission is a household user permission.
	ReadValidInstrumentsPermission Permission = "read.valid_instruments"
	// SearchValidInstrumentsPermission is a household user permission.
	SearchValidInstrumentsPermission Permission = "search.valid_instruments"
	// UpdateValidInstrumentsPermission is a household user permission.
	UpdateValidInstrumentsPermission Permission = "update.valid_instruments"
	// ArchiveValidInstrumentsPermission is a household user permission.
	ArchiveValidInstrumentsPermission Permission = "archive.valid_instruments"

	// CreateValidVesselsPermission is a household user permission.
	CreateValidVesselsPermission Permission = "create.valid_vessels"
	// ReadValidVesselsPermission is a household user permission.
	ReadValidVesselsPermission Permission = "read.valid_vessels"
	// SearchValidVesselsPermission is a household user permission.
	SearchValidVesselsPermission Permission = "search.valid_vessels"
	// UpdateValidVesselsPermission is a household user permission.
	UpdateValidVesselsPermission Permission = "update.valid_vessels"
	// ArchiveValidVesselsPermission is a household user permission.
	ArchiveValidVesselsPermission Permission = "archive.valid_vessels"

	// CreateValidIngredientsPermission is a household user permission.
	CreateValidIngredientsPermission Permission = "create.valid_ingredients"
	// ReadValidIngredientsPermission is a household user permission.
	ReadValidIngredientsPermission Permission = "read.valid_ingredients"
	// SearchValidIngredientsPermission is a household user permission.
	SearchValidIngredientsPermission Permission = "search.valid_ingredients"
	// UpdateValidIngredientsPermission is a household user permission.
	UpdateValidIngredientsPermission Permission = "update.valid_ingredients"
	// ArchiveValidIngredientsPermission is a household user permission.
	ArchiveValidIngredientsPermission Permission = "archive.valid_ingredients"

	// CreateValidIngredientGroupsPermission is a household user permission.
	CreateValidIngredientGroupsPermission Permission = "create.valid_ingredient_groups"
	// ReadValidIngredientGroupsPermission is a household user permission.
	ReadValidIngredientGroupsPermission Permission = "read.valid_ingredient_groups"
	// SearchValidIngredientGroupsPermission is a household user permission.
	SearchValidIngredientGroupsPermission Permission = "search.valid_ingredient_groups"
	// UpdateValidIngredientGroupsPermission is a household user permission.
	UpdateValidIngredientGroupsPermission Permission = "update.valid_ingredient_groups"
	// ArchiveValidIngredientGroupsPermission is a household user permission.
	ArchiveValidIngredientGroupsPermission Permission = "archive.valid_ingredient_groups"

	// CreateValidPreparationsPermission is a household user permission.
	CreateValidPreparationsPermission Permission = "create.valid_preparations"
	// ReadValidPreparationsPermission is a household user permission.
	ReadValidPreparationsPermission Permission = "read.valid_preparations"
	// SearchValidPreparationsPermission is a household user permission.
	SearchValidPreparationsPermission Permission = "search.valid_preparations"
	// UpdateValidPreparationsPermission is a household user permission.
	UpdateValidPreparationsPermission Permission = "update.valid_preparations"
	// ArchiveValidPreparationsPermission is a household user permission.
	ArchiveValidPreparationsPermission Permission = "archive.valid_preparations"

	// CreateValidMeasurementUnitsPermission is a household user permission.
	CreateValidMeasurementUnitsPermission Permission = "create.measurement_units"
	// ReadValidMeasurementUnitsPermission is a household user permission.
	ReadValidMeasurementUnitsPermission Permission = "read.measurement_units"
	// SearchValidMeasurementUnitsPermission is a household user permission.
	SearchValidMeasurementUnitsPermission Permission = "search.measurement_units"
	// UpdateValidMeasurementUnitsPermission is a household user permission.
	UpdateValidMeasurementUnitsPermission Permission = "update.measurement_units"
	// ArchiveValidMeasurementUnitsPermission is a household user permission.
	ArchiveValidMeasurementUnitsPermission Permission = "archive.measurement_units"

	// CreateValidIngredientStatesPermission is a household user permission.
	CreateValidIngredientStatesPermission Permission = "create.valid_ingredient_states"
	// ReadValidIngredientStatesPermission is a household user permission.
	ReadValidIngredientStatesPermission Permission = "read.valid_ingredient_states"
	// UpdateValidIngredientStatesPermission is a household user permission.
	UpdateValidIngredientStatesPermission Permission = "update.valid_ingredient_states"
	// ArchiveValidIngredientStatesPermission is a household user permission.
	ArchiveValidIngredientStatesPermission Permission = "archive.valid_ingredient_states"

	// CreateValidMeasurementUnitConversionsPermission is a household user permission.
	CreateValidMeasurementUnitConversionsPermission Permission = "create.measurement_conversions"
	// ReadValidMeasurementUnitConversionsPermission is a household user permission.
	ReadValidMeasurementUnitConversionsPermission Permission = "read.measurement_conversions"
	// UpdateValidMeasurementUnitConversionsPermission is a household user permission.
	UpdateValidMeasurementUnitConversionsPermission Permission = "update.measurement_conversions"
	// ArchiveValidMeasurementUnitConversionsPermission is a household user permission.
	ArchiveValidMeasurementUnitConversionsPermission Permission = "archive.measurement_conversions"

	// CreateValidIngredientPreparationsPermission is a household user permission.
	CreateValidIngredientPreparationsPermission Permission = "create.valid_ingredient_preparations"
	// ReadValidIngredientPreparationsPermission is a household user permission.
	ReadValidIngredientPreparationsPermission Permission = "read.valid_ingredient_preparations"
	// SearchValidIngredientPreparationsPermission is a household user permission.
	SearchValidIngredientPreparationsPermission Permission = "search.valid_ingredient_preparations"
	// UpdateValidIngredientPreparationsPermission is a household user permission.
	UpdateValidIngredientPreparationsPermission Permission = "update.valid_ingredient_preparations"
	// ArchiveValidIngredientPreparationsPermission is a household user permission.
	ArchiveValidIngredientPreparationsPermission Permission = "archive.valid_ingredient_preparations"

	// CreateValidIngredientStateIngredientsPermission is a household user permission.
	CreateValidIngredientStateIngredientsPermission Permission = "create.valid_ingredient_state_ingredients"
	// ReadValidIngredientStateIngredientsPermission is a household user permission.
	ReadValidIngredientStateIngredientsPermission Permission = "read.valid_ingredient_state_ingredients"
	// SearchValidIngredientStateIngredientsPermission is a household user permission.
	SearchValidIngredientStateIngredientsPermission Permission = "search.valid_ingredient_state_ingredients"
	// UpdateValidIngredientStateIngredientsPermission is a household user permission.
	UpdateValidIngredientStateIngredientsPermission Permission = "update.valid_ingredient_state_ingredients"
	// ArchiveValidIngredientStateIngredientsPermission is a household user permission.
	ArchiveValidIngredientStateIngredientsPermission Permission = "archive.valid_ingredient_state_ingredients"

	// CreateValidPreparationInstrumentsPermission is a household user permission.
	CreateValidPreparationInstrumentsPermission Permission = "create.valid_preparation_instruments"
	// ReadValidPreparationInstrumentsPermission is a household user permission.
	ReadValidPreparationInstrumentsPermission Permission = "read.valid_preparation_instruments"
	// SearchValidPreparationInstrumentsPermission is a household user permission.
	SearchValidPreparationInstrumentsPermission Permission = "search.valid_preparation_instruments"
	// UpdateValidPreparationInstrumentsPermission is a household user permission.
	UpdateValidPreparationInstrumentsPermission Permission = "update.valid_preparation_instruments"
	// ArchiveValidPreparationInstrumentsPermission is a household user permission.
	ArchiveValidPreparationInstrumentsPermission Permission = "archive.valid_preparation_instruments"

	// CreateValidPreparationVesselsPermission is a household user permission.
	CreateValidPreparationVesselsPermission Permission = "create.valid_preparation_vessels"
	// ReadValidPreparationVesselsPermission is a household user permission.
	ReadValidPreparationVesselsPermission Permission = "read.valid_preparation_vessels"
	// SearchValidPreparationVesselsPermission is a household user permission.
	SearchValidPreparationVesselsPermission Permission = "search.valid_preparation_vessels"
	// UpdateValidPreparationVesselsPermission is a household user permission.
	UpdateValidPreparationVesselsPermission Permission = "update.valid_preparation_vessels"
	// ArchiveValidPreparationVesselsPermission is a household user permission.
	ArchiveValidPreparationVesselsPermission Permission = "archive.valid_preparation_vessels"

	// CreateValidIngredientMeasurementUnitsPermission is a household user permission.
	CreateValidIngredientMeasurementUnitsPermission Permission = "create.valid_ingredient_measurement_units"
	// ReadValidIngredientMeasurementUnitsPermission is a household user permission.
	ReadValidIngredientMeasurementUnitsPermission Permission = "read.valid_ingredient_measurement_units"
	// SearchValidIngredientMeasurementUnitsPermission is a household user permission.
	SearchValidIngredientMeasurementUnitsPermission Permission = "search.valid_ingredient_measurement_units"
	// UpdateValidIngredientMeasurementUnitsPermission is a household user permission.
	UpdateValidIngredientMeasurementUnitsPermission Permission = "update.valid_ingredient_measurement_units"
	// ArchiveValidIngredientMeasurementUnitsPermission is a household user permission.
	ArchiveValidIngredientMeasurementUnitsPermission Permission = "archive.valid_ingredient_measurement_units"

	// CreateMealsPermission is a household user permission.
	CreateMealsPermission Permission = "create.meals"
	// ReadMealsPermission is a household user permission.
	ReadMealsPermission Permission = "read.meals"
	// UpdateMealsPermission is a household user permission.
	UpdateMealsPermission Permission = "update.meals"
	// ArchiveMealsPermission is a household user permission.
	ArchiveMealsPermission Permission = "archive.meals"

	// CreateRecipesPermission is a household user permission.
	CreateRecipesPermission Permission = "create.recipes"
	// ReadRecipesPermission is a household user permission.
	ReadRecipesPermission Permission = "read.recipes"
	// SearchRecipesPermission is a household user permission.
	SearchRecipesPermission Permission = "search.recipes"
	// UpdateRecipesPermission is a household user permission.
	UpdateRecipesPermission Permission = "update.recipes"
	// ArchiveRecipesPermission is a household user permission.
	ArchiveRecipesPermission Permission = "archive.recipes"

	// CreateRecipePrepTasksPermission is a household user permission.
	CreateRecipePrepTasksPermission Permission = "create.recipe_prep_tasks"
	// ReadRecipePrepTasksPermission is a household user permission.
	ReadRecipePrepTasksPermission Permission = "read.recipe_prep_tasks"
	// UpdateRecipePrepTasksPermission is a household user permission.
	UpdateRecipePrepTasksPermission Permission = "update.recipe_prep_tasks"
	// ArchiveRecipePrepTasksPermission is a household user permission.
	ArchiveRecipePrepTasksPermission Permission = "archive.recipe_prep_tasks"

	// CreateRecipeStepsPermission is a household user permission.
	CreateRecipeStepsPermission Permission = "create.recipe_steps"
	// ReadRecipeStepsPermission is a household user permission.
	ReadRecipeStepsPermission Permission = "read.recipe_steps"
	// SearchRecipeStepsPermission is a household user permission.
	SearchRecipeStepsPermission Permission = "search.recipe_steps"
	// UpdateRecipeStepsPermission is a household user permission.
	UpdateRecipeStepsPermission Permission = "update.recipe_steps"
	// ArchiveRecipeStepsPermission is a household user permission.
	ArchiveRecipeStepsPermission Permission = "archive.recipe_steps"

	// CreateRecipeStepInstrumentsPermission is a household user permission.
	CreateRecipeStepInstrumentsPermission Permission = "create.recipe_step_instruments"
	// ReadRecipeStepInstrumentsPermission is a household user permission.
	ReadRecipeStepInstrumentsPermission Permission = "read.recipe_step_instruments"
	// SearchRecipeStepInstrumentsPermission is a household user permission.
	SearchRecipeStepInstrumentsPermission Permission = "search.recipe_step_instruments"
	// UpdateRecipeStepInstrumentsPermission is a household user permission.
	UpdateRecipeStepInstrumentsPermission Permission = "update.recipe_step_instruments"
	// ArchiveRecipeStepInstrumentsPermission is a household user permission.
	ArchiveRecipeStepInstrumentsPermission Permission = "archive.recipe_step_instruments"

	// CreateRecipeStepVesselsPermission is a household user permission.
	CreateRecipeStepVesselsPermission Permission = "create.recipe_step_vessels"
	// ReadRecipeStepVesselsPermission is a household user permission.
	ReadRecipeStepVesselsPermission Permission = "read.recipe_step_vessels"
	// SearchRecipeStepVesselsPermission is a household user permission.
	SearchRecipeStepVesselsPermission Permission = "search.recipe_step_vessels"
	// UpdateRecipeStepVesselsPermission is a household user permission.
	UpdateRecipeStepVesselsPermission Permission = "update.recipe_step_vessels"
	// ArchiveRecipeStepVesselsPermission is a household user permission.
	ArchiveRecipeStepVesselsPermission Permission = "archive.recipe_step_vessels"

	// CreateRecipeStepIngredientsPermission is a household user permission.
	CreateRecipeStepIngredientsPermission Permission = "create.recipe_step_ingredients"
	// ReadRecipeStepIngredientsPermission is a household user permission.
	ReadRecipeStepIngredientsPermission Permission = "read.recipe_step_ingredients"
	// SearchRecipeStepIngredientsPermission is a household user permission.
	SearchRecipeStepIngredientsPermission Permission = "search.recipe_step_ingredients"
	// UpdateRecipeStepIngredientsPermission is a household user permission.
	UpdateRecipeStepIngredientsPermission Permission = "update.recipe_step_ingredients"
	// ArchiveRecipeStepIngredientsPermission is a household user permission.
	ArchiveRecipeStepIngredientsPermission Permission = "archive.recipe_step_ingredients"

	// CreateRecipeStepCompletionConditionsPermission is a household user permission.
	CreateRecipeStepCompletionConditionsPermission Permission = "create.recipe_step_completion_conditions"
	// ReadRecipeStepCompletionConditionsPermission is a household user permission.
	ReadRecipeStepCompletionConditionsPermission Permission = "read.recipe_step_completion_conditions"
	// SearchRecipeStepCompletionConditionsPermission is a household user permission.
	SearchRecipeStepCompletionConditionsPermission Permission = "search.recipe_step_completion_conditions"
	// UpdateRecipeStepCompletionConditionsPermission is a household user permission.
	UpdateRecipeStepCompletionConditionsPermission Permission = "update.recipe_step_completion_conditions"
	// ArchiveRecipeStepCompletionConditionsPermission is a household user permission.
	ArchiveRecipeStepCompletionConditionsPermission Permission = "archive.recipe_step_completion_conditions"

	// CreateRecipeStepProductsPermission is a household user permission.
	CreateRecipeStepProductsPermission Permission = "create.recipe_step_products"
	// ReadRecipeStepProductsPermission is a household user permission.
	ReadRecipeStepProductsPermission Permission = "read.recipe_step_products"
	// SearchRecipeStepProductsPermission is a household user permission.
	SearchRecipeStepProductsPermission Permission = "search.recipe_step_products"
	// UpdateRecipeStepProductsPermission is a household user permission.
	UpdateRecipeStepProductsPermission Permission = "update.recipe_step_products"
	// ArchiveRecipeStepProductsPermission is a household user permission.
	ArchiveRecipeStepProductsPermission Permission = "archive.recipe_step_products"

	// CreateMealPlansPermission is a household user permission.
	CreateMealPlansPermission Permission = "create.meal_plans"
	// ReadMealPlansPermission is a household user permission.
	ReadMealPlansPermission Permission = "read.meal_plans"
	// SearchMealPlansPermission is a household user permission.
	SearchMealPlansPermission Permission = "search.meal_plans"
	// UpdateMealPlansPermission is a household user permission.
	UpdateMealPlansPermission Permission = "update.meal_plans"
	// ArchiveMealPlansPermission is a household user permission.
	ArchiveMealPlansPermission Permission = "archive.meal_plans"

	// CreateMealPlanEventsPermission is a household user permission.
	CreateMealPlanEventsPermission Permission = "create.meal_plan_events"
	// ReadMealPlanEventsPermission is a household user permission.
	ReadMealPlanEventsPermission Permission = "read.meal_plan_events"
	// UpdateMealPlanEventsPermission is a household user permission.
	UpdateMealPlanEventsPermission Permission = "update.meal_plan_events"
	// ArchiveMealPlanEventsPermission is a household user permission.
	ArchiveMealPlanEventsPermission Permission = "archive.meal_plan_events"

	// CreateMealPlanOptionsPermission is a household user permission.
	CreateMealPlanOptionsPermission Permission = "create.meal_plan_options"
	// ReadMealPlanOptionsPermission is a household user permission.
	ReadMealPlanOptionsPermission Permission = "read.meal_plan_options"
	// SearchMealPlanOptionsPermission is a household user permission.
	SearchMealPlanOptionsPermission Permission = "search.meal_plan_options"
	// UpdateMealPlanOptionsPermission is a household user permission.
	UpdateMealPlanOptionsPermission Permission = "update.meal_plan_options"
	// ArchiveMealPlanOptionsPermission is a household user permission.
	ArchiveMealPlanOptionsPermission Permission = "archive.meal_plan_options"

	// CreateMealPlanGroceryListItemsPermission is a household user permission.
	CreateMealPlanGroceryListItemsPermission Permission = "create.meal_plan_grocery_list_items"
	// ReadMealPlanGroceryListItemsPermission is a household user permission.
	ReadMealPlanGroceryListItemsPermission Permission = "read.meal_plan_grocery_list_items"
	// UpdateMealPlanGroceryListItemsPermission is a household user permission.
	UpdateMealPlanGroceryListItemsPermission Permission = "update.meal_plan_grocery_list_items"
	// ArchiveMealPlanGroceryListItemsPermission is a household user permission.
	ArchiveMealPlanGroceryListItemsPermission Permission = "archive.meal_plan_grocery_list_items"

	// CreateMealPlanOptionVotesPermission is a household user permission.
	CreateMealPlanOptionVotesPermission Permission = "create.meal_plan_option_votes"
	// ReadMealPlanOptionVotesPermission is a household user permission.
	ReadMealPlanOptionVotesPermission Permission = "read.meal_plan_option_votes"
	// SearchMealPlanOptionVotesPermission is a household user permission.
	SearchMealPlanOptionVotesPermission Permission = "search.meal_plan_option_votes"
	// UpdateMealPlanOptionVotesPermission is a household user permission.
	UpdateMealPlanOptionVotesPermission Permission = "update.meal_plan_option_votes"
	// ArchiveMealPlanOptionVotesPermission is a household user permission.
	ArchiveMealPlanOptionVotesPermission Permission = "archive.meal_plan_option_votes"

	// ReadMealPlanTasksPermission is a household user permission.
	ReadMealPlanTasksPermission Permission = "read.meal_plan_tasks"
	// CreateMealPlanTasksPermission is a household user permission.
	CreateMealPlanTasksPermission Permission = "create.meal_plan_tasks"
	// UpdateMealPlanTasksPermission is a household user permission.
	UpdateMealPlanTasksPermission Permission = "update.meal_plan_tasks"

	// CreateServiceSettingsPermission is an admin user permission.
	CreateServiceSettingsPermission Permission = "create.service_settings"
	// ReadServiceSettingsPermission is an admin user permission.
	ReadServiceSettingsPermission Permission = "read.service_settings"
	// SearchServiceSettingsPermission is an admin user permission.
	SearchServiceSettingsPermission Permission = "search.service_settings"
	// ArchiveServiceSettingsPermission is an admin user permission.
	ArchiveServiceSettingsPermission Permission = "archive.service_settings"

	// CreateServiceSettingConfigurationsPermission is an admin user permission.
	CreateServiceSettingConfigurationsPermission Permission = "create.service_setting_configurations"
	// ReadServiceSettingConfigurationsPermission is an admin user permission.
	ReadServiceSettingConfigurationsPermission Permission = "read.service_setting_configurations"
	// UpdateServiceSettingConfigurationsPermission is an admin user permission.
	UpdateServiceSettingConfigurationsPermission Permission = "update.service_setting_configurations"
	// ArchiveServiceSettingConfigurationsPermission is an admin user permission.
	ArchiveServiceSettingConfigurationsPermission Permission = "archive.service_setting_configurations"

	// CreateUserIngredientPreferencesPermission is a household user permission.
	CreateUserIngredientPreferencesPermission Permission = "create.user_ingredient_preferences"
	// ReadUserIngredientPreferencesPermission is a household user permission.
	ReadUserIngredientPreferencesPermission Permission = "read.user_ingredient_preferences"
	// UpdateUserIngredientPreferencesPermission is a household user permission.
	UpdateUserIngredientPreferencesPermission Permission = "update.user_ingredient_preferences"
	// ArchiveUserIngredientPreferencesPermission is a household user permission.
	ArchiveUserIngredientPreferencesPermission Permission = "archive.user_ingredient_preferences"

	// CreateHouseholdInstrumentOwnershipsPermission is a household user permission.
	CreateHouseholdInstrumentOwnershipsPermission Permission = "create.household_instrument_ownerships"
	// ReadHouseholdInstrumentOwnershipsPermission is a household user permission.
	ReadHouseholdInstrumentOwnershipsPermission Permission = "read.household_instrument_ownerships"
	// UpdateHouseholdInstrumentOwnershipsPermission is a household user permission.
	UpdateHouseholdInstrumentOwnershipsPermission Permission = "update.household_instrument_ownerships"
	// ArchiveHouseholdInstrumentOwnershipsPermission is a household user permission.
	ArchiveHouseholdInstrumentOwnershipsPermission Permission = "archive.household_instrument_ownerships"

	// CreateRecipeRatingsPermission is a household user permission.
	CreateRecipeRatingsPermission Permission = "create.recipe_ratings"
	// ReadRecipeRatingsPermission is a household user permission.
	ReadRecipeRatingsPermission Permission = "read.recipe_ratings"
	// UpdateRecipeRatingsPermission is a household user permission.
	UpdateRecipeRatingsPermission Permission = "update.recipe_ratings"
	// ArchiveRecipeRatingsPermission is a household user permission.
	ArchiveRecipeRatingsPermission Permission = "archive.recipe_ratings"

	// CreateOAuth2ClientsPermission is a household admin permission.
	CreateOAuth2ClientsPermission Permission = "create.oauth2_clients"
	// ReadOAuth2ClientsPermission is a household admin permission.
	ReadOAuth2ClientsPermission Permission = "read.oauth2_clients"
	// ArchiveOAuth2ClientsPermission is a household admin permission.
	ArchiveOAuth2ClientsPermission Permission = "archive.oauth2_clients"

	// CreateUserNotificationsPermission is an admin user permission.
	CreateUserNotificationsPermission Permission = "create.user_notifications"
	// ReadUserNotificationsPermission is a household user permission.
	ReadUserNotificationsPermission Permission = "read.user_notifications"
	// UpdateUserNotificationsPermission is a household user permission.
	UpdateUserNotificationsPermission Permission = "update.user_notifications"
)

// ID implements the gorbac Permission interface.
func (p Permission) ID() string {
	return string(p)
}

// Match implements the gorbac Permission interface.
func (p Permission) Match(perm gorbac.Permission) bool {
	return p.ID() == perm.ID()
}

var (
	// service admin permissions.
	serviceAdminPermissions = []gorbac.Permission{
		CycleCookieSecretPermission,
		UpdateUserStatusPermission,
		ReadUserPermission,
		SearchUserPermission,
		CreateOAuth2ClientsPermission,
		ArchiveOAuth2ClientsPermission,
		ArchiveServiceSettingsPermission,
		CreateRecipesPermission,
		CreateValidInstrumentsPermission,
		UpdateValidInstrumentsPermission,
		ArchiveValidInstrumentsPermission,
		CreateValidVesselsPermission,
		UpdateValidVesselsPermission,
		ArchiveValidVesselsPermission,
		CreateValidIngredientsPermission,
		UpdateValidIngredientsPermission,
		ArchiveValidIngredientsPermission,
		CreateValidIngredientGroupsPermission,
		UpdateValidIngredientGroupsPermission,
		ArchiveValidIngredientGroupsPermission,
		CreateValidPreparationsPermission,
		UpdateValidPreparationsPermission,
		ArchiveValidPreparationsPermission,
		CreateValidMeasurementUnitsPermission,
		UpdateValidMeasurementUnitsPermission,
		ArchiveValidMeasurementUnitsPermission,
		CreateValidMeasurementUnitConversionsPermission,
		UpdateValidMeasurementUnitConversionsPermission,
		ArchiveValidMeasurementUnitConversionsPermission,
		CreateValidIngredientPreparationsPermission,
		UpdateValidIngredientPreparationsPermission,
		ArchiveValidIngredientPreparationsPermission,
		CreateValidIngredientStateIngredientsPermission,
		UpdateValidIngredientStateIngredientsPermission,
		ArchiveValidIngredientStateIngredientsPermission,
		CreateValidPreparationInstrumentsPermission,
		UpdateValidPreparationInstrumentsPermission,
		ArchiveValidPreparationInstrumentsPermission,
		CreateValidPreparationVesselsPermission,
		UpdateValidPreparationVesselsPermission,
		ArchiveValidPreparationVesselsPermission,
		CreateValidIngredientMeasurementUnitsPermission,
		UpdateValidIngredientMeasurementUnitsPermission,
		ArchiveValidIngredientMeasurementUnitsPermission,
		CreateValidIngredientStatesPermission,
		UpdateValidIngredientStatesPermission,
		ArchiveValidIngredientStatesPermission,
		CreateUserNotificationsPermission,
		// only admins can arbitrarily create these via the API, this is exclusively for integration test purposes.
		CreateServiceSettingsPermission,
		CreateMealPlanTasksPermission,
		CreateMealPlanGroceryListItemsPermission,
	}

	// household admin permissions.
	householdAdminPermissions = []gorbac.Permission{
		UpdateHouseholdPermission,
		ArchiveHouseholdPermission,
		TransferHouseholdPermission,
		InviteUserToHouseholdPermission,
		ModifyMemberPermissionsForHouseholdPermission,
		RemoveMemberHouseholdPermission,
		CreateWebhooksPermission,
		UpdateWebhooksPermission,
		ArchiveWebhooksPermission,
		CreateMealPlansPermission,
		UpdateMealPlansPermission,
		ArchiveMealPlansPermission,
		CreateMealPlanEventsPermission,
		UpdateMealPlanEventsPermission,
		ArchiveMealPlanEventsPermission,
		CreateMealPlanOptionsPermission,
		UpdateMealPlanOptionsPermission,
		ArchiveMealPlanOptionsPermission,
		CreateHouseholdInstrumentOwnershipsPermission,
		UpdateHouseholdInstrumentOwnershipsPermission,
		ArchiveHouseholdInstrumentOwnershipsPermission,
		CreateWebhookTriggerEventsPermission,
		ArchiveWebhookTriggerEventsPermission,
	}

	// household member permissions.
	householdMemberPermissions = []gorbac.Permission{
		ReadWebhooksPermission,
		ReadAuditLogEntriesPermission,
		ReadOAuth2ClientsPermission,
		ReadServiceSettingsPermission,
		SearchServiceSettingsPermission,
		CreateMealsPermission,
		ReadMealsPermission,
		UpdateMealsPermission,
		ArchiveMealsPermission,
		ReadRecipesPermission,
		SearchRecipesPermission,
		UpdateRecipesPermission,
		ArchiveRecipesPermission,
		CreateRecipeStepsPermission,
		ReadRecipeStepsPermission,
		SearchRecipeStepsPermission,
		UpdateRecipeStepsPermission,
		ArchiveRecipeStepsPermission,
		CreateRecipePrepTasksPermission,
		ReadRecipePrepTasksPermission,
		UpdateRecipePrepTasksPermission,
		ArchiveRecipePrepTasksPermission,
		CreateRecipeStepInstrumentsPermission,
		ReadRecipeStepInstrumentsPermission,
		SearchRecipeStepInstrumentsPermission,
		UpdateRecipeStepInstrumentsPermission,
		ArchiveRecipeStepInstrumentsPermission,
		CreateRecipeStepVesselsPermission,
		ReadRecipeStepVesselsPermission,
		SearchRecipeStepVesselsPermission,
		UpdateRecipeStepVesselsPermission,
		ArchiveRecipeStepVesselsPermission,
		CreateRecipeStepIngredientsPermission,
		ReadRecipeStepIngredientsPermission,
		SearchRecipeStepIngredientsPermission,
		UpdateRecipeStepIngredientsPermission,
		ArchiveRecipeStepIngredientsPermission,
		CreateRecipeStepCompletionConditionsPermission,
		ReadRecipeStepCompletionConditionsPermission,
		SearchRecipeStepCompletionConditionsPermission,
		UpdateRecipeStepCompletionConditionsPermission,
		ArchiveRecipeStepCompletionConditionsPermission,
		CreateRecipeStepProductsPermission,
		ReadRecipeStepProductsPermission,
		SearchRecipeStepProductsPermission,
		UpdateRecipeStepProductsPermission,
		ArchiveRecipeStepProductsPermission,
		ReadValidInstrumentsPermission,
		SearchValidInstrumentsPermission,
		ReadValidVesselsPermission,
		SearchValidVesselsPermission,
		ReadValidIngredientsPermission,
		SearchValidIngredientsPermission,
		ReadValidIngredientGroupsPermission,
		SearchValidIngredientGroupsPermission,
		ReadValidPreparationsPermission,
		SearchValidPreparationsPermission,
		ReadValidMeasurementUnitsPermission,
		SearchValidMeasurementUnitsPermission,
		ReadValidMeasurementUnitConversionsPermission,
		ReadValidIngredientPreparationsPermission,
		SearchValidIngredientPreparationsPermission,
		ReadValidIngredientStateIngredientsPermission,
		SearchValidIngredientStateIngredientsPermission,
		ReadValidPreparationInstrumentsPermission,
		SearchValidPreparationInstrumentsPermission,
		ReadValidPreparationVesselsPermission,
		SearchValidPreparationVesselsPermission,
		ReadValidIngredientMeasurementUnitsPermission,
		SearchValidIngredientMeasurementUnitsPermission,
		ReadMealPlansPermission,
		SearchMealPlansPermission,
		ReadMealPlanEventsPermission,
		ReadMealPlanOptionsPermission,
		SearchMealPlanOptionsPermission,
		ReadValidIngredientStatesPermission,
		ReadMealPlanGroceryListItemsPermission,
		UpdateMealPlanGroceryListItemsPermission,
		ArchiveMealPlanGroceryListItemsPermission,
		CreateMealPlanOptionVotesPermission,
		ReadMealPlanOptionVotesPermission,
		SearchMealPlanOptionVotesPermission,
		UpdateMealPlanOptionVotesPermission,
		ArchiveMealPlanOptionVotesPermission,
		CreateServiceSettingConfigurationsPermission,
		ReadServiceSettingConfigurationsPermission,
		UpdateServiceSettingConfigurationsPermission,
		ArchiveServiceSettingConfigurationsPermission,
		ReadMealPlanTasksPermission,
		UpdateMealPlanTasksPermission,
		CreateUserIngredientPreferencesPermission,
		ReadUserIngredientPreferencesPermission,
		UpdateUserIngredientPreferencesPermission,
		ArchiveUserIngredientPreferencesPermission,
		ReadHouseholdInstrumentOwnershipsPermission,
		CreateRecipeRatingsPermission,
		ReadRecipeRatingsPermission,
		UpdateRecipeRatingsPermission,
		ArchiveRecipeRatingsPermission,
		ReadUserNotificationsPermission,
		UpdateUserNotificationsPermission,
	}
)

func init() {
	// assign service admin permissions.
	for _, perm := range serviceAdminPermissions {
		must(serviceAdmin.Assign(perm))
	}

	// assign household admin permissions.
	for _, perm := range householdAdminPermissions {
		must(householdAdmin.Assign(perm))
	}

	// assign household member permissions.
	for _, perm := range householdMemberPermissions {
		must(householdMember.Assign(perm))
	}
}
