package email

import (
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestBuildGeneratedPasswordResetTokenEmail(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		user := fakes.BuildFakeUser()
		user.EmailAddressVerifiedAt = pointer.To(time.Now())
		token := fakes.BuildFakePasswordResetToken()

		actual, err := BuildGeneratedPasswordResetTokenEmail(user, token, envConfigsMap[defaultEnv])
		assert.NoError(t, err)
		assert.NotNil(t, actual)
		assert.Contains(t, actual.HTMLContent, logoURL)
	})
}

func TestBuildInviteMemberEmail(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		invitation := fakes.BuildFakeHouseholdInvitation()

		actual, err := BuildInviteMemberEmail(invitation, envConfigsMap[defaultEnv])
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})
}

func TestBuildMealPlanCreatedEmail(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		user := fakes.BuildFakeUser()
		user.EmailAddressVerifiedAt = pointer.To(time.Now())
		mealPlan := fakes.BuildFakeMealPlan()

		actual, err := BuildMealPlanCreatedEmail(user, mealPlan, envConfigsMap[defaultEnv])
		assert.NoError(t, err)
		assert.NotNil(t, actual)
		assert.Contains(t, actual.HTMLContent, logoURL)
	})
}

func TestBuildPasswordResetTokenRedeemedEmail(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		user := fakes.BuildFakeUser()
		user.EmailAddressVerifiedAt = pointer.To(time.Now())

		actual, err := BuildPasswordResetTokenRedeemedEmail(user, envConfigsMap[defaultEnv])
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})
}

func TestBuildUsernameReminderEmail(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		user := fakes.BuildFakeUser()
		user.EmailAddressVerifiedAt = pointer.To(time.Now())

		actual, err := BuildUsernameReminderEmail(user, envConfigsMap[defaultEnv])
		assert.NoError(t, err)
		assert.NotNil(t, actual)
	})
}
