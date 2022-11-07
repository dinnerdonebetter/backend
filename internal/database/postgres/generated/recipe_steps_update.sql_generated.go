// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: recipe_steps_update.sql

package generated

import (
	"context"
	"database/sql"
)

const UpdateRecipeStep = `-- name: UpdateRecipeStep :exec
UPDATE recipe_steps SET
	index = $1,
	preparation_id = $2,
	minimum_estimated_time_in_seconds = $3,
	maximum_estimated_time_in_seconds = $4,
	minimum_temperature_in_celsius = $5,
	maximum_temperature_in_celsius = $6,
	notes = $7,
	explicit_instructions = $8,
	optional = $9,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_recipe = $10
	AND id = $11
`

type UpdateRecipeStepParams struct {
	Notes                         string         `db:"notes"`
	PreparationID                 string         `db:"preparation_id"`
	BelongsToRecipe               string         `db:"belongs_to_recipe"`
	ID                            string         `db:"id"`
	ExplicitInstructions          string         `db:"explicit_instructions"`
	MaximumTemperatureInCelsius   sql.NullString `db:"maximum_temperature_in_celsius"`
	MinimumTemperatureInCelsius   sql.NullString `db:"minimum_temperature_in_celsius"`
	MinimumEstimatedTimeInSeconds sql.NullInt64  `db:"minimum_estimated_time_in_seconds"`
	MaximumEstimatedTimeInSeconds sql.NullInt64  `db:"maximum_estimated_time_in_seconds"`
	Index                         int32          `db:"index"`
	Optional                      bool           `db:"optional"`
}

func (q *Queries) UpdateRecipeStep(ctx context.Context, db DBTX, arg *UpdateRecipeStepParams) error {
	_, err := db.ExecContext(ctx, UpdateRecipeStep,
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
		arg.ID,
	)
	return err
}