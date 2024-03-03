-- name: AddToHouseholdDuringCreation :exec

INSERT INTO household_user_memberships (
	id,
	belongs_to_household,
	belongs_to_user,
	household_role
) VALUES (
	sqlc.arg(id),
	sqlc.arg(belongs_to_household),
	sqlc.arg(belongs_to_user),
	sqlc.arg(household_role)
);

-- name: ArchiveHousehold :execrows

UPDATE households SET
	last_updated_at = NOW(),
	archived_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_user = sqlc.arg(belongs_to_user)
	AND id = sqlc.arg(id);

-- name: CreateHousehold :exec

INSERT INTO households (
	id,
	name,
	billing_status,
	contact_phone,
	belongs_to_user,
	address_line_1,
	address_line_2,
	city,
	state,
	zip_code,
	country,
	latitude,
	longitude,
	webhook_hmac_secret
) VALUES (
	sqlc.arg(id),
	sqlc.arg(name),
	sqlc.arg(billing_status),
	sqlc.arg(contact_phone),
	sqlc.arg(belongs_to_user),
	sqlc.arg(address_line_1),
	sqlc.arg(address_line_2),
	sqlc.arg(city),
	sqlc.arg(state),
	sqlc.arg(zip_code),
	sqlc.arg(country),
	sqlc.arg(latitude),
	sqlc.arg(longitude),
	sqlc.arg(webhook_hmac_secret)
);

-- name: GetHouseholdByIDWithMemberships :many

SELECT
	households.id,
	households.name,
	households.billing_status,
	households.contact_phone,
	households.payment_processor_customer_id,
	households.subscription_plan_id,
	households.belongs_to_user,
	households.time_zone,
	households.address_line_1,
	households.address_line_2,
	households.city,
	households.state,
	households.zip_code,
	households.country,
	households.latitude,
	households.longitude,
	households.last_payment_provider_sync_occurred_at,
	households.webhook_hmac_secret,
	households.created_at,
	households.last_updated_at,
	households.archived_at,
	users.id as user_id,
	users.username as user_username,
	users.avatar_src as user_avatar_src,
	users.email_address as user_email_address,
	users.hashed_password as user_hashed_password,
	users.password_last_changed_at as user_password_last_changed_at,
	users.requires_password_change as user_requires_password_change,
	users.two_factor_secret as user_two_factor_secret,
	users.two_factor_secret_verified_at as user_two_factor_secret_verified_at,
	users.service_role as user_service_role,
	users.user_account_status as user_user_account_status,
	users.user_account_status_explanation as user_user_account_status_explanation,
	users.birthday as user_birthday,
	users.email_address_verification_token as user_email_address_verification_token,
	users.email_address_verified_at as user_email_address_verified_at,
	users.first_name as user_first_name,
	users.last_name as user_last_name,
	users.last_accepted_terms_of_service as user_last_accepted_terms_of_service,
	users.last_accepted_privacy_policy as user_last_accepted_privacy_policy,
	users.last_indexed_at as user_last_indexed_at,
	users.created_at as user_created_at,
	users.last_updated_at as user_last_updated_at,
	users.archived_at as user_archived_at,
	household_user_memberships.id as membership_id,
	household_user_memberships.belongs_to_household as membership_belongs_to_household,
	household_user_memberships.belongs_to_user as membership_belongs_to_user,
	household_user_memberships.default_household as membership_default_household,
	household_user_memberships.household_role as membership_household_role,
	household_user_memberships.created_at as membership_created_at,
	household_user_memberships.last_updated_at as membership_last_updated_at,
	household_user_memberships.archived_at as membership_archived_at
FROM households
	JOIN household_user_memberships ON household_user_memberships.belongs_to_household = households.id
	JOIN users ON household_user_memberships.belongs_to_user = users.id
WHERE households.archived_at IS NULL
	AND household_user_memberships.archived_at IS NULL
	AND households.id = sqlc.arg(id);

-- name: GetHouseholdsForUser :many

SELECT
	households.id,
	households.name,
	households.billing_status,
	households.contact_phone,
	households.payment_processor_customer_id,
	households.subscription_plan_id,
	households.belongs_to_user,
	households.time_zone,
	households.address_line_1,
	households.address_line_2,
	households.city,
	households.state,
	households.zip_code,
	households.country,
	households.latitude,
	households.longitude,
	households.last_payment_provider_sync_occurred_at,
	households.webhook_hmac_secret,
	households.created_at,
	households.last_updated_at,
	households.archived_at,
	(
		SELECT COUNT(households.id)
		FROM households
			JOIN household_user_memberships ON household_user_memberships.belongs_to_household = households.id
		WHERE households.archived_at IS NULL
			AND household_user_memberships.belongs_to_user = sqlc.arg(belongs_to_user)
			AND households.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
			AND households.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
			AND (
				households.last_updated_at IS NULL
				OR households.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - '999 years'::INTERVAL))
			)
			AND (
				households.last_updated_at IS NULL
				OR households.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + '999 years'::INTERVAL))
			)
	) as filtered_count,
	(
		SELECT COUNT(households.id)
		FROM households
		WHERE households.archived_at IS NULL
	) AS total_count
FROM households
	JOIN household_user_memberships ON household_user_memberships.belongs_to_household = households.id
	JOIN users ON household_user_memberships.belongs_to_user = users.id
WHERE households.archived_at IS NULL
	AND household_user_memberships.archived_at IS NULL
	AND household_user_memberships.belongs_to_user = sqlc.arg(belongs_to_user)
	AND households.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
	AND households.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
	AND (
		households.last_updated_at IS NULL
		OR households.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - '999 years'::INTERVAL))
	)
	AND (
		households.last_updated_at IS NULL
		OR households.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + '999 years'::INTERVAL))
	)
LIMIT sqlc.narg(query_limit)
OFFSET sqlc.narg(query_offset);

-- name: UpdateHousehold :execrows

UPDATE households SET
	name = sqlc.arg(name),
	contact_phone = sqlc.arg(contact_phone),
	address_line_1 = sqlc.arg(address_line_1),
	address_line_2 = sqlc.arg(address_line_2),
	city = sqlc.arg(city),
	state = sqlc.arg(state),
	zip_code = sqlc.arg(zip_code),
	country = sqlc.arg(country),
	latitude = sqlc.arg(latitude),
	longitude = sqlc.arg(longitude),
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_user = sqlc.arg(belongs_to_user)
	AND id = sqlc.arg(id);

-- name: UpdateHouseholdWebhookEncryptionKey :execrows

UPDATE households SET
	webhook_hmac_secret = sqlc.arg(webhook_hmac_secret),
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_user = sqlc.arg(belongs_to_user)
	AND id = sqlc.arg(id);
