-- name: ArchiveAPIClient :exec
UPDATE api_clients SET
	last_updated_at = NOW(),
	archived_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_user = $1 AND id = $2;