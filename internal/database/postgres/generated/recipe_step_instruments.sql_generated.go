// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: recipe_step_instruments.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const archiveRecipeStepInstrument = `-- name: ArchiveRecipeStepInstrument :execrows

UPDATE recipe_step_instruments SET archived_at = NOW() WHERE archived_at IS NULL AND belongs_to_recipe_step = $1 AND id = $2
`

type ArchiveRecipeStepInstrumentParams struct {
	BelongsToRecipeStep string
	ID                  string
}

func (q *Queries) ArchiveRecipeStepInstrument(ctx context.Context, db DBTX, arg *ArchiveRecipeStepInstrumentParams) (int64, error) {
	result, err := db.ExecContext(ctx, archiveRecipeStepInstrument, arg.BelongsToRecipeStep, arg.ID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const checkRecipeStepInstrumentExistence = `-- name: CheckRecipeStepInstrumentExistence :one

SELECT EXISTS ( SELECT recipe_step_instruments.id FROM recipe_step_instruments JOIN recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_instruments.archived_at IS NULL AND recipe_step_instruments.belongs_to_recipe_step = $1 AND recipe_step_instruments.id = $2 AND recipe_steps.archived_at IS NULL AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $1 AND recipes.archived_at IS NULL AND recipes.id = $3 )
`

type CheckRecipeStepInstrumentExistenceParams struct {
	RecipeStepID           string
	RecipeStepInstrumentID string
	RecipeID               string
}

func (q *Queries) CheckRecipeStepInstrumentExistence(ctx context.Context, db DBTX, arg *CheckRecipeStepInstrumentExistenceParams) (bool, error) {
	row := db.QueryRowContext(ctx, checkRecipeStepInstrumentExistence, arg.RecipeStepID, arg.RecipeStepInstrumentID, arg.RecipeID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createRecipeStepInstrument = `-- name: CreateRecipeStepInstrument :exec

INSERT INTO recipe_step_instruments
(id,instrument_id,recipe_step_product_id,"name",notes,preference_rank,optional,option_index,minimum_quantity,maximum_quantity,belongs_to_recipe_step)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
`

type CreateRecipeStepInstrumentParams struct {
	ID                  string
	Name                string
	Notes               string
	BelongsToRecipeStep string
	InstrumentID        sql.NullString
	RecipeStepProductID sql.NullString
	MaximumQuantity     sql.NullInt32
	PreferenceRank      int32
	OptionIndex         int32
	MinimumQuantity     int32
	Optional            bool
}

func (q *Queries) CreateRecipeStepInstrument(ctx context.Context, db DBTX, arg *CreateRecipeStepInstrumentParams) error {
	_, err := db.ExecContext(ctx, createRecipeStepInstrument,
		arg.ID,
		arg.InstrumentID,
		arg.RecipeStepProductID,
		arg.Name,
		arg.Notes,
		arg.PreferenceRank,
		arg.Optional,
		arg.OptionIndex,
		arg.MinimumQuantity,
		arg.MaximumQuantity,
		arg.BelongsToRecipeStep,
	)
	return err
}

const getRecipeStepInstrument = `-- name: GetRecipeStepInstrument :one

SELECT
	recipe_step_instruments.id,
	valid_instruments.id as valid_instrument_id,
	valid_instruments.name as valid_instrument_name,
	valid_instruments.plural_name as valid_instrument_plural_name,
	valid_instruments.description as valid_instrument_description,
	valid_instruments.icon_path as valid_instrument_icon_path,
	valid_instruments.usable_for_storage as valid_instrument_usable_for_storage,
    valid_instruments.display_in_summary_lists as valid_instrument_display_in_summary_lists,
	valid_instruments.include_in_generated_instructions as valid_instrument_include_in_generated_instructions,
	valid_instruments.slug as valid_instrument_slug,
	valid_instruments.created_at as valid_instrument_created_at,
	valid_instruments.last_updated_at as valid_instrument_last_updated_at,
	valid_instruments.archived_at as valid_instrument_archived_at,
	recipe_step_instruments.recipe_step_product_id,
	recipe_step_instruments.name,
	recipe_step_instruments.notes,
	recipe_step_instruments.preference_rank,
	recipe_step_instruments.optional,
	recipe_step_instruments.minimum_quantity,
	recipe_step_instruments.maximum_quantity,
	recipe_step_instruments.option_index,
	recipe_step_instruments.created_at,
	recipe_step_instruments.last_updated_at,
	recipe_step_instruments.archived_at,
	recipe_step_instruments.belongs_to_recipe_step
FROM recipe_step_instruments
	LEFT JOIN valid_instruments ON recipe_step_instruments.instrument_id=valid_instruments.id
	JOIN recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id
	JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
WHERE recipe_step_instruments.archived_at IS NULL
	AND recipe_step_instruments.belongs_to_recipe_step = $1
	AND recipe_step_instruments.id = $2
	AND recipe_steps.archived_at IS NULL
	AND recipe_steps.belongs_to_recipe = $3
	AND recipe_steps.id = $1
	AND recipes.archived_at IS NULL
	AND recipes.id = $3
`

type GetRecipeStepInstrumentParams struct {
	RecipeStepID           string
	RecipeStepInstrumentID string
	RecipeID               string
}

type GetRecipeStepInstrumentRow struct {
	CreatedAt                                     time.Time
	ValidInstrumentArchivedAt                     sql.NullTime
	ArchivedAt                                    sql.NullTime
	LastUpdatedAt                                 sql.NullTime
	ValidInstrumentLastUpdatedAt                  sql.NullTime
	ValidInstrumentCreatedAt                      sql.NullTime
	Name                                          string
	BelongsToRecipeStep                           string
	ID                                            string
	Notes                                         string
	ValidInstrumentDescription                    sql.NullString
	ValidInstrumentPluralName                     sql.NullString
	ValidInstrumentIconPath                       sql.NullString
	ValidInstrumentSlug                           sql.NullString
	RecipeStepProductID                           sql.NullString
	ValidInstrumentName                           sql.NullString
	ValidInstrumentID                             sql.NullString
	MaximumQuantity                               sql.NullInt32
	MinimumQuantity                               int32
	OptionIndex                                   int32
	PreferenceRank                                int32
	ValidInstrumentIncludeInGeneratedInstructions sql.NullBool
	ValidInstrumentDisplayInSummaryLists          sql.NullBool
	ValidInstrumentUsableForStorage               sql.NullBool
	Optional                                      bool
}

func (q *Queries) GetRecipeStepInstrument(ctx context.Context, db DBTX, arg *GetRecipeStepInstrumentParams) (*GetRecipeStepInstrumentRow, error) {
	row := db.QueryRowContext(ctx, getRecipeStepInstrument, arg.RecipeStepID, arg.RecipeStepInstrumentID, arg.RecipeID)
	var i GetRecipeStepInstrumentRow
	err := row.Scan(
		&i.ID,
		&i.ValidInstrumentID,
		&i.ValidInstrumentName,
		&i.ValidInstrumentPluralName,
		&i.ValidInstrumentDescription,
		&i.ValidInstrumentIconPath,
		&i.ValidInstrumentUsableForStorage,
		&i.ValidInstrumentDisplayInSummaryLists,
		&i.ValidInstrumentIncludeInGeneratedInstructions,
		&i.ValidInstrumentSlug,
		&i.ValidInstrumentCreatedAt,
		&i.ValidInstrumentLastUpdatedAt,
		&i.ValidInstrumentArchivedAt,
		&i.RecipeStepProductID,
		&i.Name,
		&i.Notes,
		&i.PreferenceRank,
		&i.Optional,
		&i.MinimumQuantity,
		&i.MaximumQuantity,
		&i.OptionIndex,
		&i.CreatedAt,
		&i.LastUpdatedAt,
		&i.ArchivedAt,
		&i.BelongsToRecipeStep,
	)
	return &i, err
}

const getRecipeStepInstruments = `-- name: GetRecipeStepInstruments :many

SELECT
    recipe_step_instruments.id,
    valid_instruments.id as valid_instrument_id,
    valid_instruments.name as valid_instrument_name,
    valid_instruments.plural_name as valid_instrument_plural_name,
    valid_instruments.description as valid_instrument_description,
    valid_instruments.icon_path as valid_instrument_icon_path,
    valid_instruments.usable_for_storage as valid_instrument_usable_for_storage,
    valid_instruments.display_in_summary_lists as valid_instrument_display_in_summary_lists,
    valid_instruments.include_in_generated_instructions as valid_instrument_include_in_generated_instructions,
    valid_instruments.slug as valid_instrument_slug,
    valid_instruments.created_at as valid_instrument_created_at,
    valid_instruments.last_updated_at as valid_instrument_last_updated_at,
    valid_instruments.archived_at as valid_instrument_archived_at,
    recipe_step_instruments.recipe_step_product_id,
    recipe_step_instruments.name,
    recipe_step_instruments.notes,
    recipe_step_instruments.preference_rank,
    recipe_step_instruments.optional,
    recipe_step_instruments.minimum_quantity,
    recipe_step_instruments.maximum_quantity,
    recipe_step_instruments.option_index,
    recipe_step_instruments.created_at,
    recipe_step_instruments.last_updated_at,
    recipe_step_instruments.archived_at,
    recipe_step_instruments.belongs_to_recipe_step,
    (
        SELECT
            COUNT(recipe_step_instruments.id)
        FROM
            recipe_step_instruments
        WHERE
            recipe_step_instruments.archived_at IS NULL
            AND recipe_step_instruments.belongs_to_recipe_step = $1
            AND recipe_step_instruments.created_at > COALESCE($2, (SELECT NOW() - interval '999 years'))
            AND recipe_step_instruments.created_at < COALESCE($3, (SELECT NOW() + interval '999 years'))
            AND (recipe_step_instruments.last_updated_at IS NULL OR recipe_step_instruments.last_updated_at > COALESCE($4, (SELECT NOW() - interval '999 years')))
            AND (recipe_step_instruments.last_updated_at IS NULL OR recipe_step_instruments.last_updated_at < COALESCE($5, (SELECT NOW() + interval '999 years')))
    ) as filtered_count,
    (
        SELECT
            COUNT(recipe_step_instruments.id)
        FROM
            recipe_step_instruments
        WHERE
            recipe_step_instruments.archived_at IS NULL
    ) as total_count
FROM recipe_step_instruments
    LEFT JOIN valid_instruments ON recipe_step_instruments.instrument_id=valid_instruments.id
    JOIN recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id
    JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
WHERE
    recipe_step_instruments.archived_at IS NULL
    AND recipe_step_instruments.belongs_to_recipe_step = $1
    AND recipe_step_instruments.created_at > COALESCE($2, (SELECT NOW() - interval '999 years'))
    AND recipe_step_instruments.created_at < COALESCE($3, (SELECT NOW() + interval '999 years'))
    AND (recipe_step_instruments.last_updated_at IS NULL OR recipe_step_instruments.last_updated_at > COALESCE($4, (SELECT NOW() - interval '999 years')))
    AND (recipe_step_instruments.last_updated_at IS NULL OR recipe_step_instruments.last_updated_at < COALESCE($5, (SELECT NOW() + interval '999 years')))
    AND recipe_steps.archived_at IS NULL
    AND recipe_steps.belongs_to_recipe = $6
    AND recipe_steps.id = $1
    AND recipes.archived_at IS NULL
    AND recipes.id = $6
    OFFSET $7
    LIMIT $8
`

type GetRecipeStepInstrumentsParams struct {
	RecipeStepID  string
	CreatedAfter  sql.NullTime
	CreatedBefore sql.NullTime
	UpdatedAfter  sql.NullTime
	UpdatedBefore sql.NullTime
	RecipeID      string
	QueryOffset   sql.NullInt32
	QueryLimit    sql.NullInt32
}

type GetRecipeStepInstrumentsRow struct {
	CreatedAt                                     time.Time
	ValidInstrumentCreatedAt                      sql.NullTime
	LastUpdatedAt                                 sql.NullTime
	ArchivedAt                                    sql.NullTime
	ValidInstrumentArchivedAt                     sql.NullTime
	ValidInstrumentLastUpdatedAt                  sql.NullTime
	Notes                                         string
	Name                                          string
	BelongsToRecipeStep                           string
	ID                                            string
	ValidInstrumentName                           sql.NullString
	ValidInstrumentIconPath                       sql.NullString
	ValidInstrumentDescription                    sql.NullString
	RecipeStepProductID                           sql.NullString
	ValidInstrumentPluralName                     sql.NullString
	ValidInstrumentID                             sql.NullString
	ValidInstrumentSlug                           sql.NullString
	TotalCount                                    int64
	FilteredCount                                 int64
	MaximumQuantity                               sql.NullInt32
	MinimumQuantity                               int32
	OptionIndex                                   int32
	PreferenceRank                                int32
	ValidInstrumentIncludeInGeneratedInstructions sql.NullBool
	ValidInstrumentDisplayInSummaryLists          sql.NullBool
	ValidInstrumentUsableForStorage               sql.NullBool
	Optional                                      bool
}

func (q *Queries) GetRecipeStepInstruments(ctx context.Context, db DBTX, arg *GetRecipeStepInstrumentsParams) ([]*GetRecipeStepInstrumentsRow, error) {
	rows, err := db.QueryContext(ctx, getRecipeStepInstruments,
		arg.RecipeStepID,
		arg.CreatedAfter,
		arg.CreatedBefore,
		arg.UpdatedAfter,
		arg.UpdatedBefore,
		arg.RecipeID,
		arg.QueryOffset,
		arg.QueryLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetRecipeStepInstrumentsRow{}
	for rows.Next() {
		var i GetRecipeStepInstrumentsRow
		if err := rows.Scan(
			&i.ID,
			&i.ValidInstrumentID,
			&i.ValidInstrumentName,
			&i.ValidInstrumentPluralName,
			&i.ValidInstrumentDescription,
			&i.ValidInstrumentIconPath,
			&i.ValidInstrumentUsableForStorage,
			&i.ValidInstrumentDisplayInSummaryLists,
			&i.ValidInstrumentIncludeInGeneratedInstructions,
			&i.ValidInstrumentSlug,
			&i.ValidInstrumentCreatedAt,
			&i.ValidInstrumentLastUpdatedAt,
			&i.ValidInstrumentArchivedAt,
			&i.RecipeStepProductID,
			&i.Name,
			&i.Notes,
			&i.PreferenceRank,
			&i.Optional,
			&i.MinimumQuantity,
			&i.MaximumQuantity,
			&i.OptionIndex,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
			&i.BelongsToRecipeStep,
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

const getRecipeStepInstrumentsForRecipe = `-- name: GetRecipeStepInstrumentsForRecipe :many

SELECT
    recipe_step_instruments.id,
    valid_instruments.id as valid_instrument_id,
    valid_instruments.name as valid_instrument_name,
    valid_instruments.plural_name as valid_instrument_plural_name,
    valid_instruments.description as valid_instrument_description,
    valid_instruments.icon_path as valid_instrument_icon_path,
    valid_instruments.usable_for_storage as valid_instrument_usable_for_storage,
    valid_instruments.display_in_summary_lists as valid_instrument_display_in_summary_lists,
    valid_instruments.include_in_generated_instructions as valid_instrument_include_in_generated_instructions,
    valid_instruments.slug as valid_instrument_slug,
    valid_instruments.created_at as valid_instrument_created_at,
    valid_instruments.last_updated_at as valid_instrument_last_updated_at,
    valid_instruments.archived_at as valid_instrument_archived_at,
    recipe_step_instruments.recipe_step_product_id,
    recipe_step_instruments.name,
    recipe_step_instruments.notes,
    recipe_step_instruments.preference_rank,
    recipe_step_instruments.optional,
    recipe_step_instruments.minimum_quantity,
    recipe_step_instruments.maximum_quantity,
    recipe_step_instruments.option_index,
    recipe_step_instruments.created_at,
    recipe_step_instruments.last_updated_at,
    recipe_step_instruments.archived_at,
    recipe_step_instruments.belongs_to_recipe_step
FROM recipe_step_instruments
	LEFT JOIN valid_instruments ON recipe_step_instruments.instrument_id=valid_instruments.id
	JOIN recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id
	JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
WHERE recipe_step_instruments.archived_at IS NULL
	AND recipe_steps.archived_at IS NULL
	AND recipe_steps.belongs_to_recipe = $1
	AND recipes.archived_at IS NULL
	AND recipes.id = $1
`

type GetRecipeStepInstrumentsForRecipeRow struct {
	CreatedAt                                     time.Time
	ValidInstrumentArchivedAt                     sql.NullTime
	ArchivedAt                                    sql.NullTime
	LastUpdatedAt                                 sql.NullTime
	ValidInstrumentLastUpdatedAt                  sql.NullTime
	ValidInstrumentCreatedAt                      sql.NullTime
	Name                                          string
	BelongsToRecipeStep                           string
	ID                                            string
	Notes                                         string
	ValidInstrumentDescription                    sql.NullString
	ValidInstrumentPluralName                     sql.NullString
	ValidInstrumentIconPath                       sql.NullString
	ValidInstrumentSlug                           sql.NullString
	RecipeStepProductID                           sql.NullString
	ValidInstrumentName                           sql.NullString
	ValidInstrumentID                             sql.NullString
	MaximumQuantity                               sql.NullInt32
	MinimumQuantity                               int32
	OptionIndex                                   int32
	PreferenceRank                                int32
	ValidInstrumentIncludeInGeneratedInstructions sql.NullBool
	ValidInstrumentDisplayInSummaryLists          sql.NullBool
	ValidInstrumentUsableForStorage               sql.NullBool
	Optional                                      bool
}

func (q *Queries) GetRecipeStepInstrumentsForRecipe(ctx context.Context, db DBTX, belongsToRecipe string) ([]*GetRecipeStepInstrumentsForRecipeRow, error) {
	rows, err := db.QueryContext(ctx, getRecipeStepInstrumentsForRecipe, belongsToRecipe)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetRecipeStepInstrumentsForRecipeRow{}
	for rows.Next() {
		var i GetRecipeStepInstrumentsForRecipeRow
		if err := rows.Scan(
			&i.ID,
			&i.ValidInstrumentID,
			&i.ValidInstrumentName,
			&i.ValidInstrumentPluralName,
			&i.ValidInstrumentDescription,
			&i.ValidInstrumentIconPath,
			&i.ValidInstrumentUsableForStorage,
			&i.ValidInstrumentDisplayInSummaryLists,
			&i.ValidInstrumentIncludeInGeneratedInstructions,
			&i.ValidInstrumentSlug,
			&i.ValidInstrumentCreatedAt,
			&i.ValidInstrumentLastUpdatedAt,
			&i.ValidInstrumentArchivedAt,
			&i.RecipeStepProductID,
			&i.Name,
			&i.Notes,
			&i.PreferenceRank,
			&i.Optional,
			&i.MinimumQuantity,
			&i.MaximumQuantity,
			&i.OptionIndex,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
			&i.BelongsToRecipeStep,
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

const updateRecipeStepInstrument = `-- name: UpdateRecipeStepInstrument :execrows

UPDATE recipe_step_instruments SET
	instrument_id = $1,
	recipe_step_product_id = $2,
	name = $3,
	notes = $4,
	preference_rank = $5,
	optional = $6,
	option_index = $7,
	minimum_quantity = $8,
	maximum_quantity = $9,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_recipe_step = $10
	AND id = $11
`

type UpdateRecipeStepInstrumentParams struct {
	Name                string
	Notes               string
	BelongsToRecipeStep string
	ID                  string
	InstrumentID        sql.NullString
	RecipeStepProductID sql.NullString
	MaximumQuantity     sql.NullInt32
	PreferenceRank      int32
	OptionIndex         int32
	MinimumQuantity     int32
	Optional            bool
}

func (q *Queries) UpdateRecipeStepInstrument(ctx context.Context, db DBTX, arg *UpdateRecipeStepInstrumentParams) (int64, error) {
	result, err := db.ExecContext(ctx, updateRecipeStepInstrument,
		arg.InstrumentID,
		arg.RecipeStepProductID,
		arg.Name,
		arg.Notes,
		arg.PreferenceRank,
		arg.Optional,
		arg.OptionIndex,
		arg.MinimumQuantity,
		arg.MaximumQuantity,
		arg.BelongsToRecipeStep,
		arg.ID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}