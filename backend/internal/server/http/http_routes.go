package http

import (
	"context"
	"fmt"
	"net/http"
	"path"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/routing"
	auditlogentriesservice "github.com/dinnerdonebetter/backend/internal/services/auditlogentries"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	householdinstrumentownershipsservice "github.com/dinnerdonebetter/backend/internal/services/householdinstrumentownerships"
	householdinvitationsservice "github.com/dinnerdonebetter/backend/internal/services/householdinvitations"
	householdsservice "github.com/dinnerdonebetter/backend/internal/services/households"
	mealplaneventsservice "github.com/dinnerdonebetter/backend/internal/services/mealplanevents"
	mealplangrocerylistitemssservice "github.com/dinnerdonebetter/backend/internal/services/mealplangrocerylistitems"
	mealplanoptionsservice "github.com/dinnerdonebetter/backend/internal/services/mealplanoptions"
	mealplanoptionvotesservice "github.com/dinnerdonebetter/backend/internal/services/mealplanoptionvotes"
	mealplansservice "github.com/dinnerdonebetter/backend/internal/services/mealplans"
	mealplantasksservice "github.com/dinnerdonebetter/backend/internal/services/mealplantasks"
	mealsservice "github.com/dinnerdonebetter/backend/internal/services/meals"
	oauth2clientsservice "github.com/dinnerdonebetter/backend/internal/services/oauth2clients"
	recipepreptasksservice "github.com/dinnerdonebetter/backend/internal/services/recipepreptasks"
	reciperatingsservice "github.com/dinnerdonebetter/backend/internal/services/reciperatings"
	recipesservice "github.com/dinnerdonebetter/backend/internal/services/recipes"
	recipestepcompletionconditionsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepcompletionconditions"
	recipestepingredientsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepingredients"
	recipestepinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepinstruments"
	recipestepproductsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepproducts"
	recipestepsservice "github.com/dinnerdonebetter/backend/internal/services/recipesteps"
	recipestepvesselsservice "github.com/dinnerdonebetter/backend/internal/services/recipestepvessels"
	servicesettingconfigurationsservice "github.com/dinnerdonebetter/backend/internal/services/servicesettingconfigurations"
	servicesettingsservice "github.com/dinnerdonebetter/backend/internal/services/servicesettings"
	useringredientpreferencesservice "github.com/dinnerdonebetter/backend/internal/services/useringredientpreferences"
	usernotificationsservice "github.com/dinnerdonebetter/backend/internal/services/usernotifications"
	usersservice "github.com/dinnerdonebetter/backend/internal/services/users"
	validingredientgroupsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientgroups"
	validingredientmeasurementunitsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientmeasurementunits"
	validingredientpreparationsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientpreparations"
	validingredientsservice "github.com/dinnerdonebetter/backend/internal/services/validingredients"
	validingredientstateingredientsservice "github.com/dinnerdonebetter/backend/internal/services/validingredientstateingredients"
	validingredientstatesservice "github.com/dinnerdonebetter/backend/internal/services/validingredientstates"
	validinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/validinstruments"
	validmeasurementconversionsservice "github.com/dinnerdonebetter/backend/internal/services/validmeasurementunitconversions"
	validmeasurementunitsservice "github.com/dinnerdonebetter/backend/internal/services/validmeasurementunits"
	validpreparationinstrumentsservice "github.com/dinnerdonebetter/backend/internal/services/validpreparationinstruments"
	validpreparationsservice "github.com/dinnerdonebetter/backend/internal/services/validpreparations"
	validpreparationvesselsservice "github.com/dinnerdonebetter/backend/internal/services/validpreparationvessels"
	validvesselsservice "github.com/dinnerdonebetter/backend/internal/services/validvessels"
	webhooksservice "github.com/dinnerdonebetter/backend/internal/services/webhooks"
)

const (
	root       = "/"
	randomRoot = "/random"
	searchRoot = "/search"
)

func buildURLVarChunk(key, pattern string) string {
	if pattern != "" {
		return fmt.Sprintf("/{%s:%s}", key, pattern)
	}
	return fmt.Sprintf("/{%s}", key)
}

func (s *server) setupRouter(ctx context.Context, router routing.Router) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	router.Route("/_meta_", func(metaRouter routing.Router) {
		// Expose a readiness check on /ready
		metaRouter.Get("/ready", func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
		})
	})

	authenticatedRouter := router.WithMiddleware(s.authService.UserAttributionMiddleware)
	authenticatedRouter.Get("/auth/status", s.authService.StatusHandler)

	router.Route("/oauth2", func(userRouter routing.Router) {
		userRouter.
			WithMiddleware(s.authService.CookieRequirementMiddleware, s.authService.UserAttributionMiddleware).
			Get("/authorize", s.authService.AuthorizeHandler)
		userRouter.Post("/token", s.authService.TokenHandler)
	})

	router.Route("/users", func(userRouter routing.Router) {
		userRouter.Post(root, s.usersService.CreateHandler)
		userRouter.Post("/login", s.authService.BuildLoginHandler(false))
		userRouter.Post("/login/admin", s.authService.BuildLoginHandler(true))
		userRouter.WithMiddleware(s.authService.UserAttributionMiddleware, s.authService.CookieRequirementMiddleware).
			Post("/logout", s.authService.EndSessionHandler)
		userRouter.Post("/username/reminder", s.usersService.RequestUsernameReminderHandler)
		userRouter.Post("/password/reset", s.usersService.CreatePasswordResetTokenHandler)
		userRouter.Post("/password/reset/redeem", s.usersService.PasswordResetTokenRedemptionHandler)
		userRouter.Post("/email_address/verify", s.usersService.VerifyUserEmailAddressHandler)
		userRouter.Post("/totp_secret/verify", s.usersService.TOTPSecretVerificationHandler)
	})

	router.Route("/auth", func(authRouter routing.Router) {
		providerRouteParam := buildURLVarChunk(authservice.AuthProviderParamKey, "")
		authRouter.Get(providerRouteParam, s.authService.SSOLoginHandler)
		authRouter.Get(path.Join(providerRouteParam, "callback"), s.authService.SSOLoginCallbackHandler)
	})

	authenticatedRouter.WithMiddleware(s.authService.AuthorizationMiddleware).Route("/api/v1", func(v1Router routing.Router) {
		adminRouter := v1Router.WithMiddleware(s.authService.ServiceAdminMiddleware)

		// Admin
		adminRouter.Route("/admin", func(adminRouter routing.Router) {
			adminRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CycleCookieSecretPermission)).
				Post("/cycle_cookie_secret", s.authService.CycleCookieSecretHandler)
			adminRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateUserStatusPermission)).
				Post("/users/status", s.adminService.UserAccountStatusChangeHandler)
		})

		// Workers
		adminRouter.Route("/workers", func(adminRouter routing.Router) {
			adminRouter.
				Post("/finalize_meal_plans", s.workerService.MealPlanFinalizationHandler)
			adminRouter.
				Post("/meal_plan_grocery_list_init", s.workerService.MealPlanGroceryListInitializationHandler)
			adminRouter.
				Post("/meal_plan_tasks", s.workerService.MealPlanTaskCreationHandler)
		})

		// Users
		v1Router.Route("/users", func(usersRouter routing.Router) {
			usersRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadUserPermission)).
				Get(root, s.usersService.ListHandler)
			usersRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.SearchUserPermission)).
				Get(searchRoot, s.usersService.UsernameSearchHandler)
			usersRouter.Post("/avatar/upload", s.usersService.AvatarUploadHandler)

			usersRouter.Get("/self", s.usersService.SelfHandler)
			usersRouter.Post("/email_address_verification", s.usersService.RequestEmailVerificationEmailHandler)
			usersRouter.Post("/permissions/check", s.usersService.PermissionsHandler)
			usersRouter.Post("/household/select", s.authService.ChangeActiveHouseholdHandler)
			usersRouter.Put("/password/new", s.usersService.UpdatePasswordHandler)
			usersRouter.Post("/totp_secret/new", s.usersService.NewTOTPSecretHandler)
			usersRouter.Put("/username", s.usersService.UpdateUserUsernameHandler)
			usersRouter.Put("/email_address", s.usersService.UpdateUserEmailAddressHandler)
			usersRouter.Put("/details", s.usersService.UpdateUserDetailsHandler)

			singleUserRoute := buildURLVarChunk(usersservice.UserIDURIParamKey, "")
			usersRouter.Route(singleUserRoute, func(singleUserRouter routing.Router) {
				singleUserRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadUserPermission)).
					Get(root, s.usersService.ReadHandler)

				singleUserRouter.Delete(root, s.usersService.ArchiveHandler)
			})
		})

		// Households
		v1Router.Route("/households", func(householdsRouter routing.Router) {
			householdsRouter.Post(root, s.householdsService.CreateHandler)
			householdsRouter.Get(root, s.householdsService.ListHandler)
			householdsRouter.Get("/current", s.householdsService.CurrentInfoHandler)

			singleUserRoute := buildURLVarChunk(householdsservice.UserIDURIParamKey, "")
			singleHouseholdRoute := buildURLVarChunk(householdsservice.HouseholdIDURIParamKey, "")
			householdsRouter.Route(singleHouseholdRoute, func(singleHouseholdRouter routing.Router) {
				singleHouseholdRouter.Get(root, s.householdsService.ReadHandler)
				singleHouseholdRouter.Put(root, s.householdsService.UpdateHandler)
				singleHouseholdRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveHouseholdPermission)).
					Delete(root, s.householdsService.ArchiveHandler)

				singleHouseholdRouter.Post("/default", s.householdsService.MarkAsDefaultHouseholdHandler)
				singleHouseholdRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.RemoveMemberHouseholdPermission)).
					Delete("/members"+singleUserRoute, s.householdsService.RemoveMemberHandler)
				singleHouseholdRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.InviteUserToHouseholdPermission)).
					Post("/invite", s.householdInvitationsService.InviteMemberHandler)
				singleHouseholdRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ModifyMemberPermissionsForHouseholdPermission)).
					Patch("/members"+singleUserRoute+"/permissions", s.householdsService.ModifyMemberPermissionsHandler)
				singleHouseholdRouter.Post("/transfer", s.householdsService.TransferHouseholdOwnershipHandler)

				singleHouseholdRouter.Route("/invitations", func(invitationsRouter routing.Router) {
					invitationsRouter.Post(root, s.householdInvitationsService.InviteMemberHandler)

					singleHouseholdInvitationRoute := buildURLVarChunk(householdinvitationsservice.HouseholdInvitationIDURIParamKey, "")
					invitationsRouter.Route(singleHouseholdInvitationRoute, func(singleHouseholdInvitationRouter routing.Router) {
						singleHouseholdInvitationRouter.Get(root, s.householdInvitationsService.ReadHandler)
					})
				})
			})

			// HouseholdInstrumentOwnerships
			householdInstrumentOwnershipsRouteWithPrefix := "/instruments"
			householdInstrumentOwnershipIDRouteParam := buildURLVarChunk(householdinstrumentownershipsservice.HouseholdInstrumentOwnershipIDURIParamKey, "")
			householdsRouter.Route(householdInstrumentOwnershipsRouteWithPrefix, func(householdInstrumentOwnershipsRouter routing.Router) {
				householdInstrumentOwnershipsRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateHouseholdInstrumentOwnershipsPermission)).
					Post(root, s.householdInstrumentOwnershipService.CreateHandler)
				householdInstrumentOwnershipsRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadHouseholdInstrumentOwnershipsPermission)).
					Get(root, s.householdInstrumentOwnershipService.ListHandler)

				householdInstrumentOwnershipsRouter.Route(householdInstrumentOwnershipIDRouteParam, func(singleHouseholdInstrumentOwnershipRouter routing.Router) {
					singleHouseholdInstrumentOwnershipRouter.
						WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadHouseholdInstrumentOwnershipsPermission)).
						Get(root, s.householdInstrumentOwnershipService.ReadHandler)
					singleHouseholdInstrumentOwnershipRouter.
						WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateHouseholdInstrumentOwnershipsPermission)).
						Put(root, s.householdInstrumentOwnershipService.UpdateHandler)
					singleHouseholdInstrumentOwnershipRouter.
						WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveHouseholdInstrumentOwnershipsPermission)).
						Delete(root, s.householdInstrumentOwnershipService.ArchiveHandler)
				})
			})
		})

		v1Router.Route("/household_invitations", func(householdInvitationsRouter routing.Router) {
			householdInvitationsRouter.Get("/sent", s.householdInvitationsService.OutboundInvitesHandler)
			householdInvitationsRouter.Get("/received", s.householdInvitationsService.InboundInvitesHandler)

			singleHouseholdInvitationRoute := buildURLVarChunk(householdinvitationsservice.HouseholdInvitationIDURIParamKey, "")
			householdInvitationsRouter.Route(singleHouseholdInvitationRoute, func(singleHouseholdInvitationRouter routing.Router) {
				singleHouseholdInvitationRouter.Get(root, s.householdInvitationsService.ReadHandler)
				singleHouseholdInvitationRouter.Put("/cancel", s.householdInvitationsService.CancelInviteHandler)
				singleHouseholdInvitationRouter.Put("/accept", s.householdInvitationsService.AcceptInviteHandler)
				singleHouseholdInvitationRouter.Put("/reject", s.householdInvitationsService.RejectInviteHandler)
			})
		})

		// OAuth2 Clients
		v1Router.Route("/oauth2_clients", func(clientRouter routing.Router) {
			clientRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadOAuth2ClientsPermission)).
				Get(root, s.oauth2ClientsService.ListHandler)
			clientRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateOAuth2ClientsPermission)).
				Post(root, s.oauth2ClientsService.CreateHandler)

			singleClientRoute := buildURLVarChunk(oauth2clientsservice.OAuth2ClientIDURIParamKey, "")
			clientRouter.Route(singleClientRoute, func(singleClientRouter routing.Router) {
				singleClientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadOAuth2ClientsPermission)).
					Get(root, s.oauth2ClientsService.ReadHandler)
				singleClientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveOAuth2ClientsPermission)).
					Delete(root, s.oauth2ClientsService.ArchiveHandler)
			})
		})

		// Webhooks
		v1Router.Route("/webhooks", func(webhookRouter routing.Router) {
			webhookRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadWebhooksPermission)).
				Get(root, s.webhooksService.ListWebhooksHandler)
			webhookRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateWebhooksPermission)).
				Post(root, s.webhooksService.CreateWebhookHandler)

			singleWebhookRoute := buildURLVarChunk(webhooksservice.WebhookIDURIParamKey, "")
			webhookRouter.Route(singleWebhookRoute, func(singleWebhookRouter routing.Router) {
				singleWebhookRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadWebhooksPermission)).
					Get(root, s.webhooksService.ReadWebhookHandler)
				singleWebhookRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveWebhooksPermission)).
					Delete(root, s.webhooksService.ArchiveWebhookHandler)

				singleWebhookTriggerEventRoute := buildURLVarChunk(webhooksservice.WebhookTriggerEventIDURIParamKey, "")
				singleWebhookRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateWebhookTriggerEventsPermission)).
					Post("/trigger_events", s.webhooksService.AddWebhookTriggerEventHandler)

				singleWebhookRouter.Route("/trigger_events"+singleWebhookTriggerEventRoute, func(singleWebhookRouter routing.Router) {
					singleWebhookRouter.
						WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveWebhookTriggerEventsPermission)).
						Delete(root, s.webhooksService.ArchiveWebhookTriggerEventHandler)
				})
			})
		})

		// Audit Log Entries
		v1Router.Route("/audit_log_entries", func(auditLogEntriesRouter routing.Router) {
			singleAuditLogEntryRoute := buildURLVarChunk(auditlogentriesservice.AuditLogEntryIDURIParamKey, "")
			auditLogEntriesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadAuditLogEntriesPermission)).
				Get(singleAuditLogEntryRoute, s.auditLogEntriesService.ReadAuditLogEntryHandler)
			auditLogEntriesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadAuditLogEntriesPermission)).
				Get("/for_user", s.auditLogEntriesService.ListUserAuditLogEntriesHandler)
			auditLogEntriesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadAuditLogEntriesPermission)).
				Get("/for_household", s.auditLogEntriesService.ListHouseholdAuditLogEntriesHandler)
		})

		// ValidInstruments
		validInstrumentPath := "valid_instruments"
		validInstrumentsRouteWithPrefix := fmt.Sprintf("/%s", validInstrumentPath)
		validInstrumentIDRouteParam := buildURLVarChunk(validinstrumentsservice.ValidInstrumentIDURIParamKey, "")
		v1Router.Route(validInstrumentsRouteWithPrefix, func(validInstrumentsRouter routing.Router) {
			validInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidInstrumentsPermission)).
				Post(root, s.validInstrumentsService.CreateHandler)
			validInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidInstrumentsPermission)).
				Get(root, s.validInstrumentsService.ListHandler)
			validInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.SearchValidInstrumentsPermission)).
				Get(searchRoot, s.validInstrumentsService.SearchHandler)
			validInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidInstrumentsPermission)).
				Get(randomRoot, s.validInstrumentsService.RandomHandler)

			validInstrumentsRouter.Route(validInstrumentIDRouteParam, func(singleValidInstrumentRouter routing.Router) {
				singleValidInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidInstrumentsPermission)).
					Get(root, s.validInstrumentsService.ReadHandler)
				singleValidInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidInstrumentsPermission)).
					Put(root, s.validInstrumentsService.UpdateHandler)
				singleValidInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidInstrumentsPermission)).
					Delete(root, s.validInstrumentsService.ArchiveHandler)
			})
		})

		// ValidVessels
		validVesselPath := "valid_vessels"
		validVesselsRouteWithPrefix := fmt.Sprintf("/%s", validVesselPath)
		validVesselIDRouteParam := buildURLVarChunk(validvesselsservice.ValidVesselIDURIParamKey, "")
		v1Router.Route(validVesselsRouteWithPrefix, func(validVesselsRouter routing.Router) {
			validVesselsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidVesselsPermission)).
				Post(root, s.validVesselsService.CreateHandler)
			validVesselsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidVesselsPermission)).
				Get(root, s.validVesselsService.ListHandler)
			validVesselsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.SearchValidVesselsPermission)).
				Get(searchRoot, s.validVesselsService.SearchHandler)
			validVesselsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidVesselsPermission)).
				Get(randomRoot, s.validVesselsService.RandomHandler)

			validVesselsRouter.Route(validVesselIDRouteParam, func(singleValidVesselRouter routing.Router) {
				singleValidVesselRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidVesselsPermission)).
					Get(root, s.validVesselsService.ReadHandler)
				singleValidVesselRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidVesselsPermission)).
					Put(root, s.validVesselsService.UpdateHandler)
				singleValidVesselRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidVesselsPermission)).
					Delete(root, s.validVesselsService.ArchiveHandler)
			})
		})

		// ValidIngredients
		validIngredientPath := "valid_ingredients"
		validIngredientsRouteWithPrefix := fmt.Sprintf("/%s", validIngredientPath)
		validIngredientIDRouteParam := buildURLVarChunk(validingredientsservice.ValidIngredientIDURIParamKey, "")
		v1Router.Route(validIngredientsRouteWithPrefix, func(validIngredientsRouter routing.Router) {
			validIngredientsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidIngredientsPermission)).
				Post(root, s.validIngredientsService.CreateHandler)
			validIngredientsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientsPermission)).
				Get(root, s.validIngredientsService.ListHandler)
			validIngredientsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientsPermission, authorization.SearchValidIngredientsPermission)).
				Get(searchRoot, s.validIngredientsService.SearchHandler)
			validIngredientsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientsPermission)).
				Get(randomRoot, s.validIngredientsService.RandomHandler)

			validIngredientsByPreparationIDSearchRouteParam := fmt.Sprintf("/by_preparation%s", buildURLVarChunk(validingredientsservice.ValidPreparationIDURIParamKey, ""))
			validIngredientsRouter.Route(validIngredientsByPreparationIDSearchRouteParam, func(byValidPreparationIDRouter routing.Router) {
				byValidPreparationIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientPreparationsPermission)).
					Get(root, s.validIngredientsService.SearchByPreparationAndIngredientNameHandler)
			})

			validIngredientsRouter.Route(validIngredientIDRouteParam, func(singleValidIngredientRouter routing.Router) {
				singleValidIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientsPermission)).
					Get(root, s.validIngredientsService.ReadHandler)
				singleValidIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientsPermission)).
					Put(root, s.validIngredientsService.UpdateHandler)
				singleValidIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientsPermission)).
					Delete(root, s.validIngredientsService.ArchiveHandler)
			})
		})

		// ValidIngredientGroups
		validIngredientGroupPath := "valid_ingredient_groups"
		validIngredientGroupsRouteWithPrefix := fmt.Sprintf("/%s", validIngredientGroupPath)
		validIngredientGroupIDRouteParam := buildURLVarChunk(validingredientgroupsservice.ValidIngredientGroupIDURIParamKey, "")
		v1Router.Route(validIngredientGroupsRouteWithPrefix, func(validIngredientGroupsRouter routing.Router) {
			validIngredientGroupsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidIngredientGroupsPermission)).
				Post(root, s.validIngredientGroupsService.CreateHandler)
			validIngredientGroupsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientGroupsPermission)).
				Get(root, s.validIngredientGroupsService.ListHandler)
			validIngredientGroupsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.SearchValidIngredientGroupsPermission)).
				Get(searchRoot, s.validIngredientGroupsService.SearchHandler)

			validIngredientGroupsRouter.Route(validIngredientGroupIDRouteParam, func(singleValidIngredientGroupRouter routing.Router) {
				singleValidIngredientGroupRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientGroupsPermission)).
					Get(root, s.validIngredientGroupsService.ReadHandler)
				singleValidIngredientGroupRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientGroupsPermission)).
					Put(root, s.validIngredientGroupsService.UpdateHandler)
				singleValidIngredientGroupRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientGroupsPermission)).
					Delete(root, s.validIngredientGroupsService.ArchiveHandler)
			})
		})

		// ValidPreparations
		validPreparationPath := "valid_preparations"
		validPreparationsRouteWithPrefix := fmt.Sprintf("/%s", validPreparationPath)
		validPreparationIDRouteParam := buildURLVarChunk(validpreparationsservice.ValidPreparationIDURIParamKey, "")
		v1Router.Route(validPreparationsRouteWithPrefix, func(validPreparationsRouter routing.Router) {
			validPreparationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidPreparationsPermission)).
				Post(root, s.validPreparationsService.CreateHandler)
			validPreparationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationsPermission)).
				Get(root, s.validPreparationsService.ListHandler)
			validPreparationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationsPermission)).
				Get(searchRoot, s.validPreparationsService.SearchHandler)
			validPreparationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationsPermission)).
				Get(randomRoot, s.validPreparationsService.RandomHandler)

			validPreparationsRouter.Route(validPreparationIDRouteParam, func(singleValidPreparationRouter routing.Router) {
				singleValidPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationsPermission)).
					Get(root, s.validPreparationsService.ReadHandler)
				singleValidPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidPreparationsPermission)).
					Put(root, s.validPreparationsService.UpdateHandler)
				singleValidPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidPreparationsPermission)).
					Delete(root, s.validPreparationsService.ArchiveHandler)
			})
		})

		// UserIngredientPreferences
		userIngredientPreferencesPath := "user_ingredient_preferences"
		userIngredientPreferencesRouteWithPrefix := fmt.Sprintf("/%s", userIngredientPreferencesPath)
		userIngredientPreferencesIDRouteParam := buildURLVarChunk(useringredientpreferencesservice.UserIngredientPreferenceIDURIParamKey, "")
		v1Router.Route(userIngredientPreferencesRouteWithPrefix, func(userIngredientPreferencesRouter routing.Router) {
			userIngredientPreferencesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateUserIngredientPreferencesPermission)).
				Post(root, s.userIngredientPreferencesService.CreateHandler)
			userIngredientPreferencesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadUserIngredientPreferencesPermission)).
				Get(root, s.userIngredientPreferencesService.ListHandler)

			userIngredientPreferencesRouter.Route(userIngredientPreferencesIDRouteParam, func(singleUserIngredientPreferenceRouter routing.Router) {
				singleUserIngredientPreferenceRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateUserIngredientPreferencesPermission)).
					Put(root, s.userIngredientPreferencesService.UpdateHandler)
				singleUserIngredientPreferenceRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveUserIngredientPreferencesPermission)).
					Delete(root, s.userIngredientPreferencesService.ArchiveHandler)
			})
		})

		// ValidMeasurementUnits
		validMeasurementUnitPath := "valid_measurement_units"
		validMeasurementUnitsRouteWithPrefix := fmt.Sprintf("/%s", validMeasurementUnitPath)
		validMeasurementUnitIDRouteParam := buildURLVarChunk(validmeasurementunitsservice.ValidMeasurementUnitIDURIParamKey, "")
		validMeasurementUnitServiceIngredientIDRouteParam := buildURLVarChunk(validmeasurementunitsservice.ValidIngredientIDURIParamKey, "")
		v1Router.Route(validMeasurementUnitsRouteWithPrefix, func(validMeasurementUnitsRouter routing.Router) {
			validMeasurementUnitsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidMeasurementUnitsPermission)).
				Post(root, s.validMeasurementUnitsService.CreateHandler)
			validMeasurementUnitsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidMeasurementUnitsPermission)).
				Get(root, s.validMeasurementUnitsService.ListHandler)
			validMeasurementUnitsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.SearchValidMeasurementUnitsPermission)).
				Get(searchRoot, s.validMeasurementUnitsService.SearchHandler)

			validMeasurementUnitsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.SearchValidMeasurementUnitsPermission)).
				Get(path.Join("/by_ingredient", validMeasurementUnitServiceIngredientIDRouteParam), s.validMeasurementUnitsService.SearchByIngredientIDHandler)

			validMeasurementUnitsRouter.Route(validMeasurementUnitIDRouteParam, func(singleValidMeasurementUnitRouter routing.Router) {
				singleValidMeasurementUnitRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidMeasurementUnitsPermission)).
					Get(root, s.validMeasurementUnitsService.ReadHandler)
				singleValidMeasurementUnitRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidMeasurementUnitsPermission)).
					Put(root, s.validMeasurementUnitsService.UpdateHandler)
				singleValidMeasurementUnitRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidMeasurementUnitsPermission)).
					Delete(root, s.validMeasurementUnitsService.ArchiveHandler)
			})
		})

		// ValidIngredientStates
		validIngredientStatePath := "valid_ingredient_states"
		validIngredientStatesRouteWithPrefix := fmt.Sprintf("/%s", validIngredientStatePath)
		validIngredientStateIDRouteParam := buildURLVarChunk(validingredientstatesservice.ValidIngredientStateIDURIParamKey, "")
		v1Router.Route(validIngredientStatesRouteWithPrefix, func(validIngredientStatesRouter routing.Router) {
			validIngredientStatesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidIngredientStatesPermission)).
				Post(root, s.validIngredientStatesService.CreateHandler)
			validIngredientStatesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientStatesPermission)).
				Get(root, s.validIngredientStatesService.ListHandler)
			validIngredientStatesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientStatesPermission)).
				Get(searchRoot, s.validIngredientStatesService.SearchHandler)

			validIngredientStatesRouter.Route(validIngredientStateIDRouteParam, func(singleValidIngredientStateRouter routing.Router) {
				singleValidIngredientStateRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientStatesPermission)).
					Get(root, s.validIngredientStatesService.ReadHandler)
				singleValidIngredientStateRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientStatesPermission)).
					Put(root, s.validIngredientStatesService.UpdateHandler)
				singleValidIngredientStateRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientStatesPermission)).
					Delete(root, s.validIngredientStatesService.ArchiveHandler)
			})
		})

		// ValidMeasurementUnitConversions
		validMeasurementUnitConversionPath := "valid_measurement_conversions"
		validMeasurementUnitConversionsRouteWithPrefix := fmt.Sprintf("/%s", validMeasurementUnitConversionPath)
		validMeasurementUnitConversionUnitIDRouteParam := buildURLVarChunk(validmeasurementconversionsservice.ValidMeasurementUnitIDURIParamKey, "")
		validMeasurementUnitConversionIDRouteParam := buildURLVarChunk(validmeasurementconversionsservice.ValidMeasurementUnitConversionIDURIParamKey, "")
		v1Router.Route(validMeasurementUnitConversionsRouteWithPrefix, func(validMeasurementUnitConversionsRouter routing.Router) {
			validMeasurementUnitConversionsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidMeasurementUnitConversionsPermission)).
				Post(root, s.validMeasurementUnitConversionsService.CreateHandler)

			validMeasurementUnitConversionsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidMeasurementUnitConversionsPermission)).
				Get(path.Join("/from_unit", validMeasurementUnitConversionUnitIDRouteParam), s.validMeasurementUnitConversionsService.FromMeasurementUnitHandler)
			validMeasurementUnitConversionsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidMeasurementUnitConversionsPermission)).
				Get(path.Join("/to_unit", validMeasurementUnitConversionUnitIDRouteParam), s.validMeasurementUnitConversionsService.ToMeasurementUnitHandler)

			validMeasurementUnitConversionsRouter.Route(validMeasurementUnitConversionIDRouteParam, func(singleValidMeasurementUnitConversionRouter routing.Router) {
				singleValidMeasurementUnitConversionRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidMeasurementUnitConversionsPermission)).
					Get(root, s.validMeasurementUnitConversionsService.ReadHandler)
				singleValidMeasurementUnitConversionRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidMeasurementUnitConversionsPermission)).
					Put(root, s.validMeasurementUnitConversionsService.UpdateHandler)
				singleValidMeasurementUnitConversionRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidMeasurementUnitConversionsPermission)).
					Delete(root, s.validMeasurementUnitConversionsService.ArchiveHandler)
			})
		})

		// ValidIngredientStateIngredients
		validIngredientStateIngredientPath := "valid_ingredient_state_ingredients"
		validIngredientStateIngredientsRouteWithPrefix := fmt.Sprintf("/%s", validIngredientStateIngredientPath)
		validIngredientStateIngredientIDRouteParam := buildURLVarChunk(validingredientstateingredientsservice.ValidIngredientStateIngredientIDURIParamKey, "")
		v1Router.Route(validIngredientStateIngredientsRouteWithPrefix, func(validIngredientStateIngredientsRouter routing.Router) {
			validIngredientStateIngredientsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidIngredientStateIngredientsPermission)).
				Post(root, s.validIngredientStateIngredientsService.CreateHandler)
			validIngredientStateIngredientsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientStateIngredientsPermission)).
				Get(root, s.validIngredientStateIngredientsService.ListHandler)

			validIngredientStateIngredientsByIngredientIDRouteParam := fmt.Sprintf("/by_ingredient%s", buildURLVarChunk(validingredientstateingredientsservice.ValidIngredientIDURIParamKey, ""))
			validIngredientStateIngredientsRouter.Route(validIngredientStateIngredientsByIngredientIDRouteParam, func(byValidIngredientIDRouter routing.Router) {
				byValidIngredientIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientStateIngredientsPermission)).
					Get(root, s.validIngredientStateIngredientsService.SearchByIngredientHandler)
			})

			validIngredientStateIngredientsByStateIngredientIDRouteParam := fmt.Sprintf("/by_ingredient_state%s", buildURLVarChunk(validingredientstateingredientsservice.ValidIngredientStateIDURIParamKey, ""))
			validIngredientStateIngredientsRouter.Route(validIngredientStateIngredientsByStateIngredientIDRouteParam, func(byValidStateIngredientIDRouter routing.Router) {
				byValidStateIngredientIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientStateIngredientsPermission)).
					Get(root, s.validIngredientStateIngredientsService.SearchByIngredientStateHandler)
			})

			validIngredientStateIngredientsRouter.Route(validIngredientStateIngredientIDRouteParam, func(singleValidIngredientStateIngredientRouter routing.Router) {
				singleValidIngredientStateIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientStateIngredientsPermission)).
					Get(root, s.validIngredientStateIngredientsService.ReadHandler)
				singleValidIngredientStateIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientStateIngredientsPermission)).
					Put(root, s.validIngredientStateIngredientsService.UpdateHandler)
				singleValidIngredientStateIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientStateIngredientsPermission)).
					Delete(root, s.validIngredientStateIngredientsService.ArchiveHandler)
			})
		})

		// ValidIngredientPreparations
		validIngredientPreparationPath := "valid_ingredient_preparations"
		validIngredientPreparationsRouteWithPrefix := fmt.Sprintf("/%s", validIngredientPreparationPath)
		validIngredientPreparationIDRouteParam := buildURLVarChunk(validingredientpreparationsservice.ValidIngredientPreparationIDURIParamKey, "")
		v1Router.Route(validIngredientPreparationsRouteWithPrefix, func(validIngredientPreparationsRouter routing.Router) {
			validIngredientPreparationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidIngredientPreparationsPermission)).
				Post(root, s.validIngredientPreparationsService.CreateHandler)
			validIngredientPreparationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientPreparationsPermission)).
				Get(root, s.validIngredientPreparationsService.ListHandler)

			validIngredientPreparationsByIngredientIDRouteParam := fmt.Sprintf("/by_ingredient%s", buildURLVarChunk(validingredientpreparationsservice.ValidIngredientIDURIParamKey, ""))
			validIngredientPreparationsRouter.Route(validIngredientPreparationsByIngredientIDRouteParam, func(byValidIngredientIDRouter routing.Router) {
				byValidIngredientIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientPreparationsPermission)).
					Get(root, s.validIngredientPreparationsService.SearchByIngredientHandler)
			})

			validIngredientPreparationsByPreparationIDRouteParam := fmt.Sprintf("/by_preparation%s", buildURLVarChunk(validingredientpreparationsservice.ValidPreparationIDURIParamKey, ""))
			validIngredientPreparationsRouter.Route(validIngredientPreparationsByPreparationIDRouteParam, func(byValidPreparationIDRouter routing.Router) {
				byValidPreparationIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientPreparationsPermission)).
					Get(root, s.validIngredientPreparationsService.SearchByPreparationHandler)
			})

			validIngredientPreparationsRouter.Route(validIngredientPreparationIDRouteParam, func(singleValidIngredientPreparationRouter routing.Router) {
				singleValidIngredientPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientPreparationsPermission)).
					Get(root, s.validIngredientPreparationsService.ReadHandler)
				singleValidIngredientPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientPreparationsPermission)).
					Put(root, s.validIngredientPreparationsService.UpdateHandler)
				singleValidIngredientPreparationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientPreparationsPermission)).
					Delete(root, s.validIngredientPreparationsService.ArchiveHandler)
			})
		})

		// ValidPreparationInstruments
		validPreparationInstrumentPath := "valid_preparation_instruments"
		validPreparationInstrumentsRouteWithPrefix := fmt.Sprintf("/%s", validPreparationInstrumentPath)
		v1Router.Route(validPreparationInstrumentsRouteWithPrefix, func(validPreparationInstrumentsRouter routing.Router) {
			validPreparationInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidPreparationInstrumentsPermission)).
				Post(root, s.validPreparationInstrumentsService.CreateHandler)
			validPreparationInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationInstrumentsPermission)).
				Get(root, s.validPreparationInstrumentsService.ListHandler)

			validPreparationInstrumentsByPreparationIDRouteParam := fmt.Sprintf("/by_preparation%s", buildURLVarChunk(validpreparationinstrumentsservice.ValidPreparationIDURIParamKey, ""))
			validPreparationInstrumentsRouter.Route(validPreparationInstrumentsByPreparationIDRouteParam, func(byValidPreparationIDRouter routing.Router) {
				byValidPreparationIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationInstrumentsPermission)).
					Get(root, s.validPreparationInstrumentsService.SearchByPreparationHandler)
			})

			validPreparationInstrumentsByInstrumentIDRouteParam := fmt.Sprintf("/by_instrument%s", buildURLVarChunk(validpreparationinstrumentsservice.ValidInstrumentIDURIParamKey, ""))
			validPreparationInstrumentsRouter.Route(validPreparationInstrumentsByInstrumentIDRouteParam, func(byValidPreparationIDRouter routing.Router) {
				byValidPreparationIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationInstrumentsPermission)).
					Get(root, s.validPreparationInstrumentsService.SearchByInstrumentHandler)
			})

			validPreparationInstrumentIDRouteParam := buildURLVarChunk(validpreparationinstrumentsservice.ValidPreparationVesselIDURIParamKey, "")
			validPreparationInstrumentsRouter.Route(validPreparationInstrumentIDRouteParam, func(singleValidPreparationInstrumentRouter routing.Router) {
				singleValidPreparationInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationInstrumentsPermission)).
					Get(root, s.validPreparationInstrumentsService.ReadHandler)
				singleValidPreparationInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidPreparationInstrumentsPermission)).
					Put(root, s.validPreparationInstrumentsService.UpdateHandler)
				singleValidPreparationInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidPreparationInstrumentsPermission)).
					Delete(root, s.validPreparationInstrumentsService.ArchiveHandler)
			})
		})

		// ValidPreparationVessels
		validPreparationVesselPath := "valid_preparation_vessels"
		validPreparationVesselsRouteWithPrefix := fmt.Sprintf("/%s", validPreparationVesselPath)
		v1Router.Route(validPreparationVesselsRouteWithPrefix, func(validPreparationVesselsRouter routing.Router) {
			validPreparationVesselsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidPreparationVesselsPermission)).
				Post(root, s.validPreparationVesselsService.CreateHandler)
			validPreparationVesselsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationVesselsPermission)).
				Get(root, s.validPreparationVesselsService.ListHandler)

			validPreparationVesselsByPreparationIDRouteParam := fmt.Sprintf("/by_preparation%s", buildURLVarChunk(validpreparationvesselsservice.ValidPreparationIDURIParamKey, ""))
			validPreparationVesselsRouter.Route(validPreparationVesselsByPreparationIDRouteParam, func(byValidPreparationIDRouter routing.Router) {
				byValidPreparationIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationVesselsPermission)).
					Get(root, s.validPreparationVesselsService.SearchByPreparationHandler)
			})

			validPreparationVesselsByInstrumentIDRouteParam := fmt.Sprintf("/by_vessel%s", buildURLVarChunk(validpreparationvesselsservice.ValidVesselIDURIParamKey, ""))
			validPreparationVesselsRouter.Route(validPreparationVesselsByInstrumentIDRouteParam, func(byValidPreparationIDRouter routing.Router) {
				byValidPreparationIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationVesselsPermission)).
					Get(root, s.validPreparationVesselsService.SearchByVesselHandler)
			})

			validPreparationVesselIDRouteParam := buildURLVarChunk(validpreparationvesselsservice.ValidPreparationVesselIDURIParamKey, "")
			validPreparationVesselsRouter.Route(validPreparationVesselIDRouteParam, func(singleValidPreparationVesselRouter routing.Router) {
				singleValidPreparationVesselRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidPreparationVesselsPermission)).
					Get(root, s.validPreparationVesselsService.ReadHandler)
				singleValidPreparationVesselRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidPreparationVesselsPermission)).
					Put(root, s.validPreparationVesselsService.UpdateHandler)
				singleValidPreparationVesselRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidPreparationVesselsPermission)).
					Delete(root, s.validPreparationVesselsService.ArchiveHandler)
			})
		})

		// ValidIngredientMeasurementUnit
		validIngredientMeasurementUnitPath := "valid_ingredient_measurement_units"
		validIngredientMeasurementUnitRouteWithPrefix := fmt.Sprintf("/%s", validIngredientMeasurementUnitPath)
		validIngredientMeasurementUnitIDRouteParam := buildURLVarChunk(validingredientmeasurementunitsservice.ValidIngredientMeasurementUnitIDURIParamKey, "")
		v1Router.Route(validIngredientMeasurementUnitRouteWithPrefix, func(validIngredientMeasurementUnitRouter routing.Router) {
			validIngredientMeasurementUnitRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateValidIngredientMeasurementUnitsPermission)).
				Post(root, s.validIngredientMeasurementUnitsService.CreateHandler)
			validIngredientMeasurementUnitRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientMeasurementUnitsPermission)).
				Get(root, s.validIngredientMeasurementUnitsService.ListHandler)

			validIngredientMeasurementUnitsByIngredientRouteParam := fmt.Sprintf("/by_ingredient%s", buildURLVarChunk(validingredientmeasurementunitsservice.ValidIngredientIDURIParamKey, ""))
			validIngredientMeasurementUnitRouter.Route(validIngredientMeasurementUnitsByIngredientRouteParam, func(byValidIngredientIDRouter routing.Router) {
				byValidIngredientIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientMeasurementUnitsPermission)).
					Get(root, s.validIngredientMeasurementUnitsService.SearchByIngredientHandler)
			})

			validIngredientMeasurementUnitsByMeasurementUnitRouteParam := fmt.Sprintf("/by_measurement_unit%s", buildURLVarChunk(validingredientmeasurementunitsservice.ValidMeasurementUnitIDURIParamKey, ""))
			validIngredientMeasurementUnitRouter.Route(validIngredientMeasurementUnitsByMeasurementUnitRouteParam, func(byValidIngredientIDRouter routing.Router) {
				byValidIngredientIDRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientMeasurementUnitsPermission)).
					Get(root, s.validIngredientMeasurementUnitsService.SearchByMeasurementUnitHandler)
			})

			validIngredientMeasurementUnitRouter.Route(validIngredientMeasurementUnitIDRouteParam, func(singleValidIngredientMeasurementUnitRouter routing.Router) {
				singleValidIngredientMeasurementUnitRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadValidIngredientMeasurementUnitsPermission)).
					Get(root, s.validIngredientMeasurementUnitsService.ReadHandler)
				singleValidIngredientMeasurementUnitRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateValidIngredientMeasurementUnitsPermission)).
					Put(root, s.validIngredientMeasurementUnitsService.UpdateHandler)
				singleValidIngredientMeasurementUnitRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveValidIngredientMeasurementUnitsPermission)).
					Delete(root, s.validIngredientMeasurementUnitsService.ArchiveHandler)
			})
		})

		// Meals
		mealPath := "meals"
		mealsRouteWithPrefix := fmt.Sprintf("/%s", mealPath)
		mealIDRouteParam := buildURLVarChunk(mealsservice.MealIDURIParamKey, "")
		v1Router.Route(mealsRouteWithPrefix, func(mealsRouter routing.Router) {
			mealsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateMealsPermission)).
				Post(root, s.mealsService.CreateMealHandler)
			mealsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealsPermission)).
				Get(root, s.mealsService.ListMealsHandler)
			mealsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealsPermission)).
				Get(searchRoot, s.mealsService.SearchMealsHandler)

			mealsRouter.Route(mealIDRouteParam, func(singleMealRouter routing.Router) {
				singleMealRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealsPermission)).
					Get(root, s.mealsService.ReadMealHandler)
				singleMealRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveMealsPermission)).
					Delete(root, s.mealsService.ArchiveMealHandler)
			})
		})

		// Components
		recipePath := "recipes"
		recipesRouteWithPrefix := fmt.Sprintf("/%s", recipePath)
		recipeIDRouteParam := buildURLVarChunk(recipesservice.RecipeIDURIParamKey, "")
		v1Router.Route(recipesRouteWithPrefix, func(recipesRouter routing.Router) {
			recipesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipesPermission)).
				Post(root, s.recipesService.CreateHandler)
			recipesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
				Get(root, s.recipesService.ListHandler)
			recipesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
				Get(searchRoot, s.recipesService.SearchHandler)

			recipesRouter.Route(recipeIDRouteParam, func(singleRecipeRouter routing.Router) {
				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
					Get(root, s.recipesService.ReadHandler)
				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
					Get("/dag", s.recipesService.DAGHandler)
				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
					Get("/mermaid", s.recipesService.MermaidHandler)
				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
					Get("/prep_steps", s.recipesService.EstimatedPrepStepsHandler)
				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipesPermission)).
					Put(root, s.recipesService.UpdateHandler)
				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipesPermission)).
					Post("/clone", s.recipesService.CloneHandler)

				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipesPermission)).
					Post("/images", s.recipesService.ImageUploadHandler)
				singleRecipeRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipesPermission)).
					Delete(root, s.recipesService.ArchiveHandler)

				// RecipeRatings
				singleRecipeRouter.Route("/ratings", func(recipeRatingsRouter routing.Router) {
					recipeRatingsRouter.
						WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipeRatingsPermission)).
						Post(root, s.recipeRatingsService.CreateHandler)
					recipeRatingsRouter.
						WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeRatingsPermission)).
						Get(root, s.recipeRatingsService.ListHandler)

					recipeRatingIDRouteParam := buildURLVarChunk(reciperatingsservice.RecipeRatingIDURIParamKey, "")
					recipeRatingsRouter.Route(recipeRatingIDRouteParam, func(singleRecipeRatingRouter routing.Router) {
						singleRecipeRatingRouter.
							WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipeRatingsPermission)).
							Get(root, s.recipeRatingsService.ReadHandler)
						singleRecipeRatingRouter.
							WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeRatingsPermission)).
							Put(root, s.recipeRatingsService.UpdateHandler)
						singleRecipeRatingRouter.
							WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipeRatingsPermission)).
							Delete(root, s.recipeRatingsService.ArchiveHandler)
					})
				})
			})
		})

		// PrepTasks
		recipePrepTaskPath := "prep_tasks"
		recipePrepTasksRoute := path.Join(
			recipePath,
			recipeIDRouteParam,
			recipePrepTaskPath,
		)
		recipePrepTasksRouteWithPrefix := fmt.Sprintf("/%s", recipePrepTasksRoute)
		recipePrepTaskIDRouteParam := buildURLVarChunk(recipepreptasksservice.RecipePrepTaskIDURIParamKey, "")
		v1Router.Route(recipePrepTasksRouteWithPrefix, func(recipePrepTasksRouter routing.Router) {
			recipePrepTasksRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipePrepTasksPermission)).
				Post(root, s.recipePrepTasksService.CreateHandler)
			recipePrepTasksRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipePrepTasksPermission)).
				Get(root, s.recipePrepTasksService.ListHandler)

			recipePrepTasksRouter.Route(recipePrepTaskIDRouteParam, func(singleRecipeStepRouter routing.Router) {
				singleRecipeStepRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipePrepTasksPermission)).
					Get(root, s.recipePrepTasksService.ReadHandler)
				singleRecipeStepRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipePrepTasksPermission)).
					Put(root, s.recipePrepTasksService.UpdateHandler)
				singleRecipeStepRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipePrepTasksPermission)).
					Delete(root, s.recipePrepTasksService.ArchiveHandler)
			})
		})

		// TaskSteps
		recipeStepPath := "steps"
		recipeStepsRoute := path.Join(
			recipePath,
			recipeIDRouteParam,
			recipeStepPath,
		)
		recipeStepsRouteWithPrefix := fmt.Sprintf("/%s", recipeStepsRoute)
		recipeStepIDRouteParam := buildURLVarChunk(recipestepsservice.RecipeStepIDURIParamKey, "")
		v1Router.Route(recipeStepsRouteWithPrefix, func(recipeStepsRouter routing.Router) {
			recipeStepsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipeStepsPermission)).
				Post(root, s.recipeStepsService.CreateHandler)
			recipeStepsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepsPermission)).
				Get(root, s.recipeStepsService.ListHandler)

			recipeStepsRouter.Route(recipeStepIDRouteParam, func(singleRecipeStepRouter routing.Router) {
				singleRecipeStepRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepsPermission)).
					Post("/images", s.recipeStepsService.ImageUploadHandler)
				singleRecipeStepRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepsPermission)).
					Get(root, s.recipeStepsService.ReadHandler)
				singleRecipeStepRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepsPermission)).
					Put(root, s.recipeStepsService.UpdateHandler)
				singleRecipeStepRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepsPermission)).
					Delete(root, s.recipeStepsService.ArchiveHandler)
			})
		})

		// RecipeStepInstruments
		recipeStepInstrumentPath := "instruments"
		recipeStepInstrumentsRoute := path.Join(
			recipePath,
			recipeIDRouteParam,
			recipeStepPath,
			recipeStepIDRouteParam,
			recipeStepInstrumentPath,
		)
		recipeStepInstrumentsRouteWithPrefix := fmt.Sprintf("/%s", recipeStepInstrumentsRoute)
		recipeStepInstrumentIDRouteParam := buildURLVarChunk(recipestepinstrumentsservice.RecipeStepInstrumentIDURIParamKey, "")
		v1Router.Route(recipeStepInstrumentsRouteWithPrefix, func(recipeStepInstrumentsRouter routing.Router) {
			recipeStepInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipeStepInstrumentsPermission)).
				Post(root, s.recipeStepInstrumentsService.CreateHandler)
			recipeStepInstrumentsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepInstrumentsPermission)).
				Get(root, s.recipeStepInstrumentsService.ListHandler)

			recipeStepInstrumentsRouter.Route(recipeStepInstrumentIDRouteParam, func(singleRecipeStepInstrumentRouter routing.Router) {
				singleRecipeStepInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepInstrumentsPermission)).
					Get(root, s.recipeStepInstrumentsService.ReadHandler)
				singleRecipeStepInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepInstrumentsPermission)).
					Put(root, s.recipeStepInstrumentsService.UpdateHandler)
				singleRecipeStepInstrumentRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepInstrumentsPermission)).
					Delete(root, s.recipeStepInstrumentsService.ArchiveHandler)
			})
		})

		// RecipeStepVessels
		recipeStepVesselPath := "vessels"
		recipeStepVesselsRoute := path.Join(
			recipePath,
			recipeIDRouteParam,
			recipeStepPath,
			recipeStepIDRouteParam,
			recipeStepVesselPath,
		)
		recipeStepVesselsRouteWithPrefix := fmt.Sprintf("/%s", recipeStepVesselsRoute)
		recipeStepVesselIDRouteParam := buildURLVarChunk(recipestepvesselsservice.RecipeStepVesselIDURIParamKey, "")
		v1Router.Route(recipeStepVesselsRouteWithPrefix, func(recipeStepVesselsRouter routing.Router) {
			recipeStepVesselsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipeStepVesselsPermission)).
				Post(root, s.recipeStepVesselsService.CreateHandler)
			recipeStepVesselsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepVesselsPermission)).
				Get(root, s.recipeStepVesselsService.ListHandler)

			recipeStepVesselsRouter.Route(recipeStepVesselIDRouteParam, func(singleRecipeStepVesselRouter routing.Router) {
				singleRecipeStepVesselRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepVesselsPermission)).
					Get(root, s.recipeStepVesselsService.ReadHandler)
				singleRecipeStepVesselRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepVesselsPermission)).
					Put(root, s.recipeStepVesselsService.UpdateHandler)
				singleRecipeStepVesselRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepVesselsPermission)).
					Delete(root, s.recipeStepVesselsService.ArchiveHandler)
			})
		})

		// RecipeStepIngredients
		recipeStepIngredientPath := "ingredients"
		recipeStepIngredientsRoute := path.Join(
			recipePath,
			recipeIDRouteParam,
			recipeStepPath,
			recipeStepIDRouteParam,
			recipeStepIngredientPath,
		)
		recipeStepIngredientsRouteWithPrefix := fmt.Sprintf("/%s", recipeStepIngredientsRoute)
		recipeStepIngredientIDRouteParam := buildURLVarChunk(recipestepingredientsservice.RecipeStepIngredientIDURIParamKey, "")
		v1Router.Route(recipeStepIngredientsRouteWithPrefix, func(recipeStepIngredientsRouter routing.Router) {
			recipeStepIngredientsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipeStepIngredientsPermission)).
				Post(root, s.recipeStepIngredientsService.CreateHandler)
			recipeStepIngredientsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepIngredientsPermission)).
				Get(root, s.recipeStepIngredientsService.ListHandler)

			recipeStepIngredientsRouter.Route(recipeStepIngredientIDRouteParam, func(singleRecipeStepIngredientRouter routing.Router) {
				singleRecipeStepIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepIngredientsPermission)).
					Get(root, s.recipeStepIngredientsService.ReadHandler)
				singleRecipeStepIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepIngredientsPermission)).
					Put(root, s.recipeStepIngredientsService.UpdateHandler)
				singleRecipeStepIngredientRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepIngredientsPermission)).
					Delete(root, s.recipeStepIngredientsService.ArchiveHandler)
			})
		})

		// RecipeStepIngredients
		recipeStepCompletionConditionPath := "completion_conditions"
		recipeStepCompletionConditionsRoute := path.Join(
			recipePath,
			recipeIDRouteParam,
			recipeStepPath,
			recipeStepIDRouteParam,
			recipeStepCompletionConditionPath,
		)
		recipeStepCompletionConditionsRouteWithPrefix := fmt.Sprintf("/%s", recipeStepCompletionConditionsRoute)
		recipeStepCompletionConditionIDRouteParam := buildURLVarChunk(recipestepcompletionconditionsservice.RecipeStepCompletionConditionIDURIParamKey, "")
		v1Router.Route(recipeStepCompletionConditionsRouteWithPrefix, func(recipeStepCompletionConditionsRouter routing.Router) {
			recipeStepCompletionConditionsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipeStepCompletionConditionsPermission)).
				Post(root, s.recipeStepCompletionConditionsService.CreateHandler)
			recipeStepCompletionConditionsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepCompletionConditionsPermission)).
				Get(root, s.recipeStepCompletionConditionsService.ListHandler)

			recipeStepCompletionConditionsRouter.Route(recipeStepCompletionConditionIDRouteParam, func(singleRecipeStepCompletionConditionRouter routing.Router) {
				singleRecipeStepCompletionConditionRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepCompletionConditionsPermission)).
					Get(root, s.recipeStepCompletionConditionsService.ReadHandler)
				singleRecipeStepCompletionConditionRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepCompletionConditionsPermission)).
					Put(root, s.recipeStepCompletionConditionsService.UpdateHandler)
				singleRecipeStepCompletionConditionRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepCompletionConditionsPermission)).
					Delete(root, s.recipeStepCompletionConditionsService.ArchiveHandler)
			})
		})

		// RecipeStepProducts
		recipeStepProductPath := "products"
		recipeStepProductsRoute := path.Join(
			recipePath,
			recipeIDRouteParam,
			recipeStepPath,
			recipeStepIDRouteParam,
			recipeStepProductPath,
		)
		recipeStepProductsRouteWithPrefix := fmt.Sprintf("/%s", recipeStepProductsRoute)
		recipeStepProductIDRouteParam := buildURLVarChunk(recipestepproductsservice.RecipeStepProductIDURIParamKey, "")
		v1Router.Route(recipeStepProductsRouteWithPrefix, func(recipeStepProductsRouter routing.Router) {
			recipeStepProductsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateRecipeStepProductsPermission)).
				Post(root, s.recipeStepProductsService.CreateHandler)
			recipeStepProductsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepProductsPermission)).
				Get(root, s.recipeStepProductsService.ListHandler)

			recipeStepProductsRouter.Route(recipeStepProductIDRouteParam, func(singleRecipeStepProductRouter routing.Router) {
				singleRecipeStepProductRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadRecipeStepProductsPermission)).
					Get(root, s.recipeStepProductsService.ReadHandler)
				singleRecipeStepProductRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateRecipeStepProductsPermission)).
					Put(root, s.recipeStepProductsService.UpdateHandler)
				singleRecipeStepProductRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveRecipeStepProductsPermission)).
					Delete(root, s.recipeStepProductsService.ArchiveHandler)
			})
		})

		// MealPlans
		mealPlanPath := "meal_plans"
		mealPlansRouteWithPrefix := fmt.Sprintf("/%s", mealPlanPath)
		mealPlanIDRouteParam := buildURLVarChunk(mealplansservice.MealPlanIDURIParamKey, "")
		v1Router.Route(mealPlansRouteWithPrefix, func(mealPlansRouter routing.Router) {
			mealPlansRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateMealPlansPermission)).
				Post(root, s.mealPlansService.CreateHandler)
			mealPlansRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlansPermission)).
				Get(root, s.mealPlansService.ListHandler)

			mealPlansRouter.Route(mealPlanIDRouteParam, func(singleMealPlanRouter routing.Router) {
				singleMealPlanRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlansPermission)).
					Get(root, s.mealPlansService.ReadHandler)
				singleMealPlanRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateMealPlansPermission)).
					Put(root, s.mealPlansService.UpdateHandler)
				singleMealPlanRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveMealPlansPermission)).
					Delete(root, s.mealPlansService.ArchiveHandler)
				singleMealPlanRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateMealPlansPermission)).
					Post("/finalize", s.mealPlansService.FinalizeHandler)
			})
		})

		// MealPlanTasks
		mealPlanTaskPath := "tasks"
		mealPlanTasksRoute := path.Join(
			mealPlanPath,
			mealPlanIDRouteParam,
			mealPlanTaskPath,
		)
		mealPlanTasksRouteWithPrefix := fmt.Sprintf("/%s", mealPlanTasksRoute)
		mealPlanTaskIDRouteParam := buildURLVarChunk(mealplantasksservice.MealPlanTaskIDURIParamKey, "")
		v1Router.Route(mealPlanTasksRouteWithPrefix, func(mealPlanTasksRouter routing.Router) {
			mealPlanTasksRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlanTasksPermission)).
				Get(root, s.mealPlanTasksService.ListByMealPlanHandler)
			mealPlanTasksRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateMealPlanTasksPermission)).
				Post(root, s.mealPlanTasksService.CreateHandler)

			mealPlanTasksRouter.Route(mealPlanTaskIDRouteParam, func(singleMealPlanTaskRouter routing.Router) {
				singleMealPlanTaskRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlanTasksPermission)).
					Get(root, s.mealPlanTasksService.ReadHandler)

				singleMealPlanTaskRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateMealPlanTasksPermission)).
					Patch(root, s.mealPlanTasksService.StatusChangeHandler)
			})
		})

		// MealPlanEvents
		mealPlanEventPath := "events"
		mealPlanEventsRoute := path.Join(
			mealPlanPath,
			mealPlanIDRouteParam,
			mealPlanEventPath,
		)
		mealPlanEventsRouteWithPrefix := fmt.Sprintf("/%s", mealPlanEventsRoute)
		mealPlanEventIDRouteParam := buildURLVarChunk(mealplaneventsservice.MealPlanEventIDURIParamKey, "")
		v1Router.Route(mealPlanEventsRouteWithPrefix, func(mealPlanEventsRouter routing.Router) {
			mealPlanEventsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateMealPlanEventsPermission)).
				Post(root, s.mealPlanEventsService.CreateHandler)
			mealPlanEventsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlanEventsPermission)).
				Get(root, s.mealPlanEventsService.ListHandler)

			mealPlanEventsRouter.Route(mealPlanEventIDRouteParam, func(singleMealPlanEventRouter routing.Router) {
				singleMealPlanEventRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlanEventsPermission)).
					Get(root, s.mealPlanEventsService.ReadHandler)
				singleMealPlanEventRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateMealPlanEventsPermission)).
					Put(root, s.mealPlanEventsService.UpdateHandler)
				singleMealPlanEventRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateMealPlanOptionVotesPermission)).
					Post("/vote", s.mealPlanOptionVotesService.CreateHandler)
				singleMealPlanEventRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveMealPlanEventsPermission)).
					Delete(root, s.mealPlanEventsService.ArchiveHandler)
			})
		})

		// MealPlanGroceryListItems
		mealPlanGroceryListItemPath := "grocery_list_items"
		mealPlanGroceryListItemsRoute := path.Join(
			mealPlanPath,
			mealPlanIDRouteParam,
			mealPlanGroceryListItemPath,
		)
		mealPlanGroceryListItemsRouteWithPrefix := fmt.Sprintf("/%s", mealPlanGroceryListItemsRoute)
		mealPlanGroceryListItemIDRouteParam := buildURLVarChunk(mealplangrocerylistitemssservice.MealPlanGroceryListItemIDURIParamKey, "")
		v1Router.Route(mealPlanGroceryListItemsRouteWithPrefix, func(mealPlanGroceryListItemsRouter routing.Router) {
			mealPlanGroceryListItemsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateMealPlanGroceryListItemsPermission)).
				Post(root, s.mealPlanGroceryListItemsService.CreateHandler)
			mealPlanGroceryListItemsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlanGroceryListItemsPermission)).
				Get(root, s.mealPlanGroceryListItemsService.ListByMealPlanHandler)

			mealPlanGroceryListItemsRouter.Route(mealPlanGroceryListItemIDRouteParam, func(singleMealPlanGroceryListItemRouter routing.Router) {
				singleMealPlanGroceryListItemRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlanGroceryListItemsPermission)).
					Get(root, s.mealPlanGroceryListItemsService.ReadHandler)
				singleMealPlanGroceryListItemRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateMealPlanGroceryListItemsPermission)).
					Put(root, s.mealPlanGroceryListItemsService.UpdateHandler)
				singleMealPlanGroceryListItemRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveMealPlanGroceryListItemsPermission)).
					Delete(root, s.mealPlanGroceryListItemsService.ArchiveHandler)
			})
		})

		// MealPlanOptions
		mealPlanOptionPath := "options"
		mealPlanOptionsRoute := path.Join(
			mealPlanPath,
			mealPlanIDRouteParam,
			mealPlanEventPath,
			mealPlanEventIDRouteParam,
			mealPlanOptionPath,
		)
		mealPlanOptionsRouteWithPrefix := fmt.Sprintf("/%s", mealPlanOptionsRoute)
		mealPlanOptionIDRouteParam := buildURLVarChunk(mealplanoptionsservice.MealPlanOptionIDURIParamKey, "")
		v1Router.Route(mealPlanOptionsRouteWithPrefix, func(mealPlanOptionsRouter routing.Router) {
			mealPlanOptionsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateMealPlanOptionsPermission)).
				Post(root, s.mealPlanOptionsService.CreateHandler)
			mealPlanOptionsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlanOptionsPermission)).
				Get(root, s.mealPlanOptionsService.ListHandler)

			mealPlanOptionsRouter.Route(mealPlanOptionIDRouteParam, func(singleMealPlanOptionRouter routing.Router) {
				singleMealPlanOptionRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlanOptionsPermission)).
					Get(root, s.mealPlanOptionsService.ReadHandler)
				singleMealPlanOptionRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateMealPlanOptionsPermission)).
					Put(root, s.mealPlanOptionsService.UpdateHandler)
				singleMealPlanOptionRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveMealPlanOptionsPermission)).
					Delete(root, s.mealPlanOptionsService.ArchiveHandler)
			})
		})

		// MealPlanOptionVotes
		mealPlanOptionVotePath := "votes"
		mealPlanOptionVotesRoute := path.Join(
			mealPlanPath,
			mealPlanIDRouteParam,
			mealPlanEventPath,
			mealPlanEventIDRouteParam,
			mealPlanOptionPath,
			mealPlanOptionIDRouteParam,
			mealPlanOptionVotePath,
		)
		mealPlanOptionVotesRouteWithPrefix := fmt.Sprintf("/%s", mealPlanOptionVotesRoute)
		mealPlanOptionVoteIDRouteParam := buildURLVarChunk(mealplanoptionvotesservice.MealPlanOptionVoteIDURIParamKey, "")
		v1Router.Route(mealPlanOptionVotesRouteWithPrefix, func(mealPlanOptionVotesRouter routing.Router) {
			mealPlanOptionVotesRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlanOptionVotesPermission)).
				Get(root, s.mealPlanOptionVotesService.ListHandler)

			mealPlanOptionVotesRouter.Route(mealPlanOptionVoteIDRouteParam, func(singleMealPlanOptionVoteRouter routing.Router) {
				singleMealPlanOptionVoteRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadMealPlanOptionVotesPermission)).
					Get(root, s.mealPlanOptionVotesService.ReadHandler)
				singleMealPlanOptionVoteRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateMealPlanOptionVotesPermission)).
					Put(root, s.mealPlanOptionVotesService.UpdateHandler)
				singleMealPlanOptionVoteRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveMealPlanOptionVotesPermission)).
					Delete(root, s.mealPlanOptionVotesService.ArchiveHandler)
			})
		})

		// ServiceSettings
		serviceSettingPath := "settings"
		serviceSettingsRouteWithPrefix := fmt.Sprintf("/%s", serviceSettingPath)
		serviceSettingIDRouteParam := buildURLVarChunk(servicesettingsservice.ServiceSettingIDURIParamKey, "")
		v1Router.Route(serviceSettingsRouteWithPrefix, func(serviceSettingsRouter routing.Router) {
			serviceSettingsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateServiceSettingsPermission)).
				Post(root, s.serviceSettingsService.CreateHandler)
			serviceSettingsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadServiceSettingsPermission)).
				Get(root, s.serviceSettingsService.ListHandler)
			serviceSettingsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadServiceSettingsPermission)).
				Get(searchRoot, s.serviceSettingsService.SearchHandler)

			serviceSettingsRouter.Route(serviceSettingIDRouteParam, func(singleServiceSettingRouter routing.Router) {
				singleServiceSettingRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadServiceSettingsPermission)).
					Get(root, s.serviceSettingsService.ReadHandler)
				singleServiceSettingRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveServiceSettingsPermission)).
					Delete(root, s.serviceSettingsService.ArchiveHandler)
			})

			serviceSettingConfigurationIDRouteParam := buildURLVarChunk(servicesettingconfigurationsservice.ServiceSettingConfigurationIDURIParamKey, "")
			serviceSettingsRouter.Route("/configurations", func(settingConfigurationRouter routing.Router) {
				settingConfigurationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateServiceSettingConfigurationsPermission)).
					Post(root, s.serviceSettingConfigurationsService.CreateHandler)
				settingConfigurationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadServiceSettingConfigurationsPermission)).
					Get("/user"+buildURLVarChunk(servicesettingconfigurationsservice.ServiceSettingConfigurationNameURIParamKey, ""), s.serviceSettingConfigurationsService.ForUserByNameHandler)
				settingConfigurationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadServiceSettingConfigurationsPermission)).
					Get("/user", s.serviceSettingConfigurationsService.ForUserHandler)
				settingConfigurationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadServiceSettingConfigurationsPermission)).
					Get("/household", s.serviceSettingConfigurationsService.ForHouseholdHandler)
				settingConfigurationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateServiceSettingConfigurationsPermission)).
					Put(serviceSettingConfigurationIDRouteParam, s.serviceSettingConfigurationsService.UpdateHandler)
				settingConfigurationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ArchiveServiceSettingConfigurationsPermission)).
					Delete(serviceSettingConfigurationIDRouteParam, s.serviceSettingConfigurationsService.ArchiveHandler)
			})
		})

		// User Notifications
		userNotificationPath := "user_notifications"
		userNotificationsRouteWithPrefix := fmt.Sprintf("/%s", userNotificationPath)
		userNotificationIDRouteParam := buildURLVarChunk(usernotificationsservice.UserNotificationIDURIParamKey, "")
		v1Router.Route(userNotificationsRouteWithPrefix, func(userNotificationsRouter routing.Router) {
			userNotificationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.CreateUserNotificationsPermission)).
				Post(root, s.userNotificationsService.CreateHandler)
			userNotificationsRouter.
				WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadUserNotificationsPermission)).
				Get(root, s.userNotificationsService.ListHandler)

			userNotificationsRouter.Route(userNotificationIDRouteParam, func(singleUserNotificationRouter routing.Router) {
				singleUserNotificationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.ReadUserNotificationsPermission)).
					Get(root, s.userNotificationsService.ReadHandler)
				singleUserNotificationRouter.
					WithMiddleware(s.authService.PermissionFilterMiddleware(authorization.UpdateUserNotificationsPermission)).
					Patch(root, s.userNotificationsService.UpdateHandler)
			})
		})
	})

	s.router = router
}
