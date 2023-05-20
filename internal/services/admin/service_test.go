package admin

import (
	"net/http"
	"testing"

	mockauthn "github.com/dinnerdonebetter/backend/internal/authentication/mock"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	mockrouting "github.com/dinnerdonebetter/backend/internal/routing/mock"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/alexedwards/scs/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func buildTestService(t *testing.T) *service {
	t.Helper()

	logger := logging.NewNoopLogger()

	rpm := mockrouting.NewRouteParamManager()
	rpm.On(
		"BuildRouteParamStringIDFetcher",
		UserIDURIParamKey,
	).Return(func(*http.Request) string { return "" })

	s := ProvideService(
		logger,
		&authservice.Config{Cookies: authservice.CookieConfig{BlockKey: "BLAHBLAHBLAHPRETENDTHISISSECRET!"}},
		&mockauthn.Authenticator{},
		&mocktypes.AdminUserDataManager{},
		scs.New(),
		encoding.ProvideServerEncoderDecoder(logger, tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON),
		rpm,
		tracing.NewNoopTracerProvider(),
	)

	mock.AssertExpectationsForObjects(t, rpm)

	srv, ok := s.(*service)
	require.True(t, ok)

	return srv
}

func TestProvideAdminService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()

		rpm := mockrouting.NewRouteParamManager()
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			UserIDURIParamKey,
		).Return(func(*http.Request) string { return "" })

		s := ProvideService(
			logger,
			&authservice.Config{Cookies: authservice.CookieConfig{BlockKey: "BLAHBLAHBLAHPRETENDTHISISSECRET!"}},
			&mockauthn.Authenticator{},
			&mocktypes.AdminUserDataManager{},
			scs.New(),
			encoding.ProvideServerEncoderDecoder(logger, tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON),
			rpm,
			tracing.NewNoopTracerProvider(),
		)

		assert.NotNil(t, s)

		mock.AssertExpectationsForObjects(t, rpm)
	})
}
