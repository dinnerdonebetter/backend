// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: valid_ingredient_measurement_units.sql

package generated

import (
	"context"
	"database/sql"
)

const ArchiveValidIngredientMeasurementUnit = `-- name: ArchiveValidIngredientMeasurementUnit :exec
UPDATE valid_ingredient_measurement_units SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1
`

func (q *Queries) ArchiveValidIngredientMeasurementUnit(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, ArchiveValidIngredientMeasurementUnit, id)
	return err
}

const CreateValidIngredientMeasurementUnit = `-- name: CreateValidIngredientMeasurementUnit :exec
INSERT INTO valid_ingredient_measurement_units (id,notes,minimum_allowable_quantity,maximum_allowable_quantity,valid_measurement_unit_id,valid_ingredient_id) VALUES ($1,$2,$3,$4,$5,$6)
`

type CreateValidIngredientMeasurementUnitParams struct {
	ID                       string
	Notes                    string
	MinimumAllowableQuantity float64
	MaximumAllowableQuantity float64
	ValidMeasurementUnitID   string
	ValidIngredientID        string
}

func (q *Queries) CreateValidIngredientMeasurementUnit(ctx context.Context, arg *CreateValidIngredientMeasurementUnitParams) error {
	_, err := q.db.ExecContext(ctx, CreateValidIngredientMeasurementUnit,
		arg.ID,
		arg.Notes,
		arg.MinimumAllowableQuantity,
		arg.MaximumAllowableQuantity,
		arg.ValidMeasurementUnitID,
		arg.ValidIngredientID,
	)
	return err
}

const GetTotalValidIngredientMeasurementUnitsCount = `-- name: GetTotalValidIngredientMeasurementUnitsCount :one
SELECT COUNT(valid_ingredient_measurement_units.id) FROM valid_ingredient_measurement_units WHERE valid_ingredient_measurement_units.archived_on IS NULL
`

func (q *Queries) GetTotalValidIngredientMeasurementUnitsCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, GetTotalValidIngredientMeasurementUnitsCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const GetValidIngredientMeasurementUnit = `-- name: GetValidIngredientMeasurementUnit :one
SELECT
    valid_ingredient_measurement_units.id,
    valid_ingredient_measurement_units.notes,
    valid_ingredient_measurement_units.minimum_allowable_quantity,
    valid_ingredient_measurement_units.maximum_allowable_quantity,
    valid_measurement_units.id,
    valid_measurement_units.name,
    valid_measurement_units.description,
    valid_measurement_units.volumetric,
    valid_measurement_units.icon_path,
    valid_measurement_units.created_on,
    valid_measurement_units.last_updated_on,
    valid_measurement_units.archived_on,
    valid_ingredients.id,
    valid_ingredients.name,
    valid_ingredients.description,
    valid_ingredients.warning,
    valid_ingredients.contains_egg,
    valid_ingredients.contains_dairy,
    valid_ingredients.contains_peanut,
    valid_ingredients.contains_tree_nut,
    valid_ingredients.contains_soy,
    valid_ingredients.contains_wheat,
    valid_ingredients.contains_shellfish,
    valid_ingredients.contains_sesame,
    valid_ingredients.contains_fish,
    valid_ingredients.contains_gluten,
    valid_ingredients.animal_flesh,
    valid_ingredients.volumetric,
    valid_ingredients.is_liquid,
    valid_ingredients.icon_path,
    valid_ingredients.created_on,
    valid_ingredients.last_updated_on,
    valid_ingredients.archived_on,
    valid_ingredient_measurement_units.created_on,
    valid_ingredient_measurement_units.last_updated_on,
    valid_ingredient_measurement_units.archived_on
FROM valid_ingredient_measurement_units
         JOIN valid_measurement_units ON valid_ingredient_measurement_units.valid_measurement_unit_id = valid_measurement_units.id
         JOIN valid_ingredients ON valid_ingredient_measurement_units.valid_ingredient_id = valid_ingredients.id
WHERE valid_ingredient_measurement_units.archived_on IS NULL
  AND valid_ingredient_measurement_units.id = $1
`

type GetValidIngredientMeasurementUnitRow struct {
	ID                       string
	Notes                    string
	MinimumAllowableQuantity float64
	MaximumAllowableQuantity float64
	ID_2                     string
	Name                     string
	Description              string
	Volumetric               sql.NullBool
	IconPath                 string
	CreatedOn                int64
	LastUpdatedOn            sql.NullInt64
	ArchivedOn               sql.NullInt64
	ID_3                     string
	Name_2                   string
	Description_2            string
	Warning                  string
	ContainsEgg              bool
	ContainsDairy            bool
	ContainsPeanut           bool
	ContainsTreeNut          bool
	ContainsSoy              bool
	ContainsWheat            bool
	ContainsShellfish        bool
	ContainsSesame           bool
	ContainsFish             bool
	ContainsGluten           bool
	AnimalFlesh              bool
	Volumetric_2             bool
	IsLiquid                 sql.NullBool
	IconPath_2               string
	CreatedOn_2              int64
	LastUpdatedOn_2          sql.NullInt64
	ArchivedOn_2             sql.NullInt64
	CreatedOn_3              int64
	LastUpdatedOn_3          sql.NullInt64
	ArchivedOn_3             sql.NullInt64
}

func (q *Queries) GetValidIngredientMeasurementUnit(ctx context.Context, id string) (*GetValidIngredientMeasurementUnitRow, error) {
	row := q.db.QueryRowContext(ctx, GetValidIngredientMeasurementUnit, id)
	var i GetValidIngredientMeasurementUnitRow
	err := row.Scan(
		&i.ID,
		&i.Notes,
		&i.MinimumAllowableQuantity,
		&i.MaximumAllowableQuantity,
		&i.ID_2,
		&i.Name,
		&i.Description,
		&i.Volumetric,
		&i.IconPath,
		&i.CreatedOn,
		&i.LastUpdatedOn,
		&i.ArchivedOn,
		&i.ID_3,
		&i.Name_2,
		&i.Description_2,
		&i.Warning,
		&i.ContainsEgg,
		&i.ContainsDairy,
		&i.ContainsPeanut,
		&i.ContainsTreeNut,
		&i.ContainsSoy,
		&i.ContainsWheat,
		&i.ContainsShellfish,
		&i.ContainsSesame,
		&i.ContainsFish,
		&i.ContainsGluten,
		&i.AnimalFlesh,
		&i.Volumetric_2,
		&i.IsLiquid,
		&i.IconPath_2,
		&i.CreatedOn_2,
		&i.LastUpdatedOn_2,
		&i.ArchivedOn_2,
		&i.CreatedOn_3,
		&i.LastUpdatedOn_3,
		&i.ArchivedOn_3,
	)
	return &i, err
}

const UpdateValidIngredientMeasurementUnit = `-- name: UpdateValidIngredientMeasurementUnit :exec
UPDATE valid_ingredient_measurement_units SET notes = $1, minimum_allowable_quantity = $2, maximum_allowable_quantity = $3, valid_measurement_unit_id = $4, valid_ingredient_id = $5, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $6
`

type UpdateValidIngredientMeasurementUnitParams struct {
	Notes                    string
	MinimumAllowableQuantity float64
	MaximumAllowableQuantity float64
	ValidMeasurementUnitID   string
	ValidIngredientID        string
	ID                       string
}

func (q *Queries) UpdateValidIngredientMeasurementUnit(ctx context.Context, arg *UpdateValidIngredientMeasurementUnitParams) error {
	_, err := q.db.ExecContext(ctx, UpdateValidIngredientMeasurementUnit,
		arg.Notes,
		arg.MinimumAllowableQuantity,
		arg.MaximumAllowableQuantity,
		arg.ValidMeasurementUnitID,
		arg.ValidIngredientID,
		arg.ID,
	)
	return err
}

const ValidIngredientMeasurementUnitExists = `-- name: ValidIngredientMeasurementUnitExists :one
SELECT EXISTS ( SELECT valid_ingredient_measurement_units.id FROM valid_ingredient_measurement_units WHERE valid_ingredient_measurement_units.archived_on IS NULL AND valid_ingredient_measurement_units.id = $1 )
`

func (q *Queries) ValidIngredientMeasurementUnitExists(ctx context.Context, id string) (bool, error) {
	row := q.db.QueryRowContext(ctx, ValidIngredientMeasurementUnitExists, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}
