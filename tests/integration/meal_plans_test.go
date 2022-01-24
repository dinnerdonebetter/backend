package integration

import (
	"context"
	"testing"
	"time"

	"github.com/prixfixeco/api_server/internal/observability/tracing"

	"github.com/prixfixeco/api_server/pkg/client/httpclient"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func checkMealPlanEquality(t *testing.T, expected, actual *types.MealPlan) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for meal plan %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.Status, actual.Status, "expected Status for meal plan %s to be %v, but it was %v", expected.ID, expected.Status, actual.Status)
	assert.Equal(t, expected.StartsAt, actual.StartsAt, "expected StartsAt for meal plan %s to be %v, but it was %v", expected.ID, expected.StartsAt, actual.StartsAt)
	assert.Equal(t, expected.EndsAt, actual.EndsAt, "expected EndsAt for meal plan %s to be %v, but it was %v", expected.ID, expected.EndsAt, actual.EndsAt)
	assert.NotZero(t, actual.CreatedOn)
}

func createMealPlanWithNotificationChannel(ctx context.Context, t *testing.T, client *httpclient.Client) *types.MealPlan {
	t.Helper()

	t.Log("creating meal plan")
	exampleMealPlan := fakes.BuildFakeMealPlan()
	for i := range exampleMealPlan.Options {
		createdMeal := createMealForTest(ctx, t, client, nil)
		exampleMealPlan.Options[i].Meal.ID = createdMeal.ID
	}

	exampleMealPlanInput := fakes.BuildFakeMealPlanCreationRequestInputFromMealPlan(exampleMealPlan)
	createdMealPlan, err := client.CreateMealPlan(ctx, exampleMealPlanInput)
	require.NotEmpty(t, createdMealPlan.ID)
	require.NoError(t, err)

	t.Logf("meal plan %q created", createdMealPlan.ID)

	createdMealPlan, err = client.GetMealPlan(ctx, createdMealPlan.ID)
	requireNotNilAndNoProblems(t, createdMealPlan, err)
	checkMealPlanEquality(t, exampleMealPlan, createdMealPlan)

	return createdMealPlan
}

var allDays = []time.Weekday{
	time.Monday,
	time.Tuesday,
	time.Wednesday,
	time.Thursday,
	time.Friday,
	time.Saturday,
	time.Sunday,
}

var allMealNames = []types.MealName{
	types.BreakfastMealName,
	types.SecondBreakfastMealName,
	types.BrunchMealName,
	types.LunchMealName,
	types.SupperMealName,
	types.DinnerMealName,
}

func byDayAndMeal(l []*types.MealPlanOption, day time.Weekday, meal types.MealName) []*types.MealPlanOption {
	out := []*types.MealPlanOption{}

	for _, o := range l {
		if o.Day == day && o.MealName == meal {
			out = append(out, o)
		}
	}

	return out
}

func (s *TestSuite) TestMealPlans_CompleteLifecycleForAllVotesReceived() {
	s.runForEachClient("should resolve the meal plan status upon receiving all votes", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create household members
			t.Logf("determining household ID")
			currentStatus, statusErr := testClients.main.UserStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantHouseholdID := currentStatus.ActiveHousehold
			t.Logf("initial household is %s; initial user ID is %s", relevantHouseholdID, s.user.ID)

			createdUsers := []*types.User{}
			createdClients := []*httpclient.Client{}

			for i := 0; i < 2; i++ {
				t.Logf("creating user to invite")
				u, _, c, _ := createUserAndClientForTest(ctx, t, nil)

				t.Logf("inviting user")
				invitation, err := testClients.main.InviteUserToHousehold(ctx, &types.HouseholdInvitationCreationRequestInput{
					FromUser:             s.user.ID,
					Note:                 t.Name(),
					ToEmail:              u.EmailAddress,
					DestinationHousehold: relevantHouseholdID,
				})
				require.NoError(t, err)

				t.Logf("checking for sent invitation")
				sentInvitations, err := testClients.main.GetPendingHouseholdInvitationsFromUser(ctx, nil)
				requireNotNilAndNoProblems(t, sentInvitations, err)
				assert.NotEmpty(t, sentInvitations.HouseholdInvitations)

				t.Logf("checking for received invitation")
				invitations, err := c.GetPendingHouseholdInvitationsForUser(ctx, nil)
				requireNotNilAndNoProblems(t, invitations, err)
				assert.NotEmpty(t, invitations.HouseholdInvitations)

				t.Logf("accepting invitation")
				require.NoError(t, c.AcceptHouseholdInvitation(ctx, relevantHouseholdID, invitation.ID, t.Name()))

				require.NoError(t, c.SwitchActiveHousehold(ctx, relevantHouseholdID))

				createdUsers = append(createdUsers, u)
				createdClients = append(createdClients, c)
			}

			// create recipes for meal plan
			createdMeals := []*types.Meal{}
			for i := 0; i < 3; i++ {
				createdMeal := createMealForTest(ctx, t, testClients.main, nil)
				createdMeals = append(createdMeals, createdMeal)
			}

			t.Log("creating meal plan")
			exampleMealPlan := &types.MealPlan{
				Notes:          t.Name(),
				Status:         types.AwaitingVotesMealPlanStatus,
				StartsAt:       uint64(time.Now().Add(24 * time.Hour).Unix()),
				EndsAt:         uint64(time.Now().Add(72 * time.Hour).Unix()),
				VotingDeadline: uint64(time.Now().Add(10 * time.Minute).Unix()),
				Options: []*types.MealPlanOption{
					{
						Meal:     types.Meal{ID: createdMeals[0].ID},
						Notes:    "option A",
						MealName: types.BreakfastMealName,
						Day:      time.Monday,
					},
					{
						Meal:     types.Meal{ID: createdMeals[1].ID},
						Notes:    "option B",
						MealName: types.BreakfastMealName,
						Day:      time.Monday,
					},
					{
						Meal:     types.Meal{ID: createdMeals[2].ID},
						Notes:    "option C",
						MealName: types.BreakfastMealName,
						Day:      time.Monday,
					},
				},
			}

			exampleMealPlanInput := fakes.BuildFakeMealPlanCreationRequestInputFromMealPlan(exampleMealPlan)
			createdMealPlan, err := testClients.main.CreateMealPlan(ctx, exampleMealPlanInput)
			require.NotEmpty(t, createdMealPlan.ID)
			require.NoError(t, err)
			t.Logf("meal plan %q created", createdMealPlan.ID)

			createdMealPlan, err = testClients.main.GetMealPlan(ctx, createdMealPlan.ID)
			requireNotNilAndNoProblems(t, createdMealPlan, err)
			checkMealPlanEquality(t, exampleMealPlan, createdMealPlan)

			userAVotes := &types.MealPlanOptionVoteCreationRequestInput{
				Votes: []*types.MealPlanOptionVoteCreationInput{
					{
						BelongsToMealPlanOption: createdMealPlan.Options[0].ID,
						Rank:                    0,
					},
					{
						BelongsToMealPlanOption: createdMealPlan.Options[1].ID,
						Rank:                    2,
					},
					{
						BelongsToMealPlanOption: createdMealPlan.Options[2].ID,
						Rank:                    1,
					},
				},
			}

			userBVotes := &types.MealPlanOptionVoteCreationRequestInput{
				Votes: []*types.MealPlanOptionVoteCreationInput{
					{
						BelongsToMealPlanOption: createdMealPlan.Options[0].ID,
						Rank:                    0,
					},
					{
						BelongsToMealPlanOption: createdMealPlan.Options[1].ID,
						Rank:                    1,
					},
					{
						BelongsToMealPlanOption: createdMealPlan.Options[2].ID,
						Rank:                    2,
					},
				},
			}

			userCVotes := &types.MealPlanOptionVoteCreationRequestInput{
				Votes: []*types.MealPlanOptionVoteCreationInput{
					{
						BelongsToMealPlanOption: createdMealPlan.Options[0].ID,
						Rank:                    1,
					},
					{
						BelongsToMealPlanOption: createdMealPlan.Options[1].ID,
						Rank:                    0,
					},
					{
						BelongsToMealPlanOption: createdMealPlan.Options[2].ID,
						Rank:                    2,
					},
				},
			}

			createdMealPlanOptionVotesA, err := createdClients[0].CreateMealPlanOptionVote(ctx, createdMealPlan.ID, userAVotes)
			require.NoError(t, err)
			require.NotNil(t, createdMealPlanOptionVotesA)
			t.Logf("meal plan option votes created for user A")

			//checkMealPlanOptionVoteEquality(t, vote, createdMealPlanOptionVotesA)

			//createdMealPlanOptionVotesA, err = createdClients[0].GetMealPlanOptionVote(ctx, createdMealPlan.ID, vote.BelongsToMealPlanOption, createdMealPlanOptionVotesA.ID)
			//requireNotNilAndNoProblems(t, createdMealPlanOptionVotesA, err)
			//require.Equal(t, vote.BelongsToMealPlanOption, createdMealPlanOptionVotesA.BelongsToMealPlanOption)
			//checkMealPlanOptionVoteEquality(t, vote, createdMealPlanOptionVotesA)

			createdMealPlanOptionVotesB, err := createdClients[1].CreateMealPlanOptionVote(ctx, createdMealPlan.ID, userBVotes)
			require.NoError(t, err)
			require.NotNil(t, createdMealPlanOptionVotesB)
			t.Logf("meal plan option votes created for user B")

			//checkMealPlanOptionVoteEquality(t, vote, createdMealPlanOptionVotesB)

			//createdMealPlanOptionVotesB, err = createdClients[1].GetMealPlanOptionVote(ctx, createdMealPlan.ID, vote.BelongsToMealPlanOption, createdMealPlanOptionVotesB.ID)
			//requireNotNilAndNoProblems(t, createdMealPlanOptionVotesB, err)
			//require.Equal(t, vote.BelongsToMealPlanOption, createdMealPlanOptionVotesB.BelongsToMealPlanOption)
			//checkMealPlanOptionVoteEquality(t, vote, createdMealPlanOptionVotesB)

			createdMealPlanOptionVotesC, err := testClients.main.CreateMealPlanOptionVote(ctx, createdMealPlan.ID, userCVotes)
			require.NoError(t, err)
			require.NotNil(t, createdMealPlanOptionVotesC)
			t.Logf("meal plan option votes created for user C")

			//checkMealPlanOptionVoteEquality(t, vote, createdMealPlanOptionVotesC)

			//createdMealPlanOptionVotesC, err = testClients.main.GetMealPlanOptionVote(ctx, createdMealPlan.ID, vote.BelongsToMealPlanOption, createdMealPlanOptionVotesC.ID)
			//requireNotNilAndNoProblems(t, createdMealPlanOptionVotesC, err)
			//require.Equal(t, vote.BelongsToMealPlanOption, createdMealPlanOptionVotesC.BelongsToMealPlanOption)
			//checkMealPlanOptionVoteEquality(t, vote, createdMealPlanOptionVotesC)

			time.Sleep(5 * time.Second)

			createdMealPlan, err = testClients.main.GetMealPlan(ctx, createdMealPlan.ID)
			requireNotNilAndNoProblems(t, createdMealPlan, err)
			assert.Equal(t, types.FinalizedMealPlanStatus, createdMealPlan.Status)

			for _, day := range allDays {
				for _, mealName := range allMealNames {
					options := byDayAndMeal(createdMealPlan.Options, day, mealName)
					if len(options) > 0 {
						selectionMade := false
						for _, opt := range options {
							if opt.Chosen {
								selectionMade = true
								break
							}
						}
						require.True(t, selectionMade)
					}
				}
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
			currentStatus, statusErr := testClients.main.UserStatus(s.ctx)
			requireNotNilAndNoProblems(t, currentStatus, statusErr)
			relevantHouseholdID := currentStatus.ActiveHousehold
			t.Logf("initial household is %s; initial user ID is %s", relevantHouseholdID, s.user.ID)

			createdUsers := []*types.User{}
			createdClients := []*httpclient.Client{}

			for i := 0; i < 2; i++ {
				t.Logf("creating user to invite")
				u, _, c, _ := createUserAndClientForTest(ctx, t, nil)

				t.Logf("inviting user")
				invitation, err := testClients.main.InviteUserToHousehold(ctx, &types.HouseholdInvitationCreationRequestInput{
					FromUser:             s.user.ID,
					Note:                 t.Name(),
					ToEmail:              u.EmailAddress,
					DestinationHousehold: relevantHouseholdID,
				})
				require.NoError(t, err)

				t.Logf("checking for sent invitation")
				sentInvitations, err := testClients.main.GetPendingHouseholdInvitationsFromUser(ctx, nil)
				requireNotNilAndNoProblems(t, sentInvitations, err)
				assert.NotEmpty(t, sentInvitations.HouseholdInvitations)

				t.Logf("checking for received invitation")
				invitations, err := c.GetPendingHouseholdInvitationsForUser(ctx, nil)
				requireNotNilAndNoProblems(t, invitations, err)
				assert.NotEmpty(t, invitations.HouseholdInvitations)

				t.Logf("accepting invitation")
				require.NoError(t, c.AcceptHouseholdInvitation(ctx, relevantHouseholdID, invitation.ID, t.Name()))

				require.NoError(t, c.SwitchActiveHousehold(ctx, relevantHouseholdID))

				createdUsers = append(createdUsers, u)
				createdClients = append(createdClients, c)
			}

			// create recipes for meal plan
			createdMeals := []*types.Meal{}
			for i := 0; i < 3; i++ {
				createdMeal := createMealForTest(ctx, t, testClients.main, nil)
				createdMeals = append(createdMeals, createdMeal)
			}

			t.Log("creating meal plan")
			exampleMealPlan := &types.MealPlan{
				Notes:          t.Name(),
				Status:         types.AwaitingVotesMealPlanStatus,
				StartsAt:       uint64(time.Now().Add(24 * time.Hour).Unix()),
				EndsAt:         uint64(time.Now().Add(72 * time.Hour).Unix()),
				VotingDeadline: uint64(time.Now().Add(10 * time.Minute).Unix()),
				Options: []*types.MealPlanOption{
					{
						Meal:     types.Meal{ID: createdMeals[0].ID},
						Notes:    "option A",
						MealName: types.BreakfastMealName,
						Day:      time.Monday,
					},
					{
						Meal:     types.Meal{ID: createdMeals[1].ID},
						Notes:    "option B",
						MealName: types.BreakfastMealName,
						Day:      time.Monday,
					},
					{
						Meal:     types.Meal{ID: createdMeals[2].ID},
						Notes:    "option C",
						MealName: types.BreakfastMealName,
						Day:      time.Monday,
					},
				},
			}

			exampleMealPlanInput := fakes.BuildFakeMealPlanCreationRequestInputFromMealPlan(exampleMealPlan)
			createdMealPlan, err := testClients.main.CreateMealPlan(ctx, exampleMealPlanInput)
			require.NotEmpty(t, createdMealPlan.ID)
			require.NoError(t, err)

			t.Logf("meal plan %q created", createdMealPlan.ID)

			createdMealPlan, err = testClients.main.GetMealPlan(ctx, createdMealPlan.ID)
			requireNotNilAndNoProblems(t, createdMealPlan, err)
			checkMealPlanEquality(t, exampleMealPlan, createdMealPlan)

			userAVotes := &types.MealPlanOptionVoteCreationRequestInput{
				Votes: []*types.MealPlanOptionVoteCreationInput{
					{
						BelongsToMealPlanOption: createdMealPlan.Options[0].ID,
						Rank:                    0,
					},
					{
						BelongsToMealPlanOption: createdMealPlan.Options[1].ID,
						Rank:                    2,
					},
					{
						BelongsToMealPlanOption: createdMealPlan.Options[2].ID,
						Rank:                    1,
					},
				},
			}

			userBVotes := &types.MealPlanOptionVoteCreationRequestInput{
				Votes: []*types.MealPlanOptionVoteCreationInput{
					{
						BelongsToMealPlanOption: createdMealPlan.Options[0].ID,
						Rank:                    0,
					},
					{
						BelongsToMealPlanOption: createdMealPlan.Options[1].ID,
						Rank:                    1,
					},
					{
						BelongsToMealPlanOption: createdMealPlan.Options[2].ID,
						Rank:                    2,
					},
				},
			}

			createdMealPlanOptionVotesA, err := createdClients[0].CreateMealPlanOptionVote(ctx, createdMealPlan.ID, userAVotes)
			require.NoError(t, err)
			require.NotNil(t, createdMealPlanOptionVotesA)
			t.Logf("meal plan option votes created for user A")

			//checkMealPlanOptionVoteEquality(t, vote, createdMealPlanOptionVotesA)

			//createdMealPlanOptionVotesA, err = createdClients[0].GetMealPlanOptionVote(ctx, createdMealPlan.ID, vote.BelongsToMealPlanOption, createdMealPlanOptionVotesA.ID)
			//requireNotNilAndNoProblems(t, createdMealPlanOptionVotesA, err)
			//require.Equal(t, vote.BelongsToMealPlanOption, createdMealPlanOptionVotesA.BelongsToMealPlanOption)
			//checkMealPlanOptionVoteEquality(t, vote, createdMealPlanOptionVotesA)

			createdMealPlanOptionVotesB, err := createdClients[1].CreateMealPlanOptionVote(ctx, createdMealPlan.ID, userBVotes)
			require.NoError(t, err)
			require.NotNil(t, createdMealPlanOptionVotesB)
			t.Logf("meal plan option votes created for user B")

			//checkMealPlanOptionVoteEquality(t, vote, createdMealPlanOptionVotesB)

			//createdMealPlanOptionVotesB, err = createdClients[1].GetMealPlanOptionVote(ctx, createdMealPlan.ID, vote.BelongsToMealPlanOption, createdMealPlanOptionVotesB.ID)
			//requireNotNilAndNoProblems(t, createdMealPlanOptionVotesB, err)
			//require.Equal(t, vote.BelongsToMealPlanOption, createdMealPlanOptionVotesB.BelongsToMealPlanOption)
			//checkMealPlanOptionVoteEquality(t, vote, createdMealPlanOptionVotesB)

			createdMealPlan, err = testClients.main.GetMealPlan(ctx, createdMealPlan.ID)
			requireNotNilAndNoProblems(t, createdMealPlan, err)
			assert.Equal(t, types.AwaitingVotesMealPlanStatus, createdMealPlan.Status)

			createdMealPlan.VotingDeadline = uint64(time.Now().Add(-10 * time.Hour).Unix())
			require.NoError(t, dbmanager.UpdateMealPlan(ctx, createdMealPlan))

			time.Sleep(5 * time.Second)

			createdMealPlan, err = testClients.main.GetMealPlan(ctx, createdMealPlan.ID)
			requireNotNilAndNoProblems(t, createdMealPlan, err)
			assert.Equal(t, types.FinalizedMealPlanStatus, createdMealPlan.Status)

			for _, day := range allDays {
				for _, mealName := range allMealNames {
					options := byDayAndMeal(createdMealPlan.Options, day, mealName)
					if len(options) > 0 {
						selectionMade := false
						for _, opt := range options {
							if opt.Chosen {
								selectionMade = true
								break
							}
						}
						require.True(t, selectionMade)
					}
				}
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
				createdMealPlan := createMealPlanWithNotificationChannel(ctx, t, testClients.main)
				expected = append(expected, createdMealPlan)
			}

			// assert meal plan list equality
			actual, err := testClients.main.GetMealPlans(ctx, nil)
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
				assert.NoError(t, testClients.main.ArchiveMealPlan(ctx, createdMealPlan.ID))
			}
		}
	})
}
