-- name: ArchiveRecipeStepCompletionCondition :execrows

UPDATE recipe_step_completion_conditions SET archived_at = NOW() WHERE archived_at IS NULL AND belongs_to_recipe_step = $1 AND id = $2;

-- name: CreateRecipeStepCompletionCondition :exec

INSERT INTO recipe_step_completion_conditions (
	id,
	belongs_to_recipe_step,
	ingredient_state,
	optional,
	notes
) VALUES ($1,$2,$3,$4,$5);

-- name: CheckRecipeStepCompletionConditionExistence :one

SELECT EXISTS ( SELECT recipe_step_completion_conditions.id FROM recipe_step_completion_conditions JOIN recipe_steps ON recipe_step_completion_conditions.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_completion_conditions.archived_at IS NULL AND recipe_step_completion_conditions.belongs_to_recipe_step = sqlc.arg(recipe_step_id) AND recipe_step_completion_conditions.id = sqlc.arg(recipe_step_completion_condition_id) AND recipe_steps.archived_at IS NULL AND recipe_steps.belongs_to_recipe = sqlc.arg(recipe_id) AND recipe_steps.id = sqlc.arg(recipe_step_id) AND recipes.archived_at IS NULL AND recipes.id = sqlc.arg(recipe_id) );

-- name: GetAllRecipeStepCompletionConditionsForRecipe :many

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
    JOIN recipe_step_completion_conditions ON recipe_step_completion_condition_ingredients.belongs_to_recipe_step_completion_condition = recipe_step_completion_conditions.id
    JOIN recipe_steps ON recipe_step_completion_conditions.belongs_to_recipe_step = recipe_steps.id
    JOIN recipes ON recipe_steps.belongs_to_recipe = recipes.id
    JOIN valid_ingredient_states ON recipe_step_completion_conditions.ingredient_state = valid_ingredient_states.id
WHERE recipe_step_completion_conditions.archived_at IS NULL
    AND recipe_step_completion_condition_ingredients.archived_at IS NULL
    AND recipe_steps.archived_at IS NULL
    AND recipes.archived_at IS NULL
    AND valid_ingredient_states.archived_at IS NULL
    AND recipes.id = $1
    GROUP BY recipe_step_completion_conditions.id,
         recipe_step_completion_condition_ingredients.id,
         valid_ingredient_states.id;

-- name: GetRecipeStepCompletionConditions :many

SELECT
    recipe_step_completion_condition_ingredients.id as recipe_step_completion_condition_ingredient_id,
    recipe_step_completion_condition_ingredients.belongs_to_recipe_step_completion_condition as recipe_step_completion_condition_ingredient_belongs_to_recipe_step_completion_condition,
    recipe_step_completion_condition_ingredients.recipe_step_ingredient as recipe_step_completion_condition_ingredient_recipe_step_ingredient,
    recipe_step_completion_condition_ingredients.created_at as recipe_step_completion_condition_ingredient_created_at,
    recipe_step_completion_condition_ingredients.last_updated_at as recipe_step_completion_condition_ingredient_last_updated_at,
    recipe_step_completion_condition_ingredients.archived_at as recipe_step_completion_condition_ingredient_archived_at,
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
    recipe_step_completion_conditions.archived_at,
    (
        SELECT
            COUNT(recipe_step_completion_conditions.id)
        FROM
            recipe_step_completion_conditions
        WHERE
            recipe_step_completion_conditions.archived_at IS NULL
            AND recipe_step_completion_conditions.belongs_to_recipe_step = sqlc.arg(recipe_step_id)
            AND recipe_step_completion_conditions.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
            AND recipe_step_completion_conditions.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
            AND (recipe_step_completion_conditions.last_updated_at IS NULL OR recipe_step_completion_conditions.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
            AND (recipe_step_completion_conditions.last_updated_at IS NULL OR recipe_step_completion_conditions.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
    ) as filtered_count,
    (
        SELECT
            COUNT(recipe_step_completion_conditions.id)
        FROM
            recipe_step_completion_conditions
        WHERE
            recipe_step_completion_conditions.archived_at IS NULL
    ) as total_count
FROM recipe_step_completion_condition_ingredients
    JOIN recipe_step_completion_conditions ON recipe_step_completion_condition_ingredients.belongs_to_recipe_step_completion_condition = recipe_step_completion_conditions.id
    JOIN recipe_steps ON recipe_step_completion_conditions.belongs_to_recipe_step = recipe_steps.id
    JOIN valid_ingredient_states ON recipe_step_completion_conditions.ingredient_state = valid_ingredient_states.id
WHERE recipe_step_completion_conditions.archived_at IS NULL
    AND recipe_step_completion_conditions.belongs_to_recipe_step = sqlc.arg(recipe_step_id)
    AND recipe_step_completion_conditions.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
    AND recipe_step_completion_conditions.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
    AND (recipe_step_completion_conditions.last_updated_at IS NULL OR recipe_step_completion_conditions.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
    AND (recipe_step_completion_conditions.last_updated_at IS NULL OR recipe_step_completion_conditions.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
    OFFSET sqlc.narg(query_offset)
    LIMIT sqlc.narg(query_limit);


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
	recipe_step_completion_conditions.archived_at,
	(
	    SELECT
	        COUNT(recipe_step_completion_conditions.id)
	    FROM
	        recipe_step_completion_conditions
	    WHERE
            recipe_step_completion_conditions.archived_at IS NULL
            AND recipe_step_completion_conditions.belongs_to_recipe_step = sqlc.arg(recipe_step_id)
            AND recipe_step_completion_conditions.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
            AND recipe_step_completion_conditions.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
            AND (recipe_step_completion_conditions.last_updated_at IS NULL OR recipe_step_completion_conditions.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
            AND (recipe_step_completion_conditions.last_updated_at IS NULL OR recipe_step_completion_conditions.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
	) as filtered_count,
	(
	    SELECT
	        COUNT(recipe_step_completion_conditions.id)
	    FROM
	        recipe_step_completion_conditions
	    WHERE
	        recipe_step_completion_conditions.archived_at IS NULL
	) as total_count
FROM recipe_step_completion_condition_ingredients
	JOIN recipe_step_completion_conditions ON recipe_step_completion_condition_ingredients.belongs_to_recipe_step_completion_condition = recipe_step_completion_conditions.id
	JOIN recipe_steps ON recipe_step_completion_conditions.belongs_to_recipe_step = recipe_steps.id
	JOIN recipes ON recipe_steps.belongs_to_recipe = recipes.id
	JOIN valid_ingredient_states ON recipe_step_completion_conditions.ingredient_state = valid_ingredient_states.id
WHERE
    recipe_step_completion_conditions.archived_at IS NULL
    AND recipe_step_completion_conditions.belongs_to_recipe_step = sqlc.arg(recipe_step_id)
	AND recipe_step_completion_conditions.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
	AND recipe_step_completion_conditions.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
	AND (recipe_step_completion_conditions.last_updated_at IS NULL OR recipe_step_completion_conditions.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
	AND (recipe_step_completion_conditions.last_updated_at IS NULL OR recipe_step_completion_conditions.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
GROUP BY
	recipe_step_completion_conditions.id,
	valid_ingredient_states.id,
	recipe_step_completion_condition_ingredients.id
ORDER BY
	recipe_step_completion_conditions.id
	LIMIT sqlc.narg(query_limit)
	OFFSET sqlc.narg(query_offset);

-- name: GetRecipeStepCompletionConditionWithIngredients :many

SELECT
    recipe_step_completion_condition_ingredients.id as recipe_step_completion_condition_ingredient_id,
    recipe_step_completion_condition_ingredients.belongs_to_recipe_step_completion_condition as recipe_step_completion_condition_ingredient_belongs_to_recipe_step_completion_condition,
    recipe_step_completion_condition_ingredients.recipe_step_ingredient as recipe_step_completion_condition_ingredient_recipe_step_ingredient,
    recipe_step_completion_condition_ingredients.created_at as recipe_step_completion_condition_ingredient_created_at,
    recipe_step_completion_condition_ingredients.last_updated_at as recipe_step_completion_condition_ingredient_last_updated_at,
    recipe_step_completion_condition_ingredients.archived_at as recipe_step_completion_condition_ingredient_archived_at,
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
    JOIN recipe_step_completion_conditions ON recipe_step_completion_condition_ingredients.belongs_to_recipe_step_completion_condition = recipe_step_completion_conditions.id
    JOIN recipe_steps ON recipe_step_completion_conditions.belongs_to_recipe_step = recipe_steps.id
    JOIN valid_ingredient_states ON recipe_step_completion_conditions.ingredient_state = valid_ingredient_states.id
WHERE recipe_step_completion_conditions.archived_at IS NULL
    AND recipe_step_completion_condition_ingredients.archived_at IS NULL
    AND recipe_step_completion_conditions.belongs_to_recipe_step = sqlc.arg(recipe_step_id)
    AND recipe_step_completion_conditions.id = sqlc.arg(recipe_step_completion_condition_id)
    AND recipe_steps.archived_at IS NULL
    AND recipe_steps.belongs_to_recipe = sqlc.arg(recipe_id)
    AND recipe_steps.id = sqlc.arg(recipe_step_id);

-- name: UpdateRecipeStepCompletionCondition :execrows

UPDATE recipe_step_completion_conditions
SET
	optional = $1,
	notes = $2,
	belongs_to_recipe_step = $3,
	ingredient_state = $4,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $5;