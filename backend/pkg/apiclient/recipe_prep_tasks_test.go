package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestRecipePrepTasks(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(recipePrepTasksTestSuite))
}

type recipePrepTasksBaseSuite struct {
	suite.Suite
	ctx                               context.Context
	exampleRecipePrepTask             *types.RecipePrepTask
	exampleRecipePrepTaskResponse     *types.APIResponse[*types.RecipePrepTask]
	exampleRecipePrepTaskListResponse *types.APIResponse[[]*types.RecipePrepTask]
	exampleRecipeID                   string
	exampleRecipePrepTaskList         []*types.RecipePrepTask
}

var _ suite.SetupTestSuite = (*recipePrepTasksBaseSuite)(nil)

func (s *recipePrepTasksBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleRecipeID = fakes.BuildFakeID()
	s.exampleRecipePrepTask = fakes.BuildFakeRecipePrepTask()
	s.exampleRecipePrepTask.BelongsToRecipe = s.exampleRecipeID
	s.exampleRecipePrepTaskResponse = &types.APIResponse[*types.RecipePrepTask]{
		Data: s.exampleRecipePrepTask,
	}

	exampleList := fakes.BuildFakeRecipePrepTaskList()
	s.exampleRecipePrepTaskList = exampleList.Data
	s.exampleRecipePrepTaskListResponse = &types.APIResponse[[]*types.RecipePrepTask]{
		Data:       s.exampleRecipePrepTaskList,
		Pagination: &exampleList.Pagination,
	}
}

type recipePrepTasksTestSuite struct {
	suite.Suite
	recipePrepTasksBaseSuite
}

func (s *recipePrepTasksTestSuite) TestClient_GetRecipePrepTask() {
	const expectedPathFormat = "/api/v1/recipes/%s/prep_tasks/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipePrepTask.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipePrepTaskResponse)
		actual, err := c.GetRecipePrepTask(s.ctx, s.exampleRecipeID, s.exampleRecipePrepTask.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleRecipePrepTask, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipePrepTask(s.ctx, "", s.exampleRecipePrepTask.ID)

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipePrepTask(s.ctx, s.exampleRecipeID, "")

		require.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipePrepTask(s.ctx, s.exampleRecipeID, s.exampleRecipePrepTask.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipePrepTask.ID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipePrepTask(s.ctx, s.exampleRecipeID, s.exampleRecipePrepTask.ID)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipePrepTasksTestSuite) TestClient_GetRecipePrepTasks() {
	const expectedPath = "/api/v1/recipes/%s/prep_tasks"

	filter := (*types.QueryFilter)(nil)

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, s.exampleRecipeID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipePrepTaskListResponse)
		actual, err := c.GetRecipePrepTasks(s.ctx, s.exampleRecipeID, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleRecipePrepTaskList, actual)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		actual, err := c.GetRecipePrepTasks(s.ctx, "", filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetRecipePrepTasks(s.ctx, s.exampleRecipeID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath, s.exampleRecipeID)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetRecipePrepTasks(s.ctx, s.exampleRecipeID, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipePrepTasksTestSuite) TestClient_CreateRecipePrepTask() {
	const expectedPath = "/api/v1/recipes/%s/prep_tasks"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeRecipePrepTaskCreationRequestInput()
		exampleInput.BelongsToRecipe = s.exampleRecipeID

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath, s.exampleRecipeID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipePrepTaskResponse)

		actual, err := c.CreateRecipePrepTask(s.ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleRecipePrepTask, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateRecipePrepTask(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.RecipePrepTaskCreationRequestInput{}

		actual, err := c.CreateRecipePrepTask(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertRecipePrepTaskToRecipePrepTaskCreationRequestInput(s.exampleRecipePrepTask)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateRecipePrepTask(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertRecipePrepTaskToRecipePrepTaskCreationRequestInput(s.exampleRecipePrepTask)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateRecipePrepTask(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *recipePrepTasksTestSuite) TestClient_UpdateRecipePrepTask() {
	const expectedPathFormat = "/api/v1/recipes/%s/prep_tasks/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipePrepTask.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipePrepTaskResponse)

		err := c.UpdateRecipePrepTask(s.ctx, s.exampleRecipePrepTask)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateRecipePrepTask(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateRecipePrepTask(s.ctx, s.exampleRecipePrepTask)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateRecipePrepTask(s.ctx, s.exampleRecipePrepTask)
		assert.Error(t, err)
	})
}

func (s *recipePrepTasksTestSuite) TestClient_ArchiveRecipePrepTask() {
	const expectedPathFormat = "/api/v1/recipes/%s/prep_tasks/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleRecipeID, s.exampleRecipePrepTask.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleRecipePrepTaskResponse)

		err := c.ArchiveRecipePrepTask(s.ctx, s.exampleRecipeID, s.exampleRecipePrepTask.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid recipe ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipePrepTask(s.ctx, "", s.exampleRecipePrepTask.ID)
		assert.Error(t, err)
	})

	s.Run("with invalid recipe step ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveRecipePrepTask(s.ctx, s.exampleRecipeID, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveRecipePrepTask(s.ctx, s.exampleRecipeID, s.exampleRecipePrepTask.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveRecipePrepTask(s.ctx, s.exampleRecipeID, s.exampleRecipePrepTask.ID)
		assert.Error(t, err)
	})
}
