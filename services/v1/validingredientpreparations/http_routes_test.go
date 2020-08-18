package validingredientpreparations

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockencoding "gitlab.com/prixfixe/prixfixe/internal/v1/encoding/mock"
	mockmetrics "gitlab.com/prixfixe/prixfixe/internal/v1/metrics/mock"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"
	mockmodels "gitlab.com/prixfixe/prixfixe/models/v1/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	mocknewsman "gitlab.com/verygoodsoftwarenotvirus/newsman/mock"
)

func TestValidIngredientPreparationsService_ListHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientPreparationList := fakemodels.BuildFakeValidIngredientPreparationList()

		validIngredientPreparationDataManager := &mockmodels.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On("GetValidIngredientPreparations", mock.Anything, mock.AnythingOfType("*models.QueryFilter")).Return(exampleValidIngredientPreparationList, nil)
		s.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidIngredientPreparationList")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager, ed)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		s := buildTestService()

		validIngredientPreparationDataManager := &mockmodels.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On("GetValidIngredientPreparations", mock.Anything, mock.AnythingOfType("*models.QueryFilter")).Return((*models.ValidIngredientPreparationList)(nil), sql.ErrNoRows)
		s.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidIngredientPreparationList")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager, ed)
	})

	T.Run("with error fetching valid ingredient preparations from database", func(t *testing.T) {
		s := buildTestService()

		validIngredientPreparationDataManager := &mockmodels.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On("GetValidIngredientPreparations", mock.Anything, mock.AnythingOfType("*models.QueryFilter")).Return((*models.ValidIngredientPreparationList)(nil), errors.New("blah"))
		s.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientPreparationList := fakemodels.BuildFakeValidIngredientPreparationList()

		validIngredientPreparationDataManager := &mockmodels.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On("GetValidIngredientPreparations", mock.Anything, mock.AnythingOfType("*models.QueryFilter")).Return(exampleValidIngredientPreparationList, nil)
		s.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidIngredientPreparationList")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ListHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager, ed)
	})
}

func TestValidIngredientPreparationsService_CreateHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleInput := fakemodels.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)

		validIngredientPreparationDataManager := &mockmodels.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On("CreateValidIngredientPreparation", mock.Anything, mock.AnythingOfType("*models.ValidIngredientPreparationCreationInput")).Return(exampleValidIngredientPreparation, nil)
		s.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.validIngredientPreparationCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidIngredientPreparation")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), createMiddlewareCtxKey, exampleInput))

		s.CreateHandler(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager, mc, r, ed)
	})

	T.Run("without input attached", func(t *testing.T) {
		s := buildTestService()

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.CreateHandler(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with error creating valid ingredient preparation", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleInput := fakemodels.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)

		validIngredientPreparationDataManager := &mockmodels.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On("CreateValidIngredientPreparation", mock.Anything, mock.AnythingOfType("*models.ValidIngredientPreparationCreationInput")).Return((*models.ValidIngredientPreparation)(nil), errors.New("blah"))
		s.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), createMiddlewareCtxKey, exampleInput))

		s.CreateHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleInput := fakemodels.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)

		validIngredientPreparationDataManager := &mockmodels.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On("CreateValidIngredientPreparation", mock.Anything, mock.AnythingOfType("*models.ValidIngredientPreparationCreationInput")).Return(exampleValidIngredientPreparation, nil)
		s.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		mc := &mockmetrics.UnitCounter{}
		mc.On("Increment", mock.Anything)
		s.validIngredientPreparationCounter = mc

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidIngredientPreparation")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), createMiddlewareCtxKey, exampleInput))

		s.CreateHandler(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager, mc, r, ed)
	})
}

func TestValidIngredientPreparationsService_ExistenceHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		s.validIngredientPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		validIngredientPreparationDataManager := &mockmodels.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On("ValidIngredientPreparationExists", mock.Anything, exampleValidIngredientPreparation.ID).Return(true, nil)
		s.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ExistenceHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})

	T.Run("with no such valid ingredient preparation in database", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		s.validIngredientPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		validIngredientPreparationDataManager := &mockmodels.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On("ValidIngredientPreparationExists", mock.Anything, exampleValidIngredientPreparation.ID).Return(false, sql.ErrNoRows)
		s.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ExistenceHandler(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})

	T.Run("with error fetching valid ingredient preparation from database", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		s.validIngredientPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		validIngredientPreparationDataManager := &mockmodels.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On("ValidIngredientPreparationExists", mock.Anything, exampleValidIngredientPreparation.ID).Return(false, errors.New("blah"))
		s.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ExistenceHandler(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})
}

func TestValidIngredientPreparationsService_ReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		s.validIngredientPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		validIngredientPreparationDataManager := &mockmodels.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On("GetValidIngredientPreparation", mock.Anything, exampleValidIngredientPreparation.ID).Return(exampleValidIngredientPreparation, nil)
		s.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidIngredientPreparation")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager, ed)
	})

	T.Run("with no such valid ingredient preparation in database", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		s.validIngredientPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		validIngredientPreparationDataManager := &mockmodels.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On("GetValidIngredientPreparation", mock.Anything, exampleValidIngredientPreparation.ID).Return((*models.ValidIngredientPreparation)(nil), sql.ErrNoRows)
		s.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})

	T.Run("with error fetching valid ingredient preparation from database", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		s.validIngredientPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		validIngredientPreparationDataManager := &mockmodels.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On("GetValidIngredientPreparation", mock.Anything, exampleValidIngredientPreparation.ID).Return((*models.ValidIngredientPreparation)(nil), errors.New("blah"))
		s.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		s.validIngredientPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		validIngredientPreparationDataManager := &mockmodels.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On("GetValidIngredientPreparation", mock.Anything, exampleValidIngredientPreparation.ID).Return(exampleValidIngredientPreparation, nil)
		s.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidIngredientPreparation")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ReadHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager, ed)
	})
}

func TestValidIngredientPreparationsService_UpdateHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleInput := fakemodels.BuildFakeValidIngredientPreparationUpdateInputFromValidIngredientPreparation(exampleValidIngredientPreparation)

		s.validIngredientPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		validIngredientPreparationDataManager := &mockmodels.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On("GetValidIngredientPreparation", mock.Anything, exampleValidIngredientPreparation.ID).Return(exampleValidIngredientPreparation, nil)
		validIngredientPreparationDataManager.On("UpdateValidIngredientPreparation", mock.Anything, mock.AnythingOfType("*models.ValidIngredientPreparation")).Return(nil)
		s.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidIngredientPreparation")).Return(nil)
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), updateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, r, validIngredientPreparationDataManager, ed)
	})

	T.Run("without update input", func(t *testing.T) {
		s := buildTestService()

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.UpdateHandler(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	T.Run("with no rows fetching valid ingredient preparation", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleInput := fakemodels.BuildFakeValidIngredientPreparationUpdateInputFromValidIngredientPreparation(exampleValidIngredientPreparation)

		s.validIngredientPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		validIngredientPreparationDataManager := &mockmodels.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On("GetValidIngredientPreparation", mock.Anything, exampleValidIngredientPreparation.ID).Return((*models.ValidIngredientPreparation)(nil), sql.ErrNoRows)
		s.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), updateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})

	T.Run("with error fetching valid ingredient preparation", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleInput := fakemodels.BuildFakeValidIngredientPreparationUpdateInputFromValidIngredientPreparation(exampleValidIngredientPreparation)

		s.validIngredientPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		validIngredientPreparationDataManager := &mockmodels.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On("GetValidIngredientPreparation", mock.Anything, exampleValidIngredientPreparation.ID).Return((*models.ValidIngredientPreparation)(nil), errors.New("blah"))
		s.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), updateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})

	T.Run("with error updating valid ingredient preparation", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleInput := fakemodels.BuildFakeValidIngredientPreparationUpdateInputFromValidIngredientPreparation(exampleValidIngredientPreparation)

		s.validIngredientPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		validIngredientPreparationDataManager := &mockmodels.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On("GetValidIngredientPreparation", mock.Anything, exampleValidIngredientPreparation.ID).Return(exampleValidIngredientPreparation, nil)
		validIngredientPreparationDataManager.On("UpdateValidIngredientPreparation", mock.Anything, mock.AnythingOfType("*models.ValidIngredientPreparation")).Return(errors.New("blah"))
		s.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), updateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})

	T.Run("with error encoding response", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleInput := fakemodels.BuildFakeValidIngredientPreparationUpdateInputFromValidIngredientPreparation(exampleValidIngredientPreparation)

		s.validIngredientPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		validIngredientPreparationDataManager := &mockmodels.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On("GetValidIngredientPreparation", mock.Anything, exampleValidIngredientPreparation.ID).Return(exampleValidIngredientPreparation, nil)
		validIngredientPreparationDataManager.On("UpdateValidIngredientPreparation", mock.Anything, mock.AnythingOfType("*models.ValidIngredientPreparation")).Return(nil)
		s.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		ed := &mockencoding.EncoderDecoder{}
		ed.On("EncodeResponse", mock.Anything, mock.AnythingOfType("*models.ValidIngredientPreparation")).Return(errors.New("blah"))
		s.encoderDecoder = ed

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		req = req.WithContext(context.WithValue(req.Context(), updateMiddlewareCtxKey, exampleInput))

		s.UpdateHandler(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		mock.AssertExpectationsForObjects(t, r, validIngredientPreparationDataManager, ed)
	})
}

func TestValidIngredientPreparationsService_ArchiveHandler(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		s.validIngredientPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		validIngredientPreparationDataManager := &mockmodels.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On("ArchiveValidIngredientPreparation", mock.Anything, exampleValidIngredientPreparation.ID).Return(nil)
		s.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		r := &mocknewsman.Reporter{}
		r.On("Report", mock.AnythingOfType("newsman.Event")).Return()
		s.reporter = r

		mc := &mockmetrics.UnitCounter{}
		mc.On("Decrement", mock.Anything).Return()
		s.validIngredientPreparationCounter = mc

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler(res, req)

		assert.Equal(t, http.StatusNoContent, res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager, mc, r)
	})

	T.Run("with no valid ingredient preparation in database", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		s.validIngredientPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		validIngredientPreparationDataManager := &mockmodels.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On("ArchiveValidIngredientPreparation", mock.Anything, exampleValidIngredientPreparation.ID).Return(sql.ErrNoRows)
		s.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		s := buildTestService()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		s.validIngredientPreparationIDFetcher = func(req *http.Request) uint64 {
			return exampleValidIngredientPreparation.ID
		}

		validIngredientPreparationDataManager := &mockmodels.ValidIngredientPreparationDataManager{}
		validIngredientPreparationDataManager.On("ArchiveValidIngredientPreparation", mock.Anything, exampleValidIngredientPreparation.ID).Return(errors.New("blah"))
		s.validIngredientPreparationDataManager = validIngredientPreparationDataManager

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			http.MethodGet,
			"http://todo.verygoodsoftwarenotvirus.ru",
			nil,
		)
		require.NotNil(t, req)
		require.NoError(t, err)

		s.ArchiveHandler(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code)

		mock.AssertExpectationsForObjects(t, validIngredientPreparationDataManager)
	})
}
