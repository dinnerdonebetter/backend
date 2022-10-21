package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prixfixeco/api_server/internal/authentication"
	"github.com/prixfixeco/api_server/internal/email"
	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/messagequeue"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/metrics"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/internal/random"
	"github.com/prixfixeco/api_server/internal/routing"
	authservice "github.com/prixfixeco/api_server/internal/services/authentication"
	"github.com/prixfixeco/api_server/internal/storage"
	"github.com/prixfixeco/api_server/internal/uploads"
	"github.com/prixfixeco/api_server/internal/uploads/images"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	serviceName        = "users_service"
	counterDescription = "number of users managed by the users service"
	counterName        = metrics.CounterName("users")
)

var _ types.UserDataService = (*service)(nil)

type (
	// RequestValidator validates request.
	RequestValidator interface {
		Validate(req *http.Request) (bool, error)
	}

	// service handles our users.
	service struct {
		userDataManager                types.UserDataManager
		householdDataManager           types.HouseholdDataManager
		householdInvitationDataManager types.HouseholdInvitationDataManager
		authSettings                   *authservice.Config
		authenticator                  authentication.Authenticator
		logger                         logging.Logger
		encoderDecoder                 encoding.ServerEncoderDecoder
		userIDFetcher                  func(*http.Request) string
		sessionContextDataFetcher      func(*http.Request) (*types.SessionContextData, error)
		userCounter                    metrics.UnitCounter
		secretGenerator                random.Generator
		imageUploadProcessor           images.ImageUploadProcessor
		uploadManager                  uploads.UploadManager
		dataChangesPublisher           messagequeue.Publisher
		tracer                         tracing.Tracer
		passwordResetTokenDataManager  types.PasswordResetTokenDataManager
		emailer                        email.Emailer
	}
)

// ProvideUsersService builds a new UsersService.
func ProvideUsersService(
	ctx context.Context,
	cfg *Config,
	authSettings *authservice.Config,
	logger logging.Logger,
	userDataManager types.UserDataManager,
	householdDataManager types.HouseholdDataManager,
	householdInvitationDataManager types.HouseholdInvitationDataManager,
	authenticator authentication.Authenticator,
	encoder encoding.ServerEncoderDecoder,
	counterProvider metrics.UnitCounterProvider,
	imageUploadProcessor images.ImageUploadProcessor,
	routeParamManager routing.RouteParamManager,
	tracerProvider tracing.TracerProvider,
	publisherProvider messagequeue.PublisherProvider,
	secretGenerator random.Generator,
	passwordResetTokenDataManager types.PasswordResetTokenDataManager,
	emailer email.Emailer,
) (types.UserDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProviderPublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up users service data changes publisher: %w", err)
	}

	uploadManager, err := storage.NewUploadManager(ctx, logger, tracerProvider, &cfg.Uploads.Storage, routeParamManager)
	if err != nil {
		return nil, fmt.Errorf("initializing uploadManager: %w", err)
	}

	s := &service{
		logger:                         logging.EnsureLogger(logger).WithName(serviceName),
		userDataManager:                userDataManager,
		householdDataManager:           householdDataManager,
		householdInvitationDataManager: householdInvitationDataManager,
		authenticator:                  authenticator,
		userIDFetcher:                  routeParamManager.BuildRouteParamStringIDFetcher(UserIDURIParamKey),
		sessionContextDataFetcher:      authservice.FetchContextFromRequest,
		encoderDecoder:                 encoder,
		authSettings:                   authSettings,
		userCounter:                    metrics.EnsureUnitCounter(counterProvider, logger, counterName, counterDescription),
		secretGenerator:                secretGenerator,
		tracer:                         tracing.NewTracer(tracerProvider.Tracer(serviceName)),
		imageUploadProcessor:           imageUploadProcessor,
		uploadManager:                  uploadManager,
		dataChangesPublisher:           dataChangesPublisher,
		passwordResetTokenDataManager:  passwordResetTokenDataManager,
		emailer:                        emailer,
	}

	return s, nil
}
