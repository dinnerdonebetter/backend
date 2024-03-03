import { Axios } from 'axios';
import format from 'string-format';

import {
  ValidIngredientGroupCreationRequestInput,
  ValidIngredientGroup,
  QueryFilter,
  ValidIngredientGroupUpdateRequestInput,
  QueryFilteredResult,
  APIResponse,
} from '@dinnerdonebetter/models';

import { backendRoutes } from './routes';

export async function createValidIngredientGroup(
  client: Axios,
  input: ValidIngredientGroupCreationRequestInput,
): Promise<ValidIngredientGroup> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<ValidIngredientGroup>>(backendRoutes.VALID_INGREDIENT_GROUPS, input);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getValidIngredientGroup(
  client: Axios,
  validIngredientGroupID: string,
): Promise<ValidIngredientGroup> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ValidIngredientGroup>>(
      format(backendRoutes.VALID_INGREDIENT_GROUP, validIngredientGroupID),
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getValidIngredientGroups(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
): Promise<QueryFilteredResult<ValidIngredientGroup>> {
  return new Promise(async function (resolve, reject) {
    const repsonse = await client.get<APIResponse<ValidIngredientGroup[]>>(
      format(backendRoutes.VALID_INGREDIENT_GROUPS),
      {
        params: filter.asRecord(),
      },
    );

    if (repsonse.data.error) {
      reject(new Error(repsonse.data.error.message));
    }

    const result = new QueryFilteredResult<ValidIngredientGroup>({
      data: repsonse.data.data,
      totalCount: repsonse.data.pagination?.totalCount,
      page: repsonse.data.pagination?.page,
      limit: repsonse.data.pagination?.limit,
      filteredCount: repsonse.data.pagination?.filteredCount,
    });

    resolve(result);
  });
}

export async function updateValidIngredientGroup(
  client: Axios,
  validIngredientGroupID: string,
  input: ValidIngredientGroupUpdateRequestInput,
): Promise<ValidIngredientGroup> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<ValidIngredientGroup>>(
      format(backendRoutes.VALID_INGREDIENT_GROUP, validIngredientGroupID),
      input,
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function deleteValidIngredientGroup(
  client: Axios,
  validIngredientGroupID: string,
): Promise<ValidIngredientGroup> {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete<APIResponse<ValidIngredientGroup>>(
      format(backendRoutes.VALID_INGREDIENT_GROUP, validIngredientGroupID),
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function searchForValidIngredientGroups(
  client: Axios,
  query: string,
  filter: QueryFilter = QueryFilter.Default(),
): Promise<ValidIngredientGroup[]> {
  return new Promise(async function (resolve, reject) {
    const p = filter.asRecord();
    p['q'] = query;

    const response = await client.get<APIResponse<ValidIngredientGroup[]>>(
      backendRoutes.VALID_INGREDIENT_GROUPS_SEARCH,
      { params: p },
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}
