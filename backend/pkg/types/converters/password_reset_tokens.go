package converters

import (
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ConvertPasswordResetTokenToPasswordResetTokenDatabaseCreationInput builds a PasswordResetTokenDatabaseCreationInput from a PasswordResetToken.
func ConvertPasswordResetTokenToPasswordResetTokenDatabaseCreationInput(input *types.PasswordResetToken) *types.PasswordResetTokenDatabaseCreationInput {
	return &types.PasswordResetTokenDatabaseCreationInput{
		ID:            input.ID,
		Token:         input.Token,
		BelongsToUser: input.BelongsToUser,
		ExpiresAt:     input.ExpiresAt,
	}
}
