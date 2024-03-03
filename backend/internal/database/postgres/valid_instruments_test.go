package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createValidInstrumentForTest(t *testing.T, ctx context.Context, exampleValidInstrument *types.ValidInstrument, dbc *Querier) *types.ValidInstrument {
	t.Helper()

	// create
	if exampleValidInstrument == nil {
		exampleValidInstrument = fakes.BuildFakeValidInstrument()
	}
	dbInput := converters.ConvertValidInstrumentToValidInstrumentDatabaseCreationInput(exampleValidInstrument)

	created, err := dbc.CreateValidInstrument(ctx, dbInput)
	exampleValidInstrument.CreatedAt = created.CreatedAt
	assert.NoError(t, err)
	assert.Equal(t, exampleValidInstrument, created)

	validInstrument, err := dbc.GetValidInstrument(ctx, created.ID)
	exampleValidInstrument.CreatedAt = validInstrument.CreatedAt

	assert.NoError(t, err)
	assert.Equal(t, validInstrument, exampleValidInstrument)

	return created
}

func TestQuerier_Integration_ValidInstruments(t *testing.T) {
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

	exampleValidInstrument := fakes.BuildFakeValidInstrument()
	createdValidInstruments := []*types.ValidInstrument{}

	// create
	createdValidInstruments = append(createdValidInstruments, createValidInstrumentForTest(t, ctx, exampleValidInstrument, dbc))

	// update
	updatedValidInstrument := fakes.BuildFakeValidInstrument()
	updatedValidInstrument.ID = createdValidInstruments[0].ID
	assert.NoError(t, dbc.UpdateValidInstrument(ctx, updatedValidInstrument))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := fakes.BuildFakeValidInstrument()
		input.Name = fmt.Sprintf("%s %d", updatedValidInstrument.Name, i)
		createdValidInstruments = append(createdValidInstruments, createValidInstrumentForTest(t, ctx, input, dbc))
	}

	// fetch as list
	validInstruments, err := dbc.GetValidInstruments(ctx, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, validInstruments.Data)
	assert.Equal(t, len(createdValidInstruments), len(validInstruments.Data))

	// fetch as list of IDs
	validInstrumentIDs := []string{}
	for _, validInstrument := range createdValidInstruments {
		validInstrumentIDs = append(validInstrumentIDs, validInstrument.ID)
	}

	byIDs, err := dbc.GetValidInstrumentsWithIDs(ctx, validInstrumentIDs)
	assert.NoError(t, err)
	assert.Equal(t, validInstruments.Data, byIDs)

	// fetch via name search
	byName, err := dbc.SearchForValidInstruments(ctx, updatedValidInstrument.Name)
	assert.NoError(t, err)
	assert.Equal(t, validInstruments.Data, byName)

	random, err := dbc.GetRandomValidInstrument(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, random)

	needingIndexing, err := dbc.GetValidInstrumentIDsThatNeedSearchIndexing(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, needingIndexing)

	// delete
	for _, validInstrument := range createdValidInstruments {
		assert.NoError(t, dbc.MarkValidInstrumentAsIndexed(ctx, validInstrument.ID))
		assert.NoError(t, dbc.ArchiveValidInstrument(ctx, validInstrument.ID))

		var exists bool
		exists, err = dbc.ValidInstrumentExists(ctx, validInstrument.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.ValidInstrument
		y, err = dbc.GetValidInstrument(ctx, validInstrument.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}

func TestQuerier_ValidInstrumentExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidInstrumentExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidInstrument(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_SearchForValidInstruments(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.SearchForValidInstruments(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidInstrument(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidInstrument(ctx, nil))
	})
}

func TestQuerier_ArchiveValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidInstrument(ctx, ""))
	})
}

func TestQuerier_MarkValidInstrumentAsIndexed(T *testing.T) {
	T.Parallel()

	T.Run("with invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkValidInstrumentAsIndexed(ctx, ""))
	})
}
