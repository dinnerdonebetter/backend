import {
  mockLogin,
  mockAdminLogin,
  mockRegister,
  mockLogout,
  mockCheckPermissions,
  mockRequestPasswordResetToken,
  mockRedeemPasswordResetToken,
} from './auth';
export { MockUserPermissionsResponseConfig } from './auth';

import {
  mockCurrentHouseholdInfo,
  mockGetHousehold,
  mockUpdateHousehold,
  mockGetHouseholds,
  mockRemoveMember,
} from './households';
export {
  MockHouseholdResponseConfig,
  MockHouseholdListResponseConfig,
  MockRemoveUserFromHouseholdResponseConfig,
} from './households';

import {
  mockGetHouseholdInvitation,
  mockAcceptHouseholdInvitation,
  mockCancelHouseholdInvitation,
  mockRejectHouseholdInvitation,
  mockPendingInvitations,
} from './household_invitations';
export {
  MockHouseholdInvitationResponseConfig,
  MockHouseholdInvitationStatusChangeResponseConfig,
  MockPendingInvitationsResponseConfig,
} from './household_invitations';

import { mockRecipe, mockRecipesList, mockRecipesSearch, mockUpdateRecipe, mockDeleteRecipe } from './recipes';
export {
  MockRecipeResponseConfig,
  MockRecipeListResponseConfig,
  MockRecipeUpdateResponseConfig,
  MockRecipeDeleteResponseConfig,
} from './recipes';

import { mockSelf, mockUser, mockUserReputationUpdate, mockUsersList, mockUsersSearch } from './users';
export {
  MockUserResponseConfig,
  MockUsersListResponseConfig,
  MockUsersSearchResponseConfig,
  MockUserReputationUpdateResponseConfig,
} from './users';

import { mockMealPlan, mockMealPlans } from './meal_plans';
export { MockMealPlanResponseConfig, MockMealPlanListResponseConfig } from './meal_plans';

import { mockMeal, mockMeals } from './meals';
export { MockMealResponseConfig, MockMealListResponseConfig } from './meals';

import {
  mockDeleteValidIngredient,
  mockUpdateValidIngredient,
  mockValidIngredient,
  mockValidIngredientsList,
  mockValidIngredientsSearch,
} from './valid_ingredients';
export {
  MockValidIngredientResponseConfig,
  MockValidIngredientListResponseConfig,
  MockValidIngredientSearchResponseConfig,
  MockValidIngredientUpdateResponseConfig,
  MockValidIngredientDeleteResponseConfig,
} from './valid_ingredients';

import {
  mockDeleteValidInstrument,
  mockUpdateValidInstrument,
  mockValidInstrument,
  mockValidInstrumentsList,
  mockValidInstrumentsSearch,
} from './valid_instruments';
export {
  MockValidInstrumentResponseConfig,
  MockValidInstrumentListResponseConfig,
  MockValidInstrumentSearchResponseConfig,
  MockValidInstrumentUpdateResponseConfig,
  MockValidInstrumentDeleteResponseConfig,
} from './valid_instruments';

import {
  mockDeleteValidPreparation,
  mockUpdateValidPreparation,
  mockValidPreparation,
  mockValidPreparationsList,
  mockValidPreparationsSearch,
} from './valid_preparations';
export {
  MockValidPreparationResponseConfig,
  MockValidPreparationListResponseConfig,
  MockValidPreparationSearchResponseConfig,
  MockValidPreparationUpdateResponseConfig,
  MockValidPreparationDeleteResponseConfig,
} from './valid_preparations';

export * from './helpers';

export const mockRoutes = {
  auth: {
    login: mockLogin,
    adminLogin: mockAdminLogin,
    register: mockRegister,
    logout: mockLogout,
    checkPermissions: mockCheckPermissions,
    requestPasswordResetToken: mockRequestPasswordResetToken,
    redeemPasswordResetToken: mockRedeemPasswordResetToken,
  },
  households: {
    get: mockGetHousehold,
    current: mockCurrentHouseholdInfo,
    removeMember: mockRemoveMember,
    list: mockGetHouseholds,
    update: mockUpdateHousehold,
  },
  householdInvitations: {
    send: () => Promise.reject('TODO: not implemented'),
    received: () => Promise.reject('TODO: not implemented'),
    sent: mockPendingInvitations,
    get: mockGetHouseholdInvitation,
    accept: mockAcceptHouseholdInvitation,
    cancel: mockCancelHouseholdInvitation,
    reject: mockRejectHouseholdInvitation,
  },
  recipes: {
    create: () => Promise.reject('TODO: not implemented'),
    get: mockRecipe,
    list: mockRecipesList,
    update: mockUpdateRecipe,
    delete: mockDeleteRecipe,
    search: mockRecipesSearch,
  },
  meals: {
    get: mockMeal,
    list: mockMeals,
    create: () => Promise.reject('TODO: not implemented'),
    update: () => Promise.reject('TODO: not implemented'),
    delete: () => Promise.reject('TODO: not implemented'),
    search: () => Promise.reject('TODO: not implemented'),
  },
  mealPlans: {
    create: () => Promise.reject('TODO: not implemented'),
    get: mockMealPlan,
    list: mockMealPlans,
    update: () => Promise.reject('TODO: not implemented'),
    delete: () => Promise.reject('TODO: not implemented'),
  },
  users: {
    self: mockSelf,
    get: mockUser,
    list: mockUsersList,
    updateReputation: mockUserReputationUpdate,
    search: mockUsersSearch,
  },
  validPreparations: {
    create: () => Promise.reject('TODO: not implemented'),
    get: mockValidPreparation,
    list: mockValidPreparationsList,
    update: mockUpdateValidPreparation,
    delete: mockDeleteValidPreparation,
    search: mockValidPreparationsSearch,
  },
  validIngredients: {
    create: () => Promise.reject('TODO: not implemented'),
    get: mockValidIngredient,
    list: mockValidIngredientsList,
    update: mockUpdateValidIngredient,
    delete: mockDeleteValidIngredient,
    search: mockValidIngredientsSearch,
  },
  validInstruments: {
    create: () => Promise.reject('TODO: not implemented'),
    get: mockValidInstrument,
    list: mockValidInstrumentsList,
    update: mockUpdateValidInstrument,
    delete: mockDeleteValidInstrument,
    search: mockValidInstrumentsSearch,
  },
};
