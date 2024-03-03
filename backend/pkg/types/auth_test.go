package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"

	"github.com/stretchr/testify/assert"
)

func TestChangeActiveHouseholdInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ChangeActiveHouseholdInput{
			HouseholdID: "123",
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestSessionContextData_AttachToLogger(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &SessionContextData{
			Requester: RequesterInfo{ServicePermissions: authorization.NewServiceRolePermissionChecker(t.Name())},
		}

		assert.NotNil(t, x.AttachToLogger(logging.NewNoopLogger()))
	})
}

func TestSessionContextData_HouseholdRolePermissionsChecker(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &SessionContextData{
			ActiveHouseholdID: t.Name(),
			HouseholdPermissions: map[string]authorization.HouseholdRolePermissionsChecker{
				t.Name(): authorization.NewHouseholdRolePermissionChecker(t.Name()),
			},
		}

		assert.NotNil(t, x.HouseholdRolePermissionsChecker())
	})
}

func TestSessionContextData_ServiceRolePermissionChecker(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &SessionContextData{
			ActiveHouseholdID: t.Name(),
			Requester: RequesterInfo{
				ServicePermissions: authorization.NewServiceRolePermissionChecker(t.Name()),
			},
		}

		assert.NotNil(t, x.ServiceRolePermissionChecker())
	})
}
