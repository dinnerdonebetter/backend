package accounts

import (
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/metrics"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/internal/routing"
	"gitlab.com/prixfixe/prixfixe/internal/search"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	counterName        metrics.CounterName = "accounts"
	counterDescription string              = "the number of accounts managed by the accounts service"
	serviceName        string              = "accounts_service"
)

var _ types.AccountDataService = (*service)(nil)

type (
	// SearchIndex is a type alias for dependency injection's sake.
	SearchIndex search.IndexManager

	// service handles to-do list accounts.
	service struct {
		logger                       logging.Logger
		accountDataManager           types.AccountDataManager
		accountMembershipDataManager types.AccountUserMembershipDataManager
		accountIDFetcher             func(*http.Request) uint64
		userIDFetcher                func(*http.Request) uint64
		sessionContextDataFetcher    func(*http.Request) (*types.SessionContextData, error)
		accountCounter               metrics.UnitCounter
		encoderDecoder               encoding.ServerEncoderDecoder
		tracer                       tracing.Tracer
	}
)

// ProvideService builds a new AccountsService.
func ProvideService(
	logger logging.Logger,
	accountDataManager types.AccountDataManager,
	accountMembershipDataManager types.AccountUserMembershipDataManager,
	encoder encoding.ServerEncoderDecoder,
	counterProvider metrics.UnitCounterProvider,
	routeParamManager routing.RouteParamManager,
) types.AccountDataService {
	return &service{
		logger:                       logging.EnsureLogger(logger).WithName(serviceName),
		accountIDFetcher:             routeParamManager.BuildRouteParamIDFetcher(logger, AccountIDURIParamKey, "account"),
		userIDFetcher:                routeParamManager.BuildRouteParamIDFetcher(logger, UserIDURIParamKey, "user"),
		sessionContextDataFetcher:    authservice.FetchContextFromRequest,
		accountDataManager:           accountDataManager,
		accountMembershipDataManager: accountMembershipDataManager,
		encoderDecoder:               encoder,
		accountCounter:               metrics.EnsureUnitCounter(counterProvider, logger, counterName, counterDescription),
		tracer:                       tracing.NewTracer(serviceName),
	}
}
