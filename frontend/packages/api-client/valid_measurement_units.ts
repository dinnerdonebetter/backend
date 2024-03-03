import { Axios } from 'axios';
import format from 'string-format';

import {
  ValidMeasurementUnitCreationRequestInput,
  ValidMeasurementUnit,
  QueryFilter,
  ValidMeasurementUnitUpdateRequestInput,
  QueryFilteredResult,
  APIResponse,
} from '@dinnerdonebetter/models';

import { backendRoutes } from './routes';

export async function createValidMeasurementUnit(
  client: Axios,
  input: ValidMeasurementUnitCreationRequestInput,
): Promise<ValidMeasurementUnit> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<ValidMeasurementUnit>>(backendRoutes.VALID_MEASUREMENT_UNITS, input);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getValidMeasurementUnit(
  client: Axios,
  validMeasurementUnitID: string,
): Promise<ValidMeasurementUnit> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ValidMeasurementUnit>>(
      format(backendRoutes.VALID_MEASUREMENT_UNIT, validMeasurementUnitID),
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getValidMeasurementUnits(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
): Promise<QueryFilteredResult<ValidMeasurementUnit>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ValidMeasurementUnit[]>>(backendRoutes.VALID_MEASUREMENT_UNITS, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<ValidMeasurementUnit>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}

export async function updateValidMeasurementUnit(
  client: Axios,
  validMeasurementUnitID: string,
  input: ValidMeasurementUnitUpdateRequestInput,
): Promise<ValidMeasurementUnit> {
  return new Promise(async function (resolve, reject) {
    const repsonse = await client.put<APIResponse<ValidMeasurementUnit>>(
      format(backendRoutes.VALID_MEASUREMENT_UNIT, validMeasurementUnitID),
      input,
    );

    if (repsonse.data.error) {
      reject(new Error(repsonse.data.error.message));
    }

    resolve(repsonse.data.data);
  });
}

export async function deleteValidMeasurementUnit(
  client: Axios,
  validMeasurementUnitID: string,
): Promise<ValidMeasurementUnit> {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete<APIResponse<ValidMeasurementUnit>>(
      format(backendRoutes.VALID_MEASUREMENT_UNIT, validMeasurementUnitID),
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function searchForValidMeasurementUnits(client: Axios, query: string): Promise<ValidMeasurementUnit[]> {
  return new Promise(async function (resolve, reject) {
    const uri = `${backendRoutes.VALID_MEASUREMENT_UNITS_SEARCH}?q=${encodeURIComponent(query)}`;
    const response = await client.get<APIResponse<ValidMeasurementUnit[]>>(uri);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function searchForValidMeasurementUnitsByIngredientID(
  client: Axios,
  validIngredientID: string,
  filter: QueryFilter = QueryFilter.Default(),
): Promise<QueryFilteredResult<ValidMeasurementUnit>> {
  return new Promise(async function (resolve, reject) {
    const uri = format(backendRoutes.VALID_MEASUREMENT_UNITS_SEARCH_BY_INGREDIENT, validIngredientID);
    const response = await client.get<APIResponse<ValidMeasurementUnit[]>>(uri, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<ValidMeasurementUnit>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}
