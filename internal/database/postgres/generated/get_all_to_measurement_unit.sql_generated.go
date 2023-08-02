// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: get_all_to_measurement_unit.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const GetAllValidMeasurementConversionsToMeasurementUnit = `-- name: GetAllValidMeasurementConversionsToMeasurementUnit :many

SELECT
	valid_measurement_conversions.id,
	valid_measurement_units_from.id,
	valid_measurement_units_from.name,
	valid_measurement_units_from.description,
	valid_measurement_units_from.volumetric,
	valid_measurement_units_from.icon_path,
	valid_measurement_units_from.universal,
	valid_measurement_units_from.metric,
	valid_measurement_units_from.imperial,
	valid_measurement_units_from.slug,
	valid_measurement_units_from.plural_name,
	valid_measurement_units_from.created_at,
	valid_measurement_units_from.last_updated_at,
	valid_measurement_units_from.archived_at,
	valid_measurement_units_to.id,
	valid_measurement_units_to.name,
	valid_measurement_units_to.description,
	valid_measurement_units_to.volumetric,
	valid_measurement_units_to.icon_path,
	valid_measurement_units_to.universal,
	valid_measurement_units_to.metric,
	valid_measurement_units_to.imperial,
	valid_measurement_units_to.slug,
	valid_measurement_units_to.plural_name,
	valid_measurement_units_to.created_at,
	valid_measurement_units_to.last_updated_at,
	valid_measurement_units_to.archived_at,
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
	valid_ingredients.animal_derived,
	valid_ingredients.plural_name,
	valid_ingredients.restrict_to_preparations,
	valid_ingredients.minimum_ideal_storage_temperature_in_celsius,
	valid_ingredients.maximum_ideal_storage_temperature_in_celsius,
	valid_ingredients.storage_instructions,
	valid_ingredients.slug,
	valid_ingredients.contains_alcohol,
	valid_ingredients.shopping_suggestions,
    valid_ingredients.is_starch,
    valid_ingredients.is_protein,
    valid_ingredients.is_grain,
    valid_ingredients.is_fruit,
    valid_ingredients.is_salt,
    valid_ingredients.is_fat,
    valid_ingredients.is_acid,
    valid_ingredients.is_heat,
	valid_ingredients.created_at,
	valid_ingredients.last_updated_at,
	valid_ingredients.archived_at,
	valid_measurement_conversions.modifier,
	valid_measurement_conversions.notes,
	valid_measurement_conversions.created_at,
	valid_measurement_conversions.last_updated_at,
	valid_measurement_conversions.archived_at
FROM valid_measurement_conversions
	     LEFT JOIN valid_ingredients ON valid_measurement_conversions.only_for_ingredient = valid_ingredients.id
	     JOIN valid_measurement_units AS valid_measurement_units_from ON valid_measurement_conversions.from_unit = valid_measurement_units_from.id
	     JOIN valid_measurement_units AS valid_measurement_units_to ON valid_measurement_conversions.to_unit = valid_measurement_units_to.id
WHERE valid_measurement_conversions.archived_at IS NULL
  AND valid_measurement_units_from.archived_at IS NULL
  AND valid_measurement_units_to.id = $1
  AND valid_measurement_units_to.archived_at IS NULL
`

type GetAllValidMeasurementConversionsToMeasurementUnitRow struct {
	CreatedAt_4                             time.Time      `db:"created_at_4"`
	CreatedAt_2                             time.Time      `db:"created_at_2"`
	CreatedAt                               time.Time      `db:"created_at"`
	LastUpdatedAt                           sql.NullTime   `db:"last_updated_at"`
	ArchivedAt                              sql.NullTime   `db:"archived_at"`
	ArchivedAt_2                            sql.NullTime   `db:"archived_at_2"`
	LastUpdatedAt_4                         sql.NullTime   `db:"last_updated_at_4"`
	LastUpdatedAt_2                         sql.NullTime   `db:"last_updated_at_2"`
	CreatedAt_3                             sql.NullTime   `db:"created_at_3"`
	ArchivedAt_4                            sql.NullTime   `db:"archived_at_4"`
	LastUpdatedAt_3                         sql.NullTime   `db:"last_updated_at_3"`
	ArchivedAt_3                            sql.NullTime   `db:"archived_at_3"`
	Slug                                    string         `db:"slug"`
	Name                                    string         `db:"name"`
	ID_3                                    string         `db:"id_3"`
	Name_2                                  string         `db:"name_2"`
	Description_2                           string         `db:"description_2"`
	Modifier                                string         `db:"modifier"`
	IconPath_2                              string         `db:"icon_path_2"`
	ID                                      string         `db:"id"`
	PluralName                              string         `db:"plural_name"`
	Notes                                   string         `db:"notes"`
	Slug_2                                  string         `db:"slug_2"`
	PluralName_2                            string         `db:"plural_name_2"`
	ID_2                                    string         `db:"id_2"`
	Description                             string         `db:"description"`
	IconPath                                string         `db:"icon_path"`
	Slug_3                                  sql.NullString `db:"slug_3"`
	Name_3                                  sql.NullString `db:"name_3"`
	Description_3                           sql.NullString `db:"description_3"`
	Warning                                 sql.NullString `db:"warning"`
	ShoppingSuggestions                     sql.NullString `db:"shopping_suggestions"`
	ID_4                                    sql.NullString `db:"id_4"`
	StorageInstructions                     sql.NullString `db:"storage_instructions"`
	MaximumIdealStorageTemperatureInCelsius sql.NullString `db:"maximum_ideal_storage_temperature_in_celsius"`
	MinimumIdealStorageTemperatureInCelsius sql.NullString `db:"minimum_ideal_storage_temperature_in_celsius"`
	PluralName_3                            sql.NullString `db:"plural_name_3"`
	IconPath_3                              sql.NullString `db:"icon_path_3"`
	ContainsSesame                          sql.NullBool   `db:"contains_sesame"`
	ContainsEgg                             sql.NullBool   `db:"contains_egg"`
	ContainsGluten                          sql.NullBool   `db:"contains_gluten"`
	AnimalFlesh                             sql.NullBool   `db:"animal_flesh"`
	Volumetric_3                            sql.NullBool   `db:"volumetric_3"`
	IsLiquid                                sql.NullBool   `db:"is_liquid"`
	ContainsShellfish                       sql.NullBool   `db:"contains_shellfish"`
	AnimalDerived                           sql.NullBool   `db:"animal_derived"`
	ContainsWheat                           sql.NullBool   `db:"contains_wheat"`
	RestrictToPreparations                  sql.NullBool   `db:"restrict_to_preparations"`
	ContainsSoy                             sql.NullBool   `db:"contains_soy"`
	ContainsTreeNut                         sql.NullBool   `db:"contains_tree_nut"`
	ContainsPeanut                          sql.NullBool   `db:"contains_peanut"`
	ContainsDairy                           sql.NullBool   `db:"contains_dairy"`
	ContainsAlcohol                         sql.NullBool   `db:"contains_alcohol"`
	ContainsFish                            sql.NullBool   `db:"contains_fish"`
	IsStarch                                sql.NullBool   `db:"is_starch"`
	IsProtein                               sql.NullBool   `db:"is_protein"`
	IsGrain                                 sql.NullBool   `db:"is_grain"`
	IsFruit                                 sql.NullBool   `db:"is_fruit"`
	IsSalt                                  sql.NullBool   `db:"is_salt"`
	IsFat                                   sql.NullBool   `db:"is_fat"`
	IsAcid                                  sql.NullBool   `db:"is_acid"`
	IsHeat                                  sql.NullBool   `db:"is_heat"`
	Volumetric                              sql.NullBool   `db:"volumetric"`
	Volumetric_2                            sql.NullBool   `db:"volumetric_2"`
	Universal_2                             bool           `db:"universal_2"`
	Metric_2                                bool           `db:"metric_2"`
	Imperial                                bool           `db:"imperial"`
	Metric                                  bool           `db:"metric"`
	Universal                               bool           `db:"universal"`
	Imperial_2                              bool           `db:"imperial_2"`
}

func (q *Queries) GetAllValidMeasurementConversionsToMeasurementUnit(ctx context.Context, db DBTX, id string) ([]*GetAllValidMeasurementConversionsToMeasurementUnitRow, error) {
	rows, err := db.QueryContext(ctx, GetAllValidMeasurementConversionsToMeasurementUnit, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetAllValidMeasurementConversionsToMeasurementUnitRow{}
	for rows.Next() {
		var i GetAllValidMeasurementConversionsToMeasurementUnitRow
		if err := rows.Scan(
			&i.ID,
			&i.ID_2,
			&i.Name,
			&i.Description,
			&i.Volumetric,
			&i.IconPath,
			&i.Universal,
			&i.Metric,
			&i.Imperial,
			&i.Slug,
			&i.PluralName,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
			&i.ID_3,
			&i.Name_2,
			&i.Description_2,
			&i.Volumetric_2,
			&i.IconPath_2,
			&i.Universal_2,
			&i.Metric_2,
			&i.Imperial_2,
			&i.Slug_2,
			&i.PluralName_2,
			&i.CreatedAt_2,
			&i.LastUpdatedAt_2,
			&i.ArchivedAt_2,
			&i.ID_4,
			&i.Name_3,
			&i.Description_3,
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
			&i.Volumetric_3,
			&i.IsLiquid,
			&i.IconPath_3,
			&i.AnimalDerived,
			&i.PluralName_3,
			&i.RestrictToPreparations,
			&i.MinimumIdealStorageTemperatureInCelsius,
			&i.MaximumIdealStorageTemperatureInCelsius,
			&i.StorageInstructions,
			&i.Slug_3,
			&i.ContainsAlcohol,
			&i.ShoppingSuggestions,
			&i.IsStarch,
			&i.IsProtein,
			&i.IsGrain,
			&i.IsFruit,
			&i.IsSalt,
			&i.IsFat,
			&i.IsAcid,
			&i.IsHeat,
			&i.CreatedAt_3,
			&i.LastUpdatedAt_3,
			&i.ArchivedAt_3,
			&i.Modifier,
			&i.Notes,
			&i.CreatedAt_4,
			&i.LastUpdatedAt_4,
			&i.ArchivedAt_4,
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
