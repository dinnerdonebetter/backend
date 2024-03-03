package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkValidIngredientEquality(t *testing.T, expected, actual *types.ValidIngredient) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for valid ingredient %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for valid ingredient %s to be %v, but it was %v", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.Warning, actual.Warning, "expected Warning for valid ingredient %s to be %v, but it was %v", expected.ID, expected.Warning, actual.Warning)
	assert.Equal(t, expected.ContainsEgg, actual.ContainsEgg, "expected ContainsEgg for valid ingredient %s to be %v, but it was %v", expected.ID, expected.ContainsEgg, actual.ContainsEgg)
	assert.Equal(t, expected.ContainsDairy, actual.ContainsDairy, "expected ContainsDairy for valid ingredient %s to be %v, but it was %v", expected.ID, expected.ContainsDairy, actual.ContainsDairy)
	assert.Equal(t, expected.ContainsPeanut, actual.ContainsPeanut, "expected ContainsPeanut for valid ingredient %s to be %v, but it was %v", expected.ID, expected.ContainsPeanut, actual.ContainsPeanut)
	assert.Equal(t, expected.ContainsTreeNut, actual.ContainsTreeNut, "expected ContainsTreeNut for valid ingredient %s to be %v, but it was %v", expected.ID, expected.ContainsTreeNut, actual.ContainsTreeNut)
	assert.Equal(t, expected.ContainsSoy, actual.ContainsSoy, "expected ContainsSoy for valid ingredient %s to be %v, but it was %v", expected.ID, expected.ContainsSoy, actual.ContainsSoy)
	assert.Equal(t, expected.ContainsWheat, actual.ContainsWheat, "expected ContainsWheat for valid ingredient %s to be %v, but it was %v", expected.ID, expected.ContainsWheat, actual.ContainsWheat)
	assert.Equal(t, expected.ContainsShellfish, actual.ContainsShellfish, "expected ContainsShellfish for valid ingredient %s to be %v, but it was %v", expected.ID, expected.ContainsShellfish, actual.ContainsShellfish)
	assert.Equal(t, expected.ContainsSesame, actual.ContainsSesame, "expected ContainsSesame for valid ingredient %s to be %v, but it was %v", expected.ID, expected.ContainsSesame, actual.ContainsSesame)
	assert.Equal(t, expected.ContainsFish, actual.ContainsFish, "expected ContainsFish for valid ingredient %s to be %v, but it was %v", expected.ID, expected.ContainsFish, actual.ContainsFish)
	assert.Equal(t, expected.ContainsGluten, actual.ContainsGluten, "expected ContainsGluten for valid ingredient %s to be %v, but it was %v", expected.ID, expected.ContainsGluten, actual.ContainsGluten)
	assert.Equal(t, expected.AnimalFlesh, actual.AnimalFlesh, "expected AnimalFlesh for valid ingredient %s to be %v, but it was %v", expected.ID, expected.AnimalFlesh, actual.AnimalFlesh)
	assert.Equal(t, expected.IsMeasuredVolumetrically, actual.IsMeasuredVolumetrically, "expected IsMeasuredVolumetrically for valid ingredient %s to be %v, but it was %v", expected.ID, expected.IsMeasuredVolumetrically, actual.IsMeasuredVolumetrically)
	assert.Equal(t, expected.IsLiquid, actual.IsLiquid, "expected IsLiquid for valid ingredient %s to be %v, but it was %v", expected.ID, expected.IsLiquid, actual.IsLiquid)
	assert.Equal(t, expected.IconPath, actual.IconPath, "expected IconPath for valid ingredient %s to be %v, but it was %v", expected.ID, expected.IconPath, actual.IconPath)
	assert.Equal(t, expected.PluralName, actual.PluralName, "expected PluralName for valid ingredient %s to be %v, but it was %v", expected.ID, expected.PluralName, actual.PluralName)
	assert.Equal(t, expected.AnimalDerived, actual.AnimalDerived, "expected AnimalDerived for valid ingredient %s to be %v, but it was %v", expected.ID, expected.AnimalDerived, actual.AnimalDerived)
	assert.Equal(t, expected.RestrictToPreparations, actual.RestrictToPreparations, "expected RestrictToPreparations for valid ingredient %s to be %v, but it was %v", expected.ID, expected.RestrictToPreparations, actual.RestrictToPreparations)
	assert.Equal(t, expected.MinimumIdealStorageTemperatureInCelsius, actual.MinimumIdealStorageTemperatureInCelsius, "expected MinimumIdealStorageTemperatureInCelsius for valid ingredient %s to be %v, but it was %v", expected.ID, expected.MinimumIdealStorageTemperatureInCelsius, actual.MinimumIdealStorageTemperatureInCelsius)
	assert.Equal(t, expected.MaximumIdealStorageTemperatureInCelsius, actual.MaximumIdealStorageTemperatureInCelsius, "expected MaximumIdealStorageTemperatureInCelsius for valid ingredient %s to be %v, but it was %v", expected.ID, expected.MaximumIdealStorageTemperatureInCelsius, actual.MaximumIdealStorageTemperatureInCelsius)
	assert.Equal(t, expected.StorageInstructions, actual.StorageInstructions, "expected StorageInstructions for valid ingredient %s to be %v, but it was %v", expected.ID, expected.StorageInstructions, actual.StorageInstructions)
	assert.Equal(t, expected.Slug, actual.Slug, "expected Slug for valid ingredient %s to be %v, but it was %v", expected.ID, expected.Slug, actual.Slug)
	assert.Equal(t, expected.ShoppingSuggestions, actual.ShoppingSuggestions, "expected ShoppingSuggestions for valid ingredient %s to be %v, but it was %v", expected.ID, expected.ShoppingSuggestions, actual.ShoppingSuggestions)
	assert.Equal(t, expected.IsStarch, actual.IsStarch, "expected IsStarch for valid ingredient %s to be %v, but it was %v", expected.ID, expected.IsStarch, actual.IsStarch)
	assert.Equal(t, expected.IsProtein, actual.IsProtein, "expected IsProtein for valid ingredient %s to be %v, but it was %v", expected.ID, expected.IsProtein, actual.IsProtein)
	assert.Equal(t, expected.IsGrain, actual.IsGrain, "expected IsGrain for valid ingredient %s to be %v, but it was %v", expected.ID, expected.IsGrain, actual.IsGrain)
	assert.Equal(t, expected.IsFruit, actual.IsFruit, "expected IsFruit for valid ingredient %s to be %v, but it was %v", expected.ID, expected.IsFruit, actual.IsFruit)
	assert.Equal(t, expected.IsSalt, actual.IsSalt, "expected IsSalt for valid ingredient %s to be %v, but it was %v", expected.ID, expected.IsSalt, actual.IsSalt)
	assert.Equal(t, expected.IsFat, actual.IsFat, "expected IsFat for valid ingredient %s to be %v, but it was %v", expected.ID, expected.IsFat, actual.IsFat)
	assert.Equal(t, expected.IsAcid, actual.IsAcid, "expected IsAcid for valid ingredient %s to be %v, but it was %v", expected.ID, expected.IsAcid, actual.IsAcid)
	assert.Equal(t, expected.IsHeat, actual.IsHeat, "expected IsHeat for valid ingredient %s to be %v, but it was %v", expected.ID, expected.IsHeat, actual.IsHeat)
	assert.Equal(t, expected.ContainsAlcohol, actual.ContainsAlcohol, "expected ContainsAlcohol for valid ingredient %s to be %v, but it was %v", expected.ID, expected.ContainsAlcohol, actual.ContainsAlcohol)
	assert.NotZero(t, actual.CreatedAt)
}

func createValidIngredientForTest(t *testing.T, ctx context.Context, adminClient *apiclient.Client) *types.ValidIngredient {
	t.Helper()

	exampleValidIngredient := fakes.BuildFakeValidIngredient()
	exampleValidIngredientInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(exampleValidIngredient)
	createdValidIngredient, err := adminClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
	require.NoError(t, err)
	checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

	createdValidIngredient, err = adminClient.GetValidIngredient(ctx, createdValidIngredient.ID)
	requireNotNilAndNoProblems(t, createdValidIngredient, err)
	checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

	return createdValidIngredient
}

func (s *TestSuite) TestValidIngredients_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdValidIngredient := createValidIngredientForTest(t, ctx, testClients.admin)

			newValidIngredient := fakes.BuildFakeValidIngredient()
			createdValidIngredient.Update(converters.ConvertValidIngredientToValidIngredientUpdateRequestInput(newValidIngredient))
			assert.NoError(t, testClients.admin.UpdateValidIngredient(ctx, createdValidIngredient))

			actual, err := testClients.admin.GetValidIngredient(ctx, createdValidIngredient.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid ingredient equality
			checkValidIngredientEquality(t, newValidIngredient, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			assert.NoError(t, testClients.admin.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
		}
	})
}

func (s *TestSuite) TestValidIngredients_GetRandom() {
	s.runForEachClient("should be able to get a random valid ingredient", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(exampleValidIngredient)
			createdValidIngredient, err := testClients.admin.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			createdValidIngredient, err = testClients.admin.GetRandomValidIngredient(ctx)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)

			assert.NoError(t, testClients.admin.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
		}
	})
}

func (s *TestSuite) TestValidIngredients_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			var expected []*types.ValidIngredient
			for i := 0; i < 5; i++ {
				exampleValidIngredient := fakes.BuildFakeValidIngredient()
				exampleValidIngredientInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(exampleValidIngredient)
				createdValidIngredient, createdValidIngredientErr := testClients.admin.CreateValidIngredient(ctx, exampleValidIngredientInput)
				require.NoError(t, createdValidIngredientErr)

				checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

				expected = append(expected, createdValidIngredient)
			}

			// assert valid ingredient list equality
			actual, err := testClients.admin.GetValidIngredients(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			for _, createdValidIngredient := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
			}
		}
	})
}

func (s *TestSuite) TestValidIngredients_Searching() {
	s.runForEachClient("should be able to be search for valid ingredients", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			var expected []*types.ValidIngredient
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredient.Name = fmt.Sprintf("example_%s", testClients.authType)
			searchQuery := exampleValidIngredient.Name
			for i := 0; i < 5; i++ {
				exampleValidIngredient.Name = fmt.Sprintf("%s %d", searchQuery, i)
				exampleValidIngredientInput := converters.ConvertValidIngredientToValidIngredientCreationRequestInput(exampleValidIngredient)
				createdValidIngredient, createdValidIngredientErr := testClients.admin.CreateValidIngredient(ctx, exampleValidIngredientInput)
				require.NoError(t, createdValidIngredientErr)
				checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

				expected = append(expected, createdValidIngredient)
			}

			exampleLimit := uint8(20)

			// assert valid ingredient list equality
			actual, err := testClients.admin.SearchValidIngredients(ctx, searchQuery, exampleLimit)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			for _, createdValidIngredient := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
			}
		}
	})
}
