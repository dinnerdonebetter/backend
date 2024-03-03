import { Page, Route } from '@playwright/test';
import { HouseholdInvitation, QueryFilteredResult } from '@dinnerdonebetter/models';
import { assertClient, assertMethod, ResponseConfig } from './helpers';

export class MockHouseholdInvitationResponseConfig extends ResponseConfig<HouseholdInvitation> {
  invitationID: string;

  constructor(invitationID: string, status: number = 200, body?: HouseholdInvitation) {
    super();

    this.invitationID = invitationID;
    this.status = status;
    this.body = body;
  }
}

export const mockGetHouseholdInvitation = (resCfg: MockHouseholdInvitationResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/household_invitations/${resCfg.invitationID}`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};

export class MockHouseholdInvitationStatusChangeResponseConfig extends ResponseConfig<void> {
  invitationID: string;

  constructor(invitationID: string, status: number = 200) {
    super();

    this.invitationID = invitationID;
    this.status = status;
  }
}

export const mockAcceptHouseholdInvitation = (resCfg: MockHouseholdInvitationStatusChangeResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/household_invitations/${resCfg.invitationID}/accept`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PUT', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};

export const mockCancelHouseholdInvitation = (resCfg: MockHouseholdInvitationStatusChangeResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/household_invitations/${resCfg.invitationID}/cancel`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PUT', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};

export const mockRejectHouseholdInvitation = (resCfg: MockHouseholdInvitationStatusChangeResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/household_invitations/${resCfg.invitationID}/reject`,
      (route: Route) => {
        const req = route.request();

        assertMethod('PUT', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};

export class MockPendingInvitationsResponseConfig extends ResponseConfig<QueryFilteredResult<HouseholdInvitation>> {}

export const mockPendingInvitations = (resCfg: MockPendingInvitationsResponseConfig) => {
  return (page: Page) =>
    page.route(
      `**/api/v1/household_invitations/sent`,
      (route: Route) => {
        const req = route.request();

        assertMethod('GET', route);
        assertClient(route);

        route.fulfill(resCfg.fulfill());
      },
      { times: resCfg.times },
    );
};
