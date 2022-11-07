// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: valid_preparation_instruments_get_one.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const GetValidPreparationInstrument = `-- name: GetValidPreparationInstrument :exec
SELECT
	valid_preparation_instruments.id,
	valid_preparation_instruments.notes,
	valid_preparations.id,
	valid_preparations.name,
	valid_preparations.description,
	valid_preparations.icon_path,
	valid_preparations.yields_nothing,
	valid_preparations.restrict_to_ingredients,
	valid_preparations.zero_ingredients_allowable,
	valid_preparations.past_tense,
	valid_preparations.created_at,
	valid_preparations.last_updated_at,
	valid_preparations.archived_at,
	valid_instruments.id,
	valid_instruments.name,
	valid_instruments.plural_name,
	valid_instruments.description,
	valid_instruments.icon_path,
	valid_instruments.usable_for_storage,
	valid_instruments.created_at,
	valid_instruments.last_updated_at,
	valid_instruments.archived_at,
	valid_preparation_instruments.created_at,
	valid_preparation_instruments.last_updated_at,
	valid_preparation_instruments.archived_at
FROM
	valid_preparation_instruments
	 JOIN valid_instruments ON valid_preparation_instruments.valid_instrument_id = valid_instruments.id
	 JOIN valid_preparations ON valid_preparation_instruments.valid_preparation_id = valid_preparations.id
WHERE
	valid_preparation_instruments.archived_at IS NULL
	AND valid_preparation_instruments.id = $1
`

type GetValidPreparationInstrumentRow struct {
	CreatedAt                time.Time    `db:"created_at"`
	CreatedAt_3              time.Time    `db:"created_at_3"`
	CreatedAt_2              time.Time    `db:"created_at_2"`
	LastUpdatedAt_3          sql.NullTime `db:"last_updated_at_3"`
	ArchivedAt_2             sql.NullTime `db:"archived_at_2"`
	LastUpdatedAt_2          sql.NullTime `db:"last_updated_at_2"`
	ArchivedAt               sql.NullTime `db:"archived_at"`
	LastUpdatedAt            sql.NullTime `db:"last_updated_at"`
	ArchivedAt_3             sql.NullTime `db:"archived_at_3"`
	IconPath                 string       `db:"icon_path"`
	Notes                    string       `db:"notes"`
	ID_2                     string       `db:"id_2"`
	Name                     string       `db:"name"`
	ID_3                     string       `db:"id_3"`
	Name_2                   string       `db:"name_2"`
	PluralName               string       `db:"plural_name"`
	Description_2            string       `db:"description_2"`
	IconPath_2               string       `db:"icon_path_2"`
	Description              string       `db:"description"`
	PastTense                string       `db:"past_tense"`
	ID                       string       `db:"id"`
	UsableForStorage         bool         `db:"usable_for_storage"`
	YieldsNothing            bool         `db:"yields_nothing"`
	RestrictToIngredients    bool         `db:"restrict_to_ingredients"`
	ZeroIngredientsAllowable bool         `db:"zero_ingredients_allowable"`
}

func (q *Queries) GetValidPreparationInstrument(ctx context.Context, db DBTX, id string) error {
	_, err := db.ExecContext(ctx, GetValidPreparationInstrument, id)
	return err
}