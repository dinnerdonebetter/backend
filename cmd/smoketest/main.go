package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/brianvoe/gofakeit/v5"
	"go.opentelemetry.io/otel/trace"

	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
	"github.com/prixfixeco/api_server/pkg/client/httpclient"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func main() {
	ctx := context.Background()
	logger := zerolog.NewZerologLogger()
	urlToUse := "https://api.prixfixe.dev"
	parsedURLToUse, err := url.Parse(urlToUse)
	if err != nil {
		panic(err)
	}

	input := &types.UserRegistrationInput{
		EmailAddress: gofakeit.Email(),
		Username:     fakes.BuildFakeUser().Username,
		Password:     gofakeit.Password(true, true, true, true, true, 64),
	}

	user, err := testutils.CreateServiceUser(ctx, urlToUse, input)
	if err != nil {
		panic(err)
	}

	cookie, err := testutils.GetLoginCookie(ctx, urlToUse, user)
	if err != nil {
		panic(err)
	}

	c, err := httpclient.NewClient(parsedURLToUse,
		trace.NewNoopTracerProvider(),
		httpclient.UsingLogger(logger),
		httpclient.UsingCookie(cookie),
	)
	if err != nil {
		panic(err)
	}

	stopChan := make(chan bool, 1)
	notificationsChan, err := c.SubscribeToNotifications(ctx, stopChan)
	if err != nil {
		panic(err)
	}

	// Create webhook.
	exampleWebhook := fakes.BuildFakeWebhook()
	exampleWebhookInput := fakes.BuildFakeWebhookCreationInputFromWebhook(exampleWebhook)
	createdWebhookID, err := c.CreateWebhook(ctx, exampleWebhookInput)
	if err != nil {
		panic(err)
	}

	logger.Info("waiting for webhook creation notification")
	<-notificationsChan

	webhook, err := c.GetWebhook(ctx, createdWebhookID)
	if err != nil {
		panic(err)
	}

	if err = c.ArchiveWebhook(ctx, createdWebhookID); err != nil {
		panic(err)
	}

	logger.Info("waiting for webhook archive notification")
	<-notificationsChan

	fmt.Println(webhook)
}
