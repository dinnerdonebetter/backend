-- name: ArchiveMeal :execrows

UPDATE meals SET archived_at = NOW() WHERE archived_at IS NULL AND created_by_user = $1 AND id = $2;

-- name: CreateMeal :exec

INSERT INTO meals (id,"name",description,min_estimated_portions,max_estimated_portions,eligible_for_meal_plans,created_by_user) VALUES ($1,$2,$3,$4,$5,$6,$7);

-- name: CheckMealExistence :one

SELECT EXISTS ( SELECT meals.id FROM meals WHERE meals.archived_at IS NULL AND meals.id = $1 );

-- name: GetMealsNeedingIndexing :many

SELECT meals.id
    FROM meals
    WHERE (meals.archived_at IS NULL)
    AND (
        (meals.last_indexed_at IS NULL)
        OR meals.last_indexed_at
            < now() - '24 hours'::INTERVAL
    );

-- name: GetMeal :many

SELECT
    meals.id,
    meals.name,
    meals.description,
    meals.min_estimated_portions,
    meals.max_estimated_portions,
    meals.eligible_for_meal_plans,
    meals.created_at,
    meals.last_updated_at,
    meals.archived_at,
    meals.created_by_user,
    meal_components.recipe_id as component_recipe_id,
    meal_components.recipe_scale as component_recipe_scale,
    meal_components.meal_component_type as component_meal_component_type
FROM meals
    JOIN meal_components ON meal_components.meal_id=meals.id
WHERE meals.archived_at IS NULL
  AND meal_components.archived_at IS NULL
  AND meals.id = $1;

-- name: GetMeals :many

SELECT
    meals.id,
    meals.name,
    meals.description,
    meals.min_estimated_portions,
    meals.max_estimated_portions,
    meals.eligible_for_meal_plans,
    meals.created_at,
    meals.last_updated_at,
    meals.archived_at,
    meals.created_by_user,
    (
        SELECT
            COUNT(meals.id)
        FROM
            meals
        WHERE
            meals.archived_at IS NULL
            AND meals.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
            AND meals.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
            AND (meals.last_updated_at IS NULL OR meals.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
            AND (meals.last_updated_at IS NULL OR meals.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
    ) as filtered_count,
    (
        SELECT
            COUNT(meals.id)
        FROM
            meals
        WHERE
            meals.archived_at IS NULL
    ) as total_count
FROM meals
WHERE meals.archived_at IS NULL
    AND meals.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
    AND meals.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
    AND (meals.last_updated_at IS NULL OR meals.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
    AND (meals.last_updated_at IS NULL OR meals.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
    OFFSET sqlc.narg(query_offset)
    LIMIT sqlc.narg(query_limit);

-- name: SearchForMeals :many

SELECT
    meals.id,
    meals.name,
    meals.description,
    meals.min_estimated_portions,
    meals.max_estimated_portions,
    meals.eligible_for_meal_plans,
    meals.created_at,
    meals.last_updated_at,
    meals.archived_at,
    meals.created_by_user,
    meal_components.recipe_id as component_recipe_id,
    meal_components.recipe_scale as component_recipe_scale,
    meal_components.meal_component_type as component_meal_component_type,
    (
        SELECT
            COUNT(meals.id)
        FROM
            meals
        WHERE
            meals.archived_at IS NULL
            AND meals.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
            AND meals.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
            AND (meals.last_updated_at IS NULL OR meals.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
            AND (meals.last_updated_at IS NULL OR meals.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
    ) as filtered_count,
    (
        SELECT
            COUNT(meals.id)
        FROM
            meals
        WHERE
            meals.archived_at IS NULL
    ) as total_count
FROM meals
    JOIN meal_components ON meal_components.meal_id=meals.id
WHERE meals.archived_at IS NULL
    AND meals.name ILIKE '%' || sqlc.arg(query) || '%'
    AND meals.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
    AND meals.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
    AND (meals.last_updated_at IS NULL OR meals.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
    AND (meals.last_updated_at IS NULL OR meals.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
    AND meal_components.archived_at IS NULL
    OFFSET sqlc.narg(query_offset)
    LIMIT sqlc.narg(query_limit);

-- name: UpdateMealLastIndexedAt :execrows

UPDATE meals SET last_indexed_at = NOW() WHERE id = $1 AND archived_at IS NULL;