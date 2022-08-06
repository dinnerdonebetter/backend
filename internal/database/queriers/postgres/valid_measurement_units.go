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
	validMeasurementUnitsOnRecipeStepIngredientsJoinClause = `valid_measurement_units ON recipe_step_ingredients.measurement_unit=valid_measurement_units.id`
	validMeasurementUnitsOnRecipeStepProductsJoinClause    = `valid_measurement_units ON recipe_step_products.measurement_unit=valid_measurement_units.id`
)

var (
	_ types.ValidMeasurementUnitDataManager = (*SQLQuerier)(nil)

	// validMeasurementUnitsTableColumns are the columns for the valid_measurement_units table.
	validMeasurementUnitsTableColumns = []string{
		"valid_measurement_units.id",
		"valid_measurement_units.name",
		"valid_measurement_units.description",
		"valid_measurement_units.volumetric",
		"valid_measurement_units.icon_path",
		"valid_measurement_units.created_on",
		"valid_measurement_units.last_updated_on",
		"valid_measurement_units.archived_on",
	}
)

// scanValidMeasurementUnit takes a database Scanner (i.e. *sql.Row) and scans the result into a valid measurement unit struct.
func (q *SQLQuerier) scanValidMeasurementUnit(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.ValidMeasurementUnit, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	x = &types.ValidMeasurementUnit{}

	targetVars := []interface{}{
		&x.ID,
		&x.Name,
		&x.Description,
		&x.Volumetric,
		&x.IconPath,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "")
	}

	return x, filteredCount, totalCount, nil
}

// scanValidMeasurementUnits takes some database rows and turns them into a slice of valid measurement units.
func (q *SQLQuerier) scanValidMeasurementUnits(ctx context.Context, rows database.ResultIterator, includeCounts bool) (validMeasurementUnits []*types.ValidMeasurementUnit, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue("include_counts", includeCounts)

	for rows.Next() {
		x, fc, tc, scanErr := q.scanValidMeasurementUnit(ctx, rows, includeCounts)
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

		validMeasurementUnits = append(validMeasurementUnits, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, logger, span, "handling rows")
	}

	return validMeasurementUnits, filteredCount, totalCount, nil
}

const validMeasurementUnitExistenceQuery = "SELECT EXISTS ( SELECT valid_measurement_units.id FROM valid_measurement_units WHERE valid_measurement_units.archived_on IS NULL AND valid_measurement_units.id = $1 )"

// ValidMeasurementUnitExists fetches whether a valid measurement unit exists from the database.
func (q *SQLQuerier) ValidMeasurementUnitExists(ctx context.Context, validMeasurementUnitID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validMeasurementUnitID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachValidMeasurementUnitIDToSpan(span, validMeasurementUnitID)

	args := []interface{}{
		validMeasurementUnitID,
	}

	result, err := q.performBooleanQuery(ctx, q.db, validMeasurementUnitExistenceQuery, args)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "performing valid measurement unit existence check")
	}

	return result, nil
}

const getValidMeasurementUnitBaseQuery = `SELECT 
	valid_measurement_units.id,
	valid_measurement_units.name,
	valid_measurement_units.description,
	valid_measurement_units.volumetric,
	valid_measurement_units.icon_path,
	valid_measurement_units.created_on, 
	valid_measurement_units.last_updated_on, 
	valid_measurement_units.archived_on 
FROM valid_measurement_units 
WHERE valid_measurement_units.archived_on IS NULL
`

const getValidMeasurementUnitQuery = getValidMeasurementUnitBaseQuery + `AND valid_measurement_units.id = $1`

// GetValidMeasurementUnit fetches a valid measurement unit from the database.
func (q *SQLQuerier) GetValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) (*types.ValidMeasurementUnit, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validMeasurementUnitID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachValidMeasurementUnitIDToSpan(span, validMeasurementUnitID)

	args := []interface{}{
		validMeasurementUnitID,
	}

	row := q.getOneRow(ctx, q.db, "valid measurement unit", getValidMeasurementUnitQuery, args)

	validMeasurementUnit, _, _, err := q.scanValidMeasurementUnit(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid measurement unit")
	}

	return validMeasurementUnit, nil
}

const getRandomValidMeasurementUnitQuery = getValidMeasurementUnitBaseQuery + `ORDER BY random() LIMIT 1`

// GetRandomValidMeasurementUnit fetches a valid measurement unit from the database.
func (q *SQLQuerier) GetRandomValidMeasurementUnit(ctx context.Context) (*types.ValidMeasurementUnit, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()
	args := []interface{}{}

	row := q.getOneRow(ctx, q.db, "valid measurement unit", getRandomValidMeasurementUnitQuery, args)

	validMeasurementUnit, _, _, err := q.scanValidMeasurementUnit(ctx, row, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid measurement unit")
	}

	return validMeasurementUnit, nil
}

const validMeasurementUnitSearchQuery = `SELECT 
	valid_measurement_units.id,
	valid_measurement_units.name,
	valid_measurement_units.description,
	valid_measurement_units.volumetric,
	valid_measurement_units.icon_path,
	valid_measurement_units.created_on,
	valid_measurement_units.last_updated_on,
	valid_measurement_units.archived_on
FROM valid_measurement_units
WHERE valid_measurement_units.name ILIKE $1
AND valid_measurement_units.archived_on IS NULL
LIMIT 50`

// SearchForValidMeasurementUnits fetches a valid measurement unit from the database.
func (q *SQLQuerier) SearchForValidMeasurementUnits(ctx context.Context, query string) ([]*types.ValidMeasurementUnit, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if query == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachValidMeasurementUnitIDToSpan(span, query)

	args := []interface{}{
		wrapQueryForILIKE(query),
	}

	rows, err := q.performReadQuery(ctx, q.db, "valid measurement units", validMeasurementUnitSearchQuery, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing valid measurement units list retrieval query")
	}

	x, _, _, err := q.scanValidMeasurementUnits(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid measurement units")
	}

	return x, nil
}

const getTotalValidMeasurementUnitsCountQuery = "SELECT COUNT(valid_measurement_units.id) FROM valid_measurement_units WHERE valid_measurement_units.archived_on IS NULL"

// GetTotalValidMeasurementUnitCount fetches the count of valid measurement units from the database that meet a particular filter.
func (q *SQLQuerier) GetTotalValidMeasurementUnitCount(ctx context.Context) (uint64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	count, err := q.performCountQuery(ctx, q.db, getTotalValidMeasurementUnitsCountQuery, "fetching count of valid measurement units")
	if err != nil {
		return 0, observability.PrepareError(err, logger, span, "querying for count of valid measurement units")
	}

	return count, nil
}

// GetValidMeasurementUnits fetches a list of valid measurement units from the database that meet a particular filter.
func (q *SQLQuerier) GetValidMeasurementUnits(ctx context.Context, filter *types.QueryFilter) (x *types.ValidMeasurementUnitList, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.ValidMeasurementUnitList{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		x.Page, x.Limit = filter.Page, filter.Limit
	}

	query, args := q.buildListQuery(ctx, "valid_measurement_units", nil, nil, nil, householdOwnershipColumn, validMeasurementUnitsTableColumns, "", false, filter, true)

	rows, err := q.performReadQuery(ctx, q.db, "validMeasurementUnits", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "executing valid measurement units list retrieval query")
	}

	if x.ValidMeasurementUnits, x.FilteredCount, x.TotalCount, err = q.scanValidMeasurementUnits(ctx, rows, true); err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid measurement units")
	}

	return x, nil
}

func (q *SQLQuerier) buildGetValidMeasurementUnitsWithIDsQuery(ctx context.Context, limit uint8, ids []string) (query string, args []interface{}) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	withIDsWhere := squirrel.Eq{
		"valid_measurement_units.id":          ids,
		"valid_measurement_units.archived_on": nil,
	}

	subqueryBuilder := q.sqlBuilder.Select(validMeasurementUnitsTableColumns...).
		From("valid_measurement_units").
		Join(fmt.Sprintf("unnest('{%s}'::text[])", joinIDs(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))

	query, args, err := q.sqlBuilder.Select(validMeasurementUnitsTableColumns...).
		FromSelect(subqueryBuilder, "valid_measurement_units").
		Where(withIDsWhere).ToSql()

	q.logQueryBuildingError(span, err)

	return query, args
}

// GetValidMeasurementUnitsWithIDs fetches valid measurement units from the database within a given set of IDs.
func (q *SQLQuerier) GetValidMeasurementUnitsWithIDs(ctx context.Context, limit uint8, ids []string) ([]*types.ValidMeasurementUnit, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

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

	query, args := q.buildGetValidMeasurementUnitsWithIDsQuery(ctx, limit, ids)

	rows, err := q.performReadQuery(ctx, q.db, "valid measurement units with IDs", query, args)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "fetching valid measurement units from database")
	}

	validMeasurementUnits, _, _, err := q.scanValidMeasurementUnits(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "scanning valid measurement units")
	}

	return validMeasurementUnits, nil
}

const validMeasurementUnitCreationQuery = "INSERT INTO valid_measurement_units (id,name,description,volumetric,icon_path) VALUES ($1,$2,$3,$4,$5)"

// CreateValidMeasurementUnit creates a valid measurement unit in the database.
func (q *SQLQuerier) CreateValidMeasurementUnit(ctx context.Context, input *types.ValidMeasurementUnitDatabaseCreationInput) (*types.ValidMeasurementUnit, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidMeasurementUnitIDKey, input.ID)

	args := []interface{}{
		input.ID,
		input.Name,
		input.Description,
		input.Volumetric,
		input.IconPath,
	}

	// create the valid measurement unit.
	if err := q.performWriteQuery(ctx, q.db, "valid measurement unit creation", validMeasurementUnitCreationQuery, args); err != nil {
		return nil, observability.PrepareError(err, logger, span, "performing valid measurement unit creation query")
	}

	x := &types.ValidMeasurementUnit{
		ID:          input.ID,
		Name:        input.Name,
		Description: input.Description,
		Volumetric:  input.Volumetric,
		IconPath:    input.IconPath,
		CreatedOn:   q.currentTime(),
	}

	tracing.AttachValidMeasurementUnitIDToSpan(span, x.ID)
	logger.Info("valid measurement unit created")

	return x, nil
}

const updateValidMeasurementUnitQuery = `
UPDATE valid_measurement_units SET 
	name = $1,
	description = $2,
	volumetric = $3,
	icon_path = $4,
	last_updated_on = extract(epoch FROM NOW()) 
WHERE archived_on IS NULL AND id = $5
`

// UpdateValidMeasurementUnit updates a particular valid measurement unit.
func (q *SQLQuerier) UpdateValidMeasurementUnit(ctx context.Context, updated *types.ValidMeasurementUnit) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidMeasurementUnitIDKey, updated.ID)
	tracing.AttachValidMeasurementUnitIDToSpan(span, updated.ID)

	args := []interface{}{
		updated.Name,
		updated.Description,
		updated.Volumetric,
		updated.IconPath,
		updated.ID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid measurement unit update", updateValidMeasurementUnitQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating valid measurement unit")
	}

	logger.Info("valid measurement unit updated")

	return nil
}

const archiveValidMeasurementUnitQuery = "UPDATE valid_measurement_units SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1"

// ArchiveValidMeasurementUnit archives a valid measurement unit from the database by its ID.
func (q *SQLQuerier) ArchiveValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validMeasurementUnitID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachValidMeasurementUnitIDToSpan(span, validMeasurementUnitID)

	args := []interface{}{
		validMeasurementUnitID,
	}

	if err := q.performWriteQuery(ctx, q.db, "valid measurement unit archive", archiveValidMeasurementUnitQuery, args); err != nil {
		return observability.PrepareError(err, logger, span, "updating valid measurement unit")
	}

	logger.Info("valid measurement unit archived")

	return nil
}