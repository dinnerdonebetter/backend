import { APIResponse, Page, Route } from '@playwright/test';
import { QueryFilter } from '@dinnerdonebetter/models';

export const clientName = 'DDB-Browser-Client';
export const clientHeaderName = 'X-Service-Client';

export interface RequestFulfillment {
  body?: string;
  contentType?: string;
  headers?: { [key: string]: string };
  path?: string;
  response?: APIResponse;
  status?: number;
}

export class ResponseConfig<T> {
  query: string;
  status: number;
  body?: T;
  filter: QueryFilter;
  times: number = 1;

  constructor(
    status: number = 200,
    query: string = '',
    filter: QueryFilter = QueryFilter.Default(),
    times: number = 1,
    body?: T,
  ) {
    this.status = status;
    if (body) this.body = body;
    this.query = query;
    this.filter = filter;
    this.times = times || 1;
  }

  fulfillRoute(route: Route): void {
    route.fulfill({
      contentType: 'application/json',
      headers: {
        'Content-Type': 'application/json',
        'X-Playwright-Mocked': 'true',
      },
      status: this.status,
      body: (this.status || 0) >= 200 && (this.status || 0) < 300 ? JSON.stringify(this.body) : '',
    });
  }

  fulfill(): RequestFulfillment {
    return {
      contentType: 'application/json',
      headers: {
        'Content-Type': 'application/json',
        'X-Playwright-Mocked': 'true',
      },
      status: this.status,
      body: (this.status || 0) >= 200 && (this.status || 0) < 300 ? JSON.stringify(this.body) : '',
    };
  }
}

export const unauthenticateAPI = (page: Page) =>
  page.route(`**/api/v1/*`, (route: Route) => {
    route.fulfill({ status: 401 });
  });

export function assertMethod(expectedMethod: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE', route: Route): void {
  const method = route.request().method();
  if (method !== expectedMethod) {
    throw new Error(`Unexpected request method: ${method}`);
  }
}

export function assertClient(route: Route) {
  route
    .request()
    .headerValue(clientHeaderName)
    .then((value) => {
      if (value && value !== clientName) {
        throw new Error(`Unexpected client: ${value}`);
      }
    });
}
