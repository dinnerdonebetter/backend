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

func createHouseholdForTest(t *testing.T, ctx context.Context, exampleHousehold *types.Household, dbc *Querier) *types.Household {
	t.Helper()

	// create
	if exampleHousehold == nil {
		exampleUser := createUserForTest(t, ctx, nil, dbc)
		exampleHousehold = fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUser.ID
	}
	exampleHousehold.PaymentProcessorCustomerID = ""
	exampleHousehold.Members = nil
	dbInput := converters.ConvertHouseholdToHouseholdDatabaseCreationInput(exampleHousehold)

	created, err := dbc.CreateHousehold(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)
	exampleHousehold.CreatedAt = created.CreatedAt
	exampleHousehold.WebhookEncryptionKey = created.WebhookEncryptionKey
	assert.Equal(t, exampleHousehold, created)

	household, err := dbc.GetHousehold(ctx, created.ID)
	require.NoError(t, err)
	require.NotNil(t, household)

	exampleHousehold.CreatedAt = household.CreatedAt
	exampleHousehold.Members = household.Members
	exampleHousehold.WebhookEncryptionKey = household.WebhookEncryptionKey

	assert.Equal(t, exampleHousehold, household)

	return created
}

func TestQuerier_Integration_Households(t *testing.T) {
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

	exampleUser := createUserForTest(t, ctx, nil, dbc)

	exampleHousehold := fakes.BuildFakeHousehold()
	exampleHousehold.Members = nil
	exampleHousehold.BelongsToUser = exampleUser.ID
	exampleHousehold.PaymentProcessorCustomerID = ""
	createdHouseholds := []*types.Household{}

	// create
	createdHouseholds = append(createdHouseholds, createHouseholdForTest(t, ctx, exampleHousehold, dbc))

	// update
	updatedHousehold := fakes.BuildFakeHousehold()
	updatedHousehold.ID = createdHouseholds[0].ID
	updatedHousehold.BelongsToUser = createdHouseholds[0].BelongsToUser
	assert.NoError(t, dbc.UpdateHousehold(ctx, updatedHousehold))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := fakes.BuildFakeHousehold()
		input.BelongsToUser = exampleUser.ID
		input.Name = fmt.Sprintf("%s %d", updatedHousehold.Name, i)
		createdHouseholds = append(createdHouseholds, createHouseholdForTest(t, ctx, input, dbc))
	}

	// fetch as list
	households, err := dbc.GetHouseholds(ctx, exampleUser.ID, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, households.Data)
	assert.GreaterOrEqual(t, len(households.Data), len(createdHouseholds))

	// delete
	for _, household := range createdHouseholds {
		assert.NoError(t, dbc.ArchiveHousehold(ctx, household.ID, exampleUser.ID))

		var y *types.Household
		y, err = dbc.GetHousehold(ctx, household.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}

	assert.NoError(t, dbc.ArchiveUser(ctx, exampleUser.ID))
}

func TestQuerier_GetHousehold(T *testing.T) {
	T.Parallel()

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleUserID := fakes.BuildFakeID()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUserID

		c, _ := buildTestClient(t)

		actual, err := c.GetHousehold(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleUserID := fakes.BuildFakeID()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUserID

		c, _ := buildTestClient(t)

		actual, err := c.GetHousehold(ctx, exampleHousehold.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetHouseholds(T *testing.T) {
	T.Parallel()

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetHouseholds(ctx, "", filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateHousehold(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateHousehold(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateHousehold(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUserID

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateHousehold(ctx, nil))
	})
}

func TestQuerier_ArchiveHousehold(T *testing.T) {
	T.Parallel()

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveHousehold(ctx, "", exampleUserID))
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveHousehold(ctx, exampleHouseholdID, ""))
	})
}
