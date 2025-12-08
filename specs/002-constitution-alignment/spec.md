# Feature Specification: Constitution Alignment Audit

**Feature Branch**: `002-constitution-alignment`
**Created**: 2025-12-07
**Status**: Draft
**Input**: User description: "ensure the codebase is aligned with the constitution"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Diagnostic Check Package Isolation Validation (Priority: P1)

A developer needs to verify that all existing diagnostic checks follow the package isolation pattern mandated by the constitution (v1.14.0). Each check should reside in its own dedicated package under `pkg/doctor/checks/<category>/<check>/` rather than coexisting with other checks in a shared category package.

**Why this priority**: This is the highest priority because it validates that the constitutional package isolation pattern is followed. Current state shows checks are already properly isolated (codeflare/, kserve/, kueue/, modelmesh/), so this story confirms compliance and documents the pattern for future checks.

**Independent Test**: Can be fully tested by examining the package structure under `pkg/doctor/checks/` and verifying each check has its own isolated package. Delivers confirmation of architectural compliance.

**Acceptance Scenarios**:

1. **Given** diagnostic checks exist in `pkg/doctor/checks/components/`, **When** examining the package structure, **Then** each check must be in its own package (e.g., `codeflare/codeflare.go`, `modelmesh/modelmesh.go`, `kserve/kserve.go`, `kueue/kueue.go`)
2. **Given** check implementations exist, **When** reviewing file organization, **Then** each check package must contain only files related to that specific check (implementation, tests, constants)
3. **Given** multiple checks exist in the same category, **When** running tests, **Then** each check must be independently testable without dependencies on sibling checks

---

### User Story 2 - Mock Centralization (Priority: P2)

A developer needs to ensure test mocks follow the centralized organization pattern using testify/mock framework, eliminating inline mock definitions scattered across test files.

**Why this priority**: This is second priority because it improves test maintainability and prevents duplication, but doesn't block functionality. The constitution mandates this refactoring.

**Independent Test**: Can be fully tested by searching for inline mock struct definitions in test files and verifying all mocks are centralized in `pkg/util/test/mocks/`. Delivers improved test quality and reusability.

**Acceptance Scenarios**:

1. **Given** test files contain mock implementations, **When** examining test code, **Then** all mocks must use testify/mock framework and be located in `pkg/util/test/mocks/<package>/`
2. **Given** the MockCheck struct exists in `selector_test.go`, **When** refactoring, **Then** it must be moved to `pkg/util/test/mocks/check/check.go` and use testify/mock
3. **Given** tests need to mock the Check interface, **When** writing tests, **Then** they must import the centralized mock from `pkg/util/test/mocks/check` rather than defining inline mocks

---

### User Story 3 - Command Package Structure Refactoring (Priority: P3)

A developer needs to refactor command packages to follow the isolation pattern defined in Principle XI, ensuring each command under `pkg/cmd/doctor/` resides in its own dedicated package. Current state shows `lint_options.go` and `upgrade_options.go` are in the parent package instead of isolated subdirectories.

**Why this priority**: This is third priority because it requires significant refactoring to achieve constitutional compliance. Current violations: lint and upgrade commands are NOT in isolated packages.

**Independent Test**: Can be fully tested by examining the `pkg/cmd/doctor/` package structure and verifying lint and upgrade commands are properly isolated in separate subdirectories. Delivers clear separation of concerns between commands.

**Acceptance Scenarios**:

1. **Given** the doctor command has lint and upgrade subcommands, **When** refactoring `pkg/cmd/doctor/`, **Then** each subcommand must have its own package (e.g., `lint/`, `upgrade/`)
2. **Given** `lint_options.go` exists in `pkg/cmd/doctor/`, **When** refactoring, **Then** it must be moved to `pkg/cmd/doctor/lint/options.go`
3. **Given** `upgrade_options.go` exists in `pkg/cmd/doctor/`, **When** refactoring, **Then** it must be moved to `pkg/cmd/doctor/upgrade/options.go`
4. **Given** `shared_options.go` exists, **When** reviewing code organization, **Then** it must remain in `pkg/cmd/doctor/shared_options.go` (parent-level shared code)

---

### User Story 4 - Code Comment Quality Review (Priority: P4)

A developer needs to identify and remove obvious comments that violate the constitution's code comment guidelines (v1.14.0), ensuring comments explain WHY rather than WHAT.

**Why this priority**: This is fourth priority because it improves code quality and readability but doesn't affect functionality. It's a cleanup task that enhances maintainability.

**Independent Test**: Can be fully tested by reviewing code comments across the codebase and identifying violations of the comment quality guidelines. Delivers cleaner, more maintainable code.

**Acceptance Scenarios**:

1. **Given** code contains comments, **When** reviewing comment content, **Then** comments must not state obvious facts (e.g., "Get the DataScienceCluster singleton")
2. **Given** code contains workarounds or non-obvious logic, **When** reviewing comments, **Then** comments must explain WHY the approach was taken, not WHAT the code does
3. **Given** exported functions exist, **When** reviewing code, **Then** godoc comments must be present for all exported identifiers

---

### User Story 5 - Gomega Assertion Pattern Compliance (Priority: P5)

A developer needs to refactor existing tests to use Gomega's struct matchers (HaveField, MatchFields) instead of individual field assertions, following the updated Principle VI guidelines.

**Why this priority**: This is fifth priority because it improves test failure diagnostics but doesn't affect test correctness. It's a quality improvement that makes debugging easier.

**Independent Test**: Can be fully tested by searching for individual field assertion patterns (e.g., `g.Expect(obj.Field).To(Equal(value))`) and verifying they're replaced with struct matchers. Delivers better error messages when tests fail.

**Acceptance Scenarios**:

1. **Given** tests validate struct fields, **When** examining test assertions, **Then** they must use `HaveField` or `MatchFields` matchers instead of accessing fields directly
2. **Given** test assertions check multiple fields, **When** reviewing test code, **Then** they must use `MatchFields(IgnoreExtras, Fields{...})` for clearer failure diagnostics
3. **Given** existing tests use individual field assertions, **When** refactoring, **Then** they must be replaced with struct field matchers while maintaining test coverage

---

### Edge Cases

- What happens when a check package needs to share utility code with other checks in the same category?
- How do we handle existing test coverage when refactoring mock implementations?
- What happens when inline comments are necessary for debugging complex logic (e.g., algorithm explanations)?
- How do we handle godoc comments that may seem "obvious" but are required by Go documentation conventions?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST verify that all diagnostic checks under `pkg/doctor/checks/` are in isolated packages following the pattern `pkg/doctor/checks/<category>/<check>/<check>.go` (current state: codeflare/, kserve/, kueue/, modelmesh/ already isolated)
- **FR-002**: The system MUST document the package isolation pattern for future diagnostic checks
- **FR-003**: All test mocks MUST use the testify/mock framework from `github.com/stretchr/testify/mock`
- **FR-004**: The MockCheck struct MUST be moved from `selector_test.go` to `pkg/util/test/mocks/check/check.go`
- **FR-005**: All inline mock struct definitions MUST be replaced with testify/mock-based implementations in centralized locations
- **FR-006**: The system MUST refactor `pkg/cmd/doctor/lint_options.go` to `pkg/cmd/doctor/lint/options.go`
- **FR-007**: The system MUST refactor `pkg/cmd/doctor/upgrade_options.go` to `pkg/cmd/doctor/upgrade/options.go`
- **FR-008**: The system MUST update cmd/doctor/ Cobra wrappers to import from new isolated command packages
- **FR-009**: The system MUST update all imports across the codebase to reference the new command package paths
- **FR-010**: All obvious code comments (stating WHAT the code does) MUST be identified and removed or rewritten to explain WHY
- **FR-011**: All godoc comments on exported identifiers MUST remain in place per Go documentation conventions
- **FR-012**: All test assertions using individual field access (e.g., `g.Expect(obj.Field).To(Equal(value))`) MUST be refactored to use Gomega struct matchers
- **FR-013**: The system MUST run `make check` after all refactoring to ensure linting and vulnerability compliance
- **FR-014**: The system MUST run `make test` after all refactoring to ensure test coverage remains stable
- **FR-015**: Package naming MUST follow the constitution's naming conventions, avoiding redundant package name repetition in type/constant names

### Key Entities

- **Diagnostic Check**: A health check implementation residing in its own isolated package, containing check logic, tests, and constants
- **Mock Implementation**: A test double using testify/mock framework, centralized in `pkg/util/test/mocks/<package>/`
- **Command Package**: A dedicated package for a single command's business logic under `pkg/cmd/<parent>/<command>/`
- **Code Comment**: Documentation explaining WHY code exists, not WHAT it does, with exceptions for godoc and public API documentation

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of diagnostic checks are isolated in dedicated packages under `pkg/doctor/checks/<category>/<check>/`
- **SC-002**: Zero inline mock struct definitions exist in test files (excluding trivial test-specific cases)
- **SC-003**: All reusable mocks use testify/mock framework and reside in `pkg/util/test/mocks/`
- **SC-004**: `pkg/cmd/doctor/` structure follows command package isolation with separate `lint/` and `upgrade/` packages
- **SC-005**: Zero obvious comments exist in the codebase (comments stating WHAT code does)
- **SC-006**: 100% of struct validation tests use Gomega struct matchers instead of individual field assertions
- **SC-007**: `make check` passes with zero linting or vulnerability issues after all refactoring
- **SC-008**: `make test` passes with stable or improved test coverage after all refactoring

## Assumptions

- The codebase is currently functional and has existing test coverage
- Refactoring can be done incrementally without breaking functionality
- Constitution v1.14.0 is the authoritative source for coding standards
- Existing tests provide adequate coverage and only need pattern refactoring, not logic changes
- The TODO items in the constitution SYNC IMPACT REPORT accurately reflect current violations