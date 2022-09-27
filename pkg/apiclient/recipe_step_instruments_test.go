package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func TestRecipeStepInstruments(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(recipeStepInstrumentsTestSuite))
}

type recipeStepInstrumentsBaseSuite struct {
	suite.Suite
	ctx                         context.Context
	exampleRecipeStepInstrument *types.RecipeStepInstrument
	exampleRecipeID             string
	exampleRecipeStepID         string
}

var _ suite.SetupTestSuite = (*recipeStepInstrumentsBaseSuite)(nil)

func (s *recipeStepInstrumentsBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleRecipeID = fakes.BuildFakeID()
	s.exampleRecipeStepID = fakes.BuildFakeID()
	s.exampleRecipeStepInstrument = fakes.BuildFakeRecipeStepInstrument()
	s.exampleRecipeStepInstrument.BelongsToRecipeStep = s.exampleRecipeStepID
}

type recipeStepInstrumentsTestSuite struct {
	suite.Suite

	recipeStepInstrumentsBaseSuite
}

func (s *recipeStepInstrumentsTestSuite) TestClient_GetRecipeStepInstrument() {
	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/instruments/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepInstrument.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStepInstrument)
		actual, err := c.GetRecipeStepInstrument(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepInstrument.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleRecipeStepInstrument, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepInstrument(s.ctx, "", s.exampleRecipeStepID, s.exampleRecipeStepInstrument.ID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepInstrument(s.ctx, s.exampleRecipeID, "", s.exampleRecipeStepInstrument.ID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step instrument ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepInstrument(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeStepInstrument(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepInstrument.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepInstrument.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipeStepInstrument(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepInstrument.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeStepInstrumentsTestSuite) TestClient_GetRecipeStepInstruments() {
	const expectedPath = "/api/v1/recipes/%s/steps/%s/instruments"

	s.Run("standard", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		exampleRecipeStepInstrumentList := fakes.BuildFakeRecipeStepInstrumentList()

		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPath, s.exampleRecipeID, s.exampleRecipeStepID)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleRecipeStepInstrumentList)
		actual, err := c.GetRecipeStepInstruments(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepInstrumentList, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepInstruments(s.ctx, "", s.exampleRecipeStepID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipeStepInstruments(s.ctx, s.exampleRecipeID, "", filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipeStepInstruments(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := (*types.QueryFilter)(nil)

		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPath, s.exampleRecipeID, s.exampleRecipeStepID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipeStepInstruments(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeStepInstrumentsTestSuite) TestClient_CreateRecipeStepInstrument() {
	const expectedPath = "/api/v1/recipes/%s/steps/%s/instruments"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepInstrumentCreationRequestInput()
		exampleInput.BelongsToRecipeStep = s.exampleRecipeStepID

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, s.exampleRecipeID, s.exampleRecipeStepID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStepInstrument)

		actual, err := c.CreateRecipeStepInstrument(s.ctx, s.exampleRecipeID, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleRecipeStepInstrument, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepInstrumentCreationRequestInput()
		exampleInput.BelongsToRecipeStep = s.exampleRecipeStepID

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateRecipeStepInstrument(s.ctx, "", exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateRecipeStepInstrument(s.ctx, s.exampleRecipeID, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.RecipeStepInstrumentCreationRequestInput{}

		actual, err := c.CreateRecipeStepInstrument(s.ctx, s.exampleRecipeID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepInstrumentCreationRequestInputFromRecipeStepInstrument(s.exampleRecipeStepInstrument)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateRecipeStepInstrument(s.ctx, s.exampleRecipeID, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipeStepInstrumentCreationRequestInputFromRecipeStepInstrument(s.exampleRecipeStepInstrument)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateRecipeStepInstrument(s.ctx, s.exampleRecipeID, exampleInput)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipeStepInstrumentsTestSuite) TestClient_UpdateRecipeStepInstrument() {
	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/instruments/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepInstrument.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipeStepInstrument)

		err := c.UpdateRecipeStepInstrument(s.ctx, s.exampleRecipeID, s.exampleRecipeStepInstrument)
		assert.NoError(t, err)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateRecipeStepInstrument(s.ctx, "", s.exampleRecipeStepInstrument)
		assert.Error(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateRecipeStepInstrument(s.ctx, s.exampleRecipeID, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateRecipeStepInstrument(s.ctx, s.exampleRecipeID, s.exampleRecipeStepInstrument)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateRecipeStepInstrument(s.ctx, s.exampleRecipeID, s.exampleRecipeStepInstrument)
		assert.Error(t, err)
	})
}

func (s *recipeStepInstrumentsTestSuite) TestClient_ArchiveRecipeStepInstrument() {
	const expectedPathFormat = "/api/v1/recipes/%s/steps/%s/instruments/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepInstrument.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		err := c.ArchiveRecipeStepInstrument(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepInstrument.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeStepInstrument(s.ctx, "", s.exampleRecipeStepID, s.exampleRecipeStepInstrument.ID)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeStepInstrument(s.ctx, s.exampleRecipeID, "", s.exampleRecipeStepInstrument.ID)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step instrument ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipeStepInstrument(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveRecipeStepInstrument(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepInstrument.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveRecipeStepInstrument(s.ctx, s.exampleRecipeID, s.exampleRecipeStepID, s.exampleRecipeStepInstrument.ID)
		assert.Error(t, err)
	})
}