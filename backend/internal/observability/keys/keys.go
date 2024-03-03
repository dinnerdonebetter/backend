package keys

const (
	// RequesterIDKey is the standard key for referring to a requesting user's ID.
	RequesterIDKey = "request.made_by"
	// HouseholdIDKey is the standard key for referring to a household ID.
	HouseholdIDKey = "household.id"
	// HouseholdInvitationIDKey is the standard key for referring to a household ID.
	HouseholdInvitationIDKey = "household_invitation.id"
	// HouseholdInvitationTokenKey is the standard key for referring to a household invitation token.
	HouseholdInvitationTokenKey = "household_invitation.token"
	// ActiveHouseholdIDKey is the standard key for referring to an active household ID.
	ActiveHouseholdIDKey = "active_household.id"
	// UserIDKey is the standard key for referring to a user ID.
	UserIDKey = "user.id"
	// UserNotificationIDKey is the standard key for referring to a user notification ID.
	UserNotificationIDKey = "user_notification.id"
	// UserEmailAddressKey is the standard key for referring to a user's email address.
	UserEmailAddressKey = "user.email_address"
	// UserIsServiceAdminKey is the standard key for referring to a user's admin status.
	UserIsServiceAdminKey = "user.is_admin"
	// UsernameKey is the standard key for referring to a username.
	UsernameKey = "user.username"
	// NameKey is the standard key for referring to a name.
	NameKey = "name"
	// SpanIDKey is the standard key for referring to a span ID.
	SpanIDKey = "span.id"
	// TraceIDKey is the standard key for referring to a trace ID.
	TraceIDKey = "trace.id"
	// FilterCreatedAfterKey is the standard key for referring to a types.QueryFilter's CreatedAfter field.
	FilterCreatedAfterKey = "query_filter.created_after"
	// FilterCreatedBeforeKey is the standard key for referring to a types.QueryFilter's CreatedBefore field.
	FilterCreatedBeforeKey = "query_filter.created_before"
	// FilterUpdatedAfterKey is the standard key for referring to a types.QueryFilter's UpdatedAfter field.
	FilterUpdatedAfterKey = "query_filter.updated_after"
	// FilterUpdatedBeforeKey is the standard key for referring to a types.QueryFilter's UpdatedAfter field.
	FilterUpdatedBeforeKey = "query_filter.updated_before"
	// FilterSortByKey is the standard key for referring to a types.QueryFilter's SortBy field.
	FilterSortByKey = "query_filter.sort_by"
	// FilterPageKey is the standard key for referring to a types.QueryFilter's page.
	FilterPageKey = "query_filter.page"
	// FilterLimitKey is the standard key for referring to a types.QueryFilter's limit.
	FilterLimitKey = "query_filter.limit"
	// FilterIsNilKey is the standard key for referring to a types.QueryFilter's null status.
	FilterIsNilKey = "query_filter.is_nil"
	// WebhookIDKey is the standard key for referring to a webhook's ID.
	WebhookIDKey = "webhook.id"
	// WebhookTriggerEventIDKey is the standard key for referring to a webhook trigger event's ID.
	WebhookTriggerEventIDKey = "webhook_trigger_event.id"
	// AuditLogEntryIDKey is the standard key for referring to an audit log entry's ID.
	AuditLogEntryIDKey = "audit_log_entry.id"
	// AuditLogEntryResourceTypesKey is the standard key for referring to an audit log entry's resource type.
	AuditLogEntryResourceTypesKey = "audit_log_entry.resource_types"
	// URLKey is the standard key for referring to a URL.
	URLKey = "url"
	// PasswordResetTokenIDKey is the standard key for referring to a password reset token's ID.
	PasswordResetTokenIDKey = "password_reset_token.id"
	// RequestHeadersKey is the standard key for referring to a http.Request's Headers.
	RequestHeadersKey = "request.headers"
	// RequestMethodKey is the standard key for referring to a http.Request's Method.
	RequestMethodKey = "request.method"
	// RequestURIKey is the standard key for referring to a http.Request's URI.
	RequestURIKey = "request.uri"
	// ResponseStatusKey is the standard key for referring to a http.Request's status.
	ResponseStatusKey = "response.status"
	// ResponseBytesWrittenKey is the standard key for referring to a http.Request's bytes written.
	ResponseBytesWrittenKey = "response.bytes_written"
	// ResponseHeadersKey is the standard key for referring to a http.Response's Headers.
	ResponseHeadersKey = "response.headers"
	// ReasonKey is the standard key for referring to a reason for a change.
	ReasonKey = "reason"
	// URLQueryKey is the standard key for referring to a URL query.
	URLQueryKey = "url.query"
	// SearchQueryKey is the standard key for referring to a search query parameter value.
	SearchQueryKey = "search_query"
	// UserAgentOSKey is the standard key for referring to a user agent's OS.
	UserAgentOSKey = "os"
	// UserAgentBotKey is the standard key for referring to a user agent's bot status.
	UserAgentBotKey = "is_bot"
	// UserAgentMobileKey is the standard key for referring to user agent's mobile status.
	UserAgentMobileKey = "is_mobile"
	// ValidationErrorKey is the standard key for referring to a struct validation error.
	ValidationErrorKey = "validation_error"

	// ValidInstrumentIDKey is the standard key for referring to a valid instrument's ID.
	ValidInstrumentIDKey = "valid_instrument.id"

	// ValidVesselIDKey is the standard key for referring to a valid vessel's ID.
	ValidVesselIDKey = "valid_vessel.id"

	// ValidIngredientIDKey is the standard key for referring to a valid ingredient's ID.
	ValidIngredientIDKey = "valid_ingredient.id"

	// ValidIngredientGroupIDKey is the standard key for referring to a valid ingredient group's ID.
	ValidIngredientGroupIDKey = "valid_ingredient_group.id"

	// ValidPreparationIDKey is the standard key for referring to a valid preparation's ID.
	ValidPreparationIDKey = "valid_preparation.id"

	// ValidIngredientStateIDKey is the standard key for referring to a valid ingredient state's ID.
	ValidIngredientStateIDKey = "valid_ingredient_state.id"

	// ValidIngredientStateIngredientIDKey is the standard key for referring to a valid ingredient state ingredient's ID.
	ValidIngredientStateIngredientIDKey = "valid_ingredient_state_ingredient.id"

	// ValidIngredientPreparationIDKey is the standard key for referring to a valid ingredient preparation's ID.
	ValidIngredientPreparationIDKey = "valid_ingredient_preparation.id"

	// ValidPreparationInstrumentIDKey is the standard key for referring to a valid instrument preparation's ID.
	ValidPreparationInstrumentIDKey = "valid_preparation_instrument.id"

	// ValidIngredientMeasurementUnitIDKey is the standard key for referring to a valid instrument preparation's ID.
	ValidIngredientMeasurementUnitIDKey = "valid_ingredient_measurement_unit.id"

	// MealIDKey is the standard key for referring to a meal's ID.
	MealIDKey = "meal.id"

	// RecipeIDKey is the standard key for referring to a recipe's ID.
	RecipeIDKey = "recipe.id"

	// RecipeStepIDKey is the standard key for referring to a recipe step's ID.
	RecipeStepIDKey = "recipe_step.id"

	// RecipePrepTaskIDKey is the standard key for referring to a recipe prep task's ID.
	RecipePrepTaskIDKey = "recipe_prep_task.id"

	// RecipeStepInstrumentIDKey is the standard key for referring to a recipe step instrument's ID.
	RecipeStepInstrumentIDKey = "recipe_step_instrument.id"

	// RecipeStepVesselIDKey is the standard key for referring to a recipe step vessel's ID.
	RecipeStepVesselIDKey = "recipe_step_vessel.id"

	// RecipeStepIngredientIDKey is the standard key for referring to a recipe step ingredient's ID.
	RecipeStepIngredientIDKey = "recipe_step_ingredient.id"

	// RecipeStepCompletionConditionIDKey is the standard key for referring to a recipe step completion condition's ID.
	RecipeStepCompletionConditionIDKey = "recipe_step_completion_condition.id"

	// RecipeStepProductIDKey is the standard key for referring to a recipe step product's ID.
	RecipeStepProductIDKey = "recipe_step_product.id"

	// MealPlanIDKey is the standard key for referring to a meal plan's ID.
	MealPlanIDKey = "meal_plan.id"

	// MealPlanEventIDKey is the standard key for referring to a meal plan event's ID.
	MealPlanEventIDKey = "meal_plan_event.id"

	// MealPlanOptionIDKey is the standard key for referring to a meal plan option's ID.
	MealPlanOptionIDKey = "meal_plan_option.id"

	// MealPlanOptionVoteIDKey is the standard key for referring to a meal plan option vote's ID.
	MealPlanOptionVoteIDKey = "meal_plan_option_vote.id"

	// ValidMeasurementUnitIDKey is the standard key for referring to a valid measurement unit's ID.
	ValidMeasurementUnitIDKey = "valid_measurement_unit.id"

	// MealPlanTaskIDKey is the standard key for referring to a meal plan task's ID.
	MealPlanTaskIDKey = "meal_plan_task.id"

	// MealPlanGroceryListItemIDKey is the standard key for referring to a meal plan grocery list item's ID.
	MealPlanGroceryListItemIDKey = "meal_plan_grocery_list_item.id"

	// ValidMeasurementUnitConversionIDKey is the standard key for referring to a valid measurement conversion's ID.
	ValidMeasurementUnitConversionIDKey = "valid_measurement_conversion.id"

	// RecipeMediaIDKey is the standard key for referring to a recipe media's ID.
	RecipeMediaIDKey = "recipe_media.id"

	// ServiceSettingIDKey is the standard key for referring to a service setting's ID.
	ServiceSettingIDKey = "service_setting.id"

	// ServiceSettingNameKey is the standard key for referring to a service setting's ID.
	ServiceSettingNameKey = "service_setting.name"

	// ServiceSettingConfigurationIDKey is the standard key for referring to a service setting configuration's ID.
	ServiceSettingConfigurationIDKey = "service_setting_configuration.id"

	// UserIngredientPreferenceIDKey is the standard key for referring to a user ingredient preference's ID.
	UserIngredientPreferenceIDKey = "user_ingredient_preference.id"

	// HouseholdInstrumentOwnershipIDKey is the standard key for referring to a household instrument ownership's ID.
	HouseholdInstrumentOwnershipIDKey = "household_instrument_ownership.id"

	// RecipeRatingIDKey is the standard key for referring to a recipe rating's ID.
	RecipeRatingIDKey = "recipe_rating.id"

	// OAuth2ClientIDKey is the standard key for referring to an OAuth2 client's database ID.
	OAuth2ClientIDKey = "oauth2_clients.id"

	// OAuth2ClientClientIDKey is the standard key for referring to an OAuth2 client's client ID.
	OAuth2ClientClientIDKey = "oauth2_clients.client_id"

	// OAuth2ClientTokenIDKey is the standard key for referring to an OAuth2 client token's ID.
	/* #nosec G101 */
	OAuth2ClientTokenIDKey = "oauth2_client_tokens.id"

	// OAuth2ClientTokenCodeKey is the standard key for referring to an OAuth2 client token's code.
	/* #nosec G101 */
	OAuth2ClientTokenCodeKey = "oauth2_client_tokens.code"

	// OAuth2ClientTokenAccessKey is the standard key for referring to an OAuth2 client token's access.
	/* #nosec G101 */
	OAuth2ClientTokenAccessKey = "oauth2_client_tokens.access"

	// OAuth2ClientTokenRefreshKey is the standard key for referring to an OAuth2 client token's refresh.
	/* #nosec G101 */
	OAuth2ClientTokenRefreshKey = "oauth2_client_tokens.refresh"

	// ValidPreparationVesselIDKey is the standard key for referring to a valid preparation vessel's ID.
	ValidPreparationVesselIDKey = "valid_preparation_vessels.id"
)
