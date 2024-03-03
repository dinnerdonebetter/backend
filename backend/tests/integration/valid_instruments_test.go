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

func checkValidInstrumentEquality(t *testing.T, expected, actual *types.ValidInstrument) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for valid instrument %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for valid instrument %s to be %v, but it was %v", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.IconPath, actual.IconPath, "expected IconPath for valid instrument %s to be %v, but it was %v", expected.ID, expected.IconPath, actual.IconPath)
	assert.Equal(t, expected.UsableForStorage, actual.UsableForStorage, "expected UsableForStorage for valid instrument %s to be %v, but it was %v", expected.ID, expected.UsableForStorage, actual.UsableForStorage)
	assert.Equal(t, expected.DisplayInSummaryLists, actual.DisplayInSummaryLists, "expected DisplayInSummaryLists for valid instrument %s to be %v, but it was %v", expected.ID, expected.DisplayInSummaryLists, actual.DisplayInSummaryLists)
	assert.Equal(t, expected.IncludeInGeneratedInstructions, actual.IncludeInGeneratedInstructions, "expected IncludeInGeneratedInstructions for valid instrument %s to be %v, but it was %v", expected.ID, expected.IncludeInGeneratedInstructions, actual.IncludeInGeneratedInstructions)
	assert.Equal(t, expected.Slug, actual.Slug, "expected UsableForStorage for valid instrument %s to be %v, but it was %v", expected.ID, expected.UsableForStorage, actual.UsableForStorage)
	assert.NotZero(t, actual.CreatedAt)
}

func createValidInstrumentForTest(t *testing.T, ctx context.Context, adminClient *apiclient.Client) *types.ValidInstrument {
	t.Helper()

	exampleValidInstrument := fakes.BuildFakeValidInstrument()
	exampleValidInstrumentInput := converters.ConvertValidInstrumentToValidInstrumentCreationRequestInput(exampleValidInstrument)
	createdValidInstrument, err := adminClient.CreateValidInstrument(ctx, exampleValidInstrumentInput)
	require.NoError(t, err)
	checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

	createdValidInstrument, err = adminClient.GetValidInstrument(ctx, createdValidInstrument.ID)
	requireNotNilAndNoProblems(t, createdValidInstrument, err)
	checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

	return createdValidInstrument
}

func (s *TestSuite) TestValidInstruments_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdValidInstrument := createValidInstrumentForTest(t, ctx, testClients.admin)

			newValidInstrument := fakes.BuildFakeValidInstrument()
			createdValidInstrument.Update(converters.ConvertValidInstrumentToValidInstrumentUpdateRequestInput(newValidInstrument))
			assert.NoError(t, testClients.admin.UpdateValidInstrument(ctx, createdValidInstrument))

			actual, err := testClients.admin.GetValidInstrument(ctx, createdValidInstrument.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid instrument equality
			checkValidInstrumentEquality(t, newValidInstrument, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			assert.NoError(t, testClients.admin.ArchiveValidInstrument(ctx, createdValidInstrument.ID))
		}
	})
}

func (s *TestSuite) TestValidInstruments_GetRandom() {
	s.runForEachClient("should be able to get a random valid instrument", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			exampleValidInstrument := fakes.BuildFakeValidInstrument()
			exampleValidInstrumentInput := converters.ConvertValidInstrumentToValidInstrumentCreationRequestInput(exampleValidInstrument)
			createdValidInstrument, err := testClients.admin.CreateValidInstrument(ctx, exampleValidInstrumentInput)
			require.NoError(t, err)
			checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

			actual, err := testClients.admin.GetRandomValidInstrument(ctx)
			requireNotNilAndNoProblems(t, actual, err)

			assert.NoError(t, testClients.admin.ArchiveValidInstrument(ctx, createdValidInstrument.ID))
		}
	})
}

func (s *TestSuite) TestValidInstruments_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			var expected []*types.ValidInstrument
			for i := 0; i < 5; i++ {
				exampleValidInstrument := fakes.BuildFakeValidInstrument()
				exampleValidInstrumentInput := converters.ConvertValidInstrumentToValidInstrumentCreationRequestInput(exampleValidInstrument)
				createdValidInstrument, createdValidInstrumentErr := testClients.admin.CreateValidInstrument(ctx, exampleValidInstrumentInput)
				require.NoError(t, createdValidInstrumentErr)

				checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

				expected = append(expected, createdValidInstrument)
			}

			// assert valid instrument list equality
			actual, err := testClients.admin.GetValidInstruments(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			for _, createdValidInstrument := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidInstrument(ctx, createdValidInstrument.ID))
			}
		}
	})
}

func (s *TestSuite) TestValidInstruments_Searching() {
	s.runForEachClient("should be able to be search for valid instruments", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			var expected []*types.ValidInstrument
			exampleValidInstrument := fakes.BuildFakeValidInstrument()
			exampleValidInstrument.Name = fmt.Sprintf("example_%s", testClients.authType)
			searchQuery := exampleValidInstrument.Name
			for i := 0; i < 5; i++ {
				exampleValidInstrument.Name = fmt.Sprintf("%s %d", searchQuery, i)
				exampleValidInstrumentInput := converters.ConvertValidInstrumentToValidInstrumentCreationRequestInput(exampleValidInstrument)
				createdValidInstrument, createdValidInstrumentErr := testClients.admin.CreateValidInstrument(ctx, exampleValidInstrumentInput)
				require.NoError(t, createdValidInstrumentErr)
				checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

				expected = append(expected, createdValidInstrument)
			}

			exampleLimit := uint8(20)

			// assert valid instrument list equality
			actual, err := testClients.admin.SearchValidInstruments(ctx, searchQuery, exampleLimit)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual),
				"expected %d to be <= %d",
				len(expected),
				len(actual),
			)

			for _, createdValidInstrument := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidInstrument(ctx, createdValidInstrument.ID))
			}
		}
	})
}
