import { Axios } from 'axios';
import format from 'string-format';

import { APIResponse, HouseholdInvitation, HouseholdInvitationUpdateRequestInput } from '@dinnerdonebetter/models';

import { backendRoutes } from './routes';

export async function getInvitation(client: Axios, invitationID: string): Promise<HouseholdInvitation> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<HouseholdInvitation>>(
      format(backendRoutes.HOUSEHOLD_INVITATION, invitationID),
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function acceptInvitation(
  client: Axios,
  invitationID: string,
  input: HouseholdInvitationUpdateRequestInput,
): Promise<HouseholdInvitation> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<HouseholdInvitation>>(
      format(backendRoutes.ACCEPT_HOUSEHOLD_INVITATION, invitationID),
      input,
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function cancelInvitation(
  client: Axios,
  invitationID: string,
  input: HouseholdInvitationUpdateRequestInput,
): Promise<HouseholdInvitation> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<HouseholdInvitation>>(
      format(backendRoutes.CANCEL_HOUSEHOLD_INVITATION, invitationID),
      input,
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function rejectInvitation(
  client: Axios,
  invitationID: string,
  input: HouseholdInvitationUpdateRequestInput,
): Promise<HouseholdInvitation> {
  return new Promise(async function (resolve, reject) {
    const response = await client.put<APIResponse<HouseholdInvitation>>(
      format(backendRoutes.REJECT_HOUSEHOLD_INVITATION, invitationID),
      input,
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}
