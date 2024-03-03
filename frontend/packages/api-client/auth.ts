import { Axios, AxiosResponse } from 'axios';

import {
  UserLoginInput,
  UserStatusResponse,
  UserRegistrationInput,
  UserPermissionsRequestInput,
  UserPermissionsResponse,
  PasswordResetTokenCreationRequestInput,
  PasswordResetTokenRedemptionRequestInput,
  UsernameReminderRequestInput,
  UserCreationResponse,
  PasswordUpdateInput,
  APIResponse,
} from '@dinnerdonebetter/models';

import { backendRoutes } from './routes';

export async function logIn(client: Axios, input: UserLoginInput): Promise<AxiosResponse<UserStatusResponse>> {
  return client.post<UserStatusResponse>(backendRoutes.LOGIN, input).then((response) => {
    if (response.status !== 202 && response.status !== 205) {
      throw new Error(`Unexpected response status: ${response.status}`);
    }

    return response;
  });
}

export async function adminLogin(client: Axios, input: UserLoginInput): Promise<AxiosResponse<UserStatusResponse>> {
  return client.post<UserStatusResponse>(backendRoutes.LOGIN_ADMIN, input).then((response) => {
    if (response.status !== 202) {
      throw new Error(`Unexpected response status: ${response.status}`);
    }

    return response;
  });
}

export async function logOut(client: Axios): Promise<AxiosResponse<UserStatusResponse>> {
  return client.post(backendRoutes.LOGOUT);
}

export async function register(client: Axios, input: UserRegistrationInput): Promise<UserCreationResponse> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<UserCreationResponse>>(backendRoutes.USER_REGISTRATION, input);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function checkPermissions(
  client: Axios,
  body: UserPermissionsRequestInput,
): Promise<UserPermissionsResponse> {
  return new Promise(async function (resolve, reject) {
    const response = await client.post<APIResponse<UserPermissionsResponse>>(backendRoutes.PERMISSIONS_CHECK, body);

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    resolve(response.data.data);
  });
}

export async function requestPasswordResetToken(
  client: Axios,
  input: PasswordResetTokenCreationRequestInput,
): Promise<AxiosResponse> {
  return client.post(backendRoutes.REQUEST_PASSWORD_RESET_TOKEN, input);
}

export async function changePassword(client: Axios, input: PasswordUpdateInput): Promise<AxiosResponse> {
  return client.put(backendRoutes.CHANGE_PASSWORD, input);
}

export async function redeemPasswordResetToken(
  client: Axios,
  input: PasswordResetTokenRedemptionRequestInput,
): Promise<AxiosResponse> {
  return client.post(backendRoutes.REDEEM_PASSWORD_RESET_TOKEN, input);
}

export async function requestUsernameReminderEmail(
  client: Axios,
  input: UsernameReminderRequestInput,
): Promise<AxiosResponse> {
  return client.post(backendRoutes.REQUEST_USERNAME_REMINDER_EMAIL, input);
}
