package mealplanoptions

import (
	"errors"
	"net/http"
	"testing"

	mockencoding "github.com/dinnerdonebetter/backend/internal/encoding/mock"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	mockrouting "github.com/dinnerdonebetter/backend/internal/routing/mock"
	mealplaneventsservice "github.com/dinnerdonebetter/backend/internal/services/mealplanevents"
	mealplansservice "github.com/dinnerdonebetter/backend/internal/services/mealplans"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildTestService() *service {
	return &service{
		logger:                    logging.NewNoopLogger(),
		mealPlanOptionDataManager: &mocktypes.MealPlanOptionDataManager{},
		mealPlanOptionIDFetcher:   func(req *http.Request) string { return "" },
		encoderDecoder:            mockencoding.NewMockEncoderDecoder(),
		tracer:                    tracing.NewTracerForTest("test"),
	}
}

func TestProvideMealPlanOptionsService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		rpm := mockrouting.NewRouteParamManager()
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			mealplansservice.MealPlanIDURIParamKey,
		).Return(func(*http.Request) string { return "" })
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			mealplaneventsservice.MealPlanEventIDURIParamKey,
		).Return(func(*http.Request) string { return "" })
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			MealPlanOptionIDURIParamKey,
		).Return(func(*http.Request) string { return "" })

		cfg := &Config{
			DataChangesTopicName: "data_changes",
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProvidePublisher", cfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

		s, err := ProvideService(
			logging.NewNoopLogger(),
			cfg,
			&mocktypes.MealPlanOptionDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			rpm,
			pp,
			tracing.NewNoopTracerProvider(),
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, rpm, pp)
	})

	T.Run("with error providing data changes producer", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			DataChangesTopicName: "data_changes",
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProvidePublisher", cfg.DataChangesTopicName).Return((*mockpublishers.Publisher)(nil), errors.New("blah"))

		s, err := ProvideService(
			logging.NewNoopLogger(),
			cfg,
			&mocktypes.MealPlanOptionDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			nil,
			pp,
			tracing.NewNoopTracerProvider(),
		)

		assert.Nil(t, s)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, pp)
	})
}
