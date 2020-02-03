package main

import (
	"context"
	"errors"
	"math/rand"
	"net/http"
	"time"

	client "gitlab.com/prixfixe/prixfixe/client/v1/http"
	randmodel "gitlab.com/prixfixe/prixfixe/tests/v1/testutil/rand/model"
)

var (
	// ErrUnavailableYet is a sentinel error value
	ErrUnavailableYet = errors.New("can't do this yet")
)

type (
	// actionFunc represents a thing you can do
	actionFunc func() (*http.Request, error)

	// Action is a wrapper struct around some important values
	Action struct {
		Action actionFunc
		Weight int
		Name   string
	}
)

// RandomAction takes a client and returns a closure which is an action
func RandomAction(c *client.V1Client) *Action {
	ctx := context.Background()
	allActions := map[string]*Action{
		"GetHealthCheck": {
			Name:   "GetHealthCheck",
			Action: c.BuildHealthCheckRequest,
			Weight: 100,
		},
		"CreateUser": {
			Name: "CreateUser",
			Action: func() (*http.Request, error) {
				ui := randmodel.RandomUserInput()
				return c.BuildCreateUserRequest(ctx, ui)
			},
			Weight: 100,
		},
	}

	for k, v := range buildInstrumentActions(c) {
		allActions[k] = v
	}

	for k, v := range buildIngredientActions(c) {
		allActions[k] = v
	}

	for k, v := range buildPreparationActions(c) {
		allActions[k] = v
	}

	for k, v := range buildRequiredPreparationInstrumentActions(c) {
		allActions[k] = v
	}

	for k, v := range buildRecipeActions(c) {
		allActions[k] = v
	}

	for k, v := range buildRecipeStepActions(c) {
		allActions[k] = v
	}

	for k, v := range buildRecipeStepInstrumentActions(c) {
		allActions[k] = v
	}

	for k, v := range buildRecipeStepIngredientActions(c) {
		allActions[k] = v
	}

	for k, v := range buildRecipeStepProductActions(c) {
		allActions[k] = v
	}

	for k, v := range buildRecipeIterationActions(c) {
		allActions[k] = v
	}

	for k, v := range buildRecipeStepEventActions(c) {
		allActions[k] = v
	}

	for k, v := range buildIterationMediaActions(c) {
		allActions[k] = v
	}

	for k, v := range buildInvitationActions(c) {
		allActions[k] = v
	}

	for k, v := range buildReportActions(c) {
		allActions[k] = v
	}

	for k, v := range buildWebhookActions(c) {
		allActions[k] = v
	}

	for k, v := range buildOAuth2ClientActions(c) {
		allActions[k] = v
	}

	totalWeight := 0
	for _, rb := range allActions {
		totalWeight += rb.Weight
	}

	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(totalWeight)

	for _, rb := range allActions {
		r -= rb.Weight
		if r <= 0 {
			return rb
		}
	}

	return nil
}
