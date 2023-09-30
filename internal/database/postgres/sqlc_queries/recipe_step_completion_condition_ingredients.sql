-- name: CreateRecipeStepCompletionConditionIngredient :exec

INSERT INTO recipe_step_completion_condition_ingredients (
	id,
	belongs_to_recipe_step_completion_condition,
	recipe_step_ingredient
) VALUES (
	sqlc.arg(id),
	sqlc.arg(belongs_to_recipe_step_completion_condition),
	sqlc.arg(recipe_step_ingredient)
);

-- name: GetAllRecipeStepCompletionConditionIngredientsForRecipeCompletionIDs :many

SELECT
	recipe_step_completion_condition_ingredients.id as recipe_step_completion_condition_ingredient_id,
	recipe_step_completion_condition_ingredients.belongs_to_recipe_step_completion_condition as recipe_step_completion_condition_ingredient_belongs_to_recipe_step_completion_condition,
	recipe_step_completion_condition_ingredients.recipe_step_ingredient as recipe_step_completion_condition_ingredient_recipe_step_ingredient,
	recipe_step_completion_condition_ingredients.created_at as recipe_step_completion_condition_ingredient_created_at,
	recipe_step_completion_condition_ingredients.last_updated_at as recipe_step_completion_condition_ingredient_last_updated_at,
	recipe_step_completion_condition_ingredients.archived_at as recipe_step_completion_condition_ingredient_archived_at,
	valid_ingredient_states.id as valid_ingredient_state_id,
	valid_ingredient_states.name as valid_ingredient_state_name,
	valid_ingredient_states.past_tense as valid_ingredient_state_past_tense,
	valid_ingredient_states.slug as valid_ingredient_state_slug,
	valid_ingredient_states.description as valid_ingredient_state_description,
	valid_ingredient_states.icon_path as valid_ingredient_state_icon_path,
	valid_ingredient_states.attribute_type as valid_ingredient_state_attribute_type,
	valid_ingredient_states.last_indexed_at as valid_ingredient_state_last_indexed_at,
	valid_ingredient_states.created_at as valid_ingredient_state_created_at,
	valid_ingredient_states.last_updated_at as valid_ingredient_state_last_updated_at,
	valid_ingredient_states.archived_at as valid_ingredient_state_archived_at
FROM recipe_step_completion_condition_ingredients
	JOIN recipe_step_completion_conditions ON recipe_step_completion_condition_ingredients.belongs_to_recipe_step_completion_condition = recipe_step_completion_conditions.id
	JOIN valid_ingredient_states ON recipe_step_completion_conditions.ingredient_state = valid_ingredient_states.id
WHERE recipe_step_completion_conditions.archived_at IS NULL
	AND recipe_step_completion_condition_ingredients.archived_at IS NULL
	AND recipe_step_completion_condition_ingredients.belongs_to_recipe_step_completion_condition = ANY(sqlc.arg(ids)::text[])
	AND valid_ingredient_states.archived_at IS NULL;
