import { test, expect, Page, Route } from '@playwright/test';

import { UserRegistrationInput } from '@dinnerdonebetter/models';

export const fakeRegister = (username: string, password: string, emailAddress: string, errorMessage: string = '') => {
  return (page: Page) =>
    page.route('**/users/', (route: Route) => {
      const requestBody = route.request().postDataJSON() as UserRegistrationInput;

      expect(requestBody.username).toEqual(username);
      expect(requestBody.password).toEqual(password);
      expect(requestBody.emailAddress).toEqual(emailAddress);

      route.fulfill({
        status: errorMessage === '' ? 201 : 400,
        body: errorMessage === '' ? '' : JSON.stringify({ message: errorMessage }),
      });
    });
};

test('register account', async ({ page }) => {
  await page.goto('/register');

  const expectedEmailAddress = 'things@stuff.com';
  const expectedUsername = 'username';
  const expectedPassword = 'hunterz2';

  await fakeRegister(expectedUsername, expectedPassword, expectedEmailAddress)(page);

  await expect(page.locator('.errors-container')).toHaveCount(0);

  const emailInput = await page.locator('input[data-qa="registration-email-address-input"]');
  await expect(emailInput).toBeEnabled();
  await emailInput.type(expectedEmailAddress);

  const usernameInput = await page.locator('input[data-qa="registration-username-input"]');
  await expect(usernameInput).toBeEnabled();
  await usernameInput.type(expectedUsername);

  const passwordInput = await page.locator('input[data-qa="registration-password-input"][type="password"]');
  await expect(passwordInput).toBeEnabled();
  await passwordInput.type(expectedPassword);

  const passwordConfirmInput = await page.locator(
    'input[data-qa="registration-password-confirm-input"][type="password"]',
  );
  await expect(passwordConfirmInput).toBeEnabled();
  await passwordConfirmInput.type(expectedPassword);

  const submitButton = await page.locator('button[data-qa="registration-button"]');
  await submitButton.click();

  await page.waitForURL('/login');
});
