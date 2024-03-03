package authorization

import (
	"encoding/gob"

	"github.com/mikespook/gorbac/v2"
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

func hasPermission(p Permission, roles ...string) bool {
	if len(roles) == 0 {
		return false
	}

	for _, r := range roles {
		if globalAuthorizer.IsGranted(r, p, nil) {
			return true
		}
	}

	return false
}
