package validpreparations

import (
	"context"
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
	"github.com/prixfixeco/api_server/internal/search"
	mocksearch "github.com/prixfixeco/api_server/internal/search/mock"
	mocktypes "github.com/prixfixeco/api_server/pkg/types/mock"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func buildTestService() *service {
	return &service{
		logger:                      logging.NewNoopLogger(),
		validPreparationDataManager: &mocktypes.ValidPreparationDataManager{},
		validPreparationIDFetcher:   func(req *http.Request) string { return "" },
		encoderDecoder:              mockencoding.NewMockEncoderDecoder(),
		search:                      &mocksearch.IndexManager{},
		tracer:                      tracing.NewTracerForTest("test"),
	}
}

func TestProvideValidPreparationsService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()

		rpm := mockrouting.NewRouteParamManager()
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			ValidPreparationIDURIParamKey,
		).Return(func(*http.Request) string { return "" })

		cfg := Config{
			SearchIndexPath:      "example/path",
			PreWritesTopicName:   "pre-writes",
			PreUpdatesTopicName:  "pre-updates",
			PreArchivesTopicName: "pre-archives",
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProviderPublisher", cfg.PreWritesTopicName).Return(&mockpublishers.Publisher{}, nil)
		pp.On("ProviderPublisher", cfg.PreUpdatesTopicName).Return(&mockpublishers.Publisher{}, nil)
		pp.On("ProviderPublisher", cfg.PreArchivesTopicName).Return(&mockpublishers.Publisher{}, nil)

		mockIndexManager := &mocksearch.IndexManager{}
		mockIndexManagerProvider := &mocksearch.IndexManagerProvider{}
		mockIndexManagerProvider.On(
			"ProvideIndexManager",
			testutils.ContextMatcher,
			logger,
			search.IndexName("valid_preparations"),
			[]string{"name", "description"},
		).Return(mockIndexManager, nil)

		s, err := ProvideService(
			ctx,
			logger,
			&cfg,
			&mocktypes.ValidPreparationDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			mockIndexManagerProvider,
			rpm,
			pp,
			trace.NewNoopTracerProvider(),
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, rpm, pp)
	})

	T.Run("with error providing pre-writes producer", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()

		cfg := Config{
			SearchIndexPath:      "example/path",
			PreWritesTopicName:   "pre-writes",
			PreUpdatesTopicName:  "pre-updates",
			PreArchivesTopicName: "pre-archives",
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProviderPublisher", cfg.PreWritesTopicName).Return((*mockpublishers.Publisher)(nil), errors.New("blah"))

		mockIndexManager := &mocksearch.IndexManager{}
		mockIndexManagerProvider := &mocksearch.IndexManagerProvider{}
		mockIndexManagerProvider.On(
			"ProvideIndexManager",
			testutils.ContextMatcher,
			logger,
			search.IndexName("valid_preparations"),
			[]string{"name", "description"},
		).Return(mockIndexManager, nil)

		s, err := ProvideService(
			ctx,
			logger,
			&cfg,
			&mocktypes.ValidPreparationDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			mockIndexManagerProvider,
			nil,
			pp,
			trace.NewNoopTracerProvider(),
		)

		assert.Nil(t, s)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, pp)
	})

	T.Run("with error providing pre-updates producer", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()

		cfg := Config{
			SearchIndexPath:      "example/path",
			PreWritesTopicName:   "pre-writes",
			PreUpdatesTopicName:  "pre-updates",
			PreArchivesTopicName: "pre-archives",
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProviderPublisher", cfg.PreWritesTopicName).Return(&mockpublishers.Publisher{}, nil)
		pp.On("ProviderPublisher", cfg.PreUpdatesTopicName).Return((*mockpublishers.Publisher)(nil), errors.New("blah"))

		mockIndexManager := &mocksearch.IndexManager{}
		mockIndexManagerProvider := &mocksearch.IndexManagerProvider{}
		mockIndexManagerProvider.On(
			"ProvideIndexManager",
			testutils.ContextMatcher,
			logger,
			search.IndexName("valid_preparations"),
			[]string{"name", "description"},
		).Return(mockIndexManager, nil)

		s, err := ProvideService(
			ctx,
			logger,
			&cfg,
			&mocktypes.ValidPreparationDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			mockIndexManagerProvider,
			nil,
			pp,
			trace.NewNoopTracerProvider(),
		)

		assert.Nil(t, s)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, pp)
	})

	T.Run("with error providing pre-archives producer", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()

		cfg := Config{
			SearchIndexPath:      "example/path",
			PreWritesTopicName:   "pre-writes",
			PreUpdatesTopicName:  "pre-updates",
			PreArchivesTopicName: "pre-archives",
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProviderPublisher", cfg.PreWritesTopicName).Return(&mockpublishers.Publisher{}, nil)
		pp.On("ProviderPublisher", cfg.PreUpdatesTopicName).Return(&mockpublishers.Publisher{}, nil)
		pp.On("ProviderPublisher", cfg.PreArchivesTopicName).Return((*mockpublishers.Publisher)(nil), errors.New("blah"))

		mockIndexManager := &mocksearch.IndexManager{}
		mockIndexManagerProvider := &mocksearch.IndexManagerProvider{}
		mockIndexManagerProvider.On(
			"ProvideIndexManager",
			testutils.ContextMatcher,
			logger,
			search.IndexName("valid_preparations"),
			[]string{"name", "description"},
		).Return(mockIndexManager, nil)

		s, err := ProvideService(
			ctx,
			logger,
			&cfg,
			&mocktypes.ValidPreparationDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			mockIndexManagerProvider,
			nil,
			pp,
			trace.NewNoopTracerProvider(),
		)

		assert.Nil(t, s)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, pp)
	})

	T.Run("with error providing index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()

		cfg := Config{
			SearchIndexPath:      "example/path",
			PreWritesTopicName:   "pre-writes",
			PreUpdatesTopicName:  "pre-updates",
			PreArchivesTopicName: "pre-archives",
		}

		mockIndexManagerProvider := &mocksearch.IndexManagerProvider{}
		mockIndexManagerProvider.On(
			"ProvideIndexManager",
			testutils.ContextMatcher,
			logger,
			search.IndexName("valid_preparations"),
			[]string{"name", "description"},
		).Return(&mocksearch.IndexManager{}, errors.New("blah"))

		s, err := ProvideService(
			ctx,
			logger,
			&cfg,
			&mocktypes.ValidPreparationDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			mockIndexManagerProvider,
			mockrouting.NewRouteParamManager(),
			nil,
			trace.NewNoopTracerProvider(),
		)

		assert.Nil(t, s)
		assert.Error(t, err)
	})
}
