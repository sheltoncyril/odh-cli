# Implementation Plan: Constitution Alignment Audit

**Branch**: `002-constitution-alignment` | **Date**: 2025-12-07 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/002-constitution-alignment/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

**REVISED 2025-12-07**: After codebase investigation, scope adjusted from original assumptions.

Refactor the odh-cli codebase to align with constitution v1.14.0 requirements. This includes: (1) **validating** existing diagnostic check package isolation (already compliant), (2) centralizing test mocks using testify/mock framework, (3) **refactoring** command package structure to isolate lint/ and upgrade/ commands, (4) removing obvious code comments, and (5) refactoring test assertions to use Gomega struct matchers.

**Key Finding**: Diagnostic checks are already properly isolated in dedicated packages (codeflare/, kserve/, kueue/, modelmesh/). Main refactoring needed is command package structure (pkg/cmd/doctor/ not following Principle XI isolation).

## Technical Context

**Language/Version**: Go 1.25.0
**Primary Dependencies**:
- github.com/stretchr/testify/mock (test mocking framework)
- github.com/onsi/gomega (assertion library)
- k8s.io/client-go/dynamic/fake (Kubernetes fake client for unit tests)
- github.com/lburgazzoli/k3s-envtest (integration testing)
- github.com/spf13/cobra v1.10.1 (CLI framework)

**Storage**: N/A (CLI tool, no persistent storage)
**Testing**: Go test with gomega assertions and testify/mock framework
**Target Platform**: Linux, macOS, Windows (kubectl plugin, cross-platform)
**Project Type**: Single CLI project (kubectl plugin)
**Performance Goals**: Refactoring quality metrics (100% package isolation, zero inline mocks, zero obvious comments)
**Constraints**: Must maintain backward compatibility, all refactoring must pass `make check` and `make test`
**Scale/Scope**: Validate 4 isolated check packages, refactor 2 command packages (lint/, upgrade/), centralize 1 mock, remove ~20 obvious comments, refactor ~164 gomega assertions

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Phase 0 Gates (Research)

- ✅ **kubectl Plugin Integration**: N/A - refactoring existing plugin, no new integration needed
- ✅ **Output Format Consistency**: N/A - no changes to command output formats
- ✅ **High-Level Resource Checks (Principle IX)**: N/A - not creating new checks, refactoring existing ones
- ✅ **Cluster-Wide Scope (Principle X)**: N/A - no changes to diagnostic scope

**Status**: PASS - All Phase 0 gates satisfied (refactoring task, not new feature development)

### Phase 1 Gates (Design)

- ✅ **Command Structure (Principle II)**: Validating existing Complete/Validate/Run pattern compliance
- ✅ **Functional Options (Principle IV)**: No changes to option patterns
- ✅ **Package Granularity**: Core refactoring objective - isolating checks per Package Granularity principle
- ✅ **Command Package Isolation (Principle XI)**: Validating existing pkg/cmd/doctor/ structure

**Status**: PASS - Design aligns with constitutional requirements

### Phase 2 Gates (Implementation)

- ⚠️ **Error Handling (Principle V)**: Verify refactored code maintains proper error wrapping
- ⚠️ **Test Coverage (Principle VI)**: Critical - all refactoring must maintain test coverage using gomega/testify patterns
- ⚠️ **Mock Organization**: Core objective - centralizing mocks per Mock Organization standard
- ⚠️ **JQ-Based Field Access (Principle VII)**: Verify no regressions during refactoring
- ⚠️ **Centralized GVK/GVR (Principle VIII)**: Verify no regressions during refactoring
- ⚠️ **Message Constants**: Verify message constants remain package-level during refactoring
- ⚠️ **Code Comments**: Core objective - removing obvious comments per Code Comments standard
- ⚠️ **Gomega Assertions**: Core objective - refactoring to struct matchers per Principle VI
- ⚠️ **make check Execution**: MANDATORY after each refactoring step
- ⚠️ **Commit Granularity**: One commit per task or grouped related tasks with T### IDs

**Status**: PENDING - Will validate during implementation phase

### Complexity Justification

No constitutional violations expected. This feature implements constitutional requirements, not exceptions to them.

## Project Structure

### Documentation (this feature)

```text
specs/002-constitution-alignment/
├── spec.md              # Feature specification
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
├── checklists/          # Quality checklists
│   └── requirements.md  # Specification quality checklist
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

**CURRENT STATE** (2025-12-07 investigation):

```text
odh-cli/
├── cmd/                          # Cobra command wrappers
│   ├── main.go
│   └── doctor/
│       ├── doctor.go
│       ├── lint.go               # Needs import update after pkg/ refactor
│       └── upgrade.go            # Needs import update after pkg/ refactor
├── pkg/
│   ├── cmd/doctor/               # Command business logic (❌ REFACTOR TARGET)
│   │   ├── shared_options.go    # ✅ Correct location (shared)
│   │   ├── lint_options.go      # ❌ MOVE → lint/options.go
│   │   ├── upgrade_options.go   # ❌ MOVE → upgrade/options.go
│   │   └── lint_integration_test.go.disabled  # ❌ MOVE → lint/
│   ├── doctor/
│   │   ├── check/                # Check framework
│   │   │   ├── check.go
│   │   │   ├── registry.go
│   │   │   ├── executor.go
│   │   │   └── selector_test.go  # ⚠️ Contains inline MockCheck (REFACTOR TARGET)
│   │   └── checks/               # Check implementations
│   │       ├── components/       # Component checks (✅ ALREADY ISOLATED)
│   │       │   ├── codeflare/    # ✅ Isolated package
│   │       │   │   ├── codeflare.go
│   │       │   │   └── codeflare_test.go
│   │       │   ├── kserve/       # ✅ Isolated package
│   │       │   │   ├── kserve.go
│   │       │   │   └── kserve_test.go
│   │       │   ├── kueue/        # ✅ Isolated package
│   │       │   │   ├── kueue.go
│   │       │   │   └── kueue_test.go
│   │       │   └── modelmesh/    # ✅ Isolated package
│   │       │       ├── modelmesh.go
│   │       │       └── modelmesh_test.go
│   │       ├── dependencies/     # Dependency checks (✅ ALREADY ISOLATED)
│   │       │   ├── certmanager/
│   │       │   ├── kueueoperator/
│   │       │   └── servicemeshoperator/
│   │       ├── services/         # Service checks (✅ ALREADY ISOLATED)
│   │       │   └── servicemesh/
│   │       └── shared/           # Shared check utilities
│   │           └── results/
│   ├── resources/                # Centralized GVK/GVR definitions
│   └── util/
│       ├── client/
│       ├── jq/
│       └── test/mocks/           # Centralized mocks (🆕 CREATE TARGET)
│           └── check/            # Mock for check.Check interface
│               └── check.go      # MockCheck implementation (to be created)
├── .specify/memory/
│   └── constitution.md           # v1.14.0 with TODO items
└── CLAUDE.md                     # Project instructions (UPDATE TARGET)
```

**TARGET STATE** (after refactoring):

```text
pkg/cmd/doctor/
├── shared_options.go             # Shared code (stays)
├── lint/                         # 🆕 Isolated lint command
│   ├── options.go                # Moved from lint_options.go
│   └── integration_test.go.disabled  # Moved from parent
└── upgrade/                      # 🆕 Isolated upgrade command
    └── options.go                # Moved from upgrade_options.go
```

**Refactoring Summary**:
1. ✅ **Diagnostic checks**: Already isolated (codeflare/, kserve/, kueue/, modelmesh/) - **VALIDATION ONLY**
2. ❌ **Command structure**: Needs refactoring to isolate lint/ and upgrade/ - **ACTIVE REFACTORING**
3. ⚠️ **MockCheck**: Move from selector_test.go to pkg/util/test/mocks/check/check.go - **ACTIVE REFACTORING**
4. ⚠️ **Comments**: ~20 obvious comments to remove across 10 files - **ACTIVE REFACTORING**
5. ⚠️ **Gomega assertions**: ~164 direct field access patterns across 11 test files - **ACTIVE REFACTORING**

## Complexity Tracking

N/A - No constitutional violations. This feature implements constitutional requirements.
