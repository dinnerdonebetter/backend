import { Axios } from 'axios';
import format from 'string-format';
import {
  QueryFilter,
  ValidIngredientPreparationCreationRequestInput,
  ValidIngredientPreparation,
  QueryFilteredResult,
  APIResponse,
} from '@dinnerdonebetter/models';

import { backendRoutes } from './routes';

export async function getValidIngredientPreparation(
  client: Axios,
  validIngredientPreparationID: string,
): Promise<ValidIngredientPreparation> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ValidIngredientPreparation>>(
      format(backendRoutes.VALID_INGREDIENT_PREPARATION, validIngredientPreparationID),
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function validIngredientPreparationsForPreparationID(
  client: Axios,
  validPreparationID: string,
  filter: QueryFilter = QueryFilter.Default(),
): Promise<QueryFilteredResult<ValidIngredientPreparation>> {
  return new Promise(async function (resolve, reject) {
    const searchURL = format(backendRoutes.VALID_INGREDIENT_PREPARATIONS_SEARCH_BY_PREPARATION_ID, validPreparationID);
    const response = await client.get<APIResponse<ValidIngredientPreparation[]>>(searchURL, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<ValidIngredientPreparation>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      filteredCount: response.data.pagination?.filteredCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}

export async function validIngredientPreparationsForIngredientID(
  client: Axios,
  validIngredientID: string,
  filter: QueryFilter = QueryFilter.Default(),
): Promise<QueryFilteredResult<ValidIngredientPreparation>> {
  return new Promise(async function (resolve, reject) {
    const searchURL = format(backendRoutes.VALID_INGREDIENT_PREPARATIONS_SEARCH_BY_INGREDIENT_ID, validIngredientID);
    const response = await client.get<APIResponse<ValidIngredientPreparation[]>>(searchURL, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<ValidIngredientPreparation>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      filteredCount: response.data.pagination?.filteredCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}

export async function createValidIngredientPreparation(
  client: Axios,
  input: ValidIngredientPreparationCreationRequestInput,
): Promise<ValidIngredientPreparation> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<ValidIngredientPreparation>>(
      backendRoutes.VALID_INGREDIENT_PREPARATIONS,
      input,
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function deleteValidIngredientPreparation(
  client: Axios,
  validIngredientPreparationID: string,
): Promise<ValidIngredientPreparation> {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete<APIResponse<ValidIngredientPreparation>>(
      format(backendRoutes.VALID_INGREDIENT_PREPARATION, validIngredientPreparationID),
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}
