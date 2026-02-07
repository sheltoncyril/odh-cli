package guardrails

import (
	"context"

	"github.com/lburgazzoli/odh-cli/pkg/lint/check"
	"github.com/lburgazzoli/odh-cli/pkg/lint/check/result"
	"github.com/lburgazzoli/odh-cli/pkg/lint/checks/shared/base"
	"github.com/lburgazzoli/odh-cli/pkg/lint/checks/shared/validate"
	"github.com/lburgazzoli/odh-cli/pkg/resources"
	"github.com/lburgazzoli/odh-cli/pkg/util/version"
)

const (
	kind = "guardrails"

	ConditionTypeOtelConfigCompatible = "OtelConfigCompatible"

	// minTargetMajorVersion is the minimum major version for this check to apply.
	minTargetMajorVersion = 3
)

// OtelMigrationCheck detects GuardrailsOrchestrator CRs using deprecated otelExporter configuration fields.
type OtelMigrationCheck struct {
	base.BaseCheck
}

func NewOtelMigrationCheck() *OtelMigrationCheck {
	return &OtelMigrationCheck{
		BaseCheck: base.BaseCheck{
			CheckGroup:       check.GroupWorkload,
			Kind:             kind,
			Type:             check.CheckTypeConfigMigration,
			CheckID:          "workloads.guardrails.otel-config-migration",
			CheckName:        "Workloads :: Guardrails :: OTEL Config Migration (3.x)",
			CheckDescription: "Detects GuardrailsOrchestrator CRs using deprecated otelExporter configuration fields that need migration",
		},
	}
}

// CanApply returns whether this check should run for the given target.
// Only applies when upgrading to 3.x or later.
func (c *OtelMigrationCheck) CanApply(_ context.Context, target check.Target) bool {
	return version.IsVersionAtLeast(target.TargetVersion, minTargetMajorVersion, 0)
}

// Validate executes the check against the provided target.
func (c *OtelMigrationCheck) Validate(
	ctx context.Context,
	target check.Target,
) (*result.DiagnosticResult, error) {
	return validate.Workloads(c, target, resources.GuardrailsOrchestrator).
		Filter(hasDeprecatedOtelFields).
		Complete(ctx, newOtelMigrationCondition)
}
