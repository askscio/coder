package dormancy

import (
	"context"
	"time"

	"cdr.dev/slog/v3"
	"github.com/coder/coder/v2/coderd/audit"
	"github.com/coder/coder/v2/coderd/database"
	agpldormancy "github.com/coder/coder/v2/coderd/dormancy"
	"github.com/coder/quartz"
)

const (
	// JobInterval is the time interval between consecutive job runs.
	JobInterval = agpldormancy.JobInterval
	// AccountDormancyPeriod defines how long user accounts can be inactive
	// before being marked as dormant.
	AccountDormancyPeriod = agpldormancy.AccountDormancyPeriod
)

// CheckInactiveUsers updates status of inactive users from active to dormant
// using default parameters.
func CheckInactiveUsers(ctx context.Context, logger slog.Logger, clk quartz.Clock, db database.Store, auditor audit.Auditor) func() {
	return agpldormancy.CheckInactiveUsers(ctx, logger, clk, db, auditor)
}

// CheckInactiveUsersWithOptions updates status of inactive users from active
// to dormant using provided parameters.
func CheckInactiveUsersWithOptions(ctx context.Context, logger slog.Logger, clk quartz.Clock, db database.Store, auditor audit.Auditor, checkInterval, dormancyPeriod time.Duration) func() {
	return agpldormancy.CheckInactiveUsersWithOptions(ctx, logger, clk, db, auditor, checkInterval, dormancyPeriod)
}
