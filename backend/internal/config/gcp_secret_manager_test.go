package config

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	analyticsconfig "github.com/dinnerdonebetter/backend/internal/analytics/config"
	"github.com/dinnerdonebetter/backend/internal/analytics/segment"
	dbconfig "github.com/dinnerdonebetter/backend/internal/database/config"
	emailconfig "github.com/dinnerdonebetter/backend/internal/email/config"
	"github.com/dinnerdonebetter/backend/internal/email/sendgrid"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/routing"
	"github.com/dinnerdonebetter/backend/internal/server/http"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	householdinvitationsservice "github.com/dinnerdonebetter/backend/internal/services/householdinvitations"
	householdsservice "github.com/dinnerdonebetter/backend/internal/services/households"
	mealplanoptionsservice "github.com/dinnerdonebetter/backend/internal/services/mealplanoptions"
	mealplanoptionvotesservice "github.com/dinnerdonebetter/backend/internal/services/mealplanoptionvotes"
	mealplansservice "github.com/dinnerdonebetter/backend/internal/services/mealplans"
	mealsservice "github.com/dinnerdonebetter/backend/internal/services/meals"
	recipepreptasksservice "github.com/dinnerdonebetter/backend/internal/services/recipepreptasks"
	recipesservice "github.com/dinnerdonebetter/backend/internal/services/recipes"
	recipestepcompletionconditionsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepcompletionconditions"
	recipestepingredientsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepingredients"
	recipestepinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepinstruments"
	recipestepproductsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepproducts"
	recipestepsservice "github.com/dinnerdonebetter/backend/internal/services/recipesteps"
	"github.com/dinnerdonebetter/backend/internal/services/servicesettingconfigurations"
	"github.com/dinnerdonebetter/backend/internal/services/servicesettings"
	usersservice "github.com/dinnerdonebetter/backend/internal/services/users"
	validingredientmeasurementunitsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientmeasurementunits"
	validingredientpreparationsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientpreparations"
	validingredientsservice "github.com/dinnerdonebetter/backend/internal/services/validingredients"
	validingredientstatesservice "github.com/dinnerdonebetter/backend/internal/services/validingredientstates"
	validinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/validinstruments"
	"github.com/dinnerdonebetter/backend/internal/services/validmeasurementunits"
	validpreparationinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/validpreparationinstruments"
	validpreparationsservice "github.com/dinnerdonebetter/backend/internal/services/validpreparations"
	webhooksservice "github.com/dinnerdonebetter/backend/internal/services/webhooks"
	testutils "github.com/dinnerdonebetter/backend/tests/utils"

	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/googleapis/gax-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockSecretVersionAccessor struct {
	mock.Mock
}

func (m *mockSecretVersionAccessor) AccessSecretVersion(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error) {
	args := m.Called(ctx, req, opts)

	return args.Get(0).(*secretmanagerpb.AccessSecretVersionResponse), args.Error(1)
}

func TestGetAPIServerConfigFromGoogleCloudRunEnvironment(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		baseConfig := &InstanceConfig{
			Observability: observability.Config{},
			Email: emailconfig.Config{
				Sendgrid: &sendgrid.Config{},
				Provider: emailconfig.ProviderSendgrid,
			},
			Analytics: analyticsconfig.Config{
				Provider: analyticsconfig.ProviderSegment,
				Segment:  &segment.Config{},
			},
			Encoding: encoding.Config{ContentType: "application/json"},
			Routing:  routing.Config{},
			Database: dbconfig.Config{},
			Meta:     MetaSettings{},
			Events:   msgconfig.Config{},
			Server: http.Config{
				StartupDeadline: time.Second,
			},
			Services: ServicesConfig{
				RecipeStepProducts:              recipestepproductsservice.Config{},
				RecipeStepInstruments:           recipestepinstrumentsservice.Config{},
				RecipeStepIngredients:           recipestepingredientsservice.Config{},
				MealPlans:                       mealplansservice.Config{},
				MealPlanOptions:                 mealplanoptionsservice.Config{},
				Households:                      householdsservice.Config{},
				HouseholdInvitations:            householdinvitationsservice.Config{},
				Webhooks:                        webhooksservice.Config{},
				Users:                           usersservice.Config{},
				RecipePrepTasks:                 recipepreptasksservice.Config{},
				RecipeStepCompletionConditions:  recipestepcompletionconditionsservice.Config{},
				ValidIngredientStates:           validingredientstatesservice.Config{},
				ServiceSettings:                 servicesettings.Config{},
				ServiceSettingConfigurations:    servicesettingconfigurations.Config{},
				ValidMeasurementUnits:           validmeasurementunits.Config{},
				ValidInstruments:                validinstrumentsservice.Config{},
				ValidIngredients:                validingredientsservice.Config{},
				ValidPreparations:               validpreparationsservice.Config{},
				MealPlanOptionVotes:             mealplanoptionvotesservice.Config{},
				ValidIngredientPreparations:     validingredientpreparationsservice.Config{},
				ValidPreparationInstruments:     validpreparationinstrumentsservice.Config{},
				ValidInstrumentMeasurementUnits: validingredientmeasurementunitsservice.Config{},
				Meals:                           mealsservice.Config{},
				Recipes: recipesservice.Config{
					PublicMediaURLPrefix: t.Name(),
				},
				RecipeSteps: recipestepsservice.Config{
					PublicMediaURLPrefix: t.Name(),
				},
				Auth: authservice.Config{
					MinimumPasswordLength: 8,
					MinimumUsernameLength: 8,
				},
			},
		}

		f, err := os.CreateTemp("", "testing.json")
		require.NoError(t, err)

		require.NoError(t, json.NewEncoder(f).Encode(baseConfig))

		require.NoError(t, os.Setenv(gcpConfigFilePathEnvVarKey, f.Name()))
		require.NoError(t, os.Setenv(gcpPortEnvVarKey, "1234"))
		require.NoError(t, os.Setenv(gcpDatabaseSocketDirEnvVarKey, "/example/blah"))
		require.NoError(t, os.Setenv(gcpDatabaseUserEnvVarKey, "user"))
		require.NoError(t, os.Setenv(gcpDatabaseUserPasswordEnvVarKey, "hunter2"))
		require.NoError(t, os.Setenv(gcpDatabaseNameEnvVarKey, "fake_db_name"))
		require.NoError(t, os.Setenv(gcpDatabaseInstanceConnNameEnvVarKey, "fake_conn_name"))
		require.NoError(t, os.Setenv(gcpCookieHashKeyEnvVarKey, "fake_cookie_hash_key"))
		require.NoError(t, os.Setenv(gcpCookieBlockKeyEnvVarKey, "fake_cookie_block_key"))
		require.NoError(t, os.Setenv(gcpSendgridTokenEnvVarKey, "fake_sendgrid_token"))
		require.NoError(t, os.Setenv(gcpPostHogKeyEnvVarKey, "fake_posthog_api_key"))
		require.NoError(t, os.Setenv(gcpSegmentTokenEnvVarKey, "fake_segment_token"))
		require.NoError(t, os.Setenv(gcpAlgoliaAPIKeyEnvVarKey, "fake_algolia_api_key"))
		require.NoError(t, os.Setenv(gcpAlgoliaAppIDEnvVarKey, "fake_algolia_app_id"))
		require.NoError(t, os.Setenv(gcpOauth2TokenEncryptionKeyEnvVarKey, "fake_oauth2_token_encryption_key"))
		require.NoError(t, os.Setenv(gcpGoogleSSOClientIDEnvVarKey, "fake_google_sso_client_id"))
		require.NoError(t, os.Setenv(gcpGoogleSSOClientSecretEnvVarKey, "fake_google_sso_client_secret"))

		ctx := context.Background()
		client := &mockSecretVersionAccessor{}

		client.On(
			"AccessSecretVersion",
			testutils.ContextMatcher,
			&secretmanagerpb.AccessSecretVersionRequest{Name: buildSecretPathForGCPSecretStore(dataChangesTopicAccessName)},
			[]gax.CallOption(nil),
		).Return(
			&secretmanagerpb.AccessSecretVersionResponse{
				Name: dataChangesTopicAccessName,
				Payload: &secretmanagerpb.SecretPayload{
					Data: []byte("this_is_the_big_secret"),
				},
			},
			nil,
		)

		cfg, err := GetAPIServerConfigFromGoogleCloudRunEnvironment(ctx, client)
		assert.NoError(t, err)
		assert.NotNil(t, cfg)

		require.NoError(t, os.Unsetenv(gcpConfigFilePathEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpPortEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpDatabaseSocketDirEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpDatabaseUserEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpDatabaseUserPasswordEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpDatabaseNameEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpDatabaseInstanceConnNameEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpCookieHashKeyEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpCookieBlockKeyEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpSendgridTokenEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpSegmentTokenEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpAlgoliaAPIKeyEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpAlgoliaAppIDEnvVarKey))

		mock.AssertExpectationsForObjects(t, client)
	})
}
