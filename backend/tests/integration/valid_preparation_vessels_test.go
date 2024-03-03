package integration

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkValidPreparationVesselEquality(t *testing.T, expected, actual *types.ValidPreparationVessel) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected StatusExplanation for valid preparation vessel %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.Preparation.ID, actual.Preparation.ID, "expected Preparation for valid preparation vessel %s to be %v, but it was %v", expected.ID, expected.Preparation.ID, actual.Preparation.ID)
	assert.Equal(t, expected.Vessel.ID, actual.Vessel.ID, "expected Vessel for valid preparation vessel %s to be %v, but it was %v", expected.ID, expected.Vessel.ID, actual.Vessel.ID)
	assert.NotZero(t, actual.CreatedAt)
}

func createValidPreparationVesselForTest(t *testing.T, ctx context.Context, adminClient *apiclient.Client) (*types.ValidPreparation, *types.ValidVessel, *types.ValidPreparationVessel) {
	t.Helper()

	createdValidPreparation := createValidPreparationForTest(t, ctx, nil, adminClient)
	createdValidVessel := createValidVesselForTest(t, ctx, nil, adminClient)

	exampleValidPreparationVessel := fakes.BuildFakeValidPreparationVessel()
	exampleValidPreparationVessel.Vessel = *createdValidVessel
	exampleValidPreparationVessel.Preparation = *createdValidPreparation
	exampleValidPreparationVesselInput := converters.ConvertValidPreparationVesselToValidPreparationVesselCreationRequestInput(exampleValidPreparationVessel)
	createdValidPreparationVessel, err := adminClient.CreateValidPreparationVessel(ctx, exampleValidPreparationVesselInput)
	require.NoError(t, err)

	checkValidPreparationVesselEquality(t, exampleValidPreparationVessel, createdValidPreparationVessel)

	createdValidPreparationVessel, err = adminClient.GetValidPreparationVessel(ctx, createdValidPreparationVessel.ID)
	requireNotNilAndNoProblems(t, createdValidPreparationVessel, err)

	return createdValidPreparation, createdValidVessel, createdValidPreparationVessel
}

func (s *TestSuite) TestValidPreparationVessels_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdValidPreparationVessel := createValidPreparationVesselForTest(t, ctx, testClients.admin)

			newValidPreparationVessel := fakes.BuildFakeValidPreparationVessel()
			newValidPreparationVessel.Vessel = createdValidPreparationVessel.Vessel
			newValidPreparationVessel.Preparation = createdValidPreparationVessel.Preparation
			createdValidPreparationVessel.Update(converters.ConvertValidPreparationVesselToValidPreparationVesselUpdateRequestInput(newValidPreparationVessel))
			assert.NoError(t, testClients.admin.UpdateValidPreparationVessel(ctx, createdValidPreparationVessel))

			actual, err := testClients.user.GetValidPreparationVessel(ctx, createdValidPreparationVessel.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid preparation vessel equality
			checkValidPreparationVesselEquality(t, newValidPreparationVessel, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			assert.NoError(t, testClients.admin.ArchiveValidPreparationVessel(ctx, createdValidPreparationVessel.ID))

			assert.NoError(t, testClients.admin.ArchiveValidVessel(ctx, createdValidPreparationVessel.Vessel.ID))

			assert.NoError(t, testClients.admin.ArchiveValidPreparation(ctx, createdValidPreparationVessel.Preparation.ID))
		}
	})
}

func (s *TestSuite) TestValidPreparationVessels_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			var expected []*types.ValidPreparationVessel
			for i := 0; i < 5; i++ {
				_, _, createdValidPreparationVessel := createValidPreparationVesselForTest(t, ctx, testClients.admin)
				expected = append(expected, createdValidPreparationVessel)
			}

			// assert valid preparation vessel list equality
			actual, err := testClients.user.GetValidPreparationVessels(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			for _, createdValidPreparationVessel := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidPreparationVessel(ctx, createdValidPreparationVessel.ID))
			}
		}
	})
}

func (s *TestSuite) TestValidPreparationVessels_Listing_ByValue() {
	s.runForEachClient("should be findable via either member of the bridge type", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdValidPreparation, createdValidVessel, createdValidPreparationVessel := createValidPreparationVesselForTest(t, ctx, testClients.admin)

			validPreparationVesselsForVessel, err := testClients.user.GetValidPreparationVesselsForVessel(ctx, createdValidVessel.ID, nil)
			requireNotNilAndNoProblems(t, validPreparationVesselsForVessel, err)

			require.Len(t, validPreparationVesselsForVessel.Data, 1)
			assert.Equal(t, validPreparationVesselsForVessel.Data[0].ID, createdValidPreparationVessel.ID)

			validPreparationVesselsForPreparation, err := testClients.user.GetValidPreparationVesselsForPreparation(ctx, createdValidPreparation.ID, nil)
			requireNotNilAndNoProblems(t, validPreparationVesselsForPreparation, err)

			require.Len(t, validPreparationVesselsForPreparation.Data, 1)
			assert.Equal(t, validPreparationVesselsForPreparation.Data[0].ID, createdValidPreparationVessel.ID)

			assert.NoError(t, testClients.admin.ArchiveValidPreparationVessel(ctx, createdValidPreparationVessel.ID))

			assert.NoError(t, testClients.admin.ArchiveValidVessel(ctx, createdValidVessel.ID))

			assert.NoError(t, testClients.admin.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
		}
	})
}
