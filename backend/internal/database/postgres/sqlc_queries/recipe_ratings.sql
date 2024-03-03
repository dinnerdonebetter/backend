-- name: ArchiveRecipeRating :execrows

UPDATE recipe_ratings SET archived_at = NOW() WHERE archived_at IS NULL AND id = sqlc.arg(id);

-- name: CreateRecipeRating :exec

INSERT INTO recipe_ratings (
	id,
	recipe_id,
	taste,
	difficulty,
	cleanup,
	instructions,
	overall,
	notes,
	by_user
) VALUES (
	sqlc.arg(id),
	sqlc.arg(recipe_id),
	sqlc.arg(taste),
	sqlc.arg(difficulty),
	sqlc.arg(cleanup),
	sqlc.arg(instructions),
	sqlc.arg(overall),
	sqlc.arg(notes),
	sqlc.arg(by_user)
);

-- name: CheckRecipeRatingExistence :one

SELECT EXISTS (
	SELECT recipe_ratings.id
	FROM recipe_ratings
	WHERE recipe_ratings.archived_at IS NULL
		AND recipe_ratings.id = sqlc.arg(id)
);

-- name: GetRecipeRatings :many

SELECT
	recipe_ratings.id,
	recipe_ratings.recipe_id,
	recipe_ratings.taste,
	recipe_ratings.difficulty,
	recipe_ratings.cleanup,
	recipe_ratings.instructions,
	recipe_ratings.overall,
	recipe_ratings.notes,
	recipe_ratings.by_user,
	recipe_ratings.created_at,
	recipe_ratings.last_updated_at,
	recipe_ratings.archived_at,
	(
		SELECT COUNT(recipe_ratings.id)
		FROM recipe_ratings
		WHERE recipe_ratings.archived_at IS NULL
			AND recipe_ratings.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
			AND recipe_ratings.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
			AND (
				recipe_ratings.last_updated_at IS NULL
				OR recipe_ratings.last_updated_at > COALESCE(sqlc.narg(updated_before), (SELECT NOW() - '999 years'::INTERVAL))
			)
			AND (
				recipe_ratings.last_updated_at IS NULL
				OR recipe_ratings.last_updated_at < COALESCE(sqlc.narg(updated_after), (SELECT NOW() + '999 years'::INTERVAL))
			)
	) AS filtered_count,
	(
		SELECT COUNT(recipe_ratings.id)
		FROM recipe_ratings
		WHERE recipe_ratings.archived_at IS NULL
	) AS total_count
FROM recipe_ratings
WHERE
	recipe_ratings.archived_at IS NULL
	AND recipe_ratings.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
	AND recipe_ratings.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
	AND (
		recipe_ratings.last_updated_at IS NULL
		OR recipe_ratings.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - '999 years'::INTERVAL))
	)
	AND (
		recipe_ratings.last_updated_at IS NULL
		OR recipe_ratings.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + '999 years'::INTERVAL))
	)
GROUP BY recipe_ratings.id
ORDER BY recipe_ratings.id
LIMIT sqlc.narg(query_limit)
OFFSET sqlc.narg(query_offset);

-- name: GetRecipeRating :one

SELECT
	recipe_ratings.id,
	recipe_ratings.recipe_id,
	recipe_ratings.taste,
	recipe_ratings.difficulty,
	recipe_ratings.cleanup,
	recipe_ratings.instructions,
	recipe_ratings.overall,
	recipe_ratings.notes,
	recipe_ratings.by_user,
	recipe_ratings.created_at,
	recipe_ratings.last_updated_at,
	recipe_ratings.archived_at
FROM recipe_ratings
WHERE recipe_ratings.archived_at IS NULL
	AND recipe_ratings.id = sqlc.arg(id);

-- name: UpdateRecipeRating :execrows

UPDATE recipe_ratings SET
	recipe_id = sqlc.arg(recipe_id),
	taste = sqlc.arg(taste),
	difficulty = sqlc.arg(difficulty),
	cleanup = sqlc.arg(cleanup),
	instructions = sqlc.arg(instructions),
	overall = sqlc.arg(overall),
	notes = sqlc.arg(notes),
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = sqlc.arg(id);
