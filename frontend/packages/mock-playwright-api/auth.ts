import type { Page, Route } from '@playwright/test';

import {
  UserStatusResponse,
  UserLoginInput,
  UserRegistrationInput,
  UserPermissionsResponse,
} from '@dinnerdonebetter/models';
import { assertClient, assertMethod, ResponseConfig } from './helpers';

export const mockLogin = (username: string, password: string, totpToken: string) => {
  return (page: Page) =>
    page.route('**/users/login', (route: Route) => {
      assertMethod('POST', route);
      assertClient(route);

      const authState = new UserStatusResponse({
        accountStatus: 'good',
        accountStatusExplanation: '',
        activeHousehold: '',
        isAuthenticated: true,
      });

      const requestBody = route.request().postDataJSON() as UserLoginInput;
      const totpTokenIsEmpty = requestBody.totpToken === '';

      if (requestBody.username !== username) {
        throw new Error('username does not match');
      }
      if (requestBody.password !== password) {
        throw new Error('password does not match');
      }

      if (!totpTokenIsEmpty) {
        if (requestBody.totpToken !== totpToken) {
          throw new Error('TOTP token does not match');
        }
      }

      route.fulfill({
        status: totpTokenIsEmpty ? 205 : 202,
        body: totpTokenIsEmpty ? '' : JSON.stringify(authState),
      });
    });
};

export const mockAdminLogin = (username: string, password: string, totpToken: string) => {
  return (page: Page) =>
    page.route('**/users/login/admin', (route: Route) => {
      const req = route.request();

      assertMethod('POST', route);
      assertClient(route);

      const authState = new UserStatusResponse({
        accountStatus: 'good',
        accountStatusExplanation: '',
        activeHousehold: '',
        isAuthenticated: true,
      });

      const requestBody = route.request().postDataJSON() as UserLoginInput;
      const totpTokenIsEmpty = requestBody.totpToken === '';

      if (requestBody.username !== username) {
        throw new Error('username does not match');
      }
      if (requestBody.password !== password) {
        throw new Error('password does not match');
      }

      if (!totpTokenIsEmpty) {
        if (requestBody.totpToken !== totpToken) {
          throw new Error('TOTP token does not match');
        }
      }

      route.fulfill({
        status: totpTokenIsEmpty ? 205 : 202,
        body: totpTokenIsEmpty ? '' : JSON.stringify(authState),
      });
    });
};

export const mockRegister = (emailAddress: string, username: string, password: string) => {
  return (page: Page) =>
    page.route(
      '**/users/register',
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        const requestBody = route.request().postDataJSON() as UserRegistrationInput;
        if (requestBody.emailAddress !== emailAddress) {
          throw new Error('username does not match');
        }
        if (requestBody.username !== username) {
          throw new Error('username does not match');
        }
        if (requestBody.password !== password) {
          throw new Error('password does not match');
        }

        route.fulfill({
          status: 201,
        });
      },
      { times: 1 },
    );
};

export const mockLogout = (statusCode = 303) => {
  return (page: Page) =>
    page.route(
      '**/users/logout',
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill({
          status: statusCode,
        });
      },
      { times: 1 },
    );
};

export class MockUserPermissionsResponseConfig extends ResponseConfig<UserPermissionsResponse> {}

export const mockCheckPermissions = (resCfg: MockUserPermissionsResponseConfig) => {
  return (page: Page) =>
    page.route(
      '**/users/permissions/check',
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};

export const mockRequestPasswordResetToken = (statusCode = 202) => {
  return (page: Page) =>
    page.route(
      '**/users/password/reset',
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill({
          status: statusCode,
        });
      },
      { times: 1 },
    );
};

export const mockRedeemPasswordResetToken = (statusCode = 202) => {
  return (page: Page) =>
    page.route(
      '**/users/password/reset/redeem',
      (route: Route) => {
        const req = route.request();

        assertMethod('POST', route);
        assertClient(route);

        route.fulfill({
          status: statusCode,
        });
      },
      { times: 1 },
    );
};
