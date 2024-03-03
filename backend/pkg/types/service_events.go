package types

const (
	// FinalizeMealPlansWithExpiredVotingPeriodsChoreType asks the worker to finalize meal plans with expired voting periods.
	FinalizeMealPlansWithExpiredVotingPeriodsChoreType = "finalize_meal_plans_with_expired_voting_periods"
	// CreateMealPlanTasksChoreType asks the worker to finalize meal plans with expired voting periods.
	CreateMealPlanTasksChoreType = "create_meal_plan_tasks"
)

type (
	// ServiceEventType enumerates valid service event types.
	ServiceEventType string

	// DataChangeMessage represents an event that asks a worker to write data to the datastore.
	DataChangeMessage struct {
		_ struct{} `json:"-"`

		HouseholdInvitation              *HouseholdInvitation            `json:"householdInvitation,omitempty"`
		ValidMeasurementUnitConversion   *ValidMeasurementUnitConversion `json:"validMeasurementUnitConversion,omitempty"`
		ValidInstrument                  *ValidInstrument                `json:"validInstrument,omitempty"`
		ValidIngredient                  *ValidIngredient                `json:"validIngredient,omitempty"`
		ValidIngredientGroup             *ValidIngredientGroup           `json:"validIngredientGroup,omitempty"`
		ValidPreparation                 *ValidPreparation               `json:"validPreparation,omitempty"`
		ValidIngredientState             *ValidIngredientState           `json:"validIngredientState,omitempty"`
		MealPlanGroceryListItem          *MealPlanGroceryListItem        `json:"mealPlanGroceryListItem,omitempty"`
		Meal                             *Meal                           `json:"meal,omitempty"`
		Context                          map[string]any                  `json:"context,omitempty"`
		Recipe                           *Recipe                         `json:"recipe,omitempty"`
		RecipePrepTask                   *RecipePrepTask                 `json:"recipePrepTask,omitempty"`
		RecipePrepTaskStep               *RecipePrepTaskStep             `json:"recipePrepTaskStep,omitempty"`
		RecipeStep                       *RecipeStep                     `json:"recipeStep,omitempty"`
		RecipeStepProduct                *RecipeStepProduct              `json:"recipeStepProduct,omitempty"`
		RecipeStepInstrument             *RecipeStepInstrument           `json:"recipeStepInstrument,omitempty"`
		RecipeStepIngredient             *RecipeStepIngredient           `json:"recipeStepIngredient,omitempty"`
		MealPlan                         *MealPlan                       `json:"mealPlan,omitempty"`
		MealPlanTask                     *MealPlanTask                   `json:"mealPlanTask,omitempty"`
		MealPlanEvent                    *MealPlanEvent                  `json:"mealPlanEvent,omitempty"`
		Household                        *Household                      `json:"household,omitempty"`
		MealPlanOption                   *MealPlanOption                 `json:"mealPlanOption,omitempty"`
		ValidIngredientMeasurementUnit   *ValidIngredientMeasurementUnit `json:"validIngredientMeasurementUnit,omitempty"`
		MealPlanOptionVote               *MealPlanOptionVote             `json:"mealPlanOptionVote,omitempty"`
		ValidPreparationInstrument       *ValidPreparationInstrument     `json:"validPreparationInstrument,omitempty"`
		Webhook                          *Webhook                        `json:"webhook,omitempty"`
		ValidIngredientPreparation       *ValidIngredientPreparation     `json:"validIngredientPreparation,omitempty"`
		ValidMeasurementUnit             *ValidMeasurementUnit           `json:"validMeasurementUnit,omitempty"`
		UserMembership                   *HouseholdUserMembership        `json:"userMembership,omitempty"`
		RecipeStepCompletionCondition    *RecipeStepCompletionCondition  `json:"recipeStepCompletionCondition,omitempty"`
		RecipeStepVessel                 *RecipeStepVessel               `json:"recipeStepVessel,omitempty"`
		PasswordResetToken               *PasswordResetToken             `json:"passwordResetToken,omitempty"`
		ValidIngredientStateIngredient   *ValidIngredientStateIngredient `json:"validIngredientStateIngredient,omitempty"`
		ServiceSetting                   *ServiceSetting                 `json:"serviceSetting,omitempty"`
		ServiceSettingConfiguration      *ServiceSettingConfiguration    `json:"serviceSettingConfiguration,omitempty"`
		HouseholdInstrumentOwnership     *HouseholdInstrumentOwnership   `json:"householdInstrumentOwnership,omitempty"`
		RecipeRating                     *RecipeRating                   `json:"recipeRating,omitempty"`
		ValidVessel                      *ValidVessel                    `json:"validVessel,omitempty"`
		ValidPreparationVessel           *ValidPreparationVessel         `json:"validPreparationVessel,omitempty"`
		UserNotification                 *UserNotification               `json:"userNotification,omitempty"`
		UserNotificationID               string                          `json:"userNotificationID"`
		RecipeStepVesselID               string                          `json:"recipeStepVesselID,omitempty"`
		HouseholdInvitationID            string                          `json:"householdInvitationID,omitempty"`
		UserID                           string                          `json:"userID"`
		HouseholdID                      string                          `json:"householdID,omitempty"`
		ValidMeasurementUnitID           string                          `json:"validMeasurementUnitID,omitempty"`
		ValidPreparationInstrumentID     string                          `json:"validPreparationInstrumentID,omitempty"`
		MealPlanOptionVoteID             string                          `json:"mealPlanOptionVoteID,omitempty"`
		ValidIngredientMeasurementUnitID string                          `json:"validIngredientMeasurementUnitID,omitempty"`
		MealPlanOptionID                 string                          `json:"mealPlanOptionID,omitempty"`
		MealPlanID                       string                          `json:"mealPlanID,omitempty"`
		MealPlanTaskID                   string                          `json:"mealPlanTaskID,omitempty"`
		RecipeStepID                     string                          `json:"recipeStepID,omitempty"`
		RecipePrepTaskID                 string                          `json:"recipePrepTaskID,omitempty"`
		RecipeID                         string                          `json:"recipeID,omitempty"`
		RecipeMediaID                    string                          `json:"recipeMediaID,omitempty"`
		MealID                           string                          `json:"mealID,omitempty"`
		MealPlanGroceryListItemID        string                          `json:"mealPlanGroceryListItemID,omitempty"`
		EventType                        ServiceEventType                `json:"messageType"`
		ValidIngredientStateIngredientID string                          `json:"validIngredientStateIngredientID"`
		ValidMeasurementUnitConversionID string                          `json:"validMeasurementUnitConversionID,omitempty"`
		ValidIngredientStateID           string                          `json:"validIngredientStateID,omitempty"`
		ValidIngredientGroupID           string                          `json:"validIngredientGroupID,omitempty"`
		MealPlanEventID                  string                          `json:"mealPlanEventID,omitempty"`
		EmailVerificationToken           string                          `json:"emailVerificationToken,omitempty"`
		UserIngredientPreferenceID       string                          `json:"userIngredientPreferenceID,omitempty"`
		HouseholdInstrumentOwnershipID   string                          `json:"householdInstrumentOwnershipID,omitempty"`
		RecipeRatingID                   string                          `json:"recipeRatingID,omitempty"`
		OAuth2ClientID                   string                          `json:"oauth2ClientID,omitempty"`
		UserIngredientPreferences        []*UserIngredientPreference     `json:"userIngredientPreference,omitempty"`
	}

	// ChoreMessage represents an event that asks a worker to perform a chore.
	ChoreMessage struct {
		_ struct{} `json:"-"`

		ChoreType string `json:"choreType"`
	}
)
