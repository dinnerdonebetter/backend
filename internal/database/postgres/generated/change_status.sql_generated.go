// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: change_status.sql

package generated

import (
	"context"
	"database/sql"
)

const ChangeMealPlanTaskStatus = `-- name: ChangeMealPlanTaskStatus :exec

UPDATE meal_plan_tasks SET completed_at = $4, status_explanation = $3, status = $2 WHERE id = $1
`

type ChangeMealPlanTaskStatusParams struct {
	ID                string         `db:"id"`
	Status            PrepStepStatus `db:"status"`
	StatusExplanation string         `db:"status_explanation"`
	CompletedAt       sql.NullTime   `db:"completed_at"`
}

func (q *Queries) ChangeMealPlanTaskStatus(ctx context.Context, db DBTX, arg *ChangeMealPlanTaskStatusParams) error {
	_, err := db.ExecContext(ctx, ChangeMealPlanTaskStatus,
		arg.ID,
		arg.Status,
		arg.StatusExplanation,
		arg.CompletedAt,
	)
	return err
}
