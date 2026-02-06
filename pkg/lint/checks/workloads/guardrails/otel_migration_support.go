package guardrails

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/lburgazzoli/odh-cli/pkg/lint/check"
	"github.com/lburgazzoli/odh-cli/pkg/lint/check/result"
)

// deprecatedOtelExpressions lists JQ expressions for the deprecated otelExporter section in 3.x.
// The entire otelExporter struct is changing and needs migration.
//
//nolint:gochecknoglobals // Package-level constant for deprecated field expressions.
var deprecatedOtelExpressions = []string{
	".spec.otelExporter",
}

func newOtelMigrationCondition(count int) result.Condition {
	if count == 0 {
		return check.NewCondition(
			ConditionTypeOtelConfigCompatible,
			metav1.ConditionTrue,
			check.ReasonVersionCompatible,
			"No GuardrailsOrchestrators found using deprecated otelExporter fields",
		)
	}

	return check.NewCondition(
		ConditionTypeOtelConfigCompatible,
		metav1.ConditionFalse,
		check.ReasonConfigurationInvalid,
		"Found %d GuardrailsOrchestrator(s) using deprecated otelExporter fields - migrate to new format before upgrading",
		count,
		check.WithImpact(result.ImpactAdvisory),
	)
}
