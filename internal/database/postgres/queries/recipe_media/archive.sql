-- name: ArchiveRecipeMedia :exec

UPDATE recipe_media SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1;
