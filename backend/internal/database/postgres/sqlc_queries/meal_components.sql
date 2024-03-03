-- name: CreateMealComponent :exec

INSERT INTO meal_components (
	id,
	meal_id,
	recipe_id,
	meal_component_type,
	recipe_scale
) VALUES (
	sqlc.arg(id),
	sqlc.arg(meal_id),
	sqlc.arg(recipe_id),
	sqlc.arg(meal_component_type),
	sqlc.arg(recipe_scale)
);
