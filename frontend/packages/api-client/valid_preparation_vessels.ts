import { Axios } from 'axios';
import format from 'string-format';

import {
  ValidPreparationVesselCreationRequestInput,
  ValidPreparationVessel,
  QueryFilteredResult,
  APIResponse,
} from '@dinnerdonebetter/models';

import { backendRoutes } from './routes';

export async function getValidPreparationVessel(
  client: Axios,
  validPreparationVesselID: string,
): Promise<ValidPreparationVessel> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ValidPreparationVessel>>(
      format(backendRoutes.VALID_PREPARATION_VESSEL, validPreparationVesselID),
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function validPreparationVesselsForPreparationID(
  client: Axios,
  validPreparationID: string,
): Promise<QueryFilteredResult<ValidPreparationVessel>> {
  return new Promise(async function (resolve, reject) {
    const searchURL = format(backendRoutes.VALID_PREPARATION_VESSELS_SEARCH_BY_PREPARATION_ID, validPreparationID);
    const response = await client.get<APIResponse<ValidPreparationVessel[]>>(searchURL);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<ValidPreparationVessel>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}

export async function validPreparationVesselsForVesselID(
  client: Axios,
  validVesselID: string,
): Promise<QueryFilteredResult<ValidPreparationVessel>> {
  return new Promise(async function (resolve, reject) {
    const searchURL = format(backendRoutes.VALID_PREPARATION_VESSELS_SEARCH_BY_VESSEL_ID, validVesselID);
    const response = await client.get<APIResponse<ValidPreparationVessel[]>>(searchURL);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<ValidPreparationVessel>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}

export async function createValidPreparationVessel(
  client: Axios,
  input: ValidPreparationVesselCreationRequestInput,
): Promise<ValidPreparationVessel> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<ValidPreparationVessel>>(
      backendRoutes.VALID_PREPARATION_VESSELS,
      input,
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function deleteValidPreparationVessel(
  client: Axios,
  validPreparationVesselID: string,
): Promise<ValidPreparationVessel> {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete<APIResponse<ValidPreparationVessel>>(
      format(backendRoutes.VALID_PREPARATION_VESSEL, validPreparationVesselID),
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}
