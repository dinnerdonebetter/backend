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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createValidPreparationForTest(t *testing.T, ctx context.Context, exampleValidPreparation *types.ValidPreparation, dbc *Querier) *types.ValidPreparation {
	t.Helper()

	// create
	if exampleValidPreparation == nil {
		exampleValidPreparation = fakes.BuildFakeValidPreparation()
	}
	dbInput := converters.ConvertValidPreparationToValidPreparationDatabaseCreationInput(exampleValidPreparation)

	created, err := dbc.CreateValidPreparation(ctx, dbInput)
	exampleValidPreparation.CreatedAt = created.CreatedAt
	assert.NoError(t, err)
	assert.Equal(t, exampleValidPreparation, created)

	validPreparation, err := dbc.GetValidPreparation(ctx, created.ID)
	exampleValidPreparation.CreatedAt = validPreparation.CreatedAt

	assert.NoError(t, err)
	assert.Equal(t, validPreparation, exampleValidPreparation)

	return created
}

func TestQuerier_Integration_ValidPreparations(t *testing.T) {
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

	exampleValidPreparation := fakes.BuildFakeValidPreparation()
	createdValidPreparations := []*types.ValidPreparation{}

	// create
	createdValidPreparations = append(createdValidPreparations, createValidPreparationForTest(t, ctx, exampleValidPreparation, dbc))

	// update
	updatedValidPreparation := fakes.BuildFakeValidPreparation()
	updatedValidPreparation.ID = createdValidPreparations[0].ID
	assert.NoError(t, dbc.UpdateValidPreparation(ctx, updatedValidPreparation))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := fakes.BuildFakeValidPreparation()
		input.Name = fmt.Sprintf("%s %d", updatedValidPreparation.Name, i)
		createdValidPreparations = append(createdValidPreparations, createValidPreparationForTest(t, ctx, input, dbc))
	}

	// fetch as list
	validPreparations, err := dbc.GetValidPreparations(ctx, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, validPreparations.Data)
	assert.Equal(t, len(createdValidPreparations), len(validPreparations.Data))

	// fetch as list of IDs
	validPreparationIDs := []string{}
	for _, validPreparation := range createdValidPreparations {
		validPreparationIDs = append(validPreparationIDs, validPreparation.ID)
	}

	byIDs, err := dbc.GetValidPreparationsWithIDs(ctx, validPreparationIDs)
	assert.NoError(t, err)
	assert.Equal(t, validPreparations.Data, byIDs)

	// fetch via name search
	byName, err := dbc.SearchForValidPreparations(ctx, updatedValidPreparation.Name)
	assert.NoError(t, err)
	assert.Equal(t, validPreparations.Data, byName)

	whatever, err := dbc.GetValidPreparationIDsThatNeedSearchIndexing(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, whatever)

	assert.NoError(t, dbc.MarkValidPreparationAsIndexed(ctx, updatedValidPreparation.ID))

	randomPreparation, err := dbc.GetRandomValidPreparation(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, randomPreparation)

	// delete
	for _, validPreparation := range createdValidPreparations {
		assert.NoError(t, dbc.ArchiveValidPreparation(ctx, validPreparation.ID))

		var exists bool
		exists, err = dbc.ValidPreparationExists(ctx, validPreparation.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.ValidPreparation
		y, err = dbc.GetValidPreparation(ctx, validPreparation.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}

func TestQuerier_ValidPreparationExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidPreparationExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidPreparation(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_SearchForValidPreparations(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.SearchForValidPreparations(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetValidPreparationsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("with nil IDs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		actual, err := c.GetValidPreparationsWithIDs(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidPreparation(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidPreparation(ctx, nil))
	})
}

func TestQuerier_ArchiveValidPreparation(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidPreparation(ctx, ""))
	})
}

func TestQuerier_MarkValidPreparationAsIndexed(T *testing.T) {
	T.Parallel()

	T.Run("with invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkValidPreparationAsIndexed(ctx, ""))
	})
}
