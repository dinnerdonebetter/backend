package indexing

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/search"
	"github.com/dinnerdonebetter/backend/internal/search/config"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	testutils "github.com/dinnerdonebetter/backend/tests/utils"

	"github.com/stretchr/testify/assert"
)

func TestHandleIndexRequest(T *testing.T) {
	T.Parallel()

	T.Run("user index type", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		tracerProvider := tracing.NewNoopTracerProvider
		searchConfig := &config.Config{}

		dataManager := database.NewMockDatabase()
		dataManager.UserDataManagerMock.On("GetUser", testutils.ContextMatcher, exampleUser.ID).Return(exampleUser, nil)
		dataManager.UserDataManagerMock.On("MarkUserAsIndexed", testutils.ContextMatcher, exampleUser.ID).Return(nil)

		indexReq := &IndexRequest{
			RowID:     exampleUser.ID,
			IndexType: search.IndexTypeUsers,
			Delete:    false,
		}

		assert.NoError(t, HandleIndexRequest(ctx, logger, tracerProvider(), searchConfig, dataManager, indexReq))
	})

	T.Run("deleting user index type", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		tracerProvider := tracing.NewNoopTracerProvider
		searchConfig := &config.Config{}

		dataManager := database.NewMockDatabase()
		dataManager.UserDataManagerMock.On("GetUser", testutils.ContextMatcher, exampleUser.ID).Return(exampleUser, nil)
		dataManager.UserDataManagerMock.On("MarkUserAsIndexed", testutils.ContextMatcher, exampleUser.ID).Return(nil)

		indexReq := &IndexRequest{
			RowID:     exampleUser.ID,
			IndexType: search.IndexTypeUsers,
			Delete:    true,
		}

		assert.NoError(t, HandleIndexRequest(ctx, logger, tracerProvider(), searchConfig, dataManager, indexReq))
	})

	T.Run("recipe index type", func(t *testing.T) {
		t.Parallel()

		exampleRecipe := fakes.BuildFakeRecipe()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		tracerProvider := tracing.NewNoopTracerProvider
		searchConfig := &config.Config{}

		dataManager := database.NewMockDatabase()
		dataManager.RecipeDataManagerMock.On("GetRecipe", testutils.ContextMatcher, exampleRecipe.ID).Return(exampleRecipe, nil)
		dataManager.RecipeDataManagerMock.On("MarkRecipeAsIndexed", testutils.ContextMatcher, exampleRecipe.ID).Return(nil)

		indexReq := &IndexRequest{
			RowID:     exampleRecipe.ID,
			IndexType: search.IndexTypeRecipes,
			Delete:    false,
		}

		assert.NoError(t, HandleIndexRequest(ctx, logger, tracerProvider(), searchConfig, dataManager, indexReq))
	})

	T.Run("meal index type", func(t *testing.T) {
		t.Parallel()

		exampleMeal := fakes.BuildFakeMeal()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		tracerProvider := tracing.NewNoopTracerProvider
		searchConfig := &config.Config{}

		dataManager := database.NewMockDatabase()
		dataManager.MealDataManagerMock.On("GetMeal", testutils.ContextMatcher, exampleMeal.ID).Return(exampleMeal, nil)
		dataManager.MealDataManagerMock.On("MarkMealAsIndexed", testutils.ContextMatcher, exampleMeal.ID).Return(nil)

		indexReq := &IndexRequest{
			RowID:     exampleMeal.ID,
			IndexType: search.IndexTypeMeals,
			Delete:    false,
		}

		assert.NoError(t, HandleIndexRequest(ctx, logger, tracerProvider(), searchConfig, dataManager, indexReq))
	})

	T.Run("valid vessel index type", func(t *testing.T) {
		t.Parallel()

		exampleValidVessel := fakes.BuildFakeValidVessel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		tracerProvider := tracing.NewNoopTracerProvider
		searchConfig := &config.Config{}

		dataManager := database.NewMockDatabase()
		dataManager.ValidVesselDataManagerMock.On("GetValidVessel", testutils.ContextMatcher, exampleValidVessel.ID).Return(exampleValidVessel, nil)
		dataManager.ValidVesselDataManagerMock.On("MarkValidVesselAsIndexed", testutils.ContextMatcher, exampleValidVessel.ID).Return(nil)

		indexReq := &IndexRequest{
			RowID:     exampleValidVessel.ID,
			IndexType: search.IndexTypeValidVessels,
			Delete:    false,
		}

		assert.NoError(t, HandleIndexRequest(ctx, logger, tracerProvider(), searchConfig, dataManager, indexReq))
	})

	T.Run("valid ingredient index type", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredient := fakes.BuildFakeValidIngredient()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		tracerProvider := tracing.NewNoopTracerProvider
		searchConfig := &config.Config{}

		dataManager := database.NewMockDatabase()
		dataManager.ValidIngredientDataManagerMock.On("GetValidIngredient", testutils.ContextMatcher, exampleValidIngredient.ID).Return(exampleValidIngredient, nil)
		dataManager.ValidIngredientDataManagerMock.On("MarkValidIngredientAsIndexed", testutils.ContextMatcher, exampleValidIngredient.ID).Return(nil)

		indexReq := &IndexRequest{
			RowID:     exampleValidIngredient.ID,
			IndexType: search.IndexTypeValidIngredients,
			Delete:    false,
		}

		assert.NoError(t, HandleIndexRequest(ctx, logger, tracerProvider(), searchConfig, dataManager, indexReq))
	})

	T.Run("valid instrument index type", func(t *testing.T) {
		t.Parallel()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		tracerProvider := tracing.NewNoopTracerProvider
		searchConfig := &config.Config{}

		dataManager := database.NewMockDatabase()
		dataManager.ValidInstrumentDataManagerMock.On("GetValidInstrument", testutils.ContextMatcher, exampleValidInstrument.ID).Return(exampleValidInstrument, nil)
		dataManager.ValidInstrumentDataManagerMock.On("MarkValidInstrumentAsIndexed", testutils.ContextMatcher, exampleValidInstrument.ID).Return(nil)

		indexReq := &IndexRequest{
			RowID:     exampleValidInstrument.ID,
			IndexType: search.IndexTypeValidInstruments,
			Delete:    false,
		}

		assert.NoError(t, HandleIndexRequest(ctx, logger, tracerProvider(), searchConfig, dataManager, indexReq))
	})

	T.Run("valid preparation index type", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparation := fakes.BuildFakeValidPreparation()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		tracerProvider := tracing.NewNoopTracerProvider
		searchConfig := &config.Config{}

		dataManager := database.NewMockDatabase()
		dataManager.ValidPreparationDataManagerMock.On("GetValidPreparation", testutils.ContextMatcher, exampleValidPreparation.ID).Return(exampleValidPreparation, nil)
		dataManager.ValidPreparationDataManagerMock.On("MarkValidPreparationAsIndexed", testutils.ContextMatcher, exampleValidPreparation.ID).Return(nil)

		indexReq := &IndexRequest{
			RowID:     exampleValidPreparation.ID,
			IndexType: search.IndexTypeValidPreparations,
			Delete:    false,
		}

		assert.NoError(t, HandleIndexRequest(ctx, logger, tracerProvider(), searchConfig, dataManager, indexReq))
	})

	T.Run("valid measurement unit index type", func(t *testing.T) {
		t.Parallel()

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		tracerProvider := tracing.NewNoopTracerProvider
		searchConfig := &config.Config{}

		dataManager := database.NewMockDatabase()
		dataManager.ValidMeasurementUnitDataManagerMock.On("GetValidMeasurementUnit", testutils.ContextMatcher, exampleValidMeasurementUnit.ID).Return(exampleValidMeasurementUnit, nil)
		dataManager.ValidMeasurementUnitDataManagerMock.On("MarkValidMeasurementUnitAsIndexed", testutils.ContextMatcher, exampleValidMeasurementUnit.ID).Return(nil)

		indexReq := &IndexRequest{
			RowID:     exampleValidMeasurementUnit.ID,
			IndexType: search.IndexTypeValidMeasurementUnits,
			Delete:    false,
		}

		assert.NoError(t, HandleIndexRequest(ctx, logger, tracerProvider(), searchConfig, dataManager, indexReq))
	})

	T.Run("valid ingredient state index type", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientState := fakes.BuildFakeValidIngredientState()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		tracerProvider := tracing.NewNoopTracerProvider
		searchConfig := &config.Config{}

		dataManager := database.NewMockDatabase()
		dataManager.ValidIngredientStateDataManagerMock.On("GetValidIngredientState", testutils.ContextMatcher, exampleValidIngredientState.ID).Return(exampleValidIngredientState, nil)
		dataManager.ValidIngredientStateDataManagerMock.On("MarkValidIngredientStateAsIndexed", testutils.ContextMatcher, exampleValidIngredientState.ID).Return(nil)

		indexReq := &IndexRequest{
			RowID:     exampleValidIngredientState.ID,
			IndexType: search.IndexTypeValidIngredientStates,
			Delete:    false,
		}

		assert.NoError(t, HandleIndexRequest(ctx, logger, tracerProvider(), searchConfig, dataManager, indexReq))
	})
}
