package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createHouseholdInvitationForTest(t *testing.T, ctx context.Context, exampleHouseholdInvitation *types.HouseholdInvitation, dbc *Querier) *types.HouseholdInvitation {
	t.Helper()

	// create
	if exampleHouseholdInvitation == nil {
		fromUser := createUserForTest(t, ctx, nil, dbc)
		toUser := createUserForTest(t, ctx, nil, dbc)
		household := createHouseholdForTest(t, ctx, nil, dbc)
		exampleHouseholdInvitation = fakes.BuildFakeHouseholdInvitation()
		exampleHouseholdInvitation.ExpiresAt = time.Now().Add(time.Hour * 24 * 7)
		exampleHouseholdInvitation.DestinationHousehold = *household
		exampleHouseholdInvitation.FromUser = *fromUser
		exampleHouseholdInvitation.ToUser = &toUser.ID
	}
	dbInput := converters.ConvertHouseholdInvitationToHouseholdInvitationDatabaseCreationInput(exampleHouseholdInvitation)

	created, err := dbc.CreateHouseholdInvitation(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)
	exampleHouseholdInvitation.CreatedAt = created.CreatedAt
	exampleHouseholdInvitation.StatusNote = created.StatusNote
	exampleHouseholdInvitation.FromUser = created.FromUser
	assert.Equal(t, exampleHouseholdInvitation.DestinationHousehold.ID, created.DestinationHousehold.ID)
	exampleHouseholdInvitation.DestinationHousehold = created.DestinationHousehold
	assert.Equal(t, exampleHouseholdInvitation, created)

	householdInvitation, err := dbc.GetHouseholdInvitationByHouseholdAndID(ctx, created.DestinationHousehold.ID, created.ID)
	assert.NoError(t, err)
	require.NotNil(t, householdInvitation)
	exampleHouseholdInvitation.CreatedAt = householdInvitation.CreatedAt
	exampleHouseholdInvitation.ExpiresAt = householdInvitation.ExpiresAt
	assert.Equal(t, exampleHouseholdInvitation.FromUser.ID, householdInvitation.FromUser.ID)
	exampleHouseholdInvitation.FromUser = householdInvitation.FromUser
	assert.Equal(t, exampleHouseholdInvitation.DestinationHousehold.ID, householdInvitation.DestinationHousehold.ID)
	exampleHouseholdInvitation.DestinationHousehold = householdInvitation.DestinationHousehold

	assert.Equal(t, householdInvitation, exampleHouseholdInvitation)

	return created
}

func TestQuerier_Integration_HouseholdInvitations(t *testing.T) {
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

	household := createHouseholdForTest(t, ctx, nil, dbc)

	fromUser := createUserForTest(t, ctx, nil, dbc)
	toUserA := createUserForTest(t, ctx, nil, dbc)
	toUserB := createUserForTest(t, ctx, nil, dbc)
	toUserC := createUserForTest(t, ctx, nil, dbc)

	toBeCancelledInput := fakes.BuildFakeHouseholdInvitation()
	toBeCancelledInput.DestinationHousehold = *household
	toBeCancelledInput.ExpiresAt = time.Now().Add(time.Hour * 24 * 7)
	toBeCancelledInput.FromUser = *fromUser
	toBeCancelledInput.ToUser = &toUserA.ID
	toBeCancelled := createHouseholdInvitationForTest(t, ctx, toBeCancelledInput, dbc)

	toBeRejectedInput := fakes.BuildFakeHouseholdInvitation()
	toBeRejectedInput.DestinationHousehold = *household
	toBeRejectedInput.ExpiresAt = time.Now().Add(time.Hour * 24 * 7)
	toBeRejectedInput.FromUser = *fromUser
	toBeRejectedInput.ToUser = &toUserB.ID
	toBeRejected := createHouseholdInvitationForTest(t, ctx, toBeRejectedInput, dbc)

	toBeAcceptedInput := fakes.BuildFakeHouseholdInvitation()
	toBeAcceptedInput.DestinationHousehold = *household
	toBeAcceptedInput.ExpiresAt = time.Now().Add(time.Hour * 24 * 7)
	toBeAcceptedInput.FromUser = *fromUser
	toBeAcceptedInput.ToUser = &toUserC.ID
	toBeAcceptedInput.ToEmail = toUserC.EmailAddress
	toBeAccepted := createHouseholdInvitationForTest(t, ctx, toBeAcceptedInput, dbc)

	outboundInvites, err := dbc.GetPendingHouseholdInvitationsFromUser(ctx, fromUser.ID, nil)
	assert.NoError(t, err)
	assert.Len(t, outboundInvites.Data, 3)

	inboundInvites, err := dbc.GetPendingHouseholdInvitationsForUser(ctx, toUserC.ID, nil)
	assert.NoError(t, err)
	assert.Len(t, inboundInvites.Data, 1)

	exists, err := dbc.HouseholdInvitationExists(ctx, toBeCancelled.ID)
	assert.NoError(t, err)
	assert.True(t, exists)

	invite, err := dbc.GetHouseholdInvitationByEmailAndToken(ctx, toUserC.EmailAddress, toBeAccepted.Token)
	assert.NoError(t, err)
	assert.NotNil(t, invite)

	// create invite for nonexistent user
	forNewUserInput := fakes.BuildFakeHouseholdInvitation()
	forNewUserInput.DestinationHousehold = *household
	forNewUserInput.ExpiresAt = time.Now().Add(time.Hour * 24 * 7)
	forNewUserInput.FromUser = *fromUser
	forNewUserInput.ToUser = nil
	forNewUserInput.ToEmail = fakes.BuildFakeUser().EmailAddress
	forNewUser := createHouseholdInvitationForTest(t, ctx, forNewUserInput, dbc)

	fakeUser := fakes.BuildFakeUser()
	fakeUser.EmailAddress = forNewUserInput.ToEmail
	dbInput := converters.ConvertUserToUserDatabaseCreationInput(fakeUser)
	dbInput.InvitationToken = forNewUser.Token
	dbInput.DestinationHouseholdID = household.ID

	createdUser, err := dbc.CreateUser(ctx, dbInput)
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)

	assert.NoError(t, dbc.CancelHouseholdInvitation(ctx, toBeCancelled.ID, "testing"))
	assert.NoError(t, dbc.RejectHouseholdInvitation(ctx, toBeRejected.ID, "testing"))
	assert.NoError(t, dbc.AcceptHouseholdInvitation(ctx, toBeAccepted.ID, toBeAccepted.Token, "testing"))
}
