# Research: Promote Lint Command to Top Level

**Feature**: 004-promote-lint-command
**Date**: 2025-12-09
**Status**: Complete

## Overview

This research validates the technical approach for promoting the `lint` command from `kubectl odh doctor lint` to `kubectl odh lint`, removing the intermediary parent command and renaming the `--version` flag to `--target-version`.

## Research Questions & Findings

### 1. Command Registration Pattern

**Question**: How does Cobra handle top-level vs nested command registration?

**Research Approach**:
- Review Cobra documentation and kubectl plugin patterns
- Examine current cmd/doctor/doctor.go and cmd/root.go
- Identify differences between nested and top-level registration

**Findings**:
- **Nested commands** (current): Parent command created first, subcommands added to parent
  ```go
  parentCmd := &cobra.Command{Use: "doctor", ...}
  lintCmd := &cobra.Command{Use: "lint", ...}
  parentCmd.AddCommand(lintCmd)
  root.AddCommand(parentCmd)
  ```

- **Top-level commands** (target): Command added directly to root
  ```go
  lintCmd := &cobra.Command{Use: "lint", ...}
  root.AddCommand(lintCmd)
  ```

**Decision**: Use direct root registration. Create `cmd/lint.go` that registers lint command to root, eliminating the doctor parent command entirely.

**Rationale**: Simpler structure, fewer indirection layers, aligns with kubectl plugin patterns where most commands are top-level.

---

### 2. Package Import Path Changes

**Question**: What files reference `pkg/cmd/doctor` or `pkg/doctor`?

**Research Approach**:
```bash
# Find all Go files importing doctor packages
grep -r "pkg/cmd/doctor" --include="*.go"
grep -r "pkg/doctor" --include="*.go"

# Find test files that might reference doctor
find . -name "*_test.go" -exec grep -l "doctor" {} \;
```

**Findings**:
- `cmd/doctor/lint.go` - imports `pkg/cmd/doctor/lint`
- `pkg/cmd/doctor/lint/lint.go` - imports `pkg/doctor/check`, `pkg/doctor/version`
- All check implementations import `pkg/doctor/check`
- Test files import doctor packages for test setup

**Impact Areas**:
1. Command registration (`cmd/root.go`)
2. Lint command implementation (`pkg/cmd/doctor/lint/`)
3. Check framework (`pkg/doctor/check/`)
4. All check implementations (`pkg/doctor/checks/`)
5. Test files across the codebase

**Decision**: Systematic package renaming with import path updates:
- `pkg/cmd/doctor` → `pkg/cmd/lint`
- `pkg/doctor` → `pkg/lint`

**Migration Strategy**:
1. Move packages using git mv to preserve history
2. Update package declarations
3. Run gofmt and goimports to fix imports automatically
4. Manually verify critical imports

---

### 3. Flag Renaming Compatibility

**Question**: Can we cleanly remove `--version` without affecting other commands?

**Research Approach**:
- Examine current flag registration in `pkg/cmd/doctor/lint/lint.go`
- Check if `--version` is used elsewhere (root command, other subcommands)
- Verify no conflicts with global flags

**Findings**:
- `--version` flag is **local to lint command only**
  - Registered in `lint.AddFlags()` method
  - Not a persistent flag on parent or root
  - Used for: upgrade readiness assessment mode

- Root command has `--version` for CLI version (different context)
  - This is a Cobra built-in flag
  - Remains unchanged

**Decision**: Safe to rename lint command's `--version` flag to `--target-version`.

**Changes Required**:
```go
// Before (pkg/cmd/doctor/lint/lint.go)
func (c *Command) AddFlags(fs *pflag.FlagSet) {
    fs.StringVar(&c.targetVersion, "version", "", "Target version for upgrade assessment")
}

// After (pkg/cmd/lint/lint.go)
func (c *Command) AddFlags(fs *pflag.FlagSet) {
    fs.StringVar(&c.targetVersion, "target-version", "", "Target version for upgrade assessment")
}
```

**Rationale**: `--target-version` is more explicit and eliminates confusion with the global `--version` flag that shows CLI version.

---

### 4. Help Text Updates

**Question**: Where are command descriptions and examples defined?

**Research Approach**:
- Locate all hardcoded strings referencing "doctor" or command paths
- Identify constants vs inline strings
- Find documentation files

**Findings**:

**Code locations**:
- `cmd/doctor/lint.go`:
  - `lintCmdName`, `lintCmdShort`, `lintCmdLong`, `lintCmdExample` constants
  - Contains examples like `kubectl odh doctor lint`

- `cmd/doctor/doctor.go`:
  - `cmdName`, `cmdShort`, `cmdLong` constants
  - Describes doctor command (TO BE REMOVED)

**Documentation locations**:
- `.specify/memory/constitution.md` (already updated to v1.16.0)
- Historical spec files (001, 002, 003) - archived, no updates needed

**Decision**: Update all help text constants to reference `kubectl odh lint` instead of `kubectl odh doctor lint`.

**Changes Required**:
1. Update `lintCmdLong` and `lintCmdExample` in new `cmd/lint.go`
2. Update all example commands to use new flag name
3. Remove doctor command constants entirely

**Example Replacements**:
```bash
# Before
kubectl odh doctor lint
kubectl odh doctor lint --version 3.0

# After
kubectl odh lint
kubectl odh lint --target-version 3.0
```

---

### 5. Package Structure Best Practices

**Question**: Should we maintain `pkg/cmd/lint` vs moving everything to `pkg/lint`?

**Research Approach**:
- Review constitution Principle II (Extensible Command Structure)
- Examine kubectl plugin patterns
- Assess current structure vs constitutional requirements

**Findings**:

**Constitutional Requirements**:
- Command definition in `cmd/` (Cobra wrappers)
- Command business logic in `pkg/cmd/<command>/` (Complete/Validate/Run)
- Domain-specific logic in `pkg/<domain>/` (optional)

**Current Structure**:
- `pkg/cmd/doctor/lint/` - Lint command implementation
- `pkg/doctor/check/` - Check framework (domain logic)
- `pkg/doctor/checks/` - Check implementations
- `pkg/doctor/version/` - Version detection

**Decision**: Two-tier package structure:
1. **pkg/cmd/lint/** - Command implementation (lint.go, lint_options.go, shared_options.go)
2. **pkg/lint/** - Domain logic (check/, version/, checks/)

**Rationale**:
- Separates command concerns (CLI handling) from domain logic (checking, version detection)
- Aligns with constitution's separation of Cobra wrappers and business logic
- Enables reuse of domain packages without command dependencies
- Standard pattern across kubectl plugins

---

## Technology Decisions

### Code Movement Strategy

**Decision**: Use `git mv` for all file moves to preserve git history.

**Commands**:
```bash
# Move command packages
git mv pkg/cmd/doctor pkg/cmd/lint

# Move domain packages
git mv pkg/doctor pkg/lint

# Move command registration
git mv cmd/doctor/lint.go cmd/lint.go

# Remove doctor command
git rm cmd/doctor/doctor.go
git rm -r cmd/doctor/
```

**Rationale**: Git history preservation aids debugging and maintains attribution.

---

### Import Path Update Strategy

**Decision**: Automated tooling + manual verification

**Tools**:
1. **gofmt** - Format code
2. **goimports** - Auto-fix import paths
3. **gopls** - IDE support for refactoring
4. **Manual grep** - Verify no missed references

**Process**:
1. Move packages
2. Update package declarations manually
3. Run `goimports -w .` to fix import statements
4. Run `make lint-fix` to apply auto-fixes
5. Run `make check` to verify

**Rationale**: Automated tooling reduces manual errors while manual verification catches edge cases.

---

## Alternatives Considered

### Alternative 1: Keep doctor command with deprecation warning

**Approach**: Maintain `kubectl odh doctor lint` with deprecation warning, add `kubectl odh lint` as alias.

**Rejected Because**:
- Increases maintenance burden (two command paths)
- Adds complexity for deprecation handling
- User specified breaking changes are acceptable
- Constitution updated to prohibit doctor parent command

---

### Alternative 2: Alias --version to --target-version

**Approach**: Support both flags with `--version` as deprecated alias.

**Rejected Because**:
- Confusing to have two flags for same functionality
- Conflicts with CLI's global --version flag semantically
- Breaking change already accepted, no need for compatibility layer
- Constitution explicitly prohibits ambiguous --version flag

---

### Alternative 3: Keep pkg/doctor package name

**Approach**: Only move command from `cmd/doctor/` to `cmd/lint/`, keep `pkg/doctor/` unchanged.

**Rejected Because**:
- Package name "doctor" no longer reflects purpose (lint diagnostics)
- Constitution requires package names match domain (lint, not doctor)
- Creates inconsistency between command name and package name
- Misleading for future developers

---

## Implementation Risks & Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|-----------|
| Missed import path | Medium | High | Automated search + `make check` catches compile errors |
| Test failures | Low | Medium | Existing tests validate behavior unchanged |
| Breaking user scripts | High | High | Clear migration docs + release notes |
| Git history loss | Low | Low | Use `git mv` preserves history |
| Performance regression | Very Low | Medium | Same code path, benchmark to verify |

---

## Success Criteria

✅ **Research Complete** when all findings documented
✅ **Decisions Made** for all unknowns
✅ **Alternatives Evaluated** and rejection rationale captured
✅ **Implementation Path Clear** with no blocking unknowns

---

## Next Steps

1. Proceed to **Phase 1: Design** to generate contracts and data model
2. Create `data-model.md` with package structure details
3. Create `contracts/command-api.md` with command interface specification
4. Create `quickstart.md` with usage examples
5. Update agent context files
6. Generate tasks in Phase 2

---

## References

- [Cobra Command Documentation](https://github.com/spf13/cobra)
- [kubectl Plugin Development Guide](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/)
- Constitution v1.16.0 - Principle XI (Lint Command Architecture)
- Constitution v1.16.0 - Principle II (Extensible Command Structure)
