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

func createHouseholdInstrumentOwnershipForTest(t *testing.T, ctx context.Context, exampleHouseholdInstrumentOwnership *types.HouseholdInstrumentOwnership, dbc *Querier) *types.HouseholdInstrumentOwnership {
	t.Helper()

	// create
	if exampleHouseholdInstrumentOwnership == nil {
		exampleHouseholdInstrumentOwnership = fakes.BuildFakeHouseholdInstrumentOwnership()
	}
	dbInput := converters.ConvertHouseholdInstrumentOwnershipToHouseholdInstrumentOwnershipDatabaseCreationInput(exampleHouseholdInstrumentOwnership)

	created, err := dbc.CreateHouseholdInstrumentOwnership(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)

	exampleHouseholdInstrumentOwnership.CreatedAt = created.CreatedAt
	assert.Equal(t, exampleHouseholdInstrumentOwnership.Instrument.ID, created.Instrument.ID)
	exampleHouseholdInstrumentOwnership.Instrument = created.Instrument
	assert.Equal(t, exampleHouseholdInstrumentOwnership, created)

	householdInstrumentOwnership, err := dbc.GetHouseholdInstrumentOwnership(ctx, created.ID, created.BelongsToHousehold)
	exampleHouseholdInstrumentOwnership.CreatedAt = householdInstrumentOwnership.CreatedAt
	assert.Equal(t, exampleHouseholdInstrumentOwnership.Instrument.ID, householdInstrumentOwnership.Instrument.ID)
	exampleHouseholdInstrumentOwnership.Instrument = householdInstrumentOwnership.Instrument

	assert.NoError(t, err)
	assert.Equal(t, householdInstrumentOwnership, exampleHouseholdInstrumentOwnership)

	return created
}

func TestQuerier_Integration_HouseholdInstrumentOwnerships(t *testing.T) {
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

	user := createUserForTest(t, ctx, nil, dbc)
	householdID, err := dbc.GetDefaultHouseholdIDForUser(ctx, user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, householdID)

	instrument := createValidInstrumentForTest(t, ctx, nil, dbc)

	exampleHouseholdInstrumentOwnership := fakes.BuildFakeHouseholdInstrumentOwnership()
	exampleHouseholdInstrumentOwnership.BelongsToHousehold = householdID
	exampleHouseholdInstrumentOwnership.Instrument = *instrument
	createdHouseholdInstrumentOwnerships := []*types.HouseholdInstrumentOwnership{}

	// create
	createdHouseholdInstrumentOwnerships = append(createdHouseholdInstrumentOwnerships, createHouseholdInstrumentOwnershipForTest(t, ctx, exampleHouseholdInstrumentOwnership, dbc))

	// update
	assert.NoError(t, dbc.UpdateHouseholdInstrumentOwnership(ctx, createdHouseholdInstrumentOwnerships[0]))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		newInstrument := createValidInstrumentForTest(t, ctx, nil, dbc)
		input := fakes.BuildFakeHouseholdInstrumentOwnership()
		input.BelongsToHousehold = householdID
		input.Instrument = *newInstrument
		createdHouseholdInstrumentOwnerships = append(createdHouseholdInstrumentOwnerships, createHouseholdInstrumentOwnershipForTest(t, ctx, input, dbc))
	}

	// fetch as list
	householdInstrumentOwnerships, err := dbc.GetHouseholdInstrumentOwnerships(ctx, householdID, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, householdInstrumentOwnerships.Data)
	assert.Equal(t, len(createdHouseholdInstrumentOwnerships), len(householdInstrumentOwnerships.Data))

	// delete
	for _, householdInstrumentOwnership := range createdHouseholdInstrumentOwnerships {
		assert.NoError(t, dbc.ArchiveHouseholdInstrumentOwnership(ctx, householdInstrumentOwnership.ID, householdID))

		var exists bool
		exists, err = dbc.HouseholdInstrumentOwnershipExists(ctx, householdInstrumentOwnership.ID, householdID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.HouseholdInstrumentOwnership
		y, err = dbc.GetHouseholdInstrumentOwnership(ctx, householdInstrumentOwnership.ID, householdID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}

func TestQuerier_HouseholdInstrumentOwnershipExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid household instrument ownership ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.HouseholdInstrumentOwnershipExists(ctx, "", exampleHouseholdID)
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetHouseholdInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("with invalid household instrument ownership ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetHouseholdInstrumentOwnership(ctx, "", exampleHouseholdID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateHouseholdInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateHouseholdInstrumentOwnership(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateHouseholdInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateHouseholdInstrumentOwnership(ctx, nil))
	})
}

func TestQuerier_ArchiveHouseholdInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("with invalid household instrument ownership ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveHouseholdInstrumentOwnership(ctx, "", exampleHouseholdID))
	})
}
