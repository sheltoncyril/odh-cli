# Quickstart: Constitution Alignment Audit

**Feature**: Constitution Alignment Audit
**Branch**: 002-constitution-alignment
**Date**: 2025-12-07

## Overview

This guide helps developers understand and implement constitutional alignment refactoring for odh-cli. The goal is to refactor existing code to comply with constitution v1.14.0 standards.

## Prerequisites

- Go 1.25.0 installed
- odh-cli repository cloned
- Branch `002-constitution-alignment` checked out
- Familiarity with constitution v1.14.0 (`.specify/memory/constitution.md`)

## Quick Start (5 minutes)

### 1. Understand the Scope

This refactoring addresses 5 TODO items from constitution v1.14.0:

1. **Package Isolation**: Move diagnostic checks to isolated packages
2. **Mock Centralization**: Move test mocks to `pkg/util/test/mocks/`
3. **Command Package Validation**: Verify command isolation
4. **Comment Quality**: Remove obvious comments
5. **Gomega Assertions**: Use struct matchers instead of field access

### 2. Set Up Development Environment

```bash
# Checkout feature branch
git checkout 002-constitution-alignment

# Verify build works
make build

# Run existing tests
make test

# Run quality checks
make check
```

### 3. Verify Current State

```bash
# Check diagnostic check file structure (BEFORE refactoring)
ls -la pkg/doctor/checks/components/
# Expected: dashboard.go, modelmesh_removal.go, kserve_serverless_removal.go

# Check for inline mocks
grep -r "type Mock" pkg/doctor/check/selector_test.go
# Expected: MockCheck struct definition

# Check command package structure
ls -la pkg/cmd/doctor/
# Expected: lint_options.go, upgrade_options.go, shared_options.go
```

## Refactoring Workflows

### Workflow 1: Package Isolation (Priority 1)

**Goal**: Move each diagnostic check to its own package.

**Steps**:

```bash
# 1. Create new package directory
mkdir -p pkg/doctor/checks/components/dashboard

# 2. Move check implementation
mv pkg/doctor/checks/components/dashboard.go \
   pkg/doctor/checks/components/dashboard/dashboard.go

# 3. Move corresponding test
mv pkg/doctor/checks/components/dashboard_test.go \
   pkg/doctor/checks/components/dashboard/dashboard_test.go

# 4. Update package declaration in moved files
# Change: package components
# To:     package dashboard

# 5. Update imports across codebase
# Find all imports: grep -r "odh-cli/pkg/doctor/checks/components" .
# Update to: odh-cli/pkg/doctor/checks/components/dashboard

# 6. Fix constant/type naming (remove package name repetition)
# Before: const dashboardCheckID = "components.dashboard"
# After:  const checkID = "components.dashboard"

# 7. Run tests
make test

# 8. Run quality checks
make check

# 9. Commit
git add pkg/doctor/checks/components/dashboard/
git commit -m "T###: Isolate dashboard check into dedicated package"
```

**Repeat for**:
- `modelmesh_removal.go` → `modelmesh/modelmesh.go`
- `kserve_serverless_removal.go` → `kserve/kserve.go`

### Workflow 2: Mock Centralization (Priority 2)

**Goal**: Move inline mocks to centralized location using testify/mock.

**Steps**:

```bash
# 1. Create mock package directory
mkdir -p pkg/util/test/mocks/check

# 2. Create centralized mock file
cat > pkg/util/test/mocks/check/check.go <<'EOF'
package mocks

import (
    "context"
    "github.com/stretchr/testify/mock"
    "github.com/lburgazzoli/odh-cli/pkg/doctor/check"
)

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
EOF

# 3. Update test to use centralized mock
# In pkg/doctor/check/selector_test.go:
# - Remove inline MockCheck struct definition
# - Import: mocks "github.com/lburgazzoli/odh-cli/pkg/util/test/mocks/check"
# - Usage: mockCheck := mocks.NewMockCheck()

# 4. Run tests
make test

# 5. Run quality checks
make check

# 6. Commit
git add pkg/util/test/mocks/check/ pkg/doctor/check/selector_test.go
git commit -m "T###: Centralize MockCheck using testify/mock framework"
```

### Workflow 3: Comment Quality Review (Priority 4)

**Goal**: Remove obvious comments, retain WHY comments.

**Steps**:

```bash
# 1. Find obvious comments
grep -r "^// Get\|^// Check\|^// Set\|^// Create" pkg/ --include="*.go"

# 2. Review each comment
# Ask: "Does this explain WHY, or just restate WHAT the code does?"

# 3. Remove obvious comments or rewrite to explain WHY
# Before: // Get the DataScienceCluster singleton
#         dsc, err := target.Client.GetDataScienceCluster(ctx)
# After:  (remove comment - code is self-explanatory)

# 4. Ensure godoc comments exist for exported functions
# Required: // NewCheck creates a dashboard check instance.
#           func NewCheck() check.Check { ... }

# 5. Run linter
make lint

# 6. Commit
git add <changed-files>
git commit -m "T###: Remove obvious comments per Code Comments standard"
```

### Workflow 4: Gomega Assertion Refactoring (Priority 5)

**Goal**: Use `HaveField()` and `MatchFields()` instead of direct field access.

**Steps**:

```bash
# 1. Find direct field access assertions
grep -r "g.Expect.*\\..*).To(" pkg/ --include="*_test.go"

# 2. Refactor to struct matchers
# Before: g.Expect(result.Status).To(Equal(check.StatusPass))
# After:  g.Expect(result).To(HaveField("Status", check.StatusPass))

# Before: g.Expect(result.Status).To(Equal(check.StatusPass))
#         g.Expect(result.Severity).To(BeNil())
#         g.Expect(result.Message).To(ContainSubstring("ready"))
# After:  g.Expect(result).To(MatchFields(IgnoreExtras, Fields{
#             "Status":   Equal(check.StatusPass),
#             "Severity": BeNil(),
#             "Message":  ContainSubstring("ready"),
#         }))

# 3. Run tests
make test

# 4. Commit
git add <changed-test-files>
git commit -m "T###: Refactor test assertions to use Gomega struct matchers"
```

## Quality Gates

After **every** refactoring step:

```bash
# 1. Run linter
make lint

# 2. Run vulnerability check
make vulncheck

# 3. Run tests
make test

# 4. Run all quality checks
make check
```

**All must pass before committing.**

## Common Issues and Solutions

### Issue 1: Import Path Updates

**Problem**: After moving files, imports are broken.

**Solution**:
```bash
# Use Go's automatic import fixing
go mod tidy
go fmt ./...

# Or use goimports
goimports -w .
```

### Issue 2: Circular Import Dependencies

**Problem**: Check packages import each other.

**Solution**:
- Checks should be independent (no cross-check imports)
- Move shared logic to `pkg/doctor/checks/shared/validation/`
- Each check package must be self-contained

### Issue 3: Test Coverage Drop

**Problem**: After refactoring, test coverage decreases.

**Solution**:
```bash
# Check coverage before refactoring
go test -coverprofile=before.out ./...
go tool cover -func=before.out

# Check coverage after refactoring
go test -coverprofile=after.out ./...
go tool cover -func=after.out

# Coverage should remain stable or improve
```

### Issue 4: Obvious vs. Necessary Comment Distinction

**Problem**: Unsure if comment is obvious or necessary.

**Solution**:
- Ask: "If I removed the comment, would code still be clear?"
- If YES → Remove comment (code is self-documenting)
- If NO → Keep comment but ensure it explains WHY, not WHAT
- Exception: godoc comments are always required for exported identifiers

## Validation Checklist

Before marking refactoring complete:

- [ ] All diagnostic checks isolated in `pkg/doctor/checks/<category>/<check>/`
- [ ] No inline mock definitions in test files
- [ ] All mocks use testify/mock and reside in `pkg/util/test/mocks/`
- [ ] Command package structure validated (lint/ and upgrade/ isolated)
- [ ] No obvious comments (stating WHAT code does)
- [ ] All exported functions have godoc comments
- [ ] All test assertions use Gomega struct matchers
- [ ] `make check` passes
- [ ] `make test` passes
- [ ] Test coverage stable or improved
- [ ] One commit per task with T### ID in message

## References

- Constitution v1.14.0: `.specify/memory/constitution.md`
- Feature Spec: `specs/002-constitution-alignment/spec.md`
- Implementation Plan: `specs/002-constitution-alignment/plan.md`
- Research Findings: `specs/002-constitution-alignment/research.md`

## Next Steps

After completing refactoring:

1. Run full quality gate: `make check && make test`
2. Update CLAUDE.md if needed (document any new patterns)
3. Create pull request from `002-constitution-alignment` to `main`
4. Reference constitution v1.14.0 in PR description
5. Mark TODO items in constitution SYNC IMPACT REPORT as complete