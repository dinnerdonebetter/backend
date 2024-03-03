/* eslint-disable no-unused-vars */
export const enum backendRoutes {
  LOGIN = '/users/login',
  LOGIN_ADMIN = '/users/login/admin',
  LOGOUT = '/users/logout',

  USER_REGISTRATION = '/users/',
  VERIFY_2FA_SECRET = '/users/totp_secret/verify',
  NEW_2FA_SECRET = '/users/totp_secret/new',
  REQUEST_PASSWORD_RESET_TOKEN = '/users/password/reset',
  CHANGE_PASSWORD = '/users/password/new',
  REQUEST_USERNAME_REMINDER_EMAIL = '/users/username/reminder',
  REDEEM_PASSWORD_RESET_TOKEN = '/users/password/reset/redeem',

  SELF = '/api/v1/users/self',
  USER = '/api/v1/users/{}',
  USERS = '/api/v1/users',
  USER_REPUTATION_UPDATE = '/api/v1/admin/users/status',
  USERS_SEARCH = '/api/v1/users/search',
  USERS_VERIFY_EMAIL_ADDRESS = '/users/email_address/verify',
  USERS_UPLOAD_NEW_AVATAR = '/api/v1/users/avatar/upload',
  USERS_REQUEST_EMAIL_VERIFICATION_EMAIL = '/api/v1/users/email_address_verification',
  PERMISSIONS_CHECK = '/api/v1/users/permissions/check',

  HOUSEHOLD = '/api/v1/households/{}',
  HOUSEHOLDS = '/api/v1/households',
  HOUSEHOLD_INFO = '/api/v1/households/current',

  HOUSEHOLD_ADD_MEMBER = '/api/v1/households/{}/invitations',
  HOUSEHOLD_REMOVE_MEMBER = '/api/v1/households/{}/members/{}',

  HOUSEHOLD_INVITATION = '/api/v1/household_invitations/{}',
  ACCEPT_HOUSEHOLD_INVITATION = '/api/v1/household_invitations/{}/accept',
  CANCEL_HOUSEHOLD_INVITATION = '/api/v1/household_invitations/{}/cancel',
  REJECT_HOUSEHOLD_INVITATION = '/api/v1/household_invitations/{}/reject',
  RECEIVED_HOUSEHOLD_INVITATIONS = '/api/v1/household_invitations/received',
  SENT_HOUSEHOLD_INVITATIONS = '/api/v1/household_invitations/sent',

  MEAL = '/api/v1/meals/{}',
  MEALS = '/api/v1/meals',
  MEALS_SEARCH = '/api/v1/meals/search',

  MEAL_PLAN = '/api/v1/meal_plans/{}',
  MEAL_PLANS = '/api/v1/meal_plans',
  MEAL_PLAN_VOTING = '/api/v1/meal_plans/{}/events/{}/vote',

  MEAL_PLAN_GROCERY_LIST_ITEM = '/api/v1/meal_plans/{}/grocery_list_items/{}',
  MEAL_PLAN_GROCERY_LIST_ITEMS = '/api/v1/meal_plans/{}/grocery_list_items',

  MEAL_PLAN_TASK = '/api/v1/meal_plans/{}/tasks/{}',
  MEAL_PLAN_TASKS = '/api/v1/meal_plans/{}/tasks',

  RECIPE = '/api/v1/recipes/{}',
  RECIPES = '/api/v1/recipes',
  RECIPES_SEARCH = '/api/v1/recipes/search',

  VALID_INGREDIENTS = '/api/v1/valid_ingredients',
  VALID_INGREDIENTS_SEARCH = '/api/v1/valid_ingredients/search',
  VALID_INGREDIENTS_SEARCH_BY_PREPARATION_ID = '/api/v1/valid_ingredients/by_preparation/{}',
  VALID_INGREDIENT = '/api/v1/valid_ingredients/{}',

  OAUTH2_CLIENTS = '/api/v1/oauth2_clients',
  OAUTH2_CLIENT = '/api/v1/oauth2_clients/{}',

  VALID_INGREDIENT_GROUPS = '/api/v1/valid_ingredient_groups',
  VALID_INGREDIENT_GROUPS_SEARCH = '/api/v1/valid_ingredient_groups/search',
  VALID_INGREDIENT_GROUPS_SEARCH_BY_PREPARATION_ID = '/api/v1/valid_ingredient_groups/by_preparation/{}',
  VALID_INGREDIENT_GROUP = '/api/v1/valid_ingredient_groups/{}',

  VALID_MEASUREMENT_UNITS = '/api/v1/valid_measurement_units',
  VALID_MEASUREMENT_UNITS_SEARCH = '/api/v1/valid_measurement_units/search',
  VALID_MEASUREMENT_UNITS_SEARCH_BY_INGREDIENT = '/api/v1/valid_measurement_units/by_ingredient/{}',
  VALID_MEASUREMENT_UNIT = '/api/v1/valid_measurement_units/{}',

  VALID_MEASUREMENT_UNIT_CONVERSIONS = '/api/v1/valid_measurement_conversions',
  VALID_MEASUREMENT_UNIT_CONVERSIONS_FROM_UNIT = '/api/v1/valid_measurement_conversions/from_unit/{}',
  VALID_MEASUREMENT_UNIT_CONVERSIONS_TO_UNIT = '/api/v1/valid_measurement_conversions/to_unit/{}',
  VALID_MEASUREMENT_UNIT_CONVERSION = '/api/v1/valid_measurement_conversions/{}',

  VALID_INSTRUMENTS = '/api/v1/valid_instruments',
  VALID_INSTRUMENTS_SEARCH = '/api/v1/valid_instruments/search',
  VALID_INSTRUMENT = '/api/v1/valid_instruments/{}',

  VALID_VESSELS = '/api/v1/valid_vessels',
  VALID_VESSELS_SEARCH = '/api/v1/valid_vessels/search',
  VALID_VESSEL = '/api/v1/valid_vessels/{}',

  VALID_INGREDIENT_STATES = '/api/v1/valid_ingredient_states',
  VALID_INGREDIENT_STATES_SEARCH = '/api/v1/valid_ingredient_states/search',
  VALID_INGREDIENT_STATE = '/api/v1/valid_ingredient_states/{}',

  VALID_PREPARATIONS = '/api/v1/valid_preparations',
  VALID_PREPARATIONS_SEARCH = '/api/v1/valid_preparations/search',
  VALID_PREPARATION = '/api/v1/valid_preparations/{}',

  VALID_PREPARATION_INSTRUMENT = '/api/v1/valid_preparation_instruments/{}',
  VALID_PREPARATION_INSTRUMENTS = '/api/v1/valid_preparation_instruments',
  VALID_PREPARATION_INSTRUMENTS_SEARCH_BY_PREPARATION_ID = '/api/v1/valid_preparation_instruments/by_preparation/{}',
  VALID_PREPARATION_INSTRUMENTS_SEARCH_BY_INSTRUMENT_ID = '/api/v1/valid_preparation_instruments/by_instrument/{}',

  VALID_PREPARATION_VESSEL = '/api/v1/valid_preparation_vessels/{}',
  VALID_PREPARATION_VESSELS = '/api/v1/valid_preparation_vessels',
  VALID_PREPARATION_VESSELS_SEARCH_BY_PREPARATION_ID = '/api/v1/valid_preparation_vessels/by_preparation/{}',
  VALID_PREPARATION_VESSELS_SEARCH_BY_VESSEL_ID = '/api/v1/valid_preparation_vessels/by_vessel/{}',

  VALID_INGREDIENT_PREPARATION = '/api/v1/valid_ingredient_preparations/{}',
  VALID_INGREDIENT_PREPARATIONS = '/api/v1/valid_ingredient_preparations',
  VALID_INGREDIENT_PREPARATIONS_SEARCH_BY_INGREDIENT_ID = '/api/v1/valid_ingredient_preparations/by_ingredient/{}',
  VALID_INGREDIENT_PREPARATIONS_SEARCH_BY_PREPARATION_ID = '/api/v1/valid_ingredient_preparations/by_preparation/{}',

  VALID_INGREDIENT_STATE_INGREDIENT = '/api/v1/valid_ingredient_state_ingredients/{}',
  VALID_INGREDIENT_STATE_INGREDIENTS = '/api/v1/valid_ingredient_state_ingredients',
  VALID_INGREDIENT_STATE_INGREDIENTS_SEARCH_BY_INGREDIENT_ID = '/api/v1/valid_ingredient_state_ingredients/by_ingredient/{}',
  VALID_INGREDIENT_STATE_INGREDIENTS_SEARCH_BY_INGREDIENT_STATE = '/api/v1/valid_ingredient_state_ingredients/by_ingredient_state/{}',

  VALID_INGREDIENT_MEASUREMENT_UNIT = '/api/v1/valid_ingredient_measurement_units/{}',
  VALID_INGREDIENT_MEASUREMENT_UNITS = '/api/v1/valid_ingredient_measurement_units',
  VALID_INGREDIENT_MEASUREMENT_UNITS_SEARCH_BY_INGREDIENT_ID = '/api/v1/valid_ingredient_measurement_units/by_ingredient/{}',
  VALID_INGREDIENT_MEASUREMENT_UNITS_SEARCH_BY_MEASUREMENT_UNIT_ID = '/api/v1/valid_ingredient_measurement_units/by_measurement_unit/{}',

  SERVICE_SETTINGS = '/api/v1/settings',
  SERVICE_SETTINGS_SEARCH = '/api/v1/settings/search',
  SERVICE_SETTING = '/api/v1/settings/{}',

  SERVICE_SETTING_CONFIGURATIONS_FOR_USER = '/api/v1/settings/configurations/user',
  SERVICE_SETTING_CONFIGURATIONS_FOR_HOUSEHOLD = '/api/v1/settings/configurations/household',
  SERVICE_SETTING_CONFIGURATION = '/api/v1/settings/configurations/{}',
}
