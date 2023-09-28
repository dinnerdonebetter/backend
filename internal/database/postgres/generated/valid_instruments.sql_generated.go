// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: valid_instruments.sql

package generated

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

const archiveValidInstrument = `-- name: ArchiveValidInstrument :execrows

UPDATE valid_instruments SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1
`

func (q *Queries) ArchiveValidInstrument(ctx context.Context, db DBTX, id string) (int64, error) {
	result, err := db.ExecContext(ctx, archiveValidInstrument, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const checkValidInstrumentExistence = `-- name: CheckValidInstrumentExistence :one

SELECT EXISTS (
	SELECT valid_instruments.id
	FROM valid_instruments
	WHERE valid_instruments.archived_at IS NULL
		AND valid_instruments.id = $1
)
`

func (q *Queries) CheckValidInstrumentExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkValidInstrumentExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createValidInstrument = `-- name: CreateValidInstrument :exec

INSERT INTO valid_instruments (
	id,
	name,
	description,
	icon_path,
	plural_name,
	usable_for_storage,
	slug,
	display_in_summary_lists,
	include_in_generated_instructions
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7,
	$8,
	$9
)
`

type CreateValidInstrumentParams struct {
	ID                             string
	Name                           string
	Description                    string
	IconPath                       string
	PluralName                     string
	Slug                           string
	UsableForStorage               bool
	DisplayInSummaryLists          bool
	IncludeInGeneratedInstructions bool
}

func (q *Queries) CreateValidInstrument(ctx context.Context, db DBTX, arg *CreateValidInstrumentParams) error {
	_, err := db.ExecContext(ctx, createValidInstrument,
		arg.ID,
		arg.Name,
		arg.Description,
		arg.IconPath,
		arg.PluralName,
		arg.UsableForStorage,
		arg.Slug,
		arg.DisplayInSummaryLists,
		arg.IncludeInGeneratedInstructions,
	)
	return err
}

const getRandomValidInstrument = `-- name: GetRandomValidInstrument :one

SELECT
	valid_instruments.id,
	valid_instruments.name,
	valid_instruments.description,
	valid_instruments.icon_path,
	valid_instruments.plural_name,
	valid_instruments.usable_for_storage,
	valid_instruments.slug,
	valid_instruments.display_in_summary_lists,
	valid_instruments.include_in_generated_instructions,
	valid_instruments.last_indexed_at,
	valid_instruments.created_at,
	valid_instruments.last_updated_at,
	valid_instruments.archived_at
FROM valid_instruments
WHERE valid_instruments.archived_at IS NULL
ORDER BY RANDOM() LIMIT 1
`

type GetRandomValidInstrumentRow struct {
	CreatedAt                      time.Time
	LastIndexedAt                  sql.NullTime
	ArchivedAt                     sql.NullTime
	LastUpdatedAt                  sql.NullTime
	IconPath                       string
	Slug                           string
	PluralName                     string
	ID                             string
	Description                    string
	Name                           string
	UsableForStorage               bool
	DisplayInSummaryLists          bool
	IncludeInGeneratedInstructions bool
}

func (q *Queries) GetRandomValidInstrument(ctx context.Context, db DBTX) (*GetRandomValidInstrumentRow, error) {
	row := db.QueryRowContext(ctx, getRandomValidInstrument)
	var i GetRandomValidInstrumentRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.IconPath,
		&i.PluralName,
		&i.UsableForStorage,
		&i.Slug,
		&i.DisplayInSummaryLists,
		&i.IncludeInGeneratedInstructions,
		&i.LastIndexedAt,
		&i.CreatedAt,
		&i.LastUpdatedAt,
		&i.ArchivedAt,
	)
	return &i, err
}

const getValidInstrument = `-- name: GetValidInstrument :one

SELECT
	valid_instruments.id,
	valid_instruments.name,
	valid_instruments.description,
	valid_instruments.icon_path,
	valid_instruments.plural_name,
	valid_instruments.usable_for_storage,
	valid_instruments.slug,
	valid_instruments.display_in_summary_lists,
	valid_instruments.include_in_generated_instructions,
	valid_instruments.last_indexed_at,
	valid_instruments.created_at,
	valid_instruments.last_updated_at,
	valid_instruments.archived_at
FROM valid_instruments
WHERE valid_instruments.archived_at IS NULL
AND valid_instruments.id = $1
`

type GetValidInstrumentRow struct {
	CreatedAt                      time.Time
	LastIndexedAt                  sql.NullTime
	ArchivedAt                     sql.NullTime
	LastUpdatedAt                  sql.NullTime
	IconPath                       string
	Slug                           string
	PluralName                     string
	ID                             string
	Description                    string
	Name                           string
	UsableForStorage               bool
	DisplayInSummaryLists          bool
	IncludeInGeneratedInstructions bool
}

func (q *Queries) GetValidInstrument(ctx context.Context, db DBTX, id string) (*GetValidInstrumentRow, error) {
	row := db.QueryRowContext(ctx, getValidInstrument, id)
	var i GetValidInstrumentRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.IconPath,
		&i.PluralName,
		&i.UsableForStorage,
		&i.Slug,
		&i.DisplayInSummaryLists,
		&i.IncludeInGeneratedInstructions,
		&i.LastIndexedAt,
		&i.CreatedAt,
		&i.LastUpdatedAt,
		&i.ArchivedAt,
	)
	return &i, err
}

const getValidInstruments = `-- name: GetValidInstruments :many

SELECT
	valid_instruments.id,
	valid_instruments.name,
	valid_instruments.description,
	valid_instruments.icon_path,
	valid_instruments.plural_name,
	valid_instruments.usable_for_storage,
	valid_instruments.slug,
	valid_instruments.display_in_summary_lists,
	valid_instruments.include_in_generated_instructions,
	valid_instruments.last_indexed_at,
	valid_instruments.created_at,
	valid_instruments.last_updated_at,
	valid_instruments.archived_at,
	(
		SELECT COUNT(valid_instruments.id)
		FROM valid_instruments
		WHERE valid_instruments.archived_at IS NULL
			AND valid_instruments.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
			AND valid_instruments.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
			AND (
				valid_instruments.last_updated_at IS NULL
				OR valid_instruments.last_updated_at > COALESCE($3, (SELECT NOW() - '999 years'::INTERVAL))
			)
			AND (
				valid_instruments.last_updated_at IS NULL
				OR valid_instruments.last_updated_at < COALESCE($4, (SELECT NOW() + '999 years'::INTERVAL))
			)
	) AS filtered_count,
	(
		SELECT COUNT(valid_instruments.id)
		FROM valid_instruments
		WHERE valid_instruments.archived_at IS NULL
	) AS total_count
FROM valid_instruments
WHERE
	valid_instruments.archived_at IS NULL
	AND valid_instruments.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
	AND valid_instruments.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
	AND (
		valid_instruments.last_updated_at IS NULL
		OR valid_instruments.last_updated_at > COALESCE($4, (SELECT NOW() - '999 years'::INTERVAL))
	)
	AND (
		valid_instruments.last_updated_at IS NULL
		OR valid_instruments.last_updated_at < COALESCE($3, (SELECT NOW() + '999 years'::INTERVAL))
	)
GROUP BY valid_instruments.id
ORDER BY valid_instruments.id
LIMIT $6
OFFSET $5
`

type GetValidInstrumentsParams struct {
	CreatedAfter  sql.NullTime
	CreatedBefore sql.NullTime
	UpdatedBefore sql.NullTime
	UpdatedAfter  sql.NullTime
	QueryOffset   sql.NullInt32
	QueryLimit    sql.NullInt32
}

type GetValidInstrumentsRow struct {
	CreatedAt                      time.Time
	LastIndexedAt                  sql.NullTime
	ArchivedAt                     sql.NullTime
	LastUpdatedAt                  sql.NullTime
	IconPath                       string
	Slug                           string
	PluralName                     string
	ID                             string
	Description                    string
	Name                           string
	FilteredCount                  int64
	TotalCount                     int64
	UsableForStorage               bool
	DisplayInSummaryLists          bool
	IncludeInGeneratedInstructions bool
}

func (q *Queries) GetValidInstruments(ctx context.Context, db DBTX, arg *GetValidInstrumentsParams) ([]*GetValidInstrumentsRow, error) {
	rows, err := db.QueryContext(ctx, getValidInstruments,
		arg.CreatedAfter,
		arg.CreatedBefore,
		arg.UpdatedBefore,
		arg.UpdatedAfter,
		arg.QueryOffset,
		arg.QueryLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetValidInstrumentsRow{}
	for rows.Next() {
		var i GetValidInstrumentsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.IconPath,
			&i.PluralName,
			&i.UsableForStorage,
			&i.Slug,
			&i.DisplayInSummaryLists,
			&i.IncludeInGeneratedInstructions,
			&i.LastIndexedAt,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
			&i.FilteredCount,
			&i.TotalCount,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getValidInstrumentsNeedingIndexing = `-- name: GetValidInstrumentsNeedingIndexing :many

SELECT valid_instruments.id
FROM valid_instruments
WHERE valid_instruments.archived_at IS NULL
	AND (
	valid_instruments.last_indexed_at IS NULL
	OR valid_instruments.last_indexed_at < NOW() - '24 hours'::INTERVAL
)
`

func (q *Queries) GetValidInstrumentsNeedingIndexing(ctx context.Context, db DBTX) ([]string, error) {
	rows, err := db.QueryContext(ctx, getValidInstrumentsNeedingIndexing)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getValidInstrumentsWithIDs = `-- name: GetValidInstrumentsWithIDs :many

SELECT
	valid_instruments.id,
	valid_instruments.name,
	valid_instruments.description,
	valid_instruments.icon_path,
	valid_instruments.plural_name,
	valid_instruments.usable_for_storage,
	valid_instruments.slug,
	valid_instruments.display_in_summary_lists,
	valid_instruments.include_in_generated_instructions,
	valid_instruments.last_indexed_at,
	valid_instruments.created_at,
	valid_instruments.last_updated_at,
	valid_instruments.archived_at
FROM valid_instruments
WHERE valid_instruments.archived_at IS NULL
	AND valid_instruments.id = ANY($1::text[])
`

type GetValidInstrumentsWithIDsRow struct {
	CreatedAt                      time.Time
	LastIndexedAt                  sql.NullTime
	ArchivedAt                     sql.NullTime
	LastUpdatedAt                  sql.NullTime
	IconPath                       string
	Slug                           string
	PluralName                     string
	ID                             string
	Description                    string
	Name                           string
	UsableForStorage               bool
	DisplayInSummaryLists          bool
	IncludeInGeneratedInstructions bool
}

func (q *Queries) GetValidInstrumentsWithIDs(ctx context.Context, db DBTX, ids []string) ([]*GetValidInstrumentsWithIDsRow, error) {
	rows, err := db.QueryContext(ctx, getValidInstrumentsWithIDs, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetValidInstrumentsWithIDsRow{}
	for rows.Next() {
		var i GetValidInstrumentsWithIDsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.IconPath,
			&i.PluralName,
			&i.UsableForStorage,
			&i.Slug,
			&i.DisplayInSummaryLists,
			&i.IncludeInGeneratedInstructions,
			&i.LastIndexedAt,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const searchForValidInstruments = `-- name: SearchForValidInstruments :many

SELECT
	valid_instruments.id,
	valid_instruments.name,
	valid_instruments.description,
	valid_instruments.icon_path,
	valid_instruments.plural_name,
	valid_instruments.usable_for_storage,
	valid_instruments.slug,
	valid_instruments.display_in_summary_lists,
	valid_instruments.include_in_generated_instructions,
	valid_instruments.last_indexed_at,
	valid_instruments.created_at,
	valid_instruments.last_updated_at,
	valid_instruments.archived_at
FROM valid_instruments
WHERE valid_instruments.name ILIKE '%' || $1::text || '%'
	AND valid_instruments.archived_at IS NULL
LIMIT 50
`

type SearchForValidInstrumentsRow struct {
	CreatedAt                      time.Time
	LastIndexedAt                  sql.NullTime
	ArchivedAt                     sql.NullTime
	LastUpdatedAt                  sql.NullTime
	IconPath                       string
	Slug                           string
	PluralName                     string
	ID                             string
	Description                    string
	Name                           string
	UsableForStorage               bool
	DisplayInSummaryLists          bool
	IncludeInGeneratedInstructions bool
}

func (q *Queries) SearchForValidInstruments(ctx context.Context, db DBTX, nameQuery string) ([]*SearchForValidInstrumentsRow, error) {
	rows, err := db.QueryContext(ctx, searchForValidInstruments, nameQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*SearchForValidInstrumentsRow{}
	for rows.Next() {
		var i SearchForValidInstrumentsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.IconPath,
			&i.PluralName,
			&i.UsableForStorage,
			&i.Slug,
			&i.DisplayInSummaryLists,
			&i.IncludeInGeneratedInstructions,
			&i.LastIndexedAt,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateValidInstrument = `-- name: UpdateValidInstrument :execrows

UPDATE valid_instruments SET
	name = $1,
	description = $2,
	icon_path = $3,
	plural_name = $4,
	usable_for_storage = $5,
	slug = $6,
	display_in_summary_lists = $7,
	include_in_generated_instructions = $8,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $9
`

type UpdateValidInstrumentParams struct {
	Name                           string
	Description                    string
	IconPath                       string
	PluralName                     string
	Slug                           string
	ID                             string
	UsableForStorage               bool
	DisplayInSummaryLists          bool
	IncludeInGeneratedInstructions bool
}

func (q *Queries) UpdateValidInstrument(ctx context.Context, db DBTX, arg *UpdateValidInstrumentParams) (int64, error) {
	result, err := db.ExecContext(ctx, updateValidInstrument,
		arg.Name,
		arg.Description,
		arg.IconPath,
		arg.PluralName,
		arg.UsableForStorage,
		arg.Slug,
		arg.DisplayInSummaryLists,
		arg.IncludeInGeneratedInstructions,
		arg.ID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const updateValidInstrumentLastIndexedAt = `-- name: UpdateValidInstrumentLastIndexedAt :execrows

UPDATE valid_instruments SET last_indexed_at = NOW() WHERE id = $1 AND archived_at IS NULL
`

func (q *Queries) UpdateValidInstrumentLastIndexedAt(ctx context.Context, db DBTX, id string) (int64, error) {
	result, err := db.ExecContext(ctx, updateValidInstrumentLastIndexedAt, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
