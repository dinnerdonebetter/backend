package main

import (
	"context"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/trace"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"

	"github.com/prixfixeco/api_server/internal/config"
	customerdataconfig "github.com/prixfixeco/api_server/internal/customerdata/config"
	"github.com/prixfixeco/api_server/internal/database/queriers/postgres"
	emailconfig "github.com/prixfixeco/api_server/internal/email/config"
	msgconfig "github.com/prixfixeco/api_server/internal/messagequeue/config"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
	"github.com/prixfixeco/api_server/internal/workers"
)

const (
	dataChangesTopicName = "data_changes"
)

func buildHandler(logger logging.Logger, worker *workers.ChoresWorker) func(ctx context.Context, sqsEvent events.SQSEvent) error {
	return func(ctx context.Context, sqsEvent events.SQSEvent) error {
		logger.Debug("handling message")
		if err := worker.HandleMessage(ctx, nil); err != nil {
			return observability.PrepareError(err, nil, nil, "handling writes message")
		}

		logger.Debug("chores performed")

		return nil
	}
}

func main() {
	ctx := context.Background()
	logger := zerolog.NewZerologLogger()
	client := &http.Client{Timeout: 10 * time.Second}

	cfg, err := config.GetConfigFromParameterStore(true)
	if err != nil {
		logger.Fatal(err)
	}
	cfg.Database.RunMigrations = false

	tracerProvider := trace.NewNoopTracerProvider()
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(xray.Propagator{})

	dataManager, err := postgres.ProvideDatabaseClient(ctx, logger, &cfg.Database, tracerProvider)
	if err != nil {
		logger.Fatal(err)
	}

	publisherProvider, err := msgconfig.ProvidePublisherProvider(logger, tracerProvider, &cfg.Events)
	if err != nil {
		logger.Fatal(err)
	}

	postChoresPublisher, err := publisherProvider.ProviderPublisher(dataChangesTopicName)
	if err != nil {
		logger.Fatal(err)
	}

	emailer, err := emailconfig.ProvideEmailer(&cfg.Email, logger, client)
	if err != nil {
		logger.Fatal(err)
	}

	cdp, err := customerdataconfig.ProvideCollector(&cfg.CustomerData, logger)
	if err != nil {
		logger.Fatal(err)
	}

	choresWorker := workers.ProvideChoresWorker(
		logger,
		dataManager,
		postChoresPublisher,
		emailer,
		cdp,
		tracerProvider,
	)

	lambda.Start(otellambda.InstrumentHandler(
		buildHandler(logger, choresWorker),
		xrayconfig.WithEventToCarrier(),
		otellambda.WithPropagator(xray.Propagator{}),
		otellambda.WithTracerProvider(tracerProvider),
		// otellambda.WithFlusher(tracerProvider),
	))
}