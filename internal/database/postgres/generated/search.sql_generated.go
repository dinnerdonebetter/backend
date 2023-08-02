// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: search.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const SearchForServiceSettings = `-- name: SearchForServiceSettings :many

SELECT
	service_settings.id,
    service_settings.name,
    service_settings.type,
    service_settings.description,
    service_settings.default_value,
    service_settings.admins_only,
    service_settings.enumeration,
    service_settings.created_at,
    service_settings.last_updated_at,
    service_settings.archived_at
FROM service_settings
WHERE service_settings.archived_at IS NULL
	AND service_settings.name ILIKE $1
LIMIT 50
`

type SearchForServiceSettingsRow struct {
	CreatedAt     time.Time      `db:"created_at"`
	LastUpdatedAt sql.NullTime   `db:"last_updated_at"`
	ArchivedAt    sql.NullTime   `db:"archived_at"`
	ID            string         `db:"id"`
	Name          string         `db:"name"`
	Type          SettingType    `db:"type"`
	Description   string         `db:"description"`
	Enumeration   string         `db:"enumeration"`
	DefaultValue  sql.NullString `db:"default_value"`
	AdminsOnly    bool           `db:"admins_only"`
}

func (q *Queries) SearchForServiceSettings(ctx context.Context, db DBTX, name string) ([]*SearchForServiceSettingsRow, error) {
	rows, err := db.QueryContext(ctx, SearchForServiceSettings, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*SearchForServiceSettingsRow{}
	for rows.Next() {
		var i SearchForServiceSettingsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Type,
			&i.Description,
			&i.DefaultValue,
			&i.AdminsOnly,
			&i.Enumeration,
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

const SearchForValidIngredientGroups = `-- name: SearchForValidIngredientGroups :many

SELECT
	valid_ingredient_groups.id,
	valid_ingredient_groups.name,
	valid_ingredient_groups.description,
	valid_ingredient_groups.slug,
	valid_ingredient_groups.created_at,
	valid_ingredient_groups.last_updated_at,
	valid_ingredient_groups.archived_at,
	valid_ingredient_group_members.id,
    valid_ingredient_group_members.belongs_to_group,
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
    valid_ingredient_group_members.created_at,
    valid_ingredient_group_members.archived_at
FROM valid_ingredient_groups
 JOIN valid_ingredient_group_members ON valid_ingredient_group_members.belongs_to_group=valid_ingredient_groups.id
  JOIN valid_ingredients ON valid_ingredients.id = valid_ingredient_group_members.valid_ingredient
WHERE valid_ingredient_groups.name ILIKE $1
AND valid_ingredient_groups.archived_at IS NULL
AND valid_ingredient_group_members.archived_at IS NULL
LIMIT 50
`

type SearchForValidIngredientGroupsRow struct {
	CreatedAt_2                             time.Time      `db:"created_at_2"`
	CreatedAt                               time.Time      `db:"created_at"`
	CreatedAt_3                             time.Time      `db:"created_at_3"`
	ArchivedAt_3                            sql.NullTime   `db:"archived_at_3"`
	LastUpdatedAt_2                         sql.NullTime   `db:"last_updated_at_2"`
	ArchivedAt_2                            sql.NullTime   `db:"archived_at_2"`
	LastUpdatedAt                           sql.NullTime   `db:"last_updated_at"`
	ArchivedAt                              sql.NullTime   `db:"archived_at"`
	Warning                                 string         `db:"warning"`
	Name                                    string         `db:"name"`
	Name_2                                  string         `db:"name_2"`
	Description_2                           string         `db:"description_2"`
	IconPath                                string         `db:"icon_path"`
	BelongsToGroup                          string         `db:"belongs_to_group"`
	ID_2                                    string         `db:"id_2"`
	Slug                                    string         `db:"slug"`
	Description                             string         `db:"description"`
	ID_3                                    string         `db:"id_3"`
	ShoppingSuggestions                     string         `db:"shopping_suggestions"`
	Slug_2                                  string         `db:"slug_2"`
	StorageInstructions                     string         `db:"storage_instructions"`
	ID                                      string         `db:"id"`
	PluralName                              string         `db:"plural_name"`
	MinimumIdealStorageTemperatureInCelsius sql.NullString `db:"minimum_ideal_storage_temperature_in_celsius"`
	MaximumIdealStorageTemperatureInCelsius sql.NullString `db:"maximum_ideal_storage_temperature_in_celsius"`
	IsLiquid                                sql.NullBool   `db:"is_liquid"`
	Volumetric                              bool           `db:"volumetric"`
	AnimalDerived                           bool           `db:"animal_derived"`
	AnimalFlesh                             bool           `db:"animal_flesh"`
	RestrictToPreparations                  bool           `db:"restrict_to_preparations"`
	ContainsGluten                          bool           `db:"contains_gluten"`
	ContainsFish                            bool           `db:"contains_fish"`
	ContainsSesame                          bool           `db:"contains_sesame"`
	ContainsShellfish                       bool           `db:"contains_shellfish"`
	ContainsAlcohol                         bool           `db:"contains_alcohol"`
	ContainsWheat                           bool           `db:"contains_wheat"`
	IsStarch                                bool           `db:"is_starch"`
	IsProtein                               bool           `db:"is_protein"`
	IsGrain                                 bool           `db:"is_grain"`
	IsFruit                                 bool           `db:"is_fruit"`
	IsSalt                                  bool           `db:"is_salt"`
	IsFat                                   bool           `db:"is_fat"`
	IsAcid                                  bool           `db:"is_acid"`
	IsHeat                                  bool           `db:"is_heat"`
	ContainsSoy                             bool           `db:"contains_soy"`
	ContainsTreeNut                         bool           `db:"contains_tree_nut"`
	ContainsPeanut                          bool           `db:"contains_peanut"`
	ContainsDairy                           bool           `db:"contains_dairy"`
	ContainsEgg                             bool           `db:"contains_egg"`
}

func (q *Queries) SearchForValidIngredientGroups(ctx context.Context, db DBTX, name string) ([]*SearchForValidIngredientGroupsRow, error) {
	rows, err := db.QueryContext(ctx, SearchForValidIngredientGroups, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*SearchForValidIngredientGroupsRow{}
	for rows.Next() {
		var i SearchForValidIngredientGroupsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Slug,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
			&i.ID_2,
			&i.BelongsToGroup,
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
			&i.IconPath,
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

const SearchForValidIngredientStates = `-- name: SearchForValidIngredientStates :many

SELECT
	valid_ingredient_states.id,
	valid_ingredient_states.name,
	valid_ingredient_states.description,
	valid_ingredient_states.icon_path,
	valid_ingredient_states.slug,
	valid_ingredient_states.past_tense,
	valid_ingredient_states.attribute_type,
	valid_ingredient_states.created_at,
	valid_ingredient_states.last_updated_at,
	valid_ingredient_states.archived_at
FROM valid_ingredient_states
WHERE valid_ingredient_states.archived_at IS NULL
	AND valid_ingredient_states.name ILIKE $1
LIMIT 50
`

type SearchForValidIngredientStatesRow struct {
	ID            string                  `db:"id"`
	Name          string                  `db:"name"`
	Description   string                  `db:"description"`
	IconPath      string                  `db:"icon_path"`
	Slug          string                  `db:"slug"`
	PastTense     string                  `db:"past_tense"`
	AttributeType IngredientAttributeType `db:"attribute_type"`
	CreatedAt     time.Time               `db:"created_at"`
	LastUpdatedAt sql.NullTime            `db:"last_updated_at"`
	ArchivedAt    sql.NullTime            `db:"archived_at"`
}

func (q *Queries) SearchForValidIngredientStates(ctx context.Context, db DBTX, name string) ([]*SearchForValidIngredientStatesRow, error) {
	rows, err := db.QueryContext(ctx, SearchForValidIngredientStates, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*SearchForValidIngredientStatesRow{}
	for rows.Next() {
		var i SearchForValidIngredientStatesRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.IconPath,
			&i.Slug,
			&i.PastTense,
			&i.AttributeType,
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

const SearchForValidIngredients = `-- name: SearchForValidIngredients :many

SELECT
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
	valid_ingredients.archived_at
FROM valid_ingredients
WHERE valid_ingredients.name ILIKE $1
	AND valid_ingredients.archived_at IS NULL
LIMIT 50
`

type SearchForValidIngredientsRow struct {
	CreatedAt                               time.Time      `db:"created_at"`
	ArchivedAt                              sql.NullTime   `db:"archived_at"`
	LastUpdatedAt                           sql.NullTime   `db:"last_updated_at"`
	Warning                                 string         `db:"warning"`
	Description                             string         `db:"description"`
	Name                                    string         `db:"name"`
	ShoppingSuggestions                     string         `db:"shopping_suggestions"`
	Slug                                    string         `db:"slug"`
	StorageInstructions                     string         `db:"storage_instructions"`
	PluralName                              string         `db:"plural_name"`
	ID                                      string         `db:"id"`
	IconPath                                string         `db:"icon_path"`
	MaximumIdealStorageTemperatureInCelsius sql.NullString `db:"maximum_ideal_storage_temperature_in_celsius"`
	MinimumIdealStorageTemperatureInCelsius sql.NullString `db:"minimum_ideal_storage_temperature_in_celsius"`
	IsLiquid                                sql.NullBool   `db:"is_liquid"`
	AnimalDerived                           bool           `db:"animal_derived"`
	ContainsTreeNut                         bool           `db:"contains_tree_nut"`
	AnimalFlesh                             bool           `db:"animal_flesh"`
	ContainsGluten                          bool           `db:"contains_gluten"`
	ContainsFish                            bool           `db:"contains_fish"`
	RestrictToPreparations                  bool           `db:"restrict_to_preparations"`
	ContainsSesame                          bool           `db:"contains_sesame"`
	ContainsShellfish                       bool           `db:"contains_shellfish"`
	ContainsWheat                           bool           `db:"contains_wheat"`
	ContainsSoy                             bool           `db:"contains_soy"`
	ContainsAlcohol                         bool           `db:"contains_alcohol"`
	Volumetric                              bool           `db:"volumetric"`
	IsStarch                                bool           `db:"is_starch"`
	IsProtein                               bool           `db:"is_protein"`
	IsGrain                                 bool           `db:"is_grain"`
	IsFruit                                 bool           `db:"is_fruit"`
	IsSalt                                  bool           `db:"is_salt"`
	IsFat                                   bool           `db:"is_fat"`
	IsAcid                                  bool           `db:"is_acid"`
	IsHeat                                  bool           `db:"is_heat"`
	ContainsPeanut                          bool           `db:"contains_peanut"`
	ContainsDairy                           bool           `db:"contains_dairy"`
	ContainsEgg                             bool           `db:"contains_egg"`
}

func (q *Queries) SearchForValidIngredients(ctx context.Context, db DBTX, name string) ([]*SearchForValidIngredientsRow, error) {
	rows, err := db.QueryContext(ctx, SearchForValidIngredients, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*SearchForValidIngredientsRow{}
	for rows.Next() {
		var i SearchForValidIngredientsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
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
			&i.IconPath,
			&i.AnimalDerived,
			&i.PluralName,
			&i.RestrictToPreparations,
			&i.MinimumIdealStorageTemperatureInCelsius,
			&i.MaximumIdealStorageTemperatureInCelsius,
			&i.StorageInstructions,
			&i.Slug,
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

const SearchForValidInstruments = `-- name: SearchForValidInstruments :many

SELECT
	valid_instruments.id,
	valid_instruments.name,
	valid_instruments.plural_name,
	valid_instruments.description,
	valid_instruments.icon_path,
	valid_instruments.usable_for_storage,
	valid_instruments.display_in_summary_lists,
	valid_instruments.include_in_generated_instructions,
	valid_instruments.slug,
	valid_instruments.created_at,
	valid_instruments.last_updated_at,
	valid_instruments.archived_at
FROM valid_instruments
WHERE valid_instruments.archived_at IS NULL
	AND valid_instruments.name ILIKE $1 LIMIT 50
`

type SearchForValidInstrumentsRow struct {
	CreatedAt                      time.Time    `db:"created_at"`
	LastUpdatedAt                  sql.NullTime `db:"last_updated_at"`
	ArchivedAt                     sql.NullTime `db:"archived_at"`
	ID                             string       `db:"id"`
	Name                           string       `db:"name"`
	PluralName                     string       `db:"plural_name"`
	Description                    string       `db:"description"`
	IconPath                       string       `db:"icon_path"`
	Slug                           string       `db:"slug"`
	UsableForStorage               bool         `db:"usable_for_storage"`
	DisplayInSummaryLists          bool         `db:"display_in_summary_lists"`
	IncludeInGeneratedInstructions bool         `db:"include_in_generated_instructions"`
}

func (q *Queries) SearchForValidInstruments(ctx context.Context, db DBTX, name string) ([]*SearchForValidInstrumentsRow, error) {
	rows, err := db.QueryContext(ctx, SearchForValidInstruments, name)
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
			&i.PluralName,
			&i.Description,
			&i.IconPath,
			&i.UsableForStorage,
			&i.DisplayInSummaryLists,
			&i.IncludeInGeneratedInstructions,
			&i.Slug,
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

const SearchForValidMeasurementUnits = `-- name: SearchForValidMeasurementUnits :many

SELECT
	valid_measurement_units.id,
	valid_measurement_units.name,
	valid_measurement_units.description,
	valid_measurement_units.volumetric,
	valid_measurement_units.icon_path,
	valid_measurement_units.universal,
	valid_measurement_units.metric,
	valid_measurement_units.imperial,
	valid_measurement_units.slug,
	valid_measurement_units.plural_name,
	valid_measurement_units.created_at,
	valid_measurement_units.last_updated_at,
	valid_measurement_units.archived_at
FROM valid_measurement_units
WHERE (valid_measurement_units.name ILIKE $1 OR valid_measurement_units.universal is TRUE)
AND valid_measurement_units.archived_at IS NULL
LIMIT 50
`

type SearchForValidMeasurementUnitsRow struct {
	CreatedAt     time.Time    `db:"created_at"`
	ArchivedAt    sql.NullTime `db:"archived_at"`
	LastUpdatedAt sql.NullTime `db:"last_updated_at"`
	PluralName    string       `db:"plural_name"`
	Name          string       `db:"name"`
	Description   string       `db:"description"`
	ID            string       `db:"id"`
	IconPath      string       `db:"icon_path"`
	Slug          string       `db:"slug"`
	Volumetric    sql.NullBool `db:"volumetric"`
	Imperial      bool         `db:"imperial"`
	Metric        bool         `db:"metric"`
	Universal     bool         `db:"universal"`
}

func (q *Queries) SearchForValidMeasurementUnits(ctx context.Context, db DBTX, name string) ([]*SearchForValidMeasurementUnitsRow, error) {
	rows, err := db.QueryContext(ctx, SearchForValidMeasurementUnits, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*SearchForValidMeasurementUnitsRow{}
	for rows.Next() {
		var i SearchForValidMeasurementUnitsRow
		if err := rows.Scan(
			&i.ID,
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

const SearchForValidPreparations = `-- name: SearchForValidPreparations :many

SELECT
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
	valid_preparations.archived_at
FROM valid_preparations
WHERE valid_preparations.archived_at IS NULL
	AND valid_preparations.name ILIKE $1
LIMIT 50
`

type SearchForValidPreparationsRow struct {
	CreatedAt                   time.Time     `db:"created_at"`
	LastUpdatedAt               sql.NullTime  `db:"last_updated_at"`
	ArchivedAt                  sql.NullTime  `db:"archived_at"`
	Name                        string        `db:"name"`
	Description                 string        `db:"description"`
	IconPath                    string        `db:"icon_path"`
	ID                          string        `db:"id"`
	Slug                        string        `db:"slug"`
	PastTense                   string        `db:"past_tense"`
	MaximumInstrumentCount      sql.NullInt32 `db:"maximum_instrument_count"`
	MaximumIngredientCount      sql.NullInt32 `db:"maximum_ingredient_count"`
	MaximumVesselCount          sql.NullInt32 `db:"maximum_vessel_count"`
	MinimumVesselCount          int32         `db:"minimum_vessel_count"`
	MinimumInstrumentCount      int32         `db:"minimum_instrument_count"`
	MinimumIngredientCount      int32         `db:"minimum_ingredient_count"`
	RestrictToIngredients       bool          `db:"restrict_to_ingredients"`
	OnlyForVessels              bool          `db:"only_for_vessels"`
	ConsumesVessel              bool          `db:"consumes_vessel"`
	ConditionExpressionRequired bool          `db:"condition_expression_required"`
	TimeEstimateRequired        bool          `db:"time_estimate_required"`
	TemperatureRequired         bool          `db:"temperature_required"`
	YieldsNothing               bool          `db:"yields_nothing"`
}

func (q *Queries) SearchForValidPreparations(ctx context.Context, db DBTX, name string) ([]*SearchForValidPreparationsRow, error) {
	rows, err := db.QueryContext(ctx, SearchForValidPreparations, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*SearchForValidPreparationsRow{}
	for rows.Next() {
		var i SearchForValidPreparationsRow
		if err := rows.Scan(
			&i.ID,
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

const SearchForValidVessels = `-- name: SearchForValidVessels :many

SELECT
	valid_vessels.id,
    valid_vessels.name,
    valid_vessels.plural_name,
    valid_vessels.description,
    valid_vessels.icon_path,
    valid_vessels.usable_for_storage,
    valid_vessels.slug,
    valid_vessels.display_in_summary_lists,
    valid_vessels.include_in_generated_instructions,
    valid_vessels.capacity,
	valid_measurement_units.id,
	valid_measurement_units.name,
	valid_measurement_units.description,
	valid_measurement_units.volumetric,
	valid_measurement_units.icon_path,
	valid_measurement_units.universal,
	valid_measurement_units.metric,
	valid_measurement_units.imperial,
	valid_measurement_units.slug,
	valid_measurement_units.plural_name,
	valid_measurement_units.created_at,
	valid_measurement_units.last_updated_at,
	valid_measurement_units.archived_at,
    valid_vessels.width_in_millimeters,
    valid_vessels.length_in_millimeters,
    valid_vessels.height_in_millimeters,
    valid_vessels.shape,
    valid_vessels.created_at,
    valid_vessels.last_updated_at,
    valid_vessels.archived_at
FROM valid_vessels
	 JOIN valid_measurement_units ON valid_vessels.capacity_unit=valid_measurement_units.id
WHERE valid_vessels.archived_at IS NULL
	AND valid_measurement_units.archived_at IS NULL
	AND valid_vessels.name ILIKE $1 LIMIT 50
`

type SearchForValidVesselsRow struct {
	CreatedAt                      time.Time      `db:"created_at"`
	CreatedAt_2                    time.Time      `db:"created_at_2"`
	ArchivedAt_2                   sql.NullTime   `db:"archived_at_2"`
	LastUpdatedAt_2                sql.NullTime   `db:"last_updated_at_2"`
	ArchivedAt                     sql.NullTime   `db:"archived_at"`
	LastUpdatedAt                  sql.NullTime   `db:"last_updated_at"`
	IconPath_2                     string         `db:"icon_path_2"`
	IconPath                       string         `db:"icon_path"`
	Name                           string         `db:"name"`
	Capacity                       string         `db:"capacity"`
	ID_2                           string         `db:"id_2"`
	Name_2                         string         `db:"name_2"`
	Description_2                  string         `db:"description_2"`
	PluralName                     string         `db:"plural_name"`
	ID                             string         `db:"id"`
	Description                    string         `db:"description"`
	Shape                          VesselShape    `db:"shape"`
	Slug                           string         `db:"slug"`
	Slug_2                         string         `db:"slug_2"`
	PluralName_2                   string         `db:"plural_name_2"`
	WidthInMillimeters             sql.NullString `db:"width_in_millimeters"`
	LengthInMillimeters            sql.NullString `db:"length_in_millimeters"`
	HeightInMillimeters            sql.NullString `db:"height_in_millimeters"`
	Volumetric                     sql.NullBool   `db:"volumetric"`
	Imperial                       bool           `db:"imperial"`
	UsableForStorage               bool           `db:"usable_for_storage"`
	DisplayInSummaryLists          bool           `db:"display_in_summary_lists"`
	Metric                         bool           `db:"metric"`
	Universal                      bool           `db:"universal"`
	IncludeInGeneratedInstructions bool           `db:"include_in_generated_instructions"`
}

func (q *Queries) SearchForValidVessels(ctx context.Context, db DBTX, name string) ([]*SearchForValidVesselsRow, error) {
	rows, err := db.QueryContext(ctx, SearchForValidVessels, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*SearchForValidVesselsRow{}
	for rows.Next() {
		var i SearchForValidVesselsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.PluralName,
			&i.Description,
			&i.IconPath,
			&i.UsableForStorage,
			&i.Slug,
			&i.DisplayInSummaryLists,
			&i.IncludeInGeneratedInstructions,
			&i.Capacity,
			&i.ID_2,
			&i.Name_2,
			&i.Description_2,
			&i.Volumetric,
			&i.IconPath_2,
			&i.Universal,
			&i.Metric,
			&i.Imperial,
			&i.Slug_2,
			&i.PluralName_2,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
			&i.WidthInMillimeters,
			&i.LengthInMillimeters,
			&i.HeightInMillimeters,
			&i.Shape,
			&i.CreatedAt_2,
			&i.LastUpdatedAt_2,
			&i.ArchivedAt_2,
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
