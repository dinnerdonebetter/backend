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

func checkValidIngredientStateEquality(t *testing.T, expected, actual *types.ValidIngredientState) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for valid ingredient state %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for valid ingredient state %s to be %v, but it was %v", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.IconPath, actual.IconPath, "expected IconPath for valid ingredient state %s to be %v, but it was %v", expected.ID, expected.IconPath, actual.IconPath)
	assert.Equal(t, expected.PastTense, actual.PastTense, "expected PastTense for valid ingredient state %s to be %v, but it was %v", expected.ID, expected.PastTense, actual.PastTense)
	assert.Equal(t, expected.Slug, actual.Slug, "expected Slug for valid ingredient state %s to be %v, but it was %v", expected.ID, expected.Slug, actual.Slug)
	assert.Equal(t, expected.AttributeType, actual.AttributeType, "expected AttributeType for valid ingredient state %s to be %v, but it was %v", expected.ID, expected.AttributeType, actual.AttributeType)
	assert.NotZero(t, actual.CreatedAt)
}

func createValidIngredientStateForTest(t *testing.T, ctx context.Context, adminClient *apiclient.Client) *types.ValidIngredientState {
	t.Helper()

	exampleValidIngredientState := fakes.BuildFakeValidIngredientState()
	exampleValidIngredientStateInput := converters.ConvertValidIngredientStateToValidIngredientStateCreationRequestInput(exampleValidIngredientState)
	createdValidIngredientState, err := adminClient.CreateValidIngredientState(ctx, exampleValidIngredientStateInput)
	require.NoError(t, err)
	checkValidIngredientStateEquality(t, exampleValidIngredientState, createdValidIngredientState)

	createdValidIngredientState, err = adminClient.GetValidIngredientState(ctx, createdValidIngredientState.ID)
	requireNotNilAndNoProblems(t, createdValidIngredientState, err)
	checkValidIngredientStateEquality(t, exampleValidIngredientState, createdValidIngredientState)

	return createdValidIngredientState
}

func (s *TestSuite) TestValidIngredientStates_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdValidIngredientState := createValidIngredientStateForTest(t, ctx, testClients.admin)

			newValidIngredientState := fakes.BuildFakeValidIngredientState()
			createdValidIngredientState.Update(converters.ConvertValidIngredientStateToValidIngredientStateUpdateRequestInput(newValidIngredientState))
			assert.NoError(t, testClients.admin.UpdateValidIngredientState(ctx, createdValidIngredientState))

			actual, err := testClients.admin.GetValidIngredientState(ctx, createdValidIngredientState.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid ingredient state equality
			checkValidIngredientStateEquality(t, newValidIngredientState, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			assert.NoError(t, testClients.admin.ArchiveValidIngredientState(ctx, createdValidIngredientState.ID))
		}
	})
}

func (s *TestSuite) TestValidIngredientStates_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			var expected []*types.ValidIngredientState
			for i := 0; i < 5; i++ {
				exampleValidIngredientState := fakes.BuildFakeValidIngredientState()
				exampleValidIngredientStateInput := converters.ConvertValidIngredientStateToValidIngredientStateCreationRequestInput(exampleValidIngredientState)
				createdValidIngredientState, createdValidIngredientStateErr := testClients.admin.CreateValidIngredientState(ctx, exampleValidIngredientStateInput)
				require.NoError(t, createdValidIngredientStateErr)

				checkValidIngredientStateEquality(t, exampleValidIngredientState, createdValidIngredientState)

				expected = append(expected, createdValidIngredientState)
			}

			// assert valid ingredient state list equality
			actual, err := testClients.admin.GetValidIngredientStates(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			for _, createdValidIngredientState := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidIngredientState(ctx, createdValidIngredientState.ID))
			}
		}
	})
}

func (s *TestSuite) TestValidIngredientStates_Searching() {
	s.runForEachClient("should be able to be search for valid ingredient states", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			var expected []*types.ValidIngredientState
			exampleValidIngredientState := fakes.BuildFakeValidIngredientState()
			exampleValidIngredientState.Name = fmt.Sprintf("example_%s", testClients.authType)
			searchQuery := exampleValidIngredientState.Name
			for i := 0; i < 5; i++ {
				exampleValidIngredientState.Name = fmt.Sprintf("%s %d", searchQuery, i)
				exampleValidIngredientStateInput := converters.ConvertValidIngredientStateToValidIngredientStateCreationRequestInput(exampleValidIngredientState)
				createdValidIngredientState, createdValidIngredientStateErr := testClients.admin.CreateValidIngredientState(ctx, exampleValidIngredientStateInput)
				require.NoError(t, createdValidIngredientStateErr)
				checkValidIngredientStateEquality(t, exampleValidIngredientState, createdValidIngredientState)

				expected = append(expected, createdValidIngredientState)
			}

			exampleLimit := uint8(20)

			// assert valid ingredient state list equality
			actual, err := testClients.admin.SearchValidIngredientStates(ctx, searchQuery, exampleLimit)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual),
				"expected %d to be <= %d",
				len(expected),
				len(actual),
			)

			for _, createdValidIngredientState := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidIngredientState(ctx, createdValidIngredientState.ID))
			}
		}
	})
}
