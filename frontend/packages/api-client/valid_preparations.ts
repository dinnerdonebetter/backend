import { Axios } from 'axios';
import format from 'string-format';

import {
  ValidPreparationCreationRequestInput,
  ValidPreparation,
  QueryFilter,
  ValidPreparationUpdateRequestInput,
  QueryFilteredResult,
  APIResponse,
} from '@dinnerdonebetter/models';

import { backendRoutes } from './routes';

export async function createValidPreparation(
  client: Axios,
  input: ValidPreparationCreationRequestInput,
): Promise<ValidPreparation> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<ValidPreparation>>(backendRoutes.VALID_PREPARATIONS, input);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getValidPreparation(client: Axios, validPreparationID: string): Promise<ValidPreparation> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ValidPreparation>>(
      format(backendRoutes.VALID_PREPARATION, validPreparationID),
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getValidPreparations(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
): Promise<QueryFilteredResult<ValidPreparation>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ValidPreparation[]>>(backendRoutes.VALID_PREPARATIONS, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<ValidPreparation>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}

export async function updateValidPreparation(
  client: Axios,
  validPreparationID: string,
  input: ValidPreparationUpdateRequestInput,
): Promise<ValidPreparation> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<ValidPreparation>>(
      format(backendRoutes.VALID_PREPARATION, validPreparationID),
      input,
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function deleteValidPreparation(client: Axios, validPreparationID: string): Promise<ValidPreparation> {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete(format(backendRoutes.VALID_PREPARATION, validPreparationID));

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function searchForValidPreparations(client: Axios, query: string): Promise<ValidPreparation[]> {
  return new Promise(async function (resolve, reject) {
    const searchURL = `${backendRoutes.VALID_PREPARATIONS_SEARCH}?q=${encodeURIComponent(query)}`;
    const response = await client.get<APIResponse<ValidPreparation[]>>(searchURL);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}
