package connectionlog

import (
	"context"
	"os"
	"sync"

	"github.com/coder/coder/v2/coderd/database"
)

// DevModeEnabled returns true if development mode for connection logs is enabled.
// Set CODER_DEV_CONNECTION_LOGS=true to enable connection logging without a license.
// This is intended for testing and evaluation purposes only.
func DevModeEnabled() bool {
	return os.Getenv("CODER_DEV_CONNECTION_LOGS") == "true"
}

// DevConnectionLogger implements the ConnectionLogger interface for development/testing.
// It stores connection logs in memory and does not require a license.
type DevConnectionLogger struct {
	*FakeConnectionLogger
}

// NewDevConnectionLogger creates a new DevConnectionLogger for development mode.
func NewDevConnectionLogger() *DevConnectionLogger {
	return &DevConnectionLogger{
		FakeConnectionLogger: NewFake(),
	}
}

func (l *DevConnectionLogger) Upsert(ctx context.Context, clog database.UpsertConnectionLogParams) error {
	return l.FakeConnectionLogger.Upsert(ctx, clog)
}

var _ ConnectionLogger = (*DevConnectionLogger)(nil)
