package dormancy

import "os"

// DevModeEnabled returns true if development mode for dormancy is enabled.
// Set CODER_DEV_DORMANCY=true to enable dormancy functionality without a license.
// This is intended for testing and evaluation purposes only.
func DevModeEnabled() bool {
	return os.Getenv("CODER_DEV_DORMANCY") == "true"
}
