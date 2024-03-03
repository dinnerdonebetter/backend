package apiclient

import (
	"context"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestHouseholdInvitations(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(householdInvitationsTestSuite))
}

type householdInvitationsTestSuite struct {
	suite.Suite
	ctx                                    context.Context
	exampleHousehold                       *types.Household
	exampleHouseholdInvitation             *types.HouseholdInvitation
	exampleHouseholdInvitationResponse     *types.APIResponse[*types.HouseholdInvitation]
	exampleUser                            *types.User
	exampleHouseholdInvitationListResponse *types.APIResponse[[]*types.HouseholdInvitation]
	exampleHouseholdInvitationList         []*types.HouseholdInvitation
}

var _ suite.SetupTestSuite = (*householdInvitationsTestSuite)(nil)

func (s *householdInvitationsTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleUser = fakes.BuildFakeUser()
	s.exampleHousehold = fakes.BuildFakeHousehold()
	s.exampleHousehold.BelongsToUser = s.exampleUser.ID
	s.exampleHouseholdInvitation = fakes.BuildFakeHouseholdInvitation()
	s.exampleHouseholdInvitation.FromUser = *s.exampleUser
	s.exampleHouseholdInvitation.ToUser = func(s string) *string { return &s }(fakes.BuildFakeUser().ID)
	s.exampleHouseholdInvitationResponse = &types.APIResponse[*types.HouseholdInvitation]{
		Data: s.exampleHouseholdInvitation,
	}
	exampleList := fakes.BuildFakeHouseholdInvitationList()
	s.exampleHouseholdInvitationList = exampleList.Data
	s.exampleHouseholdInvitationListResponse = &types.APIResponse[[]*types.HouseholdInvitation]{
		Data:       s.exampleHouseholdInvitationList,
		Pagination: &exampleList.Pagination,
	}
}

func (s *householdInvitationsTestSuite) TestClient_GetHouseholdInvitation() {
	const expectedPath = "/api/v1/households/%s/invitations/%s"

	s.Run("standard", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(true, http.MethodGet, "", expectedPath, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdInvitationResponse)

		actual, err := c.GetHouseholdInvitation(s.ctx, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})

	s.Run("with invalid household ID", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(true, http.MethodGet, "", expectedPath, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdInvitation)

		actual, err := c.GetHouseholdInvitation(s.ctx, "", s.exampleHouseholdInvitation.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	s.Run("with invalid invitation ID", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(true, http.MethodGet, "", expectedPath, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdInvitation)

		actual, err := c.GetHouseholdInvitation(s.ctx, s.exampleHousehold.ID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.GetHouseholdInvitation(s.ctx, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.GetHouseholdInvitation(s.ctx, s.exampleHousehold.ID, s.exampleHouseholdInvitation.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func (s *householdInvitationsTestSuite) TestClient_GetPendingHouseholdInvitationsForUser() {
	const expectedPath = "/api/v1/household_invitations/received"

	s.Run("standard", func() {
		t := s.T()

		filter := types.DefaultQueryFilter()
		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdInvitationListResponse)

		actual, err := c.GetPendingHouseholdInvitationsForUser(s.ctx, filter)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := types.DefaultQueryFilter()

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.GetPendingHouseholdInvitationsForUser(s.ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := types.DefaultQueryFilter()

		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.GetPendingHouseholdInvitationsForUser(s.ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func (s *householdInvitationsTestSuite) TestClient_GetPendingHouseholdInvitationsFromUser() {
	const expectedPath = "/api/v1/household_invitations/sent"

	s.Run("standard", func() {
		t := s.T()

		filter := types.DefaultQueryFilter()
		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(true, http.MethodGet, "limit=50&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdInvitationListResponse)

		actual, err := c.GetPendingHouseholdInvitationsFromUser(s.ctx, filter)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		filter := types.DefaultQueryFilter()

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.GetPendingHouseholdInvitationsFromUser(s.ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		filter := types.DefaultQueryFilter()

		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.GetPendingHouseholdInvitationsFromUser(s.ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func (s *householdInvitationsTestSuite) TestClient_AcceptHouseholdInvitation() {
	const expectedPath = "/api/v1/household_invitations/%s/accept"

	s.Run("standard", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(false, http.MethodPut, "", expectedPath, s.exampleHouseholdInvitation.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdInvitationResponse)

		assert.NoError(t, c.AcceptHouseholdInvitation(s.ctx, s.exampleHouseholdInvitation.ID, s.exampleHouseholdInvitation.Token, t.Name()))
	})

	s.Run("with invalid token", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(false, http.MethodPut, "", expectedPath, s.exampleHouseholdInvitation.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdInvitationResponse)

		assert.Error(t, c.AcceptHouseholdInvitation(s.ctx, s.exampleHouseholdInvitation.ID, "", t.Name()))
	})

	s.Run("with invalid household invitation ID", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(false, http.MethodPut, "", expectedPath, s.exampleHouseholdInvitation.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdInvitationResponse)

		assert.Error(t, c.AcceptHouseholdInvitation(s.ctx, "", s.exampleHouseholdInvitation.Token, t.Name()))
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		assert.Error(t, c.AcceptHouseholdInvitation(s.ctx, s.exampleHouseholdInvitation.ID, s.exampleHouseholdInvitation.Token, t.Name()))
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		assert.Error(t, c.AcceptHouseholdInvitation(s.ctx, s.exampleHouseholdInvitation.ID, s.exampleHouseholdInvitation.Token, t.Name()))
	})
}

func (s *householdInvitationsTestSuite) TestClient_CancelHouseholdInvitation() {
	const expectedPath = "/api/v1/household_invitations/%s/cancel"

	s.Run("standard", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(false, http.MethodPut, "", expectedPath, s.exampleHouseholdInvitation.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdInvitationResponse)

		assert.NoError(t, c.CancelHouseholdInvitation(s.ctx, s.exampleHouseholdInvitation.ID, s.exampleHouseholdInvitation.Token, t.Name()))
	})

	s.Run("with invalid household ID", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(false, http.MethodPut, "", expectedPath, s.exampleHouseholdInvitation.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdInvitationResponse)

		assert.Error(t, c.CancelHouseholdInvitation(s.ctx, "", s.exampleHouseholdInvitation.ID, t.Name()))
	})

	s.Run("with invalid household invitation ID", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(false, http.MethodPut, "", expectedPath, s.exampleHouseholdInvitation.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdInvitationResponse)

		assert.Error(t, c.CancelHouseholdInvitation(s.ctx, s.exampleHouseholdInvitation.ID, "", t.Name()))
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		assert.Error(t, c.CancelHouseholdInvitation(s.ctx, s.exampleHouseholdInvitation.ID, s.exampleHouseholdInvitation.Token, t.Name()))
	})

	s.Run("with error executing request", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		c, _ := buildTestClientThatWaitsTooLong(t)

		assert.Error(t, c.CancelHouseholdInvitation(s.ctx, s.exampleHouseholdInvitation.ID, s.exampleHouseholdInvitation.Token, t.Name()))
	})
}

func (s *householdInvitationsTestSuite) TestClient_RejectHouseholdInvitation() {
	const expectedPath = "/api/v1/household_invitations/%s/reject"

	s.Run("standard", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(false, http.MethodPut, "", expectedPath, s.exampleHouseholdInvitation.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdInvitationResponse)

		assert.NoError(t, c.RejectHouseholdInvitation(s.ctx, s.exampleHouseholdInvitation.ID, s.exampleHouseholdInvitation.Token, t.Name()))
	})

	s.Run("with invalid household ID", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(false, http.MethodPut, "", expectedPath, s.exampleHouseholdInvitation.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdInvitationResponse)

		assert.Error(t, c.RejectHouseholdInvitation(s.ctx, "", s.exampleHouseholdInvitation.ID, t.Name()))
	})

	s.Run("with invalid household invitation ID", func() {
		t := s.T()

		s.exampleHousehold.BelongsToUser = ""

		spec := newRequestSpec(false, http.MethodPut, "", expectedPath, s.exampleHouseholdInvitation.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleHouseholdInvitationResponse)

		assert.Error(t, c.RejectHouseholdInvitation(s.ctx, s.exampleHouseholdInvitation.ID, "", t.Name()))
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		assert.Error(t, c.RejectHouseholdInvitation(s.ctx, s.exampleHouseholdInvitation.ID, s.exampleHouseholdInvitation.Token, t.Name()))
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		assert.Error(t, c.RejectHouseholdInvitation(s.ctx, s.exampleHouseholdInvitation.ID, s.exampleHouseholdInvitation.Token, t.Name()))
	})
}
