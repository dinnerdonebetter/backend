import { Axios } from 'axios';
import format from 'string-format';

import {
  ServiceSettingCreationRequestInput,
  ServiceSetting,
  QueryFilter,
  ServiceSettingUpdateRequestInput,
  QueryFilteredResult,
  APIResponse,
} from '@dinnerdonebetter/models';

import { backendRoutes } from './routes';

export async function createServiceSetting(
  client: Axios,
  input: ServiceSettingCreationRequestInput,
): Promise<ServiceSetting> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<ServiceSetting>>(backendRoutes.SERVICE_SETTINGS, input);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getServiceSetting(client: Axios, serviceSettingID: string): Promise<ServiceSetting> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ServiceSetting>>(
      format(backendRoutes.SERVICE_SETTING, serviceSettingID),
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getServiceSettings(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
): Promise<QueryFilteredResult<ServiceSetting>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ServiceSetting[]>>(backendRoutes.SERVICE_SETTINGS, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<ServiceSetting>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      filteredCount: response.data.pagination?.filteredCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}

export async function updateServiceSetting(
  client: Axios,
  serviceSettingID: string,
  input: ServiceSettingUpdateRequestInput,
): Promise<ServiceSetting> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<ServiceSetting>>(
      format(backendRoutes.SERVICE_SETTING, serviceSettingID),
      input,
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function deleteServiceSetting(client: Axios, serviceSettingID: string): Promise<ServiceSetting> {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete<APIResponse<ServiceSetting>>(
      format(backendRoutes.SERVICE_SETTING, serviceSettingID),
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function searchForServiceSettings(client: Axios, query: string): Promise<ServiceSetting[]> {
  return new Promise(async function (resolve, reject) {
    const searchURL = `${backendRoutes.SERVICE_SETTINGS_SEARCH}?q=${encodeURIComponent(query)}`;
    const response = await client.get<APIResponse<ServiceSetting[]>>(searchURL);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}
