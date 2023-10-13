package types

import (
	"errors"
)

var (
	errNewPasswordSameAsOld                 = errors.New("new password cannot be the same as the old password")
	errStartsAfterItEnds                    = errors.New("invalid start and end dates")
	errOneMainMinimumRequired               = errors.New("at least one main required for meal creation")
	errInvalidVotingDeadline                = errors.New("invalid voting deadline")
	errAtLeastOneRatingRequired             = errors.New("recipe rating must have at least one rating")
	errOneInstrumentOrVesselRequired        = errors.New("at least one instrument or vessel is required")
	errInstrumentIDOrProductIndicesRequired = errors.New("either instrumentID or productOfRecipeStepIndex and productOfRecipeStepProductIndex must be set")
	errDefaultValueMustBeEnumerationValue   = errors.New("default value must be in enumeration")
	errMustBeEitherMetricOrImperial         = errors.New("cannot be both metric and imperial")
	errInvalidType                          = errors.New("unexpected type received")

	// this just ensures that we don't have any duplicated codes.
	_ = map[string]ErrorCode{
		string(ErrFetchingSessionContextData): ErrFetchingSessionContextData,
		string(ErrDecodingRequestInput):       ErrDecodingRequestInput,
		string(ErrValidatingRequestInput):     ErrValidatingRequestInput,
		string(ErrDataNotFound):               ErrDataNotFound,
		string(ErrTalkingToDatabase):          ErrTalkingToDatabase,
		string(ErrTalkingToSearchProvider):    ErrTalkingToSearchProvider,
	}
)

type (
	errorCode string
	ErrorCode errorCode
)

const (
	// ErrFetchingSessionContextData is returned when we fail to fetch session context data.
	ErrFetchingSessionContextData ErrorCode = "E101"
	// ErrDecodingRequestInput is returned when we fail to decode request input.
	ErrDecodingRequestInput ErrorCode = "E102"
	// ErrValidatingRequestInput is returned when the user provides invalid input.
	ErrValidatingRequestInput ErrorCode = "E103"
	// ErrDataNotFound is returned when we fail to find data in the database.
	ErrDataNotFound ErrorCode = "E104"
	// ErrTalkingToDatabase is returned when we fail to interact with a database.
	ErrTalkingToDatabase ErrorCode = "E105"
	// ErrTalkingToSearchProvider is returned when we fail to interact with a database.
	ErrTalkingToSearchProvider ErrorCode = "E106"
)
