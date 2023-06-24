package apiclient

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/requests"

	"github.com/gorilla/websocket"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/oauth2"
)

type option func(*Client) error

// SetOptions sets a new option on the client.
func (c *Client) SetOptions(opts ...option) error {
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return err
		}
	}

	return nil
}

// UsingJSON sets the url on the client.
func UsingJSON() func(*Client) error {
	return func(c *Client) error {
		requestBuilder, err := requests.NewBuilder(c.url, c.logger, tracing.NewNoopTracerProvider(), encoding.ProvideClientEncoder(c.logger, tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON))
		if err != nil {
			return err
		}

		c.requestBuilder = requestBuilder
		c.encoder = encoding.ProvideClientEncoder(c.logger, tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		return nil
	}
}

// UsingXML sets the url on the client.
func UsingXML() func(*Client) error {
	return func(c *Client) error {
		requestBuilder, err := requests.NewBuilder(c.url, c.logger, tracing.NewNoopTracerProvider(), encoding.ProvideClientEncoder(c.logger, tracing.NewNoopTracerProvider(), encoding.ContentTypeXML))
		if err != nil {
			return err
		}

		c.requestBuilder = requestBuilder
		c.encoder = encoding.ProvideClientEncoder(c.logger, tracing.NewNoopTracerProvider(), encoding.ContentTypeXML)

		return nil
	}
}

// UsingLogger sets the logger on the client.
func UsingLogger(logger logging.Logger) func(*Client) error {
	return func(c *Client) error {
		c.logger = logging.EnsureLogger(logger)

		return nil
	}
}

// UsingDebug sets the debug value on the client.
func UsingDebug(debug bool) func(*Client) error {
	return func(c *Client) error {
		c.debug = debug
		return nil
	}
}

// UsingTimeout sets the debug value on the client.
func UsingTimeout(timeout time.Duration) func(*Client) error {
	return func(c *Client) error {
		if timeout == 0 {
			timeout = defaultTimeout
		}

		c.authedClient.Timeout = timeout
		c.unauthenticatedClient.Timeout = timeout

		return nil
	}
}

// UsingCookie sets the authCookie value on the client.
func UsingCookie(cookie *http.Cookie) func(*Client) error {
	return func(c *Client) error {
		if cookie == nil {
			return ErrCookieRequired
		}

		crt := newCookieRoundTripper(c.logger, c.tracer, c.authedClient.Timeout, cookie)
		c.authMethod = cookieAuthMethod
		c.authedClient.Transport = crt
		c.authHeaderBuilder = crt
		c.websocketDialer = websocket.DefaultDialer
		c.authedClient = buildRetryingClient(c.authedClient, c.logger, c.tracer)

		c.logger.Debug("set client auth cookie")

		return nil
	}
}

// UsingOAuth2 sets the client to use OAuth2.
func UsingOAuth2(ctx context.Context, clientID, clientSecret string) func(*Client) error {
	genCodeChallengeS256 := func(s string) string {
		s256 := sha256.Sum256([]byte(s))
		return base64.URLEncoding.EncodeToString(s256[:])
	}

	return func(c *Client) error {
		oauth2Config := oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       []string{"household_member"},
			RedirectURL:  "http://localhost:9094/oauth2",
			Endpoint: oauth2.Endpoint{
				AuthURL:  c.URL().String() + "/oauth2/authorize",
				TokenURL: c.URL().String() + "/oauth2/token",
			},
		}

		// TODO: ephemeral server, change redirect URL?

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			oauth2Config.AuthCodeURL(
				"",
				oauth2.SetAuthURLParam("code_challenge", genCodeChallengeS256("s256example")),
				oauth2.SetAuthURLParam("code_challenge_method", "S256"),
			),
			http.NoBody,
		)
		if err != nil {
			return fmt.Errorf("failed to get oauth2 code: %w", err)
		}

		res, err := otelhttp.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("failed to get oauth2 code: %w", err)
		}
		defer func() {
			if closeErr := res.Body.Close(); closeErr != nil {
				c.logger.Error(err, "failed to close oauth2 response body")
			}
		}()

		code := res.Header.Get("code")

		token, err := oauth2Config.Exchange(ctx, code,
			oauth2.SetAuthURLParam("code_verifier", "s256example"),
		)
		if err != nil {
			return err
		}

		c.authMethod = oauth2AuthMethod
		c.authedClient.Transport = &oauth2.Transport{
			Source: oauth2.ReuseTokenSource(token, oauth2.StaticTokenSource(token)),
			Base:   otelhttp.DefaultClient.Transport,
		}

		// TODO: set authHeaderBuilder
		c.authedClient = buildRetryingClient(c.authedClient, c.logger, c.tracer)

		c.logger.Debug("set client oauth2 token")

		return nil
	}
}

// UsingPASETO sets the authCookie value on the client.
func UsingPASETO(clientID string, secretKey []byte) func(*Client) error {
	return func(c *Client) error {
		prt := newPASETORoundTripper(c, clientID, secretKey)

		c.authMethod = pasetoAuthMethod
		c.authedClient.Transport = prt
		c.authHeaderBuilder = prt
		c.websocketDialer = websocket.DefaultDialer
		c.authedClient = buildRetryingClient(c.authedClient, c.logger, c.tracer)

		return nil
	}
}
