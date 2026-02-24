package prebuilds

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/xerrors"

	"github.com/coder/coder/v2/coderd/database"
)

// DevModeEnabled returns true if development mode for prebuilds is enabled.
// Set CODER_DEV_PREBUILDS=true to enable prebuilds functionality without a license.
// This is intended for testing and evaluation purposes only.
func DevModeEnabled() bool {
	return os.Getenv("CODER_DEV_PREBUILDS") == "true"
}

// DevClaimer implements the Claimer interface for development/testing purposes.
// It provides the same functionality as the enterprise claimer without license checks.
type DevClaimer struct {
	store database.Store
}

// NewDevClaimer creates a new DevClaimer with the given database store.
func NewDevClaimer(store database.Store) *DevClaimer {
	return &DevClaimer{store: store}
}

func (c *DevClaimer) Claim(
	ctx context.Context,
	store database.Store,
	now time.Time,
	userID uuid.UUID,
	name string,
	presetID uuid.UUID,
	autostartSchedule sql.NullString,
	nextStartAt sql.NullTime,
	ttl sql.NullInt64,
) (*uuid.UUID, error) {
	result, err := store.ClaimPrebuiltWorkspace(ctx, database.ClaimPrebuiltWorkspaceParams{
		NewUserID:         userID,
		NewName:           name,
		Now:               now,
		PresetID:          presetID,
		AutostartSchedule: autostartSchedule,
		NextStartAt:       nextStartAt,
		WorkspaceTtl:      ttl,
	})
	if err != nil {
		switch {
		// No eligible prebuilds found.
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoClaimablePrebuiltWorkspaces
		default:
			return nil, xerrors.Errorf("claim prebuild for user %q: %w", userID.String(), err)
		}
	}

	return &result.ID, nil
}

var _ Claimer = &DevClaimer{}

// Note: The full DevStoreReconciler is implemented in dev_reconcile.go.
// Use NewDevStoreReconciler() to create a reconciler with full automatic pool management.
