// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: recipe_step_vessels.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const archiveRecipeStepVessel = `-- name: ArchiveRecipeStepVessel :exec

UPDATE recipe_step_vessels SET archived_at = NOW() WHERE archived_at IS NULL AND belongs_to_recipe_step = $1 AND id = $2
`

type ArchiveRecipeStepVesselParams struct {
	BelongsToRecipeStep string
	ID                  string
}

func (q *Queries) ArchiveRecipeStepVessel(ctx context.Context, db DBTX, arg *ArchiveRecipeStepVesselParams) error {
	_, err := db.ExecContext(ctx, archiveRecipeStepVessel, arg.BelongsToRecipeStep, arg.ID)
	return err
}

const checkRecipeStepVesselExistence = `-- name: CheckRecipeStepVesselExistence :one

SELECT EXISTS (
    SELECT recipe_step_vessels.id
    FROM recipe_step_vessels
        JOIN recipe_steps ON recipe_step_vessels.belongs_to_recipe_step=recipe_steps.id
        JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
    WHERE recipe_step_vessels.archived_at IS NULL
        AND recipe_step_vessels.belongs_to_recipe_step = $1
        AND recipe_step_vessels.id = $2
        AND recipe_steps.archived_at IS NULL
        AND recipe_steps.belongs_to_recipe = $3
        AND recipe_steps.id = $1
        AND recipes.archived_at IS NULL
        AND recipes.id = $3
)
`

type CheckRecipeStepVesselExistenceParams struct {
	RecipeStepID       string
	RecipeStepVesselID string
	RecipeID           string
}

func (q *Queries) CheckRecipeStepVesselExistence(ctx context.Context, db DBTX, arg *CheckRecipeStepVesselExistenceParams) (bool, error) {
	row := db.QueryRowContext(ctx, checkRecipeStepVesselExistence, arg.RecipeStepID, arg.RecipeStepVesselID, arg.RecipeID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createRecipeStepVessel = `-- name: CreateRecipeStepVessel :exec

INSERT INTO recipe_step_vessels
(id,"name",notes,belongs_to_recipe_step,recipe_step_product_id,valid_vessel_id,vessel_predicate,minimum_quantity,maximum_quantity,unavailable_after_step)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
`

type CreateRecipeStepVesselParams struct {
	ID                   string
	Name                 string
	Notes                string
	BelongsToRecipeStep  string
	VesselPredicate      string
	RecipeStepProductID  sql.NullString
	ValidVesselID        sql.NullString
	MaximumQuantity      sql.NullInt32
	MinimumQuantity      int32
	UnavailableAfterStep bool
}

func (q *Queries) CreateRecipeStepVessel(ctx context.Context, db DBTX, arg *CreateRecipeStepVesselParams) error {
	_, err := db.ExecContext(ctx, createRecipeStepVessel,
		arg.ID,
		arg.Name,
		arg.Notes,
		arg.BelongsToRecipeStep,
		arg.RecipeStepProductID,
		arg.ValidVesselID,
		arg.VesselPredicate,
		arg.MinimumQuantity,
		arg.MaximumQuantity,
		arg.UnavailableAfterStep,
	)
	return err
}

const getRecipeStepVessel = `-- name: GetRecipeStepVessel :one

SELECT
    recipe_step_vessels.id,
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
    valid_vessels.archived_at,
    recipe_step_vessels.name,
    recipe_step_vessels.notes,
    recipe_step_vessels.belongs_to_recipe_step,
    recipe_step_vessels.recipe_step_product_id,
    recipe_step_vessels.vessel_predicate,
    recipe_step_vessels.minimum_quantity,
    recipe_step_vessels.maximum_quantity,
    recipe_step_vessels.unavailable_after_step,
    recipe_step_vessels.created_at,
    recipe_step_vessels.last_updated_at,
    recipe_step_vessels.archived_at
FROM recipe_step_vessels
	LEFT JOIN valid_vessels ON recipe_step_vessels.valid_vessel_id=valid_vessels.id
    LEFT JOIN valid_measurement_units ON valid_vessels.capacity_unit=valid_measurement_units.id
	JOIN recipe_steps ON recipe_step_vessels.belongs_to_recipe_step=recipe_steps.id
	JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
WHERE recipe_step_vessels.archived_at IS NULL
	AND recipe_step_vessels.belongs_to_recipe_step = $1
	AND recipe_step_vessels.id = $2
	AND recipe_steps.archived_at IS NULL
	AND recipe_steps.belongs_to_recipe = $3
	AND recipe_steps.id = $4
	AND recipes.archived_at IS NULL
	AND recipes.id = $5
`

type GetRecipeStepVesselParams struct {
	BelongsToRecipeStep string
	ID                  string
	BelongsToRecipe     string
	ID_2                string
	ID_3                string
}

type GetRecipeStepVesselRow struct {
	CreatedAt_3                    time.Time
	LastUpdatedAt_2                sql.NullTime
	LastUpdatedAt_3                sql.NullTime
	ArchivedAt_3                   sql.NullTime
	ArchivedAt_2                   sql.NullTime
	CreatedAt                      sql.NullTime
	CreatedAt_2                    sql.NullTime
	ArchivedAt                     sql.NullTime
	LastUpdatedAt                  sql.NullTime
	VesselPredicate                string
	BelongsToRecipeStep            string
	Notes                          string
	Name_3                         string
	ID                             string
	Capacity                       sql.NullString
	IconPath                       sql.NullString
	ID_2                           sql.NullString
	Name                           sql.NullString
	PluralName                     sql.NullString
	Slug_2                         sql.NullString
	PluralName_2                   sql.NullString
	Description                    sql.NullString
	Description_2                  sql.NullString
	Name_2                         sql.NullString
	WidthInMillimeters             sql.NullString
	LengthInMillimeters            sql.NullString
	HeightInMillimeters            sql.NullString
	Shape                          NullVesselShape
	ID_3                           sql.NullString
	RecipeStepProductID            sql.NullString
	IconPath_2                     sql.NullString
	Slug                           sql.NullString
	MaximumQuantity                sql.NullInt32
	MinimumQuantity                int32
	UsableForStorage               sql.NullBool
	DisplayInSummaryLists          sql.NullBool
	IncludeInGeneratedInstructions sql.NullBool
	Volumetric                     sql.NullBool
	Imperial                       sql.NullBool
	Metric                         sql.NullBool
	Universal                      sql.NullBool
	UnavailableAfterStep           bool
}

func (q *Queries) GetRecipeStepVessel(ctx context.Context, db DBTX, arg *GetRecipeStepVesselParams) (*GetRecipeStepVesselRow, error) {
	row := db.QueryRowContext(ctx, getRecipeStepVessel,
		arg.BelongsToRecipeStep,
		arg.ID,
		arg.BelongsToRecipe,
		arg.ID_2,
		arg.ID_3,
	)
	var i GetRecipeStepVesselRow
	err := row.Scan(
		&i.ID,
		&i.ID_2,
		&i.Name,
		&i.PluralName,
		&i.Description,
		&i.IconPath,
		&i.UsableForStorage,
		&i.Slug,
		&i.DisplayInSummaryLists,
		&i.IncludeInGeneratedInstructions,
		&i.Capacity,
		&i.ID_3,
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
		&i.Name_3,
		&i.Notes,
		&i.BelongsToRecipeStep,
		&i.RecipeStepProductID,
		&i.VesselPredicate,
		&i.MinimumQuantity,
		&i.MaximumQuantity,
		&i.UnavailableAfterStep,
		&i.CreatedAt_3,
		&i.LastUpdatedAt_3,
		&i.ArchivedAt_3,
	)
	return &i, err
}

const getRecipeStepVesselsForRecipe = `-- name: GetRecipeStepVesselsForRecipe :many

SELECT
	recipe_step_vessels.id,
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
    valid_vessels.archived_at,
    recipe_step_vessels.name,
    recipe_step_vessels.notes,
    recipe_step_vessels.belongs_to_recipe_step,
    recipe_step_vessels.recipe_step_product_id,
    recipe_step_vessels.vessel_predicate,
    recipe_step_vessels.minimum_quantity,
    recipe_step_vessels.maximum_quantity,
    recipe_step_vessels.unavailable_after_step,
    recipe_step_vessels.created_at,
    recipe_step_vessels.last_updated_at,
    recipe_step_vessels.archived_at
FROM recipe_step_vessels
	LEFT JOIN valid_vessels ON recipe_step_vessels.valid_vessel_id=valid_vessels.id
    LEFT JOIN valid_measurement_units ON valid_vessels.capacity_unit=valid_measurement_units.id
	JOIN recipe_steps ON recipe_step_vessels.belongs_to_recipe_step=recipe_steps.id
	JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
WHERE recipe_step_vessels.archived_at IS NULL
	AND recipe_steps.archived_at IS NULL
	AND recipe_steps.belongs_to_recipe = $1
	AND recipes.archived_at IS NULL
	AND recipes.id = $1
`

type GetRecipeStepVesselsForRecipeRow struct {
	CreatedAt_3                    time.Time
	LastUpdatedAt_2                sql.NullTime
	LastUpdatedAt_3                sql.NullTime
	ArchivedAt_3                   sql.NullTime
	ArchivedAt_2                   sql.NullTime
	CreatedAt                      sql.NullTime
	CreatedAt_2                    sql.NullTime
	ArchivedAt                     sql.NullTime
	LastUpdatedAt                  sql.NullTime
	VesselPredicate                string
	BelongsToRecipeStep            string
	Notes                          string
	Name_3                         string
	ID                             string
	Capacity                       sql.NullString
	IconPath                       sql.NullString
	ID_2                           sql.NullString
	Name                           sql.NullString
	PluralName                     sql.NullString
	Slug_2                         sql.NullString
	PluralName_2                   sql.NullString
	Description                    sql.NullString
	Description_2                  sql.NullString
	Name_2                         sql.NullString
	WidthInMillimeters             sql.NullString
	LengthInMillimeters            sql.NullString
	HeightInMillimeters            sql.NullString
	Shape                          NullVesselShape
	ID_3                           sql.NullString
	RecipeStepProductID            sql.NullString
	IconPath_2                     sql.NullString
	Slug                           sql.NullString
	MaximumQuantity                sql.NullInt32
	MinimumQuantity                int32
	UsableForStorage               sql.NullBool
	DisplayInSummaryLists          sql.NullBool
	IncludeInGeneratedInstructions sql.NullBool
	Volumetric                     sql.NullBool
	Imperial                       sql.NullBool
	Metric                         sql.NullBool
	Universal                      sql.NullBool
	UnavailableAfterStep           bool
}

func (q *Queries) GetRecipeStepVesselsForRecipe(ctx context.Context, db DBTX, belongsToRecipe string) ([]*GetRecipeStepVesselsForRecipeRow, error) {
	rows, err := db.QueryContext(ctx, getRecipeStepVesselsForRecipe, belongsToRecipe)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetRecipeStepVesselsForRecipeRow{}
	for rows.Next() {
		var i GetRecipeStepVesselsForRecipeRow
		if err := rows.Scan(
			&i.ID,
			&i.ID_2,
			&i.Name,
			&i.PluralName,
			&i.Description,
			&i.IconPath,
			&i.UsableForStorage,
			&i.Slug,
			&i.DisplayInSummaryLists,
			&i.IncludeInGeneratedInstructions,
			&i.Capacity,
			&i.ID_3,
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
			&i.Name_3,
			&i.Notes,
			&i.BelongsToRecipeStep,
			&i.RecipeStepProductID,
			&i.VesselPredicate,
			&i.MinimumQuantity,
			&i.MaximumQuantity,
			&i.UnavailableAfterStep,
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

const updateRecipeStepVessel = `-- name: UpdateRecipeStepVessel :exec

UPDATE recipe_step_vessels SET
	name = $1,
	notes = $2,
	belongs_to_recipe_step = $3,
	recipe_step_product_id = $4,
	valid_vessel_id = $5,
	vessel_predicate = $6,
	minimum_quantity = $7,
    maximum_quantity = $8,
    unavailable_after_step = $9,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_recipe_step = $3
	AND id = $10
`

type UpdateRecipeStepVesselParams struct {
	Name                 string
	Notes                string
	RecipeStepID         string
	VesselPredicate      string
	ID                   string
	RecipeStepProductID  sql.NullString
	ValidVesselID        sql.NullString
	MaximumQuantity      sql.NullInt32
	MinimumQuantity      int32
	UnavailableAfterStep bool
}

func (q *Queries) UpdateRecipeStepVessel(ctx context.Context, db DBTX, arg *UpdateRecipeStepVesselParams) error {
	_, err := db.ExecContext(ctx, updateRecipeStepVessel,
		arg.Name,
		arg.Notes,
		arg.RecipeStepID,
		arg.RecipeStepProductID,
		arg.ValidVesselID,
		arg.VesselPredicate,
		arg.MinimumQuantity,
		arg.MaximumQuantity,
		arg.UnavailableAfterStep,
		arg.ID,
	)
	return err
}
