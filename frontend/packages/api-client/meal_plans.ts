import { Axios } from 'axios';
import format from 'string-format';

import {
  MealPlanCreationRequestInput,
  MealPlan,
  QueryFilter,
  MealPlanUpdateRequestInput,
  MealPlanOptionVoteCreationRequestInput,
  MealPlanOptionVote,
  QueryFilteredResult,
  APIResponse,
} from '@dinnerdonebetter/models';

import { backendRoutes } from './routes';

export async function createMealPlan(client: Axios, input: MealPlanCreationRequestInput): Promise<MealPlan> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<MealPlan>>(backendRoutes.MEAL_PLANS, input);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getMealPlan(client: Axios, mealPlanID: string): Promise<MealPlan> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<MealPlan>>(format(backendRoutes.MEAL_PLAN, mealPlanID));

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getMealPlans(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
): Promise<QueryFilteredResult<MealPlan>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<MealPlan[]>>(backendRoutes.MEAL_PLANS, { params: filter.asRecord() });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<MealPlan>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
      filteredCount: response.data.pagination?.filteredCount,
    });

    resolve(result);
  });
}

export async function updateMealPlan(
  client: Axios,
  mealPlanID: string,
  input: MealPlanUpdateRequestInput,
): Promise<MealPlan> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<MealPlan>>(format(backendRoutes.MEAL_PLAN, mealPlanID), input);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function deleteMealPlan(client: Axios, mealPlanID: string): Promise<MealPlan> {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete<APIResponse<MealPlan>>(format(backendRoutes.MEAL_PLAN, mealPlanID));

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function voteForMealPlan(
  client: Axios,
  mealPlanID: string,
  mealPlanEventID: string,
  votes: MealPlanOptionVoteCreationRequestInput,
): Promise<MealPlanOptionVote[]> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<MealPlanOptionVote[]>>(
      format(backendRoutes.MEAL_PLAN_VOTING, mealPlanID, mealPlanEventID),
      votes,
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}
