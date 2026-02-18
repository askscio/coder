package dormancy_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/coder/coder/v2/coderd/dormancy"
)

func TestDevModeEnabled(t *testing.T) {
	t.Run("EnabledWhenEnvVarSet", func(t *testing.T) {
		t.Setenv("CODER_DEV_DORMANCY", "true")
		require.True(t, dormancy.DevModeEnabled())
	})

	t.Run("DisabledWhenEnvVarNotSet", func(t *testing.T) {
		// Ensure the env var is not set.
		_ = os.Unsetenv("CODER_DEV_DORMANCY")
		require.False(t, dormancy.DevModeEnabled())
	})

	t.Run("DisabledWhenEnvVarSetToFalse", func(t *testing.T) {
		t.Setenv("CODER_DEV_DORMANCY", "false")
		require.False(t, dormancy.DevModeEnabled())
	})
}
