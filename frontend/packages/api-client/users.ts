import { Axios } from 'axios';
import format from 'string-format';

import {
  User,
  QueryFilter,
  UserAccountStatusUpdateInput,
  QueryFilteredResult,
  EmailAddressVerificationRequestInput,
  AvatarUpdateInput,
  TOTPSecretRefreshInput,
  APIResponse,
} from '@dinnerdonebetter/models';

import { backendRoutes } from './routes';

export async function fetchSelf(client: Axios): Promise<User> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<User>>(backendRoutes.SELF);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function requestEmailVerificationEmail(client: Axios): Promise<User> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post(backendRoutes.USERS_REQUEST_EMAIL_VERIFICATION_EMAIL);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getUser(client: Axios, userID: string): Promise<User> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<User>>(format(backendRoutes.USER, userID));

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function getUsers(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
): Promise<QueryFilteredResult<User>> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<User[]>>(backendRoutes.USERS, { params: filter.asRecord() });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<User>({
      data: response.data.data,
      filteredCount: response.data.pagination?.filteredCount,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}

export async function updateUserAccountStatus(client: Axios, input: UserAccountStatusUpdateInput): Promise<User> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<User>>(backendRoutes.USER_REPUTATION_UPDATE, input);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function searchForUsers(client: Axios, query: string): Promise<User[]> {
  return new Promise(async function (resolve, reject) {
    const response = await client.get<APIResponse<User[]>>(
      `${backendRoutes.USERS_SEARCH}?q=${encodeURIComponent(query)}`,
    );

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function verifyEmailAddress(
  client: Axios,
  verificationInput: EmailAddressVerificationRequestInput,
): Promise<User> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<User>>(backendRoutes.USERS_VERIFY_EMAIL_ADDRESS, verificationInput);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function uploadNewAvatar(client: Axios, input: AvatarUpdateInput): Promise<User> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<User>>(backendRoutes.USERS_UPLOAD_NEW_AVATAR, input);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function newTwoFactorSecret(client: Axios, input: TOTPSecretRefreshInput): Promise<User> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<User>>(backendRoutes.NEW_2FA_SECRET, input);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}
