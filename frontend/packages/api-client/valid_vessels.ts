import { Axios } from 'axios';
import format from 'string-format';

import {
  ValidVesselCreationRequestInput,
  ValidVessel,
  QueryFilter,
  ValidVesselUpdateRequestInput,
  QueryFilteredResult,
  APIResponse,
} from '@dinnerdonebetter/models';

import { backendRoutes } from './routes';

export async function createValidVessel(client: Axios, input: ValidVesselCreationRequestInput): Promise<ValidVessel> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<ValidVessel>>(backendRoutes.VALID_VESSELS, input);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getValidVessel(client: Axios, validVesselID: string): Promise<ValidVessel> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ValidVessel>>(format(backendRoutes.VALID_VESSEL, validVesselID));

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getValidVessels(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
): Promise<QueryFilteredResult<ValidVessel>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ValidVessel[]>>(backendRoutes.VALID_VESSELS, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<ValidVessel>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      filteredCount: response.data.pagination?.filteredCount,
    });

    resolve(result);
  });
}

export async function updateValidVessel(
  client: Axios,
  validVesselID: string,
  input: ValidVesselUpdateRequestInput,
): Promise<ValidVessel> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<ValidVessel>>(
      format(backendRoutes.VALID_VESSEL, validVesselID),
      input,
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function deleteValidVessel(client: Axios, validVesselID: string): Promise<ValidVessel> {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete<APIResponse<ValidVessel>>(format(backendRoutes.VALID_VESSEL, validVesselID));

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function searchForValidVessels(client: Axios, query: string): Promise<ValidVessel[]> {
  return new Promise(async function (resolve, reject) {
    const searchURL = `${backendRoutes.VALID_VESSELS_SEARCH}?q=${encodeURIComponent(query)}`;
    const response = await client.get<APIResponse<ValidVessel[]>>(searchURL);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}
