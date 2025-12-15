package migrate

import (
	"context"
	"errors"
	"fmt"

	"github.com/blang/semver/v4"
	"github.com/spf13/pflag"

	"k8s.io/cli-runtime/pkg/genericiooptions"

	"github.com/lburgazzoli/odh-cli/pkg/cmd"
	"github.com/lburgazzoli/odh-cli/pkg/lint/version"
	"github.com/lburgazzoli/odh-cli/pkg/migrate/action"
)

var _ cmd.Command = (*RunCommand)(nil)

type RunCommand struct {
	*SharedOptions

	DryRun        bool
	Prepare       bool
	Yes           bool
	BackupPath    string
	MigrationIDs  []string
	TargetVersion string

	parsedTargetVersion *semver.Version
}

func NewRunCommand(streams genericiooptions.IOStreams) *RunCommand {
	shared := NewSharedOptions(streams)

	return &RunCommand{
		SharedOptions: shared,
		BackupPath:    "./backups",
	}
}

func (c *RunCommand) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVarP(&c.Verbose, "verbose", "v", false,
		"Show detailed progress")
	fs.DurationVar(&c.Timeout, "timeout", c.Timeout,
		"Operation timeout (e.g., 10m, 30m)")

	fs.BoolVar(&c.DryRun, "dry-run", false,
		"Show what would be done without making changes")
	fs.BoolVar(&c.Prepare, "prepare", false,
		"Run pre-flight checks and backup resources (does not execute migration)")
	fs.BoolVarP(&c.Yes, "yes", "y", false,
		"Skip confirmation prompts")
	fs.StringVar(&c.BackupPath, "backup-path", c.BackupPath,
		"Path to store backup files (used with --prepare)")
	fs.StringArrayVarP(&c.MigrationIDs, "migration", "m", []string{},
		"Migration ID to execute (can be specified multiple times)")
	fs.StringVar(&c.TargetVersion, "target-version", "",
		"Target version for migration (required)")
}

func (c *RunCommand) Complete() error {
	if err := c.SharedOptions.Complete(); err != nil {
		return fmt.Errorf(msgCompletingOptions, err)
	}

	// Always enable verbose for migrate run (both dry-run and actual execution)
	c.Verbose = true

	if c.TargetVersion != "" {
		targetVer, err := semver.Parse(c.TargetVersion)
		if err != nil {
			return fmt.Errorf(msgInvalidTargetVersion, c.TargetVersion, err)
		}
		c.parsedTargetVersion = &targetVer
	}

	return nil
}

func (c *RunCommand) Validate() error {
	if err := c.SharedOptions.Validate(); err != nil {
		return fmt.Errorf(msgValidatingOptions, err)
	}

	if len(c.MigrationIDs) == 0 {
		return errors.New(msgMigrationRequired)
	}

	if c.TargetVersion == "" {
		return errors.New(msgTargetVersionRequired)
	}

	return nil
}

func (c *RunCommand) Run(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	currentVersion, err := version.Detect(ctx, c.Client)
	if err != nil {
		return fmt.Errorf(msgDetectingVersion, err)
	}

	targetVersionInfo := &version.ClusterVersion{
		Version:    c.TargetVersion,
		Source:     version.SourceManual,
		Confidence: version.ConfidenceHigh,
	}

	registry := action.GetGlobalRegistry()

	if c.Prepare {
		return c.runPrepareMode(ctx, currentVersion, targetVersionInfo, registry)
	}

	return c.runMigrationMode(ctx, currentVersion, targetVersionInfo, registry)
}

func (c *RunCommand) runPrepareMode(
	ctx context.Context,
	currentVersion *version.ClusterVersion,
	targetVersionInfo *version.ClusterVersion,
	registry *action.ActionRegistry,
) error {
	for _, migrationID := range c.MigrationIDs {
		c.IO.Errorf("Running pre-flight checks for migration: %s\n", migrationID)

		selectedAction, ok := registry.Get(migrationID)
		if !ok {
			return fmt.Errorf(msgMigrationNotFound, migrationID)
		}

		// Use verbose recorder for real-time streaming output
		recorder := action.NewVerboseRootRecorder(c.IO)
		c.IO.Errorf("\n%s:\n", migrationID)

		target := &action.ActionTarget{
			Client:         c.Client,
			CurrentVersion: currentVersion,
			TargetVersion:  targetVersionInfo,
			DryRun:         c.DryRun,
			BackupPath:     c.BackupPath,
			SkipConfirm:    c.Yes,
			Recorder:       recorder,
			IO:             c.IO,
		}

		_, err := selectedAction.Validate(ctx, target)
		if err != nil {
			return fmt.Errorf(msgPreFlightFailed, err)
		}

		// Output has already been streamed during validation
		c.IO.Fprintln()
	}

	c.IO.Errorf("Preparation complete. Run without --prepare to execute migration.")

	return nil
}

func (c *RunCommand) runMigrationMode(
	ctx context.Context,
	currentVersion *version.ClusterVersion,
	targetVersionInfo *version.ClusterVersion,
	registry *action.ActionRegistry,
) error {
	c.IO.Errorf("Current OpenShift AI version: %s", currentVersion.Version)
	c.IO.Errorf("Target OpenShift AI version: %s\n", targetVersionInfo.Version)

	for idx, migrationID := range c.MigrationIDs {
		if len(c.MigrationIDs) > 1 {
			c.IO.Errorf("\n=== Migration %d/%d: %s ===\n", idx+1, len(c.MigrationIDs), migrationID)
		}

		selectedAction, ok := registry.Get(migrationID)
		if !ok {
			return fmt.Errorf(msgMigrationNotFound, migrationID)
		}

		// Use verbose recorder for real-time streaming output
		recorder := action.NewVerboseRootRecorder(c.IO)
		c.IO.Errorf("\n%s:\n", migrationID)

		target := &action.ActionTarget{
			Client:         c.Client,
			CurrentVersion: currentVersion,
			TargetVersion:  targetVersionInfo,
			DryRun:         c.DryRun,
			BackupPath:     c.BackupPath,
			SkipConfirm:    c.Yes,
			Recorder:       recorder,
			IO:             c.IO,
		}

		if c.DryRun {
			c.IO.Errorf("DRY RUN MODE: No changes will be made to the cluster\n")
		} else if c.Yes {
			c.IO.Errorf("Running migration: %s (confirmations skipped)\n", migrationID)
		} else {
			c.IO.Errorf("Preparing migration: %s\n", migrationID)
		}

		actionResult, err := selectedAction.Execute(ctx, target)
		if err != nil {
			return fmt.Errorf(msgMigrationFailed, err)
		}

		// Output has already been streamed during execution, no need to render again
		c.IO.Fprintln()
		if !actionResult.Status.Completed {
			c.IO.Errorf("Migration %s incomplete - please review the output above", migrationID)

			return fmt.Errorf(msgMigrationHalted, migrationID)
		}
		c.IO.Errorf("Migration %s completed successfully!", migrationID)
	}

	c.IO.Fprintln()
	c.IO.Errorf("All migrations completed successfully!")

	return nil
}
