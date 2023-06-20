package integration

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func checkAPIClientEquality(t *testing.T, expected, actual *types.APIClient) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected LabelName for API client %s to be %q, but it was %q ", actual.ID, expected.Name, actual.Name)
	assert.NotEmpty(t, actual.ClientID, "expected ClientID for API client %s to not be empty, but it was", actual.ID)
	assert.Empty(t, actual.ClientSecret, "expected ClientSecret for API client %s to not be empty, but it was", actual.ID)
	assert.NotZero(t, actual.BelongsToUser, "expected CreatedByUser for API client %s to not be zero, but it was", actual.ID)
	assert.NotZero(t, actual.CreatedAt)
}

func (s *TestSuite) TestAPIClients_Creating() {
	s.runForEachClient("should be possible to create API clients", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create API client.
			exampleAPIClient := fakes.BuildFakeAPIClient()
			exampleAPIClientInput := fakes.BuildFakeAPIClientCreationInputFromClient(exampleAPIClient)
			exampleAPIClientInput.UserLoginInput = types.UserLoginInput{
				Username:  s.user.Username,
				Password:  s.user.HashedPassword,
				TOTPToken: generateTOTPTokenForUser(t, s.user),
			}

			createdAPIClient, err := testClients.user.CreateAPIClient(ctx, s.cookie, exampleAPIClientInput)
			requireNotNilAndNoProblems(t, createdAPIClient, err)

			// Assert API client equality.
			assert.NotEmpty(t, createdAPIClient.ClientID, "expected ClientID for API client %s to not be empty, but it was", createdAPIClient.ID)
			assert.NotEmpty(t, createdAPIClient.ClientSecret, "expected ClientSecret for API client %s to not be empty, but it was", createdAPIClient.ID)

			// Clean up.
			assert.NoError(t, testClients.user.ArchiveAPIClient(ctx, createdAPIClient.ID))
		}
	})
}

func (s *TestSuite) TestAPIClients_Listing() {
	s.runForEachClient("should be possible to read API clients in a list", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			const clientsToMake = 1

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create API clients.
			var expected []string
			for i := 0; i < clientsToMake; i++ {
				// Create API client.
				exampleAPIClient := fakes.BuildFakeAPIClient()
				exampleAPIClientInput := fakes.BuildFakeAPIClientCreationInputFromClient(exampleAPIClient)
				exampleAPIClientInput.UserLoginInput = types.UserLoginInput{
					Username:  s.user.Username,
					Password:  s.user.HashedPassword,
					TOTPToken: generateTOTPTokenForUser(t, s.user),
				}
				createdAPIClient, apiClientCreationErr := testClients.user.CreateAPIClient(ctx, s.cookie, exampleAPIClientInput)
				requireNotNilAndNoProblems(t, createdAPIClient, apiClientCreationErr)

				expected = append(expected, createdAPIClient.ID)
			}

			// Assert API client list equality.
			actual, err := testClients.user.GetAPIClients(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			// Clean up.
			for _, createdAPIClientID := range expected {
				assert.NoError(t, testClients.user.ArchiveAPIClient(ctx, createdAPIClientID))
			}
		}
	})
}

func (s *TestSuite) TestAPIClients_Reading_Returns404ForNonexistentAPIClient() {
	s.runForEachClient("should not be possible to read non-existent API clients", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Attempt to fetch nonexistent API client.
			_, err := testClients.user.GetAPIClient(ctx, nonexistentID)
			assert.Error(t, err)
		}
	})
}

func (s *TestSuite) TestAPIClients_Reading() {
	s.runForEachClient("should be possible to read API clients", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create API client.
			exampleAPIClient := fakes.BuildFakeAPIClient()
			exampleAPIClientInput := fakes.BuildFakeAPIClientCreationInputFromClient(exampleAPIClient)
			exampleAPIClientInput.UserLoginInput = types.UserLoginInput{
				Username:  s.user.Username,
				Password:  s.user.HashedPassword,
				TOTPToken: generateTOTPTokenForUser(t, s.user),
			}

			createdAPIClient, err := testClients.user.CreateAPIClient(ctx, s.cookie, exampleAPIClientInput)
			requireNotNilAndNoProblems(t, createdAPIClient, err)

			// Fetch API client.
			actual, err := testClients.user.GetAPIClient(ctx, createdAPIClient.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// Assert API client equality.
			checkAPIClientEquality(t, exampleAPIClient, actual)

			// Clean up API client.
			assert.NoError(t, testClients.user.ArchiveAPIClient(ctx, createdAPIClient.ID))
		}
	})
}

func (s *TestSuite) TestAPIClients_Archiving() {
	s.runForEachClient("should be possible to archive API clients", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create API client.
			exampleAPIClient := fakes.BuildFakeAPIClient()
			exampleAPIClientInput := fakes.BuildFakeAPIClientCreationInputFromClient(exampleAPIClient)
			exampleAPIClientInput.UserLoginInput = types.UserLoginInput{
				Username:  s.user.Username,
				Password:  s.user.HashedPassword,
				TOTPToken: generateTOTPTokenForUser(t, s.user),
			}

			createdAPIClient, err := testClients.user.CreateAPIClient(ctx, s.cookie, exampleAPIClientInput)
			requireNotNilAndNoProblems(t, createdAPIClient, err)

			// Clean up API client.
			assert.NoError(t, testClients.user.ArchiveAPIClient(ctx, createdAPIClient.ID))
		}
	})

	s.runForEachClient("should not be possible to archive non-existent API clients", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(context.Background(), t.Name())
			defer span.End()

			assert.Error(t, testClients.user.ArchiveAPIClient(ctx, nonexistentID))
		}
	})
}
