// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: recipe_media_get_one.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const GetRecipeMedia = `-- name: GetRecipeMedia :exec
SELECT
	recipe_media.id,
    recipe_media.belongs_to_recipe,
    recipe_media.belongs_to_recipe_step,
    recipe_media.mime_type,
    recipe_media.internal_path,
    recipe_media.external_path,
    recipe_media.index,
	recipe_media.created_at,
	recipe_media.last_updated_at,
	recipe_media.archived_at
FROM recipe_media
WHERE recipe_media.archived_at IS NULL
	AND recipe_media.id = $1
`

type GetRecipeMediaRow struct {
	CreatedAt           time.Time      `db:"created_at"`
	LastUpdatedAt       sql.NullTime   `db:"last_updated_at"`
	ArchivedAt          sql.NullTime   `db:"archived_at"`
	ID                  string         `db:"id"`
	InternalPath        string         `db:"internal_path"`
	ExternalPath        string         `db:"external_path"`
	MimeType            string         `db:"mime_type"`
	BelongsToRecipe     sql.NullString `db:"belongs_to_recipe"`
	BelongsToRecipeStep sql.NullString `db:"belongs_to_recipe_step"`
	Index               int32          `db:"index"`
}

func (q *Queries) GetRecipeMedia(ctx context.Context, db DBTX, id string) error {
	_, err := db.ExecContext(ctx, GetRecipeMedia, id)
	return err
}