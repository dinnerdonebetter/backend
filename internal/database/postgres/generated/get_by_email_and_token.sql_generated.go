// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: get_by_email_and_token.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const GetHouseholdInvitationByEmailAndToken = `-- name: GetHouseholdInvitationByEmailAndToken :one

SELECT
	household_invitations.id,
	households.id,
	households.name,
	households.billing_status,
	households.contact_phone,
	households.address_line_1,
	households.address_line_2,
	households.city,
	households.state,
	households.zip_code,
	households.country,
	households.latitude,
    households.longitude,
	households.payment_processor_customer_id,
	households.subscription_plan_id,
	households.created_at,
	households.last_updated_at,
	households.archived_at,
	households.belongs_to_user,
	household_invitations.to_email,
	household_invitations.to_user,
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
	users.created_at,
	users.last_updated_at,
	users.archived_at,
	household_invitations.to_name,
	household_invitations.status,
	household_invitations.note,
	household_invitations.status_note,
	household_invitations.token,
	household_invitations.expires_at,
	household_invitations.created_at,
	household_invitations.last_updated_at,
	household_invitations.archived_at
FROM household_invitations
	JOIN households ON household_invitations.destination_household = households.id
	JOIN users ON household_invitations.from_user = users.id
WHERE household_invitations.archived_at IS NULL
	AND household_invitations.expires_at > NOW()
	AND household_invitations.to_email = LOWER($1)
	AND household_invitations.token = $2
`

type GetHouseholdInvitationByEmailAndTokenParams struct {
	Lower string `db:"lower"`
	Token string `db:"token"`
}

type GetHouseholdInvitationByEmailAndTokenRow struct {
	ExpiresAt                    time.Time       `db:"expires_at"`
	CreatedAt_2                  time.Time       `db:"created_at_2"`
	CreatedAt                    time.Time       `db:"created_at"`
	CreatedAt_3                  time.Time       `db:"created_at_3"`
	LastUpdatedAt_2              sql.NullTime    `db:"last_updated_at_2"`
	ArchivedAt_2                 sql.NullTime    `db:"archived_at_2"`
	LastUpdatedAt_3              sql.NullTime    `db:"last_updated_at_3"`
	ArchivedAt_3                 sql.NullTime    `db:"archived_at_3"`
	Birthday                     sql.NullTime    `db:"birthday"`
	TwoFactorSecretVerifiedAt    sql.NullTime    `db:"two_factor_secret_verified_at"`
	PasswordLastChangedAt        sql.NullTime    `db:"password_last_changed_at"`
	EmailAddressVerifiedAt       sql.NullTime    `db:"email_address_verified_at"`
	ArchivedAt                   sql.NullTime    `db:"archived_at"`
	LastUpdatedAt                sql.NullTime    `db:"last_updated_at"`
	Username                     string          `db:"username"`
	TwoFactorSecret              string          `db:"two_factor_secret"`
	PaymentProcessorCustomerID   string          `db:"payment_processor_customer_id"`
	ID_2                         string          `db:"id_2"`
	BelongsToUser                string          `db:"belongs_to_user"`
	ToEmail                      string          `db:"to_email"`
	Name                         string          `db:"name"`
	ID_3                         string          `db:"id_3"`
	FirstName                    string          `db:"first_name"`
	LastName                     string          `db:"last_name"`
	ID                           string          `db:"id"`
	EmailAddress                 string          `db:"email_address"`
	BillingStatus                string          `db:"billing_status"`
	ContactPhone                 string          `db:"contact_phone"`
	HashedPassword               string          `db:"hashed_password"`
	Token                        string          `db:"token"`
	Country                      string          `db:"country"`
	StatusNote                   string          `db:"status_note"`
	ZipCode                      string          `db:"zip_code"`
	ServiceRole                  string          `db:"service_role"`
	UserAccountStatus            string          `db:"user_account_status"`
	UserAccountStatusExplanation string          `db:"user_account_status_explanation"`
	State                        string          `db:"state"`
	City                         string          `db:"city"`
	AddressLine2                 string          `db:"address_line_2"`
	AddressLine1                 string          `db:"address_line_1"`
	ToName                       string          `db:"to_name"`
	Status                       InvitationState `db:"status"`
	Note                         string          `db:"note"`
	SubscriptionPlanID           sql.NullString  `db:"subscription_plan_id"`
	AvatarSrc                    sql.NullString  `db:"avatar_src"`
	Latitude                     sql.NullString  `db:"latitude"`
	ToUser                       sql.NullString  `db:"to_user"`
	Longitude                    sql.NullString  `db:"longitude"`
	RequiresPasswordChange       bool            `db:"requires_password_change"`
}

func (q *Queries) GetHouseholdInvitationByEmailAndToken(ctx context.Context, db DBTX, arg *GetHouseholdInvitationByEmailAndTokenParams) (*GetHouseholdInvitationByEmailAndTokenRow, error) {
	row := db.QueryRowContext(ctx, GetHouseholdInvitationByEmailAndToken, arg.Lower, arg.Token)
	var i GetHouseholdInvitationByEmailAndTokenRow
	err := row.Scan(
		&i.ID,
		&i.ID_2,
		&i.Name,
		&i.BillingStatus,
		&i.ContactPhone,
		&i.AddressLine1,
		&i.AddressLine2,
		&i.City,
		&i.State,
		&i.ZipCode,
		&i.Country,
		&i.Latitude,
		&i.Longitude,
		&i.PaymentProcessorCustomerID,
		&i.SubscriptionPlanID,
		&i.CreatedAt,
		&i.LastUpdatedAt,
		&i.ArchivedAt,
		&i.BelongsToUser,
		&i.ToEmail,
		&i.ToUser,
		&i.ID_3,
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
		&i.CreatedAt_2,
		&i.LastUpdatedAt_2,
		&i.ArchivedAt_2,
		&i.ToName,
		&i.Status,
		&i.Note,
		&i.StatusNote,
		&i.Token,
		&i.ExpiresAt,
		&i.CreatedAt_3,
		&i.LastUpdatedAt_3,
		&i.ArchivedAt_3,
	)
	return &i, err
}
