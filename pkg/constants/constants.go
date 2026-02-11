package constants

// Management state values for components and services.
const (
	ManagementStateManaged   = "Managed"
	ManagementStateUnmanaged = "Unmanaged"
	ManagementStateRemoved   = "Removed"
)

// Component names used across multiple package groups.
const (
	ComponentDashboard        = "dashboard"
	ComponentKServe           = "kserve"
	ComponentTrainingOperator = "trainingoperator"
)
