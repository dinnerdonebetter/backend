// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: recipe_steps.sql

package generated

import (
	"context"
	"database/sql"
)

const ArchiveRecipeStep = `-- name: ArchiveRecipeStep :exec
UPDATE recipe_steps SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe = $1 AND id = $2
`

type ArchiveRecipeStepParams struct {
	BelongsToRecipe string
	ID              string
}

func (q *Queries) ArchiveRecipeStep(ctx context.Context, arg *ArchiveRecipeStepParams) error {
	_, err := q.db.ExecContext(ctx, ArchiveRecipeStep, arg.BelongsToRecipe, arg.ID)
	return err
}

const CreateRecipeStep = `-- name: CreateRecipeStep :exec
INSERT INTO recipe_steps (id,index,preparation_id,minimum_estimated_time_in_seconds,maximum_estimated_time_in_seconds,minimum_temperature_in_celsius,maximum_temperature_in_celsius,notes,explicit_instructions,optional,belongs_to_recipe) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
`

type CreateRecipeStepParams struct {
	ID                            string
	Index                         int32
	PreparationID                 string
	MinimumEstimatedTimeInSeconds int64
	MaximumEstimatedTimeInSeconds int64
	MinimumTemperatureInCelsius   sql.NullInt32
	MaximumTemperatureInCelsius   sql.NullInt32
	Notes                         string
	ExplicitInstructions          string
	Optional                      bool
	BelongsToRecipe               string
}

func (q *Queries) CreateRecipeStep(ctx context.Context, arg *CreateRecipeStepParams) error {
	_, err := q.db.ExecContext(ctx, CreateRecipeStep,
		arg.ID,
		arg.Index,
		arg.PreparationID,
		arg.MinimumEstimatedTimeInSeconds,
		arg.MaximumEstimatedTimeInSeconds,
		arg.MinimumTemperatureInCelsius,
		arg.MaximumTemperatureInCelsius,
		arg.Notes,
		arg.ExplicitInstructions,
		arg.Optional,
		arg.BelongsToRecipe,
	)
	return err
}

const GetRecipeStep = `-- name: GetRecipeStep :one
SELECT
    recipe_steps.id,
    recipe_steps.index,
    valid_preparations.id,
    valid_preparations.name,
    valid_preparations.description,
    valid_preparations.icon_path,
    valid_preparations.created_on,
    valid_preparations.last_updated_on,
    valid_preparations.archived_on,
    recipe_steps.minimum_estimated_time_in_seconds,
    recipe_steps.maximum_estimated_time_in_seconds,
    recipe_steps.minimum_temperature_in_celsius,
    recipe_steps.maximum_temperature_in_celsius,
    recipe_steps.explicit_instructions,
    recipe_steps.notes,
    recipe_steps.optional,
    recipe_steps.created_on,
    recipe_steps.last_updated_on,
    recipe_steps.archived_on,
    recipe_steps.belongs_to_recipe
FROM recipe_steps
         JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
         JOIN valid_preparations ON recipe_steps.preparation_id=valid_preparations.id
WHERE recipe_steps.archived_on IS NULL
  AND recipe_steps.belongs_to_recipe = $1
  AND recipe_steps.id = $2
  AND recipes.archived_on IS NULL
  AND recipes.id = $3
`

type GetRecipeStepParams struct {
	BelongsToRecipe string
	ID              string
	ID_2            string
}

type GetRecipeStepRow struct {
	ID                            string
	Index                         int32
	ID_2                          string
	Name                          string
	Description                   string
	IconPath                      string
	CreatedOn                     int64
	LastUpdatedOn                 sql.NullInt64
	ArchivedOn                    sql.NullInt64
	MinimumEstimatedTimeInSeconds int64
	MaximumEstimatedTimeInSeconds int64
	MinimumTemperatureInCelsius   sql.NullInt32
	MaximumTemperatureInCelsius   sql.NullInt32
	ExplicitInstructions          string
	Notes                         string
	Optional                      bool
	CreatedOn_2                   int64
	LastUpdatedOn_2               sql.NullInt64
	ArchivedOn_2                  sql.NullInt64
	BelongsToRecipe               string
}

func (q *Queries) GetRecipeStep(ctx context.Context, arg *GetRecipeStepParams) (*GetRecipeStepRow, error) {
	row := q.db.QueryRowContext(ctx, GetRecipeStep, arg.BelongsToRecipe, arg.ID, arg.ID_2)
	var i GetRecipeStepRow
	err := row.Scan(
		&i.ID,
		&i.Index,
		&i.ID_2,
		&i.Name,
		&i.Description,
		&i.IconPath,
		&i.CreatedOn,
		&i.LastUpdatedOn,
		&i.ArchivedOn,
		&i.MinimumEstimatedTimeInSeconds,
		&i.MaximumEstimatedTimeInSeconds,
		&i.MinimumTemperatureInCelsius,
		&i.MaximumTemperatureInCelsius,
		&i.ExplicitInstructions,
		&i.Notes,
		&i.Optional,
		&i.CreatedOn_2,
		&i.LastUpdatedOn_2,
		&i.ArchivedOn_2,
		&i.BelongsToRecipe,
	)
	return &i, err
}

const GetTotalRecipeStepsCount = `-- name: GetTotalRecipeStepsCount :one
SELECT COUNT(recipe_steps.id) FROM recipe_steps WHERE recipe_steps.archived_on IS NULL
`

func (q *Queries) GetTotalRecipeStepsCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, GetTotalRecipeStepsCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const RecipeStepExists = `-- name: RecipeStepExists :one
SELECT EXISTS ( SELECT recipe_steps.id FROM recipe_steps JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $1 AND recipe_steps.id = $2 AND recipes.archived_on IS NULL AND recipes.id = $3 )
`

type RecipeStepExistsParams struct {
	BelongsToRecipe string
	ID              string
	ID_2            string
}

func (q *Queries) RecipeStepExists(ctx context.Context, arg *RecipeStepExistsParams) (bool, error) {
	row := q.db.QueryRowContext(ctx, RecipeStepExists, arg.BelongsToRecipe, arg.ID, arg.ID_2)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const UpdateRecipeStep = `-- name: UpdateRecipeStep :exec
UPDATE recipe_steps SET
    index = $1,
    preparation_id = $2,
    minimum_estimated_time_in_seconds = $3,
    maximum_estimated_time_in_seconds = $4,
    minimum_temperature_in_celsius = $5,
    maximum_temperature_in_celsius = $6,
    explicit_instructions = $7,
    notes = $8,
    optional = $9,
    last_updated_on = extract(epoch FROM NOW())
WHERE archived_on IS NULL
  AND belongs_to_recipe = $10
  AND id = $11
`

type UpdateRecipeStepParams struct {
	Index                         int32
	PreparationID                 string
	MinimumEstimatedTimeInSeconds int64
	MaximumEstimatedTimeInSeconds int64
	MinimumTemperatureInCelsius   sql.NullInt32
	MaximumTemperatureInCelsius   sql.NullInt32
	ExplicitInstructions          string
	Notes                         string
	Optional                      bool
	BelongsToRecipe               string
	ID                            string
}

func (q *Queries) UpdateRecipeStep(ctx context.Context, arg *UpdateRecipeStepParams) error {
	_, err := q.db.ExecContext(ctx, UpdateRecipeStep,
		arg.Index,
		arg.PreparationID,
		arg.MinimumEstimatedTimeInSeconds,
		arg.MaximumEstimatedTimeInSeconds,
		arg.MinimumTemperatureInCelsius,
		arg.MaximumTemperatureInCelsius,
		arg.ExplicitInstructions,
		arg.Notes,
		arg.Optional,
		arg.BelongsToRecipe,
		arg.ID,
	)
	return err
}
