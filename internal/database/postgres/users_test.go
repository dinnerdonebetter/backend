package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/pkg/random/mock"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	testutils "github.com/dinnerdonebetter/backend/tests/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func buildMockRowsFromUsers(includeCounts bool, filteredCount uint64, users ...*types.User) *sqlmock.Rows {
	columns := usersTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, user := range users {
		rowValues := []driver.Value{
			user.ID,
			user.FirstName,
			user.LastName,
			user.Username,
			user.EmailAddress,
			user.EmailAddressVerifiedAt,
			user.AvatarSrc,
			user.HashedPassword,
			user.RequiresPasswordChange,
			user.PasswordLastChangedAt,
			user.TwoFactorSecret,
			user.TwoFactorSecretVerifiedAt,
			user.ServiceRole,
			user.AccountStatus,
			user.AccountStatusExplanation,
			user.Birthday,
			user.CreatedAt,
			user.LastUpdatedAt,
			user.ArchivedAt,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(users))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanUsers(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On(
			"Next",
		).Return(false)
		mockRows.On(
			"Err",
		).Return(errors.New("blah"))

		_, _, _, err := q.scanUsers(ctx, mockRows, false)
		assert.Error(t, err)
	})

	T.Run("logs row closing errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On(
			"Next",
		).Return(false)
		mockRows.On(
			"Err",
		).Return(nil)
		mockRows.On(
			"Close",
		).Return(errors.New("blah"))

		_, _, _, err := q.scanUsers(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_UserHasStatus(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		c, db := buildTestClient(t)
		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleStatus := string(types.GoodStandingUserAccountStatus)

		args := []any{exampleUserID, exampleStatus}

		db.ExpectQuery(formatQueryForSQLMock(userHasStatusQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.UserHasStatus(ctx, exampleUserID, exampleStatus)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		c, _ := buildTestClient(t)
		ctx := context.Background()
		exampleStatus := string(types.GoodStandingUserAccountStatus)

		actual, err := c.UserHasStatus(ctx, "", exampleStatus)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with empty statuses list", func(t *testing.T) {
		t.Parallel()

		c, _ := buildTestClient(t)
		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()

		actual, err := c.UserHasStatus(ctx, exampleUserID)
		assert.NoError(t, err)
		assert.True(t, actual)
	})

	T.Run("with error performing query", func(t *testing.T) {
		t.Parallel()

		c, db := buildTestClient(t)
		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleStatus := string(types.GoodStandingUserAccountStatus)

		args := []any{exampleUserID, exampleStatus}

		db.ExpectQuery(formatQueryForSQLMock(userHasStatusQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.UserHasStatus(ctx, exampleUserID, exampleStatus)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_getUser(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{exampleUser.ID}

		db.ExpectQuery(formatQueryForSQLMock(getUserWithVerified2FAQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromUsers(false, 0, exampleUser))

		actual, err := c.getUser(ctx, exampleUser.ID, true)
		assert.NoError(t, err)
		assert.Equal(t, exampleUser, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.getUser(ctx, "", true)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("without verified two factor secret", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{exampleUser.ID}

		db.ExpectQuery(formatQueryForSQLMock(getUserByIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromUsers(false, 0, exampleUser))

		actual, err := c.getUser(ctx, exampleUser.ID, false)
		assert.NoError(t, err)
		assert.Equal(t, exampleUser, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{exampleUser.ID}

		db.ExpectQuery(formatQueryForSQLMock(getUserWithVerified2FAQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.getUser(ctx, exampleUser.ID, true)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetUser(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{exampleUser.ID}

		db.ExpectQuery(formatQueryForSQLMock(getUserByIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromUsers(false, 0, exampleUser))

		actual, err := c.GetUser(ctx, exampleUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleUser, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetUser(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{exampleUser.ID}

		db.ExpectQuery(formatQueryForSQLMock(getUserByIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetUser(ctx, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetUserWithUnverifiedTwoFactorSecret(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{exampleUser.ID}

		db.ExpectQuery(formatQueryForSQLMock(getUserByIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromUsers(false, 0, exampleUser))

		actual, err := c.GetUserWithUnverifiedTwoFactorSecret(ctx, exampleUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleUser, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetUserWithUnverifiedTwoFactorSecret(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetUserByEmail(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{exampleUser.EmailAddress}
		db.ExpectQuery(formatQueryForSQLMock(getUserByEmailQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromUsers(false, 0, exampleUser))

		actual, err := c.GetUserByEmail(ctx, exampleUser.EmailAddress)
		assert.NoError(t, err)
		assert.Equal(t, exampleUser, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid email", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		actual, err := c.GetUserByEmail(ctx, "")
		assert.Error(t, err)
		assert.Empty(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{exampleUser.EmailAddress}
		db.ExpectQuery(formatQueryForSQLMock(getUserByEmailQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetUserByEmail(ctx, exampleUser.EmailAddress)
		assert.Error(t, err)
		assert.Empty(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetUserByUsername(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{exampleUser.Username}

		db.ExpectQuery(formatQueryForSQLMock(getUserByUsernameQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromUsers(false, 0, exampleUser))

		actual, err := c.GetUserByUsername(ctx, exampleUser.Username)
		assert.NoError(t, err)
		assert.Equal(t, exampleUser, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid username", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetUserByUsername(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("respects sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{exampleUser.Username}

		db.ExpectQuery(formatQueryForSQLMock(getUserByUsernameQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.GetUserByUsername(ctx, exampleUser.Username)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{exampleUser.Username}

		db.ExpectQuery(formatQueryForSQLMock(getUserByUsernameQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetUserByUsername(ctx, exampleUser.Username)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetAdminUserByUsername(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{exampleUser.Username}

		db.ExpectQuery(formatQueryForSQLMock(getAdminUserByUsernameQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromUsers(false, 0, exampleUser))

		actual, err := c.GetAdminUserByUsername(ctx, exampleUser.Username)
		assert.NoError(t, err)
		assert.Equal(t, exampleUser, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with empty username", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		actual, err := c.GetAdminUserByUsername(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error fetching data", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{exampleUser.Username}

		db.ExpectQuery(formatQueryForSQLMock(getAdminUserByUsernameQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetAdminUserByUsername(ctx, exampleUser.Username)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with missing user", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{exampleUser.Username}

		db.ExpectQuery(formatQueryForSQLMock(getAdminUserByUsernameQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.GetAdminUserByUsername(ctx, exampleUser.Username)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_SearchForUsersByUsername(T *testing.T) {
	T.Parallel()

	exampleUsername := fakes.BuildFakeUser().Username

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUserList := fakes.BuildFakeUserList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			wrapQueryForILIKE(exampleUsername),
		}

		db.ExpectQuery(formatQueryForSQLMock(searchForUserByUsernameQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromUsers(false, 0, exampleUserList.Data...))

		actual, err := c.SearchForUsersByUsername(ctx, exampleUsername)
		assert.NoError(t, err)
		assert.Equal(t, exampleUserList.Data, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid username", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.SearchForUsersByUsername(ctx, "")
		assert.Error(t, err)
		assert.NotNil(t, actual)
		assert.Empty(t, actual)
	})

	T.Run("respects sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			wrapQueryForILIKE(exampleUsername),
		}

		db.ExpectQuery(formatQueryForSQLMock(searchForUserByUsernameQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.SearchForUsersByUsername(ctx, exampleUsername)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, sql.ErrNoRows))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			wrapQueryForILIKE(exampleUsername),
		}

		db.ExpectQuery(formatQueryForSQLMock(searchForUserByUsernameQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.SearchForUsersByUsername(ctx, exampleUsername)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			wrapQueryForILIKE(exampleUsername),
		}

		db.ExpectQuery(formatQueryForSQLMock(searchForUserByUsernameQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.SearchForUsersByUsername(ctx, exampleUsername)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetUsers(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUserList := fakes.BuildFakeUserList()
		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "users", nil, nil, nil, "", usersTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromUsers(true, exampleUserList.FilteredCount, exampleUserList.Data...))

		actual, err := c.GetUsers(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleUserList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		exampleUserList := fakes.BuildFakeUserList()
		exampleUserList.Limit, exampleUserList.Page = 0, 0
		filter := (*types.QueryFilter)(nil)

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "users", nil, nil, nil, "", usersTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromUsers(true, exampleUserList.FilteredCount, exampleUserList.Data...))

		actual, err := c.GetUsers(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleUserList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "users", nil, nil, nil, "", usersTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetUsers(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "users", nil, nil, nil, "", usersTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetUsers(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateUser(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleCreationTime := fakes.BuildFakeTime()

		exampleUser := fakes.BuildFakeUser()
		exampleUser.TwoFactorSecretVerifiedAt = nil
		exampleUser.CreatedAt = exampleCreationTime
		exampleUser.AccountStatus = ""
		exampleUserCreationInput := converters.ConvertUserToUserDatabaseCreationInput(exampleUser)

		exampleHousehold := fakes.BuildFakeHouseholdForUser(exampleUser)
		exampleHousehold.CreatedAt = exampleCreationTime

		ctx := context.Background()
		c, db := buildTestClient(t)

		c.timeFunc = func() time.Time {
			return exampleCreationTime
		}

		db.ExpectBegin()

		exampleToken := fakes.BuildFakeID()
		secretGenerator := &mockrandom.Generator{}
		secretGenerator.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			32,
		).Return(exampleToken, nil)
		c.secretGenerator = secretGenerator

		userCreationArgs := []any{
			exampleUserCreationInput.ID,
			exampleUserCreationInput.FirstName,
			exampleUserCreationInput.LastName,
			exampleUserCreationInput.Username,
			exampleUserCreationInput.EmailAddress,
			exampleUserCreationInput.HashedPassword,
			exampleUserCreationInput.TwoFactorSecret,
			exampleUserCreationInput.AvatarSrc,
			types.UnverifiedHouseholdStatus,
			exampleUserCreationInput.Birthday,
			authorization.ServiceUserRole.String(),
			exampleToken,
		}

		db.ExpectExec(formatQueryForSQLMock(userCreationQuery)).
			WithArgs(interfaceToDriverValue(userCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		// create household for created user
		householdCreationInput := types.HouseholdCreationInputForNewUser(exampleUser)
		householdCreationInput.ID = exampleHousehold.ID
		householdCreationArgs := []any{
			&idMatcher{},
			householdCreationInput.Name,
			types.UnpaidHouseholdBillingStatus,
			householdCreationInput.ContactPhone,
			householdCreationInput.AddressLine1,
			householdCreationInput.AddressLine2,
			householdCreationInput.City,
			householdCreationInput.State,
			householdCreationInput.ZipCode,
			householdCreationInput.Country,
			householdCreationInput.Latitude,
			householdCreationInput.Longitude,
			&idMatcher{},
		}

		db.ExpectExec(formatQueryForSQLMock(householdCreationQuery)).
			WithArgs(interfaceToDriverValue(householdCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		// create household user membership for created user
		createHouseholdMembershipForNewUserArgs := []any{
			&idMatcher{},
			&idMatcher{},
			&idMatcher{},
			true,
			authorization.HouseholdAdminRole.String(),
		}

		db.ExpectExec(formatQueryForSQLMock(createHouseholdMembershipForNewUserQuery)).
			WithArgs(interfaceToDriverValue(createHouseholdMembershipForNewUserArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		// create household user membership for created user
		attachInvitationsToUsersArgs := []any{
			&idMatcher{},
			exampleUser.EmailAddress,
		}

		db.ExpectExec(formatQueryForSQLMock(attachInvitationsToUserIDQuery)).
			WithArgs(interfaceToDriverValue(attachInvitationsToUsersArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		db.ExpectCommit()

		actual, err := c.CreateUser(ctx, exampleUserCreationInput)
		assert.NoError(t, err)
		require.NotNil(t, actual)
		actual.ID = exampleUser.ID
		assert.Equal(t, exampleUser, actual)

		mock.AssertExpectationsForObjects(t, db, secretGenerator)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateUser(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error beginning transaction", func(t *testing.T) {
		t.Parallel()

		exampleCreationTime := fakes.BuildFakeTime()

		exampleUser := fakes.BuildFakeUser()
		exampleUser.TwoFactorSecretVerifiedAt = nil
		exampleUser.CreatedAt = exampleCreationTime
		exampleUser.AccountStatus = ""
		exampleUserCreationInput := converters.ConvertUserToUserDatabaseCreationInput(exampleUser)

		exampleHousehold := fakes.BuildFakeHouseholdForUser(exampleUser)
		exampleHousehold.CreatedAt = exampleCreationTime

		ctx := context.Background()
		c, db := buildTestClient(t)

		c.timeFunc = func() time.Time {
			return exampleCreationTime
		}

		exampleToken := fakes.BuildFakeID()
		secretGenerator := &mockrandom.Generator{}
		secretGenerator.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			32,
		).Return(exampleToken, nil)
		c.secretGenerator = secretGenerator

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		actual, err := c.CreateUser(ctx, exampleUserCreationInput)
		assert.Error(t, err)
		require.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, secretGenerator)
	})

	T.Run("with error executing user creation query", func(t *testing.T) {
		t.Parallel()

		exampleCreationTime := fakes.BuildFakeTime()

		exampleUser := fakes.BuildFakeUser()
		exampleUser.TwoFactorSecretVerifiedAt = nil
		exampleUser.CreatedAt = exampleCreationTime
		exampleUser.AccountStatus = ""
		exampleUserCreationInput := converters.ConvertUserToUserDatabaseCreationInput(exampleUser)

		exampleHousehold := fakes.BuildFakeHouseholdForUser(exampleUser)
		exampleHousehold.CreatedAt = exampleCreationTime

		ctx := context.Background()
		c, db := buildTestClient(t)

		c.timeFunc = func() time.Time {
			return exampleCreationTime
		}

		exampleToken := fakes.BuildFakeID()
		secretGenerator := &mockrandom.Generator{}
		secretGenerator.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			32,
		).Return(exampleToken, nil)
		c.secretGenerator = secretGenerator

		db.ExpectBegin()

		userCreationArgs := []any{
			exampleUserCreationInput.ID,
			exampleUserCreationInput.FirstName,
			exampleUserCreationInput.LastName,
			exampleUserCreationInput.Username,
			exampleUserCreationInput.EmailAddress,
			exampleUserCreationInput.HashedPassword,
			exampleUserCreationInput.TwoFactorSecret,
			exampleUserCreationInput.AvatarSrc,
			types.UnverifiedHouseholdStatus,
			exampleUserCreationInput.Birthday,
			authorization.ServiceUserRole.String(),
			exampleToken,
		}

		db.ExpectExec(formatQueryForSQLMock(userCreationQuery)).
			WithArgs(interfaceToDriverValue(userCreationArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		actual, err := c.CreateUser(ctx, exampleUserCreationInput)
		assert.Error(t, err)
		require.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, secretGenerator)
	})

	T.Run("with already existent user", func(t *testing.T) {
		t.Parallel()

		exampleCreationTime := fakes.BuildFakeTime()

		exampleUser := fakes.BuildFakeUser()
		exampleUser.TwoFactorSecretVerifiedAt = nil
		exampleUser.CreatedAt = exampleCreationTime
		exampleUser.AccountStatus = ""
		exampleUserCreationInput := converters.ConvertUserToUserDatabaseCreationInput(exampleUser)

		exampleHousehold := fakes.BuildFakeHouseholdForUser(exampleUser)
		exampleHousehold.CreatedAt = exampleCreationTime

		ctx := context.Background()
		c, db := buildTestClient(t)

		c.timeFunc = func() time.Time {
			return exampleCreationTime
		}

		exampleToken := fakes.BuildFakeID()
		secretGenerator := &mockrandom.Generator{}
		secretGenerator.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			32,
		).Return(exampleToken, nil)
		c.secretGenerator = secretGenerator

		db.ExpectBegin()

		userCreationArgs := []any{
			exampleUserCreationInput.ID,
			exampleUserCreationInput.FirstName,
			exampleUserCreationInput.LastName,
			exampleUserCreationInput.Username,
			exampleUserCreationInput.EmailAddress,
			exampleUserCreationInput.HashedPassword,
			exampleUserCreationInput.TwoFactorSecret,
			exampleUserCreationInput.AvatarSrc,
			types.UnverifiedHouseholdStatus,
			exampleUserCreationInput.Birthday,
			authorization.ServiceUserRole.String(),
			exampleToken,
		}

		db.ExpectExec(formatQueryForSQLMock(userCreationQuery)).
			WithArgs(interfaceToDriverValue(userCreationArgs)...).
			WillReturnError(&pq.Error{Code: postgresDuplicateEntryErrorCode})

		db.ExpectRollback()

		actual, err := c.CreateUser(ctx, exampleUserCreationInput)
		assert.Error(t, err)
		require.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, secretGenerator)
	})

	T.Run("with destination household", func(t *testing.T) {
		t.Parallel()

		exampleCreationTime := fakes.BuildFakeTime()

		exampleUser := fakes.BuildFakeUser()
		exampleUser.TwoFactorSecretVerifiedAt = nil
		exampleUser.CreatedAt = exampleCreationTime
		exampleUser.AccountStatus = ""
		exampleUserCreationInput := converters.ConvertUserToUserDatabaseCreationInput(exampleUser)

		exampleHousehold := fakes.BuildFakeHouseholdForUser(exampleUser)
		exampleHousehold.CreatedAt = exampleCreationTime

		exampleUserCreationInput.DestinationHouseholdID = exampleHousehold.ID
		exampleUserCreationInput.InvitationToken = t.Name()

		ctx := context.Background()
		c, db := buildTestClient(t)

		c.timeFunc = func() time.Time {
			return exampleCreationTime
		}

		exampleToken := fakes.BuildFakeID()
		secretGenerator := &mockrandom.Generator{}
		secretGenerator.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			32,
		).Return(exampleToken, nil)
		c.secretGenerator = secretGenerator

		db.ExpectBegin()

		userCreationArgs := []any{
			exampleUserCreationInput.ID,
			exampleUserCreationInput.FirstName,
			exampleUserCreationInput.LastName,
			exampleUserCreationInput.Username,
			exampleUserCreationInput.EmailAddress,
			exampleUserCreationInput.HashedPassword,
			exampleUserCreationInput.TwoFactorSecret,
			exampleUserCreationInput.AvatarSrc,
			types.UnverifiedHouseholdStatus,
			exampleUserCreationInput.Birthday,
			authorization.ServiceUserRole.String(),
			exampleToken,
		}

		db.ExpectExec(formatQueryForSQLMock(userCreationQuery)).
			WithArgs(interfaceToDriverValue(userCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		// create household for created user
		householdCreationInput := types.HouseholdCreationInputForNewUser(exampleUser)
		householdCreationInput.ID = exampleHousehold.ID
		householdCreationArgs := []any{
			&idMatcher{},
			householdCreationInput.Name,
			types.UnpaidHouseholdBillingStatus,
			householdCreationInput.ContactPhone,
			householdCreationInput.AddressLine1,
			householdCreationInput.AddressLine2,
			householdCreationInput.City,
			householdCreationInput.State,
			householdCreationInput.ZipCode,
			householdCreationInput.Country,
			householdCreationInput.Latitude,
			householdCreationInput.Longitude,
			&idMatcher{},
		}

		db.ExpectExec(formatQueryForSQLMock(householdCreationQuery)).
			WithArgs(interfaceToDriverValue(householdCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		// create household user membership for created user
		createHouseholdMembershipForNewUserArgs := []any{
			&idMatcher{},
			&idMatcher{},
			&idMatcher{},
			false,
			authorization.HouseholdAdminRole.String(),
		}

		db.ExpectExec(formatQueryForSQLMock(createHouseholdMembershipForNewUserQuery)).
			WithArgs(interfaceToDriverValue(createHouseholdMembershipForNewUserArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		args := []any{
			exampleUserCreationInput.EmailAddress,
			exampleUserCreationInput.InvitationToken,
		}

		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdInvitationByEmailAndTokenQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholdInvitations(false, 0, exampleHouseholdInvitation))

		// create household user membership for created user
		createHouseholdMembershipForNewUserArgs = []any{
			&idMatcher{},
			&idMatcher{},
			&idMatcher{},
			true,
			authorization.HouseholdMemberRole.String(),
		}

		db.ExpectExec(formatQueryForSQLMock(createHouseholdMembershipForNewUserQuery)).
			WithArgs(interfaceToDriverValue(createHouseholdMembershipForNewUserArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		invitationStatusUpdateArgs := []any{
			types.AcceptedHouseholdInvitationStatus,
			"",
			exampleHouseholdInvitation.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(setInvitationStatusQuery)).
			WithArgs(interfaceToDriverValue(invitationStatusUpdateArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		// create household user membership for created user
		attachInvitationsToUsersArgs := []any{
			&idMatcher{},
			exampleUser.EmailAddress,
		}

		db.ExpectExec(formatQueryForSQLMock(attachInvitationsToUserIDQuery)).
			WithArgs(interfaceToDriverValue(attachInvitationsToUsersArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		db.ExpectCommit()

		actual, err := c.CreateUser(ctx, exampleUserCreationInput)
		assert.NoError(t, err)
		require.NotNil(t, actual)
		actual.ID = exampleUser.ID
		assert.Equal(t, exampleUser, actual)

		mock.AssertExpectationsForObjects(t, db, secretGenerator)
	})

	T.Run("with destination household and error fetching pre-existing invitation", func(t *testing.T) {
		t.Parallel()

		exampleCreationTime := fakes.BuildFakeTime()

		exampleUser := fakes.BuildFakeUser()
		exampleUser.TwoFactorSecretVerifiedAt = nil
		exampleUser.CreatedAt = exampleCreationTime
		exampleUser.AccountStatus = ""
		exampleUserCreationInput := converters.ConvertUserToUserDatabaseCreationInput(exampleUser)

		exampleHousehold := fakes.BuildFakeHouseholdForUser(exampleUser)
		exampleHousehold.CreatedAt = exampleCreationTime

		exampleUserCreationInput.DestinationHouseholdID = exampleHousehold.ID
		exampleUserCreationInput.InvitationToken = t.Name()

		ctx := context.Background()
		c, db := buildTestClient(t)

		c.timeFunc = func() time.Time {
			return exampleCreationTime
		}

		exampleToken := fakes.BuildFakeID()
		secretGenerator := &mockrandom.Generator{}
		secretGenerator.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			32,
		).Return(exampleToken, nil)
		c.secretGenerator = secretGenerator

		db.ExpectBegin()

		userCreationArgs := []any{
			exampleUserCreationInput.ID,
			exampleUserCreationInput.FirstName,
			exampleUserCreationInput.LastName,
			exampleUserCreationInput.Username,
			exampleUserCreationInput.EmailAddress,
			exampleUserCreationInput.HashedPassword,
			exampleUserCreationInput.TwoFactorSecret,
			exampleUserCreationInput.AvatarSrc,
			types.UnverifiedHouseholdStatus,
			exampleUserCreationInput.Birthday,
			authorization.ServiceUserRole.String(),
			exampleToken,
		}

		db.ExpectExec(formatQueryForSQLMock(userCreationQuery)).
			WithArgs(interfaceToDriverValue(userCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		// create household for created user
		householdCreationInput := types.HouseholdCreationInputForNewUser(exampleUser)
		householdCreationInput.ID = exampleHousehold.ID
		householdCreationArgs := []any{
			&idMatcher{},
			householdCreationInput.Name,
			types.UnpaidHouseholdBillingStatus,
			householdCreationInput.ContactPhone,
			householdCreationInput.AddressLine1,
			householdCreationInput.AddressLine2,
			householdCreationInput.City,
			householdCreationInput.State,
			householdCreationInput.ZipCode,
			householdCreationInput.Country,
			householdCreationInput.Latitude,
			householdCreationInput.Longitude,
			&idMatcher{},
		}

		db.ExpectExec(formatQueryForSQLMock(householdCreationQuery)).
			WithArgs(interfaceToDriverValue(householdCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		// create household user membership for created user
		createHouseholdMembershipForNewUserArgs := []any{
			&idMatcher{},
			&idMatcher{},
			&idMatcher{},
			false,
			authorization.HouseholdAdminRole.String(),
		}

		db.ExpectExec(formatQueryForSQLMock(createHouseholdMembershipForNewUserQuery)).
			WithArgs(interfaceToDriverValue(createHouseholdMembershipForNewUserArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		args := []any{
			exampleUserCreationInput.EmailAddress,
			exampleUserCreationInput.InvitationToken,
		}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdInvitationByEmailAndTokenQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		actual, err := c.CreateUser(ctx, exampleUserCreationInput)
		assert.Error(t, err)
		require.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, secretGenerator)
	})

	T.Run("with destination household and error creating new membership", func(t *testing.T) {
		t.Parallel()

		exampleCreationTime := fakes.BuildFakeTime()

		exampleUser := fakes.BuildFakeUser()
		exampleUser.TwoFactorSecretVerifiedAt = nil
		exampleUser.CreatedAt = exampleCreationTime
		exampleUser.AccountStatus = ""
		exampleUserCreationInput := converters.ConvertUserToUserDatabaseCreationInput(exampleUser)

		exampleHousehold := fakes.BuildFakeHouseholdForUser(exampleUser)
		exampleHousehold.CreatedAt = exampleCreationTime

		exampleUserCreationInput.DestinationHouseholdID = exampleHousehold.ID
		exampleUserCreationInput.InvitationToken = t.Name()

		ctx := context.Background()
		c, db := buildTestClient(t)

		c.timeFunc = func() time.Time {
			return exampleCreationTime
		}

		exampleToken := fakes.BuildFakeID()
		secretGenerator := &mockrandom.Generator{}
		secretGenerator.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			32,
		).Return(exampleToken, nil)
		c.secretGenerator = secretGenerator

		db.ExpectBegin()

		userCreationArgs := []any{
			exampleUserCreationInput.ID,
			exampleUserCreationInput.FirstName,
			exampleUserCreationInput.LastName,
			exampleUserCreationInput.Username,
			exampleUserCreationInput.EmailAddress,
			exampleUserCreationInput.HashedPassword,
			exampleUserCreationInput.TwoFactorSecret,
			exampleUserCreationInput.AvatarSrc,
			types.UnverifiedHouseholdStatus,
			exampleUserCreationInput.Birthday,
			authorization.ServiceUserRole.String(),
			exampleToken,
		}

		db.ExpectExec(formatQueryForSQLMock(userCreationQuery)).
			WithArgs(interfaceToDriverValue(userCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		// create household for created user
		householdCreationInput := types.HouseholdCreationInputForNewUser(exampleUser)
		householdCreationInput.ID = exampleHousehold.ID
		householdCreationArgs := []any{
			&idMatcher{},
			householdCreationInput.Name,
			types.UnpaidHouseholdBillingStatus,
			householdCreationInput.ContactPhone,
			householdCreationInput.AddressLine1,
			householdCreationInput.AddressLine2,
			householdCreationInput.City,
			householdCreationInput.State,
			householdCreationInput.ZipCode,
			householdCreationInput.Country,
			householdCreationInput.Latitude,
			householdCreationInput.Longitude,
			&idMatcher{},
		}

		db.ExpectExec(formatQueryForSQLMock(householdCreationQuery)).
			WithArgs(interfaceToDriverValue(householdCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		// create household user membership for created user
		createHouseholdMembershipForNewUserArgs := []any{
			&idMatcher{},
			&idMatcher{},
			&idMatcher{},
			false,
			authorization.HouseholdAdminRole.String(),
		}

		db.ExpectExec(formatQueryForSQLMock(createHouseholdMembershipForNewUserQuery)).
			WithArgs(interfaceToDriverValue(createHouseholdMembershipForNewUserArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		args := []any{
			exampleUserCreationInput.EmailAddress,
			exampleUserCreationInput.InvitationToken,
		}

		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdInvitationByEmailAndTokenQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholdInvitations(false, 0, exampleHouseholdInvitation))

		// create household user membership for created user
		createHouseholdMembershipForNewUserArgs = []any{
			&idMatcher{},
			&idMatcher{},
			&idMatcher{},
			true,
			authorization.HouseholdMemberRole.String(),
		}

		db.ExpectExec(formatQueryForSQLMock(createHouseholdMembershipForNewUserQuery)).
			WithArgs(interfaceToDriverValue(createHouseholdMembershipForNewUserArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		actual, err := c.CreateUser(ctx, exampleUserCreationInput)
		assert.Error(t, err)
		require.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, secretGenerator)
	})

	T.Run("with error executing household creation query", func(t *testing.T) {
		t.Parallel()

		exampleCreationTime := fakes.BuildFakeTime()

		exampleUser := fakes.BuildFakeUser()
		exampleUser.TwoFactorSecretVerifiedAt = nil
		exampleUser.CreatedAt = exampleCreationTime
		exampleUser.AccountStatus = ""
		exampleUserCreationInput := converters.ConvertUserToUserDatabaseCreationInput(exampleUser)

		exampleHousehold := fakes.BuildFakeHouseholdForUser(exampleUser)
		exampleHousehold.CreatedAt = exampleCreationTime

		ctx := context.Background()
		c, db := buildTestClient(t)

		c.timeFunc = func() time.Time {
			return exampleCreationTime
		}

		exampleToken := fakes.BuildFakeID()
		secretGenerator := &mockrandom.Generator{}
		secretGenerator.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			32,
		).Return(exampleToken, nil)
		c.secretGenerator = secretGenerator

		db.ExpectBegin()

		userCreationArgs := []any{
			exampleUserCreationInput.ID,
			exampleUserCreationInput.FirstName,
			exampleUserCreationInput.LastName,
			exampleUserCreationInput.Username,
			exampleUserCreationInput.EmailAddress,
			exampleUserCreationInput.HashedPassword,
			exampleUserCreationInput.TwoFactorSecret,
			exampleUserCreationInput.AvatarSrc,
			types.UnverifiedHouseholdStatus,
			exampleUserCreationInput.Birthday,
			authorization.ServiceUserRole.String(),
			exampleToken,
		}

		db.ExpectExec(formatQueryForSQLMock(userCreationQuery)).
			WithArgs(interfaceToDriverValue(userCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		// create household for created user
		householdCreationInput := types.HouseholdCreationInputForNewUser(exampleUser)
		householdCreationInput.ID = exampleHousehold.ID
		householdCreationArgs := []any{
			&idMatcher{},
			householdCreationInput.Name,
			types.UnpaidHouseholdBillingStatus,
			householdCreationInput.ContactPhone,
			householdCreationInput.AddressLine1,
			householdCreationInput.AddressLine2,
			householdCreationInput.City,
			householdCreationInput.State,
			householdCreationInput.ZipCode,
			householdCreationInput.Country,
			householdCreationInput.Latitude,
			householdCreationInput.Longitude,
			&idMatcher{},
		}

		db.ExpectExec(formatQueryForSQLMock(householdCreationQuery)).
			WithArgs(interfaceToDriverValue(householdCreationArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		actual, err := c.CreateUser(ctx, exampleUserCreationInput)
		assert.Error(t, err)
		require.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, secretGenerator)
	})

	T.Run("with error creating household user membership", func(t *testing.T) {
		t.Parallel()

		exampleCreationTime := fakes.BuildFakeTime()

		exampleUser := fakes.BuildFakeUser()
		exampleUser.TwoFactorSecretVerifiedAt = nil
		exampleUser.CreatedAt = exampleCreationTime
		exampleUser.AccountStatus = ""
		exampleUserCreationInput := converters.ConvertUserToUserDatabaseCreationInput(exampleUser)

		exampleHousehold := fakes.BuildFakeHouseholdForUser(exampleUser)
		exampleHousehold.CreatedAt = exampleCreationTime

		ctx := context.Background()
		c, db := buildTestClient(t)

		c.timeFunc = func() time.Time {
			return exampleCreationTime
		}

		exampleToken := fakes.BuildFakeID()
		secretGenerator := &mockrandom.Generator{}
		secretGenerator.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			32,
		).Return(exampleToken, nil)
		c.secretGenerator = secretGenerator

		db.ExpectBegin()

		userCreationArgs := []any{
			exampleUserCreationInput.ID,
			exampleUserCreationInput.FirstName,
			exampleUserCreationInput.LastName,
			exampleUserCreationInput.Username,
			exampleUserCreationInput.EmailAddress,
			exampleUserCreationInput.HashedPassword,
			exampleUserCreationInput.TwoFactorSecret,
			exampleUserCreationInput.AvatarSrc,
			types.UnverifiedHouseholdStatus,
			exampleUserCreationInput.Birthday,
			authorization.ServiceUserRole.String(),
			exampleToken,
		}

		db.ExpectExec(formatQueryForSQLMock(userCreationQuery)).
			WithArgs(interfaceToDriverValue(userCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		// create household for created user
		householdCreationInput := types.HouseholdCreationInputForNewUser(exampleUser)
		householdCreationInput.ID = exampleHousehold.ID
		householdCreationArgs := []any{
			&idMatcher{},
			householdCreationInput.Name,
			types.UnpaidHouseholdBillingStatus,
			householdCreationInput.ContactPhone,
			householdCreationInput.AddressLine1,
			householdCreationInput.AddressLine2,
			householdCreationInput.City,
			householdCreationInput.State,
			householdCreationInput.ZipCode,
			householdCreationInput.Country,
			householdCreationInput.Latitude,
			householdCreationInput.Longitude,
			&idMatcher{},
		}

		db.ExpectExec(formatQueryForSQLMock(householdCreationQuery)).
			WithArgs(interfaceToDriverValue(householdCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		// create household user membership for created user
		createHouseholdMembershipForNewUserArgs := []any{
			&idMatcher{},
			&idMatcher{},
			&idMatcher{},
			true,
			authorization.HouseholdAdminRole.String(),
		}

		db.ExpectExec(formatQueryForSQLMock(createHouseholdMembershipForNewUserQuery)).
			WithArgs(interfaceToDriverValue(createHouseholdMembershipForNewUserArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		actual, err := c.CreateUser(ctx, exampleUserCreationInput)
		assert.Error(t, err)
		require.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, secretGenerator)
	})

	T.Run("with error attaching invitations to user", func(t *testing.T) {
		t.Parallel()

		exampleCreationTime := fakes.BuildFakeTime()

		exampleUser := fakes.BuildFakeUser()
		exampleUser.TwoFactorSecretVerifiedAt = nil
		exampleUser.CreatedAt = exampleCreationTime
		exampleUser.AccountStatus = ""
		exampleUserCreationInput := converters.ConvertUserToUserDatabaseCreationInput(exampleUser)

		exampleHousehold := fakes.BuildFakeHouseholdForUser(exampleUser)
		exampleHousehold.CreatedAt = exampleCreationTime

		ctx := context.Background()
		c, db := buildTestClient(t)

		c.timeFunc = func() time.Time {
			return exampleCreationTime
		}

		exampleToken := fakes.BuildFakeID()
		secretGenerator := &mockrandom.Generator{}
		secretGenerator.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			32,
		).Return(exampleToken, nil)
		c.secretGenerator = secretGenerator

		db.ExpectBegin()

		userCreationArgs := []any{
			exampleUserCreationInput.ID,
			exampleUserCreationInput.FirstName,
			exampleUserCreationInput.LastName,
			exampleUserCreationInput.Username,
			exampleUserCreationInput.EmailAddress,
			exampleUserCreationInput.HashedPassword,
			exampleUserCreationInput.TwoFactorSecret,
			exampleUserCreationInput.AvatarSrc,
			types.UnverifiedHouseholdStatus,
			exampleUserCreationInput.Birthday,
			authorization.ServiceUserRole.String(),
			exampleToken,
		}

		db.ExpectExec(formatQueryForSQLMock(userCreationQuery)).
			WithArgs(interfaceToDriverValue(userCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		// create household for created user
		householdCreationInput := types.HouseholdCreationInputForNewUser(exampleUser)
		householdCreationInput.ID = exampleHousehold.ID
		householdCreationArgs := []any{
			&idMatcher{},
			householdCreationInput.Name,
			types.UnpaidHouseholdBillingStatus,
			householdCreationInput.ContactPhone,
			householdCreationInput.AddressLine1,
			householdCreationInput.AddressLine2,
			householdCreationInput.City,
			householdCreationInput.State,
			householdCreationInput.ZipCode,
			householdCreationInput.Country,
			householdCreationInput.Latitude,
			householdCreationInput.Longitude,
			&idMatcher{},
		}

		db.ExpectExec(formatQueryForSQLMock(householdCreationQuery)).
			WithArgs(interfaceToDriverValue(householdCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		// create household user membership for created user
		createHouseholdMembershipForNewUserArgs := []any{
			&idMatcher{},
			&idMatcher{},
			&idMatcher{},
			true,
			authorization.HouseholdAdminRole.String(),
		}

		db.ExpectExec(formatQueryForSQLMock(createHouseholdMembershipForNewUserQuery)).
			WithArgs(interfaceToDriverValue(createHouseholdMembershipForNewUserArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		// create household user membership for created user
		attachInvitationsToUsersArgs := []any{
			&idMatcher{},
			exampleUser.EmailAddress,
		}

		db.ExpectExec(formatQueryForSQLMock(attachInvitationsToUserIDQuery)).
			WithArgs(interfaceToDriverValue(attachInvitationsToUsersArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		actual, err := c.CreateUser(ctx, exampleUserCreationInput)
		assert.Error(t, err)
		require.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, secretGenerator)
	})

	T.Run("with error committing transaction", func(t *testing.T) {
		t.Parallel()

		exampleCreationTime := fakes.BuildFakeTime()

		exampleUser := fakes.BuildFakeUser()
		exampleUser.TwoFactorSecretVerifiedAt = nil
		exampleUser.CreatedAt = exampleCreationTime
		exampleUser.AccountStatus = ""
		exampleUserCreationInput := converters.ConvertUserToUserDatabaseCreationInput(exampleUser)

		exampleHousehold := fakes.BuildFakeHouseholdForUser(exampleUser)
		exampleHousehold.CreatedAt = exampleCreationTime

		ctx := context.Background()
		c, db := buildTestClient(t)

		c.timeFunc = func() time.Time {
			return exampleCreationTime
		}

		exampleToken := fakes.BuildFakeID()
		secretGenerator := &mockrandom.Generator{}
		secretGenerator.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			32,
		).Return(exampleToken, nil)
		c.secretGenerator = secretGenerator

		db.ExpectBegin()

		userCreationArgs := []any{
			exampleUserCreationInput.ID,
			exampleUserCreationInput.FirstName,
			exampleUserCreationInput.LastName,
			exampleUserCreationInput.Username,
			exampleUserCreationInput.EmailAddress,
			exampleUserCreationInput.HashedPassword,
			exampleUserCreationInput.TwoFactorSecret,
			exampleUserCreationInput.AvatarSrc,
			types.UnverifiedHouseholdStatus,
			exampleUserCreationInput.Birthday,
			authorization.ServiceUserRole.String(),
			exampleToken,
		}

		db.ExpectExec(formatQueryForSQLMock(userCreationQuery)).
			WithArgs(interfaceToDriverValue(userCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		// create household for created user
		householdCreationInput := types.HouseholdCreationInputForNewUser(exampleUser)
		householdCreationInput.ID = exampleHousehold.ID
		householdCreationArgs := []any{
			&idMatcher{},
			householdCreationInput.Name,
			types.UnpaidHouseholdBillingStatus,
			householdCreationInput.ContactPhone,
			householdCreationInput.AddressLine1,
			householdCreationInput.AddressLine2,
			householdCreationInput.City,
			householdCreationInput.State,
			householdCreationInput.ZipCode,
			householdCreationInput.Country,
			householdCreationInput.Latitude,
			householdCreationInput.Longitude,
			&idMatcher{},
		}

		db.ExpectExec(formatQueryForSQLMock(householdCreationQuery)).
			WithArgs(interfaceToDriverValue(householdCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		// create household user membership for created user
		createHouseholdMembershipForNewUserArgs := []any{
			&idMatcher{},
			&idMatcher{},
			&idMatcher{},
			true,
			authorization.HouseholdAdminRole.String(),
		}

		db.ExpectExec(formatQueryForSQLMock(createHouseholdMembershipForNewUserQuery)).
			WithArgs(interfaceToDriverValue(createHouseholdMembershipForNewUserArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		// create household user membership for created user
		attachInvitationsToUsersArgs := []any{
			&idMatcher{},
			exampleUser.EmailAddress,
		}

		db.ExpectExec(formatQueryForSQLMock(attachInvitationsToUserIDQuery)).
			WithArgs(interfaceToDriverValue(attachInvitationsToUsersArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		actual, err := c.CreateUser(ctx, exampleUserCreationInput)
		assert.Error(t, err)
		require.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, secretGenerator)
	})

	T.Run("with error accepting household invitation", func(t *testing.T) {
		t.Parallel()

		exampleCreationTime := fakes.BuildFakeTime()

		exampleUser := fakes.BuildFakeUser()
		exampleUser.TwoFactorSecretVerifiedAt = nil
		exampleUser.CreatedAt = exampleCreationTime
		exampleUser.AccountStatus = ""
		exampleUserCreationInput := converters.ConvertUserToUserDatabaseCreationInput(exampleUser)

		exampleHousehold := fakes.BuildFakeHouseholdForUser(exampleUser)
		exampleHousehold.CreatedAt = exampleCreationTime

		exampleUserCreationInput.DestinationHouseholdID = exampleHousehold.ID
		exampleUserCreationInput.InvitationToken = t.Name()

		ctx := context.Background()
		c, db := buildTestClient(t)

		c.timeFunc = func() time.Time {
			return exampleCreationTime
		}

		exampleToken := fakes.BuildFakeID()
		secretGenerator := &mockrandom.Generator{}
		secretGenerator.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			32,
		).Return(exampleToken, nil)
		c.secretGenerator = secretGenerator

		db.ExpectBegin()

		userCreationArgs := []any{
			exampleUserCreationInput.ID,
			exampleUserCreationInput.FirstName,
			exampleUserCreationInput.LastName,
			exampleUserCreationInput.Username,
			exampleUserCreationInput.EmailAddress,
			exampleUserCreationInput.HashedPassword,
			exampleUserCreationInput.TwoFactorSecret,
			exampleUserCreationInput.AvatarSrc,
			types.UnverifiedHouseholdStatus,
			exampleUserCreationInput.Birthday,
			authorization.ServiceUserRole.String(),
			exampleToken,
		}

		db.ExpectExec(formatQueryForSQLMock(userCreationQuery)).
			WithArgs(interfaceToDriverValue(userCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		// create household for created user
		householdCreationInput := types.HouseholdCreationInputForNewUser(exampleUser)
		householdCreationInput.ID = exampleHousehold.ID
		householdCreationArgs := []any{
			&idMatcher{},
			householdCreationInput.Name,
			types.UnpaidHouseholdBillingStatus,
			householdCreationInput.ContactPhone,
			householdCreationInput.AddressLine1,
			householdCreationInput.AddressLine2,
			householdCreationInput.City,
			householdCreationInput.State,
			householdCreationInput.ZipCode,
			householdCreationInput.Country,
			householdCreationInput.Latitude,
			householdCreationInput.Longitude,
			&idMatcher{},
		}

		db.ExpectExec(formatQueryForSQLMock(householdCreationQuery)).
			WithArgs(interfaceToDriverValue(householdCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		// create household user membership for created user
		createHouseholdMembershipForNewUserArgs := []any{
			&idMatcher{},
			&idMatcher{},
			&idMatcher{},
			false,
			authorization.HouseholdAdminRole.String(),
		}

		db.ExpectExec(formatQueryForSQLMock(createHouseholdMembershipForNewUserQuery)).
			WithArgs(interfaceToDriverValue(createHouseholdMembershipForNewUserArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		args := []any{
			exampleUserCreationInput.EmailAddress,
			exampleUserCreationInput.InvitationToken,
		}

		exampleHouseholdInvitation := fakes.BuildFakeHouseholdInvitation()

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdInvitationByEmailAndTokenQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholdInvitations(false, 0, exampleHouseholdInvitation))

		// create household user membership for created user
		createHouseholdMembershipForNewUserArgs = []any{
			&idMatcher{},
			&idMatcher{},
			&idMatcher{},
			true,
			authorization.HouseholdMemberRole.String(),
		}

		db.ExpectExec(formatQueryForSQLMock(createHouseholdMembershipForNewUserQuery)).
			WithArgs(interfaceToDriverValue(createHouseholdMembershipForNewUserArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		invitationStatusUpdateArgs := []any{
			types.AcceptedHouseholdInvitationStatus,
			"",
			exampleHouseholdInvitation.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(setInvitationStatusQuery)).
			WithArgs(interfaceToDriverValue(invitationStatusUpdateArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		actual, err := c.CreateUser(ctx, exampleUserCreationInput)
		assert.Error(t, err)
		require.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, secretGenerator)
	})
}

func TestQuerier_UpdateUsername(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleNewUsername := fakes.BuildFakeUser().Username

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleNewUsername,
			exampleUser.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateUsernameQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateUserUsername(ctx, exampleUser.ID, exampleNewUsername))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with empty user", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateUserUsername(ctx, "", t.Name()))
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleNewUsername := fakes.BuildFakeUser().Username

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleNewUsername,
			exampleUser.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateUsernameQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateUserUsername(ctx, exampleUser.ID, exampleNewUsername))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateUserDetails(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleInput := fakes.BuildFakeUserDetailsUpdateInput()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.FirstName,
			exampleInput.LastName,
			exampleInput.Birthday,
			exampleUser.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateUserDetailsQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateUserDetails(ctx, exampleUser.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleInput := fakes.BuildFakeUserDetailsUpdateInput()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.FirstName,
			exampleInput.LastName,
			exampleInput.Birthday,
			exampleUser.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateUserDetailsQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateUserDetails(ctx, exampleUser.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateUserAvatar(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleInput := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		updateUserAvatarSrcArgs := []any{
			exampleInput,
			exampleUser.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateUserAvatarSrcQuery)).
			WithArgs(interfaceToDriverValue(updateUserAvatarSrcArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateUserAvatar(ctx, exampleUser.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleInput := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		updateUserAvatarSrcArgs := []any{
			exampleInput,
			exampleUser.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateUserAvatarSrcQuery)).
			WithArgs(interfaceToDriverValue(updateUserAvatarSrcArgs)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateUserAvatar(ctx, exampleUser.ID, exampleInput))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateUserPassword(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleHashedPassword := "$argon2i$v=19$m=64,t=10,p=4$RjFtMmRmU2lGYU9CMk1mMw$cuGR9AhTczPR6xDOSAMW+SvEYFyLEIS+7nlRdC9f6ys"

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleHashedPassword,
			false,
			exampleUserID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateUserPasswordQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateUserPassword(ctx, exampleUserID, exampleHashedPassword))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		exampleHashedPassword := "$argon2i$v=19$m=64,t=10,p=4$RjFtMmRmU2lGYU9CMk1mMw$cuGR9AhTczPR6xDOSAMW+SvEYFyLEIS+7nlRdC9f6ys"

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateUserPassword(ctx, "", exampleHashedPassword))
	})

	T.Run("with invalid new hash", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateUserPassword(ctx, exampleUser.ID, ""))
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleHashedPassword := "$argon2i$v=19$m=64,t=10,p=4$RjFtMmRmU2lGYU9CMk1mMw$cuGR9AhTczPR6xDOSAMW+SvEYFyLEIS+7nlRdC9f6ys"

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleHashedPassword,
			false,
			exampleUserID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateUserPasswordQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateUserPassword(ctx, exampleUserID, exampleHashedPassword))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateUserTwoFactorSecret(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			nil,
			exampleUser.TwoFactorSecret,
			exampleUser.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateUserTwoFactorSecretQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateUserTwoFactorSecret(ctx, exampleUser.ID, exampleUser.TwoFactorSecret))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateUserTwoFactorSecret(ctx, "", exampleUser.TwoFactorSecret))
	})

	T.Run("with invalid new secret", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateUserTwoFactorSecret(ctx, exampleUser.ID, ""))
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			nil,
			exampleUser.TwoFactorSecret,
			exampleUser.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateUserTwoFactorSecretQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateUserTwoFactorSecret(ctx, exampleUser.ID, exampleUser.TwoFactorSecret))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_MarkUserTwoFactorSecretAsVerified(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			types.GoodStandingUserAccountStatus,
			exampleUserID,
		}

		db.ExpectExec(formatQueryForSQLMock(markUserTwoFactorSecretAsVerified)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(sqlmock.NewResult(1, 1))

		assert.NoError(t, c.MarkUserTwoFactorSecretAsVerified(ctx, exampleUserID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkUserTwoFactorSecretAsVerified(ctx, ""))
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			types.GoodStandingUserAccountStatus,
			exampleUserID,
		}

		db.ExpectExec(formatQueryForSQLMock(markUserTwoFactorSecretAsVerified)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.MarkUserTwoFactorSecretAsVerified(ctx, exampleUserID))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_MarkUserTwoFactorSecretAsUnverified(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleSecret := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleSecret,
			exampleUserID,
		}

		db.ExpectExec(formatQueryForSQLMock(markUserTwoFactorSecretAsUnverified)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(sqlmock.NewResult(1, 1))

		assert.NoError(t, c.MarkUserTwoFactorSecretAsUnverified(ctx, exampleUserID, exampleSecret))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleSecret := fakes.BuildFakeID()
		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkUserTwoFactorSecretAsUnverified(ctx, "", exampleSecret))
	})

	T.Run("with invalid secret", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkUserTwoFactorSecretAsUnverified(ctx, exampleUserID, ""))
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleSecret := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleSecret,
			exampleUserID,
		}

		db.ExpectExec(formatQueryForSQLMock(markUserTwoFactorSecretAsUnverified)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.MarkUserTwoFactorSecretAsUnverified(ctx, exampleUserID, exampleSecret))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveUser(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		archiveUserArgs := []any{exampleUserID}

		db.ExpectExec(formatQueryForSQLMock(archiveUserQuery)).
			WithArgs(interfaceToDriverValue(archiveUserArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		archiveMembershipsArgs := []any{exampleUserID}

		db.ExpectExec(formatQueryForSQLMock(archiveMembershipsQuery)).
			WithArgs(interfaceToDriverValue(archiveMembershipsArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		db.ExpectCommit()

		assert.NoError(t, c.ArchiveUser(ctx, exampleUserID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveUser(ctx, ""))
	})

	T.Run("with error beginning transaction", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveUser(ctx, exampleUserID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing user archive query", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		archiveUserArgs := []any{exampleUserID}

		db.ExpectExec(formatQueryForSQLMock(archiveUserQuery)).
			WithArgs(interfaceToDriverValue(archiveUserArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		assert.Error(t, c.ArchiveUser(ctx, exampleUserID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing memberships archive query", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		archiveUserArgs := []any{exampleUserID}

		db.ExpectExec(formatQueryForSQLMock(archiveUserQuery)).
			WithArgs(interfaceToDriverValue(archiveUserArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		archiveMembershipsArgs := []any{exampleUserID}

		db.ExpectExec(formatQueryForSQLMock(archiveMembershipsQuery)).
			WithArgs(interfaceToDriverValue(archiveMembershipsArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		assert.Error(t, c.ArchiveUser(ctx, exampleUserID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error committing transaction", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		archiveUserArgs := []any{exampleUserID}

		db.ExpectExec(formatQueryForSQLMock(archiveUserQuery)).
			WithArgs(interfaceToDriverValue(archiveUserArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		archiveMembershipsArgs := []any{exampleUserID}

		db.ExpectExec(formatQueryForSQLMock(archiveMembershipsQuery)).
			WithArgs(interfaceToDriverValue(archiveMembershipsArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveUser(ctx, exampleUserID))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetUserByEmailAddressVerificationToken(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleInput := fakes.BuildFakeEmailAddressVerificationRequestInput()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{exampleInput.Token}

		db.ExpectQuery(formatQueryForSQLMock(getUserByEmailAddressVerificationTokenQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromUsers(false, 0, exampleUser))

		actual, err := c.GetUserByEmailAddressVerificationToken(ctx, exampleInput.Token)
		assert.NoError(t, err)
		assert.Equal(t, exampleUser, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with missing token", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetUserByEmailAddressVerificationToken(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleInput := fakes.BuildFakeEmailAddressVerificationRequestInput()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{exampleInput.Token}

		db.ExpectQuery(formatQueryForSQLMock(getUserByEmailAddressVerificationTokenQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetUserByEmailAddressVerificationToken(ctx, exampleInput.Token)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_MarkUserEmailAddressAsVerified(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleInput := fakes.BuildFakeEmailAddressVerificationRequestInput()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleUser.ID,
			exampleInput.Token,
		}

		db.ExpectExec(formatQueryForSQLMock(markEmailAddressAsVerifiedQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		err := c.MarkUserEmailAddressAsVerified(ctx, exampleUser.ID, exampleInput.Token)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with missing user ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fakes.BuildFakeEmailAddressVerificationRequestInput()

		ctx := context.Background()
		c, db := buildTestClient(t)

		err := c.MarkUserEmailAddressAsVerified(ctx, "", exampleInput.Token)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with missing token", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		c, db := buildTestClient(t)

		err := c.MarkUserEmailAddressAsVerified(ctx, exampleUser.ID, "")
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleInput := fakes.BuildFakeEmailAddressVerificationRequestInput()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleUser.ID,
			exampleInput.Token,
		}

		db.ExpectExec(formatQueryForSQLMock(markEmailAddressAsVerifiedQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		err := c.MarkUserEmailAddressAsVerified(ctx, exampleUser.ID, exampleInput.Token)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})
}
