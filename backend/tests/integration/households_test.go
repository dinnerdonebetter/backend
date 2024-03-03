package integration

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	testutils "github.com/dinnerdonebetter/backend/tests/utils"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkHouseholdEquality(t *testing.T, expected, actual *types.Household) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for household %s to be %v, but it was %v ", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.AddressLine1, actual.AddressLine1, "expected AddressLine1 for household %s to be %v, but it was %v", expected.ID, expected.AddressLine1, actual.AddressLine1)
	assert.Equal(t, expected.AddressLine2, actual.AddressLine2, "expected AddressLine2 for household %s to be %v, but it was %v", expected.ID, expected.AddressLine2, actual.AddressLine2)
	assert.Equal(t, expected.City, actual.City, "expected City for household %s to be %v, but it was %v", expected.ID, expected.City, actual.City)
	assert.Equal(t, expected.State, actual.State, "expected State for household %s to be %v, but it was %v", expected.ID, expected.State, actual.State)
	assert.Equal(t, expected.ZipCode, actual.ZipCode, "expected ZipCode for household %s to be %v, but it was %v", expected.ID, expected.ZipCode, actual.ZipCode)
	assert.Equal(t, expected.Country, actual.Country, "expected Country for household %s to be %v, but it was %v", expected.ID, expected.Country, actual.Country)
	assert.Equal(t, expected.Latitude, actual.Latitude, "expected Latitude for household %s to be %v, but it was %v", expected.ID, expected.Latitude, actual.Latitude)
	assert.Equal(t, expected.Longitude, actual.Longitude, "expected Longitude for household %s to be %v, but it was %v", expected.ID, expected.Longitude, actual.Longitude)
	assert.NotZero(t, actual.CreatedAt)
}

func (s *TestSuite) TestHouseholds_Creating() {
	s.runForEachClient("should be possible to create households", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create household.
			exampleHousehold := fakes.BuildFakeHousehold()
			exampleHouseholdInput := converters.ConvertHouseholdToHouseholdCreationRequestInput(exampleHousehold)
			createdHousehold, err := testClients.user.CreateHousehold(ctx, exampleHouseholdInput)
			requireNotNilAndNoProblems(t, createdHousehold, err)

			// Assert household equality.
			checkHouseholdEquality(t, exampleHousehold, createdHousehold)

			// Clean up.
			assert.NoError(t, testClients.user.ArchiveHousehold(ctx, createdHousehold.ID))
		}
	})
}

func (s *TestSuite) TestHouseholds_Listing() {
	s.runForEachClient("should be possible to list households", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create households.
			var expected []*types.Household
			for i := 0; i < 5; i++ {
				// Create household.
				exampleHousehold := fakes.BuildFakeHousehold()
				exampleHouseholdInput := converters.ConvertHouseholdToHouseholdCreationRequestInput(exampleHousehold)
				createdHousehold, householdCreationErr := testClients.user.CreateHousehold(ctx, exampleHouseholdInput)
				requireNotNilAndNoProblems(t, createdHousehold, householdCreationErr)

				expected = append(expected, createdHousehold)
			}

			// Assert household list equality.
			actual, err := testClients.user.GetHouseholds(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			// Clean up.
			for _, createdHousehold := range actual.Data {
				assert.NoError(t, testClients.user.ArchiveHousehold(ctx, createdHousehold.ID))
			}
		}
	})
}

func (s *TestSuite) TestHouseholds_Reading_Returns404ForNonexistentHousehold() {
	s.runForEachClient("should not be possible to read a non-existent household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Attempt to fetch nonexistent household.
			_, err := testClients.user.GetHousehold(ctx, nonexistentID)
			assert.Error(t, err)
		}
	})
}

func (s *TestSuite) TestHouseholds_Reading() {
	s.runForEachClient("should be possible to read a household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create household.
			exampleHousehold := fakes.BuildFakeHousehold()
			exampleHouseholdInput := converters.ConvertHouseholdToHouseholdCreationRequestInput(exampleHousehold)
			createdHousehold, err := testClients.user.CreateHousehold(ctx, exampleHouseholdInput)
			requireNotNilAndNoProblems(t, createdHousehold, err)

			// Fetch household.
			actual, err := testClients.user.GetHousehold(ctx, createdHousehold.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// Assert household equality.
			checkHouseholdEquality(t, exampleHousehold, actual)

			// Clean up household.
			assert.NoError(t, testClients.user.ArchiveHousehold(ctx, createdHousehold.ID))
		}
	})
}

func (s *TestSuite) TestHouseholds_Updating_Returns404ForNonexistentHousehold() {
	s.runForEachClient("should not be possible to update a non-existent household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			exampleHousehold := fakes.BuildFakeHousehold()
			exampleHousehold.ID = nonexistentID

			assert.Error(t, testClients.user.UpdateHousehold(ctx, exampleHousehold))
		}
	})
}

func (s *TestSuite) TestHouseholds_Updating() {
	s.runForEachClient("should be possible to update a household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create household.
			exampleHousehold := fakes.BuildFakeHousehold()
			exampleHouseholdInput := converters.ConvertHouseholdToHouseholdCreationRequestInput(exampleHousehold)
			createdHousehold, err := testClients.user.CreateHousehold(ctx, exampleHouseholdInput)
			requireNotNilAndNoProblems(t, createdHousehold, err)

			// Change household.
			createdHousehold.Update(converters.ConvertHouseholdToHouseholdUpdateRequestInput(exampleHousehold))
			assert.NoError(t, testClients.user.UpdateHousehold(ctx, createdHousehold))

			// Fetch household.
			actual, err := testClients.user.GetHousehold(ctx, createdHousehold.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// Assert household equality.
			checkHouseholdEquality(t, exampleHousehold, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			// Clean up household.
			assert.NoError(t, testClients.user.ArchiveHousehold(ctx, createdHousehold.ID))
		}
	})
}

func (s *TestSuite) TestHouseholds_Archiving() {
	s.runForEachClient("should be possible to archive a household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create household.
			exampleHousehold := fakes.BuildFakeHousehold()
			exampleHouseholdInput := converters.ConvertHouseholdToHouseholdCreationRequestInput(exampleHousehold)
			createdHousehold, err := testClients.user.CreateHousehold(ctx, exampleHouseholdInput)
			requireNotNilAndNoProblems(t, createdHousehold, err)

			// Clean up household.
			assert.NoError(t, testClients.user.ArchiveHousehold(ctx, createdHousehold.ID))
		}
	})
}

func (s *TestSuite) TestHouseholds_InvitingPreExistentUser() {
	s.runForEachClient("should be possible to invite an already-registered user", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			currentStatus, statusErr := testClients.user.UserStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantHouseholdID := currentStatus.ActiveHousehold

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := converters.ConvertWebhookToWebhookCreationRequestInput(exampleWebhook)
			createdWebhook, err := testClients.user.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.user.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, relevantHouseholdID, createdWebhook.BelongsToHousehold)

			u, _, c, _ := createUserAndClientForTest(ctx, t, nil)

			invitation, err := testClients.user.InviteUserToHousehold(ctx, relevantHouseholdID, &types.HouseholdInvitationCreationRequestInput{
				Note:    t.Name(),
				ToName:  t.Name(),
				ToEmail: u.EmailAddress,
			})
			require.NoError(t, err)

			sentInvitations, err := testClients.user.GetPendingHouseholdInvitationsFromUser(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.NotEmpty(t, sentInvitations.Data)

			invitations, err := c.GetPendingHouseholdInvitationsForUser(ctx, nil)
			requireNotNilAndNoProblems(t, invitations, err)
			assert.NotEmpty(t, invitations.Data)

			err = c.AcceptHouseholdInvitation(ctx, invitation.ID, invitation.Token, t.Name())
			require.NoError(t, err)

			sentInvitations, err = testClients.user.GetPendingHouseholdInvitationsFromUser(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.Empty(t, sentInvitations.Data)

			households, err := c.GetHouseholds(ctx, nil)

			var found bool
			for _, household := range households.Data {
				if !found {
					found = household.ID == relevantHouseholdID
				}
			}

			require.True(t, found)
			require.NoError(t, c.SwitchActiveHousehold(ctx, relevantHouseholdID))

			_, err = c.GetWebhook(ctx, createdWebhook.ID)
			require.NoError(t, err)
		}
	})
}

func (s *TestSuite) TestHouseholds_InvitingUserWhoSignsUpIndependently() {
	s.runForEachClient("should be possible to invite a user before they sign up", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			currentStatus, statusErr := testClients.user.UserStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantHouseholdID := currentStatus.ActiveHousehold

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := converters.ConvertWebhookToWebhookCreationRequestInput(exampleWebhook)
			createdWebhook, err := testClients.user.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.user.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, relevantHouseholdID, createdWebhook.BelongsToHousehold)

			inviteReq := &types.HouseholdInvitationCreationRequestInput{
				Note:    t.Name(),
				ToEmail: gofakeit.Email(),
			}
			invitation, err := testClients.user.InviteUserToHousehold(ctx, relevantHouseholdID, inviteReq)
			require.NoError(t, err)

			sentInvitations, err := testClients.user.GetPendingHouseholdInvitationsFromUser(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.NotEmpty(t, sentInvitations.Data)

			_, _, c, _ := createUserAndClientForTest(ctx, t, &types.UserRegistrationInput{
				EmailAddress: inviteReq.ToEmail,
				Username:     fakes.BuildFakeUser().Username,
				Password:     gofakeit.Password(true, true, true, true, false, 64),
			})

			invitations, err := c.GetPendingHouseholdInvitationsForUser(ctx, nil)
			requireNotNilAndNoProblems(t, invitations, err)
			assert.NotEmpty(t, invitations.Data)

			err = c.AcceptHouseholdInvitation(ctx, invitation.ID, invitation.Token, t.Name())
			require.NoError(t, err)

			households, err := c.GetHouseholds(ctx, nil)

			var found bool
			for _, household := range households.Data {
				if !found {
					found = household.ID == relevantHouseholdID
				}
			}

			require.True(t, found)
			require.NoError(t, c.SwitchActiveHousehold(ctx, relevantHouseholdID))

			_, err = c.GetWebhook(ctx, createdWebhook.ID)
			require.NoError(t, err)

			sentInvitations, err = testClients.user.GetPendingHouseholdInvitationsFromUser(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.Empty(t, sentInvitations.Data)
		}
	})
}

func (s *TestSuite) TestHouseholds_InvitingUserWhoSignsUpIndependentlyAndThenCancelling() {
	s.runForEachClient("should be possible to invite a user before they sign up and cancel before they can accept", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			currentStatus, statusErr := testClients.user.UserStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantHouseholdID := currentStatus.ActiveHousehold

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := converters.ConvertWebhookToWebhookCreationRequestInput(exampleWebhook)
			createdWebhook, err := testClients.user.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.user.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, relevantHouseholdID, createdWebhook.BelongsToHousehold)

			inviteReq := &types.HouseholdInvitationCreationRequestInput{
				Note:    t.Name(),
				ToEmail: gofakeit.Email(),
			}
			invitation, err := testClients.user.InviteUserToHousehold(ctx, relevantHouseholdID, inviteReq)
			require.NoError(t, err)

			sentInvitations, err := testClients.user.GetPendingHouseholdInvitationsFromUser(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.NotEmpty(t, sentInvitations.Data)

			_, _, c, _ := createUserAndClientForTest(ctx, t, &types.UserRegistrationInput{
				EmailAddress: inviteReq.ToEmail,
				Username:     fakes.BuildFakeUser().Username,
				Password:     gofakeit.Password(true, true, true, true, false, 64),
			})

			invitations, err := c.GetPendingHouseholdInvitationsForUser(ctx, nil)
			requireNotNilAndNoProblems(t, invitations, err)
			assert.NotEmpty(t, invitations.Data)

			err = testClients.user.CancelHouseholdInvitation(ctx, invitation.ID, invitation.Token, t.Name())
			require.NoError(t, err)

			sentInvitations, err = testClients.user.GetPendingHouseholdInvitationsFromUser(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.Empty(t, sentInvitations.Data)
		}
	})
}

func (s *TestSuite) TestHouseholds_InvitingNewUserWithInviteLink() {
	s.runForEachClient("should be possible to invite a user via referral link", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			currentStatus, statusErr := testClients.user.UserStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantHouseholdID := currentStatus.ActiveHousehold

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := converters.ConvertWebhookToWebhookCreationRequestInput(exampleWebhook)
			createdWebhook, err := testClients.user.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.user.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, relevantHouseholdID, createdWebhook.BelongsToHousehold)

			inviteReq := &types.HouseholdInvitationCreationRequestInput{
				Note:    t.Name(),
				ToEmail: gofakeit.Email(),
			}
			createdInvitation, err := testClients.user.InviteUserToHousehold(ctx, relevantHouseholdID, inviteReq)
			require.NoError(t, err)

			createdInvitation, err = testClients.user.GetHouseholdInvitation(ctx, relevantHouseholdID, createdInvitation.ID)
			requireNotNilAndNoProblems(t, createdInvitation, err)

			sentInvitations, err := testClients.user.GetPendingHouseholdInvitationsFromUser(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.NotEmpty(t, sentInvitations.Data)

			_, _, c, _ := createUserAndClientForTest(ctx, t, &types.UserRegistrationInput{
				EmailAddress:    inviteReq.ToEmail,
				Username:        fakes.BuildFakeUser().Username,
				Password:        gofakeit.Password(true, true, true, true, false, 64),
				InvitationID:    createdInvitation.ID,
				InvitationToken: createdInvitation.Token,
			})

			households, err := c.GetHouseholds(ctx, nil)
			require.NoError(t, err)

			var found bool
			for _, household := range households.Data {
				if !found {
					found = household.ID == relevantHouseholdID
				}
			}

			require.True(t, found)

			_, err = c.GetWebhook(ctx, createdWebhook.ID)
			require.NoError(t, err)
		}
	})
}

func (s *TestSuite) TestHouseholds_InviteCanBeCancelled() {
	s.runForEachClient("should be possible to invite an already-registered user", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			currentStatus, statusErr := testClients.user.UserStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantHouseholdID := currentStatus.ActiveHousehold

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := converters.ConvertWebhookToWebhookCreationRequestInput(exampleWebhook)
			createdWebhook, err := testClients.user.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.user.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, relevantHouseholdID, createdWebhook.BelongsToHousehold)

			inviteReq := &types.HouseholdInvitationCreationRequestInput{
				Note:    t.Name(),
				ToEmail: gofakeit.Email(),
			}
			invitation, err := testClients.user.InviteUserToHousehold(ctx, relevantHouseholdID, inviteReq)
			require.NoError(t, err)

			require.NoError(t, testClients.user.CancelHouseholdInvitation(ctx, invitation.ID, invitation.Token, t.Name()))

			sentInvitations, err := testClients.user.GetPendingHouseholdInvitationsFromUser(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.Empty(t, sentInvitations.Data)

			_, _, c, _ := createUserAndClientForTest(ctx, t, &types.UserRegistrationInput{
				EmailAddress: inviteReq.ToEmail,
				Username:     fakes.BuildFakeUser().Username,
				Password:     gofakeit.Password(true, true, true, true, false, 64),
			})

			invitations, err := c.GetPendingHouseholdInvitationsForUser(ctx, nil)
			requireNotNilAndNoProblems(t, invitations, err)
			assert.Empty(t, invitations.Data)
		}
	})
}

func (s *TestSuite) TestHouseholds_InviteCanBeRejected() {
	s.runForEachClient("should be possible to invite an already-registered user", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			currentStatus, statusErr := testClients.user.UserStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantHouseholdID := currentStatus.ActiveHousehold

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := converters.ConvertWebhookToWebhookCreationRequestInput(exampleWebhook)
			createdWebhook, err := testClients.user.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.user.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, relevantHouseholdID, createdWebhook.BelongsToHousehold)

			u, _, c, _ := createUserAndClientForTest(ctx, t, nil)

			invitation, err := testClients.user.InviteUserToHousehold(ctx, relevantHouseholdID, &types.HouseholdInvitationCreationRequestInput{
				Note:    t.Name(),
				ToEmail: u.EmailAddress,
			})
			require.NoError(t, err)

			sentInvitations, err := testClients.user.GetPendingHouseholdInvitationsFromUser(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.NotEmpty(t, sentInvitations.Data)

			invitations, err := c.GetPendingHouseholdInvitationsForUser(ctx, nil)
			requireNotNilAndNoProblems(t, invitations, err)
			assert.NotEmpty(t, invitations.Data)

			err = c.RejectHouseholdInvitation(ctx, invitation.ID, invitation.Token, t.Name())
			require.NoError(t, err)

			sentInvitations, err = testClients.user.GetPendingHouseholdInvitationsFromUser(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.Empty(t, sentInvitations.Data)
		}
	})
}

func (s *TestSuite) TestHouseholds_ChangingMemberships() {
	s.runForCookieClient("should be possible to change members of a household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			const userCount = 1

			currentStatus, statusErr := testClients.user.UserStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)

			// fetch household data
			householdCreationInput := &types.HouseholdCreationRequestInput{
				Name: fakes.BuildFakeHousehold().Name,
			}
			household, householdCreationErr := testClients.user.CreateHousehold(ctx, householdCreationInput)
			require.NoError(t, householdCreationErr)
			require.NotNil(t, household)

			require.NoError(t, testClients.user.SwitchActiveHousehold(ctx, household.ID))

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := converters.ConvertWebhookToWebhookCreationRequestInput(exampleWebhook)
			createdWebhook, err := testClients.user.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.user.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)
			require.Equal(t, household.ID, createdWebhook.BelongsToHousehold)

			// create dummy users
			users := []*types.User{}
			clients := []*apiclient.Client{}

			// create users
			for i := 0; i < userCount; i++ {
				u, _, c, _ := createUserAndClientForTest(ctx, t, nil)
				users = append(users, u)
				clients = append(clients, c)

				currentStatus, statusErr = c.UserStatus(s.ctx)
				requireNotNilAndNoProblems(t, currentStatus, statusErr)
			}

			// check that each user cannot see the unreachable webhook
			for i := 0; i < userCount; i++ {
				webhook, err := clients[i].GetWebhook(ctx, createdWebhook.ID)
				require.Nil(t, webhook)
				require.Error(t, err)
			}

			// add them to the household
			for i := 0; i < userCount; i++ {
				invitation, invitationErr := testClients.user.InviteUserToHousehold(ctx, household.ID, &types.HouseholdInvitationCreationRequestInput{
					ToEmail: users[i].EmailAddress,
					Note:    t.Name(),
				})
				require.NoError(t, invitationErr)
				require.NotEmpty(t, invitation.ID)

				invitations, fetchInvitationsErr := clients[i].GetPendingHouseholdInvitationsForUser(ctx, nil)
				requireNotNilAndNoProblems(t, invitations, fetchInvitationsErr)
				assert.NotEmpty(t, invitations.Data)

				err = clients[i].AcceptHouseholdInvitation(ctx, invitation.ID, invitation.Token, t.Name())
				require.NoError(t, err)

				require.NoError(t, clients[i].SwitchActiveHousehold(ctx, household.ID))

				currentStatus, statusErr = clients[i].UserStatus(s.ctx)
				requireNotNilAndNoProblems(t, currentStatus, statusErr)
				require.Equal(t, currentStatus.ActiveHousehold, household.ID)
			}

			// grant all permissions
			for i := 0; i < userCount; i++ {
				input := &types.ModifyUserPermissionsInput{
					Reason:  t.Name(),
					NewRole: authorization.HouseholdAdminRole.String(),
				}
				require.NoError(t, testClients.user.ModifyMemberPermissions(ctx, household.ID, users[i].ID, input))
			}

			// check that each user can see the webhook
			for i := 0; i < userCount; i++ {
				webhook, webhookRetrievalError := clients[i].GetWebhook(ctx, createdWebhook.ID)
				requireNotNilAndNoProblems(t, webhook, webhookRetrievalError)
			}

			// remove users from household
			for i := 0; i < userCount; i++ {
				require.NoError(t, testClients.user.RemoveUserFromHousehold(ctx, household.ID, users[i].ID))
			}

			// check that each user cannot see the webhook
			for i := 0; i < userCount; i++ {
				webhook, webhookRetrievalError := clients[i].GetWebhook(ctx, createdWebhook.ID)
				require.Nil(t, webhook)
				require.Error(t, webhookRetrievalError)
			}

			// Clean up.
			require.NoError(t, testClients.user.ArchiveWebhook(ctx, createdWebhook.ID))

			for i := 0; i < userCount; i++ {
				require.NoError(t, testClients.admin.ArchiveUser(ctx, users[i].ID))
			}
		}
	})
}

func (s *TestSuite) TestHouseholds_OwnershipTransfer() {
	s.runForCookieClient("should be possible to transfer ownership of a household", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create users
			futureOwner, _, futureOwnerClient, _ := createUserAndClientForTest(ctx, t, nil)

			// fetch household data
			householdCreationInput := &types.HouseholdCreationRequestInput{
				Name: fakes.BuildFakeHousehold().Name,
			}
			household, householdCreationErr := testClients.user.CreateHousehold(ctx, householdCreationInput)
			require.NoError(t, householdCreationErr)
			require.NotNil(t, household)

			require.NoError(t, testClients.user.SwitchActiveHousehold(ctx, household.ID))

			// create a webhook

			// Create webhook.
			exampleWebhook := fakes.BuildFakeWebhook()
			exampleWebhookInput := converters.ConvertWebhookToWebhookCreationRequestInput(exampleWebhook)
			createdWebhook, err := testClients.user.CreateWebhook(ctx, exampleWebhookInput)
			require.NoError(t, err)

			checkWebhookEquality(t, exampleWebhook, createdWebhook)

			createdWebhook, err = testClients.user.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, createdWebhook, err)

			require.Equal(t, household.ID, createdWebhook.BelongsToHousehold)

			// check that user cannot see the webhook
			webhook, err := futureOwnerClient.GetWebhook(ctx, createdWebhook.ID)
			require.Nil(t, webhook)
			require.Error(t, err)

			// add them to the household
			require.NoError(t, testClients.user.TransferHouseholdOwnership(ctx, household.ID, &types.HouseholdOwnershipTransferInput{
				Reason:       t.Name(),
				CurrentOwner: household.BelongsToUser,
				NewOwner:     futureOwner.ID,
			}))

			require.NoError(t, futureOwnerClient.SwitchActiveHousehold(ctx, household.ID))

			// check that user can see the webhook
			webhook, err = futureOwnerClient.GetWebhook(ctx, createdWebhook.ID)
			requireNotNilAndNoProblems(t, webhook, err)

			// check that old user cannot see the webhook
			webhook, err = testClients.user.GetWebhook(ctx, createdWebhook.ID)
			require.Nil(t, webhook)
			require.Error(t, err)

			// check that new owner can delete the webhook
			require.NoError(t, futureOwnerClient.ArchiveWebhook(ctx, createdWebhook.ID))

			// Clean up.
			require.Error(t, testClients.user.ArchiveWebhook(ctx, createdWebhook.ID))
			require.NoError(t, testClients.admin.ArchiveUser(ctx, futureOwner.ID))
		}
	})
}

func (s *TestSuite) TestHouseholds_UsersHaveBackupHouseholdCreatedForThemWhenRemovedFromLastHousehold() {
	s.runForEachClient("should be possible to invite a user via referral link", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			currentStatus, statusErr := testClients.user.UserStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantHouseholdID := currentStatus.ActiveHousehold

			inviteReq := &types.HouseholdInvitationCreationRequestInput{
				Note:    t.Name(),
				ToEmail: gofakeit.Email(),
			}
			createdInvitation, err := testClients.user.InviteUserToHousehold(ctx, relevantHouseholdID, inviteReq)
			require.NoError(t, err)

			createdInvitation, err = testClients.user.GetHouseholdInvitation(ctx, relevantHouseholdID, createdInvitation.ID)
			requireNotNilAndNoProblems(t, createdInvitation, err)

			sentInvitations, err := testClients.user.GetPendingHouseholdInvitationsFromUser(ctx, nil)
			requireNotNilAndNoProblems(t, sentInvitations, err)
			assert.NotEmpty(t, sentInvitations.Data)

			regInput := &types.UserRegistrationInput{
				EmailAddress:    inviteReq.ToEmail,
				Username:        fakes.BuildFakeUser().Username,
				Password:        gofakeit.Password(true, true, true, true, false, 64),
				InvitationID:    createdInvitation.ID,
				InvitationToken: createdInvitation.Token,
			}
			u, _, c, _ := createUserAndClientForTest(ctx, t, regInput)

			households, err := c.GetHouseholds(ctx, nil)
			require.NoError(t, err)

			assert.Len(t, households.Data, 2)

			var (
				found            bool
				otherHouseholdID string
			)

			for _, household := range households.Data {
				if household.ID == relevantHouseholdID {
					if !found {
						found = true
					}
				} else {
					otherHouseholdID = household.ID
				}
			}

			require.NotEmpty(t, otherHouseholdID)
			require.True(t, found)

			require.NoError(t, testClients.user.RemoveUserFromHousehold(ctx, relevantHouseholdID, u.ID))

			u.HashedPassword = regInput.Password

			newCookie, err := testutils.GetLoginCookie(ctx, urlToUse, u)
			require.NoError(t, err)

			require.NoError(t, c.SetOptions(apiclient.UsingCookie(newCookie)))

			household, err := c.GetCurrentHousehold(ctx)
			requireNotNilAndNoProblems(t, household, err)
			assert.NotEqual(t, relevantHouseholdID, household.ID)

			require.True(t, found)
		}
	})
}
