package migrate

const (
	msgMigrationNotFound      = "migration %q not found"
	msgInvalidTargetVersion   = "invalid target version %q: %w"
	msgTargetVersionRequired  = "--target-version flag is required"
	msgMigrationRequired      = "--migration flag is required"
	msgNoApplicableMigrations = "no migrations applicable for version %s"
	msgDetectingVersion       = "detecting cluster version: %w"
	msgMigrationHalted        = "migration halted: %s"
	msgMigrationFailed        = "migration failed: %w"
	msgPreFlightFailed        = "pre-flight validation failed: %w"
	msgCompletingOptions      = "completing shared options: %w"
	msgValidatingOptions      = "validating shared options: %w"
	msgMutuallyExclusive      = "--all and --target-version are mutually exclusive"
)
