package webhooks

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/trace"

	mockencoding "github.com/prixfixeco/api_server/internal/encoding/mock"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	mockrouting "github.com/prixfixeco/api_server/internal/routing/mock"
	mocktypes "github.com/prixfixeco/api_server/pkg/types/mock"
)

func buildTestService() *service {
	return &service{
		logger:             logging.NewNoopLogger(),
		webhookDataManager: &mocktypes.WebhookDataManager{},
		webhookIDFetcher:   func(req *http.Request) string { return "" },
		encoderDecoder:     mockencoding.NewMockEncoderDecoder(),
		tracer:             tracing.NewTracerForTest("test"),
	}
}

func TestProvideWebhooksService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		rpm := mockrouting.NewRouteParamManager()
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			WebhookIDURIParamKey,
		).Return(func(*http.Request) string { return "" })

		cfg := &Config{
			PreWritesTopicName:   "pre-writes",
			PreArchivesTopicName: "pre-archives",
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProviderPublisher", cfg.PreWritesTopicName).Return(&mockpublishers.Publisher{}, nil)
		pp.On("ProviderPublisher", cfg.PreArchivesTopicName).Return(&mockpublishers.Publisher{}, nil)

		actual, err := ProvideWebhooksService(
			logging.NewNoopLogger(),
			cfg,
			&mocktypes.WebhookDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			rpm,
			pp,
			trace.NewNoopTracerProvider(),
		)

		assert.NotNil(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, rpm, pp)
	})

	T.Run("with error providing pre-writes publisher", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			PreWritesTopicName:   "pre-writes",
			PreArchivesTopicName: "pre-archives",
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProviderPublisher", cfg.PreWritesTopicName).Return((*mockpublishers.Publisher)(nil), errors.New("blah"))

		actual, err := ProvideWebhooksService(
			logging.NewNoopLogger(),
			cfg,
			&mocktypes.WebhookDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			nil,
			pp,
			trace.NewNoopTracerProvider(),
		)

		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, pp)
	})

	T.Run("with error providing pre-archives publisher", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			PreWritesTopicName:   "pre-writes",
			PreArchivesTopicName: "pre-archives",
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProviderPublisher", cfg.PreWritesTopicName).Return(&mockpublishers.Publisher{}, nil)
		pp.On("ProviderPublisher", cfg.PreArchivesTopicName).Return((*mockpublishers.Publisher)(nil), errors.New("blah"))

		actual, err := ProvideWebhooksService(
			logging.NewNoopLogger(),
			cfg,
			&mocktypes.WebhookDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			nil,
			pp,
			trace.NewNoopTracerProvider(),
		)

		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, pp)
	})
}
