package users

import (
	"context"
	"net/http"
	"testing"

	mockauthn "github.com/dinnerdonebetter/backend/internal/authentication/mock"
	mockencoding "github.com/dinnerdonebetter/backend/internal/encoding/mock"
	"github.com/dinnerdonebetter/backend/internal/featureflags"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/objectstorage"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/random"
	"github.com/dinnerdonebetter/backend/internal/routing/chi"
	mockrouting "github.com/dinnerdonebetter/backend/internal/routing/mock"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	"github.com/dinnerdonebetter/backend/internal/uploads"
	"github.com/dinnerdonebetter/backend/internal/uploads/images"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func buildTestService(t *testing.T) *service {
	t.Helper()

	cfg := &Config{
		Uploads: uploads.Config{
			Storage: objectstorage.Config{
				FilesystemConfig: &objectstorage.FilesystemConfig{RootDirectory: t.Name()},
				BucketName:       t.Name(),
				Provider:         objectstorage.FilesystemProvider,
			},
			Debug: false,
		},
	}

	pp := &mockpublishers.ProducerProvider{}
	pp.On("ProvidePublisher", cfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

	s, err := ProvideUsersService(
		context.Background(),
		cfg,
		&authservice.Config{},
		logging.NewNoopLogger(),
		&mocktypes.UserDataManager{},
		&mocktypes.HouseholdDataManager{},
		&mocktypes.HouseholdInvitationDataManager{},
		&mocktypes.HouseholdUserMembershipDataManager{},
		&mockauthn.Authenticator{},
		mockencoding.NewMockEncoderDecoder(),
		&images.MockImageUploadProcessor{},
		chi.NewRouteParamManager(),
		tracing.NewNoopTracerProvider(),
		pp,
		random.NewGenerator(logging.NewNoopLogger(), tracing.NewNoopTracerProvider()),
		&mocktypes.PasswordResetTokenDataManager{},
		&featureflags.NoopFeatureFlagManager{},
	)

	require.NoError(t, err)

	return s.(*service)
}

func TestProvideUsersService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		rpm := mockrouting.NewRouteParamManager()
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			UserIDURIParamKey,
		).Return(func(*http.Request) string { return "" })

		cfg := &Config{
			Uploads: uploads.Config{
				Storage: objectstorage.Config{
					FilesystemConfig: &objectstorage.FilesystemConfig{RootDirectory: t.Name()},
					BucketName:       t.Name(),
					Provider:         objectstorage.FilesystemProvider,
				},
				Debug: false,
			},
		}

		rpm.On(
			"BuildRouteParamStringIDFetcher",
			cfg.Uploads.Storage.UploadFilenameKey,
		).Return(func(*http.Request) string { return "" })

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProvidePublisher", cfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

		s, err := ProvideUsersService(
			context.Background(),
			cfg,
			&authservice.Config{},
			logging.NewNoopLogger(),
			&mocktypes.UserDataManager{},
			&mocktypes.HouseholdDataManager{},
			&mocktypes.HouseholdInvitationDataManager{},
			&mocktypes.HouseholdUserMembershipDataManager{},
			&mockauthn.Authenticator{},
			mockencoding.NewMockEncoderDecoder(),
			&images.MockImageUploadProcessor{},
			rpm,
			tracing.NewNoopTracerProvider(),
			pp,
			random.NewGenerator(logging.NewNoopLogger(), tracing.NewNoopTracerProvider()),
			&mocktypes.PasswordResetTokenDataManager{},
			&featureflags.NoopFeatureFlagManager{},
		)

		assert.NotNil(t, s)
		require.NoError(t, err)

		mock.AssertExpectationsForObjects(t, rpm)
	})
}
