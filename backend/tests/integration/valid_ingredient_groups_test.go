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

func checkValidIngredientGroupEquality(t *testing.T, expected, actual *types.ValidIngredientGroup) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for valid ingredient group %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for valid ingredient group %s to be %v, but it was %v", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.Slug, actual.Slug, "expected Slug for valid ingredient group %s to be %v, but it was %v", expected.ID, expected.Slug, actual.Slug)
	assert.NotZero(t, actual.CreatedAt)
}

func createValidIngredientGroupForTest(t *testing.T, ctx context.Context, creationInput *types.ValidIngredientGroup, adminClient *apiclient.Client) *types.ValidIngredientGroup {
	t.Helper()

	createdValidIngredients := []*types.ValidIngredient{}
	for i := 0; i < 3; i++ {
		createdValidIngredients = append(createdValidIngredients, createValidIngredientForTest(t, ctx, adminClient))
	}

	exampleValidIngredientGroup := creationInput
	if exampleValidIngredientGroup == nil {
		exampleValidIngredientGroup = fakes.BuildFakeValidIngredientGroup()
	}
	exampleValidIngredientGroupInput := converters.ConvertValidIngredientGroupToValidIngredientGroupCreationRequestInput(exampleValidIngredientGroup)

	members := []*types.ValidIngredientGroupMemberCreationRequestInput{}
	for _, validIngredient := range createdValidIngredients {
		members = append(members, &types.ValidIngredientGroupMemberCreationRequestInput{
			ValidIngredientID: validIngredient.ID,
		})
	}
	exampleValidIngredientGroupInput.Members = members

	createdValidIngredientGroup, err := adminClient.CreateValidIngredientGroup(ctx, exampleValidIngredientGroupInput)
	require.NoError(t, err)

	checkValidIngredientGroupEquality(t, exampleValidIngredientGroup, createdValidIngredientGroup)

	createdValidIngredientGroup, err = adminClient.GetValidIngredientGroup(ctx, createdValidIngredientGroup.ID)
	requireNotNilAndNoProblems(t, createdValidIngredientGroup, err)
	checkValidIngredientGroupEquality(t, exampleValidIngredientGroup, createdValidIngredientGroup)

	require.Equal(t, len(createdValidIngredientGroup.Members), len(createdValidIngredients))

	return createdValidIngredientGroup
}

func (s *TestSuite) TestValidIngredientGroups_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdValidIngredientGroup := createValidIngredientGroupForTest(t, ctx, nil, testClients.admin)

			newValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()
			createdValidIngredientGroup.Update(converters.ConvertValidIngredientGroupToValidIngredientGroupUpdateRequestInput(newValidIngredientGroup))
			assert.NoError(t, testClients.admin.UpdateValidIngredientGroup(ctx, createdValidIngredientGroup))

			actual, err := testClients.admin.GetValidIngredientGroup(ctx, createdValidIngredientGroup.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid ingredient group equality
			checkValidIngredientGroupEquality(t, newValidIngredientGroup, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			assert.NoError(t, testClients.admin.ArchiveValidIngredientGroup(ctx, createdValidIngredientGroup.ID))
		}
	})
}

func (s *TestSuite) TestValidIngredientGroups_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			var expected []*types.ValidIngredientGroup
			for i := 0; i < 5; i++ {
				createdValidIngredientGroup := createValidIngredientGroupForTest(t, ctx, nil, testClients.admin)

				expected = append(expected, createdValidIngredientGroup)
			}

			// assert valid ingredient group list equality
			actual, err := testClients.admin.GetValidIngredientGroups(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			for _, createdValidIngredientGroup := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidIngredientGroup(ctx, createdValidIngredientGroup.ID))
			}
		}
	})
}

func (s *TestSuite) TestValidIngredientGroups_Searching() {
	s.runForEachClient("should be able to be search for valid ingredient groups", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			var expected []*types.ValidIngredientGroup
			exampleValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()
			exampleValidIngredientGroup.Name = fmt.Sprintf("example_%s", testClients.authType)
			searchQuery := exampleValidIngredientGroup.Name
			for i := 0; i < 5; i++ {
				exampleValidIngredientGroup.Name = fmt.Sprintf("%s %d", searchQuery, i)
				createdValidIngredientGroup := createValidIngredientGroupForTest(t, ctx, exampleValidIngredientGroup, testClients.admin)
				expected = append(expected, createdValidIngredientGroup)
			}

			exampleLimit := uint8(20)

			// assert valid ingredient group list equality
			actual, err := testClients.admin.SearchValidIngredientGroups(ctx, searchQuery, exampleLimit)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual),
				"expected %d to be <= %d",
				len(expected),
				len(actual),
			)

			for _, createdValidIngredientGroup := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidIngredientGroup(ctx, createdValidIngredientGroup.ID))
			}
		}
	})
}
