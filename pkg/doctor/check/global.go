package check

import (
	"fmt"
	"os"
)

// globalRegistry is the global check registry.
// Global variable is intentional for check auto-registration pattern via init().
//
//nolint:gochecknoglobals // Required for auto-registration pattern
var globalRegistry = NewRegistry()

// GetGlobalRegistry provides access to the shared check registry used for auto-registration.
func GetGlobalRegistry() *CheckRegistry {
	return globalRegistry
}

// MustRegisterCheck registers a check in the global registry
// Panics if registration fails.
func MustRegisterCheck(check Check) {
	if err := globalRegistry.Register(check); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to register check %s: %v\n", check.ID(), err)
		panic(err)
	}
}
