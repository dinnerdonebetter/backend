-- name: ArchiveHouseholdInstrumentOwnership :execrows

UPDATE household_instrument_ownerships SET
	archived_at = NOW()
WHERE archived_at IS NULL
	AND id = sqlc.arg(id)
	AND belongs_to_household = sqlc.arg(belongs_to_household);

-- name: CreateHouseholdInstrumentOwnership :exec

INSERT INTO household_instrument_ownerships (
	id,
	notes,
	quantity,
	valid_instrument_id,
	belongs_to_household
) VALUES (
	sqlc.arg(id),
	sqlc.arg(notes),
	sqlc.arg(quantity),
	sqlc.arg(valid_instrument_id),
	sqlc.arg(belongs_to_household)
);

-- name: CheckHouseholdInstrumentOwnershipExistence :one

SELECT EXISTS (
	SELECT household_instrument_ownerships.id
	FROM household_instrument_ownerships
	WHERE household_instrument_ownerships.archived_at IS NULL
		AND household_instrument_ownerships.id = sqlc.arg(id)
		AND household_instrument_ownerships.belongs_to_household = sqlc.arg(belongs_to_household)
);

-- name: GetHouseholdInstrumentOwnerships :many

SELECT
	household_instrument_ownerships.id,
	household_instrument_ownerships.notes,
	household_instrument_ownerships.quantity,
	valid_instruments.id as valid_instrument_id,
	valid_instruments.name as valid_instrument_name,
	valid_instruments.description as valid_instrument_description,
	valid_instruments.icon_path as valid_instrument_icon_path,
	valid_instruments.plural_name as valid_instrument_plural_name,
	valid_instruments.usable_for_storage as valid_instrument_usable_for_storage,
	valid_instruments.slug as valid_instrument_slug,
	valid_instruments.display_in_summary_lists as valid_instrument_display_in_summary_lists,
	valid_instruments.include_in_generated_instructions as valid_instrument_include_in_generated_instructions,
	valid_instruments.last_indexed_at as valid_instrument_last_indexed_at,
	valid_instruments.created_at as valid_instrument_created_at,
	valid_instruments.last_updated_at as valid_instrument_last_updated_at,
	valid_instruments.archived_at as valid_instrument_archived_at,
	household_instrument_ownerships.belongs_to_household,
	household_instrument_ownerships.created_at,
	household_instrument_ownerships.last_updated_at,
	household_instrument_ownerships.archived_at,
	(
		SELECT COUNT(household_instrument_ownerships.id)
		FROM household_instrument_ownerships
		WHERE household_instrument_ownerships.archived_at IS NULL
			AND household_instrument_ownerships.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
			AND household_instrument_ownerships.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
			AND (
				household_instrument_ownerships.last_updated_at IS NULL
				OR household_instrument_ownerships.last_updated_at > COALESCE(sqlc.narg(updated_before), (SELECT NOW() - '999 years'::INTERVAL))
			)
			AND (
				household_instrument_ownerships.last_updated_at IS NULL
				OR household_instrument_ownerships.last_updated_at < COALESCE(sqlc.narg(updated_after), (SELECT NOW() + '999 years'::INTERVAL))
			)
			AND household_instrument_ownerships.belongs_to_household = sqlc.arg(household_id)
	) AS filtered_count,
	(
		SELECT COUNT(household_instrument_ownerships.id)
		FROM household_instrument_ownerships
		WHERE household_instrument_ownerships.archived_at IS NULL
			AND household_instrument_ownerships.belongs_to_household = sqlc.arg(household_id)
	) AS total_count
FROM household_instrument_ownerships
INNER JOIN valid_instruments ON household_instrument_ownerships.valid_instrument_id = valid_instruments.id
WHERE
	household_instrument_ownerships.archived_at IS NULL
	AND household_instrument_ownerships.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
	AND household_instrument_ownerships.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
	AND (
		household_instrument_ownerships.last_updated_at IS NULL
		OR household_instrument_ownerships.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - '999 years'::INTERVAL))
	)
	AND (
		household_instrument_ownerships.last_updated_at IS NULL
		OR household_instrument_ownerships.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + '999 years'::INTERVAL))
	)
	AND household_instrument_ownerships.belongs_to_household = sqlc.arg(household_id)
GROUP BY
	household_instrument_ownerships.id,
	valid_instruments.id
ORDER BY
	household_instrument_ownerships.id
LIMIT sqlc.narg(query_limit)
OFFSET sqlc.narg(query_offset);

-- name: GetHouseholdInstrumentOwnership :one

SELECT
	household_instrument_ownerships.id,
	household_instrument_ownerships.notes,
	household_instrument_ownerships.quantity,
	valid_instruments.id as valid_instrument_id,
	valid_instruments.name as valid_instrument_name,
	valid_instruments.description as valid_instrument_description,
	valid_instruments.icon_path as valid_instrument_icon_path,
	valid_instruments.plural_name as valid_instrument_plural_name,
	valid_instruments.usable_for_storage as valid_instrument_usable_for_storage,
	valid_instruments.slug as valid_instrument_slug,
	valid_instruments.display_in_summary_lists as valid_instrument_display_in_summary_lists,
	valid_instruments.include_in_generated_instructions as valid_instrument_include_in_generated_instructions,
	valid_instruments.last_indexed_at as valid_instrument_last_indexed_at,
	valid_instruments.created_at as valid_instrument_created_at,
	valid_instruments.last_updated_at as valid_instrument_last_updated_at,
	valid_instruments.archived_at as valid_instrument_archived_at,
	household_instrument_ownerships.belongs_to_household,
	household_instrument_ownerships.created_at,
	household_instrument_ownerships.last_updated_at,
	household_instrument_ownerships.archived_at
FROM household_instrument_ownerships
INNER JOIN valid_instruments ON household_instrument_ownerships.valid_instrument_id = valid_instruments.id
WHERE household_instrument_ownerships.archived_at IS NULL
	AND household_instrument_ownerships.id = sqlc.arg(id)
	AND household_instrument_ownerships.belongs_to_household = sqlc.arg(belongs_to_household);

-- name: UpdateHouseholdInstrumentOwnership :execrows

UPDATE household_instrument_ownerships SET
	notes = sqlc.arg(notes),
	quantity = sqlc.arg(quantity),
	valid_instrument_id = sqlc.arg(valid_instrument_id),
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = sqlc.arg(id)
	AND household_instrument_ownerships.belongs_to_household = sqlc.arg(belongs_to_household);
