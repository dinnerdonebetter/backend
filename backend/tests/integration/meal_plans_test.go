package integration

import (
	"context"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkMealPlanEquality(t *testing.T, expected, actual *types.MealPlan) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for meal plan %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.Status, actual.Status, "expected Status for meal plan %s to be %v, but it was %v", expected.ID, expected.Status, actual.Status)
	assert.WithinDuration(t, expected.VotingDeadline, actual.VotingDeadline, time.Nanosecond*1000, "expected VotingDeadline for meal plan %s to be %v, but it was %v", expected.ID, expected.VotingDeadline, actual.VotingDeadline)
	assert.Equal(t, expected.TasksCreated, actual.TasksCreated, "expected TasksCreated for meal plan %s to be %v, but it was %v", expected.ID, expected.TasksCreated, actual.TasksCreated)
	assert.Equal(t, expected.ElectionMethod, actual.ElectionMethod, "expected ElectionMethod for meal plan %s to be %v, but it was %v", expected.ID, expected.ElectionMethod, actual.ElectionMethod)
	assert.Equal(t, expected.GroceryListInitialized, actual.GroceryListInitialized, "expected GroceryListInitialized for meal plan %s to be %v, but it was %v", expected.ID, expected.GroceryListInitialized, actual.GroceryListInitialized)
	assert.NotZero(t, actual.CreatedAt)
}

func createMealPlanForTest(ctx context.Context, t *testing.T, mealPlan *types.MealPlan, adminClient, client *apiclient.Client) *types.MealPlan {
	t.Helper()

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

			// create a user for the meal plan household
			_, _, householdAdminUserClient, _ := createUserAndClientForTest(ctx, t, nil)

			// create household members
			currentStatus, statusErr := householdAdminUserClient.UserStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantHouseholdID := currentStatus.ActiveHousehold

			createdUsers := []*types.User{}
			createdClients := []*apiclient.Client{}

			for i := 0; i < 2; i++ {
				u, _, c, _ := createUserAndClientForTest(ctx, t, nil)

				invitation, err := householdAdminUserClient.InviteUserToHousehold(ctx, relevantHouseholdID, &types.HouseholdInvitationCreationRequestInput{
					Note:    t.Name(),
					ToEmail: u.EmailAddress,
				})
				require.NoError(t, err)

				sentInvitations, err := householdAdminUserClient.GetPendingHouseholdInvitationsFromUser(ctx, nil)
				requireNotNilAndNoProblems(t, sentInvitations, err)
				assert.NotEmpty(t, sentInvitations.Data)

				invitations, err := c.GetPendingHouseholdInvitationsForUser(ctx, nil)
				requireNotNilAndNoProblems(t, invitations, err)
				assert.NotEmpty(t, invitations.Data)

				require.NoError(t, c.AcceptHouseholdInvitation(ctx, invitation.ID, invitation.Token, t.Name()))
				require.NoError(t, c.SwitchActiveHousehold(ctx, relevantHouseholdID))

				createdUsers = append(createdUsers, u)
				createdClients = append(createdClients, c)
			}

			// create recipes for meal plan
			createdMeals := []*types.Meal{}
			for i := 0; i < 3; i++ {
				createdMeal := createMealForTest(ctx, t, testClients.admin, householdAdminUserClient, nil)
				createdMeals = append(createdMeals, createdMeal)
			}

			const baseDeadline = 10 * time.Second
			now := time.Now()

			exampleMealPlan := &types.MealPlan{
				Notes:          t.Name(),
				Status:         string(types.MealPlanStatusAwaitingVotes),
				VotingDeadline: now.Add(baseDeadline),
				ElectionMethod: types.MealPlanElectionMethodSchulze,
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
			mealPlanCreationResult, err := householdAdminUserClient.CreateMealPlan(ctx, exampleMealPlanInput)
			require.NotEmpty(t, mealPlanCreationResult.ID)
			require.NoError(t, err)

			createdMealPlan, err := householdAdminUserClient.GetMealPlan(ctx, mealPlanCreationResult.ID)
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

			createdMealPlanOptionVotesB, err := createdClients[1].CreateMealPlanOptionVote(ctx, createdMealPlan.ID, createdMealPlanEvent.ID, userBVotes)
			require.NoError(t, err)
			require.NotNil(t, createdMealPlanOptionVotesB)

			createdMealPlanOptionVotesC, err := householdAdminUserClient.CreateMealPlanOptionVote(ctx, createdMealPlan.ID, createdMealPlanEvent.ID, userCVotes)
			require.NoError(t, err)
			require.NotNil(t, createdMealPlanOptionVotesC)

			createdMealPlan.VotingDeadline = time.Now().Add(-time.Minute)
			require.NoError(t, dbmanager.UpdateMealPlan(ctx, createdMealPlan))

			runRes, err := testClients.admin.RunFinalizeMealPlansWorker(ctx, &types.FinalizeMealPlansRequest{ReturnCount: true})
			require.NoError(t, err)
			require.NotNil(t, runRes)

			createdMealPlan, err = householdAdminUserClient.GetMealPlan(ctx, createdMealPlan.ID)
			requireNotNilAndNoProblems(t, createdMealPlan, err)
			assert.Equal(t, string(types.MealPlanStatusFinalized), createdMealPlan.Status)

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

			// create a user for the meal plan household
			_, _, householdAdminUserClient, _ := createUserAndClientForTest(ctx, t, nil)

			// create household members
			currentStatus, statusErr := householdAdminUserClient.UserStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantHouseholdID := currentStatus.ActiveHousehold

			createdUsers := []*types.User{}
			createdClients := []*apiclient.Client{}

			for i := 0; i < 2; i++ {
				u, _, c, _ := createUserAndClientForTest(ctx, t, nil)

				invitation, err := householdAdminUserClient.InviteUserToHousehold(ctx, relevantHouseholdID, &types.HouseholdInvitationCreationRequestInput{
					Note:    t.Name(),
					ToEmail: u.EmailAddress,
				})
				require.NoError(t, err)

				sentInvitations, err := householdAdminUserClient.GetPendingHouseholdInvitationsFromUser(ctx, nil)
				requireNotNilAndNoProblems(t, sentInvitations, err)
				assert.NotEmpty(t, sentInvitations.Data)

				invitations, err := c.GetPendingHouseholdInvitationsForUser(ctx, nil)
				requireNotNilAndNoProblems(t, invitations, err)
				assert.NotEmpty(t, invitations.Data)

				require.NoError(t, c.AcceptHouseholdInvitation(ctx, invitation.ID, invitation.Token, t.Name()))

				require.NoError(t, c.SwitchActiveHousehold(ctx, relevantHouseholdID))

				createdUsers = append(createdUsers, u)
				createdClients = append(createdClients, c)
			}

			// create meals for meal plan
			createdMeals := []*types.Meal{}
			for i := 0; i < 3; i++ {
				createdMeal := createMealForTest(ctx, t, testClients.admin, householdAdminUserClient, nil)
				createdMeals = append(createdMeals, createdMeal)
			}

			const baseDeadline = 10 * time.Second
			now := time.Now()

			exampleMealPlan := &types.MealPlan{
				Notes:          t.Name(),
				Status:         string(types.MealPlanStatusAwaitingVotes),
				VotingDeadline: now.Add(baseDeadline),
				ElectionMethod: types.MealPlanElectionMethodSchulze,
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
			createdMealPlan, err := householdAdminUserClient.CreateMealPlan(ctx, exampleMealPlanInput)
			require.NotEmpty(t, createdMealPlan.ID)
			require.NoError(t, err)

			createdMealPlan, err = householdAdminUserClient.GetMealPlan(ctx, createdMealPlan.ID)
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

			createdMealPlanOptionVotesB, err := createdClients[1].CreateMealPlanOptionVote(ctx, createdMealPlan.ID, createdMealPlanEvent.ID, userBVotes)
			require.NoError(t, err)
			require.NotNil(t, createdMealPlanOptionVotesB)

			createdMealPlan, err = householdAdminUserClient.GetMealPlan(ctx, createdMealPlan.ID)
			requireNotNilAndNoProblems(t, createdMealPlan, err)
			assert.Equal(t, string(types.MealPlanStatusAwaitingVotes), createdMealPlan.Status)

			createdMealPlan.VotingDeadline = time.Now().Add(-10 * time.Hour)
			require.NoError(t, dbmanager.UpdateMealPlan(ctx, createdMealPlan))

			runRes, err := testClients.admin.RunFinalizeMealPlansWorker(ctx, &types.FinalizeMealPlansRequest{ReturnCount: true})
			require.NoError(t, err)
			require.NotNil(t, runRes)

			createdMealPlan, err = householdAdminUserClient.GetMealPlan(ctx, createdMealPlan.ID)
			requireNotNilAndNoProblems(t, createdMealPlan, err)
			assert.Equal(t, string(types.MealPlanStatusFinalized), createdMealPlan.Status)

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
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			for _, createdMealPlan := range expected {
				assert.NoError(t, testClients.user.ArchiveMealPlan(ctx, createdMealPlan.ID))
			}
		}
	})
}
