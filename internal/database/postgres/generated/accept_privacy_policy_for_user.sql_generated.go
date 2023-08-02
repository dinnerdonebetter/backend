// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: accept_privacy_policy_for_user.sql

package generated

import (
	"context"
)

const AcceptPrivacyPolicyForUser = `-- name: AcceptPrivacyPolicyForUser :exec

UPDATE users SET
	last_accepted_privacy_policy = NOW()
WHERE archived_at IS NULL
	AND id = $1
`

func (q *Queries) AcceptPrivacyPolicyForUser(ctx context.Context, db DBTX, id string) error {
	_, err := db.ExecContext(ctx, AcceptPrivacyPolicyForUser, id)
	return err
}
