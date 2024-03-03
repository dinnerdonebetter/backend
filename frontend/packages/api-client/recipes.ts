import { Axios } from 'axios';
import format from 'string-format';

import {
  RecipeCreationRequestInput,
  Recipe,
  QueryFilter,
  RecipeUpdateRequestInput,
  QueryFilteredResult,
  APIResponse,
} from '@dinnerdonebetter/models';

import { backendRoutes } from './routes';

export async function createRecipe(client: Axios, input: RecipeCreationRequestInput): Promise<Recipe> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<Recipe>>(backendRoutes.RECIPES, input);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getRecipe(client: Axios, recipeID: string): Promise<Recipe> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<Recipe>>(format(backendRoutes.RECIPE, recipeID));

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getRecipes(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
): Promise<QueryFilteredResult<Recipe>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<Recipe[]>>(backendRoutes.RECIPES, { params: filter.asRecord() });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<Recipe>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
      filteredCount: response.data.pagination?.filteredCount,
    });

    resolve(result);
  });
}

export async function updateRecipe(client: Axios, recipeID: string, input: RecipeUpdateRequestInput): Promise<Recipe> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<Recipe>>(format(backendRoutes.RECIPE, recipeID), input);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function deleteRecipe(client: Axios, recipeID: string): Promise<Recipe> {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete<APIResponse<Recipe>>(format(backendRoutes.RECIPE, recipeID));

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function searchForRecipes(client: Axios, query: string): Promise<QueryFilteredResult<Recipe>> {
  return new Promise(async function (resolve, reject) {
    const searchURL = `${backendRoutes.RECIPES_SEARCH}?q=${encodeURIComponent(query)}`;
    const response = await client.get<APIResponse<Recipe[]>>(searchURL);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<Recipe>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
      filteredCount: response.data.pagination?.filteredCount,
    });

    resolve(result);
  });
}
