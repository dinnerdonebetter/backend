package codegen

import (
	"reflect"

	"github.com/dinnerdonebetter/backend/pkg/types"
)

type TypeDefinition struct {
	Type        any
	Description string
}

func (d *TypeDefinition) Name() string {
	n := reflect.Indirect(reflect.ValueOf(d.Type)).Type().Name()
	return n
}

var (
	CustomTypeMap = map[string]string{
		"ValidIngredientState.attributeType":                     "ValidIngredientStateAttributeType",
		"ValidIngredientStateCreationRequestInput.attributeType": "ValidIngredientStateAttributeType",
		"ValidIngredientStateUpdateRequestInput.attributeType":   "ValidIngredientStateAttributeType",
		"RecipeStepProduct.type":                                 "ValidRecipeStepProductType",
		"RecipeStepProductCreationRequestInput.type":             "ValidRecipeStepProductType",
		"RecipeStepProductUpdateRequestInput.type":               "ValidRecipeStepProductType",
		"MealPlanGroceryListItem.status":                         "ValidMealPlanGroceryListItemStatus",
		"MealPlanGroceryListItemCreationRequestInput.status":     "ValidMealPlanGroceryListItemStatus",
		"MealPlanGroceryListItemUpdateRequestInput.status":       "ValidMealPlanGroceryListItemStatus",
		"MealComponent.componentType":                            "MealComponentType",
		"MealComponentCreationRequestInput.componentType":        "MealComponentType",
		"MealComponentUpdateRequestInput.componentType":          "MealComponentType",
		"MealPlanTaskStatusChangeRequestInput.status":            "MealPlanTaskStatus",
		"MealPlanTask.status":                                    "MealPlanTaskStatus",
		"MealPlanTaskCreationRequestInput.status":                "MealPlanTaskStatus",
		"MealPlanTaskUpdateRequestInput.status":                  "MealPlanTaskStatus",
		"MealPlan.status":                                        "ValidMealPlanStatus",
		"MealPlanCreationRequestInput.status":                    "ValidMealPlanStatus",
		"MealPlanUpdateRequestInput.status":                      "ValidMealPlanStatus",
		"MealPlan.electionMethod":                                "ValidMealPlanElectionMethod",
		"MealPlanUpdateRequestInput.electionMethod":              "ValidMealPlanElectionMethod",
		"MealPlanCreationRequestInput.electionMethod":            "ValidMealPlanElectionMethod",
		"ValidVessel.shape":                                      "ValidVesselShapeType",
		"ValidVesselUpdateRequestInput.shape":                    "ValidVesselShapeType",
		"ValidVesselCreationRequestInput.shape":                  "ValidVesselShapeType",
	}

	DefaultEnumValues = map[string]string{
		"ValidMealPlanStatus":                "'awaiting_votes'",
		"ValidMealPlanGroceryListItemStatus": "'unknown'",
		"ValidMealPlanElectionMethod":        "'schulze'",
		"ValidIngredientStateAttributeType":  "'other'",
		"ValidRecipeStepProductType":         "'ingredient'",
		"MealComponentType":                  "'unspecified'",
		"MealPlanTaskStatus":                 "'unfinished'",
		"ValidVesselShapeType":               "'other'",
	}

	TypeDefinitionFilesToGenerate = map[string][]any{
		"admin": {
			types.ModifyUserPermissionsInput{},
		},
		"oauth2Clients": {
			types.OAuth2Client{},
			types.OAuth2ClientCreationRequestInput{},
			types.OAuth2ClientCreationResponse{},
		},
		"auth": {
			types.ChangeActiveHouseholdInput{},
			types.PasswordResetToken{},
			types.PasswordResetTokenCreationRequestInput{},
			types.PasswordResetTokenRedemptionRequestInput{},
			types.TOTPSecretRefreshInput{},
			types.TOTPSecretVerificationInput{},
			types.TOTPSecretRefreshResponse{},
			types.PasswordUpdateInput{},
		},
		"errors": {
			types.APIError{},
		},
		"householdInvitations": {
			types.HouseholdInvitation{},
			types.HouseholdInvitationUpdateRequestInput{},
			types.HouseholdInvitationCreationRequestInput{},
		},
		"households": {
			types.Household{},
			types.HouseholdCreationRequestInput{},
			types.HouseholdUpdateRequestInput{},
			types.HouseholdOwnershipTransferInput{},
		},
		"householdUserMemberships": {
			types.HouseholdUserMembership{},
			types.HouseholdUserMembershipWithUser{},
			types.HouseholdUserMembershipCreationRequestInput{},
		},
		"mealPlanEvents": {
			types.MealPlanEvent{},
			types.MealPlanEventCreationRequestInput{},
			types.MealPlanEventUpdateRequestInput{},
		},
		"mealPlanGroceryListItems": {
			types.MealPlanGroceryListItem{},
			types.MealPlanGroceryListItemCreationRequestInput{},
			types.MealPlanGroceryListItemUpdateRequestInput{},
		},
		"mealPlanOptions": {
			types.MealPlanOption{},
			types.MealPlanOptionCreationRequestInput{},
			types.MealPlanOptionUpdateRequestInput{},
		},
		"mealPlanOptionVotes": {
			types.MealPlanOptionVote{},
			types.MealPlanOptionVoteCreationInput{},
			types.MealPlanOptionVoteCreationRequestInput{},
			types.MealPlanOptionVoteUpdateRequestInput{},
		},
		"mealPlans": {
			types.MealPlan{},
			types.MealPlanCreationRequestInput{},
			types.MealPlanUpdateRequestInput{},
		},
		"mealPlanTasks": {
			types.MealPlanTask{},
			types.MealPlanTaskCreationRequestInput{},
			types.MealPlanTaskStatusChangeRequestInput{},
		},
		"meals": {
			types.Meal{},
			types.MealCreationRequestInput{},
			types.MealUpdateRequestInput{},
		},
		"recipeRatings": {
			types.RecipeRating{},
			types.RecipeRatingCreationRequestInput{},
			types.RecipeRatingUpdateRequestInput{},
		},
		"householdInstrumentOwnerships": {
			types.HouseholdInstrumentOwnership{},
			types.HouseholdInstrumentOwnershipCreationRequestInput{},
			types.HouseholdInstrumentOwnershipUpdateRequestInput{},
		},
		"mealComponents": {
			types.MealComponent{},
			types.MealComponentCreationRequestInput{},
			types.MealComponentUpdateRequestInput{},
		},
		"permissions": {
			types.UserPermissionsRequestInput{},
			types.UserPermissionsResponse{},
		},
		"recipeMedia": {
			types.RecipeMedia{},
			types.RecipeMediaCreationRequestInput{},
			types.RecipeMediaUpdateRequestInput{},
		},
		"recipePrepTasks": {
			types.RecipePrepTask{},
			types.RecipePrepTaskCreationRequestInput{},
			types.RecipePrepTaskWithinRecipeCreationRequestInput{},
			types.RecipePrepTaskUpdateRequestInput{},
		},
		"recipePrepTaskSteps": {
			types.RecipePrepTaskStep{},
			types.RecipePrepTaskStepWithinRecipeCreationRequestInput{},
			types.RecipePrepTaskStepCreationRequestInput{},
			types.RecipePrepTaskStepUpdateRequestInput{},
		},
		"recipeStepCompletionConditions": {
			types.RecipeStepCompletionCondition{},
			types.RecipeStepCompletionConditionIngredient{},
			types.RecipeStepCompletionConditionCreationRequestInput{},
			types.RecipeStepCompletionConditionIngredientCreationRequestInput{},
			types.RecipeStepCompletionConditionUpdateRequestInput{},
		},
		"recipeStepIngredients": {
			types.RecipeStepIngredient{},
			types.RecipeStepIngredientCreationRequestInput{},
			types.RecipeStepIngredientUpdateRequestInput{},
		},
		"recipeStepInstruments": {
			types.RecipeStepInstrument{},
			types.RecipeStepInstrumentCreationRequestInput{},
			types.RecipeStepInstrumentUpdateRequestInput{},
		},
		"recipeStepVessels": {
			types.RecipeStepVessel{},
			types.RecipeStepVesselCreationRequestInput{},
			types.RecipeStepVesselUpdateRequestInput{},
		},
		"recipeStepProducts": {
			types.RecipeStepProduct{},
			types.RecipeStepProductCreationRequestInput{},
			types.RecipeStepProductUpdateRequestInput{},
		},
		"recipeSteps": {
			types.RecipeStep{},
			types.RecipeStepCreationRequestInput{},
			types.RecipeStepUpdateRequestInput{},
		},
		"recipes": {
			types.Recipe{},
			types.RecipeCreationRequestInput{},
			types.RecipeUpdateRequestInput{},
		},
		"users": {
			types.UserStatusResponse{},
			types.User{},
			types.UserRegistrationInput{},
			types.UserCreationResponse{},
			types.UserLoginInput{},
			types.UsernameReminderRequestInput{},
			types.UserAccountStatusUpdateInput{},
			types.EmailAddressVerificationRequestInput{},
			types.AvatarUpdateInput{},
		},
		"userIngredientPreferences": {
			types.UserIngredientPreference{},
			types.UserIngredientPreferenceCreationRequestInput{},
			types.UserIngredientPreferenceUpdateRequestInput{},
		},
		"validIngredientMeasurementUnits": {
			types.ValidIngredientMeasurementUnit{},
			types.ValidIngredientMeasurementUnitCreationRequestInput{},
			types.ValidIngredientMeasurementUnitUpdateRequestInput{},
		},
		"validIngredientPreparations": {
			types.ValidIngredientPreparation{},
			types.ValidIngredientPreparationCreationRequestInput{},
			types.ValidIngredientPreparationUpdateRequestInput{},
		},
		"validIngredientStates": {
			types.ValidIngredientState{},
			types.ValidIngredientStateCreationRequestInput{},
			types.ValidIngredientStateUpdateRequestInput{},
		},
		"validIngredientStateIngredients": {
			types.ValidIngredientStateIngredient{},
			types.ValidIngredientStateIngredientCreationRequestInput{},
			types.ValidIngredientStateIngredientUpdateRequestInput{},
		},
		"validIngredients": {
			types.ValidIngredient{},
			types.ValidIngredientCreationRequestInput{},
			types.ValidIngredientUpdateRequestInput{},
		},
		"validIngredientGroups": {
			types.ValidIngredientGroup{},
			types.ValidIngredientGroupCreationRequestInput{},
			types.ValidIngredientGroupUpdateRequestInput{},
			types.ValidIngredientGroupMember{},
			types.ValidIngredientGroupMemberCreationRequestInput{},
		},
		"validInstruments": {
			types.ValidInstrument{},
			types.ValidInstrumentCreationRequestInput{},
			types.ValidInstrumentUpdateRequestInput{},
		},
		"validMeasurementUnitConversions": {
			types.ValidMeasurementUnitConversion{},
			types.ValidMeasurementUnitConversionCreationRequestInput{},
			types.ValidMeasurementUnitConversionUpdateRequestInput{},
		},
		"validMeasurementUnits": {
			types.ValidMeasurementUnit{},
			types.ValidMeasurementUnitCreationRequestInput{},
			types.ValidMeasurementUnitUpdateRequestInput{},
		},
		"validPreparationInstruments": {
			types.ValidPreparationInstrument{},
			types.ValidPreparationInstrumentCreationRequestInput{},
			types.ValidPreparationInstrumentUpdateRequestInput{},
		},
		"validPreparationVessels": {
			types.ValidPreparationVessel{},
			types.ValidPreparationVesselCreationRequestInput{},
			types.ValidPreparationVesselUpdateRequestInput{},
		},
		"validPreparations": {
			types.ValidPreparation{},
			types.ValidPreparationCreationRequestInput{},
			types.ValidPreparationUpdateRequestInput{},
		},
		"validVessels": {
			types.ValidVessel{},
			types.ValidVesselCreationRequestInput{},
			types.ValidVesselUpdateRequestInput{},
		},
		"serviceSetting": {
			types.ServiceSetting{},
			types.ServiceSettingCreationRequestInput{},
			types.ServiceSettingUpdateRequestInput{},
		},
		"serviceSettingConfiguration": {
			types.ServiceSettingConfiguration{},
			types.ServiceSettingConfigurationCreationRequestInput{},
			types.ServiceSettingConfigurationUpdateRequestInput{},
		},
		"webhooks": {
			types.Webhook{},
			types.WebhookTriggerEvent{},
			types.WebhookCreationRequestInput{},
		},
	}

	TypesWeCareAbout = []*TypeDefinition{
		{
			Type:        &types.ModifyUserPermissionsInput{},
			Description: "",
		},
		{
			Type:        &types.OAuth2Client{},
			Description: "",
		},
		{
			Type:        &types.OAuth2ClientCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.OAuth2ClientCreationResponse{},
			Description: "",
		},
		{
			Type:        &types.ChangeActiveHouseholdInput{},
			Description: "",
		},
		{
			Type:        &types.PasswordResetToken{},
			Description: "",
		},
		{
			Type:        &types.PasswordResetTokenCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.PasswordResetTokenRedemptionRequestInput{},
			Description: "",
		},
		{
			Type:        &types.TOTPSecretRefreshInput{},
			Description: "",
		},
		{
			Type:        &types.TOTPSecretVerificationInput{},
			Description: "",
		},
		{
			Type:        &types.TOTPSecretRefreshResponse{},
			Description: "",
		},
		{
			Type:        &types.PasswordUpdateInput{},
			Description: "",
		},
		{
			Type:        &types.APIError{},
			Description: "",
		},
		{
			Type:        &types.HouseholdInvitation{},
			Description: "",
		},
		{
			Type:        &types.HouseholdInvitationUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.HouseholdInvitationCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.Household{},
			Description: "",
		},
		{
			Type:        &types.HouseholdCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.HouseholdUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.HouseholdOwnershipTransferInput{},
			Description: "",
		},
		{
			Type:        &types.HouseholdUserMembership{},
			Description: "",
		},
		{
			Type:        &types.HouseholdUserMembershipWithUser{},
			Description: "",
		},
		{
			Type:        &types.HouseholdUserMembershipCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.MealPlanEvent{},
			Description: "",
		},
		{
			Type:        &types.MealPlanEventCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.MealPlanEventUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.MealPlanGroceryListItem{},
			Description: "",
		},
		{
			Type:        &types.MealPlanGroceryListItemCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.MealPlanGroceryListItemUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.MealPlanOption{},
			Description: "",
		},
		{
			Type:        &types.MealPlanOptionCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.MealPlanOptionUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.MealPlanOptionVote{},
			Description: "",
		},
		{
			Type:        &types.MealPlanOptionVoteCreationInput{},
			Description: "",
		},
		{
			Type:        &types.MealPlanOptionVoteCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.MealPlanOptionVoteUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.MealPlan{},
			Description: "",
		},
		{
			Type:        &types.MealPlanCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.MealPlanUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.MealPlanTask{},
			Description: "",
		},
		{
			Type:        &types.MealPlanTaskCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.MealPlanTaskStatusChangeRequestInput{},
			Description: "",
		},
		{
			Type:        &types.Meal{},
			Description: "",
		},
		{
			Type:        &types.MealCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.MealUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.RecipeRating{},
			Description: "",
		},
		{
			Type:        &types.RecipeRatingCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.RecipeRatingUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.HouseholdInstrumentOwnership{},
			Description: "",
		},
		{
			Type:        &types.HouseholdInstrumentOwnershipCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.HouseholdInstrumentOwnershipUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.MealComponent{},
			Description: "",
		},
		{
			Type:        &types.MealComponentCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.MealComponentUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.UserPermissionsRequestInput{},
			Description: "",
		},
		{
			Type:        &types.UserPermissionsResponse{},
			Description: "",
		},
		{
			Type:        &types.RecipeMedia{},
			Description: "",
		},
		{
			Type:        &types.RecipeMediaCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.RecipeMediaUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.RecipePrepTask{},
			Description: "",
		},
		{
			Type:        &types.RecipePrepTaskCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.RecipePrepTaskWithinRecipeCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.RecipePrepTaskUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.RecipePrepTaskStep{},
			Description: "",
		},
		{
			Type:        &types.RecipePrepTaskStepWithinRecipeCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.RecipePrepTaskStepCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.RecipePrepTaskStepUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.RecipeStepCompletionCondition{},
			Description: "",
		},
		{
			Type:        &types.RecipeStepCompletionConditionIngredient{},
			Description: "",
		},
		{
			Type:        &types.RecipeStepCompletionConditionCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.RecipeStepCompletionConditionIngredientCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.RecipeStepCompletionConditionUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.RecipeStepIngredient{},
			Description: "",
		},
		{
			Type:        &types.RecipeStepIngredientCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.RecipeStepIngredientUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.RecipeStepInstrument{},
			Description: "",
		},
		{
			Type:        &types.RecipeStepInstrumentCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.RecipeStepInstrumentUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.RecipeStepVessel{},
			Description: "",
		},
		{
			Type:        &types.RecipeStepVesselCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.RecipeStepVesselUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.RecipeStepProduct{},
			Description: "",
		},
		{
			Type:        &types.RecipeStepProductCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.RecipeStepProductUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.RecipeStep{},
			Description: "",
		},
		{
			Type:        &types.RecipeStepCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.RecipeStepUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.Recipe{},
			Description: "",
		},
		{
			Type:        &types.RecipeCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.RecipeUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.UserStatusResponse{},
			Description: "",
		},
		{
			Type:        &types.User{},
			Description: "",
		},
		{
			Type:        &types.UserRegistrationInput{},
			Description: "",
		},
		{
			Type:        &types.UserCreationResponse{},
			Description: "",
		},
		{
			Type:        &types.UserLoginInput{},
			Description: "",
		},
		{
			Type:        &types.UsernameReminderRequestInput{},
			Description: "",
		},
		{
			Type:        &types.UserAccountStatusUpdateInput{},
			Description: "",
		},
		{
			Type:        &types.EmailAddressVerificationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.AvatarUpdateInput{},
			Description: "",
		},
		{
			Type:        &types.UserIngredientPreference{},
			Description: "",
		},
		{
			Type:        &types.UserIngredientPreferenceCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.UserIngredientPreferenceUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ValidIngredientMeasurementUnit{},
			Description: "",
		},
		{
			Type:        &types.ValidIngredientMeasurementUnitCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ValidIngredientMeasurementUnitUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ValidIngredientPreparation{},
			Description: "",
		},
		{
			Type:        &types.ValidIngredientPreparationCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ValidIngredientPreparationUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ValidIngredientState{},
			Description: "",
		},
		{
			Type:        &types.ValidIngredientStateCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ValidIngredientStateUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ValidIngredientStateIngredient{},
			Description: "",
		},
		{
			Type:        &types.ValidIngredientStateIngredientCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ValidIngredientStateIngredientUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ValidIngredient{},
			Description: "",
		},
		{
			Type:        &types.ValidIngredientCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ValidIngredientUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ValidIngredientGroup{},
			Description: "",
		},
		{
			Type:        &types.ValidIngredientGroupCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ValidIngredientGroupUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ValidIngredientGroupMember{},
			Description: "",
		},
		{
			Type:        &types.ValidIngredientGroupMemberCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ValidInstrument{},
			Description: "",
		},
		{
			Type:        &types.ValidInstrumentCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ValidInstrumentUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ValidMeasurementUnitConversion{},
			Description: "",
		},
		{
			Type:        &types.ValidMeasurementUnitConversionCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ValidMeasurementUnitConversionUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ValidMeasurementUnit{},
			Description: "",
		},
		{
			Type:        &types.ValidMeasurementUnitCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ValidMeasurementUnitUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ValidPreparationInstrument{},
			Description: "",
		},
		{
			Type:        &types.ValidPreparationInstrumentCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ValidPreparationInstrumentUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ValidPreparationVessel{},
			Description: "",
		},
		{
			Type:        &types.ValidPreparationVesselCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ValidPreparationVesselUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ValidPreparation{},
			Description: "",
		},
		{
			Type:        &types.ValidPreparationCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ValidPreparationUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ValidVessel{},
			Description: "",
		},
		{
			Type:        &types.ValidVesselCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ValidVesselUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ServiceSetting{},
			Description: "",
		},
		{
			Type:        &types.ServiceSettingCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ServiceSettingUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ServiceSettingConfiguration{},
			Description: "",
		},
		{
			Type:        &types.ServiceSettingConfigurationCreationRequestInput{},
			Description: "",
		},
		{
			Type:        &types.ServiceSettingConfigurationUpdateRequestInput{},
			Description: "",
		},
		{
			Type:        &types.Webhook{},
			Description: "",
		},
		{
			Type:        &types.WebhookTriggerEvent{},
			Description: "",
		},
		{
			Type:        &types.WebhookCreationRequestInput{},
			Description: "",
		},
	}
)
