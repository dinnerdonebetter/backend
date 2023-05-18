package apiclient

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
	"github.com/prixfixeco/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestUsers(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(usersTestSuite))
}

type usersBaseSuite struct {
	suite.Suite

	ctx             context.Context
	exampleUser     *types.User
	exampleUserList *types.QueryFilteredResult[types.User]
}

var _ suite.SetupTestSuite = (*usersBaseSuite)(nil)

func (s *usersBaseSuite) SetupTest() {
	s.ctx = context.Background()
	s.exampleUser = fakes.BuildFakeUser()
	// the hashed passwords is never transmitted over the wire.
	s.exampleUser.HashedPassword = ""
	// the two factor secret is transmitted over the wire only on creation.
	s.exampleUser.TwoFactorSecret = ""
	// the two factor secret validation is never transmitted over the wire.
	s.exampleUser.TwoFactorSecretVerifiedAt = nil

	s.exampleUserList = fakes.BuildFakeUserList()
	for i := 0; i < len(s.exampleUserList.Data); i++ {
		// the hashed passwords is never transmitted over the wire.
		s.exampleUserList.Data[i].HashedPassword = ""
		// the two factor secret is transmitted over the wire only on creation.
		s.exampleUserList.Data[i].TwoFactorSecret = ""
		// the two factor secret validation is never transmitted over the wire.
		s.exampleUserList.Data[i].TwoFactorSecretVerifiedAt = nil
	}
}

type usersTestSuite struct {
	suite.Suite

	usersBaseSuite
}

func (s *usersTestSuite) TestClient_GetUser() {
	const expectedPathFormat = "/api/v1/users/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "", expectedPathFormat, s.exampleUser.ID)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleUser)

		actual, err := c.GetUser(s.ctx, s.exampleUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleUser, actual)
	})

	s.Run("with invalid user ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.GetUser(s.ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.GetUser(s.ctx, s.exampleUser.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.GetUser(s.ctx, s.exampleUser.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *usersTestSuite) TestClient_GetUsers() {
	const expectedPath = "/api/v1/users"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, "limit=20&page=1&sortBy=asc", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleUserList)

		actual, err := c.GetUsers(s.ctx, nil)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleUserList, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.GetUsers(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.GetUsers(s.ctx, nil)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *usersTestSuite) TestClient_SearchForUsersByUsername() {
	const expectedPath = "/api/v1/users/search"
	exampleUsername := s.exampleUser.Username

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodGet, fmt.Sprintf("q=%s", exampleUsername), expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleUserList.Data)

		actual, err := c.SearchForUsersByUsername(s.ctx, exampleUsername)
		assert.NoError(t, err)
		assert.Equal(t, s.exampleUserList.Data, actual)
	})

	s.Run("with empty query", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.SearchForUsersByUsername(s.ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		actual, err := c.SearchForUsersByUsername(s.ctx, exampleUsername)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.SearchForUsersByUsername(s.ctx, exampleUsername)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *usersTestSuite) TestClient_CreateUser() {
	const expectedPath = "/users"

	s.Run("standard", func() {
		t := s.T()

		expected := converters.ConvertUserToUserCreationResponse(s.exampleUser)
		exampleInput := fakes.BuildFakeUserRegistrationInputFromUser(s.exampleUser)
		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c := buildTestClientWithRequestBodyValidation(t, spec, &types.UserRegistrationInput{}, exampleInput, expected)

		actual, err := c.CreateUser(s.ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	s.Run("with nil input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		actual, err := c.CreateUser(s.ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeUserRegistrationInputFromUser(s.exampleUser)
		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CreateUser(s.ctx, exampleInput)

		assert.Nil(t, actual)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeUserRegistrationInputFromUser(s.exampleUser)
		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.CreateUser(s.ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func (s *usersTestSuite) TestClient_ArchiveUser() {
	const expectedPathFormat = "/api/v1/users/%s"

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(true, http.MethodDelete, "", expectedPathFormat, s.exampleUser.ID)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		err := c.ArchiveUser(s.ctx, s.exampleUser.ID)
		assert.NoError(t, err)
	})

	s.Run("with invalid user ID", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.ArchiveUser(s.ctx, "")
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.ArchiveUser(s.ctx, s.exampleUser.ID)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.ArchiveUser(s.ctx, s.exampleUser.ID)
		assert.Error(t, err)
	})
}

func (s *usersTestSuite) TestClient_UploadNewAvatar() {
	const expectedPath = "/api/v1/users/avatar/upload"
	exampleInput := fakes.BuildFakeAvatarUpdateInput()

	s.Run("standard", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithStatusCodeResponse(t, spec, http.StatusOK)

		err := c.UploadNewAvatar(s.ctx, exampleInput)
		assert.NoError(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UploadNewAvatar(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with invalid extension", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UploadNewAvatar(s.ctx, exampleInput)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		err := c.UploadNewAvatar(s.ctx, exampleInput)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		err := c.UploadNewAvatar(s.ctx, exampleInput)
		assert.Error(t, err)
	})
}

func (s *usersTestSuite) TestClient_CheckUserPermissions() {
	const expectedPath = "/api/v1/users/permissions/check"

	s.Run("standard", func() {
		t := s.T()

		exampleResponse := &types.UserPermissionsResponse{Permissions: map[string]bool{"things": true}}

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c, _ := buildTestClientWithJSONResponse(t, spec, exampleResponse)

		perms, err := c.CheckUserPermissions(s.ctx, t.Name())
		assert.NoError(t, err)
		assert.NotNil(t, perms)
		assert.NotEmpty(t, perms)
	})

	s.Run("with nil permissions ", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		perms, err := c.CheckUserPermissions(s.ctx)
		assert.Error(t, err)
		assert.Nil(t, perms)
	})

	s.Run("with invalid request builder", func() {
		t := s.T()

		c := buildTestClientWithInvalidURL(t)

		perms, err := c.CheckUserPermissions(s.ctx, t.Name())
		assert.Error(t, err)
		assert.Nil(t, perms)
	})

	s.Run("with invalid response", func() {
		t := s.T()

		spec := newRequestSpec(false, http.MethodPost, "", expectedPath)
		c := buildTestClientWithInvalidResponse(t, spec)

		perms, err := c.CheckUserPermissions(s.ctx, t.Name())
		assert.Error(t, err)
		assert.Nil(t, perms)
	})
}

func (s *usersTestSuite) TestClient_UpdateUserEmailAddress() {
	const expectedPathFormat = "/api/v1/users/email_address"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeUserEmailAddressUpdateInput()

		spec := newRequestSpec(true, http.MethodPut, "", expectedPathFormat)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleUser)

		err := c.UpdateUserEmailAddress(s.ctx, exampleInput)
		assert.NoError(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateUserEmailAddress(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeUserEmailAddressUpdateInput()
		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateUserEmailAddress(s.ctx, exampleInput)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.GetUser(s.ctx, s.exampleUser.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *usersTestSuite) TestClient_UpdateUserUsername() {
	const expectedPathFormat = "/api/v1/users/username"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeUsernameUpdateInput()

		spec := newRequestSpec(true, http.MethodPut, "", expectedPathFormat)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleUser)

		err := c.UpdateUserUsername(s.ctx, exampleInput)
		assert.NoError(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateUserUsername(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeUsernameUpdateInput()
		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateUserUsername(s.ctx, exampleInput)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.GetUser(s.ctx, s.exampleUser.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}

func (s *usersTestSuite) TestClient_UpdateUserDetails() {
	const expectedPathFormat = "/api/v1/users/details"

	s.Run("standard", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeUserDetailsUpdateInput()

		spec := newRequestSpec(true, http.MethodPut, "", expectedPathFormat)
		c, _ := buildTestClientWithJSONResponse(t, spec, s.exampleUser)

		err := c.UpdateUserDetails(s.ctx, exampleInput)
		assert.NoError(t, err)
	})

	s.Run("with invalid input", func() {
		t := s.T()

		c, _ := buildSimpleTestClient(t)

		err := c.UpdateUserDetails(s.ctx, nil)
		assert.Error(t, err)
	})

	s.Run("with error building request", func() {
		t := s.T()

		exampleInput := fakes.BuildFakeUserDetailsUpdateInput()
		c := buildTestClientWithInvalidURL(t)

		err := c.UpdateUserDetails(s.ctx, exampleInput)
		assert.Error(t, err)
	})

	s.Run("with error executing request", func() {
		t := s.T()

		c, _ := buildTestClientThatWaitsTooLong(t)

		actual, err := c.GetUser(s.ctx, s.exampleUser.ID)
		assert.Nil(t, actual)
		assert.Error(t, err)
	})
}
