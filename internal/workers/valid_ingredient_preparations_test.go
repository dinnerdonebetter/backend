package workers

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/database"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/publishers/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/search"
	mocksearch "github.com/prixfixeco/api_server/internal/search/mock"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func TestWritesWorker_createValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreWriteMessage{
			DataType:                   types.ValidIngredientPreparationDataType,
			ValidIngredientPreparation: fakes.BuildFakeValidIngredientPreparationDatabaseCreationInput(),
		}

		expectedValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientPreparationDataManager.On(
			"CreateValidIngredientPreparation",
			testutils.ContextMatcher,
			body.ValidIngredientPreparation,
		).Return(expectedValidIngredientPreparation, nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Index",
			testutils.ContextMatcher,
			expectedValidIngredientPreparation.ID,
			expectedValidIngredientPreparation,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		worker, err := ProvideWritesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.NoError(t, worker.createValidIngredientPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with error writing", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreWriteMessage{
			DataType:                   types.ValidIngredientPreparationDataType,
			ValidIngredientPreparation: fakes.BuildFakeValidIngredientPreparationDatabaseCreationInput(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientPreparationDataManager.On(
			"CreateValidIngredientPreparation",
			testutils.ContextMatcher,
			body.ValidIngredientPreparation,
		).Return((*types.ValidIngredientPreparation)(nil), errors.New("blah"))

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, nil
		}

		postArchivesPublisher := &mockpublishers.Publisher{}

		worker, err := ProvideWritesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.createValidIngredientPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error updating search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreWriteMessage{
			DataType:                   types.ValidIngredientPreparationDataType,
			ValidIngredientPreparation: fakes.BuildFakeValidIngredientPreparationDatabaseCreationInput(),
		}

		expectedValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientPreparationDataManager.On(
			"CreateValidIngredientPreparation",
			testutils.ContextMatcher,
			body.ValidIngredientPreparation,
		).Return(expectedValidIngredientPreparation, nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Index",
			testutils.ContextMatcher,
			expectedValidIngredientPreparation.ID,
			expectedValidIngredientPreparation,
		).Return(errors.New("blah"))

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		postArchivesPublisher := &mockpublishers.Publisher{}

		worker, err := ProvideWritesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.createValidIngredientPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with error publishing data change message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreWriteMessage{
			DataType:                   types.ValidIngredientPreparationDataType,
			ValidIngredientPreparation: fakes.BuildFakeValidIngredientPreparationDatabaseCreationInput(),
		}

		expectedValidIngredientPreparation := fakes.BuildFakeValidIngredientPreparation()

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientPreparationDataManager.On(
			"CreateValidIngredientPreparation",
			testutils.ContextMatcher,
			body.ValidIngredientPreparation,
		).Return(expectedValidIngredientPreparation, nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Index",
			testutils.ContextMatcher,
			expectedValidIngredientPreparation.ID,
			expectedValidIngredientPreparation,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		worker, err := ProvideWritesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.createValidIngredientPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})
}

func TestWritesWorker_updateValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreUpdateMessage{
			DataType:                   types.ValidIngredientPreparationDataType,
			ValidIngredientPreparation: fakes.BuildFakeValidIngredientPreparation(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientPreparationDataManager.On(
			"UpdateValidIngredientPreparation",
			testutils.ContextMatcher,
			body.ValidIngredientPreparation,
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Index",
			testutils.ContextMatcher,
			body.ValidIngredientPreparation.ID,
			body.ValidIngredientPreparation,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		worker, err := ProvideUpdatesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.NoError(t, worker.updateValidIngredientPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with error updating valid ingredient preparation", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreUpdateMessage{
			DataType:                   types.ValidIngredientPreparationDataType,
			ValidIngredientPreparation: fakes.BuildFakeValidIngredientPreparation(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientPreparationDataManager.On(
			"UpdateValidIngredientPreparation",
			testutils.ContextMatcher,
			body.ValidIngredientPreparation,
		).Return(errors.New("blah"))

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, nil
		}

		postArchivesPublisher := &mockpublishers.Publisher{}

		worker, err := ProvideUpdatesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.updateValidIngredientPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error updating search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreUpdateMessage{
			DataType:                   types.ValidIngredientPreparationDataType,
			ValidIngredientPreparation: fakes.BuildFakeValidIngredientPreparation(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientPreparationDataManager.On(
			"UpdateValidIngredientPreparation",
			testutils.ContextMatcher,
			body.ValidIngredientPreparation,
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Index",
			testutils.ContextMatcher,
			body.ValidIngredientPreparation.ID,
			body.ValidIngredientPreparation,
		).Return(errors.New("blah"))

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		postArchivesPublisher := &mockpublishers.Publisher{}

		worker, err := ProvideUpdatesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.updateValidIngredientPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with error publishing data change event", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreUpdateMessage{
			DataType:                   types.ValidIngredientPreparationDataType,
			ValidIngredientPreparation: fakes.BuildFakeValidIngredientPreparation(),
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientPreparationDataManager.On(
			"UpdateValidIngredientPreparation",
			testutils.ContextMatcher,
			body.ValidIngredientPreparation,
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Index",
			testutils.ContextMatcher,
			body.ValidIngredientPreparation.ID,
			body.ValidIngredientPreparation,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		worker, err := ProvideUpdatesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.updateValidIngredientPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})
}

func TestWritesWorker_archiveValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.ValidIngredientPreparationDataType,
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientPreparationDataManager.On(
			"ArchiveValidIngredientPreparation",
			testutils.ContextMatcher,
			body.ValidIngredientPreparationID,
		).Return(nil)

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.ValidIngredientPreparationID,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvideArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.NoError(t, worker.archiveValidIngredientPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with error archiving", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.ValidIngredientPreparationDataType,
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientPreparationDataManager.On(
			"ArchiveValidIngredientPreparation",
			testutils.ContextMatcher,
			body.ValidIngredientPreparationID,
		).Return(errors.New("blah"))

		postArchivesPublisher := &mockpublishers.Publisher{}
		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return nil, nil
		}

		worker, err := ProvideArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.archiveValidIngredientPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error removing from search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.ValidIngredientPreparationDataType,
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientPreparationDataManager.On(
			"ArchiveValidIngredientPreparation",
			testutils.ContextMatcher,
			body.ValidIngredientPreparationID,
		).Return(nil)

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.ValidIngredientPreparationID,
		).Return(errors.New("blah"))

		postArchivesPublisher := &mockpublishers.Publisher{}

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvideArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.archiveValidIngredientPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})

	T.Run("with error publishing post-archive message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		client := &http.Client{}

		body := &types.PreArchiveMessage{
			DataType: types.ValidIngredientPreparationDataType,
		}

		dbManager := database.NewMockDatabase()
		dbManager.ValidIngredientPreparationDataManager.On(
			"ArchiveValidIngredientPreparation",
			testutils.ContextMatcher,
			body.ValidIngredientPreparationID,
		).Return(nil)

		postArchivesPublisher := &mockpublishers.Publisher{}
		postArchivesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			mock.MatchedBy(func(message *types.DataChangeMessage) bool { return true }),
		).Return(errors.New("blah"))

		searchIndexManager := &mocksearch.IndexManager{}
		searchIndexManager.On(
			"Delete",
			testutils.ContextMatcher,
			body.ValidIngredientPreparationID,
		).Return(nil)

		searchIndexLocation := search.IndexPath(t.Name())
		searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
			return searchIndexManager, nil
		}

		worker, err := ProvideArchivesWorker(
			ctx,
			logger,
			client,
			dbManager,
			postArchivesPublisher,
			searchIndexLocation,
			searchIndexProvider,
		)
		require.NotNil(t, worker)
		require.NoError(t, err)

		assert.Error(t, worker.archiveValidIngredientPreparation(ctx, body))

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher, searchIndexManager)
	})
}
