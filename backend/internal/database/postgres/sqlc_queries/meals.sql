-- name: ArchiveMeal :execrows

UPDATE meals SET archived_at = NOW() WHERE archived_at IS NULL AND created_by_user = sqlc.arg(created_by_user) AND id = sqlc.arg(id);

-- name: CreateMeal :exec

INSERT INTO meals (
	id,
	name,
	description,
	min_estimated_portions,
	max_estimated_portions,
	eligible_for_meal_plans,
	created_by_user
) VALUES (
	sqlc.arg(id),
	sqlc.arg(name),
	sqlc.arg(description),
	sqlc.arg(min_estimated_portions),
	sqlc.arg(max_estimated_portions),
	sqlc.arg(eligible_for_meal_plans),
	sqlc.arg(created_by_user)
);

-- name: CheckMealExistence :one

SELECT EXISTS (
	SELECT meals.id
	FROM meals
	WHERE meals.archived_at IS NULL
		AND meals.id = sqlc.arg(id)
);

-- name: GetMealsNeedingIndexing :many

SELECT meals.id
	FROM meals
	WHERE meals.archived_at IS NULL
	AND (
		meals.last_indexed_at IS NULL
		OR meals.last_indexed_at < NOW() - '24 hours'::INTERVAL
	);

-- name: GetMeal :many

SELECT
	meals.id,
	meals.name,
	meals.description,
	meals.min_estimated_portions,
	meals.max_estimated_portions,
	meals.eligible_for_meal_plans,
	meals.last_indexed_at,
	meals.created_at,
	meals.last_updated_at,
	meals.archived_at,
	meals.created_by_user,
	meal_components.id as component_id,
	meal_components.meal_id as component_meal_id,
	meal_components.recipe_id as component_recipe_id,
	meal_components.meal_component_type as component_meal_component_type,
	meal_components.recipe_scale as component_recipe_scale,
	meal_components.created_at as component_created_at,
	meal_components.last_updated_at as component_last_updated_at,
	meal_components.archived_at as component_archived_at
FROM meals
	JOIN meal_components ON meal_components.meal_id=meals.id
WHERE meals.archived_at IS NULL
  AND meal_components.archived_at IS NULL
  AND meals.id = sqlc.arg(id);

-- name: GetMeals :many

SELECT
	meals.id,
	meals.name,
	meals.description,
	meals.min_estimated_portions,
	meals.max_estimated_portions,
	meals.eligible_for_meal_plans,
	meals.last_indexed_at,
	meals.created_at,
	meals.last_updated_at,
	meals.archived_at,
	meals.created_by_user,
	(
		SELECT COUNT(meals.id)
		FROM meals
		WHERE meals.archived_at IS NULL
			AND meals.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
			AND meals.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
			AND (
				meals.last_updated_at IS NULL
				OR meals.last_updated_at > COALESCE(sqlc.narg(updated_before), (SELECT NOW() - '999 years'::INTERVAL))
			)
			AND (
				meals.last_updated_at IS NULL
				OR meals.last_updated_at < COALESCE(sqlc.narg(updated_after), (SELECT NOW() + '999 years'::INTERVAL))
			)
	) AS filtered_count,
	(
		SELECT COUNT(meals.id)
		FROM meals
		WHERE meals.archived_at IS NULL
	) AS total_count
FROM meals
WHERE
	meals.archived_at IS NULL
	AND meals.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
	AND meals.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
	AND (
		meals.last_updated_at IS NULL
		OR meals.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - '999 years'::INTERVAL))
	)
	AND (
		meals.last_updated_at IS NULL
		OR meals.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + '999 years'::INTERVAL))
	)
LIMIT sqlc.narg(query_limit)
OFFSET sqlc.narg(query_offset);

-- name: SearchForMeals :many

SELECT
	meals.id,
	meals.name,
	meals.description,
	meals.min_estimated_portions,
	meals.max_estimated_portions,
	meals.eligible_for_meal_plans,
	meals.last_indexed_at,
	meals.created_at,
	meals.last_updated_at,
	meals.archived_at,
	meals.created_by_user,
	meal_components.id as component_id,
	meal_components.meal_id as component_meal_id,
	meal_components.recipe_id as component_recipe_id,
	meal_components.meal_component_type as component_meal_component_type,
	meal_components.recipe_scale as component_recipe_scale,
	meal_components.created_at as component_created_at,
	meal_components.last_updated_at as component_last_updated_at,
	meal_components.archived_at as component_archived_at,
	(
		SELECT COUNT(meals.id)
		FROM meals
		WHERE meals.archived_at IS NULL
			AND meals.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
			AND meals.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
			AND (
				meals.last_updated_at IS NULL
				OR meals.last_updated_at > COALESCE(sqlc.narg(updated_before), (SELECT NOW() - '999 years'::INTERVAL))
			)
			AND (
				meals.last_updated_at IS NULL
				OR meals.last_updated_at < COALESCE(sqlc.narg(updated_after), (SELECT NOW() + '999 years'::INTERVAL))
			)
	) AS filtered_count,
	(
		SELECT COUNT(meals.id)
		FROM meals
		WHERE meals.archived_at IS NULL
	) AS total_count
FROM meals
	JOIN meal_components ON meal_components.meal_id=meals.id
WHERE
	meals.archived_at IS NULL
	AND meals.name ILIKE '%' || sqlc.arg(query)::text || '%'
	AND meals.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
	AND meals.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
	AND (
		meals.last_updated_at IS NULL
		OR meals.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - '999 years'::INTERVAL))
	)
	AND (
		meals.last_updated_at IS NULL
		OR meals.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + '999 years'::INTERVAL))
	)
LIMIT sqlc.narg(query_limit)
OFFSET sqlc.narg(query_offset);

-- name: UpdateMealLastIndexedAt :execrows

UPDATE meals SET last_indexed_at = NOW() WHERE id = sqlc.arg(id) AND archived_at IS NULL;
