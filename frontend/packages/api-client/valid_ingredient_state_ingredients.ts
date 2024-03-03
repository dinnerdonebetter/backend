import { Axios } from 'axios';
import format from 'string-format';
import {
  QueryFilter,
  ValidIngredientStateIngredientCreationRequestInput,
  ValidIngredientStateIngredient,
  QueryFilteredResult,
  APIResponse,
} from '@dinnerdonebetter/models';

import { backendRoutes } from './routes';

export async function getValidIngredientStateIngredient(
  client: Axios,
  validIngredientStateIngredientID: string,
): Promise<ValidIngredientStateIngredient> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ValidIngredientStateIngredient>>(
      format(backendRoutes.VALID_INGREDIENT_STATE_INGREDIENT, validIngredientStateIngredientID),
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function validIngredientStateIngredientsForIngredientStateID(
  client: Axios,
  validIngredientStateID: string,
  filter: QueryFilter = QueryFilter.Default(),
): Promise<QueryFilteredResult<ValidIngredientStateIngredient>> {
  return new Promise(async function (resolve, reject) {
    const searchURL = format(
      backendRoutes.VALID_INGREDIENT_STATE_INGREDIENTS_SEARCH_BY_INGREDIENT_STATE,
      validIngredientStateID,
    );
    const response = await client.get<APIResponse<ValidIngredientStateIngredient[]>>(searchURL, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<ValidIngredientStateIngredient>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      filteredCount: response.data.pagination?.filteredCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}

export async function validIngredientStateIngredientsForIngredientID(
  client: Axios,
  validIngredientID: string,
  filter: QueryFilter = QueryFilter.Default(),
): Promise<QueryFilteredResult<ValidIngredientStateIngredient>> {
  return new Promise(async function (resolve, reject) {
    const searchURL = format(
      backendRoutes.VALID_INGREDIENT_STATE_INGREDIENTS_SEARCH_BY_INGREDIENT_ID,
      validIngredientID,
    );
    const response = await client.get<APIResponse<ValidIngredientStateIngredient[]>>(searchURL, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<ValidIngredientStateIngredient>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      filteredCount: response.data.pagination?.filteredCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}

export async function createValidIngredientStateIngredient(
  client: Axios,
  input: ValidIngredientStateIngredientCreationRequestInput,
): Promise<ValidIngredientStateIngredient> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<ValidIngredientStateIngredient>>(
      backendRoutes.VALID_INGREDIENT_STATE_INGREDIENTS,
      input,
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function deleteValidIngredientStateIngredient(
  client: Axios,
  validIngredientStateIngredientID: string,
): Promise<ValidIngredientStateIngredient> {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete<APIResponse<ValidIngredientStateIngredient>>(
      format(backendRoutes.VALID_INGREDIENT_STATE_INGREDIENT, validIngredientStateIngredientID),
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}
