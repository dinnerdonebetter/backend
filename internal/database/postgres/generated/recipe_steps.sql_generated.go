// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: recipe_steps.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const archiveRecipeStep = `-- name: ArchiveRecipeStep :execrows

UPDATE recipe_steps SET archived_at = NOW() WHERE archived_at IS NULL AND belongs_to_recipe = $1 AND id = $2
`

type ArchiveRecipeStepParams struct {
	BelongsToRecipe string
	ID              string
}

func (q *Queries) ArchiveRecipeStep(ctx context.Context, db DBTX, arg *ArchiveRecipeStepParams) (int64, error) {
	result, err := db.ExecContext(ctx, archiveRecipeStep, arg.BelongsToRecipe, arg.ID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const checkRecipeStepExistence = `-- name: CheckRecipeStepExistence :one

SELECT EXISTS (
	SELECT recipe_steps.id
	FROM recipe_steps
	JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
	WHERE recipe_steps.archived_at IS NULL
	  AND recipe_steps.belongs_to_recipe = $1
	  AND recipe_steps.id = $2
	  AND recipes.archived_at IS NULL
	  AND recipes.id = $1
)
`

type CheckRecipeStepExistenceParams struct {
	BelongsToRecipe string
	ID              string
}

func (q *Queries) CheckRecipeStepExistence(ctx context.Context, db DBTX, arg *CheckRecipeStepExistenceParams) (bool, error) {
	row := db.QueryRowContext(ctx, checkRecipeStepExistence, arg.BelongsToRecipe, arg.ID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createRecipeStep = `-- name: CreateRecipeStep :exec

INSERT INTO recipe_steps
(id,index,preparation_id,minimum_estimated_time_in_seconds,maximum_estimated_time_in_seconds,minimum_temperature_in_celsius,maximum_temperature_in_celsius,notes,explicit_instructions,condition_expression,optional,start_timer_automatically,belongs_to_recipe)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
`

type CreateRecipeStepParams struct {
	ID                            string
	BelongsToRecipe               string
	PreparationID                 string
	ConditionExpression           string
	ExplicitInstructions          string
	Notes                         string
	MaximumTemperatureInCelsius   sql.NullString
	MinimumTemperatureInCelsius   sql.NullString
	MaximumEstimatedTimeInSeconds sql.NullInt64
	MinimumEstimatedTimeInSeconds sql.NullInt64
	Index                         int32
	Optional                      bool
	StartTimerAutomatically       bool
}

func (q *Queries) CreateRecipeStep(ctx context.Context, db DBTX, arg *CreateRecipeStepParams) error {
	_, err := db.ExecContext(ctx, createRecipeStep,
		arg.ID,
		arg.Index,
		arg.PreparationID,
		arg.MinimumEstimatedTimeInSeconds,
		arg.MaximumEstimatedTimeInSeconds,
		arg.MinimumTemperatureInCelsius,
		arg.MaximumTemperatureInCelsius,
		arg.Notes,
		arg.ExplicitInstructions,
		arg.ConditionExpression,
		arg.Optional,
		arg.StartTimerAutomatically,
		arg.BelongsToRecipe,
	)
	return err
}

const getRecipeStep = `-- name: GetRecipeStep :one

SELECT
    recipe_steps.id,
    recipe_steps.index,
    valid_preparations.id as valid_preparation_id,
    valid_preparations.name as valid_preparation_name,
    valid_preparations.description as valid_preparation_description,
    valid_preparations.icon_path as valid_preparation_icon_path,
    valid_preparations.yields_nothing as valid_preparation_yields_nothing,
    valid_preparations.restrict_to_ingredients as valid_preparation_restrict_to_ingredients,
    valid_preparations.minimum_ingredient_count as valid_preparation_minimum_ingredient_count,
    valid_preparations.maximum_ingredient_count as valid_preparation_maximum_ingredient_count,
    valid_preparations.minimum_instrument_count as valid_preparation_minimum_instrument_count,
    valid_preparations.maximum_instrument_count as valid_preparation_maximum_instrument_count,
    valid_preparations.temperature_required as valid_preparation_temperature_required,
    valid_preparations.time_estimate_required as valid_preparation_time_estimate_required,
    valid_preparations.condition_expression_required as valid_preparation_condition_expression_required,
    valid_preparations.consumes_vessel as valid_preparation_consumes_vessel,
    valid_preparations.only_for_vessels as valid_preparation_only_for_vessels,
    valid_preparations.minimum_vessel_count as valid_preparation_minimum_vessel_count,
    valid_preparations.maximum_vessel_count as valid_preparation_maximum_vessel_count,
    valid_preparations.slug as valid_preparation_slug,
    valid_preparations.past_tense as valid_preparation_past_tense,
    valid_preparations.created_at as valid_preparation_created_at,
    valid_preparations.last_updated_at as valid_preparation_last_updated_at,
    valid_preparations.archived_at as valid_preparation_archived_at,
    recipe_steps.minimum_estimated_time_in_seconds,
    recipe_steps.maximum_estimated_time_in_seconds,
    recipe_steps.minimum_temperature_in_celsius,
    recipe_steps.maximum_temperature_in_celsius,
    recipe_steps.notes,
    recipe_steps.explicit_instructions,
    recipe_steps.condition_expression,
    recipe_steps.optional,
    recipe_steps.start_timer_automatically,
    recipe_steps.created_at,
    recipe_steps.last_updated_at,
    recipe_steps.archived_at,
    recipe_steps.belongs_to_recipe
FROM recipe_steps
    JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
    JOIN valid_preparations ON recipe_steps.preparation_id=valid_preparations.id
WHERE recipe_steps.archived_at IS NULL
    AND recipe_steps.belongs_to_recipe = $1
    AND recipe_steps.id = $2
    AND recipes.archived_at IS NULL
    AND recipes.id = $1
`

type GetRecipeStepParams struct {
	BelongsToRecipe string
	ID              string
}

type GetRecipeStepRow struct {
	ValidPreparationCreatedAt                   time.Time
	CreatedAt                                   time.Time
	ArchivedAt                                  sql.NullTime
	ValidPreparationLastUpdatedAt               sql.NullTime
	ValidPreparationArchivedAt                  sql.NullTime
	LastUpdatedAt                               sql.NullTime
	ConditionExpression                         string
	ValidPreparationID                          string
	ValidPreparationDescription                 string
	ID                                          string
	ValidPreparationName                        string
	ExplicitInstructions                        string
	Notes                                       string
	ValidPreparationIconPath                    string
	ValidPreparationPastTense                   string
	ValidPreparationSlug                        string
	BelongsToRecipe                             string
	MaximumTemperatureInCelsius                 sql.NullString
	MinimumTemperatureInCelsius                 sql.NullString
	MinimumEstimatedTimeInSeconds               sql.NullInt64
	MaximumEstimatedTimeInSeconds               sql.NullInt64
	ValidPreparationMaximumIngredientCount      sql.NullInt32
	ValidPreparationMaximumVesselCount          sql.NullInt32
	ValidPreparationMaximumInstrumentCount      sql.NullInt32
	ValidPreparationMinimumVesselCount          int32
	Index                                       int32
	ValidPreparationMinimumIngredientCount      int32
	ValidPreparationMinimumInstrumentCount      int32
	ValidPreparationTemperatureRequired         bool
	ValidPreparationTimeEstimateRequired        bool
	ValidPreparationConditionExpressionRequired bool
	Optional                                    bool
	StartTimerAutomatically                     bool
	ValidPreparationConsumesVessel              bool
	ValidPreparationRestrictToIngredients       bool
	ValidPreparationYieldsNothing               bool
	ValidPreparationOnlyForVessels              bool
}

func (q *Queries) GetRecipeStep(ctx context.Context, db DBTX, arg *GetRecipeStepParams) (*GetRecipeStepRow, error) {
	row := db.QueryRowContext(ctx, getRecipeStep, arg.BelongsToRecipe, arg.ID)
	var i GetRecipeStepRow
	err := row.Scan(
		&i.ID,
		&i.Index,
		&i.ValidPreparationID,
		&i.ValidPreparationName,
		&i.ValidPreparationDescription,
		&i.ValidPreparationIconPath,
		&i.ValidPreparationYieldsNothing,
		&i.ValidPreparationRestrictToIngredients,
		&i.ValidPreparationMinimumIngredientCount,
		&i.ValidPreparationMaximumIngredientCount,
		&i.ValidPreparationMinimumInstrumentCount,
		&i.ValidPreparationMaximumInstrumentCount,
		&i.ValidPreparationTemperatureRequired,
		&i.ValidPreparationTimeEstimateRequired,
		&i.ValidPreparationConditionExpressionRequired,
		&i.ValidPreparationConsumesVessel,
		&i.ValidPreparationOnlyForVessels,
		&i.ValidPreparationMinimumVesselCount,
		&i.ValidPreparationMaximumVesselCount,
		&i.ValidPreparationSlug,
		&i.ValidPreparationPastTense,
		&i.ValidPreparationCreatedAt,
		&i.ValidPreparationLastUpdatedAt,
		&i.ValidPreparationArchivedAt,
		&i.MinimumEstimatedTimeInSeconds,
		&i.MaximumEstimatedTimeInSeconds,
		&i.MinimumTemperatureInCelsius,
		&i.MaximumTemperatureInCelsius,
		&i.Notes,
		&i.ExplicitInstructions,
		&i.ConditionExpression,
		&i.Optional,
		&i.StartTimerAutomatically,
		&i.CreatedAt,
		&i.LastUpdatedAt,
		&i.ArchivedAt,
		&i.BelongsToRecipe,
	)
	return &i, err
}

const getRecipeStepByRecipeID = `-- name: GetRecipeStepByRecipeID :one

SELECT
	recipe_steps.id,
	recipe_steps.index,
    valid_preparations.id as valid_preparation_id,
    valid_preparations.name as valid_preparation_name,
    valid_preparations.description as valid_preparation_description,
    valid_preparations.icon_path as valid_preparation_icon_path,
    valid_preparations.yields_nothing as valid_preparation_yields_nothing,
    valid_preparations.restrict_to_ingredients as valid_preparation_restrict_to_ingredients,
    valid_preparations.minimum_ingredient_count as valid_preparation_minimum_ingredient_count,
    valid_preparations.maximum_ingredient_count as valid_preparation_maximum_ingredient_count,
    valid_preparations.minimum_instrument_count as valid_preparation_minimum_instrument_count,
    valid_preparations.maximum_instrument_count as valid_preparation_maximum_instrument_count,
    valid_preparations.temperature_required as valid_preparation_temperature_required,
    valid_preparations.time_estimate_required as valid_preparation_time_estimate_required,
    valid_preparations.condition_expression_required as valid_preparation_condition_expression_required,
    valid_preparations.consumes_vessel as valid_preparation_consumes_vessel,
    valid_preparations.only_for_vessels as valid_preparation_only_for_vessels,
    valid_preparations.minimum_vessel_count as valid_preparation_minimum_vessel_count,
    valid_preparations.maximum_vessel_count as valid_preparation_maximum_vessel_count,
    valid_preparations.slug as valid_preparation_slug,
    valid_preparations.past_tense as valid_preparation_past_tense,
    valid_preparations.created_at as valid_preparation_created_at,
    valid_preparations.last_updated_at as valid_preparation_last_updated_at,
    valid_preparations.archived_at as valid_preparation_archived_at,
	recipe_steps.minimum_estimated_time_in_seconds,
	recipe_steps.maximum_estimated_time_in_seconds,
	recipe_steps.minimum_temperature_in_celsius,
	recipe_steps.maximum_temperature_in_celsius,
	recipe_steps.notes,
	recipe_steps.explicit_instructions,
	recipe_steps.condition_expression,
	recipe_steps.optional,
	recipe_steps.start_timer_automatically,
	recipe_steps.created_at,
	recipe_steps.last_updated_at,
	recipe_steps.archived_at,
	recipe_steps.belongs_to_recipe
FROM recipe_steps
	JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
	JOIN valid_preparations ON recipe_steps.preparation_id=valid_preparations.id
WHERE recipe_steps.archived_at IS NULL
	AND recipe_steps.id = $1
`

type GetRecipeStepByRecipeIDRow struct {
	ValidPreparationCreatedAt                   time.Time
	CreatedAt                                   time.Time
	ArchivedAt                                  sql.NullTime
	ValidPreparationLastUpdatedAt               sql.NullTime
	ValidPreparationArchivedAt                  sql.NullTime
	LastUpdatedAt                               sql.NullTime
	ConditionExpression                         string
	ValidPreparationID                          string
	ValidPreparationDescription                 string
	ID                                          string
	ValidPreparationName                        string
	ExplicitInstructions                        string
	Notes                                       string
	ValidPreparationIconPath                    string
	ValidPreparationPastTense                   string
	ValidPreparationSlug                        string
	BelongsToRecipe                             string
	MaximumTemperatureInCelsius                 sql.NullString
	MinimumTemperatureInCelsius                 sql.NullString
	MinimumEstimatedTimeInSeconds               sql.NullInt64
	MaximumEstimatedTimeInSeconds               sql.NullInt64
	ValidPreparationMaximumIngredientCount      sql.NullInt32
	ValidPreparationMaximumVesselCount          sql.NullInt32
	ValidPreparationMaximumInstrumentCount      sql.NullInt32
	ValidPreparationMinimumVesselCount          int32
	Index                                       int32
	ValidPreparationMinimumIngredientCount      int32
	ValidPreparationMinimumInstrumentCount      int32
	ValidPreparationTemperatureRequired         bool
	ValidPreparationTimeEstimateRequired        bool
	ValidPreparationConditionExpressionRequired bool
	Optional                                    bool
	StartTimerAutomatically                     bool
	ValidPreparationConsumesVessel              bool
	ValidPreparationRestrictToIngredients       bool
	ValidPreparationYieldsNothing               bool
	ValidPreparationOnlyForVessels              bool
}

func (q *Queries) GetRecipeStepByRecipeID(ctx context.Context, db DBTX, id string) (*GetRecipeStepByRecipeIDRow, error) {
	row := db.QueryRowContext(ctx, getRecipeStepByRecipeID, id)
	var i GetRecipeStepByRecipeIDRow
	err := row.Scan(
		&i.ID,
		&i.Index,
		&i.ValidPreparationID,
		&i.ValidPreparationName,
		&i.ValidPreparationDescription,
		&i.ValidPreparationIconPath,
		&i.ValidPreparationYieldsNothing,
		&i.ValidPreparationRestrictToIngredients,
		&i.ValidPreparationMinimumIngredientCount,
		&i.ValidPreparationMaximumIngredientCount,
		&i.ValidPreparationMinimumInstrumentCount,
		&i.ValidPreparationMaximumInstrumentCount,
		&i.ValidPreparationTemperatureRequired,
		&i.ValidPreparationTimeEstimateRequired,
		&i.ValidPreparationConditionExpressionRequired,
		&i.ValidPreparationConsumesVessel,
		&i.ValidPreparationOnlyForVessels,
		&i.ValidPreparationMinimumVesselCount,
		&i.ValidPreparationMaximumVesselCount,
		&i.ValidPreparationSlug,
		&i.ValidPreparationPastTense,
		&i.ValidPreparationCreatedAt,
		&i.ValidPreparationLastUpdatedAt,
		&i.ValidPreparationArchivedAt,
		&i.MinimumEstimatedTimeInSeconds,
		&i.MaximumEstimatedTimeInSeconds,
		&i.MinimumTemperatureInCelsius,
		&i.MaximumTemperatureInCelsius,
		&i.Notes,
		&i.ExplicitInstructions,
		&i.ConditionExpression,
		&i.Optional,
		&i.StartTimerAutomatically,
		&i.CreatedAt,
		&i.LastUpdatedAt,
		&i.ArchivedAt,
		&i.BelongsToRecipe,
	)
	return &i, err
}

const getRecipeSteps = `-- name: GetRecipeSteps :many

SELECT
    recipe_steps.id,
    recipe_steps.index,
    valid_preparations.id as valid_preparation_id,
    valid_preparations.name as valid_preparation_name,
    valid_preparations.description as valid_preparation_description,
    valid_preparations.icon_path as valid_preparation_icon_path,
    valid_preparations.yields_nothing as valid_preparation_yields_nothing,
    valid_preparations.restrict_to_ingredients as valid_preparation_restrict_to_ingredients,
    valid_preparations.minimum_ingredient_count as valid_preparation_minimum_ingredient_count,
    valid_preparations.maximum_ingredient_count as valid_preparation_maximum_ingredient_count,
    valid_preparations.minimum_instrument_count as valid_preparation_minimum_instrument_count,
    valid_preparations.maximum_instrument_count as valid_preparation_maximum_instrument_count,
    valid_preparations.temperature_required as valid_preparation_temperature_required,
    valid_preparations.time_estimate_required as valid_preparation_time_estimate_required,
    valid_preparations.condition_expression_required as valid_preparation_condition_expression_required,
    valid_preparations.consumes_vessel as valid_preparation_consumes_vessel,
    valid_preparations.only_for_vessels as valid_preparation_only_for_vessels,
    valid_preparations.minimum_vessel_count as valid_preparation_minimum_vessel_count,
    valid_preparations.maximum_vessel_count as valid_preparation_maximum_vessel_count,
    valid_preparations.slug as valid_preparation_slug,
    valid_preparations.past_tense as valid_preparation_past_tense,
    valid_preparations.created_at as valid_preparation_created_at,
    valid_preparations.last_updated_at as valid_preparation_last_updated_at,
    valid_preparations.archived_at as valid_preparation_archived_at,
    recipe_steps.minimum_estimated_time_in_seconds,
    recipe_steps.maximum_estimated_time_in_seconds,
    recipe_steps.minimum_temperature_in_celsius,
    recipe_steps.maximum_temperature_in_celsius,
    recipe_steps.notes,
    recipe_steps.explicit_instructions,
    recipe_steps.condition_expression,
    recipe_steps.optional,
    recipe_steps.start_timer_automatically,
    recipe_steps.created_at,
    recipe_steps.last_updated_at,
    recipe_steps.archived_at,
    recipe_steps.belongs_to_recipe,
    (
        SELECT
            COUNT(recipe_steps.id)
        FROM
            recipe_steps
        WHERE
            recipe_steps.archived_at IS NULL
            AND recipe_steps.created_at > COALESCE($1, (SELECT NOW() - interval '999 years'))
            AND recipe_steps.created_at < COALESCE($2, (SELECT NOW() + interval '999 years'))
            AND (recipe_steps.last_updated_at IS NULL OR recipe_steps.last_updated_at > COALESCE($3, (SELECT NOW() - interval '999 years')))
            AND (recipe_steps.last_updated_at IS NULL OR recipe_steps.last_updated_at < COALESCE($4, (SELECT NOW() + interval '999 years')))
    ) as filtered_count,
    (
        SELECT
            COUNT(recipe_steps.id)
        FROM
            recipe_steps
        WHERE
            recipe_steps.archived_at IS NULL
    ) as total_count
FROM recipe_steps
    JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
    JOIN valid_preparations ON recipe_steps.preparation_id=valid_preparations.id
WHERE recipe_steps.archived_at IS NULL
    AND recipe_steps.created_at > COALESCE($1, (SELECT NOW() - interval '999 years'))
    AND recipe_steps.created_at < COALESCE($2, (SELECT NOW() + interval '999 years'))
    AND (recipe_steps.last_updated_at IS NULL OR recipe_steps.last_updated_at > COALESCE($3, (SELECT NOW() - interval '999 years')))
    AND (recipe_steps.last_updated_at IS NULL OR recipe_steps.last_updated_at < COALESCE($4, (SELECT NOW() + interval '999 years')))
    AND recipe_steps.belongs_to_recipe = $5
    AND recipes.archived_at IS NULL
    OFFSET $6
    LIMIT $7
`

type GetRecipeStepsParams struct {
	CreatedAfter  sql.NullTime
	CreatedBefore sql.NullTime
	UpdatedAfter  sql.NullTime
	UpdatedBefore sql.NullTime
	RecipeID      string
	QueryOffset   sql.NullInt32
	QueryLimit    sql.NullInt32
}

type GetRecipeStepsRow struct {
	ValidPreparationCreatedAt                   time.Time
	CreatedAt                                   time.Time
	ArchivedAt                                  sql.NullTime
	ValidPreparationLastUpdatedAt               sql.NullTime
	LastUpdatedAt                               sql.NullTime
	ValidPreparationArchivedAt                  sql.NullTime
	ValidPreparationIconPath                    string
	ValidPreparationSlug                        string
	ValidPreparationID                          string
	ID                                          string
	ValidPreparationDescription                 string
	ValidPreparationName                        string
	ConditionExpression                         string
	ExplicitInstructions                        string
	Notes                                       string
	ValidPreparationPastTense                   string
	BelongsToRecipe                             string
	MinimumTemperatureInCelsius                 sql.NullString
	MaximumTemperatureInCelsius                 sql.NullString
	MaximumEstimatedTimeInSeconds               sql.NullInt64
	MinimumEstimatedTimeInSeconds               sql.NullInt64
	FilteredCount                               int64
	TotalCount                                  int64
	ValidPreparationMaximumIngredientCount      sql.NullInt32
	ValidPreparationMaximumInstrumentCount      sql.NullInt32
	ValidPreparationMaximumVesselCount          sql.NullInt32
	ValidPreparationMinimumVesselCount          int32
	Index                                       int32
	ValidPreparationMinimumIngredientCount      int32
	ValidPreparationMinimumInstrumentCount      int32
	ValidPreparationConsumesVessel              bool
	Optional                                    bool
	StartTimerAutomatically                     bool
	ValidPreparationTemperatureRequired         bool
	ValidPreparationTimeEstimateRequired        bool
	ValidPreparationConditionExpressionRequired bool
	ValidPreparationRestrictToIngredients       bool
	ValidPreparationYieldsNothing               bool
	ValidPreparationOnlyForVessels              bool
}

func (q *Queries) GetRecipeSteps(ctx context.Context, db DBTX, arg *GetRecipeStepsParams) ([]*GetRecipeStepsRow, error) {
	rows, err := db.QueryContext(ctx, getRecipeSteps,
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
	items := []*GetRecipeStepsRow{}
	for rows.Next() {
		var i GetRecipeStepsRow
		if err := rows.Scan(
			&i.ID,
			&i.Index,
			&i.ValidPreparationID,
			&i.ValidPreparationName,
			&i.ValidPreparationDescription,
			&i.ValidPreparationIconPath,
			&i.ValidPreparationYieldsNothing,
			&i.ValidPreparationRestrictToIngredients,
			&i.ValidPreparationMinimumIngredientCount,
			&i.ValidPreparationMaximumIngredientCount,
			&i.ValidPreparationMinimumInstrumentCount,
			&i.ValidPreparationMaximumInstrumentCount,
			&i.ValidPreparationTemperatureRequired,
			&i.ValidPreparationTimeEstimateRequired,
			&i.ValidPreparationConditionExpressionRequired,
			&i.ValidPreparationConsumesVessel,
			&i.ValidPreparationOnlyForVessels,
			&i.ValidPreparationMinimumVesselCount,
			&i.ValidPreparationMaximumVesselCount,
			&i.ValidPreparationSlug,
			&i.ValidPreparationPastTense,
			&i.ValidPreparationCreatedAt,
			&i.ValidPreparationLastUpdatedAt,
			&i.ValidPreparationArchivedAt,
			&i.MinimumEstimatedTimeInSeconds,
			&i.MaximumEstimatedTimeInSeconds,
			&i.MinimumTemperatureInCelsius,
			&i.MaximumTemperatureInCelsius,
			&i.Notes,
			&i.ExplicitInstructions,
			&i.ConditionExpression,
			&i.Optional,
			&i.StartTimerAutomatically,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
			&i.BelongsToRecipe,
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

const updateRecipeStep = `-- name: UpdateRecipeStep :execrows

UPDATE recipe_steps SET
	index = $1,
	preparation_id = $2,
	minimum_estimated_time_in_seconds = $3,
	maximum_estimated_time_in_seconds = $4,
	minimum_temperature_in_celsius = $5,
	maximum_temperature_in_celsius = $6,
	notes = $7,
	explicit_instructions = $8,
	condition_expression = $9,
	optional = $10,
	start_timer_automatically = $11,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_recipe = $12
	AND id = $13
`

type UpdateRecipeStepParams struct {
	ConditionExpression           string
	PreparationID                 string
	ID                            string
	BelongsToRecipe               string
	Notes                         string
	ExplicitInstructions          string
	MinimumTemperatureInCelsius   sql.NullString
	MaximumTemperatureInCelsius   sql.NullString
	MaximumEstimatedTimeInSeconds sql.NullInt64
	MinimumEstimatedTimeInSeconds sql.NullInt64
	Index                         int32
	Optional                      bool
	StartTimerAutomatically       bool
}

func (q *Queries) UpdateRecipeStep(ctx context.Context, db DBTX, arg *UpdateRecipeStepParams) (int64, error) {
	result, err := db.ExecContext(ctx, updateRecipeStep,
		arg.Index,
		arg.PreparationID,
		arg.MinimumEstimatedTimeInSeconds,
		arg.MaximumEstimatedTimeInSeconds,
		arg.MinimumTemperatureInCelsius,
		arg.MaximumTemperatureInCelsius,
		arg.Notes,
		arg.ExplicitInstructions,
		arg.ConditionExpression,
		arg.Optional,
		arg.StartTimerAutomatically,
		arg.BelongsToRecipe,
		arg.ID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}