import { Axios } from 'axios';
import format from 'string-format';

import {
  ValidMeasurementUnitConversionCreationRequestInput,
  ValidMeasurementUnitConversion,
  QueryFilter,
  ValidMeasurementUnitConversionUpdateRequestInput,
  QueryFilteredResult,
  APIResponse,
} from '@dinnerdonebetter/models';

import { backendRoutes } from './routes';

export async function createValidMeasurementUnitConversion(
  client: Axios,
  input: ValidMeasurementUnitConversionCreationRequestInput,
): Promise<ValidMeasurementUnitConversion> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<ValidMeasurementUnitConversion>>(
      backendRoutes.VALID_MEASUREMENT_UNIT_CONVERSIONS,
      input,
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getValidMeasurementUnitConversion(
  client: Axios,
  validMeasurementUnitConversionID: string,
): Promise<ValidMeasurementUnitConversion> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ValidMeasurementUnitConversion>>(
      format(backendRoutes.VALID_MEASUREMENT_UNIT_CONVERSION, validMeasurementUnitConversionID),
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getValidMeasurementUnitConversionsFromUnit(
  client: Axios,
  validMeasurementUnitID: string,
  filter: QueryFilter = QueryFilter.Default(),
): Promise<ValidMeasurementUnitConversion[]> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ValidMeasurementUnitConversion[]>>(
      format(backendRoutes.VALID_MEASUREMENT_UNIT_CONVERSIONS_FROM_UNIT, validMeasurementUnitID),
      {
        params: filter.asRecord(),
      },
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getValidMeasurementUnitConversionsToUnit(
  client: Axios,
  validMeasurementUnitID: string,
  filter: QueryFilter = QueryFilter.Default(),
): Promise<ValidMeasurementUnitConversion[]> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ValidMeasurementUnitConversion[]>>(
      format(backendRoutes.VALID_MEASUREMENT_UNIT_CONVERSIONS_TO_UNIT, validMeasurementUnitID),
      {
        params: filter.asRecord(),
      },
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getValidMeasurementUnitConversions(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
): Promise<QueryFilteredResult<ValidMeasurementUnitConversion>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ValidMeasurementUnitConversion[]>>(
      backendRoutes.VALID_MEASUREMENT_UNIT_CONVERSIONS,
      {
        params: filter.asRecord(),
      },
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<ValidMeasurementUnitConversion>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}

export async function updateValidMeasurementUnitConversion(
  client: Axios,
  validMeasurementUnitConversionID: string,
  input: ValidMeasurementUnitConversionUpdateRequestInput,
): Promise<ValidMeasurementUnitConversion> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<ValidMeasurementUnitConversion>>(
      format(backendRoutes.VALID_MEASUREMENT_UNIT_CONVERSION, validMeasurementUnitConversionID),
      input,
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function deleteValidMeasurementUnitConversion(
  client: Axios,
  validMeasurementUnitConversionID: string,
): Promise<ValidMeasurementUnitConversion> {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete<APIResponse<ValidMeasurementUnitConversion>>(
      format(backendRoutes.VALID_MEASUREMENT_UNIT_CONVERSION, validMeasurementUnitConversionID),
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}
