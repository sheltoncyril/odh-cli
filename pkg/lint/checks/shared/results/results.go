package results

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/lburgazzoli/odh-cli/pkg/lint/check/result"
	"github.com/lburgazzoli/odh-cli/pkg/resources"
)

// SetCondition updates or adds a condition to the diagnostic result.
// If a condition with the same type already exists, it updates it.
// If no condition with that type exists, it adds a new one.
// This allows checks to potentially have multiple conditions in the future.
func SetCondition(dr *result.DiagnosticResult, condition result.Condition) {
	// Find and update existing condition of this type
	for i := range dr.Status.Conditions {
		if dr.Status.Conditions[i].Type == condition.Type {
			dr.Status.Conditions[i] = condition

			return
		}
	}

	// No existing condition found, append new one
	dr.Status.Conditions = append(dr.Status.Conditions, condition)
}

// PopulateImpactedObjects converts a list of NamespacedNames into PartialObjectMetadata
// and sets them as impacted objects on the diagnostic result.
func PopulateImpactedObjects(
	dr *result.DiagnosticResult,
	resourceType resources.ResourceType,
	names []types.NamespacedName,
) {
	dr.ImpactedObjects = make([]metav1.PartialObjectMetadata, 0, len(names))

	for _, n := range names {
		dr.ImpactedObjects = append(dr.ImpactedObjects, metav1.PartialObjectMetadata{
			TypeMeta: resourceType.TypeMeta(),
			ObjectMeta: metav1.ObjectMeta{
				Namespace: n.Namespace,
				Name:      n.Name,
			},
		})
	}
}
