package workers

import (
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// MealPlanFinalizationHandler finalizes a meal plan.
func (s *service) MealPlanFinalizationHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req).WithSpan(span)
	logger.Info("meal plan finalization worker invoked")

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	var request *types.FinalizeMealPlansRequest
	if err := s.encoderDecoder.DecodeRequest(ctx, req, &request); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request")
		errRes := types.NewAPIErrorResponse("invalid request content", types.ErrDecodingRequestInput, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusBadRequest)
		return
	}

	response := &types.FinalizeMealPlansResponse{}
	if request.ReturnCount {
		count, err := s.mealPlanFinalizationWorker.FinalizeExpiredMealPlans(ctx, nil)
		if err != nil {
			observability.AcknowledgeError(err, logger, span, "finalizing expired meal plans")
			errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
			return
		}

		response.Count = count
	} else {
		if err := s.mealPlanFinalizationWorker.FinalizeExpiredMealPlansWithoutReturningCount(ctx, nil); err != nil {
			observability.AcknowledgeError(err, logger, span, "finalizing expired meal plans")
			errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
			s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
			return
		}
	}

	logger.WithValue("finalized_count", response.Count).Info("meal plan finalization worker completed")

	responseValue := &types.APIResponse[*types.FinalizeMealPlansRequest]{
		Data:    request,
		Details: responseDetails,
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusAccepted)
}

// MealPlanGroceryListInitializationHandler initializes a grocery list for a given meal plan.
func (s *service) MealPlanGroceryListInitializationHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req).WithSpan(span)
	logger.Info("meal plan grocery list initialization worker invoked")

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	if err := s.mealPlanGroceryListInitializer.InitializeGroceryListsForFinalizedMealPlans(ctx, nil); err != nil {
		observability.AcknowledgeError(err, logger, span, "finalizing expired meal plans")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	logger.Info("meal plan grocery list initialization worker completed")

	responseValue := &types.APIResponse[any]{
		Details: responseDetails,
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusAccepted)
}

// MealPlanTaskCreationHandler creates tasks for a meal plan.
func (s *service) MealPlanTaskCreationHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req).WithSpan(span)
	logger.Info("meal plan task creation worker invoked")

	responseDetails := types.ResponseDetails{
		TraceID: span.SpanContext().TraceID().String(),
	}

	if err := s.mealPlanTaskCreatorWorker.CreateMealPlanTasksForFinalizedMealPlans(ctx, nil); err != nil {
		observability.AcknowledgeError(err, logger, span, "finalizing expired meal plans")
		errRes := types.NewAPIErrorResponse("database error", types.ErrTalkingToDatabase, responseDetails)
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, errRes, http.StatusInternalServerError)
		return
	}

	logger.Info("meal plan task creation worker completed")
	responseValue := &types.APIResponse[any]{
		Details: responseDetails,
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, responseValue, http.StatusAccepted)
}
