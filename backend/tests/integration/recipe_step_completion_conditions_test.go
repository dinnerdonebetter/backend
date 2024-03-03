package integration

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkRecipeStepCompletionConditionEquality(t *testing.T, expected, actual *types.RecipeStepCompletionCondition) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.IngredientState, actual.IngredientState, "expected IngredientState for recipe step completion condition %s to be %v, but it was %v", expected.ID, expected.IngredientState, actual.IngredientState)
	assert.Equal(t, expected.BelongsToRecipeStep, actual.BelongsToRecipeStep, "expected BelongsToRecipeStep for recipe step completion condition %s to be %v, but it was %v", expected.ID, expected.BelongsToRecipeStep, actual.BelongsToRecipeStep)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for recipe step completion condition %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.Ingredients, actual.Ingredients, "expected Ingredients for recipe step completion condition %s to be %v, but it was %v", expected.ID, expected.Ingredients, actual.Ingredients)
	assert.Equal(t, expected.Optional, actual.Optional, "expected Optional for recipe step completion condition %s to be %v, but it was %v", expected.ID, expected.Optional, actual.Optional)
	assert.NotZero(t, actual.CreatedAt)
}

func (s *TestSuite) TestRecipeStepCompletionConditions_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.user, nil)

			createdRecipeStep := createdRecipe.Steps[0]
			require.NotEmpty(t, createdRecipeStep.ID, "created recipe step ID must not be empty")

			// create ingredient state
			createdValidIngredientState := createValidIngredientStateForTest(t, ctx, testClients.admin)

			input := &types.RecipeStepCompletionConditionForExistingRecipeCreationRequestInput{
				IngredientStateID:   createdValidIngredientState.ID,
				BelongsToRecipeStep: createdRecipeStep.ID,
				Notes:               t.Name(),
				Optional:            false,
				Ingredients: []*types.RecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput{
					{
						RecipeStepIngredient: createdRecipeStep.Ingredients[0].ID,
					},
				},
			}

			createdRecipeStepCompletionCondition, err := testClients.user.CreateRecipeStepCompletionCondition(ctx, createdRecipe.ID, createdRecipeStep.ID, input)
			requireNotNilAndNoProblems(t, createdRecipeStepCompletionCondition, err)

			createdRecipeStepCompletionCondition.Notes = t.Name() + " updated"

			require.NoError(t, testClients.user.UpdateRecipeStepCompletionCondition(ctx, createdRecipe.ID, createdRecipeStepCompletionCondition))

			actual, err := testClients.user.GetRecipeStepCompletionCondition(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepCompletionCondition.ID)
			requireNotNilAndNoProblems(t, actual, err)
			actual.IngredientState = types.ValidIngredientState{
				ID: createdRecipeStepCompletionCondition.IngredientState.ID,
			}
			actual.CreatedAt = createdRecipeStepCompletionCondition.CreatedAt
			for i := range actual.Ingredients {
				actual.Ingredients[i].CreatedAt = createdRecipeStepCompletionCondition.Ingredients[i].CreatedAt
			}

			// assert recipe step completion condition equality
			checkRecipeStepCompletionConditionEquality(t, createdRecipeStepCompletionCondition, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			// assert recipe step completion condition list functionality works
			listResponse, err := testClients.user.GetRecipeStepCompletionConditions(ctx, createdRecipe.ID, createdRecipeStep.ID, types.DefaultQueryFilter())
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				1 <= len(listResponse.Data),
				"expected %d to be <= %d",
				1,
				len(listResponse.Data),
			)

			assert.NoError(t, testClients.user.ArchiveRecipeStepCompletionCondition(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepCompletionCondition.ID))

			assert.NoError(t, testClients.user.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			assert.NoError(t, testClients.user.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}
