package apiclient

import (
	"context"
	"net/http"
	"time"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/pkg/types"
)

func (c *Client) fetchAuthTokenForAPIClient(ctx context.Context, httpClient *http.Client, clientID string, secretKey []byte) (string, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if clientID == "" {
		return "", ErrEmptyInputProvided
	}

	if secretKey == nil {
		return "", ErrNilInputProvided
	}

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if httpClient.Timeout == 0 {
		httpClient.Timeout = defaultTimeout
	}

	input := &types.PASETOCreationInput{
		ClientID:    clientID,
		RequestTime: time.Now().UTC().UnixNano(),
	}

	if c.householdID != "" {
		input.HouseholdID = c.householdID
	}

	logger := c.logger.WithValue(keys.APIClientClientIDKey, clientID)
	logger.Debug("fetching auth token")

	if err := input.ValidateWithContext(ctx); err != nil {
		return "", observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildAPIClientAuthTokenRequest(ctx, input, secretKey)
	if err != nil {
		return "", observability.PrepareAndLogError(err, logger, span, "building request")
	}

	// use the default client here because we want a transport that doesn't worry about cookies or tokens.
	res, err := c.fetchResponseToRequest(ctx, httpClient, req)
	if err != nil {
		return "", observability.PrepareAndLogError(err, logger, span, "executing request")
	}

	if err = errorFromResponse(res); err != nil {
		return "", observability.PrepareAndLogError(err, logger, span, "erroneous response")
	}

	var tokenRes types.PASETOResponse

	if err = c.unmarshalBody(ctx, res, &tokenRes); err != nil {
		return "", observability.PrepareAndLogError(err, logger, span, "unmarshalling body")
	}

	logger.Debug("auth token received")

	return tokenRes.Token, nil
}