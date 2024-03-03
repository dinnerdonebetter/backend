package database

import (
	"context"
	"database/sql"
	"time"

	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/alexedwards/scs/v2"
	"github.com/stretchr/testify/mock"
)

var _ DataManager = (*MockDatabase)(nil)

// NewMockDatabase builds a mock database.
func NewMockDatabase() *MockDatabase {
	return &MockDatabase{
		HouseholdDataManagerMock:                      &mocktypes.HouseholdDataManagerMock{},
		HouseholdInvitationDataManagerMock:            &mocktypes.HouseholdInvitationDataManagerMock{},
		HouseholdUserMembershipDataManagerMock:        &mocktypes.HouseholdUserMembershipDataManagerMock{},
		ValidInstrumentDataManagerMock:                &mocktypes.ValidInstrumentDataManagerMock{},
		ValidIngredientDataManagerMock:                &mocktypes.ValidIngredientDataManagerMock{},
		ValidIngredientGroupDataManagerMock:           &mocktypes.ValidIngredientGroupDataManagerMock{},
		ValidPreparationDataManagerMock:               &mocktypes.ValidPreparationDataManagerMock{},
		ValidIngredientPreparationDataManagerMock:     &mocktypes.ValidIngredientPreparationDataManagerMock{},
		MealDataManagerMock:                           &mocktypes.MealDataManagerMock{},
		RecipeDataManagerMock:                         &mocktypes.RecipeDataManagerMock{},
		RecipeStepDataManagerMock:                     &mocktypes.RecipeStepDataManagerMock{},
		RecipeStepProductDataManagerMock:              &mocktypes.RecipeStepProductDataManagerMock{},
		RecipeStepInstrumentDataManagerMock:           &mocktypes.RecipeStepInstrumentDataManagerMock{},
		RecipeStepIngredientDataManagerMock:           &mocktypes.RecipeStepIngredientDataManagerMock{},
		MealPlanDataManagerMock:                       &mocktypes.MealPlanDataManagerMock{},
		MealPlanOptionDataManagerMock:                 &mocktypes.MealPlanOptionDataManagerMock{},
		MealPlanOptionVoteDataManagerMock:             &mocktypes.MealPlanOptionVoteDataManagerMock{},
		UserDataManagerMock:                           &mocktypes.UserDataManagerMock{},
		AdminUserDataManagerMock:                      &mocktypes.AdminUserDataManagerMock{},
		PasswordResetTokenDataManagerMock:             &mocktypes.PasswordResetTokenDataManagerMock{},
		WebhookDataManagerMock:                        &mocktypes.WebhookDataManagerMock{},
		ValidMeasurementUnitDataManagerMock:           &mocktypes.ValidMeasurementUnitDataManagerMock{},
		ValidPreparationInstrumentDataManagerMock:     &mocktypes.ValidPreparationInstrumentDataManagerMock{},
		ValidIngredientMeasurementUnitDataManagerMock: &mocktypes.ValidIngredientMeasurementUnitDataManagerMock{},
		MealPlanEventDataManagerMock:                  &mocktypes.MealPlanEventDataManagerMock{},
		MealPlanTaskDataManagerMock:                   &mocktypes.MealPlanTaskDataManagerMock{},
		RecipePrepTaskDataManagerMock:                 &mocktypes.RecipePrepTaskDataManagerMock{},
		MealPlanGroceryListItemDataManagerMock:        &mocktypes.MealPlanGroceryListItemDataManagerMock{},
		ValidMeasurementUnitConversionDataManagerMock: &mocktypes.ValidMeasurementUnitConversionDataManagerMock{},
		RecipeMediaDataManagerMock:                    &mocktypes.RecipeMediaDataManagerMock{},
		ValidIngredientStateDataManagerMock:           &mocktypes.ValidIngredientStateDataManagerMock{},
		ValidIngredientStateIngredientDataManagerMock: &mocktypes.ValidIngredientStateIngredientDataManagerMock{},
		RecipeStepCompletionConditionDataManagerMock:  &mocktypes.RecipeStepCompletionConditionDataManagerMock{},
		RecipeStepVesselDataManagerMock:               &mocktypes.RecipeStepVesselDataManagerMock{},
		ServiceSettingDataManagerMock:                 &mocktypes.ServiceSettingDataManagerMock{},
		ServiceSettingConfigurationDataManagerMock:    &mocktypes.ServiceSettingConfigurationDataManagerMock{},
		UserIngredientPreferenceDataManagerMock:       &mocktypes.UserIngredientPreferenceDataManagerMock{},
		HouseholdInstrumentOwnershipDataManagerMock:   &mocktypes.HouseholdInstrumentOwnershipDataManagerMock{},
		RecipeRatingDataManagerMock:                   &mocktypes.RecipeRatingDataManagerMock{},
		OAuth2ClientDataManagerMock:                   &mocktypes.OAuth2ClientDataManagerMock{},
		ValidVesselDataManagerMock:                    &mocktypes.ValidVesselDataManagerMock{},
		ValidPreparationVesselDataManagerMock:         &mocktypes.ValidPreparationVesselDataManagerMock{},
		UserNotificationDataManagerMock:               &mocktypes.UserNotificationDataManagerMock{},
		AuditLogEntryDataManagerMock:                  &mocktypes.AuditLogEntryDataManagerMock{},
	}
}

// MockDatabase is our mock database structure. Note, when using this in tests, you must directly access the type name of all the implicit fields.
// So `mockDB.On("GetUserByUsername"...)` is destined to fail, whereas `mockDB.UserDataManagerMock.On("GetUserByUsername"...)` would do what you want it to do.
type MockDatabase struct {
	*mocktypes.ValidIngredientStateDataManagerMock
	*mocktypes.AdminUserDataManagerMock
	*mocktypes.HouseholdUserMembershipDataManagerMock
	*mocktypes.ValidInstrumentDataManagerMock
	*mocktypes.ValidIngredientDataManagerMock
	*mocktypes.ValidIngredientGroupDataManagerMock
	*mocktypes.ValidPreparationDataManagerMock
	*mocktypes.ValidIngredientPreparationDataManagerMock
	*mocktypes.MealDataManagerMock
	*mocktypes.RecipeDataManagerMock
	*mocktypes.RecipeStepDataManagerMock
	*mocktypes.RecipeStepProductDataManagerMock
	*mocktypes.RecipeStepInstrumentDataManagerMock
	*mocktypes.RecipeStepIngredientDataManagerMock
	*mocktypes.MealPlanDataManagerMock
	*mocktypes.MealPlanOptionDataManagerMock
	*mocktypes.MealPlanOptionVoteDataManagerMock
	*mocktypes.UserDataManagerMock
	*mocktypes.PasswordResetTokenDataManagerMock
	*mocktypes.WebhookDataManagerMock
	*mocktypes.HouseholdDataManagerMock
	*mocktypes.HouseholdInvitationDataManagerMock
	*mocktypes.ValidMeasurementUnitDataManagerMock
	*mocktypes.ValidPreparationInstrumentDataManagerMock
	*mocktypes.ValidIngredientMeasurementUnitDataManagerMock
	*mocktypes.MealPlanEventDataManagerMock
	*mocktypes.MealPlanTaskDataManagerMock
	*mocktypes.RecipePrepTaskDataManagerMock
	*mocktypes.MealPlanGroceryListItemDataManagerMock
	*mocktypes.ValidMeasurementUnitConversionDataManagerMock
	*mocktypes.RecipeMediaDataManagerMock
	*mocktypes.RecipeStepCompletionConditionDataManagerMock
	*mocktypes.ValidIngredientStateIngredientDataManagerMock
	*mocktypes.RecipeStepVesselDataManagerMock
	*mocktypes.ServiceSettingDataManagerMock
	*mocktypes.ServiceSettingConfigurationDataManagerMock
	*mocktypes.UserIngredientPreferenceDataManagerMock
	*mocktypes.HouseholdInstrumentOwnershipDataManagerMock
	*mocktypes.RecipeRatingDataManagerMock
	*mocktypes.OAuth2ClientDataManagerMock
	*mocktypes.OAuth2ClientTokenDataManagerMock
	*mocktypes.ValidVesselDataManagerMock
	*mocktypes.ValidPreparationVesselDataManagerMock
	*mocktypes.UserNotificationDataManagerMock
	*mocktypes.AuditLogEntryDataManagerMock

	mock.Mock
}

// ProvideSessionStore satisfies the DataManager interface.
func (m *MockDatabase) ProvideSessionStore() scs.Store {
	return m.Called().Get(0).(scs.Store)
}

// Migrate satisfies the DataManager interface.
func (m *MockDatabase) Migrate(ctx context.Context, waitPeriod time.Duration, maxAttempts uint64) error {
	return m.Called(ctx, waitPeriod, maxAttempts).Error(0)
}

// Close satisfies the DataManager interface.
func (m *MockDatabase) Close() {
	m.Called()
}

// DB satisfies the DataManager interface.
func (m *MockDatabase) DB() *sql.DB {
	return m.Called().Get(0).(*sql.DB)
}

// IsReady satisfies the DataManager interface.
func (m *MockDatabase) IsReady(ctx context.Context, waitPeriod time.Duration, maxAttempts uint64) (ready bool) {
	return m.Called(ctx, waitPeriod, maxAttempts).Bool(0)
}

// BeginTx satisfies the DataManager interface.
func (m *MockDatabase) BeginTx(ctx context.Context, options *sql.TxOptions) (*sql.Tx, error) {
	args := m.Called(ctx, options)
	return args.Get(0).(*sql.Tx), args.Error(1)
}

var _ ResultIterator = (*MockResultIterator)(nil)

// MockResultIterator is our mock sql.Rows structure.
type MockResultIterator struct {
	mock.Mock
}

// Scan satisfies the ResultIterator interface.
func (m *MockResultIterator) Scan(dest ...any) error {
	return m.Called(dest...).Error(0)
}

// Next satisfies the ResultIterator interface.
func (m *MockResultIterator) Next() bool {
	return m.Called().Bool(0)
}

// Err satisfies the ResultIterator interface.
func (m *MockResultIterator) Err() error {
	return m.Called().Error(0)
}

// Close satisfies the ResultIterator interface.
func (m *MockResultIterator) Close() error {
	return m.Called().Error(0)
}

// MockSQLResult mocks a sql.Result.
type MockSQLResult struct {
	mock.Mock
}

// LastInsertId implements our interface.
func (m *MockSQLResult) LastInsertId() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

// RowsAffected implements our interface.
func (m *MockSQLResult) RowsAffected() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

var _ SQLQueryExecutor = (*MockQueryExecutor)(nil)

// MockQueryExecutor mocks a sql.Tx|DB.
type MockQueryExecutor struct {
	mock.Mock
}

// ExecContext is a mock function.
func (m *MockQueryExecutor) ExecContext(ctx context.Context, query string, queryArgs ...any) (sql.Result, error) {
	args := m.Called(ctx, query, queryArgs)
	return args.Get(0).(sql.Result), args.Error(1)
}

// PrepareContext is a mock function.
func (m *MockQueryExecutor) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	args := m.Called(ctx, query)
	return args.Get(0).(*sql.Stmt), args.Error(1)
}

// QueryContext is a mock function.
func (m *MockQueryExecutor) QueryContext(ctx context.Context, query string, queryArgs ...any) (*sql.Rows, error) {
	args := m.Called(ctx, query, queryArgs)
	return args.Get(0).(*sql.Rows), args.Error(1)
}

// QueryRowContext is a mock function.
func (m *MockQueryExecutor) QueryRowContext(ctx context.Context, query string, queryArgs ...any) *sql.Row {
	args := m.Called(ctx, query, queryArgs)
	return args.Get(0).(*sql.Row)
}
