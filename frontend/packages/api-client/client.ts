import axios, { AxiosInstance, AxiosError, AxiosRequestConfig, AxiosResponse } from 'axios';
import { Span } from '@opentelemetry/api';

import { buildServerSideLogger } from '@dinnerdonebetter/logger';
import {
  MealPlanTask,
  MealPlanTaskStatusChangeRequestInput,
  Household,
  HouseholdInvitationCreationRequestInput,
  HouseholdInvitationUpdateRequestInput,
  HouseholdUpdateRequestInput,
  Meal,
  MealCreationRequestInput,
  MealPlan,
  MealPlanCreationRequestInput,
  MealPlanOptionVote,
  MealPlanOptionVoteCreationRequestInput,
  MealPlanUpdateRequestInput,
  MealUpdateRequestInput,
  PasswordResetTokenCreationRequestInput,
  PasswordResetTokenRedemptionRequestInput,
  QueryFilter,
  Recipe,
  RecipeCreationRequestInput,
  RecipeUpdateRequestInput,
  User,
  UserCreationResponse,
  UserAccountStatusUpdateInput,
  UserLoginInput,
  UsernameReminderRequestInput,
  UserPermissionsRequestInput,
  UserPermissionsResponse,
  UserRegistrationInput,
  UserStatusResponse,
  ValidIngredient,
  ValidIngredientCreationRequestInput,
  ValidIngredientMeasurementUnit,
  ValidIngredientMeasurementUnitCreationRequestInput,
  ValidIngredientPreparation,
  ValidIngredientPreparationCreationRequestInput,
  ValidIngredientUpdateRequestInput,
  ValidInstrument,
  ValidInstrumentCreationRequestInput,
  ValidInstrumentUpdateRequestInput,
  ValidMeasurementUnit,
  ValidMeasurementUnitCreationRequestInput,
  ValidMeasurementUnitUpdateRequestInput,
  ValidPreparation,
  ValidPreparationCreationRequestInput,
  ValidPreparationInstrument,
  ValidPreparationInstrumentCreationRequestInput,
  ValidPreparationUpdateRequestInput,
  HouseholdInvitation,
  MealPlanGroceryListItem,
  MealPlanGroceryListItemCreationRequestInput,
  MealPlanGroceryListItemUpdateRequestInput,
  ValidIngredientState,
  ValidIngredientStateCreationRequestInput,
  ValidIngredientStateUpdateRequestInput,
  ValidMeasurementUnitConversion,
  ValidMeasurementUnitConversionCreationRequestInput,
  ValidMeasurementUnitConversionUpdateRequestInput,
  ValidIngredientStateIngredient,
  ValidIngredientStateIngredientCreationRequestInput,
  QueryFilteredResult,
  ServiceSetting,
  ServiceSettingUpdateRequestInput,
  ServiceSettingCreationRequestInput,
  ServiceSettingConfigurationCreationRequestInput,
  ServiceSettingConfiguration,
  ServiceSettingConfigurationUpdateRequestInput,
  EmailAddressVerificationRequestInput,
  PasswordUpdateInput,
  AvatarUpdateInput,
  TOTPSecretRefreshInput,
  ValidIngredientGroup,
  ValidIngredientGroupCreationRequestInput,
  ValidIngredientGroupUpdateRequestInput,
  OAuth2Client,
  OAuth2ClientCreationRequestInput,
  ValidVessel,
  ValidVesselCreationRequestInput,
  ValidVesselUpdateRequestInput,
  ValidPreparationVessel,
  ValidPreparationVesselCreationRequestInput,
} from '@dinnerdonebetter/models';

import { createMeal, getMeal, getMeals, updateMeal, deleteMeal, searchForMeals } from './meals';
import {
  createValidPreparation,
  getValidPreparation,
  getValidPreparations,
  updateValidPreparation,
  deleteValidPreparation,
  searchForValidPreparations,
} from './valid_preparations';
import {
  createServiceSetting,
  getServiceSetting,
  getServiceSettings,
  updateServiceSetting,
  deleteServiceSetting,
  searchForServiceSettings,
} from './service_settings';
import {
  createServiceSettingConfiguration,
  getServiceSettingConfigurationsForUser,
  getServiceSettingConfigurationsForHousehold,
  updateServiceSettingConfiguration,
  deleteServiceSettingConfiguration,
} from './service_setting_configurations';

import {
  validPreparationInstrumentsForPreparationID,
  validPreparationInstrumentsForInstrumentID,
  createValidPreparationInstrument,
  deleteValidPreparationInstrument,
  getValidPreparationInstrument,
} from './valid_preparation_instruments';
import {
  validPreparationVesselsForPreparationID,
  validPreparationVesselsForVesselID,
  createValidPreparationVessel,
  deleteValidPreparationVessel,
  getValidPreparationVessel,
} from './valid_preparation_vessels';
import {
  createValidMeasurementUnit,
  getValidMeasurementUnit,
  getValidMeasurementUnits,
  updateValidMeasurementUnit,
  deleteValidMeasurementUnit,
  searchForValidMeasurementUnits,
  searchForValidMeasurementUnitsByIngredientID,
} from './valid_measurement_units';
import {
  createValidInstrument,
  getValidInstrument,
  getValidInstruments,
  updateValidInstrument,
  deleteValidInstrument,
  searchForValidInstruments,
} from './valid_instruments';
import {
  createValidIngredient,
  getValidIngredient,
  getValidIngredients,
  updateValidIngredient,
  deleteValidIngredient,
  searchForValidIngredients,
  getValidIngredientsForPreparation,
} from './valid_ingredients';
import {
  createValidIngredientGroup,
  getValidIngredientGroup,
  getValidIngredientGroups,
  updateValidIngredientGroup,
  deleteValidIngredientGroup,
  searchForValidIngredientGroups,
} from './valid_ingredient_groups';
import {
  logIn,
  adminLogin,
  logOut,
  register,
  checkPermissions,
  requestPasswordResetToken,
  redeemPasswordResetToken,
  changePassword,
  requestUsernameReminderEmail,
} from './auth';
import { getInvitation, acceptInvitation, cancelInvitation, rejectInvitation } from './household_invitations';
import {
  inviteUserToHousehold,
  removeMemberFromHousehold,
  getReceivedInvites,
  getSentInvites,
  getCurrentHouseholdInfo,
  getHousehold,
  getHouseholds,
  updateHousehold,
} from './households';
import { clientName } from './constants';
import {
  createMealPlan,
  getMealPlan,
  getMealPlans,
  updateMealPlan,
  deleteMealPlan,
  voteForMealPlan,
} from './meal_plans';
import { createRecipe, getRecipe, getRecipes, updateRecipe, deleteRecipe, searchForRecipes } from './recipes';
import {
  getUser,
  getUsers,
  updateUserAccountStatus,
  searchForUsers,
  verifyEmailAddress,
  fetchSelf,
  requestEmailVerificationEmail,
  uploadNewAvatar,
  newTwoFactorSecret,
} from './users';
import {
  validIngredientMeasurementUnitsForMeasurementUnitID,
  createValidIngredientMeasurementUnit,
  deleteValidIngredientMeasurementUnit,
  validIngredientMeasurementUnitsForIngredientID,
  getValidIngredientMeasurementUnit,
} from './valid_ingredient_measurement_units';
import {
  validIngredientPreparationsForPreparationID,
  validIngredientPreparationsForIngredientID,
  createValidIngredientPreparation,
  deleteValidIngredientPreparation,
  getValidIngredientPreparation,
} from './valid_ingredient_preparations';
import { getMealPlanTask, getMealPlanTasks, updateMealPlanTaskStatus } from './meal_plan_tasks';
import {
  createMealPlanGroceryListItem,
  getMealPlanGroceryListItem,
  getMealPlanGroceryListItems,
  updateMealPlanGroceryListItem,
  deleteMealPlanGroceryListItem,
} from './meal_plan_grocery_list_items';
import {
  createValidIngredientState,
  getValidIngredientState,
  getValidIngredientStates,
  updateValidIngredientState,
  deleteValidIngredientState,
  searchForValidIngredientStates,
} from './valid_ingredient_states';
import {
  createValidMeasurementUnitConversion,
  getValidMeasurementUnitConversion,
  getValidMeasurementUnitConversions,
  updateValidMeasurementUnitConversion,
  deleteValidMeasurementUnitConversion,
  getValidMeasurementUnitConversionsFromUnit,
  getValidMeasurementUnitConversionsToUnit,
} from './valid_measurement_unit_conversions';
import {
  validIngredientStateIngredientsForIngredientStateID,
  validIngredientStateIngredientsForIngredientID,
  createValidIngredientStateIngredient,
  deleteValidIngredientStateIngredient,
  getValidIngredientStateIngredient,
} from './valid_ingredient_state_ingredients';
import { createOAuth2Client, getOAuth2Client, getOAuth2Clients, deleteOAuth2Client } from './oauth2_clients';
import {
  createValidVessel,
  getValidVessel,
  getValidVessels,
  updateValidVessel,
  deleteValidVessel,
  searchForValidVessels,
} from './valid_vessels';

const cookieName = 'ddb_api_cookie';

const logger = buildServerSideLogger('api_client');

export class DinnerDoneBetterAPIClient {
  baseURL: string;
  client: AxiosInstance;
  traceID: string;
  requestInterceptorID: number;

  constructor(baseURL: string = '', cookie?: string, traceID?: string) {
    this.baseURL = baseURL;

    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      'X-Request-Source': 'webapp',
      'X-Service-Client': clientName,
    };

    if (cookie) {
      headers['Cookie'] = `${cookieName}=${cookie}`;
    }

    this.traceID = traceID || '';

    this.client = axios.create({
      baseURL,
      timeout: 10000,
      withCredentials: true,
      crossDomain: true,
      headers,
    } as AxiosRequestConfig);

    this.requestInterceptorID = this.client.interceptors.request.use((request: AxiosRequestConfig) => {
      logger.debug(`Request: ${request.method} ${request.url}`);
      return request;
    });

    this.client.interceptors.response.use((response: AxiosResponse) => {
      logger.debug(`Request: ${response.status}`);
      return response;
    });
  }

  withSpan(span: Span): DinnerDoneBetterAPIClient {
    const spanContext = span.spanContext();
    const spanLogDetails = { spanID: spanContext.spanId, traceID: spanContext.traceId };

    this.client.interceptors.request.eject(this.requestInterceptorID);
    this.requestInterceptorID = this.client.interceptors.request.use((request: AxiosRequestConfig) => {
      logger.debug(`Request: ${request.method} ${request.url}`, spanLogDetails);

      if (this.traceID) {
        request.headers = request.headers
          ? { ...request.headers, traceparent: this.traceID }
          : { traceparent: this.traceID };
      }

      return request;
    });

    return this;
  }

  // eslint-disable-next-line no-unused-vars
  configureRouterRejectionInterceptor(redirectCallback: (_: Location) => void) {
    this.client.interceptors.response.use(
      (response: AxiosResponse) => {
        return response;
      },
      (error: AxiosError) => {
        console.debug(`Request failed: ${error.response?.status}`);
        if (error.response?.status === 401) {
          redirectCallback(window.location);
        }

        return Promise.reject(error);
      },
    );
  }

  // auth

  async logIn(input: UserLoginInput): Promise<AxiosResponse<UserStatusResponse>> {
    return logIn(this.client, input);
  }

  async adminLogin(input: UserLoginInput): Promise<AxiosResponse<UserStatusResponse>> {
    return adminLogin(this.client, input);
  }

  async logOut(): Promise<AxiosResponse<UserStatusResponse>> {
    return logOut(this.client);
  }

  async register(input: UserRegistrationInput): Promise<UserCreationResponse> {
    return register(this.client, input);
  }

  async checkPermissions(body: UserPermissionsRequestInput): Promise<UserPermissionsResponse> {
    return checkPermissions(this.client, body);
  }

  async requestPasswordResetToken(input: PasswordResetTokenCreationRequestInput): Promise<AxiosResponse> {
    return requestPasswordResetToken(this.client, input);
  }

  async redeemPasswordResetToken(input: PasswordResetTokenRedemptionRequestInput): Promise<AxiosResponse> {
    return redeemPasswordResetToken(this.client, input);
  }

  async requestUsernameReminderEmail(input: UsernameReminderRequestInput): Promise<AxiosResponse> {
    return requestUsernameReminderEmail(this.client, input);
  }

  async changePassword(input: PasswordUpdateInput): Promise<AxiosResponse> {
    return changePassword(this.client, input);
  }

  // household invitations

  async getInvitation(invitationID: string): Promise<HouseholdInvitation> {
    return getInvitation(this.client, invitationID);
  }

  async acceptInvitation(
    invitationID: string,
    input: HouseholdInvitationUpdateRequestInput,
  ): Promise<HouseholdInvitation> {
    return acceptInvitation(this.client, invitationID, input);
  }

  async cancelInvitation(
    invitationID: string,
    input: HouseholdInvitationUpdateRequestInput,
  ): Promise<HouseholdInvitation> {
    return cancelInvitation(this.client, invitationID, input);
  }

  async rejectInvitation(
    invitationID: string,
    input: HouseholdInvitationUpdateRequestInput,
  ): Promise<HouseholdInvitation> {
    return rejectInvitation(this.client, invitationID, input);
  }

  // households

  async getCurrentHouseholdInfo(): Promise<Household> {
    return getCurrentHouseholdInfo(this.client);
  }

  async getHousehold(id: string): Promise<Household> {
    return getHousehold(this.client, id);
  }

  async getHouseholds(filter: QueryFilter = QueryFilter.Default()): Promise<QueryFilteredResult<Household>> {
    return getHouseholds(this.client, filter);
  }

  async updateHousehold(householdID: string, household: HouseholdUpdateRequestInput): Promise<Household> {
    return updateHousehold(this.client, householdID, household);
  }

  async inviteUserToHousehold(
    householdID: string,
    input: HouseholdInvitationCreationRequestInput,
  ): Promise<HouseholdInvitation> {
    return inviteUserToHousehold(this.client, householdID, input);
  }

  async removeMemberFromHousehold(householdID: string, memberID: string): Promise<Household> {
    return removeMemberFromHousehold(this.client, householdID, memberID);
  }

  async getReceivedInvites(): Promise<QueryFilteredResult<HouseholdInvitation>> {
    return getReceivedInvites(this.client);
  }

  async getSentInvites(): Promise<QueryFilteredResult<HouseholdInvitation>> {
    return getSentInvites(this.client);
  }

  // meal plans

  async createMealPlan(input: MealPlanCreationRequestInput): Promise<MealPlan> {
    return createMealPlan(this.client, input);
  }

  async getMealPlan(mealPlanID: string): Promise<MealPlan> {
    return getMealPlan(this.client, mealPlanID);
  }

  async getMealPlans(filter: QueryFilter = QueryFilter.Default()): Promise<QueryFilteredResult<MealPlan>> {
    return getMealPlans(this.client, filter);
  }

  async updateMealPlan(mealPlanID: string, input: MealPlanUpdateRequestInput): Promise<MealPlan> {
    return updateMealPlan(this.client, mealPlanID, input);
  }

  async deleteMealPlan(mealPlanID: string): Promise<MealPlan> {
    return deleteMealPlan(this.client, mealPlanID);
  }

  async voteForMealPlan(
    mealPlanID: string,
    mealPlanEventID: string,
    votes: MealPlanOptionVoteCreationRequestInput,
  ): Promise<MealPlanOptionVote[]> {
    return voteForMealPlan(this.client, mealPlanID, mealPlanEventID, votes);
  }

  // meals

  async createMeal(input: MealCreationRequestInput): Promise<Meal> {
    return createMeal(this.client, input);
  }

  async getMeal(mealID: string): Promise<Meal> {
    return getMeal(this.client, mealID);
  }

  async getMeals(filter: QueryFilter = QueryFilter.Default()): Promise<QueryFilteredResult<Meal>> {
    return getMeals(this.client, filter);
  }

  async updateMeal(mealID: string, input: MealUpdateRequestInput): Promise<Meal> {
    return updateMeal(this.client, mealID, input);
  }

  async deleteMeal(mealID: string): Promise<Meal> {
    return deleteMeal(this.client, mealID);
  }

  async searchForMeals(query: string): Promise<QueryFilteredResult<Meal>> {
    return searchForMeals(this.client, query);
  }

  // recipes

  async createRecipe(input: RecipeCreationRequestInput): Promise<Recipe> {
    return createRecipe(this.client, input);
  }

  async getRecipe(recipeID: string): Promise<Recipe> {
    return getRecipe(this.client, recipeID);
  }

  async getRecipes(filter: QueryFilter = QueryFilter.Default()): Promise<QueryFilteredResult<Recipe>> {
    return getRecipes(this.client, filter);
  }

  async updateRecipe(recipeID: string, input: RecipeUpdateRequestInput): Promise<Recipe> {
    return updateRecipe(this.client, recipeID, input);
  }

  async deleteRecipe(recipeID: string): Promise<Recipe> {
    return deleteRecipe(this.client, recipeID);
  }

  async searchForRecipes(query: string): Promise<QueryFilteredResult<Recipe>> {
    return searchForRecipes(this.client, query);
  }

  // users
  async self(): Promise<User> {
    return fetchSelf(this.client);
  }

  async requestEmailVerificationEmail(): Promise<User> {
    return requestEmailVerificationEmail(this.client);
  }

  async getUser(userID: string): Promise<User> {
    return getUser(this.client, userID);
  }

  async getUsers(filter: QueryFilter = QueryFilter.Default()): Promise<QueryFilteredResult<User>> {
    return getUsers(this.client, filter);
  }

  async updateUserAccountStatus(input: UserAccountStatusUpdateInput): Promise<User> {
    return updateUserAccountStatus(this.client, input);
  }

  async newTwoFactorSecret(input: TOTPSecretRefreshInput): Promise<User> {
    return newTwoFactorSecret(this.client, input);
  }

  async searchForUsers(query: string): Promise<User[]> {
    return searchForUsers(this.client, query);
  }

  async verifyEmailAddress(input: EmailAddressVerificationRequestInput): Promise<User> {
    return verifyEmailAddress(this.client, input);
  }

  async uploadNewAvatar(input: AvatarUpdateInput): Promise<User> {
    return uploadNewAvatar(this.client, input);
  }

  // valid ingredient measurement units

  async getValidIngredientMeasurementUnit(
    validIngredientMeasurementUnitID: string,
  ): Promise<ValidIngredientMeasurementUnit> {
    return getValidIngredientMeasurementUnit(this.client, validIngredientMeasurementUnitID);
  }

  async validIngredientMeasurementUnitsForIngredientID(
    validIngredientID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredientMeasurementUnit>> {
    return validIngredientMeasurementUnitsForIngredientID(this.client, validIngredientID, filter);
  }

  async validIngredientMeasurementUnitsForMeasurementUnitID(
    validMeasurementUnitID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredientMeasurementUnit>> {
    return validIngredientMeasurementUnitsForMeasurementUnitID(this.client, validMeasurementUnitID, filter);
  }

  async createValidIngredientMeasurementUnit(
    input: ValidIngredientMeasurementUnitCreationRequestInput,
  ): Promise<ValidIngredientMeasurementUnit> {
    return createValidIngredientMeasurementUnit(this.client, input);
  }

  async deleteValidIngredientMeasurementUnit(
    validIngredientMeasurementUnitID: string,
  ): Promise<ValidIngredientMeasurementUnit> {
    return deleteValidIngredientMeasurementUnit(this.client, validIngredientMeasurementUnitID);
  }

  // valid ingredient preparations

  async getValidIngredientPreparation(validIngredientPreparationID: string): Promise<ValidIngredientPreparation> {
    return getValidIngredientPreparation(this.client, validIngredientPreparationID);
  }

  async validIngredientPreparationsForPreparationID(
    validPreparationID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredientPreparation>> {
    return validIngredientPreparationsForPreparationID(this.client, validPreparationID, filter);
  }

  async validIngredientPreparationsForIngredientID(
    validIngredientID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredientPreparation>> {
    return validIngredientPreparationsForIngredientID(this.client, validIngredientID, filter);
  }

  async createValidIngredientPreparation(
    input: ValidIngredientPreparationCreationRequestInput,
  ): Promise<ValidIngredientPreparation> {
    return createValidIngredientPreparation(this.client, input);
  }

  async deleteValidIngredientPreparation(validIngredientPreparationID: string): Promise<ValidIngredientPreparation> {
    return deleteValidIngredientPreparation(this.client, validIngredientPreparationID);
  }

  // valid ingredient state ingredients

  async getValidIngredientStateIngredient(validIngredientStateID: string): Promise<ValidIngredientStateIngredient> {
    return getValidIngredientStateIngredient(this.client, validIngredientStateID);
  }

  async validIngredientStateIngredientsForIngredientStateID(
    validIngredientStateID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredientStateIngredient>> {
    return validIngredientStateIngredientsForIngredientStateID(this.client, validIngredientStateID, filter);
  }

  async validIngredientStateIngredientsForIngredientID(
    validIngredientID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredientStateIngredient>> {
    return validIngredientStateIngredientsForIngredientID(this.client, validIngredientID, filter);
  }

  async createValidIngredientStateIngredient(
    input: ValidIngredientStateIngredientCreationRequestInput,
  ): Promise<ValidIngredientStateIngredient> {
    return createValidIngredientStateIngredient(this.client, input);
  }

  async deleteValidIngredientStateIngredient(
    validIngredientStateIngredientID: string,
  ): Promise<ValidIngredientStateIngredient> {
    return deleteValidIngredientStateIngredient(this.client, validIngredientStateIngredientID);
  }

  // valid ingredients
  async createValidIngredient(input: ValidIngredientCreationRequestInput): Promise<ValidIngredient> {
    return createValidIngredient(this.client, input);
  }

  async getValidIngredient(validIngredientID: string): Promise<ValidIngredient> {
    return getValidIngredient(this.client, validIngredientID);
  }

  async getValidIngredients(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredient>> {
    return getValidIngredients(this.client, filter);
  }

  async updateValidIngredient(
    validIngredientID: string,
    input: ValidIngredientUpdateRequestInput,
  ): Promise<ValidIngredient> {
    return updateValidIngredient(this.client, validIngredientID, input);
  }

  async deleteValidIngredient(validIngredientID: string): Promise<ValidIngredient> {
    return deleteValidIngredient(this.client, validIngredientID);
  }

  async searchForValidIngredients(query: string): Promise<QueryFilteredResult<ValidIngredient>> {
    return searchForValidIngredients(this.client, query);
  }

  async getValidIngredientsForPreparation(
    validPreparationID: string,
    query: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredient>> {
    return getValidIngredientsForPreparation(this.client, validPreparationID, query, filter);
  }

  // valid ingredient groups
  async createValidIngredientGroup(input: ValidIngredientGroupCreationRequestInput): Promise<ValidIngredientGroup> {
    return createValidIngredientGroup(this.client, input);
  }

  async getValidIngredientGroup(validIngredientGroupID: string): Promise<ValidIngredientGroup> {
    return getValidIngredientGroup(this.client, validIngredientGroupID);
  }

  async getValidIngredientGroups(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredientGroup>> {
    return getValidIngredientGroups(this.client, filter);
  }

  async updateValidIngredientGroup(
    validIngredientGroupID: string,
    input: ValidIngredientGroupUpdateRequestInput,
  ): Promise<ValidIngredientGroup> {
    return updateValidIngredientGroup(this.client, validIngredientGroupID, input);
  }

  async deleteValidIngredientGroup(validIngredientGroupID: string): Promise<ValidIngredientGroup> {
    return deleteValidIngredientGroup(this.client, validIngredientGroupID);
  }

  async searchForValidIngredientGroups(query: string): Promise<ValidIngredientGroup[]> {
    return searchForValidIngredientGroups(this.client, query);
  }

  // valid instruments
  async createValidInstrument(input: ValidInstrumentCreationRequestInput): Promise<ValidInstrument> {
    return createValidInstrument(this.client, input);
  }

  async getValidInstrument(validInstrumentID: string): Promise<ValidInstrument> {
    return getValidInstrument(this.client, validInstrumentID);
  }

  async getValidInstruments(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidInstrument>> {
    return getValidInstruments(this.client, filter);
  }

  async updateValidInstrument(
    validInstrumentID: string,
    input: ValidInstrumentUpdateRequestInput,
  ): Promise<ValidInstrument> {
    return updateValidInstrument(this.client, validInstrumentID, input);
  }

  async deleteValidInstrument(validInstrumentID: string): Promise<ValidInstrument> {
    return deleteValidInstrument(this.client, validInstrumentID);
  }

  async searchForValidInstruments(query: string): Promise<ValidInstrument[]> {
    return searchForValidInstruments(this.client, query);
  }

  // valid vessels
  async createValidVessel(input: ValidVesselCreationRequestInput): Promise<ValidVessel> {
    return createValidVessel(this.client, input);
  }

  async getValidVessel(validVesselID: string): Promise<ValidVessel> {
    return getValidVessel(this.client, validVesselID);
  }

  async getValidVessels(filter: QueryFilter = QueryFilter.Default()): Promise<QueryFilteredResult<ValidVessel>> {
    return getValidVessels(this.client, filter);
  }

  async updateValidVessel(validVesselID: string, input: ValidVesselUpdateRequestInput): Promise<ValidVessel> {
    return updateValidVessel(this.client, validVesselID, input);
  }

  async deleteValidVessel(validVesselID: string): Promise<ValidVessel> {
    return deleteValidVessel(this.client, validVesselID);
  }

  async searchForValidVessels(query: string): Promise<ValidVessel[]> {
    return searchForValidVessels(this.client, query);
  }

  // valid measurement units
  async createValidMeasurementUnit(input: ValidMeasurementUnitCreationRequestInput): Promise<ValidMeasurementUnit> {
    return createValidMeasurementUnit(this.client, input);
  }

  async getValidMeasurementUnit(validMeasurementUnitID: string): Promise<ValidMeasurementUnit> {
    return getValidMeasurementUnit(this.client, validMeasurementUnitID);
  }

  async getValidMeasurementUnits(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidMeasurementUnit>> {
    return getValidMeasurementUnits(this.client, filter);
  }

  async updateValidMeasurementUnit(
    validMeasurementUnitID: string,
    input: ValidMeasurementUnitUpdateRequestInput,
  ): Promise<ValidMeasurementUnit> {
    return updateValidMeasurementUnit(this.client, validMeasurementUnitID, input);
  }

  async deleteValidMeasurementUnit(validMeasurementUnitID: string): Promise<ValidMeasurementUnit> {
    return deleteValidMeasurementUnit(this.client, validMeasurementUnitID);
  }

  async searchForValidMeasurementUnits(query: string): Promise<ValidMeasurementUnit[]> {
    return searchForValidMeasurementUnits(this.client, query);
  }

  async searchForValidMeasurementUnitsByIngredientID(
    validIngredientID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidMeasurementUnit>> {
    return searchForValidMeasurementUnitsByIngredientID(this.client, validIngredientID, filter);
  }

  // valid measurement unit conversions
  async createValidMeasurementUnitConversion(
    input: ValidMeasurementUnitConversionCreationRequestInput,
  ): Promise<ValidMeasurementUnitConversion> {
    return createValidMeasurementUnitConversion(this.client, input);
  }

  async getValidMeasurementUnitConversion(validMeasurementUnitID: string): Promise<ValidMeasurementUnitConversion> {
    return getValidMeasurementUnitConversion(this.client, validMeasurementUnitID);
  }

  async getValidMeasurementUnitConversions(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidMeasurementUnitConversion>> {
    return getValidMeasurementUnitConversions(this.client, filter);
  }

  async getValidMeasurementUnitConversionsFromUnit(
    validMeasurementUnitID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<ValidMeasurementUnitConversion[]> {
    return getValidMeasurementUnitConversionsFromUnit(this.client, validMeasurementUnitID, filter);
  }

  async getValidMeasurementUnitConversionsToUnit(
    validMeasurementUnitID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<ValidMeasurementUnitConversion[]> {
    return getValidMeasurementUnitConversionsToUnit(this.client, validMeasurementUnitID, filter);
  }

  async updateValidMeasurementUnitConversion(
    validMeasurementUnitID: string,
    input: ValidMeasurementUnitConversionUpdateRequestInput,
  ): Promise<ValidMeasurementUnitConversion> {
    return updateValidMeasurementUnitConversion(this.client, validMeasurementUnitID, input);
  }

  async deleteValidMeasurementUnitConversion(validMeasurementUnitID: string): Promise<ValidMeasurementUnitConversion> {
    return deleteValidMeasurementUnitConversion(this.client, validMeasurementUnitID);
  }

  // valid preparation instruments
  async getValidPreparationInstrument(validPreparationInstrumentID: string): Promise<ValidPreparationInstrument> {
    return getValidPreparationInstrument(this.client, validPreparationInstrumentID);
  }

  async validPreparationInstrumentsForPreparationID(
    validPreparationID: string,
  ): Promise<QueryFilteredResult<ValidPreparationInstrument>> {
    return validPreparationInstrumentsForPreparationID(this.client, validPreparationID);
  }

  async validPreparationInstrumentsForInstrumentID(
    validInstrumentID: string,
  ): Promise<QueryFilteredResult<ValidPreparationInstrument>> {
    return validPreparationInstrumentsForInstrumentID(this.client, validInstrumentID);
  }

  async createValidPreparationInstrument(
    input: ValidPreparationInstrumentCreationRequestInput,
  ): Promise<ValidPreparationInstrument> {
    return createValidPreparationInstrument(this.client, input);
  }

  async deleteValidPreparationInstrument(validPreparationInstrumentID: string): Promise<ValidPreparationInstrument> {
    return deleteValidPreparationInstrument(this.client, validPreparationInstrumentID);
  }

  // valid preparation vessels
  async getValidPreparationVessel(validPreparationVesselID: string): Promise<ValidPreparationVessel> {
    return getValidPreparationVessel(this.client, validPreparationVesselID);
  }

  async validPreparationVesselsForPreparationID(
    validPreparationID: string,
  ): Promise<QueryFilteredResult<ValidPreparationVessel>> {
    return validPreparationVesselsForPreparationID(this.client, validPreparationID);
  }

  async validPreparationVesselsForVesselID(
    validInstrumentID: string,
  ): Promise<QueryFilteredResult<ValidPreparationVessel>> {
    return validPreparationVesselsForVesselID(this.client, validInstrumentID);
  }

  async createValidPreparationVessel(
    input: ValidPreparationVesselCreationRequestInput,
  ): Promise<ValidPreparationVessel> {
    return createValidPreparationVessel(this.client, input);
  }

  async deleteValidPreparationVessel(validPreparationVesselID: string): Promise<ValidPreparationVessel> {
    return deleteValidPreparationVessel(this.client, validPreparationVesselID);
  }

  // valid preparations
  async createValidPreparation(input: ValidPreparationCreationRequestInput): Promise<ValidPreparation> {
    return createValidPreparation(this.client, input);
  }

  async getValidPreparation(validPreparationID: string): Promise<ValidPreparation> {
    return getValidPreparation(this.client, validPreparationID);
  }

  async getValidPreparations(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidPreparation>> {
    return getValidPreparations(this.client, filter);
  }

  async updateValidPreparation(
    validPreparationID: string,
    input: ValidPreparationUpdateRequestInput,
  ): Promise<ValidPreparation> {
    return updateValidPreparation(this.client, validPreparationID, input);
  }

  async deleteValidPreparation(validPreparationID: string): Promise<ValidPreparation> {
    return deleteValidPreparation(this.client, validPreparationID);
  }

  async searchForValidPreparations(query: string): Promise<ValidPreparation[]> {
    return searchForValidPreparations(this.client, query);
  }

  // service setting
  async createServiceSetting(input: ServiceSettingCreationRequestInput): Promise<ServiceSetting> {
    return createServiceSetting(this.client, input);
  }

  async getServiceSetting(serviceSettingID: string): Promise<ServiceSetting> {
    return getServiceSetting(this.client, serviceSettingID);
  }

  async getServiceSettings(filter: QueryFilter = QueryFilter.Default()): Promise<QueryFilteredResult<ServiceSetting>> {
    return getServiceSettings(this.client, filter);
  }

  async updateServiceSetting(
    serviceSettingID: string,
    input: ServiceSettingUpdateRequestInput,
  ): Promise<ServiceSetting> {
    return updateServiceSetting(this.client, serviceSettingID, input);
  }

  async deleteServiceSetting(serviceSettingID: string): Promise<ServiceSetting> {
    return deleteServiceSetting(this.client, serviceSettingID);
  }

  async searchForServiceSettings(query: string): Promise<ServiceSetting[]> {
    return searchForServiceSettings(this.client, query);
  }

  async createServiceSettingConfiguration(
    input: ServiceSettingConfigurationCreationRequestInput,
  ): Promise<ServiceSettingConfiguration> {
    return createServiceSettingConfiguration(this.client, input);
  }

  async getServiceSettingConfigurationsForUser(): Promise<QueryFilteredResult<ServiceSettingConfiguration>> {
    return getServiceSettingConfigurationsForUser(this.client);
  }

  async getServiceSettingConfigurationsForHousehold(): Promise<QueryFilteredResult<ServiceSettingConfiguration>> {
    return getServiceSettingConfigurationsForHousehold(this.client);
  }

  async updateServiceSettingConfiguration(
    serviceSettingConfigurationID: string,
    input: ServiceSettingConfigurationUpdateRequestInput,
  ): Promise<ServiceSettingConfiguration> {
    return updateServiceSettingConfiguration(this.client, serviceSettingConfigurationID, input);
  }

  async deleteServiceSettingConfiguration(serviceSettingConfigurationID: string): Promise<ServiceSettingConfiguration> {
    return deleteServiceSettingConfiguration(this.client, serviceSettingConfigurationID);
  }

  // valid ingredient states
  async createValidIngredientState(input: ValidIngredientStateCreationRequestInput): Promise<ValidIngredientState> {
    return createValidIngredientState(this.client, input);
  }

  async getValidIngredientState(validPreparationID: string): Promise<ValidIngredientState> {
    return getValidIngredientState(this.client, validPreparationID);
  }

  async getValidIngredientStates(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredientState>> {
    return getValidIngredientStates(this.client, filter);
  }

  async updateValidIngredientState(
    validPreparationID: string,
    input: ValidIngredientStateUpdateRequestInput,
  ): Promise<ValidIngredientState> {
    return updateValidIngredientState(this.client, validPreparationID, input);
  }

  async deleteValidIngredientState(validPreparationID: string): Promise<ValidIngredientState> {
    return deleteValidIngredientState(this.client, validPreparationID);
  }

  async searchForValidIngredientStates(query: string): Promise<ValidIngredientState[]> {
    return searchForValidIngredientStates(this.client, query);
  }

  // meal plan tasks

  async getMealPlanTask(mealPlanID: string, mealPlanTaskID: string): Promise<MealPlanTask> {
    return getMealPlanTask(this.client, mealPlanID, mealPlanTaskID);
  }

  async getMealPlanTasks(mealPlanID: string): Promise<MealPlanTask[]> {
    return getMealPlanTasks(this.client, mealPlanID);
  }

  async updateMealPlanTaskStatus(
    mealPlanID: string,
    mealPlanTaskID: string,
    input: MealPlanTaskStatusChangeRequestInput,
  ): Promise<MealPlanTask> {
    return updateMealPlanTaskStatus(this.client, mealPlanID, mealPlanTaskID, input);
  }

  async createMealPlanGroceryListItem(
    mealPlanID: string,
    input: MealPlanGroceryListItemCreationRequestInput,
  ): Promise<MealPlanGroceryListItem> {
    return createMealPlanGroceryListItem(this.client, mealPlanID, input);
  }

  async getMealPlanGroceryListItem(mealPlanID: string): Promise<MealPlanGroceryListItem> {
    return getMealPlanGroceryListItem(this.client, mealPlanID);
  }

  async getMealPlanGroceryListItems(mealPlanID: string): Promise<MealPlanGroceryListItem[]> {
    return getMealPlanGroceryListItems(this.client, mealPlanID);
  }

  async updateMealPlanGroceryListItem(
    mealPlanID: string,
    mealPlanGroceryListItemID: string,
    input: MealPlanGroceryListItemUpdateRequestInput,
  ): Promise<MealPlanGroceryListItem> {
    return updateMealPlanGroceryListItem(this.client, mealPlanID, mealPlanGroceryListItemID, input);
  }

  async deleteMealPlanGroceryListItem(
    mealPlanID: string,
    mealPlanGroceryListItemID: string,
  ): Promise<MealPlanGroceryListItem> {
    return deleteMealPlanGroceryListItem(this.client, mealPlanID, mealPlanGroceryListItemID);
  }

  async createOAuth2Client(input: OAuth2ClientCreationRequestInput): Promise<OAuth2Client> {
    return createOAuth2Client(this.client, input);
  }

  async getOAuth2Client(oauth2ClientID: string): Promise<OAuth2Client> {
    return getOAuth2Client(this.client, oauth2ClientID);
  }

  async getOAuth2Clients(filter: QueryFilter = QueryFilter.Default()): Promise<QueryFilteredResult<OAuth2Client>> {
    return getOAuth2Clients(this.client, filter);
  }

  async deleteOAuth2Client(oauth2ClientID: string): Promise<OAuth2Client> {
    return deleteOAuth2Client(this.client, oauth2ClientID);
  }
}
