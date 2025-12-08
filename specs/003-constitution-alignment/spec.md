# Feature Specification: Constitution v1.15.0 Alignment

**Feature Branch**: `003-constitution-alignment`
**Created**: 2025-12-08
**Status**: Draft
**Input**: User description: "follow-constitution changes"

## Clarifications

### Session 2025-12-08

- Q: Should the refactoring be implemented incrementally (completing one user story at a time with intermediate commits) or as a coordinated change (implementing multiple related stories together before committing)? → A: Incremental by user story - Complete P1 (IOStreams) fully, then P2 (unified lint), then P3 (Command interface), etc. Each story is independently committed and tested. Note: Backward compatibility is not a concern.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - IOStreams Wrapper Reduces Output Boilerplate (Priority: P1)

As a developer maintaining command output code, I want a centralized IOStreams wrapper so that I can write cleaner code without repetitive `fmt.Fprintf(o.Out, ...)` patterns.

**Why this priority**: This is foundational infrastructure that improves code quality across all commands. It provides immediate value by reducing boilerplate and making code more maintainable without changing CLI behavior.

**Independent Test**: Can be fully tested by implementing the IOStreams wrapper in `pkg/util/iostreams/` and refactoring one command (e.g., lint) to use it. Success is verified by comparing output behavior before/after and confirming reduced code duplication.

**Acceptance Scenarios**:

1. **Given** a command needs to write formatted output, **When** the developer uses `io.Fprintf("Message: %s", value)`, **Then** the output appears on stdout with a newline automatically added
2. **Given** a command needs to write error output, **When** the developer uses `io.Errorf("Error: %v", err)`, **Then** the output appears on stderr with a newline automatically added
3. **Given** a message has no format arguments, **When** the developer uses `io.Fprintf("Static message")`, **Then** the output appears correctly without attempting string formatting

---

### User Story 2 - Unified Lint Command Simplifies User Experience (Priority: P2)

As a cluster administrator, I want a single `lint` command that can validate both current state and upgrade readiness so that I don't need to learn separate commands for related diagnostic tasks.

**Why this priority**: This significantly improves user experience by reducing command proliferation and providing a simpler mental model. The `--version` flag intuitively switches context from "lint current state" to "check upgrade readiness."

**Independent Test**: Can be tested by merging the existing lint and upgrade commands into a single lint command with optional `--version` flag. Success is verified by running `kubectl odh doctor lint` (lint mode) and `kubectl odh doctor lint --version 3.0` (upgrade mode) and confirming both behaviors work correctly.

**Acceptance Scenarios**:

1. **Given** a cluster running version 2.17, **When** an administrator runs `kubectl odh doctor lint`, **Then** the command validates the current cluster state without version comparison
2. **Given** a cluster running version 2.17, **When** an administrator runs `kubectl odh doctor lint --version 3.0`, **Then** the command assesses upgrade readiness by comparing current version 2.17 against target version 3.0
3. **Given** a check implementation, **When** it receives a CheckTarget with different CurrentVersion and Version, **Then** it executes upgrade-specific validation logic
4. **Given** a check implementation, **When** it receives a CheckTarget with matching CurrentVersion and Version, **Then** it executes current-state validation logic

---

### User Story 3 - Command Interface Standardizes Command Structure (Priority: P3)

As a developer adding new commands, I want a Command interface with standard methods so that all commands follow a consistent pattern and are easier to test.

**Why this priority**: This establishes architectural consistency but doesn't change user-facing behavior. It's important for long-term maintainability but can be implemented gradually.

**Independent Test**: Can be tested by defining the Command interface and refactoring the lint command to implement it. Success is verified by confirming the command passes all existing tests and the new structure supports easier testing through the AddFlags method.

**Acceptance Scenarios**:

1. **Given** a new command implementation, **When** it implements the Command interface, **Then** it provides Complete(), Validate(), Run(), and AddFlags() methods
2. **Given** a command's flags, **When** defined via AddFlags() method, **Then** they are registered with the pflag.FlagSet and can be tested independently
3. **Given** a command struct, **When** it's constructed via NewCommand(), **Then** it initializes its own SharedOptions internally

---

### User Story 4 - Flexible Initialization Supports Both Simple and Complex Cases (Priority: P4)

As a developer creating command instances, I want both struct-based and functional option patterns available so that I can use simple struct initialization for common cases and functional options for complex scenarios.

**Why this priority**: This provides developer flexibility without forcing complexity on simple cases. It's a nice-to-have that improves ergonomics but isn't critical for functionality.

**Independent Test**: Can be tested by implementing both initialization patterns for the lint command. Success is verified by creating command instances using both approaches and confirming identical behavior.

**Acceptance Scenarios**:

1. **Given** a simple command configuration, **When** a developer uses struct initialization `lint.NewCommand(lint.CommandOptions{Shared: opts, TargetVersion: "3.0"})`, **Then** the command is created with all specified options
2. **Given** a complex command configuration, **When** a developer uses functional options `lint.NewCommand(lint.WithShared(opts), lint.WithTargetVersion("3.0"))`, **Then** the command is created with the same options as struct initialization
3. **Given** a command constructor, **When** it accepts a CommandOptions struct, **Then** it also provides With* functional option functions for all configurable fields

---

### User Story 5 - Deprecated API Audit Modernizes Codebase (Priority: P5)

As a maintainer, I want all deprecated API usage identified and replaced so that the codebase stays modern, secure, and maintainable without relying on deprecated dependencies.

**Why this priority**: This is technical debt cleanup that improves long-term maintainability. While important, it doesn't add user-facing value and can be addressed incrementally.

**Independent Test**: Can be tested by auditing the codebase for deprecated API usage (IDE warnings, godoc markers) and replacing instances with modern alternatives. Success is verified by running all tests and confirming no deprecation warnings remain.

**Acceptance Scenarios**:

1. **Given** a codebase using deprecated APIs, **When** an audit is performed using IDE deprecation warnings, **Then** all deprecated API usage is identified
2. **Given** a deprecated API with a modern alternative, **When** the code is refactored, **Then** it uses the modern API and includes no deprecation warnings
3. **Given** a deprecated API without an alternative, **When** the code is reviewed, **Then** it includes a comment explaining why the deprecated API is required and references a tracking issue

---

### Edge Cases

- What happens when a check needs both lint and upgrade logic but they share significant code? (Answer: Extract shared logic into helper functions; keep mode-specific logic in separate code paths)
- How does the IOStreams wrapper handle nil writers? (Answer: Should validate writers in SharedOptions initialization and return error if nil)
- What happens when a command uses both struct and functional options simultaneously? (Answer: Functional options would need to be applied after struct initialization, or choose one pattern per constructor)
- How do existing tests handle the Command interface refactoring? (Answer: Tests should continue to work as the interface matches the existing Complete/Validate/Run pattern; AddFlags tests need to be added)
- What happens if intermediate commits during incremental implementation break the build? (Answer: Each user story must leave the codebase in a working state with all tests passing; backward compatibility is not a concern so internal API changes are acceptable between stories)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide an IOStreams wrapper in `pkg/util/iostreams/` with methods Fprintf, Fprintln, Errorf, and Errorln
- **FR-002**: IOStreams wrapper MUST automatically append newlines to all output (internally using fmt.Fprintln)
- **FR-003**: IOStreams wrapper MUST automatically select correct output stream (Out for Fprintf/Fprintln, ErrOut for Errorf/Errorln)
- **FR-004**: IOStreams wrapper formatting methods MUST only call fmt.Sprintf when args are provided
- **FR-005**: System MUST merge lint and upgrade commands into single lint command in `pkg/cmd/doctor/lint/`
- **FR-006**: Lint command MUST accept optional `--version` flag to specify target version
- **FR-007**: Lint command MUST validate current cluster state when `--version` is omitted
- **FR-008**: Lint command MUST assess upgrade readiness when `--version` differs from current version
- **FR-009**: Check implementations MUST detect execution mode by comparing `target.CurrentVersion` with `target.Version`
- **FR-010**: System MUST define a Command interface with methods: Complete() error, Validate() error, Run(ctx context.Context) error, AddFlags(fs *pflag.FlagSet)
- **FR-011**: Command structs MUST be named "Command" (not "Options")
- **FR-012**: Command constructors MUST be named "NewCommand()" (not "NewOptions()")
- **FR-013**: Command implementation files MUST be named `<command>.go` (not `options.go`)
- **FR-014**: Commands MUST initialize their own SharedOptions internally
- **FR-015**: System MUST support struct-based command initialization via CommandOptions struct
- **FR-016**: System MUST support functional option pattern via With* functions
- **FR-017**: System MUST audit codebase for deprecated API usage using IDE warnings and godoc markers
- **FR-018**: Code MUST replace deprecated APIs with modern alternatives where available
- **FR-019**: Code using unavoidable deprecated APIs MUST include explanatory comments with tracking issue references
- **FR-020**: Upgraded commands MUST remove the separate upgrade subcommand and its associated files
- **FR-021**: Implementation MUST follow user story priority order (P1 → P2 → P3 → P4 → P5) with each story completed and committed independently

### Key Entities

- **IOStreams Wrapper**: Utility struct providing formatted output methods that automatically manage newlines and stream selection. Contains In, Out, and ErrOut writers.
- **Command Interface**: Standard interface defining the lifecycle methods all commands must implement (Complete, Validate, Run, AddFlags).
- **CommandOptions Struct**: Configuration struct for struct-based initialization containing all configurable command properties.
- **Functional Options**: Functions following the With* naming pattern that configure Command instances via closures.
- **CheckTarget**: Execution context passed to checks containing CurrentVersion and Version fields for mode detection.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All commands using repetitive `fmt.Fprintf(o.Out, ...)` patterns are refactored to use IOStreams wrapper, reducing output-related lines of code by at least 30%
- **SC-002**: Users can validate current cluster state using `kubectl odh doctor lint` without specifying version flag
- **SC-003**: Users can assess upgrade readiness using `kubectl odh doctor lint --version <target>` with identical output format as lint mode
- **SC-004**: All checks adapt behavior based on version context without duplicating check implementations
- **SC-005**: All commands implement the Command interface with AddFlags method for centralized flag registration
- **SC-006**: Developer can create command instances using either struct initialization or functional options with identical results
- **SC-007**: Codebase audit identifies 100% of deprecated API usage and provides replacement plan for each instance
- **SC-008**: All quality gates (`make check`) pass after refactoring with zero new linting issues
- **SC-009**: Existing tests continue to pass without modification (except for new AddFlags tests)
- **SC-010**: Code review confirms no implementation details leak into user-facing command behavior changes
- **SC-011**: Each user story (P1 through P5) is delivered as an independent, working commit with all tests passing
