// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: get_with_ids.sql

package generated

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

const getValidInstrumentWithIDs = `-- name: GetValidInstrumentWithIDs :many

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
    AND valid_instruments.id = ANY($1::text[])
`

type GetValidInstrumentWithIDsRow struct {
	ID                             string
	Name                           string
	PluralName                     string
	Description                    string
	IconPath                       string
	UsableForStorage               bool
	DisplayInSummaryLists          bool
	IncludeInGeneratedInstructions bool
	Slug                           string
	CreatedAt                      time.Time
	LastUpdatedAt                  sql.NullTime
	ArchivedAt                     sql.NullTime
}

func (q *Queries) GetValidInstrumentWithIDs(ctx context.Context, db DBTX, ids []string) ([]*GetValidInstrumentWithIDsRow, error) {
	rows, err := db.QueryContext(ctx, getValidInstrumentWithIDs, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetValidInstrumentWithIDsRow{}
	for rows.Next() {
		var i GetValidInstrumentWithIDsRow
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

const getValidMeasurementUnitsWithIDs = `-- name: GetValidMeasurementUnitsWithIDs :many

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
WHERE valid_measurement_units.archived_at IS NULL
	AND valid_measurement_units.id = ANY($1::text[])
`

type GetValidMeasurementUnitsWithIDsRow struct {
	ID            string
	Name          string
	Description   string
	Volumetric    sql.NullBool
	IconPath      string
	Universal     bool
	Metric        bool
	Imperial      bool
	Slug          string
	PluralName    string
	CreatedAt     time.Time
	LastUpdatedAt sql.NullTime
	ArchivedAt    sql.NullTime
}

func (q *Queries) GetValidMeasurementUnitsWithIDs(ctx context.Context, db DBTX, ids []string) ([]*GetValidMeasurementUnitsWithIDsRow, error) {
	rows, err := db.QueryContext(ctx, getValidMeasurementUnitsWithIDs, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetValidMeasurementUnitsWithIDsRow{}
	for rows.Next() {
		var i GetValidMeasurementUnitsWithIDsRow
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

const getValidPreparationsWithIDs = `-- name: GetValidPreparationsWithIDs :many

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
	AND valid_preparations.id = ANY($1::text[])
`

type GetValidPreparationsWithIDsRow struct {
	ID                          string
	Name                        string
	Description                 string
	IconPath                    string
	YieldsNothing               bool
	RestrictToIngredients       bool
	MinimumIngredientCount      int32
	MaximumIngredientCount      sql.NullInt32
	MinimumInstrumentCount      int32
	MaximumInstrumentCount      sql.NullInt32
	TemperatureRequired         bool
	TimeEstimateRequired        bool
	ConditionExpressionRequired bool
	ConsumesVessel              bool
	OnlyForVessels              bool
	MinimumVesselCount          int32
	MaximumVesselCount          sql.NullInt32
	Slug                        string
	PastTense                   string
	CreatedAt                   time.Time
	LastUpdatedAt               sql.NullTime
	ArchivedAt                  sql.NullTime
}

func (q *Queries) GetValidPreparationsWithIDs(ctx context.Context, db DBTX, dollar_1 []string) ([]*GetValidPreparationsWithIDsRow, error) {
	rows, err := db.QueryContext(ctx, getValidPreparationsWithIDs, pq.Array(dollar_1))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetValidPreparationsWithIDsRow{}
	for rows.Next() {
		var i GetValidPreparationsWithIDsRow
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

const getValidVesselsWithIDs = `-- name: GetValidVesselsWithIDs :many

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
    valid_vessels.capacity::float,
    valid_measurement_units.id as valid_measurement_unit_id,
    valid_measurement_units.name as valid_measurement_unit_name,
    valid_measurement_units.description as valid_measurement_unit_description,
    valid_measurement_units.volumetric as valid_measurement_unit_volumetric,
    valid_measurement_units.icon_path as valid_measurement_unit_icon_path,
    valid_measurement_units.universal as valid_measurement_unit_universal,
    valid_measurement_units.metric as valid_measurement_unit_metric,
    valid_measurement_units.imperial as valid_measurement_unit_imperial,
    valid_measurement_units.slug as valid_measurement_unit_slug,
    valid_measurement_units.plural_name as valid_measurement_unit_plural_name,
    valid_measurement_units.created_at as valid_measurement_unit_created_at,
    valid_measurement_units.last_updated_at as valid_measurement_unit_last_updated_at,
    valid_measurement_units.archived_at as valid_measurement_unit_archived_at,
    valid_vessels.width_in_millimeters::float,
    valid_vessels.length_in_millimeters::float,
    valid_vessels.height_in_millimeters::float,
    valid_vessels.shape,
    valid_vessels.created_at,
    valid_vessels.last_updated_at,
    valid_vessels.archived_at
FROM valid_vessels
    JOIN valid_measurement_units ON valid_vessels.capacity_unit=valid_measurement_units.id
WHERE valid_vessels.archived_at IS NULL
  AND valid_measurement_units.archived_at IS NULL
  AND valid_vessels.id = ANY($1::text[])
`

type GetValidVesselsWithIDsRow struct {
	ID                                string
	Name                              string
	PluralName                        string
	Description                       string
	IconPath                          string
	UsableForStorage                  bool
	Slug                              string
	DisplayInSummaryLists             bool
	IncludeInGeneratedInstructions    bool
	ValidVesselsCapacity              float64
	ValidMeasurementUnitID            string
	ValidMeasurementUnitName          string
	ValidMeasurementUnitDescription   string
	ValidMeasurementUnitVolumetric    sql.NullBool
	ValidMeasurementUnitIconPath      string
	ValidMeasurementUnitUniversal     bool
	ValidMeasurementUnitMetric        bool
	ValidMeasurementUnitImperial      bool
	ValidMeasurementUnitSlug          string
	ValidMeasurementUnitPluralName    string
	ValidMeasurementUnitCreatedAt     time.Time
	ValidMeasurementUnitLastUpdatedAt sql.NullTime
	ValidMeasurementUnitArchivedAt    sql.NullTime
	ValidVesselsWidthInMillimeters    float64
	ValidVesselsLengthInMillimeters   float64
	ValidVesselsHeightInMillimeters   float64
	Shape                             VesselShape
	CreatedAt                         time.Time
	LastUpdatedAt                     sql.NullTime
	ArchivedAt                        sql.NullTime
}

func (q *Queries) GetValidVesselsWithIDs(ctx context.Context, db DBTX, ids []string) ([]*GetValidVesselsWithIDsRow, error) {
	rows, err := db.QueryContext(ctx, getValidVesselsWithIDs, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetValidVesselsWithIDsRow{}
	for rows.Next() {
		var i GetValidVesselsWithIDsRow
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
			&i.ValidVesselsCapacity,
			&i.ValidMeasurementUnitID,
			&i.ValidMeasurementUnitName,
			&i.ValidMeasurementUnitDescription,
			&i.ValidMeasurementUnitVolumetric,
			&i.ValidMeasurementUnitIconPath,
			&i.ValidMeasurementUnitUniversal,
			&i.ValidMeasurementUnitMetric,
			&i.ValidMeasurementUnitImperial,
			&i.ValidMeasurementUnitSlug,
			&i.ValidMeasurementUnitPluralName,
			&i.ValidMeasurementUnitCreatedAt,
			&i.ValidMeasurementUnitLastUpdatedAt,
			&i.ValidMeasurementUnitArchivedAt,
			&i.ValidVesselsWidthInMillimeters,
			&i.ValidVesselsLengthInMillimeters,
			&i.ValidVesselsHeightInMillimeters,
			&i.Shape,
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
