-- name: AttachHouseholdInvitationsToUserID :exec

UPDATE household_invitations SET
    to_user = sqlc.arg(user_id),
    last_updated_at = NOW()
WHERE archived_at IS NULL
  AND to_email = LOWER(sqlc.arg(email_address));

-- name: CreateHouseholdInvitation :exec

INSERT INTO household_invitations (id,from_user,to_user,to_name,note,to_email,token,destination_household,expires_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9);

-- name: CheckHouseholdInvitationExistence :one

SELECT EXISTS ( SELECT household_invitations.id FROM household_invitations WHERE household_invitations.archived_at IS NULL AND household_invitations.id = $1 );

-- name: GetHouseholdInvitationByEmailAndToken :one

SELECT
	household_invitations.id,
	households.id as household_id,
	households.name as household_name,
	households.billing_status as household_billing_status,
	households.contact_phone as household_contact_phone,
	households.address_line_1 as household_address_line_1,
	households.address_line_2 as household_address_line_2,
	households.city as household_city,
	households.state as household_state,
	households.zip_code as household_zip_code,
	households.country as household_country,
	households.latitude as household_latitude,
    households.longitude as household_longitude,
	households.payment_processor_customer_id as household_payment_processor_customer_id,
	households.subscription_plan_id as household_subscription_plan_id,
	households.created_at as household_created_at,
	households.last_updated_at as household_last_updated_at,
	households.archived_at as household_archived_at,
	households.belongs_to_user as household_belongs_to_user,
	household_invitations.to_email,
	household_invitations.to_user,
    users.id as user_id,
    users.first_name as user_first_name,
    users.last_name as user_last_name,
    users.username as user_username,
    users.email_address as user_email_address,
    users.email_address_verified_at as user_email_address_verified_at,
    users.avatar_src as user_avatar_src,
    users.hashed_password as user_hashed_password,
    users.requires_password_change as user_requires_password_change,
    users.password_last_changed_at as user_password_last_changed_at,
    users.two_factor_secret as user_two_factor_secret,
    users.two_factor_secret_verified_at as user_two_factor_secret_verified_at,
    users.service_role as user_service_role,
    users.user_account_status as user_user_account_status,
    users.user_account_status_explanation as user_user_account_status_explanation,
    users.birthday as user_birthday,
    users.last_accepted_terms_of_service as user_last_accepted_terms_of_service,
    users.last_accepted_privacy_policy as user_last_accepted_privacy_policy,
    users.created_at as user_created_at,
    users.last_updated_at as user_last_updated_at,
    users.archived_at as user_archived_at,
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
	AND household_invitations.to_email = LOWER(sqlc.arg(email_address))
	AND household_invitations.token = sqlc.arg(token);

-- name: GetHouseholdInvitationByHouseholdAndID :one

SELECT
	household_invitations.id,
    households.id as household_id,
    households.name as household_name,
    households.billing_status as household_billing_status,
    households.contact_phone as household_contact_phone,
    households.address_line_1 as household_address_line_1,
    households.address_line_2 as household_address_line_2,
    households.city as household_city,
    households.state as household_state,
    households.zip_code as household_zip_code,
    households.country as household_country,
    households.latitude as household_latitude,
    households.longitude as household_longitude,
    households.payment_processor_customer_id as household_payment_processor_customer_id,
    households.subscription_plan_id as household_subscription_plan_id,
    households.created_at as household_created_at,
    households.last_updated_at as household_last_updated_at,
    households.archived_at as household_archived_at,
    households.belongs_to_user as household_belongs_to_user,
	household_invitations.to_email,
	household_invitations.to_user,
    users.id as user_id,
    users.first_name as user_first_name,
    users.last_name as user_last_name,
    users.username as user_username,
    users.email_address as user_email_address,
    users.email_address_verified_at as user_email_address_verified_at,
    users.avatar_src as user_avatar_src,
    users.hashed_password as user_hashed_password,
    users.requires_password_change as user_requires_password_change,
    users.password_last_changed_at as user_password_last_changed_at,
    users.two_factor_secret as user_two_factor_secret,
    users.two_factor_secret_verified_at as user_two_factor_secret_verified_at,
    users.service_role as user_service_role,
    users.user_account_status as user_user_account_status,
    users.user_account_status_explanation as user_user_account_status_explanation,
    users.birthday as user_birthday,
    users.last_accepted_terms_of_service as user_last_accepted_terms_of_service,
    users.last_accepted_privacy_policy as user_last_accepted_privacy_policy,
    users.created_at as user_created_at,
    users.last_updated_at as user_last_updated_at,
    users.archived_at as user_archived_at,
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
      AND household_invitations.expires_at > NOW()
	AND household_invitations.destination_household = $1
	AND household_invitations.id = $2;

-- name: GetHouseholdInvitationByTokenAndID :one

SELECT
	household_invitations.id,
    households.id as household_id,
    households.name as household_name,
    households.billing_status as household_billing_status,
    households.contact_phone as household_contact_phone,
    households.address_line_1 as household_address_line_1,
    households.address_line_2 as household_address_line_2,
    households.city as household_city,
    households.state as household_state,
    households.zip_code as household_zip_code,
    households.country as household_country,
    households.latitude as household_latitude,
    households.longitude as household_longitude,
    households.payment_processor_customer_id as household_payment_processor_customer_id,
    households.subscription_plan_id as household_subscription_plan_id,
    households.created_at as household_created_at,
    households.last_updated_at as household_last_updated_at,
    households.archived_at as household_archived_at,
    households.belongs_to_user as household_belongs_to_user,
    household_invitations.to_email,
    household_invitations.to_user,
    users.id as user_id,
    users.first_name as user_first_name,
    users.last_name as user_last_name,
    users.username as user_username,
    users.email_address as user_email_address,
    users.email_address_verified_at as user_email_address_verified_at,
    users.avatar_src as user_avatar_src,
    users.hashed_password as user_hashed_password,
    users.requires_password_change as user_requires_password_change,
    users.password_last_changed_at as user_password_last_changed_at,
    users.two_factor_secret as user_two_factor_secret,
    users.two_factor_secret_verified_at as user_two_factor_secret_verified_at,
    users.service_role as user_service_role,
    users.user_account_status as user_user_account_status,
    users.user_account_status_explanation as user_user_account_status_explanation,
    users.birthday as user_birthday,
    users.last_accepted_terms_of_service as user_last_accepted_terms_of_service,
    users.last_accepted_privacy_policy as user_last_accepted_privacy_policy,
    users.created_at as user_created_at,
    users.last_updated_at as user_last_updated_at,
    users.archived_at as user_archived_at,
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
	AND household_invitations.token = $1
	AND household_invitations.id = $2;

-- name: GetPendingInvitesFromUser :many

SELECT
    household_invitations.id,
    households.id as household_id,
    households.name as household_name,
    households.billing_status as household_billing_status,
    households.contact_phone as household_contact_phone,
    households.address_line_1 as household_address_line_1,
    households.address_line_2 as household_address_line_2,
    households.city as household_city,
    households.state as household_state,
    households.zip_code as household_zip_code,
    households.country as household_country,
    households.latitude as household_latitude,
    households.longitude as household_longitude,
    households.payment_processor_customer_id as household_payment_processor_customer_id,
    households.subscription_plan_id as household_subscription_plan_id,
    households.created_at as household_created_at,
    households.last_updated_at as household_last_updated_at,
    households.archived_at as household_archived_at,
    households.belongs_to_user as household_belongs_to_user,
    household_invitations.to_email,
    household_invitations.to_user,
    users.id as user_id,
    users.first_name as user_first_name,
    users.last_name as user_last_name,
    users.username as user_username,
    users.email_address as user_email_address,
    users.email_address_verified_at as user_email_address_verified_at,
    users.avatar_src as user_avatar_src,
    users.hashed_password as user_hashed_password,
    users.requires_password_change as user_requires_password_change,
    users.password_last_changed_at as user_password_last_changed_at,
    users.two_factor_secret as user_two_factor_secret,
    users.two_factor_secret_verified_at as user_two_factor_secret_verified_at,
    users.service_role as user_service_role,
    users.user_account_status as user_user_account_status,
    users.user_account_status_explanation as user_user_account_status_explanation,
    users.birthday as user_birthday,
    users.last_accepted_terms_of_service as user_last_accepted_terms_of_service,
    users.last_accepted_privacy_policy as user_last_accepted_privacy_policy,
    users.created_at as user_created_at,
    users.last_updated_at as user_last_updated_at,
    users.archived_at as user_archived_at,
    household_invitations.to_name,
    household_invitations.status,
    household_invitations.note,
    household_invitations.status_note,
    household_invitations.token,
    household_invitations.expires_at,
    household_invitations.created_at,
    household_invitations.last_updated_at,
    household_invitations.archived_at,
    (
        SELECT COUNT(household_invitations.id)
        FROM household_invitations
        WHERE household_invitations.archived_at IS NULL
          AND household_invitations.from_user = sqlc.arg(user_id)
          AND household_invitations.status = sqlc.arg(status)
          AND household_invitations.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
          AND household_invitations.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
          AND (household_invitations.last_updated_at IS NULL OR household_invitations.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
          AND (household_invitations.last_updated_at IS NULL OR household_invitations.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
    ) as filtered_count,
    (
        SELECT COUNT(household_invitations.id)
        FROM household_invitations
        WHERE household_invitations.archived_at IS NULL
          AND household_invitations.from_user = sqlc.arg(user_id)
          AND household_invitations.status = sqlc.arg(status)
    ) as total_count
FROM household_invitations
    JOIN households ON household_invitations.destination_household = households.id
    JOIN users ON household_invitations.from_user = users.id
WHERE household_invitations.archived_at IS NULL
    AND household_invitations.from_user = sqlc.arg(user_id)
    AND household_invitations.status = sqlc.arg(status)
    AND household_invitations.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
    AND household_invitations.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
    AND (household_invitations.last_updated_at IS NULL OR household_invitations.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
    AND (household_invitations.last_updated_at IS NULL OR household_invitations.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
    OFFSET sqlc.narg(query_offset)
    LIMIT sqlc.narg(query_limit);

-- name: GetPendingInvitesForUser :many

SELECT
    household_invitations.id,
    households.id as household_id,
    households.name as household_name,
    households.billing_status as household_billing_status,
    households.contact_phone as household_contact_phone,
    households.address_line_1 as household_address_line_1,
    households.address_line_2 as household_address_line_2,
    households.city as household_city,
    households.state as household_state,
    households.zip_code as household_zip_code,
    households.country as household_country,
    households.latitude as household_latitude,
    households.longitude as household_longitude,
    households.payment_processor_customer_id as household_payment_processor_customer_id,
    households.subscription_plan_id as household_subscription_plan_id,
    households.created_at as household_created_at,
    households.last_updated_at as household_last_updated_at,
    households.archived_at as household_archived_at,
    households.belongs_to_user as household_belongs_to_user,
    household_invitations.to_email,
    household_invitations.to_user,
    users.id as user_id,
    users.first_name as user_first_name,
    users.last_name as user_last_name,
    users.username as user_username,
    users.email_address as user_email_address,
    users.email_address_verified_at as user_email_address_verified_at,
    users.avatar_src as user_avatar_src,
    users.hashed_password as user_hashed_password,
    users.requires_password_change as user_requires_password_change,
    users.password_last_changed_at as user_password_last_changed_at,
    users.two_factor_secret as user_two_factor_secret,
    users.two_factor_secret_verified_at as user_two_factor_secret_verified_at,
    users.service_role as user_service_role,
    users.user_account_status as user_user_account_status,
    users.user_account_status_explanation as user_user_account_status_explanation,
    users.birthday as user_birthday,
    users.last_accepted_terms_of_service as user_last_accepted_terms_of_service,
    users.last_accepted_privacy_policy as user_last_accepted_privacy_policy,
    users.created_at as user_created_at,
    users.last_updated_at as user_last_updated_at,
    users.archived_at as user_archived_at,
    household_invitations.to_name,
    household_invitations.status,
    household_invitations.note,
    household_invitations.status_note,
    household_invitations.token,
    household_invitations.expires_at,
    household_invitations.created_at,
    household_invitations.last_updated_at,
    household_invitations.archived_at,
    (
        SELECT COUNT(household_invitations.id)
        FROM household_invitations
        WHERE household_invitations.archived_at IS NULL
          AND household_invitations.to_user = sqlc.arg(user_id)
          AND household_invitations.status = sqlc.arg(status)
          AND household_invitations.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
          AND household_invitations.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
          AND (household_invitations.last_updated_at IS NULL OR household_invitations.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
          AND (household_invitations.last_updated_at IS NULL OR household_invitations.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
    ) as filtered_count,
    (
        SELECT COUNT(household_invitations.id)
        FROM household_invitations
        WHERE household_invitations.archived_at IS NULL
          AND household_invitations.to_user = sqlc.arg(user_id)
          AND household_invitations.status = sqlc.arg(status)
    ) as total_count
FROM household_invitations
    JOIN households ON household_invitations.destination_household = households.id
    JOIN users ON household_invitations.from_user = users.id
WHERE household_invitations.archived_at IS NULL
    AND household_invitations.to_user = sqlc.arg(user_id)
    AND household_invitations.status = sqlc.arg(status)
    AND household_invitations.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
    AND household_invitations.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
    AND (household_invitations.last_updated_at IS NULL OR household_invitations.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
    AND (household_invitations.last_updated_at IS NULL OR household_invitations.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
    OFFSET sqlc.narg(query_offset)
    LIMIT sqlc.narg(query_limit);

-- name: SetHouseholdInvitationStatus :exec

UPDATE household_invitations SET
	status = $1,
	status_note = $2,
	last_updated_at = NOW(),
	archived_at = NOW()
WHERE archived_at IS NULL
	AND id = $3;