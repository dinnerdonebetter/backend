-- name: CreateAuditLogEntry :exec

INSERT INTO audit_log_entries (
	id,
	resource_type,
	relevant_id,
	event_type,
	changes,
	belongs_to_user,
	belongs_to_household
) VALUES (
	sqlc.arg(id),
	sqlc.arg(resource_type),
	sqlc.arg(relevant_id),
	sqlc.arg(event_type),
	sqlc.arg(changes),
	sqlc.narg(belongs_to_user),
	sqlc.narg(belongs_to_household)
);

-- name: GetAuditLogEntry :one

SELECT
	audit_log_entries.id,
	audit_log_entries.resource_type,
	audit_log_entries.relevant_id,
	audit_log_entries.event_type,
	audit_log_entries.changes,
	audit_log_entries.belongs_to_user,
	audit_log_entries.belongs_to_household,
	audit_log_entries.created_at
FROM audit_log_entries
WHERE audit_log_entries.id = sqlc.arg(id);

-- name: GetAuditLogEntriesForUser :many

SELECT
	audit_log_entries.id,
	audit_log_entries.resource_type,
	audit_log_entries.relevant_id,
	audit_log_entries.event_type,
	audit_log_entries.changes,
	audit_log_entries.belongs_to_user,
	audit_log_entries.belongs_to_household,
	audit_log_entries.created_at,
	(
		SELECT COUNT(audit_log_entries.id)
		FROM audit_log_entries
		WHERE audit_log_entries.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
			AND audit_log_entries.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
			AND audit_log_entries.belongs_to_user = sqlc.arg(belongs_to_user)
	) AS filtered_count,
	(
		SELECT COUNT(audit_log_entries.id)
		FROM audit_log_entries
		WHERE
			audit_log_entries.belongs_to_user = sqlc.arg(belongs_to_user)
	) AS total_count
FROM audit_log_entries
WHERE audit_log_entries.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
	AND audit_log_entries.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
	AND audit_log_entries.belongs_to_user = sqlc.arg(belongs_to_user)
LIMIT sqlc.narg(query_limit)
OFFSET sqlc.narg(query_offset);

-- name: GetAuditLogEntriesForUserAndResourceType :many

SELECT
	audit_log_entries.id,
	audit_log_entries.resource_type,
	audit_log_entries.relevant_id,
	audit_log_entries.event_type,
	audit_log_entries.changes,
	audit_log_entries.belongs_to_user,
	audit_log_entries.belongs_to_household,
	audit_log_entries.created_at,
	(
		SELECT COUNT(audit_log_entries.id)
		FROM audit_log_entries
		WHERE audit_log_entries.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
			AND audit_log_entries.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
			AND audit_log_entries.belongs_to_user = sqlc.arg(belongs_to_user)
			AND audit_log_entries.resource_type = ANY(sqlc.arg(resources)::text[])
	) AS filtered_count,
	(
		SELECT COUNT(audit_log_entries.id)
		FROM audit_log_entries
		WHERE
			audit_log_entries.belongs_to_user = sqlc.arg(belongs_to_user)
			AND audit_log_entries.resource_type = ANY(sqlc.arg(resources)::text[])
	) AS total_count
FROM audit_log_entries
WHERE audit_log_entries.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
	AND audit_log_entries.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
	AND audit_log_entries.belongs_to_user = sqlc.arg(belongs_to_user)
	AND audit_log_entries.resource_type = ANY(sqlc.arg(resources)::text[])
LIMIT sqlc.narg(query_limit)
OFFSET sqlc.narg(query_offset);

-- name: GetAuditLogEntriesForHousehold :many

SELECT
	audit_log_entries.id,
	audit_log_entries.resource_type,
	audit_log_entries.relevant_id,
	audit_log_entries.event_type,
	audit_log_entries.changes,
	audit_log_entries.belongs_to_user,
	audit_log_entries.belongs_to_household,
	audit_log_entries.created_at,
	(
		SELECT COUNT(audit_log_entries.id)
		FROM audit_log_entries
		WHERE audit_log_entries.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
			AND audit_log_entries.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
			AND audit_log_entries.belongs_to_household = sqlc.arg(belongs_to_household)
	) AS filtered_count,
	(
		SELECT COUNT(audit_log_entries.id)
		FROM audit_log_entries
		WHERE
			audit_log_entries.belongs_to_household = sqlc.arg(belongs_to_household)
	) AS total_count
FROM audit_log_entries
WHERE audit_log_entries.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
	AND audit_log_entries.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
	AND audit_log_entries.belongs_to_household = sqlc.arg(belongs_to_household)
LIMIT sqlc.narg(query_limit)
OFFSET sqlc.narg(query_offset);

-- name: GetAuditLogEntriesForHouseholdAndResourceType :many

SELECT
	audit_log_entries.id,
	audit_log_entries.resource_type,
	audit_log_entries.relevant_id,
	audit_log_entries.event_type,
	audit_log_entries.changes,
	audit_log_entries.belongs_to_user,
	audit_log_entries.belongs_to_household,
	audit_log_entries.created_at,
	(
		SELECT COUNT(audit_log_entries.id)
		FROM audit_log_entries
		WHERE audit_log_entries.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
			AND audit_log_entries.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
			AND audit_log_entries.belongs_to_household = sqlc.arg(belongs_to_household)
			AND audit_log_entries.resource_type = ANY(sqlc.arg(resources)::text[])
	) AS filtered_count,
	(
		SELECT COUNT(audit_log_entries.id)
		FROM audit_log_entries
		WHERE
			audit_log_entries.belongs_to_household = sqlc.arg(belongs_to_household)
			AND audit_log_entries.resource_type = ANY(sqlc.arg(resources)::text[])
	) AS total_count
FROM audit_log_entries
WHERE audit_log_entries.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
	AND audit_log_entries.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
	AND audit_log_entries.belongs_to_household = sqlc.arg(belongs_to_household)
	AND audit_log_entries.resource_type = ANY(sqlc.arg(resources)::text[])
LIMIT sqlc.narg(query_limit)
OFFSET sqlc.narg(query_offset);
