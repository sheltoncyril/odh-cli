# Tasks: Constitution v1.15.0 Alignment

**Input**: Design documents from `/home/luca/work/dev/openshift-ai/odh-cli/specs/003-constitution-alignment/`
**Prerequisites**: plan.md (required), spec.md (required for user stories)

**Tests**: Test tasks are included per constitutional requirement (Principle VI: Test-First Development)

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story in priority order (P1 → P2 → P3 → P4 → P5).

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

This is a CLI tool refactoring with single project structure:
- Command wrappers: `cmd/doctor/`
- Business logic: `pkg/cmd/doctor/`
- Utilities: `pkg/util/`
- Tests: Co-located with implementation files

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Minimal setup for refactoring task (no new project initialization needed)

- [x] T001 Run `make lint-fix` to auto-fix existing trivial linting issues before refactoring
- [x] T002 Verify all existing tests pass with `make test` to establish baseline
- [x] T003 Verify all quality gates pass with `make check` to establish baseline

**Checkpoint**: Baseline established - ready to begin incremental refactoring

---

## Phase 2: User Story 1 - IOStreams Wrapper (Priority: P1) 🎯 MVP

**Goal**: Create centralized IOStreams wrapper in `pkg/util/iostreams/` to eliminate repetitive `fmt.Fprintf(o.Out, ...)` patterns and reduce output-related code by 30%

**Independent Test**: Implement IOStreams wrapper and refactor lint command to use it. Verify output behavior identical before/after and measure code reduction in output-related lines.

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T004 [P] [US1] Create unit tests for IOStreams.Fprintf in pkg/util/iostreams/iostreams_test.go (test formatting with and without args)
- [ ] T005 [P] [US1] Create unit tests for IOStreams.Fprintln in pkg/util/iostreams/iostreams_test.go (test plain output)
- [ ] T006 [P] [US1] Create unit tests for IOStreams.Errorf in pkg/util/iostreams/iostreams_test.go (test error output to stderr)
- [ ] T007 [P] [US1] Create unit tests for IOStreams.Errorln in pkg/util/iostreams/iostreams_test.go (test plain error output)
- [ ] T008 [US1] Create unit tests for nil writer validation in pkg/util/iostreams/iostreams_test.go (edge case from spec)

### Implementation for User Story 1

- [ ] T009 [US1] Create pkg/util/iostreams/ directory
- [ ] T010 [US1] Implement IOStreams struct with In, Out, ErrOut fields in pkg/util/iostreams/iostreams.go
- [ ] T011 [P] [US1] Implement Fprintf method (auto-newline, conditional fmt.Sprintf) in pkg/util/iostreams/iostreams.go
- [ ] T012 [P] [US1] Implement Fprintln method (direct output) in pkg/util/iostreams/iostreams.go
- [ ] T013 [P] [US1] Implement Errorf method (stderr, auto-newline, conditional fmt.Sprintf) in pkg/util/iostreams/iostreams.go
- [ ] T014 [P] [US1] Implement Errorln method (stderr, direct output) in pkg/util/iostreams/iostreams.go
- [ ] T015 [US1] Update SharedOptions in pkg/cmd/doctor/shared_options.go to embed IOStreams wrapper
- [ ] T016 [US1] Refactor lint command output calls in pkg/cmd/doctor/lint/options.go to use IOStreams wrapper instead of fmt.Fprintf
- [ ] T017 [US1] Measure code reduction: count output-related lines before/after refactoring (target: ≥30% reduction per SC-001)
- [ ] T018 [US1] Run existing lint tests to verify output behavior unchanged
- [ ] T019 [US1] Run `make lint-fix` to auto-fix any formatting issues
- [ ] T020 [US1] Run `make check` to verify quality gates pass
- [ ] T021 [US1] Commit with message: "feat(US1): implement IOStreams wrapper and refactor lint command output\n\nT004-T020: IOStreams wrapper with Fprintf/Fprintln/Errorf/Errorln\n- Created pkg/util/iostreams/ with automatic newline handling\n- Refactored lint command to use wrapper, reducing output code by X%\n- All tests passing, output behavior preserved\n\n🤖 Generated with Claude Code\nCo-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"

**Checkpoint**: User Story 1 complete. IOStreams wrapper functional, lint command refactored, 30% code reduction achieved, all tests passing.

---

## Phase 3: User Story 2 - Unified Lint Command (Priority: P2)

**Goal**: Merge separate lint and upgrade commands into single `lint` command with optional `--version` flag. Remove upgrade subcommand entirely.

**Independent Test**: Run `kubectl odh doctor lint` (lint mode) and `kubectl odh doctor lint --version 3.0` (upgrade mode). Verify both modes work correctly with same output format.

### Tests for User Story 2

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T022 [P] [US2] Create unit test for lint mode (no --version flag) in pkg/cmd/doctor/lint/lint_test.go
- [ ] T023 [P] [US2] Create unit test for upgrade mode (with --version flag) in pkg/cmd/doctor/lint/lint_test.go
- [ ] T024 [P] [US2] Create unit test verifying CheckTarget.CurrentVersion == CheckTarget.Version in lint mode in pkg/cmd/doctor/lint/lint_test.go
- [ ] T025 [P] [US2] Create unit test verifying CheckTarget.CurrentVersion != CheckTarget.Version in upgrade mode in pkg/cmd/doctor/lint/lint_test.go
- [ ] T026 [US2] Create integration test for both lint and upgrade modes in pkg/cmd/doctor/lint/lint_test.go

### Implementation for User Story 2

- [ ] T027 [US2] Copy upgrade command tests from pkg/cmd/doctor/upgrade/ to pkg/cmd/doctor/lint/ for preservation
- [ ] T028 [US2] Add --version flag to lint command in pkg/cmd/doctor/lint/options.go (string field for target version)
- [ ] T029 [US2] Update lint.Complete() to set target version from flag or default to current version in pkg/cmd/doctor/lint/options.go
- [ ] T030 [US2] Update lint.Validate() to accept optional version in pkg/cmd/doctor/lint/options.go
- [ ] T031 [US2] Merge upgrade command logic into lint.Run() with version-based mode detection in pkg/cmd/doctor/lint/options.go
- [ ] T032 [US2] Update cmd/doctor/lint.go Cobra wrapper to register --version flag
- [ ] T033 [US2] Update check executor to populate CheckTarget with both CurrentVersion and Version in pkg/doctor/check/executor.go
- [ ] T034 [US2] Delete cmd/doctor/upgrade.go (Cobra wrapper for upgrade command)
- [ ] T035 [US2] Delete pkg/cmd/doctor/upgrade/ directory entirely
- [ ] T036 [US2] Run all lint tests (original + merged upgrade tests) to verify both modes work
- [ ] T037 [US2] Run manual integration test: `kubectl odh doctor lint` (should validate current state)
- [ ] T038 [US2] Run manual integration test: `kubectl odh doctor lint --version 3.0` (should assess upgrade readiness)
- [ ] T039 [US2] Run `make lint-fix` to auto-fix any formatting issues
- [ ] T040 [US2] Run `make check` to verify quality gates pass
- [ ] T041 [US2] Commit with message: "feat(US2): merge lint and upgrade commands into unified lint command\n\nT022-T040: Unified lint command with optional --version flag\n- Added --version flag to lint command for upgrade-readiness mode\n- Merged upgrade command logic into lint.Run() with mode detection\n- Deleted separate upgrade subcommand (cmd/doctor/upgrade.go, pkg/cmd/doctor/upgrade/)\n- Both modes tested and working correctly\n\n🤖 Generated with Claude Code\nCo-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"

**Checkpoint**: User Story 2 complete. Lint and upgrade merged, --version flag working, upgrade command removed, all tests passing.

---

## Phase 4: User Story 3 - Command Interface (Priority: P3)

**Goal**: Define standard Command interface with Complete/Validate/Run/AddFlags methods. Refactor lint command to implement it. Rename Options → Command, options.go → lint.go.

**Independent Test**: Verify lint command implements Command interface, has AddFlags method for flag registration, and passes all existing tests with new structure.

### Tests for User Story 3

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T042 [P] [US3] Create unit test for AddFlags method in pkg/cmd/doctor/lint/lint_test.go (verify flags registered correctly)
- [ ] T043 [P] [US3] Create unit test verifying lint.Command implements cmd.Command interface in pkg/cmd/doctor/lint/lint_test.go
- [ ] T044 [US3] Create unit test for NewCommand() constructor initialization in pkg/cmd/doctor/lint/lint_test.go

### Implementation for User Story 3

- [ ] T045 [US3] Define Command interface in pkg/cmd/command.go with Complete() error, Validate() error, Run(ctx context.Context) error, AddFlags(fs *pflag.FlagSet) methods
- [ ] T046 [US3] Rename pkg/cmd/doctor/lint/options.go to pkg/cmd/doctor/lint/lint.go (per constitution Principle II)
- [ ] T047 [US3] Rename Options struct to Command struct in pkg/cmd/doctor/lint/lint.go (per FR-011)
- [ ] T048 [US3] Rename NewOptions() function to NewCommand() in pkg/cmd/doctor/lint/lint.go (per FR-012)
- [ ] T049 [US3] Implement AddFlags method on lint.Command to register flags with pflag.FlagSet in pkg/cmd/doctor/lint/lint.go
- [ ] T050 [US3] Update cmd/doctor/lint.go to call command.AddFlags(cmd.Flags()) instead of inline flag registration
- [ ] T051 [US3] Verify lint.Command satisfies cmd.Command interface (add compile-time check)
- [ ] T052 [US3] Update all references from Options to Command in pkg/cmd/doctor/lint/lint.go
- [ ] T053 [US3] Update SharedOptions initialization to be done internally by Command per FR-014 in pkg/cmd/doctor/lint/lint.go
- [ ] T054 [US3] Run all lint tests to verify refactored structure works correctly
- [ ] T055 [US3] Run `make lint-fix` to auto-fix any formatting issues
- [ ] T056 [US3] Run `make check` to verify quality gates pass
- [ ] T057 [US3] Commit with message: "feat(US3): introduce Command interface and refactor lint command structure\n\nT042-T056: Command interface standardization\n- Defined cmd.Command interface with Complete/Validate/Run/AddFlags\n- Renamed lint.Options → lint.Command, NewOptions → NewCommand\n- Renamed options.go → lint.go per constitutional naming convention\n- Implemented AddFlags method for centralized flag registration\n- Command initializes SharedOptions internally\n\n🤖 Generated with Claude Code\nCo-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"

**Checkpoint**: User Story 3 complete. Command interface defined, lint command refactored, AddFlags tested, all tests passing.

---

## Phase 5: User Story 4 - Flexible Initialization (Priority: P4)

**Goal**: Support both struct-based initialization (preferred) and functional options patterns for command initialization. Create CommandOptions struct and With* functions.

**Independent Test**: Create command instances using both struct initialization and functional options. Verify both produce identical behavior.

### Tests for User Story 4

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T058 [P] [US4] Create unit test for struct-based initialization in pkg/cmd/doctor/lint/lint_test.go
- [ ] T059 [P] [US4] Create unit test for functional options initialization in pkg/cmd/doctor/lint/lint_test.go
- [ ] T060 [US4] Create unit test verifying both initialization patterns produce identical Command state in pkg/cmd/doctor/lint/lint_test.go

### Implementation for User Story 4

- [ ] T061 [US4] Create pkg/cmd/doctor/lint/lint_options.go file (per convention for options code)
- [ ] T062 [P] [US4] Define CommandOptions struct with Shared and TargetVersion fields in pkg/cmd/doctor/lint/lint_options.go
- [ ] T063 [P] [US4] Define CommandOption type as func(*Command) in pkg/cmd/doctor/lint/lint_options.go
- [ ] T064 [P] [US4] Implement WithShared(shared *doctor.SharedOptions) CommandOption function in pkg/cmd/doctor/lint/lint_options.go
- [ ] T065 [P] [US4] Implement WithTargetVersion(version string) CommandOption function in pkg/cmd/doctor/lint/lint_options.go
- [ ] T066 [US4] Update NewCommand() to accept CommandOptions struct as primary parameter in pkg/cmd/doctor/lint/lint.go
- [ ] T067 [US4] Create NewCommandWithOptions(...CommandOption) constructor for functional options pattern in pkg/cmd/doctor/lint/lint.go
- [ ] T068 [US4] Update cmd/doctor/lint.go to use struct-based initialization (preferred pattern)
- [ ] T069 [US4] Run all lint tests with both initialization patterns
- [ ] T070 [US4] Run `make lint-fix` to auto-fix any formatting issues
- [ ] T071 [US4] Run `make check` to verify quality gates pass
- [ ] T072 [US4] Commit with message: "feat(US4): add flexible initialization patterns for command construction\n\nT058-T071: Struct-based and functional options support\n- Created CommandOptions struct for simple initialization (preferred)\n- Implemented With* functional option functions for complex cases\n- Both patterns tested and produce identical behavior (SC-006)\n- Options code in lint_options.go per convention\n\n🤖 Generated with Claude Code\nCo-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"

**Checkpoint**: User Story 4 complete. Both initialization patterns working, tests verify identical behavior, all tests passing.

---

## Phase 6: User Story 5 - Deprecated API Audit (Priority: P5)

**Goal**: Identify all deprecated API usage throughout codebase using IDE warnings and godoc markers. Replace deprecated APIs with modern alternatives where available. Document unavoidable deprecated usage with explanatory comments.

**Independent Test**: Audit complete, no deprecation warnings remain (or documented with tracking issues). All tests pass after replacements.

### Implementation for User Story 5

- [ ] T073 [US5] Run codebase audit for deprecated API usage using IDE warnings in GoLand/VSCode
- [ ] T074 [US5] Run codebase audit for deprecated API usage using godoc markers (search for "// Deprecated:")
- [ ] T075 [US5] Create deprecation audit report listing all identified deprecated APIs with locations and replacement options
- [ ] T076 [US5] Replace deprecated API #1 with modern alternative (if available) - update specific file based on audit
- [ ] T077 [US5] Run tests after replacement #1 to verify no breakage
- [ ] T078 [US5] Replace deprecated API #2 with modern alternative (if available) - update specific file based on audit
- [ ] T079 [US5] Run tests after replacement #2 to verify no breakage
- [ ] T080 [US5] Replace deprecated API #3 with modern alternative (if available) - update specific file based on audit
- [ ] T081 [US5] Run tests after replacement #3 to verify no breakage
- [ ] T082 [US5] For unavoidable deprecated APIs, add explanatory comments with tracking issue references per FR-019
- [ ] T083 [US5] Verify zero deprecation warnings remain (or all documented with justifications)
- [ ] T084 [US5] Run `make lint-fix` to auto-fix any formatting issues
- [ ] T085 [US5] Run `make check` to verify quality gates pass
- [ ] T086 [US5] Commit with message: "feat(US5): audit and replace deprecated API usage across codebase\n\nT073-T085: Deprecated API cleanup\n- Audited codebase using IDE warnings and godoc markers\n- Replaced X deprecated APIs with modern alternatives\n- Documented Y unavoidable deprecated APIs with explanatory comments\n- 100% deprecated API coverage achieved (SC-007)\n\n🤖 Generated with Claude Code\nCo-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"

**Checkpoint**: User Story 5 complete. Deprecated API audit done, replacements made, unavoidable usage documented, all tests passing.

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Final improvements and validation affecting multiple user stories

- [ ] T087 [P] Update CLAUDE.md with constitution v1.15.0 changes (IOStreams wrapper, Command interface, merged lint command)
- [ ] T088 [P] Update constitution.md Follow-up TODOs section to mark all items as completed
- [ ] T089 Verify all 11 success criteria met (SC-001 through SC-011) from spec.md
- [ ] T090 Verify all constitutional gates pass from plan.md (Phase 2 Implementation checklist)
- [ ] T091 Run full test suite with `make test` to ensure all tests pass across all user stories
- [ ] T092 Run quality gates with `make check` to ensure zero linting issues
- [ ] T093 Verify incremental delivery: 5 independent commits (one per user story) with all tests passing
- [ ] T094 Final review: Ensure no implementation details leaked into CLI behavior (SC-010)

**Checkpoint**: Constitution v1.15.0 alignment complete. All user stories delivered incrementally, all success criteria met, all tests passing.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **User Story 1 (Phase 2)**: Depends on Setup completion
- **User Story 2 (Phase 3)**: Depends on User Story 1 completion (can optionally use IOStreams wrapper)
- **User Story 3 (Phase 4)**: Depends on User Story 2 completion (requires merged lint command structure)
- **User Story 4 (Phase 5)**: Depends on User Story 3 completion (requires Command interface)
- **User Story 5 (Phase 6)**: Can start after User Story 1-4 complete (independent audit task)
- **Polish (Phase 7)**: Depends on all user stories being complete

### User Story Dependencies

```text
Setup (Phase 1)
  ↓
User Story 1 (P1) - IOStreams Wrapper [No dependencies on other stories]
  ↓
User Story 2 (P2) - Unified Lint Command [Can use IOStreams from P1]
  ↓
User Story 3 (P3) - Command Interface [Requires P2 lint command structure]
  ↓
User Story 4 (P4) - Flexible Initialization [Requires P3 Command interface]
  ↓
User Story 5 (P5) - Deprecated API Audit [Independent after P1-P4]
  ↓
Polish (Phase 7)
```

### Within Each User Story

- Tests MUST be written and FAIL before implementation (test-first per constitution)
- Implementation tasks follow test tasks
- Quality gates (`make lint-fix`, `make check`) run after implementation
- Commit completes user story as independent increment
- Story complete before moving to next priority

### Parallel Opportunities

**Within User Story 1 (IOStreams Wrapper)**:
- T004-T008: All unit tests can run in parallel
- T011-T014: All IOStreams methods can be implemented in parallel

**Within User Story 2 (Unified Lint Command)**:
- T022-T026: All unit/integration tests can run in parallel
- T027-T033: Sequential (each builds on previous)

**Within User Story 3 (Command Interface)**:
- T042-T044: All unit tests can run in parallel
- T045-T053: Sequential refactoring steps

**Within User Story 4 (Flexible Initialization)**:
- T058-T060: All unit tests can run in parallel
- T062-T065: All option definitions can be implemented in parallel

**Within User Story 5 (Deprecated API Audit)**:
- T073-T075: Audit tasks sequential (audit → report)
- T076-T081: Replacements sequential (one at a time with test validation)

**Within Polish Phase**:
- T087-T088: Documentation updates can run in parallel

---

## Parallel Example: User Story 1 (IOStreams Wrapper)

```bash
# Launch all tests for User Story 1 together:
Task: "Create unit tests for IOStreams.Fprintf in pkg/util/iostreams/iostreams_test.go"
Task: "Create unit tests for IOStreams.Fprintln in pkg/util/iostreams/iostreams_test.go"
Task: "Create unit tests for IOStreams.Errorf in pkg/util/iostreams/iostreams_test.go"
Task: "Create unit tests for IOStreams.Errorln in pkg/util/iostreams/iostreams_test.go"

# Launch all method implementations together:
Task: "Implement Fprintf method in pkg/util/iostreams/iostreams.go"
Task: "Implement Fprintln method in pkg/util/iostreams/iostreams.go"
Task: "Implement Errorf method in pkg/util/iostreams/iostreams.go"
Task: "Implement Errorln method in pkg/util/iostreams/iostreams.go"
```

---

## Implementation Strategy

### Incremental Delivery (Required)

Per clarification session 2025-12-08 and FR-021, implementation MUST follow user story priority order with independent commits:

1. **Complete User Story 1 (P1)** → Commit → Verify tests pass → Verify `make check` passes
2. **Complete User Story 2 (P2)** → Commit → Verify tests pass → Verify `make check` passes
3. **Complete User Story 3 (P3)** → Commit → Verify tests pass → Verify `make check` passes
4. **Complete User Story 4 (P4)** → Commit → Verify tests pass → Verify `make check` passes
5. **Complete User Story 5 (P5)** → Commit → Verify tests pass → Verify `make check` passes
6. **Complete Polish** → Final validation → All 11 success criteria verified

Each commit must:
- Leave codebase in working state
- Pass all quality gates (`make check`)
- Follow constitution commit message format with task IDs
- Demonstrate measurable progress toward success criteria

### Backward Compatibility

**NOT a concern** - internal API changes are acceptable between user stories. Focus on:
- All tests passing after each story
- CLI behavior preserved (no breaking changes to user-facing commands)
- Constitutional compliance maintained

---

## Notes

- **[P]** tasks = different files, no dependencies, can run in parallel
- **[Story]** label maps task to specific user story for traceability (US1-US5)
- Each user story should be independently completable and testable
- Verify tests fail before implementing (test-first development per constitution)
- Commit after each user story completion (incremental delivery per clarification)
- Run `make lint-fix` before `make check` per constitutional workflow
- Stop at each checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence

---

## Success Criteria Mapping

| Success Criterion | Verified By | Task(s) |
|-------------------|-------------|---------|
| SC-001: 30% code reduction in output lines | User Story 1 | T017 |
| SC-002: `kubectl odh doctor lint` validates current state | User Story 2 | T037 |
| SC-003: `kubectl odh doctor lint --version` assesses upgrade | User Story 2 | T038 |
| SC-004: Checks adapt behavior based on version | User Story 2 | T024, T025 |
| SC-005: Commands implement Command interface with AddFlags | User Story 3 | T051, T054 |
| SC-006: Both initialization patterns work identically | User Story 4 | T060, T069 |
| SC-007: 100% deprecated API audit coverage | User Story 5 | T083 |
| SC-008: All quality gates pass | All Stories | T020, T040, T056, T071, T085, T092 |
| SC-009: Existing tests pass | All Stories | T018, T036, T054, T069, T091 |
| SC-010: No implementation details leak to CLI | Polish | T094 |
| SC-011: Each story delivered as independent commit | All Stories | T021, T041, T057, T072, T086 |

---

## Task Count Summary

- **Total Tasks**: 94
- **Phase 1 (Setup)**: 3 tasks
- **Phase 2 (User Story 1 - IOStreams)**: 18 tasks (5 test, 13 implementation)
- **Phase 3 (User Story 2 - Unified Lint)**: 20 tasks (5 test, 15 implementation)
- **Phase 4 (User Story 3 - Command Interface)**: 16 tasks (3 test, 13 implementation)
- **Phase 5 (User Story 4 - Flexible Init)**: 15 tasks (3 test, 12 implementation)
- **Phase 6 (User Story 5 - Deprecated API)**: 14 tasks (0 test, 14 implementation)
- **Phase 7 (Polish)**: 8 tasks

**Parallel Opportunities**: 25 tasks marked [P] can be executed in parallel within their respective phases.

**MVP Scope**: User Story 1 only (IOStreams wrapper) provides immediate value with 30% code reduction.