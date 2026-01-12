package result

import (
	"errors"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Validation error messages.
	errMsgGroupEmpty               = "group must not be empty"
	errMsgKindEmpty                = "kind must not be empty"
	errMsgNameEmpty                = "name must not be empty"
	errMsgConditionsEmpty          = "status.conditions must contain at least one condition"
	errMsgConditionTypeEmpty       = "condition with empty type found"
	errMsgConditionReasonEmpty     = "condition %q has empty reason"
	errMsgConditionInvalidStatus   = "condition %q has invalid status (must be True, False, or Unknown)"
	errMsgAnnotationInvalidFormat  = "annotation key %q must be in domain/key format (e.g., openshiftai.io/version)"
	errMsgConditionInvalidSeverity = "condition %q has invalid severity (must be critical, warning, or info)"
)

// Severity represents the impact level of a diagnostic condition.
type Severity string

// Severity levels for diagnostic conditions.
const (
	SeverityCritical Severity = "critical"
	SeverityWarning  Severity = "warning"
	SeverityInfo     Severity = "" // Empty string to omit from output (omitempty)
)

// Condition represents a diagnostic condition with severity level.
// It embeds metav1.Condition and adds a Severity field to indicate
// the impact level of the condition result.
type Condition struct {
	metav1.Condition `json:",inline" yaml:",inline"`

	// Severity indicates the impact level: "critical", "warning", or "info".
	// If empty, severity is derived from Status (False=critical, Unknown=warning, True=info).
	Severity Severity `json:"severity,omitempty" yaml:"severity,omitempty"`
}

// DiagnosticSpec describes what the check validates.
type DiagnosticSpec struct {
	// Description provides a detailed explanation of the check purpose and significance
	Description string `json:"description" yaml:"description"`
}

// DiagnosticStatus contains the condition-based validation results.
type DiagnosticStatus struct {
	// Conditions is an array of validation conditions ordered by execution sequence
	Conditions []Condition `json:"conditions" yaml:"conditions"`
}

// DiagnosticResult represents a diagnostic check result with flattened metadata fields.
type DiagnosticResult struct {
	// Group is the diagnostic target category (e.g., "components", "services", "workloads")
	Group string `json:"group" yaml:"group"`

	// Kind is the specific target being checked (e.g., "kserve", "auth", "cert-manager")
	Kind string `json:"kind" yaml:"kind"`

	// Name is the check identifier (e.g., "version-compatibility", "configuration-valid")
	Name string `json:"name" yaml:"name"`

	// Annotations contains optional key-value metadata with domain-qualified keys
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`

	// Spec describes what the check validates
	Spec DiagnosticSpec `json:"spec" yaml:"spec"`

	// Status contains the condition-based validation results
	Status DiagnosticStatus `json:"status" yaml:"status"`

	// ImpactedObjects contains references to resources impacted by this diagnostic.
	// Uses PartialObjectMetadata to store minimal object info with optional annotations
	// for additional context (e.g., deployment mode, configuration details).
	ImpactedObjects []metav1.PartialObjectMetadata `json:"impactedObjects,omitempty" yaml:"impactedObjects,omitempty"`
}

// isValidAnnotationKey validates that an annotation key follows the domain/key format.
// Valid examples: openshiftai.io/version, example.com/name
// Invalid examples: version, /name, example.com/.
func isValidAnnotationKey(key string) bool {
	// Must contain exactly one '/' separating domain and key
	parts := strings.Split(key, "/")
	if len(parts) != 2 {
		return false
	}

	domain, name := parts[0], parts[1]

	// Domain and name must both be non-empty
	if domain == "" || name == "" {
		return false
	}

	// Domain must contain at least one dot (e.g., example.com, openshiftai.io)
	if !strings.Contains(domain, ".") {
		return false
	}

	return true
}

// validateCondition validates a single condition.
func validateCondition(condition *Condition) error {
	if condition.Type == "" {
		return errors.New(errMsgConditionTypeEmpty)
	}
	if condition.Status != metav1.ConditionTrue &&
		condition.Status != metav1.ConditionFalse &&
		condition.Status != metav1.ConditionUnknown {
		return fmt.Errorf(errMsgConditionInvalidStatus, condition.Type)
	}
	if condition.Reason == "" {
		return fmt.Errorf(errMsgConditionReasonEmpty, condition.Type)
	}
	// Validate severity if provided (empty string is valid for SeverityInfo)
	if condition.Severity != SeverityCritical &&
		condition.Severity != SeverityWarning &&
		condition.Severity != SeverityInfo {
		return fmt.Errorf(errMsgConditionInvalidSeverity, condition.Type)
	}

	return nil
}

// Validate checks if the diagnostic result is valid.
func (r *DiagnosticResult) Validate() error {
	// Validate flattened metadata fields
	if r.Group == "" {
		return errors.New(errMsgGroupEmpty)
	}
	if r.Kind == "" {
		return errors.New(errMsgKindEmpty)
	}
	if r.Name == "" {
		return errors.New(errMsgNameEmpty)
	}

	// Validate annotation keys follow domain/key format
	for key := range r.Annotations {
		if !isValidAnnotationKey(key) {
			return fmt.Errorf(errMsgAnnotationInvalidFormat, key)
		}
	}

	// Validate conditions array
	if len(r.Status.Conditions) == 0 {
		return errors.New(errMsgConditionsEmpty)
	}

	// Validate each condition
	for i := range r.Status.Conditions {
		if err := validateCondition(&r.Status.Conditions[i]); err != nil {
			return err
		}
	}

	return nil
}

// New creates a new diagnostic result.
func New(
	group string,
	kind string,
	name string,
	description string,
) *DiagnosticResult {
	return &DiagnosticResult{
		Group:       group,
		Kind:        kind,
		Name:        name,
		Annotations: make(map[string]string),
		Spec: DiagnosticSpec{
			Description: description,
		},
		Status: DiagnosticStatus{
			Conditions: []Condition{},
		},
	}
}

// IsFailing returns true if any condition has status False or Unknown.
func (r *DiagnosticResult) IsFailing() bool {
	for _, cond := range r.Status.Conditions {
		if cond.Status == metav1.ConditionFalse || cond.Status == metav1.ConditionUnknown {
			return true
		}
	}

	return false
}

// GetMessage returns a summary message from all conditions.
func (r *DiagnosticResult) GetMessage() string {
	if len(r.Status.Conditions) == 0 {
		return ""
	}
	// Return the first condition's message as the primary message
	return r.Status.Conditions[0].Message
}

// GetSeverity returns the severity level based on condition severities.
// If a condition has an explicit Severity field set, that value is used.
// Otherwise, severity is derived from the condition status:
// Critical: ConditionFalse (indicates a failure)
// Warning: ConditionUnknown (indicates an error or inability to determine)
// Info: ConditionTrue (indicates success/informational)
// Returns the highest severity found across all conditions (critical > warning > info).
// Returns nil if there are no conditions.
func (r *DiagnosticResult) GetSeverity() *string {
	if len(r.Status.Conditions) == 0 {
		return nil
	}

	var hasCritical bool
	var hasWarning bool

	for _, cond := range r.Status.Conditions {
		severity := cond.Severity
		// If severity not explicitly set, derive from status
		if severity == "" {
			switch cond.Status {
			case metav1.ConditionFalse:
				severity = SeverityCritical
			case metav1.ConditionUnknown:
				severity = SeverityWarning
			case metav1.ConditionTrue:
				severity = SeverityInfo
			default:
				continue
			}
		}

		// Track highest severity (critical > warning > info)
		switch severity {
		case SeverityCritical:
			hasCritical = true
		case SeverityWarning:
			hasWarning = true
		case SeverityInfo:
			// Info severity doesn't affect flags (lowest priority)
		}
	}

	var result Severity
	if hasCritical {
		result = SeverityCritical
	} else if hasWarning {
		result = SeverityWarning
	} else {
		result = SeverityInfo
	}

	resultStr := string(result)

	return &resultStr
}

// GetRemediation returns remediation guidance.
// Currently returns empty string as remediation is not part of the CR pattern.
// Remediation can be inferred from condition reasons and messages.
func (r *DiagnosticResult) GetRemediation() string {
	return ""
}

// GetStatusString returns a string representation of the overall status.
// Pass: All conditions are True
// Fail: Any condition is False
// Error: Any condition is Unknown.
func (r *DiagnosticResult) GetStatusString() string {
	if len(r.Status.Conditions) == 0 {
		return "Unknown"
	}

	for _, cond := range r.Status.Conditions {
		if cond.Status == metav1.ConditionFalse {
			return "Fail"
		}
		if cond.Status == metav1.ConditionUnknown {
			return "Error"
		}
	}

	return "Pass"
}

// DiagnosticResultList represents a list of diagnostic results.
type DiagnosticResultList struct {
	ClusterVersion *string             `json:"clusterVersion,omitempty" yaml:"clusterVersion,omitempty"`
	TargetVersion  *string             `json:"targetVersion,omitempty"  yaml:"targetVersion,omitempty"`
	Results        []*DiagnosticResult `json:"results"                  yaml:"results"`
}

// NewDiagnosticResultList creates a new list.
func NewDiagnosticResultList(clusterVersion *string, targetVersion *string) *DiagnosticResultList {
	return &DiagnosticResultList{
		ClusterVersion: clusterVersion,
		TargetVersion:  targetVersion,
		Results:        make([]*DiagnosticResult, 0),
	}
}
