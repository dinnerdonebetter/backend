// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: meals_exists.sql

package generated

import (
	"context"
)

const MealExists = `-- name: MealExists :exec
SELECT EXISTS ( SELECT meals.id FROM meals WHERE meals.archived_at IS NULL AND meals.id = $1 )
`

func (q *Queries) MealExists(ctx context.Context, db DBTX, id string) error {
	_, err := db.ExecContext(ctx, MealExists, id)
	return err
}