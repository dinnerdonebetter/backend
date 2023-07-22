package authorization

import (
	"encoding/gob"
)

type (
	// HouseholdRole describes a role a user has for a household context.
	HouseholdRole role

	// HouseholdRolePermissionsChecker checks permissions for one or more household Roles.
	HouseholdRolePermissionsChecker interface {
		HasPermission(Permission) bool
	}
)

const (
	// HouseholdMemberRole is a role for a plain household participant.
	HouseholdMemberRole HouseholdRole = iota
	// HouseholdAdminRole is a role for someone who can manipulate the specifics of a household.
	HouseholdAdminRole HouseholdRole = iota

	householdAdminRoleName  = "household_admin"
	householdMemberRoleName = "household_member"
)

var (
	householdAdmin  = gorbac.NewStdRole(householdAdminRoleName)
	householdMember = gorbac.NewStdRole(householdMemberRoleName)
)

type householdRoleCollection struct {
	Roles []string
}

func init() {
	gob.Register(householdRoleCollection{})
}

// NewHouseholdRolePermissionChecker returns a new checker for a set of Roles.
func NewHouseholdRolePermissionChecker(roles ...string) HouseholdRolePermissionsChecker {
	return &householdRoleCollection{
		Roles: roles,
	}
}

func (r HouseholdRole) String() string {
	switch r {
	case HouseholdMemberRole:
		return householdMemberRoleName
	case HouseholdAdminRole:
		return householdAdminRoleName
	default:
		return ""
	}
}

// HasPermission returns whether a user can do something or not.
func (r householdRoleCollection) HasPermission(p Permission) bool {
	return hasPermission(p, r.Roles...)
}

// CanUpdateHouseholds returns whether a user can update households or not.
func (r householdRoleCollection) CanUpdateHouseholds() bool {
	return hasPermission(UpdateHouseholdPermission, r.Roles...)
}

// CanDeleteHouseholds returns whether a user can delete households or not.
func (r householdRoleCollection) CanDeleteHouseholds() bool {
	return hasPermission(ArchiveHouseholdPermission, r.Roles...)
}

// CanAddMemberToHouseholds returns whether a user can add members to households or not.
func (r householdRoleCollection) CanAddMemberToHouseholds() bool {
	return hasPermission(InviteUserToHouseholdPermission, r.Roles...)
}

// CanRemoveMemberFromHouseholds returns whether a user can remove members from households or not.
func (r householdRoleCollection) CanRemoveMemberFromHouseholds() bool {
	return hasPermission(RemoveMemberHouseholdPermission, r.Roles...)
}

// CanTransferHouseholdToNewOwner returns whether a user can transfer a household to a new owner or not.
func (r householdRoleCollection) CanTransferHouseholdToNewOwner() bool {
	return hasPermission(TransferHouseholdPermission, r.Roles...)
}

// CanCreateWebhooks returns whether a user can create webhooks or not.
func (r householdRoleCollection) CanCreateWebhooks() bool {
	return hasPermission(CreateWebhooksPermission, r.Roles...)
}

// CanSeeWebhooks returns whether a user can view webhooks or not.
func (r householdRoleCollection) CanSeeWebhooks() bool {
	return hasPermission(ReadWebhooksPermission, r.Roles...)
}

// CanUpdateWebhooks returns whether a user can update webhooks or not.
func (r householdRoleCollection) CanUpdateWebhooks() bool {
	return hasPermission(UpdateWebhooksPermission, r.Roles...)
}

// CanArchiveWebhooks returns whether a user can delete webhooks or not.
func (r householdRoleCollection) CanArchiveWebhooks() bool {
	return hasPermission(ArchiveWebhooksPermission, r.Roles...)
}

// PermissionSummary renders a permission summary.
func (r householdRoleCollection) PermissionSummary() map[string]bool {
	return map[string]bool{
		"CanUpdateHouseholds":            r.CanUpdateHouseholds(),
		"CanDeleteHouseholds":            r.CanDeleteHouseholds(),
		"CanAddMemberToHouseholds":       r.CanAddMemberToHouseholds(),
		"CanRemoveMemberFromHouseholds":  r.CanRemoveMemberFromHouseholds(),
		"CanTransferHouseholdToNewOwner": r.CanTransferHouseholdToNewOwner(),
		"CanCreateWebhooks":              r.CanCreateWebhooks(),
		"CanSeeWebhooks":                 r.CanSeeWebhooks(),
		"CanUpdateWebhooks":              r.CanUpdateWebhooks(),
		"CanArchiveWebhooks":             r.CanArchiveWebhooks(),
	}
}
