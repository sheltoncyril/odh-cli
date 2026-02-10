package results_test

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/lburgazzoli/odh-cli/pkg/lint/check"
	"github.com/lburgazzoli/odh-cli/pkg/lint/check/result"
	"github.com/lburgazzoli/odh-cli/pkg/lint/checks/shared/results"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

func TestSetCondition_AddNew(t *testing.T) {
	g := NewWithT(t)

	dr := result.New("component", "test", "check", "description")
	results.SetCondition(dr, check.NewCondition(check.ConditionTypeCompatible, metav1.ConditionTrue, check.WithReason("TestReason"), check.WithMessage("test message")))

	g.Expect(dr.Status.Conditions).To(HaveLen(1))
	g.Expect(dr.Status.Conditions[0].Condition).To(MatchFields(IgnoreExtras, Fields{
		"Type":    Equal(check.ConditionTypeCompatible),
		"Status":  Equal(metav1.ConditionTrue),
		"Reason":  Equal("TestReason"),
		"Message": Equal("test message"),
	}))
}

func TestSetCondition_UpdateExisting(t *testing.T) {
	g := NewWithT(t)

	dr := result.New("component", "test", "check", "description")

	// First call adds new condition
	results.SetCondition(dr, check.NewCondition(check.ConditionTypeCompatible, metav1.ConditionTrue, check.WithReason("reason1"), check.WithMessage("message1")))
	g.Expect(dr.Status.Conditions).To(HaveLen(1))

	// Second call with same type updates existing
	results.SetCondition(dr, check.NewCondition(check.ConditionTypeCompatible, metav1.ConditionFalse, check.WithReason("reason2"), check.WithMessage("message2")))
	g.Expect(dr.Status.Conditions).To(HaveLen(1))
	g.Expect(dr.Status.Conditions[0].Condition).To(MatchFields(IgnoreExtras, Fields{
		"Status":  Equal(metav1.ConditionFalse),
		"Reason":  Equal("reason2"),
		"Message": Equal("message2"),
	}))
}

func TestSetCondition_MultipleConditionTypes(t *testing.T) {
	g := NewWithT(t)

	dr := result.New("component", "test", "check", "description")

	// Add first condition type
	results.SetCondition(dr, check.NewCondition(check.ConditionTypeCompatible, metav1.ConditionTrue, check.WithReason("reason1"), check.WithMessage("message1")))
	g.Expect(dr.Status.Conditions).To(HaveLen(1))

	// Add different condition type
	results.SetCondition(dr, check.NewCondition(check.ConditionTypeAvailable, metav1.ConditionTrue, check.WithReason("reason2"), check.WithMessage("message2")))
	g.Expect(dr.Status.Conditions).To(HaveLen(2))

	// Update first condition type
	results.SetCondition(dr, check.NewCondition(check.ConditionTypeCompatible, metav1.ConditionFalse, check.WithReason("reason3"), check.WithMessage("message3")))
	g.Expect(dr.Status.Conditions).To(HaveLen(2))

	// Verify both conditions exist with correct values
	g.Expect(dr.Status.Conditions[0].Condition).To(MatchFields(IgnoreExtras, Fields{
		"Type":   Equal(check.ConditionTypeCompatible),
		"Status": Equal(metav1.ConditionFalse),
	}))
	g.Expect(dr.Status.Conditions[1].Condition).To(MatchFields(IgnoreExtras, Fields{
		"Type":   Equal(check.ConditionTypeAvailable),
		"Status": Equal(metav1.ConditionTrue),
	}))
}

func TestSetCompatibilitySuccess_Simple(t *testing.T) {
	g := NewWithT(t)

	dr := result.New("component", "test", "check", "description")
	results.SetCondition(dr, check.NewCondition(
		check.ConditionTypeCompatible,
		metav1.ConditionTrue,
		check.WithReason(check.ReasonVersionCompatible),
		check.WithMessage("simple success message"),
	))

	g.Expect(dr.Status.Conditions).To(HaveLen(1))
	g.Expect(dr.Status.Conditions[0].Condition).To(MatchFields(IgnoreExtras, Fields{
		"Type":    Equal(check.ConditionTypeCompatible),
		"Status":  Equal(metav1.ConditionTrue),
		"Reason":  Equal(check.ReasonVersionCompatible),
		"Message": Equal("simple success message"),
	}))
}

func TestSetCompatibilitySuccess_WithFormatting(t *testing.T) {
	g := NewWithT(t)

	dr := result.New("component", "test", "check", "description")
	results.SetCondition(dr, check.NewCondition(
		check.ConditionTypeCompatible,
		metav1.ConditionTrue,
		check.WithReason(check.ReasonVersionCompatible),
		check.WithMessage("State: %s is compatible with version %s", "Managed", "3.0"),
	))

	g.Expect(dr.Status.Conditions).To(HaveLen(1))
	g.Expect(dr.Status.Conditions[0].Condition).To(MatchFields(IgnoreExtras, Fields{
		"Type":    Equal(check.ConditionTypeCompatible),
		"Status":  Equal(metav1.ConditionTrue),
		"Reason":  Equal(check.ReasonVersionCompatible),
		"Message": Equal("State: Managed is compatible with version 3.0"),
	}))
}

func TestSetCompatibilityFailure_Simple(t *testing.T) {
	g := NewWithT(t)

	dr := result.New("component", "test", "check", "description")
	results.SetCondition(dr, check.NewCondition(
		check.ConditionTypeCompatible,
		metav1.ConditionFalse,
		check.WithReason(check.ReasonVersionIncompatible),
		check.WithMessage("simple failure message"),
	))

	g.Expect(dr.Status.Conditions).To(HaveLen(1))
	g.Expect(dr.Status.Conditions[0].Condition).To(MatchFields(IgnoreExtras, Fields{
		"Type":    Equal(check.ConditionTypeCompatible),
		"Status":  Equal(metav1.ConditionFalse),
		"Reason":  Equal(check.ReasonVersionIncompatible),
		"Message": Equal("simple failure message"),
	}))
}

func TestSetCompatibilityFailure_WithFormatting(t *testing.T) {
	g := NewWithT(t)

	dr := result.New("component", "test", "check", "description")
	results.SetCondition(dr, check.NewCondition(
		check.ConditionTypeCompatible,
		metav1.ConditionFalse,
		check.WithReason(check.ReasonVersionIncompatible),
		check.WithMessage("State: %s is incompatible with version %s", "Removed", "3.0"),
	))

	g.Expect(dr.Status.Conditions).To(HaveLen(1))
	g.Expect(dr.Status.Conditions[0].Condition).To(MatchFields(IgnoreExtras, Fields{
		"Type":    Equal(check.ConditionTypeCompatible),
		"Status":  Equal(metav1.ConditionFalse),
		"Reason":  Equal(check.ReasonVersionIncompatible),
		"Message": Equal("State: Removed is incompatible with version 3.0"),
	}))
}
