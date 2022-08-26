package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	recipeStepsOnRecipeStepProductsJoinClause = "recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id"
)

var (
	_ types.RecipeStepProductDataManager = (*SQLQuerier)(nil)

	// recipeStepProductsTableColumns are the columns for the recipe_step_products table.
	recipeStepProductsTableColumns = []string{
		"recipe_step_products.id",
		"recipe_step_products.name",
		"recipe_step_products.type",
		"valid_measurement_units.id",
		"valid_measurement_units.name",
		"valid_measurement_units.description",
		"valid_measurement_units.volumetric",
		"valid_measurement_units.icon_path",
		"valid_measurement_units.universal",
		"valid_measurement_units.metric",
		"valid_measurement_units.imperial",
		"valid_measurement_units.plural_name",
		"valid_measurement_units.created_on",
		"valid_measurement_units.last_updated_on",
		"valid_measurement_units.archived_on",
		"recipe_step_products.minimum_quantity_value",
		"recipe_step_products.maximum_quantity_value",
		"recipe_step_products.quantity_notes",
		"recipe_step_products.compostable",
		"recipe_step_products.maximum_storage_duration_in_seconds",
		"recipe_step_products.minimum_storage_temperature_in_celsius",
		"recipe_step_products.maximum_storage_temperature_in_celsius",
		"recipe_step_products.created_on",
		"recipe_step_products.last_updated_on",
		"recipe_step_products.archived_on",
		"recipe_step_products.belongs_to_recipe_step",
	}

	getRecipeStepProductsJoins = []string{
		recipeStepsOnRecipeStepProductsJoinClause,
		recipesOnRecipeStepsJoinClause,
		validMeasurementUnitsOnRecipeStepProductsJoinClause,
	}
)

// scanRecipeStepProduct takes a database Scanner (i.e. *sql.Row) and scans the result into a recipe step product struct.
func (q *SQLQuerier) scanRecipeStepProduct(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.RecipeStepProduct, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.RecipeStepProduct{}

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.Type,
		&x.MeasurementUnit.ID,
		&x.MeasurementUnit.Name,
		&x.MeasurementUnit.Description,
		&x.MeasurementUnit.Volumetric,
		&x.MeasurementUnit.IconPath,
		&x.MeasurementUnit.Universal,
		&x.MeasurementUnit.Metric,
		&x.MeasurementUnit.Imperial,
		&x.MeasurementUnit.PluralName,
		&x.MeasurementUnit.CreatedOn,
		&x.MeasurementUnit.LastUpdatedOn,
		&x.MeasurementUnit.ArchivedOn,
		&x.MinimumQuantityValue,
		&x.MaximumQuantityValue,
		&x.QuantityNotes,
		&x.Compostable,
		&x.MaximumStorageDurationInSeconds,
		&x.MinimumStorageTemperatureInCelsius,
		&x.MaximumStorageTemperatureInCelsius,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToRecipeStep,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "")
	}

	return x, filteredCount, totalCount, nil
}

// scanRecipeStepProducts takes some database rows and turns them into a slice of recipe step products.
func (q *SQLQuerier) scanRecipeStepProducts(ctx context.Context, rows database.ResultIterator, includeCounts bool) (recipeStepProducts []*types.RecipeStepProduct, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanRecipeStepProduct(ctx, rows, includeCounts)
		if scanErr != nil {
			return nil, 0, 0, scanErr
		}

		if includeCounts {
			if filteredCount == 0 {
				filteredCount = fc
			}

			if totalCount == 0 {
				totalCount = tc
			}
		}

		recipeStepProducts = append(recipeStepProducts, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return recipeStepProducts, filteredCount, totalCount, nil
}

const recipeStepProductExistenceQuery = "SELECT EXISTS ( SELECT recipe_step_products.id FROM recipe_step_products JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_products.archived_on IS NULL AND recipe_step_products.belongs_to_recipe_step = $1 AND recipe_step_products.id = $2 AND recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.archived_on IS NULL AND recipes.id = $5 )"

// RecipeStepProductExists fetches whether a recipe step product exists from the database.
func (q *SQLQuerier) RecipeStepProductExists(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepProductID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	args := []interface{}{
		recipeStepID,
		recipeStepProductID,
		recipeID,
		recipeStepID,
		recipeID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, recipeStepProductExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing recipe step product existence check")
	}

	return result, nil
}

const getRecipeStepProductQuery = `SELECT
	recipe_step_products.id,
	recipe_step_products.name,
	recipe_step_products.type,
	valid_measurement_units.id,
	valid_measurement_units.name,
	valid_measurement_units.description,
	valid_measurement_units.volumetric,
	valid_measurement_units.icon_path,
	valid_measurement_units.universal,
	valid_measurement_units.metric,
	valid_measurement_units.imperial,
	valid_measurement_units.plural_name,
	valid_measurement_units.created_on,
	valid_measurement_units.last_updated_on,
	valid_measurement_units.archived_on,
	recipe_step_products.minimum_quantity_value,
	recipe_step_products.maximum_quantity_value,
	recipe_step_products.quantity_notes,
	recipe_step_products.compostable,
	recipe_step_products.maximum_storage_duration_in_seconds,
	recipe_step_products.minimum_storage_temperature_in_celsius,
	recipe_step_products.maximum_storage_temperature_in_celsius,
	recipe_step_products.created_on,
	recipe_step_products.last_updated_on,
	recipe_step_products.archived_on,
	recipe_step_products.belongs_to_recipe_step
FROM recipe_step_products
JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id
JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
JOIN valid_measurement_units ON recipe_step_products.measurement_unit=valid_measurement_units.id
WHERE recipe_step_products.archived_on IS NULL
AND recipe_step_products.belongs_to_recipe_step = $1 
AND recipe_step_products.id = $2
AND recipe_steps.archived_on IS NULL
AND recipe_steps.belongs_to_recipe = $3
AND recipe_steps.id = $4
AND recipes.archived_on IS NULL 
AND recipes.id = $5
`

// GetRecipeStepProduct fetches a recipe step product from the database.
func (q *SQLQuerier) GetRecipeStepProduct(ctx context.Context, recipeID, recipeStepID, recipeStepProductID string) (*types.RecipeStepProduct, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepProductID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	args := []interface{}{
		recipeStepID,
		recipeStepProductID,
		recipeID,
		recipeStepID,
		recipeID,
	}

	row := q.getOneRow(ctx, q.db, "recipeStepProduct", getRecipeStepProductQuery, args)

	recipeStepProduct, _, _, err := q.scanRecipeStepProduct(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipeStepProduct")
	}

	return recipeStepProduct, nil
}

const getTotalRecipeStepProductsCountQuery = "SELECT COUNT(recipe_step_products.id) FROM recipe_step_products WHERE recipe_step_products.archived_on IS NULL"

// GetTotalRecipeStepProductCount fetches the count of recipe step products from the database that meet a particular filter.
func (q *SQLQuerier) GetTotalRecipeStepProductCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	count, err := q.performCountQuery(ctx, q.db, getTotalRecipeStepProductsCountQuery, "fetching count of recipe step products")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of recipe step products")
	}

	return count, nil
}

const getRecipeStepProductsForRecipeQuery = `SELECT
	recipe_step_products.id,
	recipe_step_products.name,
	recipe_step_products.type,
	valid_measurement_units.id,
	valid_measurement_units.name,
	valid_measurement_units.description,
	valid_measurement_units.volumetric,
	valid_measurement_units.icon_path,
	valid_measurement_units.universal,
	valid_measurement_units.metric,
	valid_measurement_units.imperial,
	valid_measurement_units.plural_name,
	valid_measurement_units.created_on,
	valid_measurement_units.last_updated_on,
	valid_measurement_units.archived_on,
	recipe_step_products.minimum_quantity_value,
	recipe_step_products.maximum_quantity_value,
	recipe_step_products.quantity_notes,
	recipe_step_products.compostable,
	recipe_step_products.maximum_storage_duration_in_seconds,
	recipe_step_products.minimum_storage_temperature_in_celsius,
	recipe_step_products.maximum_storage_temperature_in_celsius,
	recipe_step_products.created_on,
	recipe_step_products.last_updated_on,
	recipe_step_products.archived_on,
	recipe_step_products.belongs_to_recipe_step
FROM recipe_step_products
JOIN recipe_steps ON recipe_step_products.belongs_to_recipe_step=recipe_steps.id
JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
LEFT OUTER JOIN valid_measurement_units ON recipe_step_products.measurement_unit=valid_measurement_units.id
WHERE recipe_step_products.archived_on IS NULL
AND recipe_steps.archived_on IS NULL
AND recipe_steps.belongs_to_recipe = $1
AND recipes.archived_on IS NULL 
AND recipes.id = $2
`

// getRecipeStepProductsForRecipe fetches a list of recipe step products from the database that meet a particular filter.
func (q *SQLQuerier) getRecipeStepProductsForRecipe(ctx context.Context, recipeID string) ([]*types.RecipeStepProduct, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	args := []interface{}{
		recipeID,
		recipeID,
	}

	rows, err := q.performReadQuery(ctx, q.db, "recipe step products", getRecipeStepProductsForRecipeQuery, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing recipe step products list retrieval query")
	}

	recipeStepProducts, _, _, err := q.scanRecipeStepProducts(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipe step products")
	}

	return recipeStepProducts, nil
}

// GetRecipeStepProducts fetches a list of recipe step products from the database that meet a particular filter.
func (q *SQLQuerier) GetRecipeStepProducts(ctx context.Context, recipeID, recipeStepID string, filter *types.QueryFilter) (x *types.RecipeStepProductList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeIDKey, recipeID)
	tracing.AttachRecipeIDToSpan(span, recipeID)

	if recipeStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	x = &types.RecipeStepProductList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		if filter.Page != nil {
			x.Page = *filter.Page
		}

		if filter.Limit != nil {
			x.Limit = *filter.Limit
		}
	}

	query, args := q.buildListQuery(ctx, "recipe_step_products", getRecipeStepProductsJoins, []string{"valid_measurement_units.id"}, nil, householdOwnershipColumn, recipeStepProductsTableColumns, "", false, filter, true)

	rows, err := q.performReadQuery(ctx, q.db, "recipe step products", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing recipe step products list retrieval query")
	}

	if x.RecipeStepProducts, x.FilteredCount, x.TotalCount, err = q.scanRecipeStepProducts(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipe step products")
	}

	return x, nil
}

func (q *SQLQuerier) buildGetRecipeStepProductsWithIDsQuery(ctx context.Context, recipeStepID string, limit uint8, ids []string) (query string, args []interface{}) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	withIDsWhere := squirrel.Eq{
		"recipe_step_products.id":                     ids,
		"recipe_step_products.archived_on":            nil,
		"recipe_step_products.belongs_to_recipe_step": recipeStepID,
	}

	subqueryBuilder := q.sqlBuilder.Select(recipeStepProductsTableColumns...).
		From("recipe_step_products").
		Join(fmt.Sprintf("unnest('{%s}'::text[])", joinIDs(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))

	query, args, err := q.sqlBuilder.Select(recipeStepProductsTableColumns...).
		FromSelect(subqueryBuilder, "recipe_step_products").
		Where(withIDsWhere).ToSql()

	q.logQueryBuildingError(span, err)

	return query, args
}

// GetRecipeStepProductsWithIDs fetches recipe step products from the database within a given set of IDs.
func (q *SQLQuerier) GetRecipeStepProductsWithIDs(ctx context.Context, recipeStepID string, limit uint8, ids []string) ([]*types.RecipeStepProduct, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeStepID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if ids == nil {
		return nil, ErrNilInputProvided
	}

	if limit == 0 {
		limit = uint8(types.DefaultLimit)
	}

	logger = logger.WithValues(map[string]interface{}{
		"limit":    limit,
		"id_count": len(ids),
	})

	query, args := q.buildGetRecipeStepProductsWithIDsQuery(ctx, recipeStepID, limit, ids)

	rows, err := q.performReadQuery(ctx, q.db, "recipe step products with IDs", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching recipe step products from database")
	}

	recipeStepProducts, _, _, err := q.scanRecipeStepProducts(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning recipe step products")
	}

	return recipeStepProducts, nil
}

const recipeStepProductCreationQuery = `INSERT INTO recipe_step_products
    (id,name,type,measurement_unit,minimum_quantity_value,maximum_quantity_value,quantity_notes,compostable,maximum_storage_duration_in_seconds,minimum_storage_temperature_in_celsius,maximum_storage_temperature_in_celsius,belongs_to_recipe_step) 
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`

// CreateRecipeStepProduct creates a recipe step product in the database.
func (q *SQLQuerier) createRecipeStepProduct(ctx context.Context, db database.SQLQueryExecutor, input *types.RecipeStepProductDatabaseCreationInput) (*types.RecipeStepProduct, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeStepProductIDKey, input.ID)

	args := []interface{}{
		input.ID,
		input.Name,
		input.Type,
		input.MeasurementUnitID,
		input.MinimumQuantityValue,
		input.MaximumQuantityValue,
		input.QuantityNotes,
		input.Compostable,
		input.MaximumStorageDurationInSeconds,
		input.MinimumStorageTemperatureInCelsius,
		input.MaximumStorageTemperatureInCelsius,
		input.BelongsToRecipeStep,
	}

	// create the recipe step product.
	if err := q.performWriteQuery(ctx, db, "recipe step product creation", recipeStepProductCreationQuery, args); err != nil {
		return nil, observability.PrepareError(err, logger, span, "performing recipe step product creation query")
	}

	x := &types.RecipeStepProduct{
		ID:                                 input.ID,
		Name:                               input.Name,
		Type:                               input.Type,
		MeasurementUnit:                    types.ValidMeasurementUnit{ID: input.MeasurementUnitID},
		MinimumQuantityValue:               input.MinimumQuantityValue,
		MaximumQuantityValue:               input.MaximumQuantityValue,
		QuantityNotes:                      input.QuantityNotes,
		BelongsToRecipeStep:                input.BelongsToRecipeStep,
		Compostable:                        input.Compostable,
		MaximumStorageDurationInSeconds:    input.MaximumStorageDurationInSeconds,
		MinimumStorageTemperatureInCelsius: input.MinimumStorageTemperatureInCelsius,
		MaximumStorageTemperatureInCelsius: input.MaximumStorageTemperatureInCelsius,
		CreatedOn:                          q.currentTime(),
	}

	tracing.AttachRecipeStepProductIDToSpan(span, x.ID)

	return x, nil
}

// CreateRecipeStepProduct creates a recipe step product in the database.
func (q *SQLQuerier) CreateRecipeStepProduct(ctx context.Context, input *types.RecipeStepProductDatabaseCreationInput) (*types.RecipeStepProduct, error) {
	return q.createRecipeStepProduct(ctx, q.db, input)
}

const updateRecipeStepProductQuery = `
UPDATE recipe_step_products
SET 
    name = $1,
    type = $2,
    measurement_unit = $3,
    minimum_quantity_value = $4,
    maximum_quantity_value = $5,
    quantity_notes = $6,
    compostable = $7,
	maximum_storage_duration_in_seconds = $8,
	minimum_storage_temperature_in_celsius = $9,
	maximum_storage_temperature_in_celsius = $10,
    last_updated_on = extract(epoch FROM NOW())
WHERE archived_on IS NULL
  AND belongs_to_recipe_step = $11
  AND id = $12
`

// UpdateRecipeStepProduct updates a particular recipe step product.
func (q *SQLQuerier) UpdateRecipeStepProduct(ctx context.Context, updated *types.RecipeStepProduct) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.RecipeStepProductIDKey, updated.ID)
	tracing.AttachRecipeStepProductIDToSpan(span, updated.ID)

	args := []interface{}{
		updated.Name,
		updated.Type,
		updated.MeasurementUnit.ID,
		updated.MinimumQuantityValue,
		updated.MaximumQuantityValue,
		updated.QuantityNotes,
		updated.Compostable,
		updated.MaximumStorageDurationInSeconds,
		updated.MinimumStorageTemperatureInCelsius,
		updated.MaximumStorageTemperatureInCelsius,
		updated.BelongsToRecipeStep,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe step product update", updateRecipeStepProductQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating recipe step product")
	}

	logger.Info("recipe step product updated")

	return nil
}

const archiveRecipeStepProductQuery = "UPDATE recipe_step_products SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe_step = $1 AND id = $2"

// ArchiveRecipeStepProduct archives a recipe step product from the database by its ID.
func (q *SQLQuerier) ArchiveRecipeStepProduct(ctx context.Context, recipeStepID, recipeStepProductID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if recipeStepID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepIDKey, recipeStepID)
	tracing.AttachRecipeStepIDToSpan(span, recipeStepID)

	if recipeStepProductID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeStepProductIDKey, recipeStepProductID)
	tracing.AttachRecipeStepProductIDToSpan(span, recipeStepProductID)

	args := []interface{}{
		recipeStepID,
		recipeStepProductID,
	}

	if err := q.performWriteQuery(ctx, q.db, "recipe step product archive", archiveRecipeStepProductQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating recipe step product")
	}

	logger.Info("recipe step product archived")

	return nil
}