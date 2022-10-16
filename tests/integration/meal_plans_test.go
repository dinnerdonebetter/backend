package integration

import (
	"context"
	"github.com/prixfixeco/api_server/pkg/types/converters"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/apiclient"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func checkMealPlanEquality(t *testing.T, expected, actual *types.MealPlan) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected StatusExplanation for meal plan %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.Status, actual.Status, "expected CreationExplanation for meal plan %s to be %v, but it was %v", expected.ID, expected.Status, actual.Status)
	assert.WithinDuration(t, expected.VotingDeadline, actual.VotingDeadline, time.Nanosecond*1000, "expected VotingDeadline for meal plan %s to be %v, but it was %v", expected.ID, expected.VotingDeadline, actual.VotingDeadline)
	assert.NotZero(t, actual.CreatedAt)
}

func createMealPlanForTest(ctx context.Context, t *testing.T, mealPlan *types.MealPlan, adminClient, client *apiclient.Client) *types.MealPlan {
	t.Helper()

	t.Log("creating meal plan")
	if mealPlan == nil {
		mealPlan = fakes.BuildFakeMealPlan()
		for i, evt := range mealPlan.Events {
			for j := range evt.Options {
				createdMeal := createMealForTest(ctx, t, adminClient, client, nil)
				mealPlan.Events[i].Options[j].Meal.ID = createdMeal.ID
				mealPlan.Events[i].Options[j].AssignedCook = nil
			}
		}
	}

	exampleMealPlanInput := converters.ConvertMealPlanToMealPlanCreationRequestInput(mealPlan)
	createdMealPlan, err := client.CreateMealPlan(ctx, exampleMealPlanInput)
	require.NoError(t, err)
	require.NotEmpty(t, createdMealPlan.ID)

	t.Logf("meal plan %q created", createdMealPlan.ID)

	createdMealPlan, err = client.GetMealPlan(ctx, createdMealPlan.ID)
	requireNotNilAndNoProblems(t, createdMealPlan, err)
	checkMealPlanEquality(t, mealPlan, createdMealPlan)

	return createdMealPlan
}

func (s *TestSuite) TestMealPlans_CompleteLifecycleForAllVotesReceived() {
	s.runForEachClient("should resolve the meal plan status upon receiving all votes", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create household members
			t.Logf("determining household ID")
			currentStatus, statusErr := testClients.user.UserStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantHouseholdID := currentStatus.ActiveHousehold
			t.Logf("initial household is %s; initial user ID is %s", relevantHouseholdID, s.user.ID)

			createdUsers := []*types.User{}
			createdClients := []*apiclient.Client{}

			for i := 0; i < 2; i++ {
				t.Logf("creating user to invite")
				u, _, c, _ := createUserAndClientForTest(ctx, t, nil)

				t.Logf("inviting user")
				invitation, err := testClients.user.InviteUserToHousehold(ctx, &types.HouseholdInvitationCreationRequestInput{
					FromUser:               s.user.ID,
					Note:                   t.Name(),
					ToEmail:                u.EmailAddress,
					DestinationHouseholdID: relevantHouseholdID,
				})
				require.NoError(t, err)

				t.Logf("checking for sent invitation")
				sentInvitations, err := testClients.user.GetPendingHouseholdInvitationsFromUser(ctx, nil)
				requireNotNilAndNoProblems(t, sentInvitations, err)
				assert.NotEmpty(t, sentInvitations.HouseholdInvitations)

				t.Logf("checking for received invitation")
				invitations, err := c.GetPendingHouseholdInvitationsForUser(ctx, nil)
				requireNotNilAndNoProblems(t, invitations, err)
				assert.NotEmpty(t, invitations.HouseholdInvitations)

				t.Logf("accepting invitation")
				require.NoError(t, c.AcceptHouseholdInvitation(ctx, invitation.ID, invitation.Token, t.Name()))

				require.NoError(t, c.SwitchActiveHousehold(ctx, relevantHouseholdID))

				createdUsers = append(createdUsers, u)
				createdClients = append(createdClients, c)
			}

			// create recipes for meal plan
			createdMeals := []*types.Meal{}
			for i := 0; i < 3; i++ {
				createdMeal := createMealForTest(ctx, t, testClients.admin, testClients.user, nil)
				createdMeals = append(createdMeals, createdMeal)
			}

			const baseDeadline = 10 * time.Second
			now := time.Now()

			t.Log("creating meal plan")
			exampleMealPlan := &types.MealPlan{
				Notes:          t.Name(),
				Status:         types.AwaitingVotesMealPlanStatus,
				VotingDeadline: now.Add(baseDeadline),
				Events: []*types.MealPlanEvent{
					{
						StartsAt: now.Add(24 * time.Hour),
						EndsAt:   now.Add(72 * time.Hour),
						MealName: types.BreakfastMealName,
						Options: []*types.MealPlanOption{
							{
								Meal:  types.Meal{ID: createdMeals[0].ID},
								Notes: "option A",
							},
							{
								Meal:  types.Meal{ID: createdMeals[1].ID},
								Notes: "option B",
							},
							{
								Meal:  types.Meal{ID: createdMeals[2].ID},
								Notes: "option C",
							},
						},
					},
				},
			}

			exampleMealPlanInput := converters.ConvertMealPlanToMealPlanCreationRequestInput(exampleMealPlan)
			createdMealPlan, err := testClients.user.CreateMealPlan(ctx, exampleMealPlanInput)
			require.NotEmpty(t, createdMealPlan.ID)
			require.NoError(t, err)
			t.Logf("meal plan %q created", createdMealPlan.ID)

			createdMealPlan, err = testClients.user.GetMealPlan(ctx, createdMealPlan.ID)
			requireNotNilAndNoProblems(t, createdMealPlan, err)
			checkMealPlanEquality(t, exampleMealPlan, createdMealPlan)

			require.NotEmpty(t, createdMealPlan.Events)
			require.NotEmpty(t, createdMealPlan.Events[0].Options)

			createdMealPlanEvent := createdMealPlan.Events[0]
			require.NotNil(t, createdMealPlanEvent)

			userAVotes := &types.MealPlanOptionVoteCreationRequestInput{
				Votes: []*types.MealPlanOptionVoteCreationInput{
					{
						BelongsToMealPlanOption: createdMealPlanEvent.Options[0].ID,
						Rank:                    0,
					},
					{
						BelongsToMealPlanOption: createdMealPlanEvent.Options[1].ID,
						Rank:                    2,
					},
					{
						BelongsToMealPlanOption: createdMealPlanEvent.Options[2].ID,
						Rank:                    1,
					},
				},
			}

			userBVotes := &types.MealPlanOptionVoteCreationRequestInput{
				Votes: []*types.MealPlanOptionVoteCreationInput{
					{
						BelongsToMealPlanOption: createdMealPlanEvent.Options[0].ID,
						Rank:                    0,
					},
					{
						BelongsToMealPlanOption: createdMealPlanEvent.Options[1].ID,
						Rank:                    1,
					},
					{
						BelongsToMealPlanOption: createdMealPlanEvent.Options[2].ID,
						Rank:                    2,
					},
				},
			}

			userCVotes := &types.MealPlanOptionVoteCreationRequestInput{
				Votes: []*types.MealPlanOptionVoteCreationInput{
					{
						BelongsToMealPlanOption: createdMealPlanEvent.Options[0].ID,
						Rank:                    1,
					},
					{
						BelongsToMealPlanOption: createdMealPlanEvent.Options[1].ID,
						Rank:                    0,
					},
					{
						BelongsToMealPlanOption: createdMealPlanEvent.Options[2].ID,
						Rank:                    2,
					},
				},
			}

			createdMealPlanOptionVotesA, err := createdClients[0].CreateMealPlanOptionVote(ctx, createdMealPlan.ID, createdMealPlanEvent.ID, userAVotes)
			require.NoError(t, err)
			require.NotNil(t, createdMealPlanOptionVotesA)
			t.Logf("meal plan option votes created for user A")

			createdMealPlanOptionVotesB, err := createdClients[1].CreateMealPlanOptionVote(ctx, createdMealPlan.ID, createdMealPlanEvent.ID, userBVotes)
			require.NoError(t, err)
			require.NotNil(t, createdMealPlanOptionVotesB)
			t.Logf("meal plan option votes created for user B")

			createdMealPlanOptionVotesC, err := testClients.user.CreateMealPlanOptionVote(ctx, createdMealPlan.ID, createdMealPlanEvent.ID, userCVotes)
			require.NoError(t, err)
			require.NotNil(t, createdMealPlanOptionVotesC)
			t.Logf("meal plan option votes created for user C")

			time.Sleep(baseDeadline * 2)

			createdMealPlan, err = testClients.user.GetMealPlan(ctx, createdMealPlan.ID)
			requireNotNilAndNoProblems(t, createdMealPlan, err)
			assert.Equal(t, types.FinalizedMealPlanStatus, createdMealPlan.Status)

			for _, event := range createdMealPlan.Events {
				selectionMade := false
				for _, opt := range event.Options {
					if opt.Chosen {
						selectionMade = true
						break
					}
				}
				require.True(t, selectionMade)
			}
		}
	})
}

func (s *TestSuite) TestMealPlans_CompleteLifecycleForSomeVotesReceived() {
	s.runForEachClient("should resolve the meal plan status upon voting deadline expiry", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create household members
			t.Logf("determining household ID")
			currentStatus, statusErr := testClients.user.UserStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantHouseholdID := currentStatus.ActiveHousehold
			t.Logf("initial household is %s; initial user ID is %s", relevantHouseholdID, s.user.ID)

			createdUsers := []*types.User{}
			createdClients := []*apiclient.Client{}

			for i := 0; i < 2; i++ {
				t.Logf("creating user to invite")
				u, _, c, _ := createUserAndClientForTest(ctx, t, nil)

				t.Logf("inviting user")
				invitation, err := testClients.user.InviteUserToHousehold(ctx, &types.HouseholdInvitationCreationRequestInput{
					FromUser:               s.user.ID,
					Note:                   t.Name(),
					ToEmail:                u.EmailAddress,
					DestinationHouseholdID: relevantHouseholdID,
				})
				require.NoError(t, err)

				t.Logf("checking for sent invitation")
				sentInvitations, err := testClients.user.GetPendingHouseholdInvitationsFromUser(ctx, nil)
				requireNotNilAndNoProblems(t, sentInvitations, err)
				assert.NotEmpty(t, sentInvitations.HouseholdInvitations)

				t.Logf("checking for received invitation")
				invitations, err := c.GetPendingHouseholdInvitationsForUser(ctx, nil)
				requireNotNilAndNoProblems(t, invitations, err)
				assert.NotEmpty(t, invitations.HouseholdInvitations)

				t.Logf("accepting invitation")
				require.NoError(t, c.AcceptHouseholdInvitation(ctx, invitation.ID, invitation.Token, t.Name()))

				require.NoError(t, c.SwitchActiveHousehold(ctx, relevantHouseholdID))

				createdUsers = append(createdUsers, u)
				createdClients = append(createdClients, c)
			}

			// create recipes for meal plan
			createdMeals := []*types.Meal{}
			for i := 0; i < 3; i++ {
				createdMeal := createMealForTest(ctx, t, testClients.admin, testClients.user, nil)
				createdMeals = append(createdMeals, createdMeal)
			}

			const baseDeadline = 10 * time.Second
			now := time.Now()

			t.Log("creating meal plan")
			exampleMealPlan := &types.MealPlan{
				Notes:          t.Name(),
				Status:         types.AwaitingVotesMealPlanStatus,
				VotingDeadline: now.Add(baseDeadline),
				Events: []*types.MealPlanEvent{
					{
						StartsAt: now.Add(24 * time.Hour),
						EndsAt:   now.Add(72 * time.Hour),
						MealName: types.BreakfastMealName,
						Options: []*types.MealPlanOption{
							{
								Meal:  types.Meal{ID: createdMeals[0].ID},
								Notes: "option A",
							},
							{
								Meal:  types.Meal{ID: createdMeals[1].ID},
								Notes: "option B",
							},
							{
								Meal:  types.Meal{ID: createdMeals[2].ID},
								Notes: "option C",
							},
						},
					},
				},
			}

			exampleMealPlanInput := converters.ConvertMealPlanToMealPlanCreationRequestInput(exampleMealPlan)
			createdMealPlan, err := testClients.user.CreateMealPlan(ctx, exampleMealPlanInput)
			require.NotEmpty(t, createdMealPlan.ID)
			require.NoError(t, err)

			t.Logf("meal plan %q created", createdMealPlan.ID)

			createdMealPlan, err = testClients.user.GetMealPlan(ctx, createdMealPlan.ID)
			requireNotNilAndNoProblems(t, createdMealPlan, err)
			checkMealPlanEquality(t, exampleMealPlan, createdMealPlan)

			createdMealPlanEvent := createdMealPlan.Events[0]

			userAVotes := &types.MealPlanOptionVoteCreationRequestInput{
				Votes: []*types.MealPlanOptionVoteCreationInput{
					{
						BelongsToMealPlanOption: createdMealPlanEvent.Options[0].ID,
						Rank:                    0,
					},
					{
						BelongsToMealPlanOption: createdMealPlanEvent.Options[1].ID,
						Rank:                    2,
					},
					{
						BelongsToMealPlanOption: createdMealPlanEvent.Options[2].ID,
						Rank:                    1,
					},
				},
			}

			userBVotes := &types.MealPlanOptionVoteCreationRequestInput{
				Votes: []*types.MealPlanOptionVoteCreationInput{
					{
						BelongsToMealPlanOption: createdMealPlanEvent.Options[0].ID,
						Rank:                    0,
					},
					{
						BelongsToMealPlanOption: createdMealPlanEvent.Options[1].ID,
						Rank:                    1,
					},
					{
						BelongsToMealPlanOption: createdMealPlanEvent.Options[2].ID,
						Rank:                    2,
					},
				},
			}

			createdMealPlanOptionVotesA, err := createdClients[0].CreateMealPlanOptionVote(ctx, createdMealPlan.ID, createdMealPlanEvent.ID, userAVotes)
			require.NoError(t, err)
			require.NotNil(t, createdMealPlanOptionVotesA)
			t.Logf("meal plan option votes created for user A")

			createdMealPlanOptionVotesB, err := createdClients[1].CreateMealPlanOptionVote(ctx, createdMealPlan.ID, createdMealPlanEvent.ID, userBVotes)
			require.NoError(t, err)
			require.NotNil(t, createdMealPlanOptionVotesB)
			t.Logf("meal plan option votes created for user B")

			createdMealPlan, err = testClients.user.GetMealPlan(ctx, createdMealPlan.ID)
			requireNotNilAndNoProblems(t, createdMealPlan, err)
			assert.Equal(t, types.AwaitingVotesMealPlanStatus, createdMealPlan.Status)

			createdMealPlan.VotingDeadline = time.Now().Add(-10 * time.Hour)
			require.NoError(t, dbmanager.UpdateMealPlan(ctx, createdMealPlan))

			time.Sleep(baseDeadline * 2)

			createdMealPlan, err = testClients.user.GetMealPlan(ctx, createdMealPlan.ID)
			requireNotNilAndNoProblems(t, createdMealPlan, err)
			assert.Equal(t, types.FinalizedMealPlanStatus, createdMealPlan.Status)

			for _, event := range createdMealPlan.Events {
				selectionMade := false
				for _, opt := range event.Options {
					if opt.Chosen {
						selectionMade = true
						break
					}
				}
				require.True(t, selectionMade)
			}
		}
	})
}

func (s *TestSuite) TestMealPlans_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating meal plans")
			var expected []*types.MealPlan
			for i := 0; i < 5; i++ {
				createdMealPlan := createMealPlanForTest(ctx, t, nil, testClients.admin, testClients.user)
				expected = append(expected, createdMealPlan)
			}

			// assert meal plan list equality
			actual, err := testClients.user.GetMealPlans(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.MealPlans),
				"expected %d to be <= %d",
				len(expected),
				len(actual.MealPlans),
			)

			t.Log("cleaning up")
			for _, createdMealPlan := range expected {
				assert.NoError(t, testClients.user.ArchiveMealPlan(ctx, createdMealPlan.ID))
			}
		}
	})
}
