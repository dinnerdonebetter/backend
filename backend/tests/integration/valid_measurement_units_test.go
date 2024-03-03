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

func checkValidMeasurementUnitEquality(t *testing.T, expected, actual *types.ValidMeasurementUnit) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for valid measurement unit %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for valid measurement unit %s to be %v, but it was %v", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.Volumetric, actual.Volumetric, "expected Volumetric for valid measurement unit %s to be %v, but it was %v", expected.ID, expected.Volumetric, actual.Volumetric)
	assert.Equal(t, expected.IconPath, actual.IconPath, "expected IconPath for valid measurement unit %s to be %v, but it was %v", expected.ID, expected.IconPath, actual.IconPath)
	assert.Equal(t, expected.Universal, actual.Universal, "expected Universal for valid measurement unit %s to be %v, but it was %v", expected.ID, expected.Universal, actual.Universal)
	assert.Equal(t, expected.Metric, actual.Metric, "expected Metric for valid measurement unit %s to be %v, but it was %v", expected.ID, expected.Metric, actual.Metric)
	assert.Equal(t, expected.Imperial, actual.Imperial, "expected Imperial for valid measurement unit %s to be %v, but it was %v", expected.ID, expected.Imperial, actual.Imperial)
	assert.Equal(t, expected.Slug, actual.Slug, "expected Slug for valid measurement unit %s to be %v, but it was %v", expected.ID, expected.Slug, actual.Slug)
	assert.Equal(t, expected.PluralName, actual.PluralName, "expected PluralName for valid measurement unit %s to be %v, but it was %v", expected.ID, expected.PluralName, actual.PluralName)
	assert.NotZero(t, actual.CreatedAt)
}

func createValidMeasurementUnitForTest(t *testing.T, ctx context.Context, adminClient *apiclient.Client) *types.ValidMeasurementUnit {
	exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
	exampleValidMeasurementUnitInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(exampleValidMeasurementUnit)
	createdValidMeasurementUnit, err := adminClient.CreateValidMeasurementUnit(ctx, exampleValidMeasurementUnitInput)
	require.NoError(t, err)
	checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

	createdValidMeasurementUnit, err = adminClient.GetValidMeasurementUnit(ctx, createdValidMeasurementUnit.ID)
	requireNotNilAndNoProblems(t, createdValidMeasurementUnit, err)
	checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

	return createdValidMeasurementUnit
}

func (s *TestSuite) TestValidMeasurementUnits_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdValidMeasurementUnit := createValidMeasurementUnitForTest(t, ctx, testClients.admin)

			newValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
			createdValidMeasurementUnit.Update(converters.ConvertValidMeasurementUnitToValidMeasurementUnitUpdateRequestInput(newValidMeasurementUnit))
			assert.NoError(t, testClients.admin.UpdateValidMeasurementUnit(ctx, createdValidMeasurementUnit))

			actual, err := testClients.admin.GetValidMeasurementUnit(ctx, createdValidMeasurementUnit.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid measurement unit equality
			checkValidMeasurementUnitEquality(t, newValidMeasurementUnit, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			assert.NoError(t, testClients.admin.ArchiveValidMeasurementUnit(ctx, createdValidMeasurementUnit.ID))
		}
	})
}

func (s *TestSuite) TestValidMeasurementUnits_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			var expected []*types.ValidMeasurementUnit
			for i := 0; i < 5; i++ {
				exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
				exampleValidMeasurementUnitInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(exampleValidMeasurementUnit)
				createdValidMeasurementUnit, createdValidMeasurementUnitErr := testClients.admin.CreateValidMeasurementUnit(ctx, exampleValidMeasurementUnitInput)
				require.NoError(t, createdValidMeasurementUnitErr)

				checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

				expected = append(expected, createdValidMeasurementUnit)
			}

			// assert valid measurement unit list equality
			actual, err := testClients.admin.GetValidMeasurementUnits(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			for _, createdValidMeasurementUnit := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidMeasurementUnit(ctx, createdValidMeasurementUnit.ID))
			}
		}
	})
}

func (s *TestSuite) TestValidMeasurementUnits_Searching() {
	s.runForEachClient("should be able to be search for valid measurement units", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			var expected []*types.ValidMeasurementUnit
			exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
			exampleValidMeasurementUnit.Name = fmt.Sprintf("example_%s", testClients.authType)
			searchQuery := exampleValidMeasurementUnit.Name
			for i := 0; i < 5; i++ {
				exampleValidMeasurementUnit.Name = fmt.Sprintf("%s %d", searchQuery, i)
				exampleValidMeasurementUnitInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(exampleValidMeasurementUnit)
				createdValidMeasurementUnit, createdValidMeasurementUnitErr := testClients.admin.CreateValidMeasurementUnit(ctx, exampleValidMeasurementUnitInput)
				require.NoError(t, createdValidMeasurementUnitErr)
				checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

				expected = append(expected, createdValidMeasurementUnit)
			}

			exampleLimit := uint8(20)

			// assert valid measurement unit list equality
			actual, err := testClients.admin.SearchValidMeasurementUnits(ctx, searchQuery, exampleLimit)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual),
				"expected %d to be <= %d",
				len(expected),
				len(actual),
			)

			for _, createdValidMeasurementUnit := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidMeasurementUnit(ctx, createdValidMeasurementUnit.ID))
			}
		}
	})
}
