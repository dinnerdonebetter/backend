package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createValidIngredientGroupForTest(t *testing.T, ctx context.Context, exampleValidIngredientGroup *types.ValidIngredientGroup, dbc *Querier) *types.ValidIngredientGroup {
	t.Helper()

	// create
	if exampleValidIngredientGroup == nil {
		exampleValidIngredientGroup = fakes.BuildFakeValidIngredientGroup()
	}
	dbInput := converters.ConvertValidIngredientGroupToValidIngredientGroupDatabaseCreationInput(exampleValidIngredientGroup)

	created, err := dbc.CreateValidIngredientGroup(ctx, dbInput)
	require.NoError(t, err)
	require.NotNil(t, created)
	exampleValidIngredientGroup.CreatedAt = created.CreatedAt
	for i := range exampleValidIngredientGroup.Members {
		exampleValidIngredientGroup.Members[i].CreatedAt = created.Members[i].CreatedAt
		exampleValidIngredientGroup.Members[i].ValidIngredient = created.Members[i].ValidIngredient
	}
	assert.Equal(t, exampleValidIngredientGroup, created)

	validIngredientGroup, err := dbc.GetValidIngredientGroup(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, validIngredientGroup)

	exampleValidIngredientGroup.CreatedAt = validIngredientGroup.CreatedAt
	for i := range exampleValidIngredientGroup.Members {
		exampleValidIngredientGroup.Members[i].CreatedAt = validIngredientGroup.Members[i].CreatedAt
		exampleValidIngredientGroup.Members[i].ValidIngredient = validIngredientGroup.Members[i].ValidIngredient
	}

	assert.NoError(t, err)
	assert.Equal(t, validIngredientGroup, exampleValidIngredientGroup)

	return created
}

func TestQuerier_Integration_ValidIngredientGroups(t *testing.T) {
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

	exampleValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()
	exampleValidIngredientGroup.Members = []*types.ValidIngredientGroupMember{}
	for i := 0; i < exampleQuantity; i++ {
		exampleValidIngredientGroup.Members = append(exampleValidIngredientGroup.Members, &types.ValidIngredientGroupMember{
			ID:              identifiers.New(),
			BelongsToGroup:  exampleValidIngredientGroup.ID,
			ValidIngredient: *createValidIngredientForTest(t, ctx, nil, dbc),
		})
	}
	createdValidIngredientGroups := []*types.ValidIngredientGroup{}

	// create
	createdValidIngredientGroups = append(createdValidIngredientGroups, createValidIngredientGroupForTest(t, ctx, exampleValidIngredientGroup, dbc))

	// update
	updatedValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()
	updatedValidIngredientGroup.ID = createdValidIngredientGroups[0].ID
	assert.NoError(t, dbc.UpdateValidIngredientGroup(ctx, updatedValidIngredientGroup))

	// fetch as list
	validIngredientGroups, err := dbc.GetValidIngredientGroups(ctx, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, validIngredientGroups.Data)
	assert.Equal(t, len(createdValidIngredientGroups), len(validIngredientGroups.Data))

	// fetch via name search
	byName, err := dbc.SearchForValidIngredientGroups(ctx, updatedValidIngredientGroup.Name, nil)
	assert.NoError(t, err)
	assert.Equal(t, validIngredientGroups.Data, byName)

	// delete
	for _, validIngredientGroup := range createdValidIngredientGroups {
		assert.NoError(t, dbc.ArchiveValidIngredientGroup(ctx, validIngredientGroup.ID))

		var exists bool
		exists, err = dbc.ValidIngredientGroupExists(ctx, validIngredientGroup.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.ValidIngredientGroup
		y, err = dbc.GetValidIngredientGroup(ctx, validIngredientGroup.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}

func TestQuerier_ValidIngredientGroupExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient group ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidIngredientGroupExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetValidIngredientGroup(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient group ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidIngredientGroup(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_SearchForValidIngredientGroups(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient group ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)
		filter := types.DefaultQueryFilter()

		actual, err := c.SearchForValidIngredientGroups(ctx, "", filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateValidIngredientGroup(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidIngredientGroup(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateValidIngredientGroup(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidIngredientGroup(ctx, nil))
	})
}

func TestQuerier_ArchiveValidIngredientGroup(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient group ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidIngredientGroup(ctx, ""))
	})
}
