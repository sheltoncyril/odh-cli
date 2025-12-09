# Feature Specification: Promote Lint Command to Top Level

**Feature Branch**: `004-promote-lint-command`
**Created**: 2025-12-09
**Status**: Draft
**Input**: User description: "remove the doctor subcommand and move the lint one level above so one can run kubectl odh lint or kubectl odh lint --version, also rename --version to --target-version"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Direct Lint Command Access (Priority: P1)

Users can execute the lint command directly at the top level without the "doctor" intermediary command, making the CLI more intuitive and reducing typing overhead for the most commonly used diagnostic functionality.

**Why this priority**: This is the core functionality of the feature - simplifying command access. The lint command is the primary diagnostic tool, and users should be able to access it with minimal typing. This provides immediate value and can be tested independently.

**Independent Test**: Can be fully tested by executing `kubectl odh lint` and verifying it validates the cluster, and delivers the same functionality currently available via `kubectl odh doctor lint`.

**Acceptance Scenarios**:

1. **Given** a user wants to validate their cluster, **When** they run `kubectl odh lint`, **Then** the system performs current state validation and displays validation results
2. **Given** a user wants help with the lint command, **When** they run `kubectl odh lint --help`, **Then** the system displays comprehensive help information including usage, flags, and examples
3. **Given** a user tries the old command structure, **When** they run `kubectl odh doctor lint`, **Then** the system displays an error indicating the command is not found

---

### User Story 2 - Upgrade Assessment with Clear Flag Name (Priority: P2)

Users can assess upgrade readiness using a clearly named `--target-version` flag that explicitly communicates its purpose, replacing the ambiguous `--version` flag.

**Why this priority**: This improves user experience by making the command's intent clearer. While important for usability, it depends on P1 being complete (the command being promoted). It can be tested independently by verifying the flag works correctly.

**Independent Test**: Can be fully tested by running `kubectl odh lint --target-version 3.0` and verifying it assesses upgrade readiness to version 3.0, delivering upgrade assessment functionality.

**Acceptance Scenarios**:

1. **Given** a user wants to assess upgrade readiness, **When** they run `kubectl odh lint --target-version 3.0`, **Then** the system analyzes the cluster for upgrade compatibility to version 3.0
2. **Given** a user provides an invalid version, **When** they run `kubectl odh lint --target-version invalid`, **Then** the system displays a clear error message indicating the version format is invalid
3. **Given** a user tries the old flag name, **When** they run `kubectl odh lint --version 3.0`, **Then** the system displays an error indicating the flag is not recognized

---

### Edge Cases

- What happens when a user runs `kubectl odh doctor` (accessing the now-removed parent command)?
- What error message is shown when users try the old `kubectl odh doctor lint` command?
- How are man pages and documentation updated to reflect the new command structure?
- What happens if users specify an empty value for `--target-version`?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST allow users to execute lint validation directly via `kubectl odh lint` without requiring the "doctor" intermediary command
- **FR-002**: System MUST accept `--target-version` flag to specify the target version for upgrade assessment
- **FR-003**: System MUST maintain all existing lint functionality when accessed via the new top-level command
- **FR-004**: System MUST remove the `kubectl odh doctor` command entirely
- **FR-005**: System MUST NOT recognize the `--version` flag in the context of the lint command
- **FR-006**: System MUST update all help text, examples, and command descriptions to reflect the new top-level command structure
- **FR-007**: System MUST preserve all existing output formats (JSON, YAML, text) when using the promoted command
- **FR-008**: System MUST maintain all existing flag combinations and options (except `--version`) with the promoted command
- **FR-009**: System MUST display clear error messages when users attempt to use removed commands or flags

### Key Entities

- **Lint Command**: The primary diagnostic command that validates cluster state or assesses upgrade readiness. Promoted from subcommand to top-level access.
- **Target Version**: The version of OpenShift AI that the user wants to assess upgrade compatibility for, specified via the `--target-version` flag.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can successfully execute lint validation using `kubectl odh lint` with identical results to the previous `kubectl odh doctor lint` command
- **SC-002**: Users can specify target version using `kubectl odh lint --target-version X.Y` and receive upgrade assessment results
- **SC-003**: All existing lint command functionality (output formats, flag combinations, check filtering) works identically with the promoted command
- **SC-004**: The `kubectl odh doctor` command no longer exists and returns a command not found error
- **SC-005**: The `--version` flag is not recognized by the lint command and returns a clear error message
- **SC-006**: Command execution time remains unchanged (within 5%) compared to the previous command structure
- **SC-007**: Documentation and help text accurately reflect the new command structure for all user-facing materials

### Assumptions

- This is a breaking change and users with existing scripts/automation will need to update them
- The `--version` flag removal applies only to the lint command context, not globally to the CLI
- All existing lint command validation logic remains unchanged - only the command routing and flag names are affected
- Error handling and exit codes remain consistent with current behavior
- Migration documentation will be provided separately to help users update their workflows
