// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: meal_plan_tasks_get_one.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const GetMealPlanTask = `-- name: GetMealPlanTask :exec
SELECT
    meal_plan_tasks.id,
    meal_plan_options.id,
    meal_plan_options.assigned_cook,
    meal_plan_options.assigned_dishwasher,
    meal_plan_options.chosen,
    meal_plan_options.tiebroken,
    meal_plan_options.meal_id,
    meal_plan_options.notes,
    meal_plan_options.created_at,
    meal_plan_options.last_updated_at,
    meal_plan_options.archived_at,
    meal_plan_options.belongs_to_meal_plan_event,
    recipe_prep_tasks.id,
    recipe_prep_tasks.notes,
    recipe_prep_tasks.explicit_storage_instructions,
    recipe_prep_tasks.minimum_time_buffer_before_recipe_in_seconds,
    recipe_prep_tasks.maximum_time_buffer_before_recipe_in_seconds,
    recipe_prep_tasks.storage_type,
    recipe_prep_tasks.minimum_storage_temperature_in_celsius,
    recipe_prep_tasks.maximum_storage_temperature_in_celsius,
    recipe_prep_tasks.belongs_to_recipe,
    recipe_prep_tasks.created_at,
    recipe_prep_tasks.last_updated_at,
    recipe_prep_tasks.archived_at,
    recipe_prep_task_steps.id,
    recipe_prep_task_steps.belongs_to_recipe_step,
    recipe_prep_task_steps.belongs_to_recipe_prep_task,
    recipe_prep_task_steps.satisfies_recipe_step,
    meal_plan_tasks.created_at,
    meal_plan_tasks.completed_at,
    meal_plan_tasks.status,
    meal_plan_tasks.creation_explanation,
    meal_plan_tasks.status_explanation,
    meal_plan_tasks.assigned_to_user
FROM meal_plan_tasks
    FULL OUTER JOIN meal_plan_options ON meal_plan_tasks.belongs_to_meal_plan_option=meal_plan_options.id
    FULL OUTER JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event=meal_plan_events.id
    FULL OUTER JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan=meal_plans.id
    FULL OUTER JOIN meals ON meal_plan_options.meal_id=meals.id
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
	CreatedAt                              time.Time                `db:"created_at"`
	CreatedAt_3                            time.Time                `db:"created_at_3"`
	CreatedAt_2                            time.Time                `db:"created_at_2"`
	LastUpdatedAt_2                        sql.NullTime             `db:"last_updated_at_2"`
	CompletedAt                            sql.NullTime             `db:"completed_at"`
	ArchivedAt_2                           sql.NullTime             `db:"archived_at_2"`
	LastUpdatedAt                          sql.NullTime             `db:"last_updated_at"`
	ArchivedAt                             sql.NullTime             `db:"archived_at"`
	Notes                                  string                   `db:"notes"`
	MealID                                 string                   `db:"meal_id"`
	CreationExplanation                    string                   `db:"creation_explanation"`
	StatusExplanation                      string                   `db:"status_explanation"`
	ID_3                                   string                   `db:"id_3"`
	ID_4                                   string                   `db:"id_4"`
	ExplicitStorageInstructions            string                   `db:"explicit_storage_instructions"`
	Status                                 PrepStepStatus           `db:"status"`
	BelongsToRecipeStep                    string                   `db:"belongs_to_recipe_step"`
	ID                                     string                   `db:"id"`
	ID_2                                   string                   `db:"id_2"`
	BelongsToRecipe                        string                   `db:"belongs_to_recipe"`
	Notes_2                                string                   `db:"notes_2"`
	BelongsToRecipePrepTask                string                   `db:"belongs_to_recipe_prep_task"`
	MaximumStorageTemperatureInCelsius     sql.NullString           `db:"maximum_storage_temperature_in_celsius"`
	MinimumStorageTemperatureInCelsius     sql.NullString           `db:"minimum_storage_temperature_in_celsius"`
	BelongsToMealPlanEvent                 sql.NullString           `db:"belongs_to_meal_plan_event"`
	AssignedDishwasher                     sql.NullString           `db:"assigned_dishwasher"`
	AssignedCook                           sql.NullString           `db:"assigned_cook"`
	StorageType                            NullStorageContainerType `db:"storage_type"`
	AssignedToUser                         sql.NullString           `db:"assigned_to_user"`
	MaximumTimeBufferBeforeRecipeInSeconds sql.NullInt32            `db:"maximum_time_buffer_before_recipe_in_seconds"`
	MinimumTimeBufferBeforeRecipeInSeconds int32                    `db:"minimum_time_buffer_before_recipe_in_seconds"`
	Tiebroken                              bool                     `db:"tiebroken"`
	Chosen                                 bool                     `db:"chosen"`
	SatisfiesRecipeStep                    bool                     `db:"satisfies_recipe_step"`
}

func (q *Queries) GetMealPlanTask(ctx context.Context, db DBTX, id string) error {
	_, err := db.ExecContext(ctx, GetMealPlanTask, id)
	return err
}