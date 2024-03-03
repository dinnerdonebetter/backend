import { Axios } from 'axios';
import format from 'string-format';

import { APIResponse, MealPlanTask, MealPlanTaskStatusChangeRequestInput } from '@dinnerdonebetter/models';

import { backendRoutes } from './routes';

export async function getMealPlanTask(
  client: Axios,
  mealPlanID: string,
  mealPlanTaskID: string,
): Promise<MealPlanTask> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<MealPlanTask>>(
      format(backendRoutes.MEAL_PLAN_TASKS, mealPlanID, mealPlanTaskID),
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getMealPlanTasks(client: Axios, mealPlanID: string): Promise<MealPlanTask[]> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<MealPlanTask[]>>(format(backendRoutes.MEAL_PLAN_TASKS, mealPlanID));

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function updateMealPlanTaskStatus(
  client: Axios,
  mealPlanID: string,
  mealPlanTaskID: string,
  input: MealPlanTaskStatusChangeRequestInput,
): Promise<MealPlanTask> {
  return new Promise(async function (resolve, reject) {
    const response = await client.patch<APIResponse<MealPlanTask>>(
      format(backendRoutes.MEAL_PLAN_TASK, mealPlanID, mealPlanTaskID),
      input,
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}
