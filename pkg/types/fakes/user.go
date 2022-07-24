package fakes

import (
	"encoding/base32"
	"fmt"
	"log"
	"time"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/pquerna/otp/totp"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/internal/authorization"
	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeUser builds a faked User.
func BuildFakeUser() *types.User {
	return &types.User{
		ID:                        ksuid.New().String(),
		EmailAddress:              fake.Email(),
		Username:                  fake.Password(true, true, true, false, false, 32),
		BirthDay:                  func(x uint8) *uint8 { return &x }(fake.Uint8()),
		BirthMonth:                func(x uint8) *uint8 { return &x }(fake.Uint8()),
		AccountStatus:             types.GoodStandingUserAccountStatus,
		TwoFactorSecret:           base32.StdEncoding.EncodeToString([]byte(fake.Password(false, true, true, false, false, 32))),
		TwoFactorSecretVerifiedOn: func(i uint64) *uint64 { return &i }(uint64(uint32(fake.Date().Unix()))),
		ServiceRoles:              []string{authorization.ServiceUserRole.String()},
		CreatedOn:                 uint64(uint32(fake.Date().Unix())),
	}
}

// BuildUserCreationResponseFromUser builds a faked UserCreationResponse.
func BuildUserCreationResponseFromUser(user *types.User) *types.UserCreationResponse {
	return &types.UserCreationResponse{
		CreatedUserID: user.ID,
		Username:      user.Username,
		CreatedOn:     user.CreatedOn,
	}
}

// BuildFakeUserList builds a faked UserList.
func BuildFakeUserList() *types.UserList {
	var examples []*types.User
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeUser())
	}

	return &types.UserList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Users: examples,
	}
}

// BuildFakeUserCreationInput builds a faked UserRegistrationInput.
func BuildFakeUserCreationInput() *types.UserRegistrationInput {
	exampleUser := BuildFakeUser()

	return &types.UserRegistrationInput{
		Username:     exampleUser.Username,
		EmailAddress: fake.Email(),
		Password:     fake.Password(true, true, true, true, false, 32),
		BirthDay:     exampleUser.BirthDay,
		BirthMonth:   exampleUser.BirthMonth,
	}
}

// BuildFakeUserRegistrationInputFromUser builds a faked UserRegistrationInput.
func BuildFakeUserRegistrationInputFromUser(user *types.User) *types.UserRegistrationInput {
	return &types.UserRegistrationInput{
		Username:     user.Username,
		EmailAddress: user.EmailAddress,
		Password:     fake.Password(true, true, true, true, false, 32),
		BirthDay:     user.BirthDay,
		BirthMonth:   user.BirthMonth,
	}
}

// BuildFakeUserRegistrationInputWithInviteFromUser builds a faked UserRegistrationInput.
func BuildFakeUserRegistrationInputWithInviteFromUser(user *types.User) *types.UserRegistrationInput {
	return &types.UserRegistrationInput{
		Username:        user.Username,
		EmailAddress:    user.EmailAddress,
		Password:        fake.Password(true, true, true, true, false, 32),
		BirthDay:        user.BirthDay,
		BirthMonth:      user.BirthMonth,
		InvitationToken: fake.UUID(),
		InvitationID:    BuildFakeID(),
	}
}

// BuildFakeUserDataStoreCreationInputFromUser builds a faked UserDatabaseCreationInput.
func BuildFakeUserDataStoreCreationInputFromUser(user *types.User) *types.UserDatabaseCreationInput {
	return &types.UserDatabaseCreationInput{
		ID:              user.ID,
		EmailAddress:    user.EmailAddress,
		Username:        user.Username,
		HashedPassword:  user.HashedPassword,
		TwoFactorSecret: user.TwoFactorSecret,
		BirthDay:        user.BirthDay,
		BirthMonth:      user.BirthMonth,
	}
}

// BuildFakeUserAccountStatusUpdateInputFromUser builds a faked UserAccountStatusUpdateInput.
func BuildFakeUserAccountStatusUpdateInputFromUser(user *types.User) *types.UserAccountStatusUpdateInput {
	return &types.UserAccountStatusUpdateInput{
		TargetUserID: ksuid.New().String(),
		NewStatus:    user.AccountStatus,
		Reason:       fake.Sentence(10),
	}
}

// BuildFakeUserRegistrationInput builds a faked UserLoginInput.
func BuildFakeUserRegistrationInput() *types.UserRegistrationInput {
	return &types.UserRegistrationInput{
		Username:     fake.Username(),
		Password:     fake.Password(true, true, true, true, false, 32),
		EmailAddress: fake.Email(),
		BirthDay:     func(x uint8) *uint8 { return &x }(fake.Uint8()),
		BirthMonth:   func(x uint8) *uint8 { return &x }(fake.Uint8()),
	}
}

// BuildFakeUserLoginInputFromUser builds a faked UserLoginInput.
func BuildFakeUserLoginInputFromUser(user *types.User) *types.UserLoginInput {
	return &types.UserLoginInput{
		Username:  user.Username,
		Password:  fake.Password(true, true, true, true, false, 32),
		TOTPToken: fmt.Sprintf("0%s", fake.Zip()),
	}
}

// BuildFakePasswordUpdateInput builds a faked PasswordUpdateInput.
func BuildFakePasswordUpdateInput() *types.PasswordUpdateInput {
	return &types.PasswordUpdateInput{
		NewPassword:     fake.Password(true, true, true, true, false, 32),
		CurrentPassword: fake.Password(true, true, true, true, false, 32),
		TOTPToken:       fmt.Sprintf("0%s", fake.Zip()),
	}
}

// BuildFakeTOTPSecretRefreshInput builds a faked TOTPSecretRefreshInput.
func BuildFakeTOTPSecretRefreshInput() *types.TOTPSecretRefreshInput {
	return &types.TOTPSecretRefreshInput{
		CurrentPassword: fake.Password(true, true, true, true, false, 32),
		TOTPToken:       fmt.Sprintf("0%s", fake.Zip()),
	}
}

// BuildFakeUserPermissionsRequestInput builds a faked UserPermissionsRequestInput.
func BuildFakeUserPermissionsRequestInput() *types.UserPermissionsRequestInput {
	return &types.UserPermissionsRequestInput{
		Permissions: []string{
			fake.LoremIpsumSentence(exampleQuantity),
			fake.LoremIpsumSentence(exampleQuantity),
			fake.LoremIpsumSentence(exampleQuantity),
		},
	}
}

// BuildFakeTOTPSecretVerificationInput builds a faked TOTPSecretVerificationInput.
func BuildFakeTOTPSecretVerificationInput() *types.TOTPSecretVerificationInput {
	user := BuildFakeUser()

	token, err := totp.GenerateCode(user.TwoFactorSecret, time.Now().UTC())
	if err != nil {
		log.Panicf("error generating TOTP token for fakes user: %v", err)
	}

	return &types.TOTPSecretVerificationInput{
		UserID:    user.ID,
		TOTPToken: token,
	}
}

// BuildFakeTOTPSecretVerificationInputForUser builds a faked TOTPSecretVerificationInput for a given user.
func BuildFakeTOTPSecretVerificationInputForUser(user *types.User) *types.TOTPSecretVerificationInput {
	token, err := totp.GenerateCode(user.TwoFactorSecret, time.Now().UTC())
	if err != nil {
		log.Panicf("error generating TOTP token for fakes user: %v", err)
	}

	return &types.TOTPSecretVerificationInput{
		UserID:    user.ID,
		TOTPToken: token,
	}
}

// BuildFakePasswordResetToken builds a faked PasswordResetToken.
func BuildFakePasswordResetToken() *types.PasswordResetToken {
	fakeDate := fake.Date()

	return &types.PasswordResetToken{
		ID:            ksuid.New().String(),
		Token:         fake.UUID(),
		BelongsToUser: ksuid.New().String(),
		ExpiresAt:     uint64(uint32(fakeDate.Add(30 * time.Minute).Unix())),
		CreatedOn:     uint64(uint32(fakeDate.Unix())),
	}
}

// BuildFakePasswordResetTokenCreationRequestInput builds a faked PasswordResetTokenCreationRequestInput.
func BuildFakePasswordResetTokenCreationRequestInput() *types.PasswordResetTokenCreationRequestInput {
	return &types.PasswordResetTokenCreationRequestInput{EmailAddress: fake.Email()}
}

// BuildFakePasswordResetTokenRedemptionRequestInput builds a faked PasswordResetTokenRedemptionRequestInput.
func BuildFakePasswordResetTokenRedemptionRequestInput() *types.PasswordResetTokenRedemptionRequestInput {
	return &types.PasswordResetTokenRedemptionRequestInput{
		Token:       fake.LoremIpsumSentence(exampleQuantity),
		NewPassword: fake.Password(true, true, true, true, false, 32),
	}
}
