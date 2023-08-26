// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: meal_plan_tasks.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const changeMealPlanTaskStatus = `-- name: ChangeMealPlanTaskStatus :exec

UPDATE meal_plan_tasks SET completed_at = $4, status_explanation = $3, status = $2 WHERE id = $1
`

type ChangeMealPlanTaskStatusParams struct {
	ID                string
	Status            PrepStepStatus
	StatusExplanation string
	CompletedAt       sql.NullTime
}

func (q *Queries) ChangeMealPlanTaskStatus(ctx context.Context, db DBTX, arg *ChangeMealPlanTaskStatusParams) error {
	_, err := db.ExecContext(ctx, changeMealPlanTaskStatus,
		arg.ID,
		arg.Status,
		arg.StatusExplanation,
		arg.CompletedAt,
	)
	return err
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
	MealPlanID     string
	MealPlanTaskID string
}

func (q *Queries) CheckMealPlanTaskExistence(ctx context.Context, db DBTX, arg *CheckMealPlanTaskExistenceParams) (bool, error) {
	row := db.QueryRowContext(ctx, checkMealPlanTaskExistence, arg.MealPlanID, arg.MealPlanTaskID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createMealPlanTask = `-- name: CreateMealPlanTask :exec

INSERT INTO meal_plan_tasks (id,status,status_explanation,creation_explanation,belongs_to_meal_plan_option,belongs_to_recipe_prep_task,assigned_to_user)
VALUES ($1,$2,$3,$4,$5,$6,$7)
`

type CreateMealPlanTaskParams struct {
	ID                      string
	Status                  PrepStepStatus
	StatusExplanation       string
	CreationExplanation     string
	BelongsToMealPlanOption string
	BelongsToRecipePrepTask string
	AssignedToUser          sql.NullString
}

func (q *Queries) CreateMealPlanTask(ctx context.Context, db DBTX, arg *CreateMealPlanTaskParams) error {
	_, err := db.ExecContext(ctx, createMealPlanTask,
		arg.ID,
		arg.Status,
		arg.StatusExplanation,
		arg.CreationExplanation,
		arg.BelongsToMealPlanOption,
		arg.BelongsToRecipePrepTask,
		arg.AssignedToUser,
	)
	return err
}

const getMealPlanTask = `-- name: GetMealPlanTask :one

SELECT
	meal_plan_tasks.id,
	meal_plan_options.id as meal_plan_option_id,
	meal_plan_options.assigned_cook as meal_plan_option_assigned_cook,
	meal_plan_options.assigned_dishwasher as meal_plan_option_assigned_dishwasher,
	meal_plan_options.chosen as meal_plan_option_chosen,
	meal_plan_options.tiebroken as meal_plan_option_tiebroken,
    meal_plan_options.meal_scale as meal_plan_option_meal_scale,
	meal_plan_options.meal_id as meal_plan_option_meal_id,
	meal_plan_options.notes as meal_plan_option_notes,
	meal_plan_options.created_at as meal_plan_option_created_at,
	meal_plan_options.last_updated_at as meal_plan_option_last_updated_at,
	meal_plan_options.archived_at as meal_plan_option_archived_at,
	meal_plan_options.belongs_to_meal_plan_event as meal_plan_option_belongs_to_meal_plan_event,
	recipe_prep_tasks.id as prep_task_id,
	recipe_prep_tasks.name as prep_task_name,
	recipe_prep_tasks.description as prep_task_description,
	recipe_prep_tasks.notes as prep_task_notes,
	recipe_prep_tasks.optional as prep_task_optional,
	recipe_prep_tasks.explicit_storage_instructions as prep_task_explicit_storage_instructions,
	recipe_prep_tasks.minimum_time_buffer_before_recipe_in_seconds as prep_task_minimum_time_buffer_before_recipe_in_seconds,
	recipe_prep_tasks.maximum_time_buffer_before_recipe_in_seconds as prep_task_maximum_time_buffer_before_recipe_in_seconds,
	recipe_prep_tasks.storage_type as prep_task_storage_type,
	recipe_prep_tasks.minimum_storage_temperature_in_celsius as prep_task_minimum_storage_temperature_in_celsius,
	recipe_prep_tasks.maximum_storage_temperature_in_celsius as prep_task_maximum_storage_temperature_in_celsius,
	recipe_prep_tasks.belongs_to_recipe as prep_task_belongs_to_recipe,
	recipe_prep_tasks.created_at as prep_task_created_at,
	recipe_prep_tasks.last_updated_at as prep_task_last_updated_at,
	recipe_prep_tasks.archived_at as prep_task_archived_at,
	recipe_prep_task_steps.id as prep_task_step_id,
	recipe_prep_task_steps.belongs_to_recipe_step as prep_task_step_belongs_to_recipe_step,
	recipe_prep_task_steps.belongs_to_recipe_prep_task as prep_task_step_belongs_to_recipe_prep_task,
	recipe_prep_task_steps.satisfies_recipe_step as prep_task_step_satisfies_recipe_step,
	meal_plan_tasks.created_at,
	meal_plan_tasks.last_updated_at,
	meal_plan_tasks.completed_at,
	meal_plan_tasks.status,
	meal_plan_tasks.creation_explanation,
	meal_plan_tasks.status_explanation,
	meal_plan_tasks.assigned_to_user
FROM meal_plan_tasks
	JOIN meal_plan_options ON meal_plan_tasks.belongs_to_meal_plan_option=meal_plan_options.id
	JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event=meal_plan_events.id
	JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan=meal_plans.id
	JOIN meals ON meal_plan_options.meal_id=meals.id
	JOIN recipe_prep_tasks ON meal_plan_tasks.belongs_to_recipe_prep_task=recipe_prep_tasks.id
	JOIN recipe_prep_task_steps ON recipe_prep_task_steps.belongs_to_recipe_prep_task=recipe_prep_tasks.id
	JOIN recipe_steps ON recipe_prep_task_steps.belongs_to_recipe_step=recipe_steps.id
WHERE meal_plan_options.archived_at IS NULL
	AND meal_plan_events.archived_at IS NULL
	AND meal_plans.archived_at IS NULL
	AND meals.archived_at IS NULL
	AND recipe_steps.archived_at IS NULL
	AND meal_plan_tasks.id = $1
`

type GetMealPlanTaskRow struct {
	MealPlanOptionCreatedAt                        time.Time
	CreatedAt                                      time.Time
	PrepTaskCreatedAt                              time.Time
	PrepTaskArchivedAt                             sql.NullTime
	MealPlanOptionLastUpdatedAt                    sql.NullTime
	CompletedAt                                    sql.NullTime
	LastUpdatedAt                                  sql.NullTime
	PrepTaskLastUpdatedAt                          sql.NullTime
	MealPlanOptionArchivedAt                       sql.NullTime
	PrepTaskStepID                                 string
	PrepTaskExplicitStorageInstructions            string
	MealPlanOptionNotes                            string
	MealPlanOptionID                               string
	PrepTaskID                                     string
	PrepTaskName                                   string
	PrepTaskDescription                            string
	PrepTaskNotes                                  string
	CreationExplanation                            string
	MealPlanOptionMealID                           string
	PrepTaskStepBelongsToRecipePrepTask            string
	MealPlanOptionMealScale                        string
	StatusExplanation                              string
	Status                                         PrepStepStatus
	PrepTaskStepBelongsToRecipeStep                string
	PrepTaskBelongsToRecipe                        string
	ID                                             string
	PrepTaskMinimumStorageTemperatureInCelsius     sql.NullString
	MealPlanOptionAssignedCook                     sql.NullString
	PrepTaskStorageType                            NullStorageContainerType
	PrepTaskMaximumStorageTemperatureInCelsius     sql.NullString
	AssignedToUser                                 sql.NullString
	MealPlanOptionAssignedDishwasher               sql.NullString
	MealPlanOptionBelongsToMealPlanEvent           sql.NullString
	PrepTaskMaximumTimeBufferBeforeRecipeInSeconds sql.NullInt32
	PrepTaskMinimumTimeBufferBeforeRecipeInSeconds int32
	PrepTaskOptional                               bool
	PrepTaskStepSatisfiesRecipeStep                bool
	MealPlanOptionChosen                           bool
	MealPlanOptionTiebroken                        bool
}

func (q *Queries) GetMealPlanTask(ctx context.Context, db DBTX, mealPlanTaskID string) (*GetMealPlanTaskRow, error) {
	row := db.QueryRowContext(ctx, getMealPlanTask, mealPlanTaskID)
	var i GetMealPlanTaskRow
	err := row.Scan(
		&i.ID,
		&i.MealPlanOptionID,
		&i.MealPlanOptionAssignedCook,
		&i.MealPlanOptionAssignedDishwasher,
		&i.MealPlanOptionChosen,
		&i.MealPlanOptionTiebroken,
		&i.MealPlanOptionMealScale,
		&i.MealPlanOptionMealID,
		&i.MealPlanOptionNotes,
		&i.MealPlanOptionCreatedAt,
		&i.MealPlanOptionLastUpdatedAt,
		&i.MealPlanOptionArchivedAt,
		&i.MealPlanOptionBelongsToMealPlanEvent,
		&i.PrepTaskID,
		&i.PrepTaskName,
		&i.PrepTaskDescription,
		&i.PrepTaskNotes,
		&i.PrepTaskOptional,
		&i.PrepTaskExplicitStorageInstructions,
		&i.PrepTaskMinimumTimeBufferBeforeRecipeInSeconds,
		&i.PrepTaskMaximumTimeBufferBeforeRecipeInSeconds,
		&i.PrepTaskStorageType,
		&i.PrepTaskMinimumStorageTemperatureInCelsius,
		&i.PrepTaskMaximumStorageTemperatureInCelsius,
		&i.PrepTaskBelongsToRecipe,
		&i.PrepTaskCreatedAt,
		&i.PrepTaskLastUpdatedAt,
		&i.PrepTaskArchivedAt,
		&i.PrepTaskStepID,
		&i.PrepTaskStepBelongsToRecipeStep,
		&i.PrepTaskStepBelongsToRecipePrepTask,
		&i.PrepTaskStepSatisfiesRecipeStep,
		&i.CreatedAt,
		&i.LastUpdatedAt,
		&i.CompletedAt,
		&i.Status,
		&i.CreationExplanation,
		&i.StatusExplanation,
		&i.AssignedToUser,
	)
	return &i, err
}

const listAllMealPlanTasksByMealPlan = `-- name: ListAllMealPlanTasksByMealPlan :many

SELECT
    meal_plan_tasks.id,
    meal_plan_options.id as meal_plan_option_id,
    meal_plan_options.assigned_cook as meal_plan_option_assigned_cook,
    meal_plan_options.assigned_dishwasher as meal_plan_option_assigned_dishwasher,
    meal_plan_options.chosen as meal_plan_option_chosen,
    meal_plan_options.tiebroken as meal_plan_option_tiebroken,
    meal_plan_options.meal_scale as meal_plan_option_meal_scale,
    meal_plan_options.meal_id as meal_plan_option_meal_id,
    meal_plan_options.notes as meal_plan_option_notes,
    meal_plan_options.created_at as meal_plan_option_created_at,
    meal_plan_options.last_updated_at as meal_plan_option_last_updated_at,
    meal_plan_options.archived_at as meal_plan_option_archived_at,
    meal_plan_options.belongs_to_meal_plan_event as meal_plan_option_belongs_to_meal_plan_event,
    recipe_prep_tasks.id as prep_task_id,
    recipe_prep_tasks.name as prep_task_name,
    recipe_prep_tasks.description as prep_task_description,
    recipe_prep_tasks.notes as prep_task_notes,
    recipe_prep_tasks.optional as prep_task_optional,
    recipe_prep_tasks.explicit_storage_instructions as prep_task_explicit_storage_instructions,
    recipe_prep_tasks.minimum_time_buffer_before_recipe_in_seconds as prep_task_minimum_time_buffer_before_recipe_in_seconds,
    recipe_prep_tasks.maximum_time_buffer_before_recipe_in_seconds as prep_task_maximum_time_buffer_before_recipe_in_seconds,
    recipe_prep_tasks.storage_type as prep_task_storage_type,
    recipe_prep_tasks.minimum_storage_temperature_in_celsius as prep_task_minimum_storage_temperature_in_celsius,
    recipe_prep_tasks.maximum_storage_temperature_in_celsius as prep_task_maximum_storage_temperature_in_celsius,
    recipe_prep_tasks.belongs_to_recipe as prep_task_belongs_to_recipe,
    recipe_prep_tasks.created_at as prep_task_created_at,
    recipe_prep_tasks.last_updated_at as prep_task_last_updated_at,
    recipe_prep_tasks.archived_at as prep_task_archived_at,
    recipe_prep_task_steps.id as prep_task_step_id,
    recipe_prep_task_steps.belongs_to_recipe_step as prep_task_step_belongs_to_recipe_step,
    recipe_prep_task_steps.belongs_to_recipe_prep_task as prep_task_step_belongs_to_recipe_prep_task,
    recipe_prep_task_steps.satisfies_recipe_step as prep_task_step_satisfies_recipe_step,
    meal_plan_tasks.created_at,
    meal_plan_tasks.last_updated_at,
    meal_plan_tasks.completed_at,
    meal_plan_tasks.status,
    meal_plan_tasks.creation_explanation,
    meal_plan_tasks.status_explanation,
    meal_plan_tasks.assigned_to_user
FROM meal_plan_tasks
	JOIN meal_plan_options ON meal_plan_tasks.belongs_to_meal_plan_option=meal_plan_options.id
	JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event=meal_plan_events.id
	JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan=meal_plans.id
	JOIN meals ON meal_plan_options.meal_id=meals.id
	JOIN recipe_prep_tasks ON meal_plan_tasks.belongs_to_recipe_prep_task=recipe_prep_tasks.id
	JOIN recipe_prep_task_steps ON recipe_prep_task_steps.belongs_to_recipe_prep_task=recipe_prep_tasks.id
	JOIN recipe_steps ON recipe_prep_task_steps.belongs_to_recipe_step=recipe_steps.id
WHERE meal_plan_options.archived_at IS NULL
	AND meal_plan_events.archived_at IS NULL
	AND meal_plans.archived_at IS NULL
	AND meals.archived_at IS NULL
	AND recipe_steps.archived_at IS NULL
	AND meal_plans.id = $1
`

type ListAllMealPlanTasksByMealPlanRow struct {
	MealPlanOptionCreatedAt                        time.Time
	CreatedAt                                      time.Time
	PrepTaskCreatedAt                              time.Time
	PrepTaskArchivedAt                             sql.NullTime
	MealPlanOptionLastUpdatedAt                    sql.NullTime
	CompletedAt                                    sql.NullTime
	LastUpdatedAt                                  sql.NullTime
	PrepTaskLastUpdatedAt                          sql.NullTime
	MealPlanOptionArchivedAt                       sql.NullTime
	PrepTaskStepID                                 string
	PrepTaskExplicitStorageInstructions            string
	MealPlanOptionNotes                            string
	MealPlanOptionID                               string
	PrepTaskID                                     string
	PrepTaskName                                   string
	PrepTaskDescription                            string
	PrepTaskNotes                                  string
	CreationExplanation                            string
	MealPlanOptionMealID                           string
	PrepTaskStepBelongsToRecipePrepTask            string
	MealPlanOptionMealScale                        string
	StatusExplanation                              string
	Status                                         PrepStepStatus
	PrepTaskStepBelongsToRecipeStep                string
	PrepTaskBelongsToRecipe                        string
	ID                                             string
	PrepTaskMinimumStorageTemperatureInCelsius     sql.NullString
	MealPlanOptionAssignedCook                     sql.NullString
	PrepTaskStorageType                            NullStorageContainerType
	PrepTaskMaximumStorageTemperatureInCelsius     sql.NullString
	AssignedToUser                                 sql.NullString
	MealPlanOptionAssignedDishwasher               sql.NullString
	MealPlanOptionBelongsToMealPlanEvent           sql.NullString
	PrepTaskMaximumTimeBufferBeforeRecipeInSeconds sql.NullInt32
	PrepTaskMinimumTimeBufferBeforeRecipeInSeconds int32
	PrepTaskOptional                               bool
	PrepTaskStepSatisfiesRecipeStep                bool
	MealPlanOptionChosen                           bool
	MealPlanOptionTiebroken                        bool
}

func (q *Queries) ListAllMealPlanTasksByMealPlan(ctx context.Context, db DBTX, mealPlanID string) ([]*ListAllMealPlanTasksByMealPlanRow, error) {
	rows, err := db.QueryContext(ctx, listAllMealPlanTasksByMealPlan, mealPlanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*ListAllMealPlanTasksByMealPlanRow{}
	for rows.Next() {
		var i ListAllMealPlanTasksByMealPlanRow
		if err := rows.Scan(
			&i.ID,
			&i.MealPlanOptionID,
			&i.MealPlanOptionAssignedCook,
			&i.MealPlanOptionAssignedDishwasher,
			&i.MealPlanOptionChosen,
			&i.MealPlanOptionTiebroken,
			&i.MealPlanOptionMealScale,
			&i.MealPlanOptionMealID,
			&i.MealPlanOptionNotes,
			&i.MealPlanOptionCreatedAt,
			&i.MealPlanOptionLastUpdatedAt,
			&i.MealPlanOptionArchivedAt,
			&i.MealPlanOptionBelongsToMealPlanEvent,
			&i.PrepTaskID,
			&i.PrepTaskName,
			&i.PrepTaskDescription,
			&i.PrepTaskNotes,
			&i.PrepTaskOptional,
			&i.PrepTaskExplicitStorageInstructions,
			&i.PrepTaskMinimumTimeBufferBeforeRecipeInSeconds,
			&i.PrepTaskMaximumTimeBufferBeforeRecipeInSeconds,
			&i.PrepTaskStorageType,
			&i.PrepTaskMinimumStorageTemperatureInCelsius,
			&i.PrepTaskMaximumStorageTemperatureInCelsius,
			&i.PrepTaskBelongsToRecipe,
			&i.PrepTaskCreatedAt,
			&i.PrepTaskLastUpdatedAt,
			&i.PrepTaskArchivedAt,
			&i.PrepTaskStepID,
			&i.PrepTaskStepBelongsToRecipeStep,
			&i.PrepTaskStepBelongsToRecipePrepTask,
			&i.PrepTaskStepSatisfiesRecipeStep,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.CompletedAt,
			&i.Status,
			&i.CreationExplanation,
			&i.StatusExplanation,
			&i.AssignedToUser,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listIncompleteMealPlanTasksByMealPlanOption = `-- name: ListIncompleteMealPlanTasksByMealPlanOption :many

SELECT
	meal_plan_tasks.id,
	meal_plan_options.id,
	meal_plan_options.assigned_cook,
	meal_plan_options.assigned_dishwasher,
	meal_plan_options.chosen,
	meal_plan_options.tiebroken,
    meal_plan_options.meal_scale,
	meal_plan_options.meal_id,
	meal_plan_options.notes,
	meal_plan_options.created_at,
	meal_plan_options.last_updated_at,
	meal_plan_options.archived_at,
	meal_plan_options.belongs_to_meal_plan_event,
	recipe_steps.id,
	recipe_steps.index,
	valid_preparations.id,
	valid_preparations.name,
	valid_preparations.description,
	valid_preparations.icon_path,
	valid_preparations.yields_nothing,
	valid_preparations.restrict_to_ingredients,
	valid_preparations.minimum_ingredient_count,
	valid_preparations.maximum_ingredient_count,
	valid_preparations.minimum_instrument_count,
	valid_preparations.maximum_instrument_count,
	valid_preparations.temperature_required,
	valid_preparations.time_estimate_required,
	valid_preparations.condition_expression_required,
    valid_preparations.consumes_vessel,
    valid_preparations.only_for_vessels,
    valid_preparations.minimum_vessel_count,
    valid_preparations.maximum_vessel_count,
	valid_preparations.slug,
	valid_preparations.past_tense,
	valid_preparations.created_at,
	valid_preparations.last_updated_at,
	valid_preparations.archived_at,
	recipe_steps.minimum_estimated_time_in_seconds,
	recipe_steps.maximum_estimated_time_in_seconds,
	recipe_steps.minimum_temperature_in_celsius,
	recipe_steps.maximum_temperature_in_celsius,
	recipe_steps.notes,
	recipe_steps.explicit_instructions,
	recipe_steps.condition_expression,
	recipe_steps.optional,
	recipe_steps.start_timer_automatically,
	recipe_steps.created_at,
	recipe_steps.last_updated_at,
	recipe_steps.archived_at,
	recipe_steps.belongs_to_recipe,
	meal_plan_tasks.assigned_to_user,
	meal_plan_tasks.status,
	meal_plan_tasks.status_explanation,
	meal_plan_tasks.creation_explanation,
	meal_plan_tasks.created_at,
	meal_plan_tasks.completed_at
FROM meal_plan_tasks
	 FULL OUTER JOIN meal_plan_options ON meal_plan_tasks.belongs_to_meal_plan_option=meal_plan_options.id
	 FULL OUTER JOIN meal_plans ON meal_plan_options.belongs_to_meal_plan=meal_plans.id
	 FULL OUTER JOIN meals ON meal_plan_options.meal_id=meals.id
	 JOIN recipe_steps ON meal_plan_tasks.satisfies_recipe_step=recipe_steps.id
	 JOIN valid_preparations ON recipe_steps.preparation_id=valid_preparations.id
WHERE meal_plan_tasks.belongs_to_meal_plan_option = $1
AND meal_plan_tasks.completed_at IS NULL
`

type ListIncompleteMealPlanTasksByMealPlanOptionRow struct {
	CreatedAt_3                   time.Time
	CreatedAt_2                   time.Time
	ArchivedAt                    sql.NullTime
	CompletedAt                   sql.NullTime
	CreatedAt_4                   sql.NullTime
	ArchivedAt_3                  sql.NullTime
	LastUpdatedAt_3               sql.NullTime
	ArchivedAt_2                  sql.NullTime
	LastUpdatedAt_2               sql.NullTime
	CreatedAt                     sql.NullTime
	LastUpdatedAt                 sql.NullTime
	ExplicitInstructions          string
	BelongsToRecipe               string
	ID_3                          string
	PastTense                     string
	ID_4                          string
	Name                          string
	Description                   string
	IconPath                      string
	Notes_2                       string
	Slug                          string
	ConditionExpression           string
	AssignedCook                  sql.NullString
	StatusExplanation             sql.NullString
	MaximumTemperatureInCelsius   sql.NullString
	Status                        NullPrepStepStatus
	AssignedToUser                sql.NullString
	BelongsToMealPlanEvent        sql.NullString
	AssignedDishwasher            sql.NullString
	MealScale                     sql.NullString
	ID                            sql.NullString
	MinimumTemperatureInCelsius   sql.NullString
	CreationExplanation           sql.NullString
	ID_2                          sql.NullString
	Notes                         sql.NullString
	MealID                        sql.NullString
	MaximumEstimatedTimeInSeconds sql.NullInt64
	MinimumEstimatedTimeInSeconds sql.NullInt64
	MaximumIngredientCount        sql.NullInt32
	MaximumVesselCount            sql.NullInt32
	MaximumInstrumentCount        sql.NullInt32
	MinimumInstrumentCount        int32
	Index                         int32
	MinimumIngredientCount        int32
	MinimumVesselCount            int32
	Chosen                        sql.NullBool
	Tiebroken                     sql.NullBool
	ConsumesVessel                bool
	StartTimerAutomatically       bool
	ConditionExpressionRequired   bool
	TimeEstimateRequired          bool
	TemperatureRequired           bool
	Optional                      bool
	RestrictToIngredients         bool
	YieldsNothing                 bool
	OnlyForVessels                bool
}

func (q *Queries) ListIncompleteMealPlanTasksByMealPlanOption(ctx context.Context, db DBTX, belongsToMealPlanOption string) ([]*ListIncompleteMealPlanTasksByMealPlanOptionRow, error) {
	rows, err := db.QueryContext(ctx, listIncompleteMealPlanTasksByMealPlanOption, belongsToMealPlanOption)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*ListIncompleteMealPlanTasksByMealPlanOptionRow{}
	for rows.Next() {
		var i ListIncompleteMealPlanTasksByMealPlanOptionRow
		if err := rows.Scan(
			&i.ID,
			&i.ID_2,
			&i.AssignedCook,
			&i.AssignedDishwasher,
			&i.Chosen,
			&i.Tiebroken,
			&i.MealScale,
			&i.MealID,
			&i.Notes,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
			&i.BelongsToMealPlanEvent,
			&i.ID_3,
			&i.Index,
			&i.ID_4,
			&i.Name,
			&i.Description,
			&i.IconPath,
			&i.YieldsNothing,
			&i.RestrictToIngredients,
			&i.MinimumIngredientCount,
			&i.MaximumIngredientCount,
			&i.MinimumInstrumentCount,
			&i.MaximumInstrumentCount,
			&i.TemperatureRequired,
			&i.TimeEstimateRequired,
			&i.ConditionExpressionRequired,
			&i.ConsumesVessel,
			&i.OnlyForVessels,
			&i.MinimumVesselCount,
			&i.MaximumVesselCount,
			&i.Slug,
			&i.PastTense,
			&i.CreatedAt_2,
			&i.LastUpdatedAt_2,
			&i.ArchivedAt_2,
			&i.MinimumEstimatedTimeInSeconds,
			&i.MaximumEstimatedTimeInSeconds,
			&i.MinimumTemperatureInCelsius,
			&i.MaximumTemperatureInCelsius,
			&i.Notes_2,
			&i.ExplicitInstructions,
			&i.ConditionExpression,
			&i.Optional,
			&i.StartTimerAutomatically,
			&i.CreatedAt_3,
			&i.LastUpdatedAt_3,
			&i.ArchivedAt_3,
			&i.BelongsToRecipe,
			&i.AssignedToUser,
			&i.Status,
			&i.StatusExplanation,
			&i.CreationExplanation,
			&i.CreatedAt_4,
			&i.CompletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
