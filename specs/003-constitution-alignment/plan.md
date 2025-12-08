# Implementation Plan: Constitution v1.15.0 Alignment

**Branch**: `003-constitution-alignment` | **Date**: 2025-12-08 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/home/luca/work/dev/openshift-ai/odh-cli/specs/003-constitution-alignment/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

This plan implements constitution v1.15.0 architectural changes through five independent user stories delivered incrementally:

1. **P1 - IOStreams Wrapper**: Create centralized output utility in `pkg/util/iostreams/` to eliminate repetitive `fmt.Fprintf(o.Out, ...)` patterns across all commands
2. **P2 - Unified Lint Command**: Merge separate lint and upgrade commands into single `lint` command with optional `--version` flag for upgrade-readiness assessment
3. **P3 - Command Interface**: Define standard Command interface with Complete/Validate/Run/AddFlags methods and refactor lint command to implement it
4. **P4 - Flexible Initialization**: Support both struct-based and functional options patterns for command initialization
5. **P5 - Deprecated API Audit**: Identify and replace all deprecated API usage throughout codebase

Each story will be completed independently with all tests passing before proceeding to the next, enabling safe incremental delivery without backward compatibility concerns.

## Technical Context

**Language/Version**: Go 1.25.0
**Primary Dependencies**:
- github.com/spf13/cobra v1.10.1 (CLI framework)
- github.com/spf13/pflag (flag management)
- k8s.io/cli-runtime v0.34.1 (kubectl integration)
- k8s.io/client-go v0.34.1 (Kubernetes client)

**Storage**: N/A (CLI tool, no persistent storage)
**Testing**:
- Unit tests: fake client from k8s.io/client-go/dynamic/fake
- Integration tests: k3s-envtest (github.com/lburgazzoli/k3s-envtest)
- Mocking: github.com/stretchr/testify/mock
- Assertions: Gomega (vanilla, no Ginkgo)

**Target Platform**: Linux/macOS/Windows (kubectl plugin)
**Project Type**: Single project (CLI tool)
**Performance Goals**: N/A (refactoring preserves existing performance)
**Constraints**:
- All quality gates (`make check`) must pass after each user story
- 30% reduction in output-related code (IOStreams wrapper)
- Zero new linting issues
- Existing tests must continue passing (except new AddFlags tests)

**Scale/Scope**:
- ~10 commands to refactor for IOStreams wrapper
- 2 commands to merge (lint + upgrade → lint)
- 1 Command interface definition
- ~5-10 deprecated API instances to audit and replace

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Phase 0 (Research) - ✅ **PASS**

✅ **kubectl Plugin Integration** (Principle I): Refactoring preserves existing kubectl plugin integration. No changes to binary name or kubeconfig usage.

✅ **High-Level Resource Checks** (Principle IX): Refactoring does not modify check targets. Existing checks remain focused on high-level CRs (DataScienceCluster, Notebook, InferenceService, etc.).

✅ **Cluster-Wide Diagnostic Scope** (Principle X): Refactoring does not modify diagnostic scope. No namespace filtering introduced.

### Phase 1 (Design) - ✅ **PASS**

✅ **Extensible Command Structure** (Principle II):
- Plan introduces Command interface with Complete/Validate/Run/AddFlags methods (User Story 3)
- Renames Options → Command, NewOptions → NewCommand
- Renames options.go → lint.go
- Commands initialize SharedOptions internally

✅ **Consistent Output Formats** (Principle III): Refactoring preserves existing `-o/--output` flag behavior for table/JSON/YAML formats.

✅ **Flexible Initialization Patterns** (Principle IV):
- Plan supports both struct-based initialization (preferred) and functional options (User Story 4)
- Options code in lint_options.go per convention

✅ **Strict Error Handling** (Principle V): Refactoring preserves existing error wrapping with fmt.Errorf(%w) and context propagation.

✅ **Test-First Development** (Principle VI):
- Each user story includes unit and integration tests
- Uses fake client for unit tests, k3s-envtest for integration
- Uses testify/mock for mocking (already adopted)
- Uses Gomega MatchFields pattern for struct assertions (already adopted)

✅ **JQ-Based Field Access** (Principle VII): Refactoring does not modify field access patterns. Existing JQ usage preserved.

✅ **Centralized Resource Type Definitions** (Principle VIII): Refactoring does not modify GVK/GVR definitions. Existing pkg/resources/types.go unchanged.

✅ **Doctor Command Architecture** (Principle XI):
- Plan merges lint and upgrade into single lint command (User Story 2)
- Optional --version flag switches between lint/upgrade modes
- Checks detect mode by comparing target.CurrentVersion vs target.Version

✅ **IOStreams Wrapper** (Development Standards): Plan implements IOStreams wrapper in pkg/util/iostreams/ (User Story 1) per constitutional requirement.

✅ **Deprecated API Avoidance** (Development Standards): Plan includes codebase audit (User Story 5) to identify and replace deprecated APIs.

✅ **Lint-Fix-First** (Development Workflow): Each user story will run `make lint-fix` before manual issue resolution per constitutional requirement.

### Phase 2 (Implementation) - To be verified during execution

- Error handling with fmt.Errorf(%w)
- Test coverage with fake client + k3s-envtest
- testify/mock for mocking (mocks in pkg/util/test/mocks)
- JQ-based field access for unstructured objects
- Centralized GVK/GVR definitions
- User-facing messages as package-level constants
- `make check` execution after each user story
- Full linting compliance
- One commit per user story with task IDs

## Project Structure

### Documentation (this feature)

```text
specs/003-constitution-alignment/
├── plan.md              # This file (/speckit.plan command output)
├── spec.md              # Feature specification (already created)
├── checklists/
│   └── requirements.md  # Specification quality checklist (already created)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

**Note**: No research.md, data-model.md, contracts/, or quickstart.md needed for this refactoring task. All technical decisions are specified in constitution v1.15.0.

### Source Code (repository root)

**Existing Structure** (unchanged):
```text
cmd/
├── main.go
└── doctor/
    ├── doctor.go         # Root doctor command
    ├── lint.go           # Lint command registration (Cobra wrapper)
    └── upgrade.go        # Upgrade command registration (TO BE REMOVED in P2)

pkg/
├── cmd/doctor/
│   ├── shared_options.go      # SharedOptions struct
│   ├── shared_options_test.go
│   ├── lint/
│   │   └── options.go         # Lint command business logic (TO BE REFACTORED)
│   └── upgrade/
│       └── options.go         # Upgrade command business logic (TO BE REMOVED in P2)
├── doctor/
│   ├── check/                 # Check framework
│   │   ├── check.go
│   │   ├── executor.go
│   │   └── registry.go
│   └── checks/                # Check implementations
│       ├── components/
│       ├── services/
│       └── workloads/
├── resources/
│   └── types.go               # Centralized GVK/GVR definitions (unchanged)
└── util/
    ├── client/
    ├── jq/
    └── test/mocks/            # Centralized mocks
```

**New Structure After Refactoring**:
```text
cmd/
├── main.go
└── doctor/
    ├── doctor.go              # Root doctor command (unchanged)
    └── lint.go                # Lint command registration (MODIFIED for --version flag)

pkg/
├── cmd/
│   ├── command.go             # NEW: Command interface definition (P3)
│   └── doctor/
│       ├── shared_options.go  # MODIFIED: embed IOStreams wrapper (P1)
│       ├── shared_options_test.go
│       └── lint/
│           ├── lint.go        # RENAMED from options.go, implements Command interface (P2, P3)
│           ├── lint_options.go # NEW: CommandOptions struct and functional options (P4)
│           └── lint_test.go   # MODIFIED: add AddFlags tests (P3)
├── doctor/
│   ├── check/                 # MODIFIED: checks detect mode via version comparison (P2)
│   │   ├── check.go
│   │   ├── executor.go
│   │   └── registry.go
│   └── checks/                # Check implementations (unchanged)
│       ├── components/
│       ├── services/
│       └── workloads/
├── resources/
│   └── types.go               # Unchanged
└── util/
    ├── client/
    ├── iostreams/             # NEW: IOStreams wrapper (P1)
    │   ├── iostreams.go
    │   └── iostreams_test.go
    ├── jq/
    └── test/mocks/
```

**Structure Decision**: Single project CLI tool structure preserved. New `pkg/util/iostreams/` package added for IOStreams wrapper (P1). Command interface defined in `pkg/cmd/command.go` (P3). Separate upgrade command removed and merged into lint (P2). File renamed from `options.go` → `lint.go` per constitution (P3).

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

**No violations**. All constitutional principles are followed:
- ✅ Command interface matches Principle II requirements
- ✅ IOStreams wrapper follows Development Standards
- ✅ Incremental delivery (P1→P2→P3→P4→P5) follows clarified implementation strategy
- ✅ Merged lint command follows Principle XI architecture
- ✅ All existing constitutional principles preserved (kubectl integration, test-first, JQ-based access, centralized GVK/GVR, etc.)

## Implementation Strategy

### Incremental Delivery Approach

Per clarification session 2025-12-08, implementation MUST follow user story priority order (P1 → P2 → P3 → P4 → P5) with each story completed and committed independently. Backward compatibility is not a concern, allowing internal API changes between stories.

Each user story MUST:
1. Leave codebase in working state with all tests passing
2. Pass all quality gates (`make check`)
3. Be committed with task ID following constitution commit granularity standards
4. Demonstrate measurable progress toward success criteria

### User Story Dependencies

```text
P1 (IOStreams Wrapper)
  ↓ [Can start immediately, no dependencies]
P2 (Unified Lint Command)
  ↓ [Can use IOStreams from P1, but not required]
P3 (Command Interface)
  ↓ [Requires P2 lint command structure to exist]
P4 (Flexible Initialization)
  ↓ [Requires P3 Command interface to exist]
P5 (Deprecated API Audit)
  ↓ [Independent, can run anytime after P1-P4 complete]
```

### Risk Mitigation

| Risk | Mitigation | Contingency |
|------|------------|-------------|
| IOStreams wrapper changes break existing output | Compare output before/after in tests; use Gomega matchers to verify exact output format | Revert to direct fmt.Fprintf if breaking changes found |
| Merging lint + upgrade creates test conflicts | Run existing lint and upgrade tests separately before merging; ensure both pass | Keep separate command structure if merge too complex |
| Command interface refactoring breaks existing tests | Preserve existing Complete/Validate/Run signatures; AddFlags is additive | Make AddFlags optional with default implementation |
| Functional options add complexity | Implement struct initialization first; add functional options as optional enhancement | Skip functional options if struct initialization sufficient |
| Deprecated API replacement introduces bugs | Replace one deprecated API at a time; run full test suite after each replacement | Document unavoidable deprecated APIs with TODO comments |

### Success Validation Per User Story

**P1 - IOStreams Wrapper**:
- ✅ pkg/util/iostreams/iostreams.go exists with Fprintf, Fprintln, Errorf, Errorln methods
- ✅ At least one command (e.g., lint) refactored to use IOStreams wrapper
- ✅ Output behavior identical to before refactoring (verified via tests)
- ✅ Code reduction of at least 30% in output-related lines (SC-001)
- ✅ All tests pass, `make check` passes

**P2 - Unified Lint Command**:
- ✅ cmd/doctor/upgrade.go deleted
- ✅ pkg/cmd/doctor/upgrade/ directory deleted
- ✅ pkg/cmd/doctor/lint/lint.go accepts optional --version flag
- ✅ `kubectl odh doctor lint` validates current state (no --version)
- ✅ `kubectl odh doctor lint --version 3.0` assesses upgrade readiness
- ✅ All existing lint and upgrade tests passing
- ✅ All tests pass, `make check` passes

**P3 - Command Interface**:
- ✅ pkg/cmd/command.go defines Command interface with Complete/Validate/Run/AddFlags
- ✅ pkg/cmd/doctor/lint/options.go renamed to lint.go
- ✅ lint.Command struct (not Options) implements Command interface
- ✅ lint.NewCommand() constructor (not NewOptions())
- ✅ AddFlags method registers flags with pflag.FlagSet
- ✅ New tests for AddFlags method
- ✅ All tests pass, `make check` passes

**P4 - Flexible Initialization**:
- ✅ pkg/cmd/doctor/lint/lint_options.go defines CommandOptions struct
- ✅ lint.NewCommand() accepts CommandOptions struct
- ✅ With* functional option functions defined (WithShared, WithTargetVersion, etc.)
- ✅ Both initialization patterns tested and produce identical behavior (SC-006)
- ✅ All tests pass, `make check` passes

**P5 - Deprecated API Audit**:
- ✅ All deprecated API usage identified via IDE warnings and godoc markers
- ✅ Deprecated APIs replaced with modern alternatives where available
- ✅ Unavoidable deprecated APIs documented with explanatory comments and tracking issues
- ✅ Zero deprecation warnings remain (SC-007)
- ✅ All tests pass, `make check` passes

## Phase 0: Research (N/A for Refactoring)

**Status**: ✅ **SKIPPED** - All technical decisions specified in constitution v1.15.0

This refactoring task implements constitutional requirements with no open technical questions:
- IOStreams wrapper design specified in constitution Development Standards
- Command interface signature specified in constitution Principle II
- Merged lint command architecture specified in constitution Principle XI
- Initialization patterns specified in constitution Principle IV
- Deprecated API handling specified in constitution Development Standards

All implementation details are deterministic from constitutional requirements. No research phase needed.

## Phase 1: Design (N/A for Refactoring)

**Status**: ✅ **SKIPPED** - No new data models, API contracts, or quickstart guides needed

This refactoring task modifies existing command structure and introduces utility infrastructure. Key design artifacts:

### IOStreams Wrapper Interface (P1)

Already specified in constitution:
```go
package iostreams

type IOStreams struct {
    In     io.Reader
    Out    io.Writer
    ErrOut io.Writer
}

func (s *IOStreams) Fprintf(format string, args ...any)
func (s *IOStreams) Fprintln(args ...any)
func (s *IOStreams) Errorf(format string, args ...any)
func (s *IOStreams) Errorln(args ...any)
```

### Command Interface (P3)

Already specified in constitution:
```go
package cmd

type Command interface {
    Complete() error
    Validate() error
    Run(ctx context.Context) error
    AddFlags(fs *pflag.FlagSet)
}
```

### CommandOptions Struct (P4)

```go
package lint

type CommandOptions struct {
    Shared        *doctor.SharedOptions
    TargetVersion string
}

type CommandOption func(*Command)

func WithShared(shared *doctor.SharedOptions) CommandOption
func WithTargetVersion(version string) CommandOption
```

### Version Detection Pattern (P2)

Checks detect lint vs upgrade mode:
```go
func (c *Check) Run(ctx context.Context, target *check.CheckTarget) check.CheckResult {
    isLintMode := target.Version.Version == target.CurrentVersion.Version
    isUpgradeMode := target.Version.Version != target.CurrentVersion.Version

    if isUpgradeMode {
        return checkUpgradeReadiness(target)
    }
    return checkCurrentState(target)
}
```

**No contracts/, data-model.md, or quickstart.md generated** - refactoring task does not introduce new APIs, data models, or user-facing features requiring quickstart documentation.

## Phase 2: Task Generation

**Status**: ⏳ **PENDING** - Use `/speckit.tasks` command to generate tasks.md

Task generation will create detailed implementation tasks for each user story following the incremental delivery strategy (P1 → P2 → P3 → P4 → P5).

Expected task structure:
- **User Story 1 (P1)**: 3-5 tasks for IOStreams wrapper implementation and command refactoring
- **User Story 2 (P2)**: 5-7 tasks for merging lint/upgrade commands and adding --version flag
- **User Story 3 (P3)**: 4-6 tasks for Command interface definition and lint refactoring
- **User Story 4 (P4)**: 3-5 tasks for struct/functional options implementation
- **User Story 5 (P5)**: 2-4 tasks for deprecated API audit and replacement

Each task will follow constitution commit granularity (T###: description format) and include:
- Clear acceptance criteria
- Dependencies on prior tasks
- Expected test coverage
- Constitution check compliance

## Delivery Milestones

| Milestone | User Story | Deliverable | Success Criteria |
|-----------|------------|-------------|------------------|
| M1 | P1 | IOStreams wrapper + refactored command output | SC-001: 30% code reduction, all tests pass |
| M2 | P2 | Unified lint command with --version flag | SC-002, SC-003, SC-004: lint and upgrade modes work correctly |
| M3 | P3 | Command interface + lint refactoring | SC-005: Command interface implemented, AddFlags tested |
| M4 | P4 | Flexible initialization patterns | SC-006: Both struct and functional options work identically |
| M5 | P5 | Deprecated API audit complete | SC-007: 100% deprecated APIs identified and replaced/documented |
| **Final** | All | Constitution v1.15.0 compliance | SC-008, SC-009, SC-010, SC-011: All quality gates pass, tests pass, incremental delivery complete |

## Next Steps

1. ✅ **Plan Complete**: Review this implementation plan
2. ⏳ **Task Generation**: Run `/speckit.tasks` to create detailed tasks.md with granular implementation steps
3. ⏳ **Implementation**: Execute tasks in user story priority order (P1 → P2 → P3 → P4 → P5)
4. ⏳ **Validation**: Verify each user story against success criteria before proceeding to next
5. ⏳ **Final Review**: Ensure all 11 success criteria met and constitutional compliance verified