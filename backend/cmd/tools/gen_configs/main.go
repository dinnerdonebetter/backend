package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"

	analyticsconfig "github.com/dinnerdonebetter/backend/internal/analytics/config"
	"github.com/dinnerdonebetter/backend/internal/analytics/segment"
	"github.com/dinnerdonebetter/backend/internal/config"
	dbconfig "github.com/dinnerdonebetter/backend/internal/database/config"
	emailconfig "github.com/dinnerdonebetter/backend/internal/email/config"
	"github.com/dinnerdonebetter/backend/internal/email/sendgrid"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/messagequeue/redis"
	"github.com/dinnerdonebetter/backend/internal/objectstorage"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	logcfg "github.com/dinnerdonebetter/backend/internal/observability/logging/config"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing/cloudtrace"
	tracingcfg "github.com/dinnerdonebetter/backend/internal/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing/oteltracehttp"
	"github.com/dinnerdonebetter/backend/internal/routing"
	"github.com/dinnerdonebetter/backend/internal/search/algolia"
	searchcfg "github.com/dinnerdonebetter/backend/internal/search/config"
	"github.com/dinnerdonebetter/backend/internal/server/http"
	auditlogentriesservice "github.com/dinnerdonebetter/backend/internal/services/auditlogentries"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	householdinstrumentownershipsservice "github.com/dinnerdonebetter/backend/internal/services/householdinstrumentownerships"
	householdinvitationsservice "github.com/dinnerdonebetter/backend/internal/services/householdinvitations"
	householdsservice "github.com/dinnerdonebetter/backend/internal/services/households"
	mealplaneventsservice "github.com/dinnerdonebetter/backend/internal/services/mealplanevents"
	"github.com/dinnerdonebetter/backend/internal/services/mealplangrocerylistitems"
	mealplanoptionsservice "github.com/dinnerdonebetter/backend/internal/services/mealplanoptions"
	mealplanoptionvotesservice "github.com/dinnerdonebetter/backend/internal/services/mealplanoptionvotes"
	mealplansservice "github.com/dinnerdonebetter/backend/internal/services/mealplans"
	"github.com/dinnerdonebetter/backend/internal/services/mealplantasks"
	mealsservice "github.com/dinnerdonebetter/backend/internal/services/meals"
	oauth2clientsservice "github.com/dinnerdonebetter/backend/internal/services/oauth2clients"
	recipepreptasksservice "github.com/dinnerdonebetter/backend/internal/services/recipepreptasks"
	reciperatingsservice "github.com/dinnerdonebetter/backend/internal/services/reciperatings"
	recipesservice "github.com/dinnerdonebetter/backend/internal/services/recipes"
	recipestepcompletionconditionsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepcompletionconditions"
	recipestepingredientsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepingredients"
	recipestepinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepinstruments"
	recipestepproductsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepproducts"
	recipestepsservice "github.com/dinnerdonebetter/backend/internal/services/recipesteps"
	recipestepvesselsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepvessels"
	"github.com/dinnerdonebetter/backend/internal/services/servicesettingconfigurations"
	"github.com/dinnerdonebetter/backend/internal/services/servicesettings"
	useringredientpreferencesservice "github.com/dinnerdonebetter/backend/internal/services/useringredientpreferences"
	usernotificationsservice "github.com/dinnerdonebetter/backend/internal/services/usernotifications"
	usersservice "github.com/dinnerdonebetter/backend/internal/services/users"
	validingredientgroupsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientgroups"
	validingredientmeasurementunitsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientmeasurementunits"
	validingredientpreparationsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientpreparations"
	validingredientsservice "github.com/dinnerdonebetter/backend/internal/services/validingredients"
	validingredientstateingredientsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientstateingredients"
	validingredientstatesservice "github.com/dinnerdonebetter/backend/internal/services/validingredientstates"
	validinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/validinstruments"
	validmeasurementconversionsservice "github.com/dinnerdonebetter/backend/internal/services/validmeasurementunitconversions"
	validmeasurementunitsservice "github.com/dinnerdonebetter/backend/internal/services/validmeasurementunits"
	validpreparationinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/validpreparationinstruments"
	validpreparationsservice "github.com/dinnerdonebetter/backend/internal/services/validpreparations"
	validpreparationvesselsservice "github.com/dinnerdonebetter/backend/internal/services/validpreparationvessels"
	validvesselsservice "github.com/dinnerdonebetter/backend/internal/services/validvessels"
	webhooksservice "github.com/dinnerdonebetter/backend/internal/services/webhooks"
	workersservice "github.com/dinnerdonebetter/backend/internal/services/workers"
	"github.com/dinnerdonebetter/backend/internal/uploads"
)

const (
	defaultPort         = 8000
	defaultCookieDomain = ".whatever.gov"
	/* #nosec G101 */
	debugCookieSecret = "HEREISA32CHARSECRETWHICHISMADEUP"
	/* #nosec G101 */
	debugCookieSigningKey    = "DIFFERENT32CHARSECRETTHATIMADEUP"
	devPostgresDBConnDetails = "postgres://dbuser:hunter2@pgdatabase:5432/dinner-done-better?sslmode=disable"
	defaultCookieName        = authservice.DefaultCookieName

	// run modes.
	developmentEnv = "development"
	testingEnv     = "testing"

	// message provider topics.
	dataChangesTopicName = "data_changes"

	maxAttempts = 50

	contentTypeJSON               = "application/json"
	workerQueueAddress            = "worker_queue:6379"
	localOAuth2TokenEncryptionKey = debugCookieSecret
)

var (
	localRoutingConfig = routing.Config{
		Provider:               routing.ChiProvider,
		EnableCORSForLocalhost: true,
		SilenceRouteLogging:    false,
	}

	devRoutingConfig = routing.Config{
		Provider:               routing.ChiProvider,
		EnableCORSForLocalhost: true,
		SilenceRouteLogging:    false,
	}

	devEnvLogConfig = logcfg.Config{
		Level:    logging.DebugLevel,
		Provider: logcfg.ProviderSlog,
	}

	localLogConfig = logcfg.Config{
		Level:    logging.DebugLevel,
		Provider: logcfg.ProviderSlog,
	}

	localServer = http.Config{
		Debug:           true,
		HTTPPort:        defaultPort,
		StartupDeadline: time.Minute,
	}

	localCookies = authservice.CookieConfig{
		Name:       defaultCookieName,
		Domain:     defaultCookieDomain,
		HashKey:    debugCookieSecret,
		BlockKey:   debugCookieSigningKey,
		Lifetime:   authservice.DefaultCookieLifetime,
		SecureOnly: false,
	}

	localTracingConfig = tracingcfg.Config{
		Provider: tracingcfg.ProviderOtel,
		Otel: &oteltracehttp.Config{
			SpanCollectionProbability: 1,
			CollectorEndpoint:         "http://tracing-server:14268/api/traces",
			ServiceName:               "dinner_done_better_service",
		},
	}
)

func saveConfig(ctx context.Context, outputPath string, cfg *config.InstanceConfig, indent, validate bool) error {
	/* #nosec G301 */
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o777); err != nil {
		// okay, who gives a shit?
		_ = err
	}

	if validate {
		if err := cfg.ValidateWithContext(ctx, true); err != nil {
			return err
		}
	}

	var (
		output []byte
		err    error
	)

	if indent {
		output, err = json.MarshalIndent(cfg, "", "\t")
	} else {
		output, err = json.Marshal(cfg)
	}

	if err != nil {
		return err
	}

	/* #nosec G306 */
	return os.WriteFile(outputPath, output, 0o644)
}

type configFunc func(ctx context.Context, filePath string) error

var files = map[string]configFunc{
	"environments/local/config_files/service-config.json":             buildLocalDevelopmentServiceConfig(false),
	"environments/local/config_files/service-config-local.json":       buildLocalDevelopmentServiceConfig(true),
	"environments/testing/config_files/integration-tests-config.json": integrationTestConfig,
	"environments/dev/config_files/service-config.json":               devEnvironmentServerConfig,
}

func buildDevEnvironmentServerConfig() *config.InstanceConfig {
	cookieConfig := authservice.CookieConfig{
		Name:       defaultCookieName,
		Domain:     ".dinnerdonebetter.dev",
		Lifetime:   (24 * time.Hour) * 30,
		SecureOnly: true,
	}

	emailConfig := emailconfig.Config{
		Provider: emailconfig.ProviderSendgrid,
		Sendgrid: &sendgrid.Config{},
	}

	analyticsConfig := analyticsconfig.Config{
		Provider: analyticsconfig.ProviderSegment,
		Segment:  &segment.Config{APIToken: ""},
	}

	cfg := &config.InstanceConfig{
		Routing: devRoutingConfig,
		Meta: config.MetaSettings{
			Debug:   true,
			RunMode: developmentEnv,
		},
		Encoding: encoding.Config{
			ContentType: contentTypeJSON,
		},
		Events: msgconfig.Config{
			Consumers: msgconfig.MessageQueueConfig{
				Provider: msgconfig.ProviderPubSub,
			},
			Publishers: msgconfig.MessageQueueConfig{
				Provider: msgconfig.ProviderPubSub,
			},
		},
		Email:     emailConfig,
		Analytics: analyticsConfig,
		Server: http.Config{
			Debug:           true,
			HTTPPort:        defaultPort,
			StartupDeadline: time.Minute,
		},
		Search: searchcfg.Config{
			Algolia:  &algolia.Config{},
			Provider: searchcfg.AlgoliaProvider,
		},
		Database: dbconfig.Config{
			Debug:           true,
			LogQueries:      true,
			RunMigrations:   true,
			MaxPingAttempts: maxAttempts,
			PingWaitPeriod:  time.Second,
		},
		Observability: observability.Config{
			Logging: devEnvLogConfig,
			Tracing: tracingcfg.Config{
				Provider: tracingcfg.ProviderCloudTrace,
				CloudTrace: &cloudtrace.Config{
					ProjectID:                 "dinner-done-better-dev",
					ServiceName:               "dinner_done_better_api",
					SpanCollectionProbability: 1,
				},
			},
		},
		Services: config.ServicesConfig{
			AuditLogEntries: auditlogentriesservice.Config{},
			Auth: authservice.Config{
				OAuth2: authservice.OAuth2Config{
					Domain:               "https://dinnerdonebetter.dev",
					AccessTokenLifespan:  time.Hour,
					RefreshTokenLifespan: time.Hour,
					Debug:                false,
				},
				Cookies:               cookieConfig,
				Debug:                 true,
				EnableUserSignup:      true,
				MinimumUsernameLength: 3,
				MinimumPasswordLength: 8,
			},
			Users: usersservice.Config{
				DataChangesTopicName: dataChangesTopicName,
				PublicMediaURLPrefix: "https://media.dinnerdonebetter.dev/avatars",
				Uploads: uploads.Config{
					Debug: true,
					Storage: objectstorage.Config{
						UploadFilenameKey: "avatar",
						Provider:          objectstorage.GCPCloudStorageProvider,
						BucketName:        "media.dinnerdonebetter.dev",
						BucketPrefix:      "avatars/",
						GCPConfig: &objectstorage.GCPConfig{
							BucketName: "media.dinnerdonebetter.dev",
						},
					},
				},
			},
			Recipes: recipesservice.Config{
				// note, this should effectively be "https://media.dinnerdonebetter.dev" + bucket prefix
				UseSearchService:     true,
				PublicMediaURLPrefix: "https://media.dinnerdonebetter.dev/recipe_media",
				Uploads: uploads.Config{
					Debug: true,
					Storage: objectstorage.Config{
						UploadFilenameKey: "recipe_media",
						Provider:          objectstorage.GCPCloudStorageProvider,
						BucketName:        "media.dinnerdonebetter.dev",
						BucketPrefix:      "recipe_media/",
						GCPConfig: &objectstorage.GCPConfig{
							BucketName: "media.dinnerdonebetter.dev",
						},
					},
				},
			},
			RecipeSteps: recipestepsservice.Config{
				// note, this should effectively be "https://media.dinnerdonebetter.dev" + bucket prefix
				PublicMediaURLPrefix: "https://media.dinnerdonebetter.dev/recipe_media",
				Uploads: uploads.Config{
					Debug: true,
					Storage: objectstorage.Config{
						UploadFilenameKey: "recipe_media",
						Provider:          objectstorage.GCPCloudStorageProvider,
						BucketName:        "media.dinnerdonebetter.dev",
						BucketPrefix:      "recipe_media/",
						GCPConfig: &objectstorage.GCPConfig{
							BucketName: "media.dinnerdonebetter.dev",
						},
					},
				},
			},
			ValidIngredients: validingredientsservice.Config{
				UseSearchService: true,
			},
			ValidIngredientStates: validingredientstatesservice.Config{
				UseSearchService: true,
			},
			ValidInstruments: validinstrumentsservice.Config{
				UseSearchService: true,
			},
			ValidVessels: validvesselsservice.Config{
				UseSearchService: true,
			},
			ValidMeasurementUnits: validmeasurementunitsservice.Config{
				UseSearchService: true,
			},
			ValidPreparations: validpreparationsservice.Config{
				UseSearchService: true,
			},
			Meals: mealsservice.Config{
				UseSearchService: true,
			},
		},
	}

	return cfg
}

func devEnvironmentServerConfig(ctx context.Context, filePath string) error {
	cfg := buildDevEnvironmentServerConfig()

	return saveConfig(ctx, filePath, cfg, false, false)
}

func buildDevConfig() *config.InstanceConfig {
	return &config.InstanceConfig{
		Routing: localRoutingConfig,
		Meta: config.MetaSettings{
			Debug:   true,
			RunMode: developmentEnv,
		},
		Encoding: encoding.Config{
			ContentType: contentTypeJSON,
		},
		Events: msgconfig.Config{
			Consumers: msgconfig.MessageQueueConfig{
				Provider: msgconfig.ProviderRedis,
				Redis: redis.Config{
					QueueAddresses: []string{workerQueueAddress},
				},
			},
			Publishers: msgconfig.MessageQueueConfig{
				Provider: msgconfig.ProviderRedis,
				Redis: redis.Config{
					QueueAddresses: []string{workerQueueAddress},
				},
			},
		},
		Search: searchcfg.Config{
			Algolia:  &algolia.Config{},
			Provider: searchcfg.AlgoliaProvider,
		},
		Server: localServer,
		Database: dbconfig.Config{
			OAuth2TokenEncryptionKey: localOAuth2TokenEncryptionKey,
			Debug:                    true,
			RunMigrations:            true,
			LogQueries:               true,
			MaxPingAttempts:          maxAttempts,
			PingWaitPeriod:           time.Second,
			ConnectionDetails:        devPostgresDBConnDetails,
		},
		Observability: observability.Config{
			Logging: localLogConfig,
			Tracing: localTracingConfig,
		},
		Services: config.ServicesConfig{
			AuditLogEntries: auditlogentriesservice.Config{},
			Auth: authservice.Config{
				OAuth2: authservice.OAuth2Config{
					Domain:               "http://localhost:9000",
					AccessTokenLifespan:  time.Hour,
					RefreshTokenLifespan: time.Hour,
					Debug:                false,
				},
				SSO: authservice.SSOConfigs{
					Google: authservice.GoogleSSOConfig{
						CallbackURL: "https://app.dinnerdonebetter.dev/auth/google/callback",
					},
				},
				Cookies:               localCookies,
				Debug:                 true,
				EnableUserSignup:      true,
				MinimumUsernameLength: 3,
				MinimumPasswordLength: 8,
				DataChangesTopicName:  dataChangesTopicName,
			},
			Users: usersservice.Config{
				DataChangesTopicName: dataChangesTopicName,
				Uploads: uploads.Config{
					Debug: true,
					Storage: objectstorage.Config{
						UploadFilenameKey: "avatar",
						Provider:          objectstorage.FilesystemProvider,
						BucketName:        "avatars",
						FilesystemConfig: &objectstorage.FilesystemConfig{
							RootDirectory: "/uploads",
						},
					},
				},
			},
			Households: householdsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			HouseholdInvitations: householdinvitationsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			Webhooks: webhooksservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidInstruments: validinstrumentsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidVessels: validvesselsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidPreparationVessels: validpreparationvesselsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidIngredients: validingredientsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidIngredientGroups: validingredientgroupsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidPreparations: validpreparationsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			UserIngredientPreferences: useringredientpreferencesservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidIngredientStates: validingredientstatesservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidMeasurementUnits: validmeasurementunitsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidIngredientStateIngredients: validingredientstateingredientsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidIngredientPreparations: validingredientpreparationsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidPreparationInstruments: validpreparationinstrumentsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidInstrumentMeasurementUnits: validingredientmeasurementunitsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			Meals: mealsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			Recipes: recipesservice.Config{
				DataChangesTopicName: dataChangesTopicName,
				PublicMediaURLPrefix: "https://example.website.lol",
				Uploads: uploads.Config{
					Debug: true,
					Storage: objectstorage.Config{
						UploadFilenameKey: "recipe_media",
						Provider:          objectstorage.FilesystemProvider,
						BucketName:        "recipe_media",
						FilesystemConfig: &objectstorage.FilesystemConfig{
							RootDirectory: "/uploads",
						},
					},
				},
			},
			RecipeSteps: recipestepsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
				PublicMediaURLPrefix: "https://example.website.lol",
				Uploads: uploads.Config{
					Debug: true,
					Storage: objectstorage.Config{
						UploadFilenameKey: "recipe_media",
						Provider:          objectstorage.FilesystemProvider,
						BucketName:        "recipe_media",
						FilesystemConfig: &objectstorage.FilesystemConfig{
							RootDirectory: "/uploads",
						},
					},
				},
			},
			RecipeStepProducts: recipestepproductsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipeStepInstruments: recipestepinstrumentsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipeStepVessels: recipestepvesselsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipeStepIngredients: recipestepingredientsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipeStepCompletionConditions: recipestepcompletionconditionsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlans: mealplansservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlanEvents: mealplaneventsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlanOptions: mealplanoptionsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlanOptionVotes: mealplanoptionvotesservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlanTasks: mealplantasks.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipePrepTasks: recipepreptasksservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlanGroceryListItems: mealplangrocerylistitems.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidMeasurementUnitConversions: validmeasurementconversionsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ServiceSettings: servicesettings.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ServiceSettingConfigurations: servicesettingconfigurations.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipeRatings: reciperatingsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			HouseholdInstrumentOwnerships: householdinstrumentownershipsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			OAuth2Clients: oauth2clientsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			Workers: workersservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			UserNotifications: usernotificationsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
		},
	}
}

func buildLocalDevelopmentServiceConfig(local bool) func(context.Context, string) error {
	const localUploadsDir = "artifacts/uploads"
	const localRedisAddr = "localhost:6379"
	return func(ctx context.Context, filePath string) error {
		cfg := buildDevConfig()

		if local {
			cfg.Database.ConnectionDetails = "postgres://dbuser:hunter2@localhost:5432/dinner-done-better?sslmode=disable"
			cfg.Events.Consumers.Redis.QueueAddresses = []string{localRedisAddr}
			cfg.Events.Publishers.Redis.QueueAddresses = []string{localRedisAddr}
			cfg.Services.Users.Uploads.Storage.FilesystemConfig.RootDirectory = localUploadsDir
			cfg.Services.Recipes.Uploads.Storage.FilesystemConfig.RootDirectory = localUploadsDir
			cfg.Services.RecipeSteps.Uploads.Storage.FilesystemConfig.RootDirectory = localUploadsDir
		}

		return saveConfig(ctx, filePath, cfg, true, true)
	}
}

func buildIntegrationTestsConfig() *config.InstanceConfig {
	return &config.InstanceConfig{
		Routing: localRoutingConfig,
		Meta: config.MetaSettings{
			Debug:   false,
			RunMode: testingEnv,
		},
		Events: msgconfig.Config{
			Consumers: msgconfig.MessageQueueConfig{
				Provider: msgconfig.ProviderRedis,
				Redis: redis.Config{
					QueueAddresses: []string{workerQueueAddress},
				},
			},
			Publishers: msgconfig.MessageQueueConfig{
				Provider: msgconfig.ProviderRedis,
				Redis: redis.Config{
					QueueAddresses: []string{workerQueueAddress},
				},
			},
		},
		Encoding: encoding.Config{
			ContentType: contentTypeJSON,
		},
		Server: http.Config{
			Debug:           false,
			HTTPPort:        defaultPort,
			StartupDeadline: time.Minute,
		},
		Database: dbconfig.Config{
			OAuth2TokenEncryptionKey: localOAuth2TokenEncryptionKey,
			Debug:                    true,
			RunMigrations:            true,
			LogQueries:               true,
			MaxPingAttempts:          maxAttempts,
			PingWaitPeriod:           1500 * time.Millisecond,
			ConnectionDetails:        devPostgresDBConnDetails,
		},
		Observability: observability.Config{
			Logging: logcfg.Config{
				Level:    logging.InfoLevel,
				Provider: logcfg.ProviderSlog,
			},
			Tracing: localTracingConfig,
		},
		Services: config.ServicesConfig{
			AuditLogEntries: auditlogentriesservice.Config{},
			Auth: authservice.Config{
				OAuth2: authservice.OAuth2Config{
					Domain:               "http://localhost:9000",
					AccessTokenLifespan:  time.Hour,
					RefreshTokenLifespan: time.Hour,
					Debug:                false,
				},
				Cookies: authservice.CookieConfig{
					Name:       defaultCookieName,
					Domain:     defaultCookieDomain,
					HashKey:    debugCookieSecret,
					BlockKey:   debugCookieSigningKey,
					Lifetime:   authservice.DefaultCookieLifetime,
					SecureOnly: false,
				},
				Debug:                 false,
				EnableUserSignup:      true,
				MinimumUsernameLength: 3,
				MinimumPasswordLength: 8,
				DataChangesTopicName:  dataChangesTopicName,
			},
			Users: usersservice.Config{
				DataChangesTopicName: dataChangesTopicName,
				Uploads: uploads.Config{
					Debug: false,
					Storage: objectstorage.Config{
						Provider:   "memory",
						BucketName: "avatars",
						S3Config:   nil,
					},
				},
			},
			Households: householdsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			HouseholdInvitations: householdinvitationsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			Webhooks: webhooksservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidInstruments: validinstrumentsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidVessels: validvesselsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidPreparationVessels: validpreparationvesselsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidIngredients: validingredientsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidIngredientGroups: validingredientgroupsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidPreparations: validpreparationsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			UserIngredientPreferences: useringredientpreferencesservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidIngredientStates: validingredientstatesservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidMeasurementUnits: validmeasurementunitsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidIngredientStateIngredients: validingredientstateingredientsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidIngredientPreparations: validingredientpreparationsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidPreparationInstruments: validpreparationinstrumentsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidInstrumentMeasurementUnits: validingredientmeasurementunitsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			Meals: mealsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			Recipes: recipesservice.Config{
				DataChangesTopicName: dataChangesTopicName,
				PublicMediaURLPrefix: "https://media.example.website/lol",
				Uploads: uploads.Config{
					Debug: false,
					Storage: objectstorage.Config{
						Provider:   "memory",
						BucketName: "recipes",
						S3Config:   nil,
					},
				},
			},
			RecipeSteps: recipestepsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
				PublicMediaURLPrefix: "https://media.example.website/lol",
				Uploads: uploads.Config{
					Debug: false,
					Storage: objectstorage.Config{
						Provider:   "memory",
						BucketName: "recipes",
						S3Config:   nil,
					},
				},
			},
			RecipeStepInstruments: recipestepinstrumentsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipeStepVessels: recipestepvesselsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipeStepIngredients: recipestepingredientsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipeStepCompletionConditions: recipestepcompletionconditionsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipeStepProducts: recipestepproductsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlans: mealplansservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlanEvents: mealplaneventsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlanOptions: mealplanoptionsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlanOptionVotes: mealplanoptionvotesservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlanTasks: mealplantasks.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipePrepTasks: recipepreptasksservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			MealPlanGroceryListItems: mealplangrocerylistitems.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ValidMeasurementUnitConversions: validmeasurementconversionsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ServiceSettings: servicesettings.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			ServiceSettingConfigurations: servicesettingconfigurations.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			RecipeRatings: reciperatingsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			HouseholdInstrumentOwnerships: householdinstrumentownershipsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			Workers: workersservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
			UserNotifications: usernotificationsservice.Config{
				DataChangesTopicName: dataChangesTopicName,
			},
		},
	}
}

func integrationTestConfig(ctx context.Context, filePath string) error {
	cfg := buildIntegrationTestsConfig()

	return saveConfig(ctx, filePath, cfg, true, true)
}

func main() {
	ctx := context.Background()

	for filePath, fun := range files {
		if err := fun(ctx, filePath); err != nil {
			log.Fatalf("error rendering %s: %v", filePath, err)
		}
	}
}
