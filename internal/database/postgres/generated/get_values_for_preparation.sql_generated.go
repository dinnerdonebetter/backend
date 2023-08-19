// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: get_values_for_preparation.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const getValidIngredientPreparationsForPreparation = `-- name: GetValidIngredientPreparationsForPreparation :many

SELECT
  valid_ingredient_preparations.id,
  valid_ingredient_preparations.notes,
  valid_preparations.id,
  valid_preparations.name,
  valid_preparations.description,
  valid_preparations.icon_path,
  valid_preparations.yields_nothing,
  valid_preparations.restrict_to_ingredients,
  valid_preparations.minimum_ingredient_count,
  valid_preparations.maximum_ingredient_count,
  valid_preparations.minimum_instrument_count,
  valid_preparations.maximum_instrument_count,
  valid_preparations.temperature_required,
  valid_preparations.time_estimate_required,
  valid_preparations.condition_expression_required,
  valid_preparations.consumes_vessel,
  valid_preparations.only_for_vessels,
  valid_preparations.minimum_vessel_count,
  valid_preparations.maximum_vessel_count,
  valid_preparations.slug,
  valid_preparations.past_tense,
  valid_preparations.created_at,
  valid_preparations.last_updated_at,
  valid_preparations.archived_at,
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
  valid_ingredient_preparations.created_at,
  valid_ingredient_preparations.last_updated_at,
  valid_ingredient_preparations.archived_at
FROM
  valid_ingredient_preparations
  JOIN valid_ingredients ON valid_ingredient_preparations.valid_ingredient_id = valid_ingredients.id
  JOIN valid_preparations ON valid_ingredient_preparations.valid_preparation_id = valid_preparations.id
WHERE
  valid_ingredient_preparations.archived_at IS NULL
  AND valid_ingredient_preparations.valid_preparation_id IN ($1)
LIMIT
  20
`

type GetValidIngredientPreparationsForPreparationRow struct {
	ID                                      string
	Notes                                   string
	ID_2                                    string
	Name                                    string
	Description                             string
	IconPath                                string
	YieldsNothing                           bool
	RestrictToIngredients                   bool
	MinimumIngredientCount                  int32
	MaximumIngredientCount                  sql.NullInt32
	MinimumInstrumentCount                  int32
	MaximumInstrumentCount                  sql.NullInt32
	TemperatureRequired                     bool
	TimeEstimateRequired                    bool
	ConditionExpressionRequired             bool
	ConsumesVessel                          bool
	OnlyForVessels                          bool
	MinimumVesselCount                      int32
	MaximumVesselCount                      sql.NullInt32
	Slug                                    string
	PastTense                               string
	CreatedAt                               time.Time
	LastUpdatedAt                           sql.NullTime
	ArchivedAt                              sql.NullTime
	ID_3                                    string
	Name_2                                  string
	Description_2                           string
	Warning                                 string
	ContainsEgg                             bool
	ContainsDairy                           bool
	ContainsPeanut                          bool
	ContainsTreeNut                         bool
	ContainsSoy                             bool
	ContainsWheat                           bool
	ContainsShellfish                       bool
	ContainsSesame                          bool
	ContainsFish                            bool
	ContainsGluten                          bool
	AnimalFlesh                             bool
	Volumetric                              bool
	IsLiquid                                sql.NullBool
	IconPath_2                              string
	AnimalDerived                           bool
	PluralName                              string
	RestrictToPreparations                  bool
	MinimumIdealStorageTemperatureInCelsius sql.NullString
	MaximumIdealStorageTemperatureInCelsius sql.NullString
	StorageInstructions                     string
	Slug_2                                  string
	ContainsAlcohol                         bool
	ShoppingSuggestions                     string
	IsStarch                                bool
	IsProtein                               bool
	IsGrain                                 bool
	IsFruit                                 bool
	IsSalt                                  bool
	IsFat                                   bool
	IsAcid                                  bool
	IsHeat                                  bool
	CreatedAt_2                             time.Time
	LastUpdatedAt_2                         sql.NullTime
	ArchivedAt_2                            sql.NullTime
	CreatedAt_3                             time.Time
	LastUpdatedAt_3                         sql.NullTime
	ArchivedAt_3                            sql.NullTime
}

func (q *Queries) GetValidIngredientPreparationsForPreparation(ctx context.Context, db DBTX, validPreparationID string) ([]*GetValidIngredientPreparationsForPreparationRow, error) {
	rows, err := db.QueryContext(ctx, getValidIngredientPreparationsForPreparation, validPreparationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetValidIngredientPreparationsForPreparationRow{}
	for rows.Next() {
		var i GetValidIngredientPreparationsForPreparationRow
		if err := rows.Scan(
			&i.ID,
			&i.Notes,
			&i.ID_2,
			&i.Name,
			&i.Description,
			&i.IconPath,
			&i.YieldsNothing,
			&i.RestrictToIngredients,
			&i.MinimumIngredientCount,
			&i.MaximumIngredientCount,
			&i.MinimumInstrumentCount,
			&i.MaximumInstrumentCount,
			&i.TemperatureRequired,
			&i.TimeEstimateRequired,
			&i.ConditionExpressionRequired,
			&i.ConsumesVessel,
			&i.OnlyForVessels,
			&i.MinimumVesselCount,
			&i.MaximumVesselCount,
			&i.Slug,
			&i.PastTense,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
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
			&i.Volumetric,
			&i.IsLiquid,
			&i.IconPath_2,
			&i.AnimalDerived,
			&i.PluralName,
			&i.RestrictToPreparations,
			&i.MinimumIdealStorageTemperatureInCelsius,
			&i.MaximumIdealStorageTemperatureInCelsius,
			&i.StorageInstructions,
			&i.Slug_2,
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
			&i.CreatedAt_2,
			&i.LastUpdatedAt_2,
			&i.ArchivedAt_2,
			&i.CreatedAt_3,
			&i.LastUpdatedAt_3,
			&i.ArchivedAt_3,
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
