package action

import (
	"context"

	"github.com/blang/semver/v4"

	"github.com/lburgazzoli/odh-cli/pkg/migrate/action/result"
	"github.com/lburgazzoli/odh-cli/pkg/util/client"
	"github.com/lburgazzoli/odh-cli/pkg/util/iostreams"
)

type ActionGroup string

const (
	GroupMigration  ActionGroup = "migration"
	GroupBackup     ActionGroup = "backup"
	GroupValidation ActionGroup = "validation"
)

type Action interface {
	ID() string
	Name() string
	Description() string
	Group() ActionGroup

	// CanApply returns whether this action should run for the given target context.
	// Actions can use target.CurrentVersion, target.TargetVersion, or target.Client for filtering.
	CanApply(target *ActionTarget) bool
	Validate(ctx context.Context, target *ActionTarget) (*result.ActionResult, error)
	Execute(ctx context.Context, target *ActionTarget) (*result.ActionResult, error)
}

// ActionTarget holds all context needed for executing migration actions.
type ActionTarget struct {
	Client         *client.Client
	CurrentVersion *semver.Version // Version being migrated FROM
	TargetVersion  *semver.Version // Version being migrated TO
	DryRun         bool
	SkipConfirm    bool
	Recorder       StepRecorder
	IO             iostreams.Interface
}
