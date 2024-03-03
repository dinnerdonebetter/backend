package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createValidPreparationInstrumentForTest(t *testing.T, ctx context.Context, exampleValidPreparationInstrument *types.ValidPreparationInstrument, dbc *Querier) *types.ValidPreparationInstrument {
	t.Helper()

	// create
	if exampleValidPreparationInstrument == nil {
		exampleValidInstrument := createValidInstrumentForTest(t, ctx, nil, dbc)
		exampleValidPreparation := createValidPreparationForTest(t, ctx, nil, dbc)
		exampleValidPreparationInstrument = fakes.BuildFakeValidPreparationInstrument()
		exampleValidPreparationInstrument.Instrument = *exampleValidInstrument
		exampleValidPreparationInstrument.Preparation = *exampleValidPreparation
	}

	dbInput := converters.ConvertValidPreparationInstrumentToValidPreparationInstrumentDatabaseCreationInput(exampleValidPreparationInstrument)

	created, err := dbc.CreateValidPreparationInstrument(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)
	exampleValidPreparationInstrument.CreatedAt = created.CreatedAt
	exampleValidPreparationInstrument.Preparation = types.ValidPreparation{ID: exampleValidPreparationInstrument.Preparation.ID}
	exampleValidPreparationInstrument.Instrument = types.ValidInstrument{ID: exampleValidPreparationInstrument.Instrument.ID}
	assert.Equal(t, exampleValidPreparationInstrument, created)

	validPreparationInstrument, err := dbc.GetValidPreparationInstrument(ctx, created.ID)
	exampleValidPreparationInstrument.CreatedAt = validPreparationInstrument.CreatedAt
	exampleValidPreparationInstrument.Preparation = validPreparationInstrument.Preparation
	exampleValidPreparationInstrument.Instrument = validPreparationInstrument.Instrument

	assert.NoError(t, err)
	assert.Equal(t, validPreparationInstrument, exampleValidPreparationInstrument)

	return created
}

func TestQuerier_Integration_ValidPreparationInstruments(t *testing.T) {
	if !runningContainerTests {
		t.SkipNow()
	}

	ctx := context.Background()
	dbc, container := buildDatabaseClientForTest(t, ctx)

	databaseURI, err := container.ConnectionString(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, databaseURI)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	exampleValidInstrument := createValidInstrumentForTest(t, ctx, nil, dbc)
	exampleValidPreparation := createValidPreparationForTest(t, ctx, nil, dbc)
	exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
	exampleValidPreparationInstrument.Preparation = *exampleValidPreparation
	exampleValidPreparationInstrument.Instrument = *exampleValidInstrument
	createdValidPreparationInstruments := []*types.ValidPreparationInstrument{}

	// create
	createdValidPreparationInstruments = append(createdValidPreparationInstruments, createValidPreparationInstrumentForTest(t, ctx, exampleValidPreparationInstrument, dbc))

	// update
	updatedValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
	updatedValidPreparationInstrument.ID = createdValidPreparationInstruments[0].ID
	updatedValidPreparationInstrument.Preparation = createdValidPreparationInstruments[0].Preparation
	updatedValidPreparationInstrument.Instrument = createdValidPreparationInstruments[0].Instrument
	assert.NoError(t, dbc.UpdateValidPreparationInstrument(ctx, updatedValidPreparationInstrument))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := fakes.BuildFakeValidPreparationInstrument()
		input.Preparation = createdValidPreparationInstruments[0].Preparation
		input.Instrument = createdValidPreparationInstruments[0].Instrument
		createdValidPreparationInstruments = append(createdValidPreparationInstruments, createValidPreparationInstrumentForTest(t, ctx, input, dbc))
	}

	// fetch as list
	validPreparationInstruments, err := dbc.GetValidPreparationInstruments(ctx, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, validPreparationInstruments.Data)
	assert.Equal(t, len(createdValidPreparationInstruments), len(validPreparationInstruments.Data))

	forPreparation, err := dbc.GetValidPreparationInstrumentsForPreparation(ctx, createdValidPreparationInstruments[0].Preparation.ID, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, forPreparation.Data)

	forInstrument, err := dbc.GetValidPreparationInstrumentsForInstrument(ctx, createdValidPreparationInstruments[0].Instrument.ID, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, forInstrument.Data)

	// delete
	for _, validPreparationInstrument := range createdValidPreparationInstruments {
		assert.NoError(t, dbc.ArchiveValidPreparationInstrument(ctx, validPreparationInstrument.ID))

		var exists bool
		exists, err = dbc.ValidPreparationInstrumentExists(ctx, validPreparationInstrument.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.ValidPreparationInstrument
		y, err = dbc.GetValidPreparationInstrument(ctx, validPreparationInstrument.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}

func TestQuerier_ValidPreparationInstrumentExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid preparation instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidPreparationInstrumentExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid preparation instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidPreparationInstrument(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidPreparationInstrument(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidPreparationInstrument(ctx, nil))
	})
}

func TestQuerier_ArchiveValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid preparation instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidPreparationInstrument(ctx, ""))
	})
}
