package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	logcfg "github.com/dinnerdonebetter/backend/internal/observability/logging/config"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/server/http/utils"
	"github.com/dinnerdonebetter/backend/pkg/apiclient"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	testutils "github.com/dinnerdonebetter/backend/tests/utils"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/require"
)

func logJSON(t *testing.T, x any) {
	t.Helper()

	rawBytes, err := json.Marshal(x)
	require.NoError(t, err)

	t.Log(string(rawBytes))
}

func requireNotNilAndNoProblems(t *testing.T, i any, err error) {
	t.Helper()

	require.NoError(t, err)
	require.NotNil(t, i)
}

func createUserAndClientForTest(ctx context.Context, t *testing.T, input *types.UserRegistrationInput) (user *types.User, cookie *http.Cookie, client, oauthedClient *apiclient.Client) {
	t.Helper()

	if input == nil {
		input = &types.UserRegistrationInput{
			EmailAddress: gofakeit.Email(),
			Username:     fakes.BuildFakeUser().Username,
			Password:     gofakeit.Password(true, true, true, true, false, 64),
		}
	}

	user, err := testutils.CreateServiceUser(ctx, urlToUse, input)
	require.NoError(t, err)

	cookie, err = testutils.GetLoginCookie(ctx, urlToUse, user)
	require.NoError(t, err)

	code, err := totp.GenerateCode(strings.ToUpper(user.TwoFactorSecret), time.Now().UTC())
	require.NoError(t, err)

	loginInput := &types.UserLoginInput{
		Username:  user.Username,
		Password:  user.HashedPassword,
		TOTPToken: code,
	}

	client, err = initializeCookiePoweredClient(ctx, loginInput)
	require.NoError(t, err)

	oauthedClient, err = initializeOAuth2PoweredClient(ctx, cookie)
	require.NoError(t, err)

	return user, cookie, client, oauthedClient
}

func initializeCookiePoweredClient(ctx context.Context, loginInput *types.UserLoginInput) (*apiclient.Client, error) {
	if parsedURLToUse == nil {
		panic("url not set!")
	}

	logger := (&logcfg.Config{Provider: logcfg.ProviderSlog}).ProvideLogger()

	c, err := apiclient.NewClient(
		parsedURLToUse,
		tracing.NewNoopTracerProvider(),
		apiclient.UsingLogger(logger),
		apiclient.UsingTracingProvider(tracing.NewNoopTracerProvider()),
		apiclient.UsingURL(urlToUse),
		apiclient.UsingLogin(ctx, loginInput),
	)
	if err != nil {
		return nil, err
	}

	if debug {
		if setOptionErr := c.SetOptions(apiclient.UsingDebug(true)); setOptionErr != nil {
			return nil, setOptionErr
		}
	}

	return c, nil
}

func initializeOAuth2PoweredClient(ctx context.Context, cookie *http.Cookie) (*apiclient.Client, error) {
	if parsedURLToUse == nil {
		panic("url not set!")
	}

	logger := (&logcfg.Config{Provider: logcfg.ProviderSlog}).ProvideLogger()

	c, err := apiclient.NewClient(
		parsedURLToUse,
		tracing.NewNoopTracerProvider(),
		apiclient.UsingLogger(logger),
		apiclient.UsingTracingProvider(tracing.NewNoopTracerProvider()),
		apiclient.UsingURL(urlToUse),
		apiclient.UsingOAuth2(ctx, createdClientID, createdClientSecret, cookie),
	)
	if err != nil {
		return nil, err
	}

	if debug {
		if setOptionErr := c.SetOptions(apiclient.UsingDebug(true)); setOptionErr != nil {
			return nil, setOptionErr
		}
	}

	return c, nil
}

func buildSimpleClient(t *testing.T) *apiclient.Client {
	t.Helper()

	c, err := apiclient.NewClient(
		parsedURLToUse,
		tracing.NewNoopTracerProvider(),
		apiclient.UsingTracingProvider(tracing.NewNoopTracerProvider()),
		apiclient.UsingURL(urlToUse),
	)
	require.NoError(t, err)

	return c
}

func generateTOTPTokenForUser(t *testing.T, u *types.User) string {
	t.Helper()

	code, err := totp.GenerateCode(u.TwoFactorSecret, time.Now().UTC())
	require.NotEmpty(t, code)
	require.NoError(t, err)

	return code
}

func buildAdminCookieAndOAuthedClients(ctx context.Context, t *testing.T) (cookieClient *apiclient.Client, oauthedClient *apiclient.Client) {
	t.Helper()

	ctx, span := tracing.StartSpan(ctx)
	defer span.End()

	u := serverutils.DetermineServiceURL()
	urlToUse = u.String()

	logger := (&logcfg.Config{Provider: logcfg.ProviderSlog}).ProvideLogger()
	logger.WithValue(keys.URLKey, urlToUse).Info("checking server")

	serverutils.EnsureServerIsUp(ctx, urlToUse)

	adminCode, err := totp.GenerateCode(strings.ToUpper(premadeAdminUser.TwoFactorSecret), time.Now().UTC())
	require.NoError(t, err)

	loginInput := &types.UserLoginInput{
		Username:  premadeAdminUser.Username,
		Password:  premadeAdminUser.HashedPassword,
		TOTPToken: adminCode,
	}

	adminCookieClient, err := initializeCookiePoweredClient(ctx, loginInput)
	require.NoError(t, err)

	cookie, err := testutils.GetLoginCookie(ctx, urlToUse, premadeAdminUser)
	require.NoError(t, err)

	oauthedClient, err = initializeOAuth2PoweredClient(ctx, cookie)
	require.NoError(t, err)

	return adminCookieClient, oauthedClient
}
