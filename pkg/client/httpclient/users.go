package httpclient

import (
	"context"
	"fmt"
	"strings"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

// GetUser retrieves a user.
func (c *Client) GetUser(ctx context.Context, userID string) (*types.User, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}

	req, err := c.requestBuilder.BuildGetUserRequest(ctx, userID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building get user request")
	}

	var user *types.User
	if err = c.fetchAndUnmarshal(ctx, req, &user); err != nil {
		return nil, observability.PrepareError(err, span, "fetching user")
	}

	return user, nil
}

// GetUsers retrieves a list of users.
func (c *Client) GetUsers(ctx context.Context, filter *types.QueryFilter) (*types.UserList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetUsersRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building users list request")
	}

	var users *types.UserList
	if err = c.fetchAndUnmarshal(ctx, req, &users); err != nil {
		return nil, observability.PrepareError(err, span, "retrieving users")
	}

	return users, nil
}

// SearchForUsersByUsername searches for a user from a list of users by their username.
func (c *Client) SearchForUsersByUsername(ctx context.Context, username string) ([]*types.User, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if username == "" {
		return nil, ErrEmptyUsernameProvided
	}

	req, err := c.requestBuilder.BuildSearchForUsersByUsernameRequest(ctx, username)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building username search request")
	}

	var users []*types.User
	if err = c.fetchAndUnmarshal(ctx, req, &users); err != nil {
		return nil, observability.PrepareError(err, span, "searching for users")
	}

	return users, nil
}

// CreateUser creates a new user.
func (c *Client) CreateUser(ctx context.Context, input *types.UserRegistrationInput) (*types.UserCreationResponse, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	// deliberately not validating here
	// maybe I should make a client-side validate method vs a server-side?

	req, err := c.requestBuilder.BuildCreateUserRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building create user request")
	}

	var user *types.UserCreationResponse
	if err = c.fetchAndUnmarshalWithoutAuthentication(ctx, req, &user); err != nil {
		return nil, observability.PrepareError(err, span, "creating user")
	}

	return user, nil
}

// ArchiveUser archives a user.
func (c *Client) ArchiveUser(ctx context.Context, userID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return ErrInvalidIDProvided
	}

	req, err := c.requestBuilder.BuildArchiveUserRequest(ctx, userID)
	if err != nil {
		return observability.PrepareError(err, span, "building archive user request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareError(err, span, "archiving user")
	}

	return nil
}

const (
	pngExtension  = "pngExtension"
	jpegExtension = "jpegExtension"
	gifExtension  = "gifExtension"
)

// UploadNewAvatar uploads a new avatar.
func (c *Client) UploadNewAvatar(ctx context.Context, avatar []byte, extension string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if len(avatar) == 0 {
		return fmt.Errorf("%w: %d", ErrInvalidAvatarSize, len(avatar))
	}

	ex := strings.ToLower(strings.TrimSpace(extension))
	if ex != jpegExtension && ex != pngExtension && ex != gifExtension {
		err := fmt.Errorf("%s: %w", extension, ErrInvalidImageExtension)
		return observability.PrepareError(err, span, "uploading avatar")
	}

	req, err := c.requestBuilder.BuildAvatarUploadRequest(ctx, avatar, extension)
	if err != nil {
		return observability.PrepareError(err, span, "building avatar upload request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareError(err, span, "uploading avatar")
	}

	return nil
}

// CheckUserPermissions checks if a user has certain permissions.
func (c *Client) CheckUserPermissions(ctx context.Context, permissions ...string) (*types.UserPermissionsResponse, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if permissions == nil {
		return nil, ErrNilInputProvided
	}

	req, err := c.requestBuilder.BuildCheckUserPermissionsRequests(ctx, permissions...)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building permission check request")
	}

	var res *types.UserPermissionsResponse
	if err = c.fetchAndUnmarshal(ctx, req, &res); err != nil {
		return nil, observability.PrepareError(err, span, "checking permission")
	}

	return res, nil
}
