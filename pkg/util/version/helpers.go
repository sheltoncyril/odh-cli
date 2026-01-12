package version

import "github.com/blang/semver/v4"

// IsUpgradeFrom2xTo3x checks if the versions represent an upgrade from 2.x to 3.x specifically.
// Future major versions (4.x+) may have different compatibility requirements.
// Returns false if either version is nil.
func IsUpgradeFrom2xTo3x(from *semver.Version, to *semver.Version) bool {
	if from == nil || to == nil {
		return false
	}

	return from.Major == 2 && to.Major == 3
}
