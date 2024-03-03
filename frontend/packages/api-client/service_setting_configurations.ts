import { Axios } from 'axios';
import format from 'string-format';

import {
  ServiceSettingConfigurationCreationRequestInput,
  ServiceSettingConfiguration,
  ServiceSettingConfigurationUpdateRequestInput,
  QueryFilteredResult,
  APIResponse,
} from '@dinnerdonebetter/models';

import { backendRoutes } from './routes';

export async function createServiceSettingConfiguration(
  client: Axios,
  input: ServiceSettingConfigurationCreationRequestInput,
): Promise<ServiceSettingConfiguration> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<ServiceSettingConfiguration>>(backendRoutes.SERVICE_SETTINGS, input);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getServiceSettingConfigurationsForUser(
  client: Axios,
): Promise<QueryFilteredResult<ServiceSettingConfiguration>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ServiceSettingConfiguration[]>>(
      backendRoutes.SERVICE_SETTING_CONFIGURATIONS_FOR_USER,
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<ServiceSettingConfiguration>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      filteredCount: response.data.pagination?.filteredCount,
      limit: response.data.pagination?.limit,
      page: response.data.pagination?.page,
    });

    resolve(result);
  });
}

export async function getServiceSettingConfigurationsForHousehold(
  client: Axios,
): Promise<QueryFilteredResult<ServiceSettingConfiguration>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<ServiceSettingConfiguration[]>>(
      backendRoutes.SERVICE_SETTING_CONFIGURATIONS_FOR_HOUSEHOLD,
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<ServiceSettingConfiguration>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      filteredCount: response.data.pagination?.filteredCount,
      limit: response.data.pagination?.limit,
      page: response.data.pagination?.page,
    });

    resolve(result);
  });
}

export async function updateServiceSettingConfiguration(
  client: Axios,
  serviceSettingConfigurationID: string,
  input: ServiceSettingConfigurationUpdateRequestInput,
): Promise<ServiceSettingConfiguration> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<ServiceSettingConfiguration>>(
      format(backendRoutes.SERVICE_SETTING_CONFIGURATION, serviceSettingConfigurationID),
      input,
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function deleteServiceSettingConfiguration(
  client: Axios,
  serviceSettingConfigurationID: string,
): Promise<ServiceSettingConfiguration> {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete<APIResponse<ServiceSettingConfiguration>>(
      format(backendRoutes.SERVICE_SETTING_CONFIGURATION, serviceSettingConfigurationID),
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}
