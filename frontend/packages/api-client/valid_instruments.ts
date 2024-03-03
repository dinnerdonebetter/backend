import { Axios } from 'axios';
import format from 'string-format';

import {
  ValidInstrumentCreationRequestInput,
  ValidInstrument,
  QueryFilter,
  ValidInstrumentUpdateRequestInput,
  QueryFilteredResult,
  APIResponse,
} from '@dinnerdonebetter/models';

import { backendRoutes } from './routes';

export async function createValidInstrument(
  client: Axios,
  input: ValidInstrumentCreationRequestInput,
): Promise<ValidInstrument> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<ValidInstrument>>(backendRoutes.VALID_INSTRUMENTS, input);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getValidInstrument(client: Axios, validInstrumentID: string): Promise<ValidInstrument> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ValidInstrument>>(
      format(backendRoutes.VALID_INSTRUMENT, validInstrumentID),
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getValidInstruments(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
): Promise<QueryFilteredResult<ValidInstrument>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ValidInstrument[]>>(backendRoutes.VALID_INSTRUMENTS, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<ValidInstrument>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}

export async function updateValidInstrument(
  client: Axios,
  validInstrumentID: string,
  input: ValidInstrumentUpdateRequestInput,
): Promise<ValidInstrument> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<ValidInstrument>>(
      format(backendRoutes.VALID_INSTRUMENT, validInstrumentID),
      input,
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function deleteValidInstrument(client: Axios, validInstrumentID: string): Promise<ValidInstrument> {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete<APIResponse<ValidInstrument>>(
      format(backendRoutes.VALID_INSTRUMENT, validInstrumentID),
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function searchForValidInstruments(client: Axios, query: string): Promise<ValidInstrument[]> {
  return new Promise(async function (resolve, reject) {
    const searchURL = `${backendRoutes.VALID_INSTRUMENTS_SEARCH}?q=${encodeURIComponent(query)}`;
    const response = await client.get<APIResponse<ValidInstrument[]>>(searchURL);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}
