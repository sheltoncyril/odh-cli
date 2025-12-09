# Data Model: Promote Lint Command to Top Level

**Feature**: 004-promote-lint-command
**Date**: 2025-12-09

## Overview

This document defines the data structures and package organization for the promoted lint command. Since this is a refactoring task, there are no new data entities—only package restructuring and field renaming.

## Package Structure

### Command Layer (pkg/cmd/lint/)

Location for CLI command implementation following constitution Principle II.

```go
// pkg/cmd/lint/lint.go
package lint

import (
    "context"
    "github.com/spf13/pflag"
    "k8s.io/cli-runtime/pkg/genericiooptions"
)

// Command implements the lint command following the Command interface pattern
type Command struct {
    shared        *SharedOptions
    targetVersion string  // Target version for upgrade assessment (renamed from --version)

    // Additional fields remain unchanged
    checksPattern string
    severity      string
    output        string
}

// NewCommand creates a new lint command with struct initialization (PREFERRED)
func NewCommand(opts CommandOptions) *Command {
    return &Command{
        shared:        opts.Shared,
        targetVersion: opts.TargetVersion,
        checksPattern: opts.ChecksPattern,
        severity:      opts.Severity,
        output:        opts.Output,
    }
}

// Interface methods (constitutional requirement)
func (c *Command) Complete() error
func (c *Command) Validate() error
func (c *Command) Run(ctx context.Context) error
func (c *Command) AddFlags(fs *pflag.FlagSet)
```

### Options Structures (pkg/cmd/lint/lint_options.go)

```go
// pkg/cmd/lint/lint_options.go
package lint

// CommandOptions for struct-based initialization (PREFERRED pattern)
type CommandOptions struct {
    Shared        *SharedOptions
    TargetVersion string  // Renamed from Version
    ChecksPattern string
    Severity      string
    Output        string
}

// CommandOption for functional options pattern (ADVANCED pattern)
type CommandOption func(*Command)

func WithTargetVersion(version string) CommandOption {
    return func(c *Command) {
        c.targetVersion = version
    }
}

func WithShared(shared *SharedOptions) CommandOption {
    return func(c *Command) {
        c.shared = shared
    }
}

// Additional With* functions for other fields...
```

### Shared Options (pkg/cmd/lint/shared_options.go)

Moved from `pkg/cmd/doctor/shared_options.go` without modification.

```go
// pkg/cmd/lint/shared_options.go
package lint

import (
    "k8s.io/cli-runtime/pkg/genericclioptions"
    "k8s.io/client-go/dynamic"
    "github.com/lburgazzoli/odh-cli/pkg/util/iostreams"
)

type SharedOptions struct {
    ConfigFlags *genericclioptions.ConfigFlags
    IOStreams   *iostreams.IOStreams

    // Populated during Complete()
    client dynamic.Interface
}

func NewSharedOptions(io *iostreams.IOStreams) *SharedOptions
func (o *SharedOptions) Complete() error
func (o *SharedOptions) Client() dynamic.Interface
```

## Domain Layer (pkg/lint/)

### Check Framework (pkg/lint/check/)

Moved from `pkg/doctor/check/` without modification to internal logic.

```go
// pkg/lint/check/check.go
package check

type Check interface {
    ID() string
    Category() CheckCategory
    Run(ctx context.Context, target *CheckTarget) CheckResult
}

type CheckTarget struct {
    Client         dynamic.Interface
    Version        *version.Version  // Current cluster version
    CurrentVersion *version.Version  // Same as Version
    // Note: Import path changes to pkg/lint/version
}

type CheckResult struct {
    CheckID     string
    Status      CheckStatus
    Message     string
    Severity    *CheckSeverity
    Remediation string
}

// Additional check framework types unchanged...
```

### Version Detection (pkg/lint/version/)

Moved from `pkg/doctor/version/` without modification.

```go
// pkg/lint/version/version.go
package version

type Version struct {
    Version string  // e.g., "2.15.0"
    // Additional fields...
}

func Detect(ctx context.Context, client dynamic.Interface) (*Version, error)
```

### Check Implementations (pkg/lint/checks/)

Moved from `pkg/doctor/checks/` with updated import paths only.

```go
// pkg/lint/checks/components/dashboard/dashboard.go
package dashboard

import (
    "context"
    "github.com/lburgazzoli/odh-cli/pkg/lint/check"  // Updated import
)

type Check struct {
    // Fields unchanged
}

func (c *Check) Run(ctx context.Context, target *check.CheckTarget) check.CheckResult {
    // Implementation unchanged
}
```

## Command Registration Layer (cmd/)

### Lint Command Registration (cmd/lint.go)

New file replacing `cmd/doctor/lint.go`.

```go
// cmd/lint.go
package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "k8s.io/cli-runtime/pkg/genericclioptions"
    "k8s.io/cli-runtime/pkg/genericiooptions"

    "github.com/lburgazzoli/odh-cli/pkg/cmd/lint"  // Updated import
)

const (
    lintCmdName  = "lint"
    lintCmdShort = "Validate current OpenShift AI installation or assess upgrade readiness"
    lintCmdLong  = `
Validates the current OpenShift AI installation or assesses upgrade readiness.

LINT MODE (without --target-version):
  Validates the current cluster state and reports configuration issues.

UPGRADE MODE (with --target-version):
  Assesses upgrade readiness by comparing current version against target version.

...
`
    lintCmdExample = `
  # Validate current cluster state
  kubectl odh lint

  # Assess upgrade readiness for version 3.0
  kubectl odh lint --target-version 3.0

...
`
)

func AddLintCommand(parent *cobra.Command, flags *genericclioptions.ConfigFlags) {
    streams := genericiooptions.IOStreams{
        In:     parent.InOrStdin(),
        Out:    parent.OutOrStdout(),
        ErrOut: parent.ErrOrStderr(),
    }

    command := lint.NewCommand(lint.CommandOptions{
        Shared: lint.NewSharedOptions(&streams),
    })
    command.ConfigFlags = flags

    cmd := &cobra.Command{
        Use:     lintCmdName,
        Short:   lintCmdShort,
        Long:    lintCmdLong,
        Example: lintCmdExample,
        RunE: func(cmd *cobra.Command, _ []string) error {
            if err := command.Complete(); err != nil {
                return fmt.Errorf("completing lint command: %w", err)
            }
            if err := command.Validate(); err != nil {
                return fmt.Errorf("validating lint command: %w", err)
            }
            return command.Run(cmd.Context())
        },
    }

    command.AddFlags(cmd.Flags())
    parent.AddCommand(cmd)
}
```

### Root Command Updates (cmd/root.go)

```go
// cmd/root.go - Updated to register lint directly
package cmd

import (
    "github.com/spf13/cobra"
    "k8s.io/cli-runtime/pkg/genericclioptions"
)

func Execute() {
    flags := genericclioptions.NewConfigFlags(true)
    rootCmd := &cobra.Command{
        Use:   "kubectl-odh",
        Short: "OpenShift AI CLI plugin",
    }

    // Register lint command directly (no doctor parent)
    AddLintCommand(rootCmd, flags)

    // Other commands...

    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}
```

## Data Flow

### Lint Mode (Current State Validation)

```
User Input: kubectl odh lint
↓
cmd/lint.go → AddLintCommand()
↓
pkg/cmd/lint/lint.go → Command{targetVersion: ""}
↓
Complete() → Populate client, detect current version
↓
Validate() → Verify flags
↓
Run() → Execute checks in lint mode
↓
Output → Display results (table/JSON/YAML)
```

### Upgrade Mode (Target Version Specified)

```
User Input: kubectl odh lint --target-version 3.0
↓
cmd/lint.go → AddLintCommand()
↓
pkg/cmd/lint/lint.go → Command{targetVersion: "3.0"}
↓
Complete() → Populate client, detect current version, parse target
↓
Validate() → Verify version format
↓
Run() → Execute checks in upgrade mode (compare current vs target)
↓
Output → Display upgrade readiness results
```

## Field Changes Summary

| Location | Old Field | New Field | Type | Notes |
|----------|-----------|-----------|------|-------|
| Command struct | `version string` | `targetVersion string` | string | Internal field rename |
| AddFlags() | `--version` | `--target-version` | flag | CLI flag rename |
| CommandOptions | `Version string` | `TargetVersion string` | string | Struct field rename |
| Package path | `pkg/cmd/doctor/lint` | `pkg/cmd/lint` | import | Package relocation |
| Package path | `pkg/doctor` | `pkg/lint` | import | Package relocation |

## Import Path Migration Matrix

| Old Import | New Import | Affected Files |
|------------|------------|----------------|
| `github.com/lburgazzoli/odh-cli/pkg/cmd/doctor/lint` | `github.com/lburgazzoli/odh-cli/pkg/cmd/lint` | cmd/root.go, tests |
| `github.com/lburgazzoli/odh-cli/pkg/doctor/check` | `github.com/lburgazzoli/odh-cli/pkg/lint/check` | All check implementations |
| `github.com/lburgazzoli/odh-cli/pkg/doctor/version` | `github.com/lburgazzoli/odh-cli/pkg/lint/version` | pkg/cmd/lint/lint.go, checks |
| `github.com/lburgazzoli/odh-cli/pkg/doctor/checks/*` | `github.com/lburgazzoli/odh-cli/pkg/lint/checks/*` | Test files, registry |

## Validation Rules

### Command Validation (unchanged)

- **targetVersion**: Must be valid semantic version format (e.g., "2.15.0", "3.0")
- **checksPattern**: Optional glob pattern (e.g., "components/*", "*dashboard*")
- **severity**: Must be one of: critical, warning, info (if specified)
- **output**: Must be one of: table, json, yaml

### Behavior Validation

- **Lint mode** (no target version): Validate current cluster state
- **Upgrade mode** (target version specified): Compare current vs target for breaking changes

## State Transitions

No state machine—command is stateless. Each execution:
1. Complete (populate fields, create clients)
2. Validate (verify inputs)
3. Run (execute checks)
4. Return results

## Testing Considerations

### Package Move Testing

- Verify all imports update correctly
- Run `go mod tidy` to clean dependencies
- Verify no circular dependencies

### Field Rename Testing

- Test `--target-version` flag parsing
- Test empty/missing target version (lint mode)
- Test invalid target version (error handling)
- Verify existing tests pass with renamed fields

### Integration Testing

- End-to-end: `kubectl odh lint` (no target)
- End-to-end: `kubectl odh lint --target-version 3.0`
- Verify identical behavior to old command
- Verify removed command returns error

## Constitutional Compliance

✅ **Principle II** (Extensible Command Structure): Command implements interface methods, file named lint.go

✅ **Principle IV** (Flexible Initialization): Supports both struct and functional options

✅ **Package Granularity**: Clear separation between cmd/lint/ (command) and pkg/lint/ (domain)

✅ **Code Organization**: Follows cmd/ and pkg/ structure with domain-specific packages
