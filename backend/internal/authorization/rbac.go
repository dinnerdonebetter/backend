package authorization

import (
	"github.com/mikespook/gorbac/v2"
)

type (
	role int
)

var (
	globalAuthorizer *gorbac.RBAC
)

func init() {
	globalAuthorizer = initializeRBAC()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func initializeRBAC() *gorbac.RBAC {
	rbac := gorbac.New()

	must(rbac.Add(serviceUser))
	must(rbac.Add(serviceAdmin))
	must(rbac.Add(householdAdmin))
	must(rbac.Add(householdMember))

	must(rbac.SetParent(householdAdminRoleName, householdMemberRoleName))
	must(rbac.SetParent(serviceAdminRoleName, householdAdminRoleName))

	return rbac
}
