// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: get_by_username.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const getUserByUsername = `-- name: GetUserByUsername :one

SELECT
	users.id,
	users.first_name,
	users.last_name,
	users.username,
	users.email_address,
	users.email_address_verified_at,
	users.avatar_src,
	users.hashed_password,
	users.requires_password_change,
	users.password_last_changed_at,
	users.two_factor_secret,
	users.two_factor_secret_verified_at,
	users.service_role,
	users.user_account_status,
	users.user_account_status_explanation,
	users.birthday,
	users.last_accepted_terms_of_service,
    users.last_accepted_privacy_policy,
	users.created_at,
	users.last_updated_at,
	users.archived_at
FROM users
WHERE users.archived_at IS NULL
	AND users.username = $1
`

type GetUserByUsernameRow struct {
	ID                           string
	FirstName                    string
	LastName                     string
	Username                     string
	EmailAddress                 string
	EmailAddressVerifiedAt       sql.NullTime
	AvatarSrc                    sql.NullString
	HashedPassword               string
	RequiresPasswordChange       bool
	PasswordLastChangedAt        sql.NullTime
	TwoFactorSecret              string
	TwoFactorSecretVerifiedAt    sql.NullTime
	ServiceRole                  string
	UserAccountStatus            string
	UserAccountStatusExplanation string
	Birthday                     sql.NullTime
	LastAcceptedTermsOfService   sql.NullTime
	LastAcceptedPrivacyPolicy    sql.NullTime
	CreatedAt                    time.Time
	LastUpdatedAt                sql.NullTime
	ArchivedAt                   sql.NullTime
}

func (q *Queries) GetUserByUsername(ctx context.Context, db DBTX, username string) (*GetUserByUsernameRow, error) {
	row := db.QueryRowContext(ctx, getUserByUsername, username)
	var i GetUserByUsernameRow
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Username,
		&i.EmailAddress,
		&i.EmailAddressVerifiedAt,
		&i.AvatarSrc,
		&i.HashedPassword,
		&i.RequiresPasswordChange,
		&i.PasswordLastChangedAt,
		&i.TwoFactorSecret,
		&i.TwoFactorSecretVerifiedAt,
		&i.ServiceRole,
		&i.UserAccountStatus,
		&i.UserAccountStatusExplanation,
		&i.Birthday,
		&i.LastAcceptedTermsOfService,
		&i.LastAcceptedPrivacyPolicy,
		&i.CreatedAt,
		&i.LastUpdatedAt,
		&i.ArchivedAt,
	)
	return &i, err
}
