package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	dbconfig "github.com/dinnerdonebetter/backend/internal/database/config"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/observability/logging/config"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"
	"github.com/dinnerdonebetter/backend/internal/search"
	"github.com/dinnerdonebetter/backend/internal/search/algolia"
	searchcfg "github.com/dinnerdonebetter/backend/internal/search/config"
	"github.com/dinnerdonebetter/backend/internal/search/indexing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

func main() {
	indicesPtr := flag.String("indices", "", "indices to initialize")
	wipePtr := flag.Bool("wipe", false, "whether to wipe the indices or not")

	flag.Parse()

	indices := strings.Split(*indicesPtr, ",")
	if len(indices) == 0 {
		log.Fatal("indices are required")
	}

	ctx := context.Background()

	logger := (&loggingcfg.Config{Level: logging.DebugLevel, Provider: loggingcfg.ProviderSlog}).ProvideLogger()

	tracerProvider := tracing.NewNoopTracerProvider()

	cfg := &searchcfg.Config{
		Provider: searchcfg.AlgoliaProvider,
		Algolia: &algolia.Config{
			AppID:  os.Getenv("ALGOLIA_APP_ID"),
			APIKey: os.Getenv("ALGOLIA_API_KEY"),
		},
	}

	dbConfig := &dbconfig.Config{
		ConnectionDetails: os.Getenv("DATABASE_URL"),
	}

	dataManager, err := postgres.ProvideDatabaseClient(ctx, logger, tracerProvider, dbConfig)
	if dataManager != nil {
		defer dataManager.Close()
	}

	if err != nil {
		log.Println(fmt.Errorf("initializing database client: %w", err))
		return
	}

	var (
		im               search.IndexManager
		indexRequestChan = make(chan *indexing.IndexRequest)
		wipeOnce         sync.Once
		waitGroup        sync.WaitGroup
	)

	go func() {
		for x := range indexRequestChan {
			wipeOnce.Do(func() {
				if *wipePtr {
					log.Println("wiping index")
					if err = im.Wipe(ctx); err != nil {
						log.Println(fmt.Errorf("wiping index: %w", err))
						return
					}
					log.Println("wiped index")
				}
			})

			if err = indexing.HandleIndexRequest(ctx, logger, tracerProvider, cfg, dataManager, x); err != nil {
				observability.AcknowledgeError(err, logger, nil, "indexing row")
			}

			waitGroup.Done()
		}
	}()

	for i, index := range indices {
		if i > 0 {
			waitGroup.Wait()
		}

		filter := types.DefaultQueryFilter()
		filter.Limit = pointer.To(uint8(50))
		thresholdMet := false

		switch index {
		case search.IndexTypeRecipes:
			im, err = searchcfg.ProvideIndex[types.RecipeSearchSubset](ctx, logger, tracerProvider, cfg, index)
			if err != nil {
				observability.AcknowledgeError(err, logger, nil, "initializing index manager")
				return
			}

			for !thresholdMet {
				var data *types.QueryFilteredResult[types.Recipe]
				data, err = dataManager.GetRecipes(ctx, filter)
				if err != nil {
					log.Println(fmt.Errorf("getting Recipe data: %w", err))
					return
				}

				for _, x := range data.Data {
					indexRequestChan <- &indexing.IndexRequest{
						RowID:     x.ID,
						IndexType: search.IndexTypeRecipes,
					}
					waitGroup.Add(1)
				}

				thresholdMet = len(data.Data) == 0
				*filter.Page++
			}
		case search.IndexTypeMeals:
			im, err = searchcfg.ProvideIndex[types.MealSearchSubset](ctx, logger, tracerProvider, cfg, index)
			if err != nil {
				observability.AcknowledgeError(err, logger, nil, "initializing index manager")
				return
			}

			for !thresholdMet {
				var data *types.QueryFilteredResult[types.Meal]
				data, err = dataManager.GetMeals(ctx, filter)
				if err != nil {
					log.Println(fmt.Errorf("getting Meal data: %w", err))
					return
				}

				for _, x := range data.Data {
					indexRequestChan <- &indexing.IndexRequest{
						RowID:     x.ID,
						IndexType: search.IndexTypeMeals,
					}
					waitGroup.Add(1)
				}

				thresholdMet = len(data.Data) == 0
				*filter.Page++
			}
		case search.IndexTypeValidIngredients:
			im, err = searchcfg.ProvideIndex[types.ValidIngredientSearchSubset](ctx, logger, tracerProvider, cfg, index)
			if err != nil {
				observability.AcknowledgeError(err, logger, nil, "initializing index manager")
				return
			}

			for !thresholdMet {
				var data *types.QueryFilteredResult[types.ValidIngredient]
				data, err = dataManager.GetValidIngredients(ctx, filter)
				if err != nil {
					log.Println(fmt.Errorf("getting ValidIngredient data: %w", err))
					return
				}

				for _, x := range data.Data {
					indexRequestChan <- &indexing.IndexRequest{
						RowID:     x.ID,
						IndexType: search.IndexTypeValidIngredients,
					}
					waitGroup.Add(1)
				}

				thresholdMet = len(data.Data) == 0
				*filter.Page++
			}
		case search.IndexTypeValidInstruments:
			im, err = searchcfg.ProvideIndex[types.ValidInstrumentSearchSubset](ctx, logger, tracerProvider, cfg, index)
			if err != nil {
				observability.AcknowledgeError(err, logger, nil, "initializing index manager")
				return
			}

			for !thresholdMet {
				var data *types.QueryFilteredResult[types.ValidInstrument]
				data, err = dataManager.GetValidInstruments(ctx, filter)
				if err != nil {
					log.Println(fmt.Errorf("getting ValidInstrument data: %w", err))
					return
				}

				for _, x := range data.Data {
					indexRequestChan <- &indexing.IndexRequest{
						RowID:     x.ID,
						IndexType: search.IndexTypeValidInstruments,
					}
					waitGroup.Add(1)
				}

				thresholdMet = len(data.Data) == 0
				*filter.Page++
			}
		case search.IndexTypeValidMeasurementUnits:
			im, err = searchcfg.ProvideIndex[types.ValidMeasurementUnitSearchSubset](ctx, logger, tracerProvider, cfg, index)
			if err != nil {
				observability.AcknowledgeError(err, logger, nil, "initializing index manager")
				return
			}

			for !thresholdMet {
				var data *types.QueryFilteredResult[types.ValidMeasurementUnit]
				data, err = dataManager.GetValidMeasurementUnits(ctx, filter)
				if err != nil {
					log.Println(fmt.Errorf("getting ValidMeasurementUnit data: %w", err))
					return
				}

				for _, x := range data.Data {
					indexRequestChan <- &indexing.IndexRequest{
						RowID:     x.ID,
						IndexType: search.IndexTypeValidMeasurementUnits,
					}
					waitGroup.Add(1)
				}

				thresholdMet = len(data.Data) == 0
				*filter.Page++
			}
		case search.IndexTypeValidPreparations:
			im, err = searchcfg.ProvideIndex[types.ValidPreparationSearchSubset](ctx, logger, tracerProvider, cfg, index)
			if err != nil {
				observability.AcknowledgeError(err, logger, nil, "initializing index manager")
				return
			}

			for !thresholdMet {
				var data *types.QueryFilteredResult[types.ValidPreparation]
				data, err = dataManager.GetValidPreparations(ctx, filter)
				if err != nil {
					log.Println(fmt.Errorf("getting ValidPreparation data: %w", err))
					return
				}

				for _, x := range data.Data {
					indexRequestChan <- &indexing.IndexRequest{
						RowID:     x.ID,
						IndexType: search.IndexTypeValidPreparations,
					}
					waitGroup.Add(1)
				}

				thresholdMet = len(data.Data) == 0
				*filter.Page++
			}
		case search.IndexTypeValidIngredientStates:
			im, err = searchcfg.ProvideIndex[types.ValidIngredientStateSearchSubset](ctx, logger, tracerProvider, cfg, index)
			if err != nil {
				observability.AcknowledgeError(err, logger, nil, "initializing index manager")
				return
			}

			for !thresholdMet {
				var data *types.QueryFilteredResult[types.ValidIngredientState]
				data, err = dataManager.GetValidIngredientStates(ctx, filter)
				if err != nil {
					log.Println(fmt.Errorf("getting ValidIngredientState data: %w", err))
					return
				}

				for _, x := range data.Data {
					indexRequestChan <- &indexing.IndexRequest{
						RowID:     x.ID,
						IndexType: search.IndexTypeValidIngredientStates,
					}
					waitGroup.Add(1)
				}

				thresholdMet = len(data.Data) == 0
				*filter.Page++
			}
		case search.IndexTypeValidVessels:
			im, err = searchcfg.ProvideIndex[types.ValidVesselSearchSubset](ctx, logger, tracerProvider, cfg, index)
			if err != nil {
				observability.AcknowledgeError(err, logger, nil, "initializing index manager")
				return
			}

			for !thresholdMet {
				var data *types.QueryFilteredResult[types.ValidVessel]
				data, err = dataManager.GetValidVessels(ctx, filter)
				if err != nil {
					log.Println(fmt.Errorf("getting ValidVessel data: %w", err))
					return
				}

				for _, x := range data.Data {
					indexRequestChan <- &indexing.IndexRequest{
						RowID:     x.ID,
						IndexType: search.IndexTypeValidVessels,
					}
					waitGroup.Add(1)
				}

				thresholdMet = len(data.Data) == 0
				*filter.Page++
			}
		case search.IndexTypeUsers:
			im, err = searchcfg.ProvideIndex[types.UserSearchSubset](ctx, logger, tracerProvider, cfg, index)
			if err != nil {
				observability.AcknowledgeError(err, logger, nil, "initializing index manager")
				return
			}

			for !thresholdMet {
				var data *types.QueryFilteredResult[types.User]
				data, err = dataManager.GetUsers(ctx, filter)
				if err != nil {
					log.Println(fmt.Errorf("getting user data: %w", err))
					return
				}

				for _, x := range data.Data {
					indexRequestChan <- &indexing.IndexRequest{
						RowID:     x.ID,
						IndexType: search.IndexTypeUsers,
					}
					waitGroup.Add(1)
				}

				thresholdMet = len(data.Data) == 0
				*filter.Page++
			}
		}
	}
}
