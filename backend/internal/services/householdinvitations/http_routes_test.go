package householdinvitations

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	randommock "github.com/dinnerdonebetter/backend/internal/pkg/random/mock"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"
	testutils "github.com/dinnerdonebetter/backend/tests/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_service_InviteMemberHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeHouseholdInvitationCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		udm := &mocktypes.UserDataManagerMock{}
		udm.On(
			"GetUserByEmail",
			testutils.ContextMatcher,
			strings.TrimSpace(strings.ToLower(exampleInput.ToEmail)),
		).Return(helper.exampleUser, nil)
		helper.service.userDataManager = udm

		sg := &randommock.Generator{}
		sg.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			64,
		).Return(t.Name(), nil)
		helper.service.secretGenerator = sg

		dbManager := database.NewMockDatabase()
		dbManager.HouseholdInvitationDataManagerMock.On(
			"CreateHouseholdInvitation",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.HouseholdInvitationDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleHouseholdInvitation, nil)
		helper.service.householdInvitationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.InviteMemberHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		helper.exampleHouseholdInvitation.FromUser.TwoFactorSecret = ""
		helper.exampleHouseholdInvitation.DestinationHousehold.WebhookEncryptionKey = ""
		assert.Equal(t, actual.Data, helper.exampleHouseholdInvitation)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, udm, sg, dbManager, dataChangesPublisher)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)

		helper.service.sessionContextDataFetcher = func(*http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		helper.service.InviteMemberHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := &types.HouseholdInvitationCreationRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.InviteMemberHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error generating invitation token", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeHouseholdInvitationCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		sg := &randommock.Generator{}
		sg.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			64,
		).Return("", errors.New("blah"))
		helper.service.secretGenerator = sg

		helper.service.InviteMemberHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, sg)
	})

	T.Run("with error fetching user ID", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeHouseholdInvitationCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		sg := &randommock.Generator{}
		sg.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			64,
		).Return(t.Name(), nil)
		helper.service.secretGenerator = sg

		udm := &mocktypes.UserDataManagerMock{}
		udm.On(
			"GetUserByEmail",
			testutils.ContextMatcher,
			strings.TrimSpace(strings.ToLower(exampleInput.ToEmail)),
		).Return((*types.User)(nil), errors.New("blah"))
		helper.service.userDataManager = udm

		helper.service.InviteMemberHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, udm, sg)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeHouseholdInvitationCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		udm := &mocktypes.UserDataManagerMock{}
		udm.On(
			"GetUserByEmail",
			testutils.ContextMatcher,
			strings.TrimSpace(strings.ToLower(exampleInput.ToEmail)),
		).Return(helper.exampleUser, nil)
		helper.service.userDataManager = udm

		sg := &randommock.Generator{}
		sg.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			64,
		).Return(t.Name(), nil)
		helper.service.secretGenerator = sg

		dbManager := database.NewMockDatabase()
		dbManager.HouseholdInvitationDataManagerMock.On(
			"CreateHouseholdInvitation",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.HouseholdInvitationDatabaseCreationInput) bool { return true }),
		).Return((*types.HouseholdInvitation)(nil), errors.New("blah"))
		helper.service.householdInvitationDataManager = dbManager

		helper.service.InviteMemberHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, udm, sg, dbManager)
	})

	T.Run("with error publishing message", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeHouseholdInvitationCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		udm := &mocktypes.UserDataManagerMock{}
		udm.On(
			"GetUserByEmail",
			testutils.ContextMatcher,
			strings.TrimSpace(strings.ToLower(exampleInput.ToEmail)),
		).Return(helper.exampleUser, nil)
		helper.service.userDataManager = udm

		sg := &randommock.Generator{}
		sg.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			64,
		).Return(t.Name(), nil)
		helper.service.secretGenerator = sg

		dbManager := database.NewMockDatabase()
		dbManager.HouseholdInvitationDataManagerMock.On(
			"CreateHouseholdInvitation",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.HouseholdInvitationDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleHouseholdInvitation, nil)
		helper.service.householdInvitationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.InviteMemberHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		helper.exampleHouseholdInvitation.FromUser.TwoFactorSecret = ""
		helper.exampleHouseholdInvitation.DestinationHousehold.WebhookEncryptionKey = ""
		assert.Equal(t, actual.Data, helper.exampleHouseholdInvitation)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, udm, sg, dbManager, dataChangesPublisher)
	})

	T.Run("with error collecting data", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeHouseholdInvitationCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		udm := &mocktypes.UserDataManagerMock{}
		udm.On(
			"GetUserByEmail",
			testutils.ContextMatcher,
			strings.TrimSpace(strings.ToLower(exampleInput.ToEmail)),
		).Return(helper.exampleUser, nil)
		helper.service.userDataManager = udm

		sg := &randommock.Generator{}
		sg.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			64,
		).Return(t.Name(), nil)
		helper.service.secretGenerator = sg

		dbManager := database.NewMockDatabase()
		dbManager.HouseholdInvitationDataManagerMock.On(
			"CreateHouseholdInvitation",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.HouseholdInvitationDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleHouseholdInvitation, nil)
		helper.service.householdInvitationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.InviteMemberHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		helper.exampleHouseholdInvitation.FromUser.TwoFactorSecret = ""
		helper.exampleHouseholdInvitation.DestinationHousehold.WebhookEncryptionKey = ""
		assert.Equal(t, actual.Data, helper.exampleHouseholdInvitation)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, udm, sg, dbManager, dataChangesPublisher)
	})

	T.Run("with error publishing email request", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := fakes.BuildFakeHouseholdInvitationCreationRequestInput()
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		udm := &mocktypes.UserDataManagerMock{}
		udm.On(
			"GetUserByEmail",
			testutils.ContextMatcher,
			strings.TrimSpace(strings.ToLower(exampleInput.ToEmail)),
		).Return(helper.exampleUser, nil)
		helper.service.userDataManager = udm

		sg := &randommock.Generator{}
		sg.On(
			"GenerateBase64EncodedString",
			testutils.ContextMatcher,
			64,
		).Return(t.Name(), nil)
		helper.service.secretGenerator = sg

		dbManager := database.NewMockDatabase()
		dbManager.HouseholdInvitationDataManagerMock.On(
			"CreateHouseholdInvitation",
			testutils.ContextMatcher,
			mock.MatchedBy(func(*types.HouseholdInvitationDatabaseCreationInput) bool { return true }),
		).Return(helper.exampleHouseholdInvitation, nil)
		helper.service.householdInvitationDataManager = dbManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.InviteMemberHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusCreated, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		helper.exampleHouseholdInvitation.FromUser.TwoFactorSecret = ""
		helper.exampleHouseholdInvitation.DestinationHousehold.WebhookEncryptionKey = ""
		assert.Equal(t, actual.Data, helper.exampleHouseholdInvitation)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, udm, sg, dbManager, dataChangesPublisher)
	})
}

func Test_service_ReadHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)

		wd := &mocktypes.HouseholdInvitationDataManagerMock{}
		wd.On(
			"GetHouseholdInvitationByHouseholdAndID",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			helper.exampleHouseholdInvitation.ID,
		).Return(helper.exampleHouseholdInvitation, nil)
		helper.service.householdInvitationDataManager = wd

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		helper.exampleHouseholdInvitation.FromUser.TwoFactorSecret = ""
		helper.exampleHouseholdInvitation.DestinationHousehold.WebhookEncryptionKey = ""
		assert.Equal(t, actual.Data, helper.exampleHouseholdInvitation)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, wd)
	})

	T.Run("with error retrieving session context data", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.sessionContextDataFetcher = testutils.BrokenSessionContextDataFetcher

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no such household invitation in database", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)

		wd := &mocktypes.HouseholdInvitationDataManagerMock{}
		wd.On(
			"GetHouseholdInvitationByHouseholdAndID",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			helper.exampleHouseholdInvitation.ID,
		).Return((*types.HouseholdInvitation)(nil), sql.ErrNoRows)
		helper.service.householdInvitationDataManager = wd

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, wd)
	})

	T.Run("with error fetching household invitation from database", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)

		wd := &mocktypes.HouseholdInvitationDataManagerMock{}
		wd.On(
			"GetHouseholdInvitationByHouseholdAndID",
			testutils.ContextMatcher,
			helper.exampleHousehold.ID,
			helper.exampleHouseholdInvitation.ID,
		).Return((*types.HouseholdInvitation)(nil), errors.New("blah"))
		helper.service.householdInvitationDataManager = wd

		helper.service.ReadHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, wd)
	})
}

func Test_service_InboundInvitesHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		exampleHouseholdInvitations := fakes.BuildFakeHouseholdInvitationList()

		wd := &mocktypes.HouseholdInvitationDataManagerMock{}
		wd.On(
			"GetPendingHouseholdInvitationsForUser",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleHouseholdInvitations, nil)
		helper.service.householdInvitationDataManager = wd

		helper.service.InboundInvitesHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.NotEmpty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, wd)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.sessionContextDataFetcher = func(*http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		helper.service.InboundInvitesHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error reading from database", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)

		wd := &mocktypes.HouseholdInvitationDataManagerMock{}
		wd.On(
			"GetPendingHouseholdInvitationsForUser",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.HouseholdInvitation])(nil), errors.New("blah"))
		helper.service.householdInvitationDataManager = wd

		helper.service.InboundInvitesHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, wd)
	})
}

func Test_service_OutboundInvitesHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		exampleHouseholdInvitations := fakes.BuildFakeHouseholdInvitationList()

		wd := &mocktypes.HouseholdInvitationDataManagerMock{}
		wd.On(
			"GetPendingHouseholdInvitationsFromUser",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return(exampleHouseholdInvitations, nil)
		helper.service.householdInvitationDataManager = wd

		helper.service.OutboundInvitesHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)
		var actual *types.APIResponse[[]*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.NotEmpty(t, actual.Data)
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, wd)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.sessionContextDataFetcher = func(*http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		helper.service.OutboundInvitesHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error reading from database", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)

		wd := &mocktypes.HouseholdInvitationDataManagerMock{}
		wd.On(
			"GetPendingHouseholdInvitationsFromUser",
			testutils.ContextMatcher,
			helper.exampleUser.ID,
			mock.IsType(&types.QueryFilter{}),
		).Return((*types.QueryFilteredResult[types.HouseholdInvitation])(nil), errors.New("blah"))
		helper.service.householdInvitationDataManager = wd

		helper.service.OutboundInvitesHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, wd)
	})
}

func Test_service_AcceptInviteHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManagerMock{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return(helper.exampleHouseholdInvitation, nil)

		dataManager.On(
			"AcceptHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.ID,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.Note,
		).Return(nil)
		helper.service.householdInvitationDataManager = dataManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.AcceptInviteHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusAccepted, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, dataManager, dataChangesPublisher)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.sessionContextDataFetcher = func(*http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		helper.service.AcceptInviteHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader([]byte("{ bad json lol")))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.AcceptInviteHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := &types.HouseholdUpdateRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.AcceptInviteHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no matching invitation", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManagerMock{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return((*types.HouseholdInvitation)(nil), sql.ErrNoRows)
		helper.service.householdInvitationDataManager = dataManager

		helper.service.AcceptInviteHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dataManager)
	})

	T.Run("with error reading invitation", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManagerMock{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return((*types.HouseholdInvitation)(nil), errors.New("blah"))
		helper.service.householdInvitationDataManager = dataManager

		helper.service.AcceptInviteHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dataManager)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManagerMock{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return(helper.exampleHouseholdInvitation, nil)

		dataManager.On(
			"AcceptHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.ID,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.Note,
		).Return(errors.New("blah"))
		helper.service.householdInvitationDataManager = dataManager

		helper.service.AcceptInviteHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dataManager)
	})

	T.Run("with error publishing service event", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManagerMock{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return(helper.exampleHouseholdInvitation, nil)

		dataManager.On(
			"AcceptHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.ID,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.Note,
		).Return(nil)
		helper.service.householdInvitationDataManager = dataManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.AcceptInviteHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusAccepted, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, dataManager, dataChangesPublisher)
	})
}

func Test_service_CancelInviteHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManagerMock{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return(helper.exampleHouseholdInvitation, nil)

		dataManager.On(
			"CancelHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.ID,
			helper.exampleHouseholdInvitation.Note,
		).Return(nil)
		helper.service.householdInvitationDataManager = dataManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CancelInviteHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusAccepted, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, dataManager, dataChangesPublisher)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := &types.HouseholdUpdateRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CancelInviteHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with unparseable input", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader([]byte("{ bad json lol")))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.CancelInviteHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManagerMock{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return((*types.HouseholdInvitation)(nil), sql.ErrNoRows)
		helper.service.householdInvitationDataManager = dataManager

		helper.service.CancelInviteHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dataManager)
	})

	T.Run("with error fetching invitation", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManagerMock{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return((*types.HouseholdInvitation)(nil), errors.New("blah"))
		helper.service.householdInvitationDataManager = dataManager

		helper.service.CancelInviteHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dataManager)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.sessionContextDataFetcher = func(*http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		helper.service.CancelInviteHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManagerMock{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return(helper.exampleHouseholdInvitation, nil)

		dataManager.On(
			"CancelHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.ID,
			helper.exampleHouseholdInvitation.Note,
		).Return(errors.New("blah"))
		helper.service.householdInvitationDataManager = dataManager

		helper.service.CancelInviteHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dataManager)
	})

	T.Run("with error publishing service event", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManagerMock{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return(helper.exampleHouseholdInvitation, nil)

		dataManager.On(
			"CancelHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.ID,
			helper.exampleHouseholdInvitation.Note,
		).Return(nil)
		helper.service.householdInvitationDataManager = dataManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.CancelInviteHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusAccepted, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, dataManager, dataChangesPublisher)
	})
}

func Test_service_RejectInviteHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManagerMock{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return(helper.exampleHouseholdInvitation, nil)

		dataManager.On(
			"RejectHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.ID,
			helper.exampleHouseholdInvitation.Note,
		).Return(nil)
		helper.service.householdInvitationDataManager = dataManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(nil)
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.RejectInviteHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusAccepted, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, dataManager, dataChangesPublisher)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := &types.HouseholdInvitationUpdateRequestInput{}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.RejectInviteHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManagerMock{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return((*types.HouseholdInvitation)(nil), sql.ErrNoRows)
		helper.service.householdInvitationDataManager = dataManager

		helper.service.RejectInviteHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusNotFound, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dataManager)
	})

	T.Run("with error fetching invitation", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManagerMock{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return((*types.HouseholdInvitation)(nil), errors.New("blah"))

		helper.service.householdInvitationDataManager = dataManager

		helper.service.RejectInviteHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dataManager)
	})

	T.Run("with error fetching session context data", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.sessionContextDataFetcher = func(*http.Request) (*types.SessionContextData, error) {
			return nil, errors.New("blah")
		}

		helper.service.RejectInviteHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusUnauthorized, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader([]byte("{ bad json lol")))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		helper.service.RejectInviteHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusBadRequest, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManagerMock{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return(helper.exampleHouseholdInvitation, nil)

		dataManager.On(
			"RejectHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.ID,
			helper.exampleHouseholdInvitation.Note,
		).Return(errors.New("blah"))
		helper.service.householdInvitationDataManager = dataManager

		helper.service.RejectInviteHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusInternalServerError, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.Empty(t, actual.Data)
		assert.Error(t, actual.Error)

		mock.AssertExpectationsForObjects(t, dataManager)
	})

	T.Run("with error publishing service event", func(t *testing.T) {
		t.Parallel()

		helper := newTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		exampleInput := converters.ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(helper.exampleHouseholdInvitation)
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		var err error
		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(jsonBytes))
		require.NoError(t, err)
		require.NotNil(t, helper.req)

		dataManager := &mocktypes.HouseholdInvitationDataManagerMock{}
		dataManager.On(
			"GetHouseholdInvitationByTokenAndID",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.Token,
			helper.exampleHouseholdInvitation.ID,
		).Return(helper.exampleHouseholdInvitation, nil)

		dataManager.On(
			"RejectHouseholdInvitation",
			testutils.ContextMatcher,
			helper.exampleHouseholdInvitation.ID,
			helper.exampleHouseholdInvitation.Note,
		).Return(nil)
		helper.service.householdInvitationDataManager = dataManager

		dataChangesPublisher := &mockpublishers.Publisher{}
		dataChangesPublisher.On(
			"Publish",
			testutils.ContextMatcher,
			testutils.DataChangeMessageMatcher,
		).Return(errors.New("blah"))
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.RejectInviteHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusAccepted, helper.res.Code)
		var actual *types.APIResponse[*types.HouseholdInvitation]
		require.NoError(t, helper.service.encoderDecoder.DecodeBytes(helper.ctx, helper.res.Body.Bytes(), &actual))
		assert.NoError(t, actual.Error.AsError())

		mock.AssertExpectationsForObjects(t, dataManager, dataChangesPublisher)
	})
}
