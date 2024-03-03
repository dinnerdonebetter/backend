package householdinvitations

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/email"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/random"
	"github.com/dinnerdonebetter/backend/internal/routing"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	householdsservice "github.com/dinnerdonebetter/backend/internal/services/households"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "household_invitations_service"
)

var (
	_ types.HouseholdInvitationDataService = (*service)(nil)
)

type (
	// service handles webhooks.
	service struct {
		logger                         logging.Logger
		userDataManager                types.UserDataManager
		householdInvitationDataManager types.HouseholdInvitationDataManager
		tracer                         tracing.Tracer
		encoderDecoder                 encoding.ServerEncoderDecoder
		emailer                        email.Emailer
		secretGenerator                random.Generator
		dataChangesPublisher           messagequeue.Publisher
		householdIDFetcher             func(*http.Request) string
		householdInvitationIDFetcher   func(*http.Request) string
		sessionContextDataFetcher      func(*http.Request) (*types.SessionContextData, error)
	}
)

// ProvideHouseholdInvitationsService builds a new HouseholdInvitationDataService.
func ProvideHouseholdInvitationsService(
	logger logging.Logger,
	cfg *Config,
	userDataManager types.UserDataManager,
	householdInvitationDataManager types.HouseholdInvitationDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
	emailer email.Emailer,
	secretGenerator random.Generator,
) (types.HouseholdInvitationDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up household invitations service data changes publisher: %w", err)
	}

	s := &service{
		logger:                         logging.EnsureLogger(logger).WithName(serviceName),
		userDataManager:                userDataManager,
		householdInvitationDataManager: householdInvitationDataManager,
		encoderDecoder:                 encoder,
		dataChangesPublisher:           dataChangesPublisher,
		emailer:                        emailer,
		secretGenerator:                secretGenerator,
		sessionContextDataFetcher:      authservice.FetchContextFromRequest,
		householdIDFetcher:             routeParamManager.BuildRouteParamStringIDFetcher(householdsservice.HouseholdIDURIParamKey),
		householdInvitationIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(HouseholdInvitationIDURIParamKey),
		tracer:                         tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return s, nil
}
