package launchdarkly

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dinnerdonebetter/backend/internal/featureflags"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/launchdarkly/go-sdk-common/v3/ldcontext"
	ld "github.com/launchdarkly/go-server-sdk/v6"
	"github.com/launchdarkly/go-server-sdk/v6/ldcomponents"
)

var (
	ErrMissingHTTPClient = errors.New("missing HTTP launchDarklyClient")
	ErrNilConfig         = errors.New("missing config")
	ErrMissingSDKKey     = errors.New("missing SDK key")
)

type (
	Config struct {
		SDKKey      string        `json:"sdkKey" mapstructure:"sdk_key" toml:"sdk_key"`
		InitTimeout time.Duration `json:"initTimeout" mapstructure:"init_timeout" toml:"init_timeout"`
	}

	launchDarklyClient interface {
		BoolVariation(key string, context ldcontext.Context, defaultVal bool) (bool, error)
		Identify(context ldcontext.Context) error
	}

	// featureFlagManager implements the feature flag interface.
	featureFlagManager struct {
		launchDarklyClient launchDarklyClient
		logger             logging.Logger
		tracer             tracing.Tracer
	}
)

// NewFeatureFlagManager constructs a new featureFlagManager.
func NewFeatureFlagManager(cfg *Config, logger logging.Logger, tracerProvider tracing.TracerProvider, httpClient *http.Client, configModifiers ...func(ld.Config) ld.Config) (featureflags.FeatureFlagManager, error) {
	if httpClient == nil {
		return nil, ErrMissingHTTPClient
	}

	if cfg == nil {
		return nil, ErrNilConfig
	}

	if cfg.SDKKey == "" {
		return nil, ErrMissingSDKKey
	}

	if cfg.InitTimeout == time.Duration(0) {
		cfg.InitTimeout = 5 * time.Second
	}

	ldConfig := ld.Config{
		HTTP: ldcomponents.HTTPConfiguration().HTTPClientFactory(func() *http.Client { return httpClient }),
	}

	for _, modifier := range configModifiers {
		ldConfig = modifier(ldConfig)
	}

	client, err := ld.MakeCustomClient(
		cfg.SDKKey,
		ldConfig,
		cfg.InitTimeout,
	)
	if err != nil {
		return nil, fmt.Errorf("error initializing LaunchDarkly launchDarklyClient: %w", err)
	}

	ffm := &featureFlagManager{
		logger:             logger,
		tracer:             tracing.NewTracer(tracerProvider.Tracer("launchdarkly_feature_flag_manager")),
		launchDarklyClient: client,
	}

	return ffm, nil
}

// CanUseFeature returns whether a user can use a feature or not.
func (f *featureFlagManager) CanUseFeature(ctx context.Context, userID, feature string) (bool, error) {
	_, span := tracing.StartSpan(ctx)
	defer span.End()

	return f.launchDarklyClient.BoolVariation(feature, ldcontext.New(userID), false)
}

// Identify identifies a user in LaunchDarkly.
func (f *featureFlagManager) Identify(ctx context.Context, user *types.User) error {
	_, span := tracing.StartSpan(ctx)
	defer span.End()

	return f.launchDarklyClient.Identify(
		ldcontext.NewBuilderFromContext(ldcontext.New(user.ID)).
			Name(user.Username).
			SetString("email", user.EmailAddress).
			SetString("first_name", user.FirstName).
			SetString("last_name", user.LastName).
			Build(),
	)
}
