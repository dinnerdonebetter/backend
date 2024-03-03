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

func TestUserIngredientPreferences(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(userIngredientPreferencesTestSuite))
}

type userIngredientPreferencesBaseSuite struct {
	suite.Suite
	ctx                                         context.Context
	exampleUserIngredientPreference             *types.UserIngredientPreference
	exampleUserIngredientPreferenceResponse     *types.APIResponse[*types.UserIngredientPreference]
	exampleUserIngredientPreferenceListResponse *types.APIResponse[[]*types.UserIngredientPreference]
	exampleUserIngredientPreferenceList         []*types.UserIngredientPreference
}

var _ suite.SetupTestSuite = (*userIngredientPreferencesBaseSuite)(nil)

func (s *userIngredientPreferencesBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleUserIngredientPreference = fakes.BuildFakeUserIngredientPreference()
	s.exampleUserIngredientPreferenceResponse = &types.APIResponse[*types.UserIngredientPreference]{
		Data: s.exampleUserIngredientPreference,
	}

	exampleList := fakes.BuildFakeUserIngredientPreferenceList()
	s.exampleUserIngredientPreferenceList = exampleList.Data
	s.exampleUserIngredientPreferenceListResponse = &types.APIResponse[[]*types.UserIngredientPreference]{
		Data:       s.exampleUserIngredientPreferenceList,
		Pagination: &exampleList.Pagination,
	}
}

type userIngredientPreferencesTestSuite struct {
	suite.Suite
	userIngredientPreferencesBaseSuite
}

func (s *userIngredientPreferencesTestSuite) TestClient_GetUserIngredientPreferences() {
	const expectedPath = "/api/v1/user_ingredient_preferences"

	filter := (*types.QueryFilter)(nil)

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleUserIngredientPreferenceListResponse)
		actual, err := c.GetUserIngredientPreferences(s.ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleUserIngredientPreferenceList, actual.Data)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetUserIngredientPreferences(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)
		actual, err := c.GetUserIngredientPreferences(s.ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *userIngredientPreferencesTestSuite) TestClient_CreateUserIngredientPreference() {
	const expectedPath = "/api/v1/user_ingredient_preferences"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeUserIngredientPreferenceCreationRequestInput()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleUserIngredientPreferenceListResponse)

		actual, err := c.CreateUserIngredientPreference(s.ctx, exampleInput)
		require.NotEmpty(t, actual)
		assert.NoError(t, err)

		assert.Equal(t, s.exampleUserIngredientPreferenceList, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateUserIngredientPreference(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)
		exampleInput := &types.UserIngredientPreferenceCreationRequestInput{}

		actual, err := c.CreateUserIngredientPreference(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := converters.ConvertUserIngredientPreferenceToUserIngredientPreferenceCreationRequestInput(s.exampleUserIngredientPreference)

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.CreateUserIngredientPreference(s.ctx, exampleInput)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := converters.ConvertUserIngredientPreferenceToUserIngredientPreferenceCreationRequestInput(s.exampleUserIngredientPreference)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateUserIngredientPreference(s.ctx, exampleInput)
		assert.Empty(t, actual)
		assert.Error(t, err)
	})
}

func (s *userIngredientPreferencesTestSuite) TestClient_UpdateUserIngredientPreference() {
	const expectedPathFormat = "/api/v1/user_ingredient_preferences/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPut, "", expectedPathFormat, s.exampleUserIngredientPreference.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleUserIngredientPreferenceResponse)

		err := c.UpdateUserIngredientPreference(s.ctx, s.exampleUserIngredientPreference)
		assert.NoError(t, err)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateUserIngredientPreference(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateUserIngredientPreference(s.ctx, s.exampleUserIngredientPreference)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UpdateUserIngredientPreference(s.ctx, s.exampleUserIngredientPreference)
		assert.Error(t, err)
	})
}

func (s *userIngredientPreferencesTestSuite) TestClient_ArchiveUserIngredientPreference() {
	const expectedPathFormat = "/api/v1/user_ingredient_preferences/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleUserIngredientPreference.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleUserIngredientPreferenceResponse)

		err := c.ArchiveUserIngredientPreference(s.ctx, s.exampleUserIngredientPreference.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid user ingredient preference ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveUserIngredientPreference(s.ctx, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveUserIngredientPreference(s.ctx, s.exampleUserIngredientPreference.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveUserIngredientPreference(s.ctx, s.exampleUserIngredientPreference.ID)
		assert.Error(t, err)
	})
}
