package check

import "github.com/lburgazzoli/odh-cli/pkg/util/version"

// IsUpgradeFrom2xTo3x checks if target represents an upgrade from 2.x to 3.x or later.
// This is the most common version check pattern across checks in the codebase.
// Returns false if target is nil, or if either version is nil.
func IsUpgradeFrom2xTo3x(target *CheckTarget) bool {
	if target == nil {
		return false
	}

	return version.IsUpgradeFrom2xTo3x(target.CurrentVersion, target.Version)
}
