import { Axios } from 'axios';
import format from 'string-format';

import {
  WebhookCreationRequestInput,
  Webhook,
  QueryFilter,
  QueryFilteredResult,
  APIResponse,
} from '@dinnerdonebetter/models';

import { backendRoutes } from './routes';

export async function createWebhook(client: Axios, input: WebhookCreationRequestInput): Promise<Webhook> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<Webhook>>(backendRoutes.VALID_PREPARATIONS, input);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getWebhook(client: Axios, WebhookID: string): Promise<Webhook> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<Webhook>>(format(backendRoutes.VALID_PREPARATION, WebhookID));

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getWebhooks(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
): Promise<QueryFilteredResult<Webhook>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<Webhook[]>>(backendRoutes.VALID_PREPARATIONS, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<Webhook>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      filteredCount: response.data.pagination?.filteredCount,
    });

    resolve(result);
  });
}

export async function deleteWebhook(client: Axios, WebhookID: string): Promise<Webhook> {
  return new Promise(async function (resolve, reject) {
    const response = await client.delete<APIResponse<Webhook>>(format(backendRoutes.VALID_PREPARATION, WebhookID));

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function searchForWebhooks(client: Axios, query: string): Promise<Webhook[]> {
  return new Promise(async function (resolve, reject) {
    const searchURL = `${backendRoutes.VALID_PREPARATIONS_SEARCH}?q=${encodeURIComponent(query)}`;
    const response = await client.get<APIResponse<Webhook[]>>(searchURL);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}
