// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: redeem.sql

package generated

import (
	"context"
)

const RedeemPasswordResetToken = `-- name: RedeemPasswordResetToken :exec

UPDATE password_reset_tokens SET redeemed_at = NOW() WHERE redeemed_at IS NULL AND id = $1
`

func (q *Queries) RedeemPasswordResetToken(ctx context.Context, db DBTX, id string) error {
	_, err := db.ExecContext(ctx, RedeemPasswordResetToken, id)
	return err
}
