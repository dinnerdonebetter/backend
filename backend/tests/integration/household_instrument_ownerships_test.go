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

func checkHouseholdInstrumentOwnershipEquality(t *testing.T, expected, actual *types.HouseholdInstrumentOwnership) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for household instrument ownership %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.Instrument.ID, actual.Instrument.ID, "expected Vessel.ID for household instrument ownership %s to be %v, but it was %v", expected.ID, expected.Instrument.ID, actual.Instrument.ID)
	assert.Equal(t, expected.Quantity, actual.Quantity, "expected Quantity for household instrument ownership %s to be %v, but it was %v", expected.ID, expected.Quantity, actual.Quantity)
	assert.NotZero(t, actual.CreatedAt)
}

func (s *TestSuite) TestHouseholdInstrumentOwnerships_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdValidInstrument := createValidInstrumentForTest(t, ctx, testClients.admin)

			exampleHouseholdInstrumentOwnership := fakes.BuildFakeHouseholdInstrumentOwnership()
			exampleHouseholdInstrumentOwnership.Instrument = *createdValidInstrument
			exampleHouseholdInstrumentOwnershipInput := converters.ConvertHouseholdInstrumentOwnershipToHouseholdInstrumentOwnershipCreationRequestInput(exampleHouseholdInstrumentOwnership)
			createdHouseholdInstrumentOwnership, err := testClients.admin.CreateHouseholdInstrumentOwnership(ctx, exampleHouseholdInstrumentOwnershipInput)
			require.NoError(t, err)
			checkHouseholdInstrumentOwnershipEquality(t, exampleHouseholdInstrumentOwnership, createdHouseholdInstrumentOwnership)

			newHouseholdInstrumentOwnership := fakes.BuildFakeHouseholdInstrumentOwnership()
			newHouseholdInstrumentOwnership.Instrument = *createdValidInstrument
			createdHouseholdInstrumentOwnership.Update(converters.ConvertHouseholdInstrumentOwnershipToHouseholdInstrumentOwnershipUpdateRequestInput(newHouseholdInstrumentOwnership))
			assert.NoError(t, testClients.admin.UpdateHouseholdInstrumentOwnership(ctx, createdHouseholdInstrumentOwnership))

			actual, err := testClients.admin.GetHouseholdInstrumentOwnership(ctx, createdHouseholdInstrumentOwnership.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert household instrument ownership equality
			checkHouseholdInstrumentOwnershipEquality(t, newHouseholdInstrumentOwnership, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			assert.NoError(t, testClients.admin.ArchiveHouseholdInstrumentOwnership(ctx, createdHouseholdInstrumentOwnership.ID))
		}
	})
}

func (s *TestSuite) TestHouseholdInstrumentOwnerships_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			var expected []*types.HouseholdInstrumentOwnership
			for i := 0; i < 5; i++ {
				createdValidInstrument := createValidInstrumentForTest(t, ctx, testClients.admin)

				exampleHouseholdInstrumentOwnership := fakes.BuildFakeHouseholdInstrumentOwnership()
				exampleHouseholdInstrumentOwnership.Instrument = *createdValidInstrument
				exampleHouseholdInstrumentOwnershipInput := converters.ConvertHouseholdInstrumentOwnershipToHouseholdInstrumentOwnershipCreationRequestInput(exampleHouseholdInstrumentOwnership)
				createdHouseholdInstrumentOwnership, err := testClients.admin.CreateHouseholdInstrumentOwnership(ctx, exampleHouseholdInstrumentOwnershipInput)
				require.NoError(t, err)
				checkHouseholdInstrumentOwnershipEquality(t, exampleHouseholdInstrumentOwnership, createdHouseholdInstrumentOwnership)

				expected = append(expected, createdHouseholdInstrumentOwnership)
			}

			// assert household instrument ownership list equality
			actual, err := testClients.admin.GetHouseholdInstrumentOwnerships(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			for _, createdHouseholdInstrumentOwnership := range expected {
				assert.NoError(t, testClients.admin.ArchiveHouseholdInstrumentOwnership(ctx, createdHouseholdInstrumentOwnership.ID))
			}
		}
	})
}
