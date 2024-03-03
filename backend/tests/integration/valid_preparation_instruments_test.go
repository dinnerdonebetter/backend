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

func checkValidPreparationInstrumentEquality(t *testing.T, expected, actual *types.ValidPreparationInstrument) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected StatusExplanation for valid preparation instrument %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.Preparation.ID, actual.Preparation.ID, "expected Preparation for valid preparation instrument %s to be %v, but it was %v", expected.ID, expected.Preparation.ID, actual.Preparation.ID)
	assert.Equal(t, expected.Instrument.ID, actual.Instrument.ID, "expected Vessel for valid preparation instrument %s to be %v, but it was %v", expected.ID, expected.Instrument.ID, actual.Instrument.ID)
	assert.NotZero(t, actual.CreatedAt)
}

func (s *TestSuite) TestValidPreparationInstruments_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(exampleValidPreparation)
			createdValidPreparation, err := testClients.admin.CreateValidPreparation(ctx, exampleValidPreparationInput)
			require.NoError(t, err)

			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			createdValidPreparation, err = testClients.user.GetValidPreparation(ctx, createdValidPreparation.ID)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			exampleValidInstrument := fakes.BuildFakeValidInstrument()
			exampleValidInstrumentInput := converters.ConvertValidInstrumentToValidInstrumentCreationRequestInput(exampleValidInstrument)
			createdValidInstrument, err := testClients.admin.CreateValidInstrument(ctx, exampleValidInstrumentInput)
			require.NoError(t, err)

			checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

			createdValidInstrument, err = testClients.user.GetValidInstrument(ctx, createdValidInstrument.ID)
			requireNotNilAndNoProblems(t, createdValidInstrument, err)
			checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

			exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
			exampleValidPreparationInstrument.Instrument = *createdValidInstrument
			exampleValidPreparationInstrument.Preparation = *createdValidPreparation
			exampleValidPreparationInstrumentInput := converters.ConvertValidPreparationInstrumentToValidPreparationInstrumentCreationRequestInput(exampleValidPreparationInstrument)
			createdValidPreparationInstrument, err := testClients.admin.CreateValidPreparationInstrument(ctx, exampleValidPreparationInstrumentInput)
			require.NoError(t, err)

			checkValidPreparationInstrumentEquality(t, exampleValidPreparationInstrument, createdValidPreparationInstrument)

			createdValidPreparationInstrument, err = testClients.user.GetValidPreparationInstrument(ctx, createdValidPreparationInstrument.ID)
			requireNotNilAndNoProblems(t, createdValidPreparationInstrument, err)

			checkValidPreparationInstrumentEquality(t, exampleValidPreparationInstrument, createdValidPreparationInstrument)

			newValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
			newValidPreparationInstrument.Instrument = *createdValidInstrument
			newValidPreparationInstrument.Preparation = *createdValidPreparation
			createdValidPreparationInstrument.Update(converters.ConvertValidPreparationInstrumentToValidPreparationInstrumentUpdateRequestInput(newValidPreparationInstrument))
			assert.NoError(t, testClients.admin.UpdateValidPreparationInstrument(ctx, createdValidPreparationInstrument))

			actual, err := testClients.user.GetValidPreparationInstrument(ctx, createdValidPreparationInstrument.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid preparation instrument equality
			checkValidPreparationInstrumentEquality(t, newValidPreparationInstrument, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			assert.NoError(t, testClients.admin.ArchiveValidPreparationInstrument(ctx, createdValidPreparationInstrument.ID))

			assert.NoError(t, testClients.admin.ArchiveValidInstrument(ctx, createdValidInstrument.ID))

			assert.NoError(t, testClients.admin.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
		}
	})
}

func (s *TestSuite) TestValidPreparationInstruments_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			var expected []*types.ValidPreparationInstrument
			for i := 0; i < 5; i++ {
				exampleValidPreparation := fakes.BuildFakeValidPreparation()
				exampleValidPreparationInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(exampleValidPreparation)
				createdValidPreparation, err := testClients.admin.CreateValidPreparation(ctx, exampleValidPreparationInput)
				require.NoError(t, err)

				checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

				createdValidPreparation, err = testClients.user.GetValidPreparation(ctx, createdValidPreparation.ID)
				requireNotNilAndNoProblems(t, createdValidPreparation, err)
				checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

				exampleValidInstrument := fakes.BuildFakeValidInstrument()
				exampleValidInstrumentInput := converters.ConvertValidInstrumentToValidInstrumentCreationRequestInput(exampleValidInstrument)
				createdValidInstrument, err := testClients.admin.CreateValidInstrument(ctx, exampleValidInstrumentInput)
				require.NoError(t, err)

				checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

				createdValidInstrument, err = testClients.user.GetValidInstrument(ctx, createdValidInstrument.ID)
				requireNotNilAndNoProblems(t, createdValidInstrument, err)
				checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

				exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
				exampleValidPreparationInstrument.Instrument = *createdValidInstrument
				exampleValidPreparationInstrument.Preparation = *createdValidPreparation
				exampleValidPreparationInstrumentInput := converters.ConvertValidPreparationInstrumentToValidPreparationInstrumentCreationRequestInput(exampleValidPreparationInstrument)
				createdValidPreparationInstrument, createdValidPreparationInstrumentErr := testClients.admin.CreateValidPreparationInstrument(ctx, exampleValidPreparationInstrumentInput)
				require.NoError(t, createdValidPreparationInstrumentErr)

				checkValidPreparationInstrumentEquality(t, exampleValidPreparationInstrument, createdValidPreparationInstrument)

				expected = append(expected, createdValidPreparationInstrument)
			}

			// assert valid preparation instrument list equality
			actual, err := testClients.user.GetValidPreparationInstruments(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			for _, createdValidPreparationInstrument := range expected {
				assert.NoError(t, testClients.admin.ArchiveValidPreparationInstrument(ctx, createdValidPreparationInstrument.ID))
			}
		}
	})
}

func (s *TestSuite) TestValidPreparationInstruments_Listing_ByValue() {
	s.runForEachClient("should be findable via either member of the bridge type", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := converters.ConvertValidPreparationToValidPreparationCreationRequestInput(exampleValidPreparation)
			createdValidPreparation, err := testClients.admin.CreateValidPreparation(ctx, exampleValidPreparationInput)
			require.NoError(t, err)

			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			createdValidPreparation, err = testClients.user.GetValidPreparation(ctx, createdValidPreparation.ID)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			exampleValidInstrument := fakes.BuildFakeValidInstrument()
			exampleValidInstrumentInput := converters.ConvertValidInstrumentToValidInstrumentCreationRequestInput(exampleValidInstrument)
			createdValidInstrument, err := testClients.admin.CreateValidInstrument(ctx, exampleValidInstrumentInput)
			require.NoError(t, err)

			checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

			createdValidInstrument, err = testClients.user.GetValidInstrument(ctx, createdValidInstrument.ID)
			requireNotNilAndNoProblems(t, createdValidInstrument, err)
			checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

			exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
			exampleValidPreparationInstrument.Instrument = *createdValidInstrument
			exampleValidPreparationInstrument.Preparation = *createdValidPreparation
			exampleValidPreparationInstrumentInput := converters.ConvertValidPreparationInstrumentToValidPreparationInstrumentCreationRequestInput(exampleValidPreparationInstrument)
			createdValidPreparationInstrument, err := testClients.admin.CreateValidPreparationInstrument(ctx, exampleValidPreparationInstrumentInput)
			require.NoError(t, err)

			checkValidPreparationInstrumentEquality(t, exampleValidPreparationInstrument, createdValidPreparationInstrument)

			validPreparationInstrumentsForInstrument, err := testClients.user.GetValidPreparationInstrumentsForInstrument(ctx, createdValidInstrument.ID, nil)
			requireNotNilAndNoProblems(t, validPreparationInstrumentsForInstrument, err)

			require.Len(t, validPreparationInstrumentsForInstrument.Data, 1)
			assert.Equal(t, validPreparationInstrumentsForInstrument.Data[0].ID, createdValidPreparationInstrument.ID)

			validPreparationInstrumentsForPreparation, err := testClients.user.GetValidPreparationInstrumentsForPreparation(ctx, createdValidPreparation.ID, nil)
			requireNotNilAndNoProblems(t, validPreparationInstrumentsForPreparation, err)

			require.Len(t, validPreparationInstrumentsForPreparation.Data, 1)
			assert.Equal(t, validPreparationInstrumentsForPreparation.Data[0].ID, createdValidPreparationInstrument.ID)

			assert.NoError(t, testClients.admin.ArchiveValidPreparationInstrument(ctx, createdValidPreparationInstrument.ID))

			assert.NoError(t, testClients.admin.ArchiveValidInstrument(ctx, createdValidInstrument.ID))

			assert.NoError(t, testClients.admin.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
		}
	})
}
