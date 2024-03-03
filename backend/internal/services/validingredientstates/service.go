package validingredientstates

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing"
	"github.com/dinnerdonebetter/backend/internal/search"
	searchcfg "github.com/dinnerdonebetter/backend/internal/search/config"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "valid_preparations_service"
)

var _ types.ValidIngredientStateDataService = (*service)(nil)

type (
	// service handles valid ingredient states.
	service struct {
		cfg                             *Config
		logger                          logging.Logger
		validIngredientStateDataManager types.ValidIngredientStateDataManager
		validIngredientStateIDFetcher   func(*http.Request) string
		sessionContextDataFetcher       func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher            messagequeue.Publisher
		encoderDecoder                  encoding.ServerEncoderDecoder
		tracer                          tracing.Tracer
		searchIndex                     search.IndexSearcher[types.ValidIngredientStateSearchSubset]
	}
)

// ProvideService builds a new ValidIngredientStatesService.
func ProvideService(
	ctx context.Context,
	logger logging.Logger,
	cfg *Config,
	searchConfig *searchcfg.Config,
	validIngredientStateDataManager types.ValidIngredientStateDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.ValidIngredientStateDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up valid ingredient states service data changes publisher: %w", err)
	}

	searchIndex, err := searchcfg.ProvideIndex[types.ValidIngredientStateSearchSubset](ctx, logger, tracerProvider, searchConfig, search.IndexTypeValidIngredientStates)
	if err != nil {
		return nil, observability.PrepareError(err, nil, "initializing valid ingredient state index manager")
	}

	svc := &service{
		cfg:                             cfg,
		logger:                          logging.EnsureLogger(logger).WithName(serviceName),
		validIngredientStateIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(ValidIngredientStateIDURIParamKey),
		sessionContextDataFetcher:       authservice.FetchContextFromRequest,
		validIngredientStateDataManager: validIngredientStateDataManager,
		dataChangesPublisher:            dataChangesPublisher,
		encoderDecoder:                  encoder,
		tracer:                          tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		searchIndex:                     searchIndex,
	}

	return svc, nil
}
