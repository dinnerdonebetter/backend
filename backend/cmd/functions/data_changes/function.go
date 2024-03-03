package datachangesfunction

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/dinnerdonebetter/backend/internal/analytics"
	analyticsconfig "github.com/dinnerdonebetter/backend/internal/analytics/config"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	"github.com/dinnerdonebetter/backend/internal/email"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/observability/logging/config"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/search"
	"github.com/dinnerdonebetter/backend/internal/search/indexing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	_ "github.com/KimMachineGun/automemlimit"
	"github.com/cloudevents/sdk-go/v2/event"
	"go.opentelemetry.io/otel"
	_ "go.uber.org/automaxprocs"
)

func init() {
	// Register a CloudEvent function with the Functions Framework
	functions.CloudEvent("ProcessDataChange", ProcessDataChange)
}

var (
	errRequiredDataIsNil = errors.New("required data is nil")

	nonWebhookEventTypes = []types.ServiceEventType{
		types.UserSignedUpCustomerEventType,
		types.UserArchivedCustomerEventType,
		types.TwoFactorSecretVerifiedCustomerEventType,
		types.TwoFactorDeactivatedCustomerEventType,
		types.TwoFactorSecretChangedCustomerEventType,
		types.PasswordResetTokenCreatedEventType,
		types.PasswordResetTokenRedeemedEventType,
		types.PasswordChangedEventType,
		types.EmailAddressChangedEventType,
		types.UsernameChangedEventType,
		types.UserDetailsChangedEventType,
		types.UsernameReminderRequestedEventType,
		types.UserLoggedInCustomerEventType,
		types.UserLoggedOutCustomerEventType,
		types.UserChangedActiveHouseholdCustomerEventType,
		types.UserEmailAddressVerifiedEventType,
		types.UserEmailAddressVerificationEmailRequestedEventType,
		types.HouseholdMemberRemovedCustomerEventType,
		types.HouseholdMembershipPermissionsUpdatedCustomerEventType,
		types.HouseholdOwnershipTransferredCustomerEventType,
		types.OAuth2ClientCreatedCustomerEventType,
		types.OAuth2ClientArchivedCustomerEventType,
	}
)

// MessagePublishedData contains the full Pub/Sub message
// See the documentation for more details:
// https://cloud.google.com/eventarc/docs/cloudevents#pubsub
type MessagePublishedData struct {
	Message PubSubMessage
}

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// ProcessDataChange handles a data change.
func ProcessDataChange(ctx context.Context, e event.Event) error {
	if strings.TrimSpace(strings.ToLower(os.Getenv("CEASE_OPERATION"))) == "true" {
		slog.Info("CEASE_OPERATION is set to true, exiting")
		return nil
	}

	logger := (&loggingcfg.Config{Level: logging.DebugLevel, Provider: loggingcfg.ProviderSlog}).ProvideLogger()

	var msg MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %w", err)
	}

	cfg, err := config.GetDataChangesWorkerConfigFromGoogleCloudSecretManager(ctx)
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}

	tracerProvider, initializeTracerErr := cfg.Observability.Tracing.ProvideTracerProvider(ctx, logger)
	if initializeTracerErr != nil {
		logger.Error(initializeTracerErr, "initializing tracer")
	}
	otel.SetTracerProvider(tracerProvider)

	tracer := tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("data_changes_job"))

	ctx, span := tracer.StartSpan(ctx)
	defer span.End()

	analyticsEventReporter, err := analyticsconfig.ProvideEventReporter(&cfg.Analytics, logger, tracerProvider)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "setting up customer data collector")
	}

	defer analyticsEventReporter.Close()

	publisherProvider, err := msgconfig.ProvidePublisherProvider(ctx, logger, tracerProvider, &cfg.Events)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "configuring queue manager")
	}

	defer publisherProvider.Close()

	outboundEmailsPublisher, err := publisherProvider.ProvidePublisher(os.Getenv("OUTBOUND_EMAILS_TOPIC_NAME"))
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "configuring outbound emails publisher")
	}

	defer outboundEmailsPublisher.Stop()

	searchDataIndexPublisher, err := publisherProvider.ProvidePublisher(os.Getenv("SEARCH_INDEXING_TOPIC_NAME"))
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "configuring search indexing publisher")
	}

	defer searchDataIndexPublisher.Stop()

	webhookExecutionRequestPublisher, err := publisherProvider.ProvidePublisher(os.Getenv("WEBHOOK_EXECUTION_REQUESTS_TOPIC_NAME"))
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "configuring search indexing publisher")
	}

	defer webhookExecutionRequestPublisher.Stop()

	// manual db timeout until I find out what's wrong
	dbConnectionContext, cancel := context.WithTimeout(ctx, 15*time.Second)
	dataManager, err := postgres.ProvideDatabaseClient(dbConnectionContext, logger, tracerProvider, &cfg.Database)
	if err != nil {
		cancel()
		return observability.PrepareAndLogError(err, logger, span, "establishing database connection")
	}

	cancel()
	defer dataManager.Close()

	var changeMessage types.DataChangeMessage
	if err = json.Unmarshal(msg.Message.Data, &changeMessage); err != nil {
		logger = logger.WithValue("raw_data", msg.Message.Data)
		return observability.PrepareAndLogError(err, logger, span, "unmarshalling data change message")
	}

	logger = logger.WithValue("event_type", changeMessage.EventType)

	if changeMessage.UserID != "" && changeMessage.EventType != "" {
		if err = analyticsEventReporter.EventOccurred(ctx, changeMessage.EventType, changeMessage.UserID, changeMessage.Context); err != nil {
			observability.AcknowledgeError(err, logger, span, "notifying customer data platform")
		}
	}

	var wg sync.WaitGroup

	go func() {
		wg.Add(1)
		if changeMessage.HouseholdID != "" && !slices.Contains(nonWebhookEventTypes, changeMessage.EventType) {
			var relevantWebhooks []*types.Webhook
			relevantWebhooks, err = dataManager.GetWebhooksForHouseholdAndEvent(ctx, changeMessage.HouseholdID, changeMessage.EventType)
			if err != nil {
				observability.AcknowledgeError(err, logger, span, "getting webhooks")
			}

			for _, webhook := range relevantWebhooks {
				if err = webhookExecutionRequestPublisher.Publish(ctx, &types.WebhookExecutionRequest{
					WebhookID:   webhook.ID,
					HouseholdID: changeMessage.HouseholdID,
					Payload:     changeMessage,
				}); err != nil {
					observability.AcknowledgeError(err, logger, span, "publishing webhook execution request")
				}
			}
		}
		wg.Done()
	}()

	go func() {
		wg.Add(1)
		if err = handleOutboundNotifications(ctx, logger, tracer, dataManager, outboundEmailsPublisher, webhookExecutionRequestPublisher, analyticsEventReporter, &changeMessage); err != nil {
			observability.AcknowledgeError(err, logger, span, "notifying customer(s)")
		}
		wg.Done()
	}()

	go func() {
		wg.Add(1)
		if err = handleSearchIndexUpdates(ctx, logger, tracer, searchDataIndexPublisher, &changeMessage); err != nil {
			observability.AcknowledgeError(err, logger, span, "updating search index)")
		}
		wg.Done()
	}()

	wg.Wait()

	return nil
}

func handleSearchIndexUpdates(
	ctx context.Context,
	l logging.Logger,
	tracer tracing.Tracer,
	searchDataIndexPublisher messagequeue.Publisher,
	changeMessage *types.DataChangeMessage,
) error {
	ctx, span := tracer.StartSpan(ctx)
	defer span.End()

	logger := l.WithValue("event_type", changeMessage.EventType)

	switch changeMessage.EventType {
	case types.UserSignedUpCustomerEventType,
		types.UserArchivedCustomerEventType,
		types.EmailAddressChangedEventType,
		types.UsernameChangedEventType,
		types.UserDetailsChangedEventType,
		types.UserEmailAddressVerifiedEventType:
		if changeMessage.UserID == "" {
			observability.AcknowledgeError(errRequiredDataIsNil, logger, span, "updating search index for User")
		}

		if err := searchDataIndexPublisher.Publish(ctx, &indexing.IndexRequest{
			RowID:     changeMessage.UserID,
			IndexType: search.IndexTypeUsers,
			Delete:    changeMessage.EventType == types.UserArchivedCustomerEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case types.RecipeCreatedCustomerEventType,
		types.RecipeUpdatedCustomerEventType,
		types.RecipeArchivedCustomerEventType:
		if changeMessage.Recipe == nil {
			observability.AcknowledgeError(errRequiredDataIsNil, logger, span, "updating search index for Recipe")
		}

		if err := searchDataIndexPublisher.Publish(ctx, &indexing.IndexRequest{
			RowID:     changeMessage.Recipe.ID,
			IndexType: search.IndexTypeRecipes,
			Delete:    changeMessage.EventType == types.RecipeArchivedCustomerEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case types.MealCreatedCustomerEventType,
		types.MealUpdatedCustomerEventType,
		types.MealArchivedCustomerEventType:
		if changeMessage.Meal == nil {
			observability.AcknowledgeError(errRequiredDataIsNil, logger, span, "updating search index for Meal")
		}

		if err := searchDataIndexPublisher.Publish(ctx, &indexing.IndexRequest{
			RowID:     changeMessage.Meal.ID,
			IndexType: search.IndexTypeRecipes,
			Delete:    changeMessage.EventType == types.MealArchivedCustomerEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case types.ValidIngredientCreatedCustomerEventType,
		types.ValidIngredientUpdatedCustomerEventType,
		types.ValidIngredientArchivedCustomerEventType:
		if changeMessage.ValidIngredient == nil {
			observability.AcknowledgeError(errRequiredDataIsNil, logger, span, "updating search index for ValidIngredient")
		}

		if err := searchDataIndexPublisher.Publish(ctx, &indexing.IndexRequest{
			RowID:     changeMessage.ValidIngredient.ID,
			IndexType: search.IndexTypeRecipes,
			Delete:    changeMessage.EventType == types.ValidIngredientArchivedCustomerEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case types.ValidInstrumentCreatedCustomerEventType,
		types.ValidInstrumentUpdatedCustomerEventType,
		types.ValidInstrumentArchivedCustomerEventType:
		if changeMessage.ValidInstrument == nil {
			observability.AcknowledgeError(errRequiredDataIsNil, logger, span, "updating search index for ValidInstrument")
		}

		if err := searchDataIndexPublisher.Publish(ctx, &indexing.IndexRequest{
			RowID:     changeMessage.ValidInstrument.ID,
			IndexType: search.IndexTypeRecipes,
			Delete:    changeMessage.EventType == types.ValidInstrumentArchivedCustomerEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case types.ValidMeasurementUnitCreatedCustomerEventType,
		types.ValidMeasurementUnitUpdatedCustomerEventType,
		types.ValidMeasurementUnitArchivedCustomerEventType:
		if changeMessage.ValidMeasurementUnit == nil {
			observability.AcknowledgeError(errRequiredDataIsNil, logger, span, "updating search index for ValidMeasurementUnit")
		}

		if err := searchDataIndexPublisher.Publish(ctx, &indexing.IndexRequest{
			RowID:     changeMessage.ValidMeasurementUnit.ID,
			IndexType: search.IndexTypeRecipes,
			Delete:    changeMessage.EventType == types.ValidMeasurementUnitArchivedCustomerEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case types.ValidPreparationCreatedCustomerEventType,
		types.ValidPreparationUpdatedCustomerEventType,
		types.ValidPreparationArchivedCustomerEventType:
		if changeMessage.ValidPreparation == nil {
			observability.AcknowledgeError(errRequiredDataIsNil, logger, span, "updating search index for ValidPreparation")
		}

		if err := searchDataIndexPublisher.Publish(ctx, &indexing.IndexRequest{
			RowID:     changeMessage.ValidPreparation.ID,
			IndexType: search.IndexTypeRecipes,
			Delete:    changeMessage.EventType == types.ValidPreparationArchivedCustomerEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case types.ValidIngredientStateCreatedCustomerEventType,
		types.ValidIngredientStateUpdatedCustomerEventType,
		types.ValidIngredientStateArchivedCustomerEventType:
		if changeMessage.ValidIngredientState == nil {
			observability.AcknowledgeError(errRequiredDataIsNil, logger, span, "updating search index for ValidIngredientState")
		}

		if err := searchDataIndexPublisher.Publish(ctx, &indexing.IndexRequest{
			RowID:     changeMessage.ValidIngredientState.ID,
			IndexType: search.IndexTypeRecipes,
			Delete:    changeMessage.EventType == types.ValidIngredientStateArchivedCustomerEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case types.ValidIngredientMeasurementUnitCreatedCustomerEventType,
		types.ValidIngredientMeasurementUnitUpdatedCustomerEventType,
		types.ValidIngredientMeasurementUnitArchivedCustomerEventType:
		if changeMessage.ValidIngredientMeasurementUnit == nil {
			observability.AcknowledgeError(errRequiredDataIsNil, logger, span, "updating search index for ValidIngredientMeasurementUnit")
		}

		if err := searchDataIndexPublisher.Publish(ctx, &indexing.IndexRequest{
			RowID:     changeMessage.ValidIngredientMeasurementUnit.ID,
			IndexType: search.IndexTypeRecipes,
			Delete:    changeMessage.EventType == types.ValidIngredientMeasurementUnitArchivedCustomerEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case types.ValidPreparationInstrumentCreatedCustomerEventType,
		types.ValidPreparationInstrumentUpdatedCustomerEventType,
		types.ValidPreparationInstrumentArchivedCustomerEventType:
		if changeMessage.ValidPreparationInstrument == nil {
			observability.AcknowledgeError(errRequiredDataIsNil, logger, span, "updating search index for ValidPreparationInstrument")
		}

		if err := searchDataIndexPublisher.Publish(ctx, &indexing.IndexRequest{
			RowID:     changeMessage.ValidPreparationInstrument.ID,
			IndexType: search.IndexTypeRecipes,
			Delete:    changeMessage.EventType == types.ValidPreparationInstrumentArchivedCustomerEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case types.ValidIngredientPreparationCreatedCustomerEventType,
		types.ValidIngredientPreparationUpdatedCustomerEventType,
		types.ValidIngredientPreparationArchivedCustomerEventType:
		if changeMessage.ValidIngredientPreparation == nil {
			observability.AcknowledgeError(errRequiredDataIsNil, logger, span, "updating search index for ValidIngredientPreparation")
		}

		if err := searchDataIndexPublisher.Publish(ctx, &indexing.IndexRequest{
			RowID:     changeMessage.ValidIngredientPreparation.ID,
			IndexType: search.IndexTypeRecipes,
			Delete:    changeMessage.EventType == types.ValidIngredientPreparationArchivedCustomerEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	default:
		logger.Debug("event type not handled for search indexing")
		return nil
	}
}

func handleOutboundNotifications(
	ctx context.Context,
	l logging.Logger,
	tracer tracing.Tracer,
	dataManager database.DataManager,
	outboundEmailsPublisher,
	webhookExecutionRequestsPublisher messagequeue.Publisher,
	analyticsEventReporter analytics.EventReporter,
	changeMessage *types.DataChangeMessage,
) error {
	ctx, span := tracer.StartSpan(ctx)
	defer span.End()

	var (
		emailType string
		edrs      []*email.DeliveryRequest
	)

	logger := l.WithValue("event_type", changeMessage.EventType)

	switch changeMessage.EventType {
	case types.UserSignedUpCustomerEventType:
		emailType = "user signup"
		if err := analyticsEventReporter.AddUser(ctx, changeMessage.UserID, changeMessage.Context); err != nil {
			observability.AcknowledgeError(err, logger, span, "notifying customer data platform")
		}

		edrs = append(edrs, &email.DeliveryRequest{
			UserID:                 changeMessage.UserID,
			Template:               email.TemplateTypeVerifyEmailAddress,
			EmailVerificationToken: changeMessage.EmailVerificationToken,
		})
	case types.UserEmailAddressVerificationEmailRequestedEventType:
		emailType = "email address verification"

		edrs = append(edrs, &email.DeliveryRequest{
			UserID:                 changeMessage.UserID,
			Template:               email.TemplateTypeVerifyEmailAddress,
			EmailVerificationToken: changeMessage.EmailVerificationToken,
		})
	case types.MealPlanCreatedCustomerEventType:
		emailType = "meal plan created"
		mealPlan := changeMessage.MealPlan
		if mealPlan == nil {
			return observability.PrepareError(fmt.Errorf("meal plan is nil"), span, "publishing meal plan created email")
		}

		household, err := dataManager.GetHousehold(ctx, mealPlan.BelongsToHousehold)
		if err != nil {
			return observability.PrepareError(err, span, "getting household")
		}

		for _, member := range household.Members {
			if member.BelongsToUser.EmailAddressVerifiedAt != nil {
				edrs = append(edrs, &email.DeliveryRequest{
					UserID:   member.BelongsToUser.ID,
					Template: email.TemplateTypeMealPlanCreated,
					MealPlan: mealPlan,
				})
			}
		}
	case types.PasswordResetTokenCreatedEventType:
		emailType = "password reset request"
		if changeMessage.PasswordResetToken == nil {
			return observability.PrepareError(fmt.Errorf("password reset token is nil"), span, "publishing password reset token email")
		}

		edrs = append(edrs, &email.DeliveryRequest{
			UserID:             changeMessage.UserID,
			Template:           email.TemplateTypePasswordResetTokenCreated,
			PasswordResetToken: changeMessage.PasswordResetToken,
		})

	case types.UsernameReminderRequestedEventType:
		emailType = "username reminder"
		edrs = append(edrs, &email.DeliveryRequest{
			UserID:   changeMessage.UserID,
			Template: email.TemplateTypeUsernameReminder,
		})

	case types.PasswordResetTokenRedeemedEventType:
		emailType = "password reset token redeemed"
		edrs = append(edrs, &email.DeliveryRequest{
			UserID:   changeMessage.UserID,
			Template: email.TemplateTypePasswordResetTokenRedeemed,
		})

	case types.PasswordChangedEventType:
		emailType = "password reset token redeemed"
		edrs = append(edrs, &email.DeliveryRequest{
			UserID:   changeMessage.UserID,
			Template: email.TemplateTypePasswordReset,
		})

	case types.HouseholdInvitationCreatedCustomerEventType:
		emailType = "household invitation created"
		if changeMessage.HouseholdInvitation == nil {
			return observability.PrepareError(fmt.Errorf("household invitation is nil"), span, "publishing password reset token redemption email")
		}

		edrs = append(edrs, &email.DeliveryRequest{
			UserID:     changeMessage.UserID,
			Template:   email.TemplateTypeInvite,
			Invitation: changeMessage.HouseholdInvitation,
		})
	}

	if len(edrs) == 0 {
		logger.WithValue("email_type", emailType).WithValue("outbound_emails_to_send", len(edrs)).Info("publishing email requests")
	}

	for _, edr := range edrs {
		if err := outboundEmailsPublisher.Publish(ctx, edr); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing %s request email", emailType)
		}
	}

	return nil
}
