package audit

import (
	"context"
	"os"

	"github.com/coder/coder/v2/coderd/database"
)

// DevModeEnabled returns true if development mode for audit logs is enabled.
// Set CODER_DEV_AUDIT_LOGS=true to enable audit functionality without a license.
// This is intended for testing and evaluation purposes only.
func DevModeEnabled() bool {
	return os.Getenv("CODER_DEV_AUDIT_LOGS") == "true"
}

// DevAuditor implements the Auditor interface for development/testing purposes.
// It provides the same functionality as MockAuditor, storing audit logs in memory.
type DevAuditor struct {
	*MockAuditor
}

// NewDevAuditor creates a new DevAuditor for development mode.
// This auditor stores logs in memory and does not require a license.
func NewDevAuditor() *DevAuditor {
	return &DevAuditor{
		MockAuditor: NewMock(),
	}
}

func (a *DevAuditor) Export(ctx context.Context, alog database.AuditLog) error {
	return a.MockAuditor.Export(ctx, alog)
}

var _ Auditor = (*DevAuditor)(nil)
