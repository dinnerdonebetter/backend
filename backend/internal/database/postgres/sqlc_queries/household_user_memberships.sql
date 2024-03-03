-- name: AddUserToHousehold :exec

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

-- name: CreateHouseholdUserMembershipForNewUser :exec

INSERT INTO household_user_memberships (
	id,
	belongs_to_household,
	belongs_to_user,
	default_household,
	household_role
) VALUES (
	sqlc.arg(id),
	sqlc.arg(belongs_to_household),
	sqlc.arg(belongs_to_user),
	sqlc.arg(default_household),
	sqlc.arg(household_role)
);

-- name: GetDefaultHouseholdIDForUser :one

SELECT households.id
FROM households
	JOIN household_user_memberships ON household_user_memberships.belongs_to_household = households.id
WHERE household_user_memberships.belongs_to_user = sqlc.arg(belongs_to_user)
	AND household_user_memberships.default_household = TRUE;

-- name: GetHouseholdUserMembershipsForUser :many

SELECT
	household_user_memberships.id,
	household_user_memberships.belongs_to_household,
	household_user_memberships.belongs_to_user,
	household_user_memberships.default_household,
	household_user_memberships.household_role,
	household_user_memberships.created_at,
	household_user_memberships.last_updated_at,
	household_user_memberships.archived_at
FROM household_user_memberships
	JOIN households ON households.id = household_user_memberships.belongs_to_household
WHERE household_user_memberships.archived_at IS NULL
	AND household_user_memberships.belongs_to_user = sqlc.arg(belongs_to_user);

-- name: MarkHouseholdUserMembershipAsUserDefault :exec

UPDATE household_user_memberships SET
	default_household = (belongs_to_user = sqlc.arg(belongs_to_user) AND belongs_to_household = sqlc.arg(belongs_to_household))
WHERE archived_at IS NULL
	AND belongs_to_user = sqlc.arg(belongs_to_user);

-- name: ModifyHouseholdUserPermissions :exec

UPDATE household_user_memberships SET
	household_role = sqlc.arg(household_role)
WHERE belongs_to_household = sqlc.arg(belongs_to_household)
	AND belongs_to_user = sqlc.arg(belongs_to_user);

-- name: RemoveUserFromHousehold :exec

UPDATE household_user_memberships SET
	archived_at = NOW(),
	default_household = 'false'
WHERE household_user_memberships.archived_at IS NULL
	AND household_user_memberships.belongs_to_household = sqlc.arg(belongs_to_household)
	AND household_user_memberships.belongs_to_user = sqlc.arg(belongs_to_user);

-- name: TransferHouseholdMembership :exec

UPDATE household_user_memberships SET
	belongs_to_user = sqlc.arg(belongs_to_user)
WHERE archived_at IS NULL
	AND belongs_to_household = sqlc.arg(belongs_to_household)
	AND belongs_to_user = sqlc.arg(belongs_to_user);

-- name: TransferHouseholdOwnership :exec

UPDATE households SET
	belongs_to_user = sqlc.arg(new_owner)
WHERE archived_at IS NULL
	AND belongs_to_user = sqlc.arg(old_owner)
	AND id = sqlc.arg(household_id);

-- name: UserIsHouseholdMember :one

SELECT EXISTS (
	SELECT household_user_memberships.id
	FROM household_user_memberships
	WHERE household_user_memberships.archived_at IS NULL
		AND household_user_memberships.belongs_to_household = sqlc.arg(belongs_to_household)
		AND household_user_memberships.belongs_to_user = sqlc.arg(belongs_to_user)
);
