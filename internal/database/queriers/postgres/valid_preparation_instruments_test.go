package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func buildMockRowsFromValidPreparationInstruments(includeCounts bool, filteredCount uint64, validPreparationInstruments ...*types.ValidPreparationInstrument) *sqlmock.Rows {
	columns := fullValidPreparationInstrumentsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range validPreparationInstruments {
		rowValues := []driver.Value{
			x.ID,
			x.Notes,
			x.Preparation.ID,
			x.Preparation.Name,
			x.Preparation.Description,
			x.Preparation.IconPath,
			x.Preparation.CreatedOn,
			x.Preparation.LastUpdatedOn,
			x.Preparation.ArchivedOn,
			x.Instrument.ID,
			x.Instrument.Name,
			x.Instrument.Variant,
			x.Instrument.Description,
			x.Instrument.IconPath,
			x.Instrument.CreatedOn,
			x.Instrument.LastUpdatedOn,
			x.Instrument.ArchivedOn,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(validPreparationInstruments))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanValidPreparationInstruments(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanValidPreparationInstruments(ctx, mockRows, false)
		assert.Error(t, err)
	})

	T.Run("logs row closing errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, _, _, err := q.scanValidPreparationInstruments(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_ValidPreparationInstrumentExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleValidPreparationInstrument.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validPreparationInstrumentExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.ValidPreparationInstrumentExists(ctx, exampleValidPreparationInstrument.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid preparation instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidPreparationInstrumentExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleValidPreparationInstrument.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validPreparationInstrumentExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.ValidPreparationInstrumentExists(ctx, exampleValidPreparationInstrument.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleValidPreparationInstrument.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validPreparationInstrumentExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.ValidPreparationInstrumentExists(ctx, exampleValidPreparationInstrument.ID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleValidPreparationInstrument.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidPreparationInstrumentQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidPreparationInstruments(false, 0, exampleValidPreparationInstrument))

		actual, err := c.GetValidPreparationInstrument(ctx, exampleValidPreparationInstrument.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparationInstrument, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid preparation instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidPreparationInstrument(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleValidPreparationInstrument.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidPreparationInstrumentQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidPreparationInstrument(ctx, exampleValidPreparationInstrument.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetTotalValidPreparationInstrumentCount(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleCount := uint64(123)

		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getTotalValidPreparationInstrumentsCountQuery)).
			WithArgs().
			WillReturnRows(newCountDBRowResponse(uint64(123)))

		actual, err := c.GetTotalValidPreparationInstrumentCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getTotalValidPreparationInstrumentsCountQuery)).
			WithArgs().
			WillReturnError(errors.New("blah"))

		actual, err := c.GetTotalValidPreparationInstrumentCount(ctx)
		assert.Error(t, err)
		assert.Zero(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidPreparationInstruments(T *testing.T) {
	T.Parallel()

	joins := []string{
		validInstrumentsOnValidPreparationInstrumentsJoinClause,
		validPreparationsOnValidPreparationInstrumentsJoinClause,
	}

	groupBys := []string{
		"valid_preparations.id",
		"valid_instruments.id",
		"valid_preparation_instruments.id",
	}

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleValidPreparationInstrumentList := fakes.BuildFakeValidPreparationInstrumentList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_preparation_instruments", joins, groupBys, nil, householdOwnershipColumn, fullValidPreparationInstrumentsTableColumns, "", false, filter, true)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidPreparationInstruments(true, exampleValidPreparationInstrumentList.FilteredCount, exampleValidPreparationInstrumentList.ValidPreparationInstruments...))

		actual, err := c.GetValidPreparationInstruments(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparationInstrumentList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := (*types.QueryFilter)(nil)
		exampleValidPreparationInstrumentList := fakes.BuildFakeValidPreparationInstrumentList()
		exampleValidPreparationInstrumentList.Page = 0
		exampleValidPreparationInstrumentList.Limit = 0

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_preparation_instruments", joins, groupBys, nil, householdOwnershipColumn, fullValidPreparationInstrumentsTableColumns, "", false, filter, true)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidPreparationInstruments(true, exampleValidPreparationInstrumentList.FilteredCount, exampleValidPreparationInstrumentList.ValidPreparationInstruments...))

		actual, err := c.GetValidPreparationInstruments(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparationInstrumentList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_preparation_instruments", joins, groupBys, nil, householdOwnershipColumn, fullValidPreparationInstrumentsTableColumns, "", false, filter, true)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidPreparationInstruments(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_preparation_instruments", joins, groupBys, nil, householdOwnershipColumn, fullValidPreparationInstrumentsTableColumns, "", false, filter, true)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetValidPreparationInstruments(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		exampleValidPreparationInstrument.ID = "1"
		exampleValidPreparationInstrument.Preparation = types.ValidPreparation{ID: exampleValidPreparationInstrument.Preparation.ID}
		exampleValidPreparationInstrument.Instrument = types.ValidInstrument{ID: exampleValidPreparationInstrument.Instrument.ID}

		exampleInput := fakes.BuildFakeValidPreparationInstrumentDatabaseCreationInputFromValidPreparationInstrument(exampleValidPreparationInstrument)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.Notes,
			exampleInput.ValidPreparationID,
			exampleInput.ValidInstrumentID,
		}

		db.ExpectExec(formatQueryForSQLMock(validPreparationInstrumentCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() uint64 {
			return exampleValidPreparationInstrument.CreatedOn
		}

		actual, err := c.CreateValidPreparationInstrument(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparationInstrument, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidPreparationInstrument(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New(t.Name())
		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		exampleInput := fakes.BuildFakeValidPreparationInstrumentDatabaseCreationInputFromValidPreparationInstrument(exampleValidPreparationInstrument)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.Notes,
			exampleInput.ValidPreparationID,
			exampleInput.ValidInstrumentID,
		}

		db.ExpectExec(formatQueryForSQLMock(validPreparationInstrumentCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() uint64 {
			return exampleValidPreparationInstrument.CreatedOn
		}

		actual, err := c.CreateValidPreparationInstrument(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleValidPreparationInstrument.Notes,
			exampleValidPreparationInstrument.Preparation.ID,
			exampleValidPreparationInstrument.Instrument.ID,
			exampleValidPreparationInstrument.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidPreparationInstrumentQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateValidPreparationInstrument(ctx, exampleValidPreparationInstrument))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidPreparationInstrument(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleValidPreparationInstrument.Notes,
			exampleValidPreparationInstrument.Preparation.ID,
			exampleValidPreparationInstrument.Instrument.ID,
			exampleValidPreparationInstrument.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidPreparationInstrumentQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateValidPreparationInstrument(ctx, exampleValidPreparationInstrument))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleValidPreparationInstrument.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveValidPreparationInstrumentQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.ArchiveValidPreparationInstrument(ctx, exampleValidPreparationInstrument.ID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid preparation instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidPreparationInstrument(ctx, ""))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleValidPreparationInstrument.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveValidPreparationInstrumentQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveValidPreparationInstrument(ctx, exampleValidPreparationInstrument.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}