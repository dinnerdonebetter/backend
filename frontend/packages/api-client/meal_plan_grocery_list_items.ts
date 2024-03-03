import { Axios } from 'axios';
import format from 'string-format';

import {
  MealPlanGroceryListItemCreationRequestInput,
  MealPlanGroceryListItem,
  MealPlanGroceryListItemUpdateRequestInput,
  APIResponse,
} from '@dinnerdonebetter/models';

import { backendRoutes } from './routes';

export async function createMealPlanGroceryListItem(
  client: Axios,
  mealPlanID: string,
  input: MealPlanGroceryListItemCreationRequestInput,
): Promise<MealPlanGroceryListItem> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<MealPlanGroceryListItem>>(
      format(backendRoutes.MEAL_PLAN_GROCERY_LIST_ITEMS, mealPlanID),
      input,
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getMealPlanGroceryListItem(client: Axios, mealPlanID: string): Promise<MealPlanGroceryListItem> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<MealPlanGroceryListItem>>(
      format(backendRoutes.MEAL_PLAN_GROCERY_LIST_ITEM, mealPlanID),
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getMealPlanGroceryListItems(
  client: Axios,
  mealPlanID: string,
): Promise<MealPlanGroceryListItem[]> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<MealPlanGroceryListItem[]>>(
      format(backendRoutes.MEAL_PLAN_GROCERY_LIST_ITEMS, mealPlanID),
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function updateMealPlanGroceryListItem(
  client: Axios,
  mealPlanID: string,
  mealPlanGroceryListItemID: string,
  input: MealPlanGroceryListItemUpdateRequestInput,
): Promise<MealPlanGroceryListItem> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<MealPlanGroceryListItem>>(
      format(backendRoutes.MEAL_PLAN_GROCERY_LIST_ITEM, mealPlanID, mealPlanGroceryListItemID),
      input,
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function deleteMealPlanGroceryListItem(
  client: Axios,
  mealPlanID: string,
  mealPlanGroceryListItemID: string,
): Promise<MealPlanGroceryListItem> {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete<APIResponse<MealPlanGroceryListItem>>(
      format(backendRoutes.MEAL_PLAN_GROCERY_LIST_ITEM, mealPlanID, mealPlanGroceryListItemID),
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}
