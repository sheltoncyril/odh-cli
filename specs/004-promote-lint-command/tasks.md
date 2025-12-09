# Tasks: Promote Lint Command to Top Level

**Input**: Design documents from `/specs/004-promote-lint-command/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/command-api.md

**Tests**: Test tasks are included as this is a refactoring that requires regression testing to ensure behavioral equivalence.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)
- Include exact file paths in descriptions

## Path Conventions

- **Go CLI project**: `cmd/`, `pkg/cmd/`, `pkg/` at repository root
- Command registration: `cmd/<command>.go`
- Command implementation: `pkg/cmd/<command>/<command>.go`
- Domain logic: `pkg/<domain>/`

## Phase 1: Setup (Package Structure Preparation)

**Purpose**: Prepare new package structure without breaking existing code

- [X] T001 Create new package directory structure at pkg/cmd/lint/
- [X] T002 Create new package directory structure at pkg/lint/
- [X] T003 [P] Create cmd/lint.go placeholder file for top-level command registration
- [X] T004 [P] Verify no naming conflicts with existing packages

---

## Phase 2: Foundational (Core Package Migration)

**Purpose**: Move domain packages (check framework, version detection, checks) that ALL user stories depend on

**⚠️ CRITICAL**: User story implementation cannot begin until domain packages are relocated

- [X] T005 Move pkg/doctor/check/ to pkg/lint/check/ preserving git history
- [X] T006 Move pkg/doctor/version/ to pkg/lint/version/ preserving git history
- [X] T007 Move pkg/doctor/checks/ to pkg/lint/checks/ preserving git history
- [X] T008 Update package declarations in pkg/lint/check/*.go from "package check" (if nested) to correct package name
- [X] T009 Update package declarations in pkg/lint/version/*.go from "package version" (if nested) to correct package name
- [X] T010 Update package declarations in pkg/lint/checks/**/*.go to maintain correct package names
- [X] T011 Run goimports -w . to auto-fix import paths for moved packages
- [X] T012 Run make lint-fix to apply automated fixes
- [X] T013 Run make check to verify no compilation errors after package moves

**Checkpoint**: Foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - Direct Lint Command Access (Priority: P1) 🎯 MVP

**Goal**: Users can execute `kubectl odh lint` directly at top level, removing the doctor intermediary command

**Independent Test**: Execute `kubectl odh lint` and verify it validates the cluster with identical results to old `kubectl odh doctor lint`

### Implementation for User Story 1

- [X] T014 [P] [US1] Move pkg/cmd/doctor/lint/lint.go to pkg/cmd/lint/lint.go preserving git history
- [X] T015 [P] [US1] Move pkg/cmd/doctor/lint/lint_options.go to pkg/cmd/lint/lint_options.go preserving git history
- [X] T016 [P] [US1] Move pkg/cmd/doctor/lint/lint_test.go to pkg/cmd/lint/lint_test.go preserving git history
- [X] T017 [P] [US1] Move pkg/cmd/doctor/shared_options.go to pkg/cmd/lint/shared_options.go preserving git history
- [X] T018 [P] [US1] Move pkg/cmd/doctor/shared_options_test.go to pkg/cmd/lint/shared_options_test.go preserving git history
- [X] T019 [US1] Update package declaration in pkg/cmd/lint/lint.go from "package lint" to ensure correct package name
- [X] T020 [US1] Update package declaration in pkg/cmd/lint/lint_options.go to match package lint
- [X] T021 [US1] Update package declaration in pkg/cmd/lint/shared_options.go to match package lint
- [X] T022 [US1] Update import paths in pkg/cmd/lint/lint.go to reference pkg/lint/check, pkg/lint/version, pkg/lint/checks
- [X] T023 [US1] Update import paths in pkg/cmd/lint/lint_test.go to reference pkg/lint packages
- [X] T024 [US1] Update import paths in pkg/cmd/lint/shared_options_test.go to reference pkg/lint packages
- [X] T025 [US1] Implement top-level lint command registration in cmd/lint.go (create AddLintCommand function)
- [X] T026 [US1] Update cmd/root.go to call AddLintCommand directly instead of doctor.AddCommand
- [X] T027 [US1] Update help text constants in cmd/lint.go to use "kubectl odh lint" instead of "kubectl odh doctor lint"
- [X] T028 [US1] Update examples in cmd/lint.go to show "kubectl odh lint" command usage
- [X] T029 [US1] Update lintCmdLong description in cmd/lint.go to remove doctor command references
- [X] T030 [US1] Remove cmd/doctor/doctor.go file entirely
- [X] T031 [US1] Remove cmd/doctor/lint.go file entirely
- [X] T032 [US1] Remove cmd/doctor/ directory
- [X] T033 [US1] Run goimports -w . to fix any remaining import path issues
- [X] T034 [US1] Run make lint-fix to apply automated formatting and fixes
- [X] T035 [US1] Run make check to verify compilation, linting, and tests pass

### Tests for User Story 1

- [X] T036 [US1] Verify kubectl odh lint executes successfully (integration test)
- [X] T037 [US1] Verify kubectl odh lint --help displays correct help text with new command path
- [X] T038 [US1] Verify kubectl odh doctor lint returns command not found error
- [X] T039 [US1] Verify kubectl odh doctor returns command not found error
- [X] T040 [US1] Verify all existing lint functionality works identically (output formats, check execution)
- [X] T041 [US1] Verify command execution time is within 5% of baseline (performance test)

**Checkpoint**: At this point, `kubectl odh lint` should be fully functional, and `kubectl odh doctor lint` should be removed

---

## Phase 4: User Story 2 - Upgrade Assessment with Clear Flag Name (Priority: P2)

**Goal**: Users can assess upgrade readiness using `--target-version` flag instead of ambiguous `--version`

**Independent Test**: Run `kubectl odh lint --target-version 3.0` and verify it assesses upgrade readiness correctly

### Implementation for User Story 2

- [X] T042 [US2] Update AddFlags method in pkg/cmd/lint/lint.go to register --target-version flag instead of --version
- [X] T043 [US2] Update targetVersion field references in pkg/cmd/lint/lint.go to use new flag name
- [X] T044 [US2] Update CommandOptions struct in pkg/cmd/lint/lint_options.go to rename Version field to TargetVersion
- [X] T045 [US2] Update all references to opts.Version in pkg/cmd/lint/lint.go to opts.TargetVersion
- [X] T046 [US2] Update help text in cmd/lint.go to reference --target-version in examples
- [X] T047 [US2] Update lintCmdExample constant in cmd/lint.go to show --target-version usage
- [X] T048 [US2] Update error messages in pkg/cmd/lint/lint.go to reference --target-version for clarity
- [X] T049 [US2] Search codebase for any remaining "--version" string references in lint command context and update
- [X] T050 [US2] Run make lint-fix to apply automated fixes
- [X] T051 [US2] Run make check to verify compilation and tests pass

### Tests for User Story 2

- [X] T052 [US2] Verify kubectl odh lint --target-version 3.0 executes upgrade assessment correctly
- [X] T053 [US2] Verify kubectl odh lint --target-version invalid displays clear error message
- [X] T054 [US2] Verify kubectl odh lint --version 3.0 returns unknown flag error
- [X] T055 [US2] Verify --target-version flag parsing works correctly in unit tests (update pkg/cmd/lint/lint_test.go)
- [X] T056 [US2] Verify upgrade mode behavior is identical to old --version flag (regression test)

**Checkpoint**: Both user stories should now be fully functional with new command path and flag name

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Final cleanup and verification across all changes

- [X] T057 [P] Run grep -r "pkg/cmd/doctor" to verify no remaining import references to old package
- [X] T058 [P] Run grep -r "pkg/doctor" to verify no remaining import references to old domain package
- [X] T059 [P] Run grep -r "kubectl odh doctor lint" in codebase to find any remaining command references
- [X] T060 [P] Run grep -r '"doctor"' in cmd/ pkg/ to find any remaining string references
- [X] T061 [P] Verify go.mod has no stale references (run go mod tidy)
- [X] T062 [P] Update constitution reference in code comments if any reference old command structure
- [X] T063 Run make check to verify all linting, testing, and vulnerability checks pass
- [X] T064 Run full integration test suite to verify no regressions
- [X] T065 Manually test kubectl odh lint with various flag combinations
- [X] T066 Manually test kubectl odh lint --target-version with various versions
- [X] T067 Verify help text accuracy by running kubectl odh lint --help
- [X] T068 Verify error messages are clear for removed commands/flags
- [X] T069 Performance benchmark comparison against baseline (verify SC-006: within 5%)
- [X] T070 Final code review for any missed references or cleanup opportunities

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup (Phase 1) - BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational (Phase 2) completion
- **User Story 2 (Phase 4)**: Depends on User Story 1 (Phase 3) completion - flag rename requires command promotion first
- **Polish (Phase 5)**: Depends on both user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Depends on Foundational phase - promotes lint to top level
- **User Story 2 (P2)**: Depends on User Story 1 - renames flag in already-promoted command

### Within Each User Story

**User Story 1**:
- File moves (T014-T018) can happen in parallel [P]
- Package declaration updates (T019-T021) must wait for moves
- Import path updates (T022-T024) must wait for package declarations
- Command registration (T025-T032) can happen after import fixes
- Cleanup and verification (T033-T035) must be sequential
- Tests (T036-T041) validate the completed story

**User Story 2**:
- Flag rename tasks (T042-T049) can mostly happen in parallel within same story
- Cleanup (T050-T051) must be sequential
- Tests (T052-T056) validate the flag rename

### Parallel Opportunities

**Phase 1** (Setup):
- T003 and T004 can run in parallel [P]

**Phase 2** (Foundational):
- T005, T006, T007 (package moves) must be sequential to avoid conflicts
- T008, T009, T010 (package declarations) can happen in parallel after moves

**Phase 3** (User Story 1):
- T014, T015, T016, T017, T018 (file moves) can run in parallel [P]

**Phase 5** (Polish):
- T057, T058, T059, T060, T061, T062 (grep searches) can all run in parallel [P]

---

## Parallel Example: User Story 1 File Moves

```bash
# Launch all file move operations together (they touch different files):
Task: "Move pkg/cmd/doctor/lint/lint.go to pkg/cmd/lint/lint.go"
Task: "Move pkg/cmd/doctor/lint/lint_options.go to pkg/cmd/lint/lint_options.go"
Task: "Move pkg/cmd/doctor/lint/lint_test.go to pkg/cmd/lint/lint_test.go"
Task: "Move pkg/cmd/doctor/shared_options.go to pkg/cmd/lint/shared_options.go"
Task: "Move pkg/cmd/doctor/shared_options_test.go to pkg/cmd/lint/shared_options_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup → Package structure ready
2. Complete Phase 2: Foundational → Domain packages moved
3. Complete Phase 3: User Story 1 → Command promoted to top level
4. **STOP and VALIDATE**: Test `kubectl odh lint` independently
5. Verify old command path removed
6. Deploy/demo if ready (breaking change requires communication)

### Incremental Delivery

1. Complete Setup + Foundational → Package structure ready
2. Add User Story 1 → Test independently → Command accessible at top level (MVP!)
3. Add User Story 2 → Test independently → Flag name clarified
4. Polish → Final verification and cleanup

### Sequential Strategy (Recommended)

Since User Story 2 depends on User Story 1 (flag exists in promoted command):

1. Complete Setup (Phase 1)
2. Complete Foundational (Phase 2)
3. Complete User Story 1 (Phase 3) → Validate independently
4. Complete User Story 2 (Phase 4) → Validate flag rename
5. Complete Polish (Phase 5) → Final verification

---

## Constitutional Compliance Checkpoints

After each phase, verify constitutional requirements:

**After User Story 1 (T035, T041)**:
- ✅ Command implements Complete/Validate/Run/AddFlags (Principle II)
- ✅ Package structure follows pkg/cmd/lint/ pattern (Code Organization)
- ✅ All output formats preserved (Principle III)
- ✅ make check passes (Quality Gates)

**After User Story 2 (T051, T056)**:
- ✅ Flag name clarity improved (Lint Command Architecture)
- ✅ Help text updated (FR-006)
- ✅ Error messages clear (FR-009)

**After Polish (T063-T070)**:
- ✅ No doctor command references remain
- ✅ All imports updated correctly
- ✅ Performance within 5% baseline (SC-006)
- ✅ Full regression test suite passes

---

## Notes

- [P] tasks = different files or independent operations, no dependencies
- [Story] label maps task to specific user story (US1 or US2)
- Each user story builds on previous (US2 depends on US1)
- Use `git mv` for all file moves to preserve history (T005-T007, T014-T018)
- Run `make lint-fix` after each major change group to auto-fix formatting
- Run `make check` frequently to catch issues early
- Commit after each task or logical group with format: `T###: Description`
- This is a breaking change - coordinate with users before deployment
- All existing lint validation logic remains unchanged (FR-003)
- Zero changes to check implementation or execution (Constraint from plan.md)

---

## Summary

**Total Tasks**: 70
**User Story 1 Tasks**: 28 (T014-T041)
**User Story 2 Tasks**: 15 (T042-T056)
**Setup Tasks**: 4 (T001-T004)
**Foundational Tasks**: 9 (T005-T013)
**Polish Tasks**: 14 (T057-T070)

**Parallel Opportunities**:
- Setup: 2 parallel tasks
- Foundational: 3 parallel groups (moves, then declarations)
- User Story 1: 5 file moves in parallel, 6 grep searches in polish
- User Story 2: Most flag rename tasks can be parallelized

**MVP Scope**: Phase 1 + Phase 2 + Phase 3 (User Story 1 only)
- Delivers: `kubectl odh lint` command at top level
- Removes: `kubectl odh doctor` parent command
- Validates: Identical functionality to old command path

**Full Feature**: All phases
- Adds: Clear `--target-version` flag name
- Removes: Ambiguous `--version` flag
- Validates: Complete migration with no regressions
