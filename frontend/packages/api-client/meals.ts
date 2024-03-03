import { Axios } from 'axios';
import format from 'string-format';

import {
  MealCreationRequestInput,
  Meal,
  MealUpdateRequestInput,
  QueryFilter,
  QueryFilteredResult,
  APIResponse,
} from '@dinnerdonebetter/models';

import { backendRoutes } from './routes';

export async function createMeal(client: Axios, input: MealCreationRequestInput): Promise<Meal> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<Meal>>(backendRoutes.MEALS, input);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getMeal(client: Axios, mealID: string): Promise<Meal> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<Meal>>(format(backendRoutes.MEAL, mealID));

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getMeals(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
): Promise<QueryFilteredResult<Meal>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<Meal[]>>(backendRoutes.MEALS, { params: filter.asRecord() });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<Meal>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
      filteredCount: response.data.pagination?.filteredCount,
    });

    resolve(result);
  });
}

export async function updateMeal(client: Axios, mealID: string, input: MealUpdateRequestInput): Promise<Meal> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<Meal>>(format(backendRoutes.MEAL, mealID), input);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function deleteMeal(client: Axios, mealID: string): Promise<Meal> {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete<APIResponse<Meal>>(format(backendRoutes.MEAL, mealID));

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function searchForMeals(client: Axios, query: string): Promise<QueryFilteredResult<Meal>> {
  return new Promise(async function (resolve, reject) {
    const searchURL = `${backendRoutes.MEALS_SEARCH}?q=${encodeURIComponent(query)}`;
    const response = await client.get<APIResponse<Meal[]>>(searchURL);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<Meal>({
      data: response.data.data,
      filteredCount: response.data.pagination?.filteredCount,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}
