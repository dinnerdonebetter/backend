package recipes

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/encoding/mock"
	"github.com/dinnerdonebetter/backend/internal/features/recipeanalysis"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/objectstorage"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	mockrouting "github.com/dinnerdonebetter/backend/internal/routing/mock"
	searchcfg "github.com/dinnerdonebetter/backend/internal/search/config"
	"github.com/dinnerdonebetter/backend/internal/uploads"
	"github.com/dinnerdonebetter/backend/internal/uploads/images"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildTestService() *service {
	return &service{
		logger:            logging.NewNoopLogger(),
		recipeDataManager: &mocktypes.RecipeDataManagerMock{},
		recipeIDFetcher:   func(req *http.Request) string { return "" },
		encoderDecoder:    encoding.ProvideServerEncoderDecoder(nil, nil, encoding.ContentTypeJSON),
		tracer:            tracing.NewTracerForTest("test"),
		cfg: &Config{
			UseSearchService: false,
		},
	}
}

func TestProvideRecipesService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		rpm := mockrouting.NewRouteParamManager()
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			RecipeIDURIParamKey,
		).Return(func(*http.Request) string { return "" })

		cfg := &Config{
			Uploads: uploads.Config{
				Storage: objectstorage.Config{
					BucketName: t.Name(),
					Provider:   objectstorage.MemoryProvider,
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

		s, err := ProvideService(
			context.Background(),
			logging.NewNoopLogger(),
			cfg,
			&searchcfg.Config{},
			&mocktypes.RecipeDataManagerMock{},
			&mocktypes.RecipeMediaDataManagerMock{},
			&recipeanalysis.MockRecipeAnalyzer{},
			mockencoding.NewMockEncoderDecoder(),
			rpm,
			pp,
			&images.MockImageUploadProcessor{},
			tracing.NewNoopTracerProvider(),
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, rpm, pp)
	})

	T.Run("with error providing data changes producer", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			Uploads: uploads.Config{
				Storage: objectstorage.Config{
					BucketName: t.Name(),
					Provider:   objectstorage.MemoryProvider,
				},
				Debug: false,
			},
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProvidePublisher", cfg.DataChangesTopicName).Return((*mockpublishers.Publisher)(nil), errors.New("blah"))

		s, err := ProvideService(
			context.Background(),
			logging.NewNoopLogger(),
			cfg,
			&searchcfg.Config{},
			&mocktypes.RecipeDataManagerMock{},
			&mocktypes.RecipeMediaDataManagerMock{},
			&recipeanalysis.MockRecipeAnalyzer{},
			mockencoding.NewMockEncoderDecoder(),
			nil,
			pp,
			&images.MockImageUploadProcessor{},
			tracing.NewNoopTracerProvider(),
		)

		assert.Nil(t, s)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, pp)
	})
}
