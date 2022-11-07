// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: recipe_step_instruments_exists.sql

package generated

import (
	"context"
)

const RecipeStepInstrumentExists = `-- name: RecipeStepInstrumentExists :exec
SELECT EXISTS ( SELECT recipe_step_instruments.id FROM recipe_step_instruments JOIN recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_instruments.archived_at IS NULL AND recipe_step_instruments.belongs_to_recipe_step = $1 AND recipe_step_instruments.id = $2 AND recipe_steps.archived_at IS NULL AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.archived_at IS NULL AND recipes.id = $5 )
`

type RecipeStepInstrumentExistsParams struct {
	BelongsToRecipeStep string `db:"belongs_to_recipe_step"`
	ID                  string `db:"id"`
	BelongsToRecipe     string `db:"belongs_to_recipe"`
	ID_2                string `db:"id_2"`
	ID_3                string `db:"id_3"`
}

func (q *Queries) RecipeStepInstrumentExists(ctx context.Context, db DBTX, arg *RecipeStepInstrumentExistsParams) error {
	_, err := db.ExecContext(ctx, RecipeStepInstrumentExists,
		arg.BelongsToRecipeStep,
		arg.ID,
		arg.BelongsToRecipe,
		arg.ID_2,
		arg.ID_3,
	)
	return err
}