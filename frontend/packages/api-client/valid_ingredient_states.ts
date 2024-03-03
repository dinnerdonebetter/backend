import { Axios } from 'axios';
import format from 'string-format';

import {
  ValidIngredientStateCreationRequestInput,
  ValidIngredientState,
  QueryFilter,
  ValidIngredientStateUpdateRequestInput,
  QueryFilteredResult,
  APIResponse,
} from '@dinnerdonebetter/models';

import { backendRoutes } from './routes';

export async function createValidIngredientState(
  client: Axios,
  input: ValidIngredientStateCreationRequestInput,
): Promise<ValidIngredientState> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<ValidIngredientState>>(backendRoutes.VALID_INGREDIENT_STATES, input);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getValidIngredientState(
  client: Axios,
  validIngredientStateID: string,
): Promise<ValidIngredientState> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ValidIngredientState>>(
      format(backendRoutes.VALID_INGREDIENT_STATE, validIngredientStateID),
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getValidIngredientStates(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
): Promise<QueryFilteredResult<ValidIngredientState>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ValidIngredientState[]>>(backendRoutes.VALID_INGREDIENT_STATES, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<ValidIngredientState>({
      data: response.data.data,
      filteredCount: response.data.pagination?.filteredCount,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}

export async function updateValidIngredientState(
  client: Axios,
  validIngredientStateID: string,
  input: ValidIngredientStateUpdateRequestInput,
): Promise<ValidIngredientState> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<ValidIngredientState>>(
      format(backendRoutes.VALID_INGREDIENT_STATE, validIngredientStateID),
      input,
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function deleteValidIngredientState(
  client: Axios,
  validIngredientStateID: string,
): Promise<ValidIngredientState> {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete<APIResponse<ValidIngredientState>>(
      format(backendRoutes.VALID_INGREDIENT_STATE, validIngredientStateID),
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function searchForValidIngredientStates(client: Axios, query: string): Promise<ValidIngredientState[]> {
  return new Promise(async function (resolve, reject) {
    const searchURL = `${backendRoutes.VALID_INGREDIENT_STATES_SEARCH}?q=${encodeURIComponent(query)}`;
    const response = await client.get<APIResponse<ValidIngredientState[]>>(searchURL);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}
