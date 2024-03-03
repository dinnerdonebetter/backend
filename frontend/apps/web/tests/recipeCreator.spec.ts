import { test, expect } from '@playwright/test';

import { Recipe } from '@dinnerdonebetter/models';

test('recipe creator functions as expected', async ({ page }) => {
  await page.goto('/recipes/new');

  const exampleRecipe = new Recipe({
    name: 'Test recipe',
    description: 'test recipe description',
    source: 'https://www.grandmascookbook.pizza',
    minimumEstimatedPortions: 4,
  });

  const recipeNameInput = await page.locator('input[data-qa="recipe-name-input"]');
  await expect(recipeNameInput).toBeEnabled();
  await recipeNameInput.type(exampleRecipe.name);

  const recipeDescriptionInput = await page.locator('textarea[data-qa="recipe-description-input"]');
  await expect(recipeDescriptionInput).toBeEnabled();
  await recipeDescriptionInput.type(exampleRecipe.description);

  const recipeSourceInput = await page.locator('input[data-qa="recipe-source-input"]');
  await expect(recipeSourceInput).toBeEnabled();
  await recipeSourceInput.type(exampleRecipe.source);

  const recipeYieldsPortionsInput = await page.locator('input[data-qa="recipe-minimum-estimated-portions-input"]');
  await expect(recipeYieldsPortionsInput).toBeEnabled();
  await recipeYieldsPortionsInput.type(exampleRecipe.minimumEstimatedPortions.toString());
});
