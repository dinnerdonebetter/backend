/*
Command index_initializer is a CLI that takes in some data via flags about your
database and the type you want to index, and hydrates a Bleve index full of that type.
This tool is to be used in the event of some data corruption that takes the search index
out of commission.
*/
package main

import (
	"context"
	"log"
	"time"

	config "gitlab.com/prixfixe/prixfixe/internal/config"
	"gitlab.com/prixfixe/prixfixe/internal/database"
	dbconfig "gitlab.com/prixfixe/prixfixe/internal/database/config"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/search"
	"gitlab.com/prixfixe/prixfixe/internal/search/bleve"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	flag "github.com/spf13/pflag"
)

var (
	indexOutputPath string
	typeName        string

	dbConnectionDetails string
	databaseType        string

	deadline time.Duration

	validTypeNames = map[string]struct{}{
		"valid_instrument":             {},
		"valid_preparation":            {},
		"valid_ingredient":             {},
		"valid_ingredient_preparation": {},
		"valid_preparation_instrument": {},
		"recipe":                       {},
		"recipe_step":                  {},
		"recipe_step_ingredient":       {},
		"recipe_step_product":          {},
		"invitation":                   {},
		"report":                       {},
	}

	validDatabaseTypes = map[string]struct{}{
		dbconfig.PostgresProvider: {},
	}
)

const (
	outputPathVerboseFlagName   = "output"
	dbConnectionVerboseFlagName = "db_connection"
	dbTypeVerboseFlagName       = "db_type"
)

func init() {
	flag.StringVarP(&indexOutputPath, outputPathVerboseFlagName, "o", "", "output path for bleve index")
	flag.StringVarP(&typeName, "type", "t", "", "which type to create bleve index for")

	flag.StringVarP(&dbConnectionDetails, dbConnectionVerboseFlagName, "c", "", "connection string for the relevant database")
	flag.StringVarP(&databaseType, dbTypeVerboseFlagName, "b", "", "which type of database to connect to")

	flag.DurationVarP(&deadline, "deadline", "d", time.Minute, "amount of time to spend adding to the index")
}

func main() {
	flag.Parse()
	logger := logging.ProvideLogger(logging.Config{
		Name:     "search_index_initializer",
		Provider: logging.ProviderZerolog,
	})
	ctx := context.Background()

	if indexOutputPath == "" {
		log.Fatalf("No output path specified, please provide one via the --%s flag", outputPathVerboseFlagName)
		return
	} else if _, ok := validTypeNames[typeName]; !ok {
		log.Fatalf("Invalid type name %q specified", typeName)
		return
	} else if dbConnectionDetails == "" {
		log.Fatalf("No database connection details %q specified, please provide one via the --%s flag", dbConnectionDetails, dbConnectionVerboseFlagName)
		return
	} else if _, ok := validDatabaseTypes[databaseType]; !ok {
		log.Fatalf("Invalid database type %q specified, please provide one via the --%s flag", databaseType, dbTypeVerboseFlagName)
		return
	}

	im, err := bleve.NewBleveIndexManager(search.IndexPath(indexOutputPath), search.IndexName(typeName), logger)
	if err != nil {
		log.Fatal(err)
	}

	cfg := &config.InstanceConfig{
		Database: dbconfig.Config{
			MetricsCollectionInterval: time.Second,
			Provider:                  databaseType,
			ConnectionDetails:         database.ConnectionDetails(dbConnectionDetails),
		},
	}

	// connect to our database.
	logger.Debug("connecting to database")
	rawDB, err := dbconfig.ProvideDatabaseConnection(logger, &cfg.Database)
	if err != nil {
		log.Fatalf("error establishing connection to database: %v", err)
	}

	// establish the database client.
	logger.Debug("setting up database client")
	dbClient, err := config.ProvideDatabaseClient(ctx, logger, rawDB, cfg)
	if err != nil {
		log.Fatalf("error initializing database client: %v", err)
	}

	switch typeName {
	case "valid_instrument":
		outputChan := make(chan []*types.ValidInstrument)
		if queryErr := dbClient.GetAllValidInstruments(ctx, outputChan, 1000); queryErr != nil {
			log.Fatalf("error fetching valid instruments from database: %v", err)
		}

		for {
			select {
			case validInstruments := <-outputChan:
				for _, x := range validInstruments {
					if searchIndexErr := im.Index(ctx, x.ID, x); searchIndexErr != nil {
						logger.WithValue("id", x.ID).Error(searchIndexErr, "error adding to search index")
					}
				}
			case <-time.After(deadline):
				logger.Info("terminating")
				return
			}
		}
	case "valid_preparation":
		outputChan := make(chan []*types.ValidPreparation)
		if queryErr := dbClient.GetAllValidPreparations(ctx, outputChan, 1000); queryErr != nil {
			log.Fatalf("error fetching valid preparations from database: %v", err)
		}

		for {
			select {
			case validPreparations := <-outputChan:
				for _, x := range validPreparations {
					if searchIndexErr := im.Index(ctx, x.ID, x); searchIndexErr != nil {
						logger.WithValue("id", x.ID).Error(searchIndexErr, "error adding to search index")
					}
				}
			case <-time.After(deadline):
				logger.Info("terminating")
				return
			}
		}
	case "valid_ingredient":
		outputChan := make(chan []*types.ValidIngredient)
		if queryErr := dbClient.GetAllValidIngredients(ctx, outputChan, 1000); queryErr != nil {
			log.Fatalf("error fetching valid ingredients from database: %v", err)
		}

		for {
			select {
			case validIngredients := <-outputChan:
				for _, x := range validIngredients {
					if searchIndexErr := im.Index(ctx, x.ID, x); searchIndexErr != nil {
						logger.WithValue("id", x.ID).Error(searchIndexErr, "error adding to search index")
					}
				}
			case <-time.After(deadline):
				logger.Info("terminating")
				return
			}
		}

	default:
		log.Fatal("this should never occur")
	}
}
