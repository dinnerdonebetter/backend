package featureflags

import (
	"context"

	"github.com/prixfixeco/backend/pkg/types"
)

type (
	// FeatureFlagManager manages feature flags.
	FeatureFlagManager interface {
		Identify(ctx context.Context, user *types.User) error
		CanUseFeature(ctx context.Context, userID, feature string) (bool, error)
	}
)

func NewNoopFeatureFlagManager() FeatureFlagManager {
	return &NoopFeatureFlagManager{}
}

// NoopFeatureFlagManager is a no-op FeatureFlagManager.
type NoopFeatureFlagManager struct{}

// Identify implements the FeatureFlagManager interface.
func (m *NoopFeatureFlagManager) Identify(context.Context, *types.User) error {
	return nil
}

// CanUseFeature implements the FeatureFlagManager interface.
func (*NoopFeatureFlagManager) CanUseFeature(context.Context, string, string) (bool, error) {
	return false, nil
}
