// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: exists.sql

package generated

import (
	"context"
	"database/sql"
)

const checkHouseholdInstrumentOwnershipExistence = `-- name: CheckHouseholdInstrumentOwnershipExistence :one

SELECT EXISTS ( SELECT household_instrument_ownerships.id FROM household_instrument_ownerships WHERE household_instrument_ownerships.archived_at IS NULL AND household_instrument_ownerships.id = $1 AND household_instrument_ownerships.belongs_to_household = $2 )
`

type CheckHouseholdInstrumentOwnershipExistenceParams struct {
	ID                 string
	BelongsToHousehold string
}

func (q *Queries) CheckHouseholdInstrumentOwnershipExistence(ctx context.Context, db DBTX, arg *CheckHouseholdInstrumentOwnershipExistenceParams) (bool, error) {
	row := db.QueryRowContext(ctx, checkHouseholdInstrumentOwnershipExistence, arg.ID, arg.BelongsToHousehold)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkHouseholdInvitationExistence = `-- name: CheckHouseholdInvitationExistence :one

SELECT EXISTS ( SELECT household_invitations.id FROM household_invitations WHERE household_invitations.archived_at IS NULL AND household_invitations.id = $1 )
`

func (q *Queries) CheckHouseholdInvitationExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkHouseholdInvitationExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkMealExistence = `-- name: CheckMealExistence :one

SELECT EXISTS ( SELECT meals.id FROM meals WHERE meals.archived_at IS NULL AND meals.id = $1 )
`

func (q *Queries) CheckMealExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkMealExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkMealPlanEventExistence = `-- name: CheckMealPlanEventExistence :one

SELECT EXISTS ( SELECT meal_plan_events.id FROM meal_plan_events WHERE meal_plan_events.archived_at IS NULL AND meal_plan_events.id = $1 )
`

func (q *Queries) CheckMealPlanEventExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkMealPlanEventExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkMealPlanExistence = `-- name: CheckMealPlanExistence :one

SELECT EXISTS ( SELECT meal_plans.id FROM meal_plans WHERE meal_plans.archived_at IS NULL AND meal_plans.id = $1 )
`

func (q *Queries) CheckMealPlanExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkMealPlanExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkMealPlanGroceryListItemExistence = `-- name: CheckMealPlanGroceryListItemExistence :one

SELECT EXISTS ( SELECT meal_plan_grocery_list_items.id FROM meal_plan_grocery_list_items WHERE meal_plan_grocery_list_items.archived_at IS NULL AND meal_plan_grocery_list_items.id = $1 )
`

func (q *Queries) CheckMealPlanGroceryListItemExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkMealPlanGroceryListItemExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkMealPlanOptionExistence = `-- name: CheckMealPlanOptionExistence :one

SELECT EXISTS (
	SELECT
	 meal_plan_options.id
	FROM
	 meal_plan_options
		JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event = meal_plan_events.id
		JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan = meal_plans.id
	WHERE
	 meal_plan_options.archived_at IS NULL
	AND meal_plan_options.belongs_to_meal_plan_event = $2
	AND meal_plan_options.id = $3
	AND meal_plan_events.archived_at IS NULL
	AND meal_plan_events.belongs_to_meal_plan = $1
	AND meal_plan_events.id = $2
	AND meal_plans.archived_at IS NULL
	AND meal_plans.id = $1
)
`

type CheckMealPlanOptionExistenceParams struct {
	BelongsToMealPlan      string
	BelongsToMealPlanEvent sql.NullString
	ID                     string
}

func (q *Queries) CheckMealPlanOptionExistence(ctx context.Context, db DBTX, arg *CheckMealPlanOptionExistenceParams) (bool, error) {
	row := db.QueryRowContext(ctx, checkMealPlanOptionExistence, arg.BelongsToMealPlan, arg.BelongsToMealPlanEvent, arg.ID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkMealPlanOptionVoteExistence = `-- name: CheckMealPlanOptionVoteExistence :one

SELECT EXISTS (
	SELECT
	 meal_plan_option_votes.id
	FROM
	 meal_plan_option_votes
		JOIN meal_plan_options ON meal_plan_option_votes.belongs_to_meal_plan_option=meal_plan_options.id
		JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event=meal_plan_events.id
		JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan=meal_plans.id
	WHERE meal_plan_option_votes.archived_at IS NULL
	AND meal_plan_option_votes.belongs_to_meal_plan_option = $1
	AND meal_plan_option_votes.id = $2
	AND meal_plan_options.archived_at IS NULL
	AND meal_plan_options.belongs_to_meal_plan_event = $3
	AND meal_plan_events.archived_at IS NULL
	AND meal_plan_events.belongs_to_meal_plan = $4
	AND meal_plan_options.id = $1
	AND meal_plans.archived_at IS NULL
	AND meal_plans.id = $4
)
`

type CheckMealPlanOptionVoteExistenceParams struct {
	BelongsToMealPlanOption string
	ID                      string
	BelongsToMealPlanEvent  sql.NullString
	BelongsToMealPlan       string
}

func (q *Queries) CheckMealPlanOptionVoteExistence(ctx context.Context, db DBTX, arg *CheckMealPlanOptionVoteExistenceParams) (bool, error) {
	row := db.QueryRowContext(ctx, checkMealPlanOptionVoteExistence,
		arg.BelongsToMealPlanOption,
		arg.ID,
		arg.BelongsToMealPlanEvent,
		arg.BelongsToMealPlan,
	)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkMealPlanTaskExistence = `-- name: CheckMealPlanTaskExistence :one

SELECT EXISTS (
	SELECT meal_plan_tasks.id
	FROM meal_plan_tasks
		FULL OUTER JOIN meal_plan_options ON meal_plan_tasks.belongs_to_meal_plan_option=meal_plan_options.id
		FULL OUTER JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event=meal_plan_events.id
		FULL OUTER JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan=meal_plans.id
	WHERE meal_plan_tasks.completed_at IS NULL
		AND meal_plans.id = $1
		AND meal_plans.archived_at IS NULL
		AND meal_plan_tasks.id = $2
)
`

type CheckMealPlanTaskExistenceParams struct {
	ID   string
	ID_2 string
}

func (q *Queries) CheckMealPlanTaskExistence(ctx context.Context, db DBTX, arg *CheckMealPlanTaskExistenceParams) (bool, error) {
	row := db.QueryRowContext(ctx, checkMealPlanTaskExistence, arg.ID, arg.ID_2)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkOAuth2ClientTokenExistence = `-- name: CheckOAuth2ClientTokenExistence :one

SELECT EXISTS ( SELECT oauth2_client_tokens.id FROM oauth2_client_tokens WHERE oauth2_client_tokens.archived_at IS NULL AND oauth2_client_tokens.id = $1 )
`

func (q *Queries) CheckOAuth2ClientTokenExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkOAuth2ClientTokenExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkRecipeExistence = `-- name: CheckRecipeExistence :one

SELECT EXISTS ( SELECT recipes.id FROM recipes WHERE recipes.archived_at IS NULL AND recipes.id = $1 )
`

func (q *Queries) CheckRecipeExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkRecipeExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkRecipeMediaExistence = `-- name: CheckRecipeMediaExistence :one

SELECT EXISTS ( SELECT recipe_media.id FROM recipe_media WHERE recipe_media.archived_at IS NULL AND recipe_media.id = $1 )
`

func (q *Queries) CheckRecipeMediaExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkRecipeMediaExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkRecipePrepTaskExistence = `-- name: CheckRecipePrepTaskExistence :one

SELECT EXISTS (
	SELECT recipe_prep_tasks.id
	FROM recipe_prep_tasks
	JOIN recipes ON recipe_prep_tasks.belongs_to_recipe=recipes.id
	WHERE recipe_prep_tasks.archived_at IS NULL
	  AND recipe_prep_tasks.belongs_to_recipe = $1
	  AND recipe_prep_tasks.id = $2
	  AND recipes.archived_at IS NULL
	  AND recipes.id = $1
)
`

type CheckRecipePrepTaskExistenceParams struct {
	BelongsToRecipe string
	ID              string
}

func (q *Queries) CheckRecipePrepTaskExistence(ctx context.Context, db DBTX, arg *CheckRecipePrepTaskExistenceParams) (bool, error) {
	row := db.QueryRowContext(ctx, checkRecipePrepTaskExistence, arg.BelongsToRecipe, arg.ID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkRecipeRatingExistence = `-- name: CheckRecipeRatingExistence :one

SELECT EXISTS ( SELECT recipe_ratings.id FROM recipe_ratings WHERE recipe_ratings.archived_at IS NULL AND recipe_ratings.id = $1 )
`

func (q *Queries) CheckRecipeRatingExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkRecipeRatingExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkRecipeStepCompletionConditionExistence = `-- name: CheckRecipeStepCompletionConditionExistence :one

SELECT EXISTS ( SELECT recipe_step_completion_conditions.id FROM recipe_step_completion_conditions JOIN recipe_steps ON recipe_step_completion_conditions.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_completion_conditions.archived_at IS NULL AND recipe_step_completion_conditions.belongs_to_recipe_step = $1 AND recipe_step_completion_conditions.id = $2 AND recipe_steps.archived_at IS NULL AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.archived_at IS NULL AND recipes.id = $5 )
`

type CheckRecipeStepCompletionConditionExistenceParams struct {
	BelongsToRecipeStep string
	ID                  string
	BelongsToRecipe     string
	ID_2                string
	ID_3                string
}

func (q *Queries) CheckRecipeStepCompletionConditionExistence(ctx context.Context, db DBTX, arg *CheckRecipeStepCompletionConditionExistenceParams) (bool, error) {
	row := db.QueryRowContext(ctx, checkRecipeStepCompletionConditionExistence,
		arg.BelongsToRecipeStep,
		arg.ID,
		arg.BelongsToRecipe,
		arg.ID_2,
		arg.ID_3,
	)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkRecipeStepExistence = `-- name: CheckRecipeStepExistence :one

SELECT EXISTS (
	SELECT recipe_steps.id
	FROM recipe_steps
	JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
	WHERE recipe_steps.archived_at IS NULL
	  AND recipe_steps.belongs_to_recipe = $1
	  AND recipe_steps.id = $2
	  AND recipes.archived_at IS NULL
	  AND recipes.id = $1
)
`

type CheckRecipeStepExistenceParams struct {
	BelongsToRecipe string
	ID              string
}

func (q *Queries) CheckRecipeStepExistence(ctx context.Context, db DBTX, arg *CheckRecipeStepExistenceParams) (bool, error) {
	row := db.QueryRowContext(ctx, checkRecipeStepExistence, arg.BelongsToRecipe, arg.ID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkRecipeStepIngredientExistence = `-- name: CheckRecipeStepIngredientExistence :one

SELECT EXISTS ( SELECT recipe_step_ingredients.id FROM recipe_step_ingredients JOIN recipe_steps ON recipe_step_ingredients.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_ingredients.archived_at IS NULL AND recipe_step_ingredients.belongs_to_recipe_step = $1 AND recipe_step_ingredients.id = $2 AND recipe_steps.archived_at IS NULL AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.archived_at IS NULL AND recipes.id = $5 )
`

type CheckRecipeStepIngredientExistenceParams struct {
	BelongsToRecipeStep string
	ID                  string
	BelongsToRecipe     string
	ID_2                string
	ID_3                string
}

func (q *Queries) CheckRecipeStepIngredientExistence(ctx context.Context, db DBTX, arg *CheckRecipeStepIngredientExistenceParams) (bool, error) {
	row := db.QueryRowContext(ctx, checkRecipeStepIngredientExistence,
		arg.BelongsToRecipeStep,
		arg.ID,
		arg.BelongsToRecipe,
		arg.ID_2,
		arg.ID_3,
	)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkRecipeStepInstrumentExistence = `-- name: CheckRecipeStepInstrumentExistence :one

SELECT EXISTS ( SELECT recipe_step_instruments.id FROM recipe_step_instruments JOIN recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_instruments.archived_at IS NULL AND recipe_step_instruments.belongs_to_recipe_step = $1 AND recipe_step_instruments.id = $2 AND recipe_steps.archived_at IS NULL AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.archived_at IS NULL AND recipes.id = $5 )
`

type CheckRecipeStepInstrumentExistenceParams struct {
	BelongsToRecipeStep string
	ID                  string
	BelongsToRecipe     string
	ID_2                string
	ID_3                string
}

func (q *Queries) CheckRecipeStepInstrumentExistence(ctx context.Context, db DBTX, arg *CheckRecipeStepInstrumentExistenceParams) (bool, error) {
	row := db.QueryRowContext(ctx, checkRecipeStepInstrumentExistence,
		arg.BelongsToRecipeStep,
		arg.ID,
		arg.BelongsToRecipe,
		arg.ID_2,
		arg.ID_3,
	)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkRecipeStepProductExistence = `-- name: CheckRecipeStepProductExistence :one

SELECT EXISTS ( SELECT recipe_step_products.id FROM recipe_step_products JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_products.archived_at IS NULL AND recipe_step_products.belongs_to_recipe_step = $1 AND recipe_step_products.id = $2 AND recipe_steps.archived_at IS NULL AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.archived_at IS NULL AND recipes.id = $5 )
`

type CheckRecipeStepProductExistenceParams struct {
	BelongsToRecipeStep string
	ID                  string
	BelongsToRecipe     string
	ID_2                string
	ID_3                string
}

func (q *Queries) CheckRecipeStepProductExistence(ctx context.Context, db DBTX, arg *CheckRecipeStepProductExistenceParams) (bool, error) {
	row := db.QueryRowContext(ctx, checkRecipeStepProductExistence,
		arg.BelongsToRecipeStep,
		arg.ID,
		arg.BelongsToRecipe,
		arg.ID_2,
		arg.ID_3,
	)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkRecipeStepVesselExistence = `-- name: CheckRecipeStepVesselExistence :one

SELECT EXISTS (
    SELECT recipe_step_vessels.id
    FROM recipe_step_vessels
        JOIN recipe_steps ON recipe_step_vessels.belongs_to_recipe_step=recipe_steps.id
        JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
    WHERE recipe_step_vessels.archived_at IS NULL
        AND recipe_step_vessels.belongs_to_recipe_step = $1
        AND recipe_step_vessels.id = $2
        AND recipe_steps.archived_at IS NULL
        AND recipe_steps.belongs_to_recipe = $3
        AND recipe_steps.id = $4
        AND recipes.archived_at IS NULL
        AND recipes.id = $5
)
`

type CheckRecipeStepVesselExistenceParams struct {
	BelongsToRecipeStep string
	ID                  string
	BelongsToRecipe     string
	ID_2                string
	ID_3                string
}

func (q *Queries) CheckRecipeStepVesselExistence(ctx context.Context, db DBTX, arg *CheckRecipeStepVesselExistenceParams) (bool, error) {
	row := db.QueryRowContext(ctx, checkRecipeStepVesselExistence,
		arg.BelongsToRecipeStep,
		arg.ID,
		arg.BelongsToRecipe,
		arg.ID_2,
		arg.ID_3,
	)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkServiceSettingConfigurationExistence = `-- name: CheckServiceSettingConfigurationExistence :one

SELECT EXISTS ( SELECT service_setting_configurations.id FROM service_setting_configurations WHERE service_setting_configurations.archived_at IS NULL AND service_setting_configurations.id = $1 )
`

func (q *Queries) CheckServiceSettingConfigurationExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkServiceSettingConfigurationExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkServiceSettingExistence = `-- name: CheckServiceSettingExistence :one

SELECT EXISTS ( SELECT service_settings.id FROM service_settings WHERE service_settings.archived_at IS NULL AND service_settings.id = $1 )
`

func (q *Queries) CheckServiceSettingExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkServiceSettingExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkUserIngredientPreferenceExistence = `-- name: CheckUserIngredientPreferenceExistence :one

SELECT EXISTS ( SELECT user_ingredient_preferences.id FROM user_ingredient_preferences WHERE user_ingredient_preferences.archived_at IS NULL AND user_ingredient_preferences.id = $1 AND user_ingredient_preferences.belongs_to_user = $2 )
`

type CheckUserIngredientPreferenceExistenceParams struct {
	ID            string
	BelongsToUser string
}

func (q *Queries) CheckUserIngredientPreferenceExistence(ctx context.Context, db DBTX, arg *CheckUserIngredientPreferenceExistenceParams) (bool, error) {
	row := db.QueryRowContext(ctx, checkUserIngredientPreferenceExistence, arg.ID, arg.BelongsToUser)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkValidIngredientExistence = `-- name: CheckValidIngredientExistence :one

SELECT EXISTS ( SELECT valid_ingredients.id FROM valid_ingredients WHERE valid_ingredients.archived_at IS NULL AND valid_ingredients.id = $1 )
`

func (q *Queries) CheckValidIngredientExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkValidIngredientExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkValidIngredientGroupExistence = `-- name: CheckValidIngredientGroupExistence :one

SELECT EXISTS ( SELECT valid_ingredient_groups.id FROM valid_ingredient_groups WHERE valid_ingredient_groups.archived_at IS NULL AND valid_ingredient_groups.id = $1 )
`

func (q *Queries) CheckValidIngredientGroupExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkValidIngredientGroupExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkValidIngredientMeasurementUnitExistence = `-- name: CheckValidIngredientMeasurementUnitExistence :one

SELECT EXISTS ( SELECT valid_ingredient_measurement_units.id FROM valid_ingredient_measurement_units WHERE valid_ingredient_measurement_units.archived_at IS NULL AND valid_ingredient_measurement_units.id = $1 )
`

func (q *Queries) CheckValidIngredientMeasurementUnitExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkValidIngredientMeasurementUnitExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkValidIngredientPreparationExistence = `-- name: CheckValidIngredientPreparationExistence :one

SELECT EXISTS ( SELECT valid_ingredient_preparations.id FROM valid_ingredient_preparations WHERE valid_ingredient_preparations.archived_at IS NULL AND valid_ingredient_preparations.id = $1 )
`

func (q *Queries) CheckValidIngredientPreparationExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkValidIngredientPreparationExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkValidIngredientStateExistence = `-- name: CheckValidIngredientStateExistence :one

SELECT EXISTS ( SELECT valid_ingredient_states.id FROM valid_ingredient_states WHERE valid_ingredient_states.archived_at IS NULL AND valid_ingredient_states.id = $1 )
`

func (q *Queries) CheckValidIngredientStateExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkValidIngredientStateExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkValidIngredientStateIngredientExistence = `-- name: CheckValidIngredientStateIngredientExistence :one

SELECT EXISTS ( SELECT valid_ingredient_state_ingredients.id FROM valid_ingredient_state_ingredients WHERE valid_ingredient_state_ingredients.archived_at IS NULL AND valid_ingredient_state_ingredients.id = $1 )
`

func (q *Queries) CheckValidIngredientStateIngredientExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkValidIngredientStateIngredientExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkValidInstrumentExistence = `-- name: CheckValidInstrumentExistence :one

SELECT EXISTS ( SELECT valid_instruments.id FROM valid_instruments WHERE valid_instruments.archived_at IS NULL AND valid_instruments.id = $1 )
`

func (q *Queries) CheckValidInstrumentExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkValidInstrumentExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkValidMeasurementUnitConversionExistence = `-- name: CheckValidMeasurementUnitConversionExistence :one

SELECT EXISTS ( SELECT valid_measurement_conversions.id FROM valid_measurement_conversions WHERE valid_measurement_conversions.archived_at IS NULL AND valid_measurement_conversions.id = $1 )
`

func (q *Queries) CheckValidMeasurementUnitConversionExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkValidMeasurementUnitConversionExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkValidMeasurementUnitExistence = `-- name: CheckValidMeasurementUnitExistence :one

SELECT EXISTS ( SELECT valid_measurement_units.id FROM valid_measurement_units WHERE valid_measurement_units.archived_at IS NULL AND valid_measurement_units.id = $1 )
`

func (q *Queries) CheckValidMeasurementUnitExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkValidMeasurementUnitExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkValidPreparationExistence = `-- name: CheckValidPreparationExistence :one

SELECT EXISTS ( SELECT valid_preparations.id FROM valid_preparations WHERE valid_preparations.archived_at IS NULL AND valid_preparations.id = $1 )
`

func (q *Queries) CheckValidPreparationExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkValidPreparationExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkValidPreparationInstrumentExistence = `-- name: CheckValidPreparationInstrumentExistence :one

SELECT EXISTS ( SELECT valid_preparation_instruments.id FROM valid_preparation_instruments WHERE valid_preparation_instruments.archived_at IS NULL AND valid_preparation_instruments.id = $1 )
`

func (q *Queries) CheckValidPreparationInstrumentExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkValidPreparationInstrumentExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkValidPreparationVesselExistence = `-- name: CheckValidPreparationVesselExistence :one

SELECT EXISTS ( SELECT valid_preparation_vessels.id FROM valid_preparation_vessels WHERE valid_preparation_vessels.archived_at IS NULL AND valid_preparation_vessels.id = $1 )
`

func (q *Queries) CheckValidPreparationVesselExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkValidPreparationVesselExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkValidVesselExistence = `-- name: CheckValidVesselExistence :one

SELECT EXISTS ( SELECT valid_vessels.id FROM valid_vessels WHERE valid_vessels.archived_at IS NULL AND valid_vessels.id = $1 )
`

func (q *Queries) CheckValidVesselExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkValidVesselExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkWebhookExistence = `-- name: CheckWebhookExistence :one

SELECT EXISTS (
	SELECT webhooks.id
	FROM webhooks
	WHERE webhooks.archived_at IS NULL
	  AND webhooks.belongs_to_household = $1
	  AND webhooks.id = $2
)
`

type CheckWebhookExistenceParams struct {
	BelongsToHousehold string
	ID                 string
}

func (q *Queries) CheckWebhookExistence(ctx context.Context, db DBTX, arg *CheckWebhookExistenceParams) (bool, error) {
	row := db.QueryRowContext(ctx, checkWebhookExistence, arg.BelongsToHousehold, arg.ID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}
