import { expect, Page } from '@playwright/test';

export enum KeyPresses {
  ArrowUp = 'ArrowUp',
  ArrowDown = 'ArrowDown',
  ArrowLeft = 'ArrowLeft',
  ArrowRight = 'ArrowRight',
  Enter = 'Enter',
  Escape = 'Escape',
  Backspace = 'Backspace',
  Tab = 'Tab',
  Space = ' ',
  Shift = 'Shift',
  Control = 'Control',
  Alt = 'Alt',
  Meta = 'Meta',
  Delete = 'Delete',
  Insert = 'Insert',
  Home = 'Home',
  End = 'End',
  PageUp = 'PageUp',
  PageDown = 'PageDown',
  F1 = 'F1',
  F2 = 'F2',
  F3 = 'F3',
  F4 = 'F4',
  F5 = 'F5',
  F6 = 'F6',
  F7 = 'F7',
  F8 = 'F8',
  F9 = 'F9',
  F10 = 'F10',
  F11 = 'F11',
  F12 = 'F12',

  // Keypad keys
  Numpad0 = '0',
  Numpad1 = '1',
  Numpad2 = '2',
  Numpad3 = '3',
  Numpad4 = '4',
  Numpad5 = '5',
  Numpad6 = '6',
  Numpad7 = '7',
  Numpad8 = '8',
  Numpad9 = '9',
  NumpadAdd = 'Add',
  NumpadSubtract = 'Subtract',
  NumpadMultiply = 'Multiply',
  NumpadDivide = 'Divide',
  NumpadDecimal = 'Decimal',
}

type validElementTag =
  | 'div'
  | 'img'
  | 'span'
  | 'input'
  | 'tr'
  | 'td'
  | 'th'
  | 'button'
  | 'a'
  | 'select'
  | 'option'
  | 'textarea'
  | 'h1'
  | 'h2'
  | 'h3'
  | 'h4'
  | 'h5'
  | 'h6';

export const selector = (element: validElementTag, selector: string, addendum: string = ''): string => {
  return `${element}[data-qa="${selector}"]${addendum}`;
};

export const sleep = async (millis: number) => await new Promise((resolve) => setTimeout(resolve, millis));

export const checkDropdown = async (page: Page, sel: string, checkChild: boolean = true): Promise<void> => {
  await expect(page.locator(selector('div', sel))).toBeVisible();
  await page.locator(selector('div', sel)).hover();
  if (checkChild) {
    await expect(page.locator(selector('a', `${sel}-action-0`))).toBeVisible();
  }
};
