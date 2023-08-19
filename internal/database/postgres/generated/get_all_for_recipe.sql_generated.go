// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: get_all_for_recipe.sql

package generated

import (
	"context"
	"database/sql"
)

const getAllRecipeStepCompletionConditionsForRecipe = `-- name: GetAllRecipeStepCompletionConditionsForRecipe :many

SELECT
	recipe_step_completion_condition_ingredients.id as recipe_step_completion_condition_ingredient_id,
	recipe_step_completion_condition_ingredients.belongs_to_recipe_step_completion_condition as recipe_step_completion_condition_ingredient_belongs_to_recipe_step_completion_condition,
	recipe_step_completion_condition_ingredients.recipe_step_ingredient as recipe_step_completion_condition_ingredient_recipe_step_ingredient,
	recipe_step_completion_conditions.id,
	recipe_step_completion_conditions.belongs_to_recipe_step,
	valid_ingredient_states.id as valid_ingredient_state_id,
	valid_ingredient_states.name as valid_ingredient_state_name,
	valid_ingredient_states.description as valid_ingredient_state_description,
	valid_ingredient_states.icon_path as valid_ingredient_state_icon_path,
	valid_ingredient_states.slug as valid_ingredient_state_slug,
	valid_ingredient_states.past_tense as valid_ingredient_state_past_tense,
	valid_ingredient_states.attribute_type as valid_ingredient_state_attribute_type,
	valid_ingredient_states.created_at as valid_ingredient_state_created_at,
	valid_ingredient_states.last_updated_at as valid_ingredient_state_last_updated_at,
	valid_ingredient_states.archived_at as valid_ingredient_state_archived_at,
	recipe_step_completion_conditions.optional,
	recipe_step_completion_conditions.notes,
	recipe_step_completion_conditions.created_at,
	recipe_step_completion_conditions.last_updated_at,
	recipe_step_completion_conditions.archived_at
FROM recipe_step_completion_condition_ingredients
	LEFT JOIN recipe_step_completion_conditions ON recipe_step_completion_condition_ingredients.belongs_to_recipe_step_completion_condition = recipe_step_completion_conditions.id
	LEFT JOIN recipe_steps ON recipe_step_completion_conditions.belongs_to_recipe_step = recipe_steps.id
	LEFT JOIN recipes ON recipe_steps.belongs_to_recipe = recipes.id
	LEFT JOIN valid_ingredient_states ON recipe_step_completion_conditions.ingredient_state = valid_ingredient_states.id
WHERE recipe_step_completion_conditions.archived_at IS NULL
	AND recipe_step_completion_condition_ingredients.archived_at IS NULL
    AND recipe_steps.archived_at IS NULL
	AND recipes.archived_at IS NULL
    AND valid_ingredient_states.archived_at IS NULL
	AND recipes.id = $1
GROUP BY recipe_step_completion_conditions.id,
	     recipe_step_completion_condition_ingredients.id,
	     valid_ingredient_states.id
`

type GetAllRecipeStepCompletionConditionsForRecipeRow struct {
	RecipeStepCompletionConditionIngredientID                string
	RecipeStepCompletionConditionIngredientBelongsToRecipeS  string
	RecipeStepCompletionConditionIngredientRecipeStepIngredi string
	ID                                                       sql.NullString
	BelongsToRecipeStep                                      sql.NullString
	ValidIngredientStateID                                   sql.NullString
	ValidIngredientStateName                                 sql.NullString
	ValidIngredientStateDescription                          sql.NullString
	ValidIngredientStateIconPath                             sql.NullString
	ValidIngredientStateSlug                                 sql.NullString
	ValidIngredientStatePastTense                            sql.NullString
	ValidIngredientStateAttributeType                        NullIngredientAttributeType
	ValidIngredientStateCreatedAt                            sql.NullTime
	ValidIngredientStateLastUpdatedAt                        sql.NullTime
	ValidIngredientStateArchivedAt                           sql.NullTime
	Optional                                                 sql.NullBool
	Notes                                                    sql.NullString
	CreatedAt                                                sql.NullTime
	LastUpdatedAt                                            sql.NullTime
	ArchivedAt                                               sql.NullTime
}

func (q *Queries) GetAllRecipeStepCompletionConditionsForRecipe(ctx context.Context, db DBTX, id string) ([]*GetAllRecipeStepCompletionConditionsForRecipeRow, error) {
	rows, err := db.QueryContext(ctx, getAllRecipeStepCompletionConditionsForRecipe, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetAllRecipeStepCompletionConditionsForRecipeRow{}
	for rows.Next() {
		var i GetAllRecipeStepCompletionConditionsForRecipeRow
		if err := rows.Scan(
			&i.RecipeStepCompletionConditionIngredientID,
			&i.RecipeStepCompletionConditionIngredientBelongsToRecipeS,
			&i.RecipeStepCompletionConditionIngredientRecipeStepIngredi,
			&i.ID,
			&i.BelongsToRecipeStep,
			&i.ValidIngredientStateID,
			&i.ValidIngredientStateName,
			&i.ValidIngredientStateDescription,
			&i.ValidIngredientStateIconPath,
			&i.ValidIngredientStateSlug,
			&i.ValidIngredientStatePastTense,
			&i.ValidIngredientStateAttributeType,
			&i.ValidIngredientStateCreatedAt,
			&i.ValidIngredientStateLastUpdatedAt,
			&i.ValidIngredientStateArchivedAt,
			&i.Optional,
			&i.Notes,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
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
