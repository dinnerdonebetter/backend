package websockets

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/encoding"
	mockencoding "github.com/prixfixeco/api_server/internal/encoding/mock"
	mockconsumers "github.com/prixfixeco/api_server/internal/messagequeue/consumers/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	authservice "github.com/prixfixeco/api_server/internal/services/authentication"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func buildTestService() *service {
	return &service{
		cookieName:     "testing",
		logger:         logging.NewNoopLogger(),
		encoderDecoder: mockencoding.NewMockEncoderDecoder(),
		tracer:         tracing.NewTracer("test"),
		connections:    map[string][]websocketConnection{},
	}
}

func TestProvideService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		authCfg := &authservice.Config{}
		logger := logging.NewNoopLogger()
		encoder := encoding.ProvideServerEncoderDecoder(logger, encoding.ContentTypeJSON)

		consumer := &mockconsumers.Consumer{}
		consumer.On("Consume", chan bool(nil), chan error(nil))

		consumerProvider := &mockconsumers.ConsumerProvider{}
		consumerProvider.On(
			"ProviderConsumer",
			testutils.ContextMatcher,
			dataChangesTopicName,
			mock.Anything,
		).Return(consumer, nil)

		actual, err := ProvideService(
			ctx,
			authCfg,
			logger,
			encoder,
			consumerProvider,
		)

		require.NoError(t, err)
		require.NotNil(t, actual)

		mock.AssertExpectationsForObjects(t, consumerProvider)
	})

	T.Run("with consumer provider error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		authCfg := &authservice.Config{}
		logger := logging.NewNoopLogger()
		encoder := encoding.ProvideServerEncoderDecoder(logger, encoding.ContentTypeJSON)

		consumerProvider := &mockconsumers.ConsumerProvider{}
		consumerProvider.On(
			"ProviderConsumer",
			testutils.ContextMatcher,
			dataChangesTopicName,
			mock.Anything,
		).Return(&mockconsumers.Consumer{}, errors.New("blah"))

		actual, err := ProvideService(
			ctx,
			authCfg,
			logger,
			encoder,
			consumerProvider,
		)

		require.Error(t, err)
		require.Nil(t, actual)
	})
}

func Test_buildWebsocketErrorFunc(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		encoder := encoding.ProvideServerEncoderDecoder(nil, encoding.ContentTypeJSON)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		buildWebsocketErrorFunc(encoder)(res, req, 200, errors.New("blah"))
	})
}
