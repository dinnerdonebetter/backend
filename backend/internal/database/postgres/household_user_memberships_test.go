package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestQuerier_Integration_HouseholdUserMemberships(t *testing.T) {
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
	households, err := dbc.GetHouseholds(ctx, exampleUser.ID, nil)
	assert.NoError(t, err)
	assert.Len(t, households.Data, 1)

	household := fakes.BuildFakeHousehold()
	household.BelongsToUser = exampleUser.ID
	exampleHousehold := createHouseholdForTest(t, ctx, household, dbc)

	memberUserIDs := []string{exampleUser.ID}

	for i := 0; i < exampleQuantity; i++ {
		newMember := createUserForTest(t, ctx, nil, dbc)
		assert.NoError(t, dbc.addUserToHousehold(ctx, dbc.db, &types.HouseholdUserMembershipDatabaseCreationInput{
			ID:            identifiers.New(),
			Reason:        "testing",
			UserID:        newMember.ID,
			HouseholdID:   exampleHousehold.ID,
			HouseholdRole: "household_member",
		}))
		memberUserIDs = append(memberUserIDs, newMember.ID)
	}

	household, err = dbc.GetHousehold(ctx, exampleHousehold.ID)
	assert.NoError(t, err)

	householdMemberUserIDs := []string{}
	for _, member := range household.Members {
		householdMemberUserIDs = append(householdMemberUserIDs, member.BelongsToUser.ID)
	}

	assert.Subset(t, memberUserIDs, householdMemberUserIDs)

	isMember, err := dbc.UserIsMemberOfHousehold(ctx, memberUserIDs[0], exampleHousehold.ID)
	assert.NoError(t, err)
	assert.True(t, isMember)

	assert.NoError(t, dbc.MarkHouseholdAsUserDefault(ctx, memberUserIDs[1], exampleHousehold.ID))
	defaultHouseholdID, err := dbc.GetDefaultHouseholdIDForUser(ctx, memberUserIDs[1])
	assert.NoError(t, err)
	assert.Equal(t, exampleHousehold.ID, defaultHouseholdID)

	sessionCtxData, err := dbc.BuildSessionContextDataForUser(ctx, memberUserIDs[1])
	assert.NoError(t, err)
	assert.NotNil(t, sessionCtxData)
	assert.Equal(t, exampleHousehold.ID, sessionCtxData.ActiveHouseholdID)

	assert.NoError(t, dbc.RemoveUserFromHousehold(ctx, memberUserIDs[len(memberUserIDs)-1], exampleHousehold.ID))

	assert.NoError(t, dbc.TransferHouseholdOwnership(ctx, exampleHousehold.ID, &types.HouseholdOwnershipTransferInput{
		Reason:       "testing",
		CurrentOwner: exampleUser.ID,
		NewOwner:     memberUserIDs[1],
	}))

	assert.NoError(t, dbc.ModifyUserPermissions(ctx, exampleHousehold.ID, memberUserIDs[0], &types.ModifyUserPermissionsInput{
		Reason:  "testing",
		NewRole: "household_admin",
	}))
}

func TestQuerier_BuildSessionContextDataForUser(T *testing.T) {
	T.Parallel()

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.BuildSessionContextDataForUser(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetDefaultHouseholdIDForUser(T *testing.T) {
	T.Parallel()

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetDefaultHouseholdIDForUser(ctx, "")
		assert.Error(t, err)
		assert.Zero(t, actual)
	})
}

func TestQuerier_MarkHouseholdAsUserDefault(T *testing.T) {
	T.Parallel()

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHousehold := fakes.BuildFakeHousehold()

		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkHouseholdAsUserDefault(ctx, "", exampleHousehold.ID))
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()

		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkHouseholdAsUserDefault(ctx, exampleUser.ID, ""))
	})
}

func TestQuerier_UserIsMemberOfHousehold(T *testing.T) {
	T.Parallel()

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.UserIsMemberOfHousehold(ctx, "", exampleHouseholdID)
		assert.False(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.UserIsMemberOfHousehold(ctx, exampleUserID, "")
		assert.False(t, actual)
		assert.Error(t, err)
	})
}

func TestQuerier_ModifyUserPermissions(T *testing.T) {
	T.Parallel()

	T.Run("with invalid household id", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeUserPermissionModificationInput()

		c, _ := buildTestClient(t)

		assert.Error(t, c.ModifyUserPermissions(ctx, "", exampleUserID, exampleInput))
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		assert.Error(t, c.ModifyUserPermissions(ctx, exampleHouseholdID, exampleUserID, nil))
	})
}

func TestSQLQuerier_addUserToHousehold(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, db := buildTestClient(t)

		err := c.addUserToHousehold(ctx, c.db, nil)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_RemoveUserFromHousehold(T *testing.T) {
	T.Parallel()

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHousehold := fakes.BuildFakeHousehold()

		c, _ := buildTestClient(t)

		assert.Error(t, c.RemoveUserFromHousehold(ctx, "", exampleHousehold.ID))
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()

		c, _ := buildTestClient(t)

		assert.Error(t, c.RemoveUserFromHousehold(ctx, exampleUser.ID, ""))
	})

	T.Run("with error creating transaction", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdID := fakes.BuildFakeID()

		c, db := buildTestClient(t)

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		assert.Error(t, c.RemoveUserFromHousehold(ctx, exampleUserID, exampleHouseholdID))
	})
}
