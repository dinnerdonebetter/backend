package integration

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkUserIngredientPreferenceEquality(t *testing.T, expected, actual *types.UserIngredientPreference) {
	t.Helper()

	assert.NotZero(t, actual.ID)

	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for user ingredient preference %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.Ingredient.ID, actual.Ingredient.ID, "expected IngredientID for user ingredient preference %s to be %v, but it was %v", expected.ID, expected.Ingredient.ID, actual.Ingredient.ID)
	assert.Equal(t, expected.Rating, actual.Rating, "expected Rating for user ingredient preference %s to be %v, but it was %v", expected.ID, expected.Rating, actual.Rating)
	assert.Equal(t, expected.Allergy, actual.Allergy, "expected Allergy for user ingredient preference %s to be %v, but it was %v", expected.ID, expected.Allergy, actual.Allergy)

	assert.NotZero(t, actual.CreatedAt)
}

func (s *TestSuite) TestUserIngredientPreferences_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdValidIngredient := createValidIngredientForTest(t, ctx, testClients.admin)

			exampleUserIngredientPreference := fakes.BuildFakeUserIngredientPreference()
			exampleUserIngredientPreference.Ingredient = *createdValidIngredient
			exampleUserIngredientPreferenceInput := converters.ConvertUserIngredientPreferenceToUserIngredientPreferenceCreationRequestInput(exampleUserIngredientPreference)
			created, err := testClients.admin.CreateUserIngredientPreference(ctx, exampleUserIngredientPreferenceInput)
			require.NoError(t, err)

			require.Len(t, created, 1)
			createdUserIngredientPreference := created[0]

			checkUserIngredientPreferenceEquality(t, exampleUserIngredientPreference, createdUserIngredientPreference)

			createdUserIngredientPreferences, err := testClients.admin.GetUserIngredientPreferences(ctx, types.DefaultQueryFilter())
			requireNotNilAndNoProblems(t, createdUserIngredientPreference, err)
			checkUserIngredientPreferenceEquality(t, exampleUserIngredientPreference, createdUserIngredientPreference)
			require.NotEmpty(t, createdUserIngredientPreferences.Data, "expected to find at least one user ingredient preference")
			createdUserIngredientPreference = createdUserIngredientPreferences.Data[0]

			createdValidIngredient2 := createValidIngredientForTest(t, ctx, testClients.admin)
			newUserIngredientPreference := fakes.BuildFakeUserIngredientPreference()
			newUserIngredientPreference.Ingredient = *createdValidIngredient2
			createdUserIngredientPreference.Update(converters.ConvertUserIngredientPreferenceToUserIngredientPreferenceUpdateRequestInput(newUserIngredientPreference))
			assert.NoError(t, testClients.admin.UpdateUserIngredientPreference(ctx, createdUserIngredientPreference))

			newResults, err := testClients.admin.GetUserIngredientPreferences(ctx, types.DefaultQueryFilter())
			requireNotNilAndNoProblems(t, newResults, err)
			require.NotEmpty(t, newResults.Data, "expected to find at least one user ingredient preference")
			actual := newResults.Data[0]

			// assert user ingredient preference equality
			checkUserIngredientPreferenceEquality(t, newUserIngredientPreference, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			assert.NoError(t, testClients.admin.ArchiveUserIngredientPreference(ctx, createdUserIngredientPreference.ID))
		}
	})
}

func (s *TestSuite) TestUserIngredientPreferences_CreatedFromIngredientGroup() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdValidIngredientGroup := createValidIngredientGroupForTest(t, ctx, nil, testClients.admin)
			logJSON(t, createdValidIngredientGroup)

			exampleUserIngredientPreference := fakes.BuildFakeUserIngredientPreference()
			exampleUserIngredientPreferenceInput := converters.ConvertUserIngredientPreferenceToUserIngredientPreferenceCreationRequestInput(exampleUserIngredientPreference)
			exampleUserIngredientPreferenceInput.ValidIngredientGroupID = createdValidIngredientGroup.ID
			exampleUserIngredientPreferenceInput.ValidIngredientID = ""

			logJSON(t, exampleUserIngredientPreferenceInput)
			logJSON(t, createdValidIngredientGroup)

			created, err := testClients.admin.CreateUserIngredientPreference(ctx, exampleUserIngredientPreferenceInput)
			require.NoError(t, err)

			require.Equal(t, len(created), len(createdValidIngredientGroup.Members))

			for _, createdUserIngredientPreference := range created {
				assert.NoError(t, testClients.admin.ArchiveUserIngredientPreference(ctx, createdUserIngredientPreference.ID))
			}
		}
	})
}

func (s *TestSuite) TestUserIngredientPreferences_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			var expected []*types.UserIngredientPreference
			for i := 0; i < 5; i++ {
				createdValidIngredient := createValidIngredientForTest(t, ctx, testClients.admin)
				exampleUserIngredientPreference := fakes.BuildFakeUserIngredientPreference()
				exampleUserIngredientPreference.Ingredient = *createdValidIngredient
				exampleUserIngredientPreferenceInput := converters.ConvertUserIngredientPreferenceToUserIngredientPreferenceCreationRequestInput(exampleUserIngredientPreference)
				created, createdUserIngredientPreferenceErr := testClients.admin.CreateUserIngredientPreference(ctx, exampleUserIngredientPreferenceInput)
				require.NoError(t, createdUserIngredientPreferenceErr)

				require.Len(t, created, 1)
				createdUserIngredientPreference := created[0]

				checkUserIngredientPreferenceEquality(t, exampleUserIngredientPreference, createdUserIngredientPreference)

				expected = append(expected, createdUserIngredientPreference)
			}

			// assert user ingredient preference list equality
			actual, err := testClients.admin.GetUserIngredientPreferences(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			for _, createdUserIngredientPreference := range expected {
				assert.NoError(t, testClients.admin.ArchiveUserIngredientPreference(ctx, createdUserIngredientPreference.ID))
			}
		}
	})
}
