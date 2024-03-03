import { Axios } from 'axios';
import format from 'string-format';

import {
  Household,
  QueryFilter,
  HouseholdUpdateRequestInput,
  HouseholdInvitationCreationRequestInput,
  HouseholdInvitation,
  QueryFilteredResult,
  APIResponse,
} from '@dinnerdonebetter/models';

import { backendRoutes } from './routes';

export async function getCurrentHouseholdInfo(client: Axios): Promise<Household> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<Household>>(backendRoutes.HOUSEHOLD_INFO);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getHousehold(client: Axios, id: string): Promise<Household> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<Household>>(format(backendRoutes.HOUSEHOLD, id));

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getHouseholds(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
): Promise<QueryFilteredResult<Household>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<Household[]>>(backendRoutes.HOUSEHOLDS, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<Household>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      filteredCount: response.data.pagination?.filteredCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}

export async function updateHousehold(
  client: Axios,
  householdID: string,
  household: HouseholdUpdateRequestInput,
): Promise<Household> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<Household>>(format(backendRoutes.HOUSEHOLD, householdID), household);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function inviteUserToHousehold(
  client: Axios,
  householdID: string,
  input: HouseholdInvitationCreationRequestInput,
): Promise<HouseholdInvitation> {
  return new Promise(async function (resolve, reject) {
    const uri = format(backendRoutes.HOUSEHOLD_ADD_MEMBER, householdID);

    const response = await client.post<APIResponse<HouseholdInvitation>>(uri, input);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function removeMemberFromHousehold(
  client: Axios,
  householdID: string,
  memberID: string,
): Promise<Household> {
  return new Promise(async function (resolve, reject) {
    const uri = format(backendRoutes.HOUSEHOLD_REMOVE_MEMBER, householdID, memberID);
    const response = await client.delete<APIResponse<Household>>(uri);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getReceivedInvites(client: Axios): Promise<QueryFilteredResult<HouseholdInvitation>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<HouseholdInvitation[]>>(backendRoutes.RECEIVED_HOUSEHOLD_INVITATIONS);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<HouseholdInvitation>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      filteredCount: response.data.pagination?.filteredCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}

export async function getSentInvites(client: Axios): Promise<QueryFilteredResult<HouseholdInvitation>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<HouseholdInvitation[]>>(backendRoutes.SENT_HOUSEHOLD_INVITATIONS);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<HouseholdInvitation>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      filteredCount: response.data.pagination?.filteredCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}
