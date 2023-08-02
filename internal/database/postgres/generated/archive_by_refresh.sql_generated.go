// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: archive_by_refresh.sql

package generated

import (
	"context"
)

const ArchiveOAuth2ClientTokenByRefresh = `-- name: ArchiveOAuth2ClientTokenByRefresh :exec

DELETE FROM oauth2_client_tokens WHERE refresh = $1
`

func (q *Queries) ArchiveOAuth2ClientTokenByRefresh(ctx context.Context, db DBTX, refresh string) error {
	_, err := db.ExecContext(ctx, ArchiveOAuth2ClientTokenByRefresh, refresh)
	return err
}