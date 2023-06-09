package searchdataindexscheduler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/logging/zerolog"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/search/indexing"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"go.opentelemetry.io/otel"
	_ "go.uber.org/automaxprocs"
)

func init() {
	// Register a CloudEvent function with the Functions Framework
	functions.CloudEvent("ScheduleIndexOperation", ScheduleIndexOperation)
}

// ScheduleIndexOperation handles a search index schedule request.
func ScheduleIndexOperation(ctx context.Context, _ event.Event) error {
	logger := zerolog.NewZerologLogger(logging.DebugLevel)

	if strings.TrimSpace(strings.ToLower(os.Getenv("CEASE_OPERATION"))) == "true" {
		logger.Info("CEASE_OPERATION is set to true, exiting")
		return nil
	}

	cfg, err := config.GetSearchDataIndexSchedulerConfigFromGoogleCloudSecretManager(ctx)
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}

	tracerProvider, err := cfg.Observability.Tracing.ProvideTracerProvider(ctx, logger)
	if err != nil {
		logger.Error(err, "initializing tracer")
	}
	otel.SetTracerProvider(tracerProvider)

	tracer := tracing.NewTracer(tracerProvider.Tracer("search_indexer_cloud_function"))
	ctx, span := tracer.StartSpan(ctx)
	defer span.End()

	// manual db timeout until I find out what's wrong
	dbConnectionContext, cancel := context.WithTimeout(ctx, 15*time.Second)
	dataManager, err := postgres.ProvideDatabaseClient(dbConnectionContext, logger, &cfg.Database, tracerProvider)
	if err != nil {
		cancel()
		return observability.PrepareAndLogError(err, logger, span, "establishing database connection")
	}

	cancel()
	defer dataManager.Close()

	publisherProvider, err := msgconfig.ProvidePublisherProvider(logger, tracerProvider, &cfg.Events)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "configuring queue manager")
	}

	defer publisherProvider.Close()

	searchDataIndexPublisher, err := publisherProvider.ProvidePublisher(os.Getenv("SEARCH_INDEXING_TOPIC_NAME"))
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "configuring search indexing publisher")
	}

	defer searchDataIndexPublisher.Stop()

	var ids []string

	// figure out what records to join
	chosenIndex := indexing.AllIndexTypes[rand.Intn(len(indexing.AllIndexTypes))]
	logger = logger.WithValue("chosen_index_type", chosenIndex)

	logger.Info("index type chosen")

	var actionFunc func(context.Context) ([]string, error)

	switch chosenIndex {
	case indexing.IndexTypeValidPreparations:
		actionFunc = dataManager.GetValidPreparationIDsThatNeedSearchIndexing
	case indexing.IndexTypeRecipes:
		actionFunc = dataManager.GetRecipeIDsThatNeedSearchIndexing
	case indexing.IndexTypeMeals:
		actionFunc = dataManager.GetMealIDsThatNeedSearchIndexing
	case indexing.IndexTypeValidIngredients:
		actionFunc = dataManager.GetValidIngredientIDsThatNeedSearchIndexing
	case indexing.IndexTypeValidInstruments:
		actionFunc = dataManager.GetValidInstrumentIDsThatNeedSearchIndexing
	case indexing.IndexTypeValidMeasurementUnits:
		actionFunc = dataManager.GetValidMeasurementUnitIDsThatNeedSearchIndexing
	case indexing.IndexTypeValidIngredientStates:
		actionFunc = dataManager.GetValidIngredientStateIDsThatNeedSearchIndexing
	default:
		logger.Info("unhandled index type chosen, exiting")
		return nil
	}

	if actionFunc != nil {
		ids, err = actionFunc(ctx)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				observability.AcknowledgeError(err, logger, span, "getting valid ingredient state IDs that need search indexing")
			}
			return nil
		}
	} else {
		logger.Info("unspecified action function, exiting")
		return nil
	}

	if len(ids) > 0 {
		logger.WithValue("count", len(ids)).Info("publishing search index requests")
	}

	for _, id := range ids {
		indexReq := &indexing.IndexRequest{
			RowID:     id,
			IndexType: chosenIndex,
		}
		if err = searchDataIndexPublisher.Publish(ctx, indexReq); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing search index request")
		}
	}

	return nil

}