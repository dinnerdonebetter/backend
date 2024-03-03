package recipes

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/features/recipeanalysis"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/objectstorage"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing"
	"github.com/dinnerdonebetter/backend/internal/search"
	searchcfg "github.com/dinnerdonebetter/backend/internal/search/config"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	"github.com/dinnerdonebetter/backend/internal/uploads"
	"github.com/dinnerdonebetter/backend/internal/uploads/images"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "recipes_service"
)

var _ types.RecipeDataService = (*service)(nil)

type (
	// service handles recipes.
	service struct {
		logger                    logging.Logger
		tracer                    tracing.Tracer
		recipeDataManager         types.RecipeDataManager
		recipeMediaDataManager    types.RecipeMediaDataManager
		recipeAnalyzer            recipeanalysis.RecipeAnalyzer
		imageUploadProcessor      images.MediaUploadProcessor
		encoderDecoder            encoding.ServerEncoderDecoder
		dataChangesPublisher      messagequeue.Publisher
		searchIndex               search.IndexSearcher[types.RecipeSearchSubset]
		uploadManager             uploads.UploadManager
		timeFunc                  func() time.Time
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		recipeIDFetcher           func(*http.Request) string
		cfg                       *Config
	}
)

func defaultTimeFunc() time.Time {
	return time.Now()
}

var errInvalidConfig = errors.New("config cannot be nil")

// ProvideService builds a new RecipesService.
func ProvideService(
	ctx context.Context,
	logger logging.Logger,
	cfg *Config,
	searchConfig *searchcfg.Config,
	recipeDataManager types.RecipeDataManager,
	recipeMediaDataManager types.RecipeMediaDataManager,
	recipeGrapher recipeanalysis.RecipeAnalyzer,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	imageUploadProcessor images.MediaUploadProcessor,
	tracerProvider tracing.TracerProvider,
) (types.RecipeDataService, error) {
	if cfg == nil {
		return nil, errInvalidConfig
	}

	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up recipe service data changes publisher: %w", err)
	}

	uploader, err := objectstorage.NewUploadManager(ctx, logger, tracerProvider, &cfg.Uploads.Storage, routeParamManager)
	if err != nil {
		return nil, fmt.Errorf("initializing recipe service upload manager: %w", err)
	}

	searchIndex, err := searchcfg.ProvideIndex[types.RecipeSearchSubset](ctx, logger, tracerProvider, searchConfig, search.IndexTypeRecipes)
	if err != nil {
		return nil, observability.PrepareError(err, nil, "initializing recipe index manager")
	}

	svc := &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		recipeIDFetcher:           routeParamManager.BuildRouteParamStringIDFetcher(RecipeIDURIParamKey),
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		recipeDataManager:         recipeDataManager,
		cfg:                       cfg,
		recipeMediaDataManager:    recipeMediaDataManager,
		dataChangesPublisher:      dataChangesPublisher,
		encoderDecoder:            encoder,
		timeFunc:                  defaultTimeFunc,
		recipeAnalyzer:            recipeGrapher,
		uploadManager:             uploader,
		imageUploadProcessor:      imageUploadProcessor,
		tracer:                    tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		searchIndex:               searchIndex,
	}

	return svc, nil
}
