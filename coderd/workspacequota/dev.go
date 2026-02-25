package workspacequota

import (
	"context"
	"database/sql"
	"errors"
	"os"

	"github.com/google/uuid"

	"cdr.dev/slog/v3"
	"github.com/coder/coder/v2/coderd/database"
	"github.com/coder/coder/v2/provisionerd/proto"
)

// DevModeEnabled returns true if development mode for quotas is enabled.
// Set CODER_DEV_QUOTAS=true to enable quota functionality without a license.
// This is intended for testing and evaluation purposes only.
func DevModeEnabled() bool {
	return os.Getenv("CODER_DEV_QUOTAS") == "true"
}

// DevCommitter implements the proto.QuotaCommitter interface for development/testing.
// It provides the same quota enforcement functionality as the enterprise committer
// without requiring a license.
type DevCommitter struct {
	Log      slog.Logger
	Database database.Store
}

// NewDevCommitter creates a new DevCommitter for development mode.
func NewDevCommitter(log slog.Logger, db database.Store) *DevCommitter {
	return &DevCommitter{
		Log:      log,
		Database: db,
	}
}

func (c *DevCommitter) CommitQuota(
	ctx context.Context, request *proto.CommitQuotaRequest,
) (*proto.CommitQuotaResponse, error) {
	jobID, err := uuid.Parse(request.JobId)
	if err != nil {
		return nil, err
	}

	nextBuild, err := c.Database.GetWorkspaceBuildByJobID(ctx, jobID)
	if err != nil {
		return nil, err
	}

	workspace, err := c.Database.GetWorkspaceByID(ctx, nextBuild.WorkspaceID)
	if err != nil {
		return nil, err
	}

	var (
		consumed int64
		budget   int64
		permit   bool
	)
	err = c.Database.InTx(func(s database.Store) error {
		var err error
		consumed, err = s.GetQuotaConsumedForUser(ctx, database.GetQuotaConsumedForUserParams{
			OwnerID:        workspace.OwnerID,
			OrganizationID: workspace.OrganizationID,
		})
		if err != nil {
			return err
		}

		budget, err = s.GetQuotaAllowanceForUser(ctx, database.GetQuotaAllowanceForUserParams{
			UserID:         workspace.OwnerID,
			OrganizationID: workspace.OrganizationID,
		})
		if err != nil {
			return err
		}

		// If the new build will reduce overall quota consumption, then we
		// allow it even if the user is over quota.
		netIncrease := true
		prevBuild, err := s.GetWorkspaceBuildByWorkspaceIDAndBuildNumber(ctx, database.GetWorkspaceBuildByWorkspaceIDAndBuildNumberParams{
			WorkspaceID: workspace.ID,
			BuildNumber: nextBuild.BuildNumber - 1,
		})
		if err == nil {
			netIncrease = request.DailyCost >= prevBuild.DailyCost
			c.Log.Debug(
				ctx, "previous build cost",
				slog.F("prev_cost", prevBuild.DailyCost),
				slog.F("next_cost", request.DailyCost),
				slog.F("net_increase", netIncrease),
			)
		} else if !errors.Is(err, sql.ErrNoRows) {
			return err
		}

		newConsumed := int64(request.DailyCost) + consumed
		if newConsumed > budget && netIncrease {
			c.Log.Debug(
				ctx, "over quota, rejecting",
				slog.F("prev_consumed", consumed),
				slog.F("next_consumed", newConsumed),
				slog.F("budget", budget),
			)
			permit = false
			return nil
		}

		permit = true
		return nil
	}, nil)
	if err != nil {
		return nil, err
	}

	return &proto.CommitQuotaResponse{
		Ok:                permit,
		CreditsConsumed:   consumed,
		Budget:            budget,
		DailyCostIncrease: request.DailyCost,
	}, nil
}

var _ proto.QuotaCommitter = (*DevCommitter)(nil)
