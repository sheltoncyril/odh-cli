package guardrails

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/lburgazzoli/odh-cli/pkg/lint/check"
	"github.com/lburgazzoli/odh-cli/pkg/lint/check/result"
	"github.com/lburgazzoli/odh-cli/pkg/lint/checks/shared/validate"
	"github.com/lburgazzoli/odh-cli/pkg/util/inspect"
)

// deprecatedOtelExpressions lists JQ expressions for the deprecated otelExporter section in 3.x.
// The entire otelExporter struct is changing and needs migration.
//
//nolint:gochecknoglobals // Package-level constant for deprecated field expressions.
var deprecatedOtelExpressions = []string{
	".spec.otelExporter",
}

// hasDeprecatedOtelFields returns true if the object contains any deprecated otelExporter fields.
func hasDeprecatedOtelFields(obj *unstructured.Unstructured) (bool, error) {
	found, err := inspect.HasFields(obj.Object, deprecatedOtelExpressions...)
	if err != nil {
		//nolint:wrapcheck // WorkloadBuilder.Run wraps filter errors with resource-type context.
		return false, err
	}

	return len(found) > 0, nil
}

func newOtelMigrationCondition(
	_ context.Context,
	req *validate.WorkloadRequest[*unstructured.Unstructured],
) ([]result.Condition, error) {
	count := len(req.Items)

	if count == 0 {
		return []result.Condition{check.NewCondition(
			ConditionTypeOtelConfigCompatible,
			metav1.ConditionTrue,
			check.ReasonVersionCompatible,
			"No GuardrailsOrchestrators found using deprecated otelExporter fields",
		)}, nil
	}

	return []result.Condition{check.NewCondition(
		ConditionTypeOtelConfigCompatible,
		metav1.ConditionFalse,
		check.ReasonConfigurationInvalid,
		"Found %d GuardrailsOrchestrator(s) using deprecated otelExporter fields - migrate to new format before upgrading",
		count,
		check.WithImpact(result.ImpactAdvisory),
	)}, nil
}
