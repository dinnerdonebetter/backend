package recipesteps

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/objectstorage"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	recipesservice "github.com/dinnerdonebetter/backend/internal/services/recipes"
	"github.com/dinnerdonebetter/backend/internal/uploads"
	"github.com/dinnerdonebetter/backend/internal/uploads/images"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "recipe_steps_service"
)

var _ types.RecipeStepDataService = (*service)(nil)

type (
	// service handles recipe steps.
	service struct {
		logger                    logging.Logger
		dataChangesPublisher      messagequeue.Publisher
		recipeStepDataManager     types.RecipeStepDataManager
		recipeMediaDataManager    types.RecipeMediaDataManager
		uploadManager             uploads.UploadManager
		tracer                    tracing.Tracer
		encoderDecoder            encoding.ServerEncoderDecoder
		imageUploadProcessor      images.MediaUploadProcessor
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		recipeIDFetcher           func(*http.Request) string
		recipeStepIDFetcher       func(*http.Request) string
		cfg                       Config
	}
)

// ProvideService builds a new RecipeStepsService.
func ProvideService(
	ctx context.Context,
	logger logging.Logger,
	cfg *Config,
	recipeStepDataManager types.RecipeStepDataManager,
	recipeMediaDataManager types.RecipeMediaDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
	imageUploadProcessor images.MediaUploadProcessor,
) (types.RecipeStepDataService, error) {
	if cfg == nil {
		return nil, errors.New("nil config provided to recipe steps service")
	}

	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up recipe steps service data changes publisher: %w", err)
	}

	uploader, err := objectstorage.NewUploadManager(ctx, logger, tracerProvider, &cfg.Uploads.Storage, routeParamManager)
	if err != nil {
		return nil, fmt.Errorf("initializing recipe steps service upload manager: %w", err)
	}

	svc := &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		recipeIDFetcher:           routeParamManager.BuildRouteParamStringIDFetcher(recipesservice.RecipeIDURIParamKey),
		recipeStepIDFetcher:       routeParamManager.BuildRouteParamStringIDFetcher(RecipeStepIDURIParamKey),
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		recipeStepDataManager:     recipeStepDataManager,
		recipeMediaDataManager:    recipeMediaDataManager,
		dataChangesPublisher:      dataChangesPublisher,
		encoderDecoder:            encoder,
		cfg:                       *cfg,
		imageUploadProcessor:      imageUploadProcessor,
		uploadManager:             uploader,
		tracer:                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc, nil
}
