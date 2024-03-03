import type { Page, Route } from '@playwright/test';

import { QueryFilteredResult, User } from '@dinnerdonebetter/models';
import { spellWord } from './utils';
import { assertClient, assertMethod, methods, ResponseConfig } from './helpers';

export class MockUserResponseConfig extends ResponseConfig<User> {
  userID: string;

  constructor(userID: string, status: number = 200, body?: User) {
    super();

    this.userID = userID;
    this.status = status;
    this.body = body;
  }
}

export const mockSelf = (resCfg: ResponseConfig<User>) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/users/self`,
      (route: Route) => {
        const req = route.request();

        assertMethod(methods.GET, route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};

export const mockUser = (resCfg: MockUserResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/users/${resCfg.userID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod(methods.GET, route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};

export class MockUsersListResponseConfig extends ResponseConfig<QueryFilteredResult<User>> {}

export const mockUsersList = (resCfg: MockUsersListResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/users?${resCfg.filter.asURLSearchParams()}`,
      (route: Route) => {
        const req = route.request();

        assertMethod(methods.GET, route);
        assertClient(route);

        if (resCfg.body && resCfg.filter) resCfg.body.limit = resCfg.filter.limit;
        if (resCfg.body && resCfg.filter) resCfg.body.page = resCfg.filter.page;

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};

export class MockUsersSearchResponseConfig extends ResponseConfig<User[]> {}

export const mockUsersSearch = (resCfg: MockUsersSearchResponseConfig) => {
  return (page: Page) => {
    for (const word of spellWord(resCfg.query)) {
      page.route(
        `**/api/v1/users/search?q=${encodeURIComponent(word)}`,
        (route: Route) => {
          const req = route.request();

          assertMethod(methods.GET, route);
          assertClient(route);

          const rv = resCfg.fulfill();

          if (word !== resCfg.query) {
            rv.body = JSON.stringify([]);
          }

          route.fulfill(rv);
        },
        { times: resCfg.times },
      );
    }
  };
};

export class MockUserReputationUpdateResponseConfig extends ResponseConfig<User> {}

export const mockUserReputationUpdate = (resCfg: MockUserReputationUpdateResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/admin/users/status`,
      (route: Route) => {
        const req = route.request();

        assertMethod(methods.GET, route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
