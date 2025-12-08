# Data Model: Constitution Alignment Audit

**Feature**: Constitution Alignment Audit
**Branch**: 002-constitution-alignment
**Date**: 2025-12-07

## Overview

This refactoring task does not introduce new data models. It reorganizes existing code structures to align with constitutional principles. This document describes the structural entities being refactored.

## Structural Entities

### 1. Diagnostic Check Package

**Description**: A self-contained package containing a single diagnostic check implementation.

**Structure**:
```
pkg/doctor/checks/<category>/<check>/
├── <check>.go       # Check implementation (struct, methods, constants)
└── <check>_test.go  # Unit tests for the check
```

**Constraints**:
- Package name must match check domain (e.g., `modelmesh`, `kserve`, `dashboard`)
- File name must match package name (e.g., `modelmesh.go` in `modelmesh/` package)
- No cross-check dependencies (each check independently testable)
- Package constants must not repeat package name (e.g., `checkID`, not `modelmeshCheckID`)

**Examples**:
- `pkg/doctor/checks/components/dashboard/dashboard.go` - Dashboard check
- `pkg/doctor/checks/components/modelmesh/modelmesh.go` - Modelmesh removal check
- `pkg/doctor/checks/components/kserve/kserve.go` - KServe serverless removal check

### 2. Centralized Mock

**Description**: A test double using testify/mock framework, placed in centralized location for reuse.

**Structure**:
```go
package mocks

import "github.com/stretchr/testify/mock"

type MockCheck struct {
    mock.Mock
}

func NewMockCheck() *MockCheck {
    return &MockCheck{}
}

func (m *MockCheck) ID() string {
    args := m.Called()
    return args.String(0)
}

func (m *MockCheck) Category() check.CheckCategory {
    args := m.Called()
    return args.Get(0).(check.CheckCategory)
}

func (m *MockCheck) Run(ctx context.Context, target check.CheckTarget) check.Result {
    args := m.Called(ctx, target)
    return args.Get(0).(check.Result)
}
```

**Location**: `pkg/util/test/mocks/<package>/<interface>.go`

**Constraints**:
- Must use `testify/mock` framework (embed `mock.Mock`)
- Package name must be `mocks` (not `mocks_<package>`)
- Must provide constructor function (`NewMock<Type>()`)
- All interface methods must use `m.Called(args...)` pattern

### 3. Command Package

**Description**: Isolated package containing a single command's business logic.

**Structure**:
```
pkg/cmd/doctor/
├── shared.go              # Shared options, types, utilities (parent-level)
├── lint/
│   ├── options.go         # LintOptions struct
│   ├── run.go             # Lint command business logic
│   └── options_test.go    # Lint tests
└── upgrade/
    ├── options.go         # UpgradeOptions struct
    ├── run.go             # Upgrade command business logic
    └── options_test.go    # Upgrade tests
```

**Constraints**:
- Each command in its own package (no mixing)
- Shared code in parent-level `shared.go` only
- Each command package contains its own `options.go`
- No circular dependencies between sibling commands

### 4. Code Comment

**Description**: Documentation explaining WHY code exists (not WHAT it does).

**Categories**:

**Prohibited** (obvious comments):
```go
// Get the DataScienceCluster singleton
dsc, err := target.Client.GetDataScienceCluster(ctx)
```

**Required** (WHY comments):
```go
// ServiceMesh is deprecated in 3.x but Unmanaged state must still be checked
// because users may have manually deployed service mesh operators
if managementStateStr == "Managed" || managementStateStr == "Unmanaged" {
```

**Required** (godoc):
```go
// NewCheck creates a new dashboard check instance.
func NewCheck() check.Check { ... }
```

**Constraints**:
- Comments must explain non-obvious WHY, not obvious WHAT
- All exported identifiers must have godoc comments
- Security-sensitive code must be commented
- Workarounds must reference issue/bug numbers

### 5. Gomega Test Assertion

**Description**: Struct field validation using Gomega matchers.

**Patterns**:

**Single field**:
```go
g.Expect(result).To(HaveField("Status", check.StatusPass))
```

**Multiple fields**:
```go
g.Expect(result).To(MatchFields(IgnoreExtras, Fields{
    "Status":   Equal(check.StatusPass),
    "Severity": BeNil(),
    "Message":  ContainSubstring("ready"),
}))
```

**Constraints**:
- No direct field access in assertions (e.g., `g.Expect(result.Status)`)
- Use `HaveField()` for single field checks
- Use `MatchFields()` for multiple field checks
- Provides better failure diagnostics with full struct context

## Refactoring Mapping

| Current State | Target State | Rationale |
|---------------|--------------|-----------|
| `pkg/doctor/checks/components/dashboard.go` | `pkg/doctor/checks/components/dashboard/dashboard.go` | Package isolation per Principle |
| `pkg/doctor/checks/components/modelmesh_removal.go` | `pkg/doctor/checks/components/modelmesh/modelmesh.go` | Package isolation per Principle |
| `pkg/doctor/checks/components/kserve_serverless_removal.go` | `pkg/doctor/checks/components/kserve/kserve.go` | Package isolation per Principle |
| `pkg/doctor/check/selector_test.go` (MockCheck inline) | `pkg/util/test/mocks/check/check.go` | Mock centralization standard |
| `g.Expect(result.Field).To(...)` | `g.Expect(result).To(HaveField("Field", ...))` | Gomega struct matchers |
| `// Get something` comments | Remove or rewrite to explain WHY | Code comment quality standard |

## No New Data Introduced

This refactoring does not:
- Add new custom resources or Kubernetes objects
- Introduce new configuration formats
- Create new API contracts
- Modify data persistence patterns

All changes are structural (code organization) and stylistic (comment quality, test patterns).