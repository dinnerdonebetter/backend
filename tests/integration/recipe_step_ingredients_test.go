package integration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"
)

func checkRecipeStepIngredientEquality(t *testing.T, expected, actual *types.RecipeStepIngredient) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.IngredientID, actual.IngredientID, "expected IngredientID for recipe step ingredient %s to be %v, but it was %v", expected.ID, expected.IngredientID, actual.IngredientID)
	assert.Equal(t, expected.QuantityType, actual.QuantityType, "expected QuantityType for recipe step ingredient %s to be %v, but it was %v", expected.ID, expected.QuantityType, actual.QuantityType)
	assert.Equal(t, expected.QuantityValue, actual.QuantityValue, "expected QuantityValue for recipe step ingredient %s to be %v, but it was %v", expected.ID, expected.QuantityValue, actual.QuantityValue)
	assert.Equal(t, expected.QuantityNotes, actual.QuantityNotes, "expected QuantityNotes for recipe step ingredient %s to be %v, but it was %v", expected.ID, expected.QuantityNotes, actual.QuantityNotes)
	assert.Equal(t, expected.ProductOfRecipeStep, actual.ProductOfRecipeStep, "expected ProductOfRecipeStep for recipe step ingredient %s to be %v, but it was %v", expected.ID, expected.ProductOfRecipeStep, actual.ProductOfRecipeStep)
	assert.Equal(t, expected.IngredientNotes, actual.IngredientNotes, "expected IngredientNotes for recipe step ingredient %s to be %v, but it was %v", expected.ID, expected.IngredientNotes, actual.IngredientNotes)
	assert.NotZero(t, actual.CreatedOn)
}

// convertRecipeStepIngredientToRecipeStepIngredientUpdateInput creates an RecipeStepIngredientUpdateRequestInput struct from a recipe step ingredient.
func convertRecipeStepIngredientToRecipeStepIngredientUpdateInput(x *types.RecipeStepIngredient) *types.RecipeStepIngredientUpdateRequestInput {
	return &types.RecipeStepIngredientUpdateRequestInput{
		IngredientID:    x.IngredientID,
		QuantityType:    x.QuantityType,
		QuantityValue:   x.QuantityValue,
		QuantityNotes:   x.QuantityNotes,
		ProductOfRecipe: x.ProductOfRecipeStep,
		IngredientNotes: x.IngredientNotes,
	}
}

func (s *TestSuite) TestRecipeStepIngredients_CompleteLifecycle() {
	s.runForCookieClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			stopChan := make(chan bool, 1)
			notificationsChan, err := testClients.main.SubscribeToDataChangeNotifications(ctx, stopChan)
			require.NotNil(t, notificationsChan)
			require.NoError(t, err)

			var n *types.DataChangeMessage

			createdValidIngredient, _, createdRecipe := createRecipeWithNotificationChannel(ctx, t, notificationsChan, testClients.main)

			var (
				createdRecipeStepID,
				createdRecipeStepIngredientID string
			)
			for _, step := range createdRecipe.Steps {
				createdRecipeStepID = step.ID
				for _, ingredient := range step.Ingredients {
					createdRecipeStepIngredientID = ingredient.ID
					break
				}
			}

			t.Log("fetching changed recipe step ingredient")
			createdRecipeStepIngredient, err := testClients.main.GetRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepIngredientID)
			requireNotNilAndNoProblems(t, createdRecipeStepIngredient, err)

			t.Log("changing recipe step ingredient")
			newRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
			newRecipeStepIngredient.BelongsToRecipeStep = createdRecipeStepID
			newRecipeStepIngredient.IngredientID = &createdValidIngredient.ID
			createdRecipeStepIngredient.Update(convertRecipeStepIngredientToRecipeStepIngredientUpdateInput(newRecipeStepIngredient))
			assert.NoError(t, testClients.main.UpdateRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStepIngredient))

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.RecipeStepIngredientDataType)

			t.Log("fetching changed recipe step ingredient")
			actual, err := testClients.main.GetRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepIngredientID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step ingredient equality
			checkRecipeStepIngredientEquality(t, newRecipeStepIngredient, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up recipe step ingredient")
			assert.NoError(t, testClients.main.ArchiveRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepIngredientID))

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})

	s.runForPASETOClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			var checkFunc func() bool
			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdValidIngredient, _, createdRecipe := createRecipeWithPolling(ctx, t, testClients.main)

			createdRecipeStepID := createdRecipe.Steps[0].ID
			var createdRecipeStepIngredientID string
			for _, step := range createdRecipe.Steps {
				createdRecipeStepID = step.ID
				for _, ingredient := range step.Ingredients {
					createdRecipeStepIngredientID = ingredient.ID
					break
				}
				break
			}

			var (
				createdRecipeStepIngredient *types.RecipeStepIngredient
				err                         error
			)
			checkFunc = func() bool {
				createdRecipeStepIngredient, err = testClients.main.GetRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepIngredientID)
				return assert.NotNil(t, createdRecipeStepIngredient) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)

			// change recipe step ingredient
			newRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
			newRecipeStepIngredient.BelongsToRecipeStep = createdRecipeStepID
			newRecipeStepIngredient.IngredientID = &createdValidIngredient.ID
			createdRecipeStepIngredient.Update(convertRecipeStepIngredientToRecipeStepIngredientUpdateInput(newRecipeStepIngredient))
			assert.NoError(t, testClients.main.UpdateRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStepIngredient))

			time.Sleep(2 * time.Second)

			// retrieve changed recipe step ingredient
			var actual *types.RecipeStepIngredient
			checkFunc = func() bool {
				actual, err = testClients.main.GetRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepIngredientID)
				return assert.NotNil(t, createdRecipeStepIngredient) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)

			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step ingredient equality
			checkRecipeStepIngredientEquality(t, newRecipeStepIngredient, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up recipe step ingredient")
			assert.NoError(t, testClients.main.ArchiveRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepIngredientID))

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeStepIngredients_Listing() {
	s.runForCookieClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			stopChan := make(chan bool, 1)
			notificationsChan, err := testClients.main.SubscribeToDataChangeNotifications(ctx, stopChan)
			require.NotNil(t, notificationsChan)
			require.NoError(t, err)

			var n *types.DataChangeMessage

			createdValidIngredient, _, createdRecipe := createRecipeWithNotificationChannel(ctx, t, notificationsChan, testClients.main)

			var (
				createdRecipeStepID string
			)
			for _, step := range createdRecipe.Steps {
				createdRecipeStepID = step.ID
				break
			}

			t.Log("creating recipe step ingredients")
			var expected []*types.RecipeStepIngredient
			for i := 0; i < 5; i++ {
				exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
				exampleRecipeStepIngredient.BelongsToRecipeStep = createdRecipeStepID
				exampleRecipeStepIngredient.IngredientID = &createdValidIngredient.ID
				exampleRecipeStepIngredientInput := fakes.BuildFakeRecipeStepIngredientCreationRequestInputFromRecipeStepIngredient(exampleRecipeStepIngredient)
				createdRecipeStepIngredientID, createdRecipeStepIngredientErr := testClients.main.CreateRecipeStepIngredient(ctx, createdRecipe.ID, exampleRecipeStepIngredientInput)
				require.NoError(t, createdRecipeStepIngredientErr)
				t.Logf("recipe step ingredient %q created", createdRecipeStepIngredientID)

				n = <-notificationsChan
				assert.Equal(t, n.DataType, types.RecipeStepIngredientDataType)
				require.NotNil(t, n.RecipeStepIngredient)
				checkRecipeStepIngredientEquality(t, exampleRecipeStepIngredient, n.RecipeStepIngredient)

				createdRecipeStepIngredient, createdRecipeStepIngredientErr := testClients.main.GetRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepIngredientID)
				requireNotNilAndNoProblems(t, createdRecipeStepIngredient, createdRecipeStepIngredientErr)
				require.Equal(t, createdRecipeStepID, createdRecipeStepIngredient.BelongsToRecipeStep)

				expected = append(expected, createdRecipeStepIngredient)
			}

			// assert recipe step ingredient list equality
			actual, err := testClients.main.GetRecipeStepIngredients(ctx, createdRecipe.ID, createdRecipeStepID, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.RecipeStepIngredients),
				"expected %d to be <= %d",
				len(expected),
				len(actual.RecipeStepIngredients),
			)

			t.Log("cleaning up")
			for _, createdRecipeStepIngredient := range expected {
				assert.NoError(t, testClients.main.ArchiveRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepIngredient.ID))
			}

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})

	s.runForPASETOClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			var checkFunc func() bool
			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdValidIngredient, _, createdRecipe := createRecipeWithPolling(ctx, t, testClients.main)

			var (
				createdRecipeStepID string
			)
			for _, step := range createdRecipe.Steps {
				createdRecipeStepID = step.ID
				break
			}

			t.Log("creating recipe step ingredients")
			var expected []*types.RecipeStepIngredient
			for i := 0; i < 5; i++ {
				exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
				exampleRecipeStepIngredient.BelongsToRecipeStep = createdRecipeStepID
				exampleRecipeStepIngredient.IngredientID = &createdValidIngredient.ID
				exampleRecipeStepIngredientInput := fakes.BuildFakeRecipeStepIngredientCreationRequestInputFromRecipeStepIngredient(exampleRecipeStepIngredient)
				createdRecipeStepIngredientID, createdRecipeStepIngredientErr := testClients.main.CreateRecipeStepIngredient(ctx, createdRecipe.ID, exampleRecipeStepIngredientInput)
				require.NoError(t, createdRecipeStepIngredientErr)

				var createdRecipeStepIngredient *types.RecipeStepIngredient
				checkFunc = func() bool {
					createdRecipeStepIngredient, createdRecipeStepIngredientErr = testClients.main.GetRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepIngredientID)
					return assert.NotNil(t, createdRecipeStepIngredient) && assert.NoError(t, createdRecipeStepIngredientErr)
				}
				assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
				checkRecipeStepIngredientEquality(t, exampleRecipeStepIngredient, createdRecipeStepIngredient)

				expected = append(expected, createdRecipeStepIngredient)
			}

			// assert recipe step ingredient list equality
			actual, err := testClients.main.GetRecipeStepIngredients(ctx, createdRecipe.ID, createdRecipeStepID, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.RecipeStepIngredients),
				"expected %d to be <= %d",
				len(expected),
				len(actual.RecipeStepIngredients),
			)

			t.Log("cleaning up")
			for _, createdRecipeStepIngredient := range expected {
				assert.NoError(t, testClients.main.ArchiveRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepIngredient.ID))
			}

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}
