# Command API Contract: Lint Command

**Feature**: 004-promote-lint-command
**Date**: 2025-12-09
**Version**: 1.0.0

## Overview

This document defines the public API contract for the `kubectl odh lint` command. This contract ensures behavioral compatibility during and after the refactoring from `kubectl odh doctor lint` to `kubectl odh lint`.

## Command Interface

### Package Export

**Import Path**: `github.com/lburgazzoli/odh-cli/pkg/cmd/lint`

**Exported Types**:
```go
type Command struct { /* private fields */ }
type CommandOptions struct { /* public fields */ }
type CommandOption func(*Command)
type SharedOptions struct { /* public fields */ }
```

**Exported Functions**:
```go
func NewCommand(opts CommandOptions) *Command
func NewCommand(options ...CommandOption) *Command
func NewSharedOptions(io *iostreams.IOStreams) *SharedOptions
```

### Command Methods

All commands MUST implement these methods (constitutional requirement):

```go
// Complete populates fields and creates clients
func (c *Command) Complete() error

// Validate verifies all required fields and constraints
func (c *Command) Validate() error

// Run executes the lint validation logic
func (c *Command) Run(ctx context.Context) error

// AddFlags registers command-specific flags
func (c *Command) AddFlags(fs *pflag.FlagSet)
```

## CLI Interface

### Command Path

**Old (removed)**:
```bash
kubectl odh doctor lint [flags]
```

**New (required)**:
```bash
kubectl odh lint [flags]
```

### Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--target-version` | | string | "" | Target version for upgrade assessment |
| `--output` | `-o` | string | "table" | Output format: table, json, yaml |
| `--checks` | | string | "*" | Glob pattern to filter checks |
| `--severity` | | string | "" | Filter by severity: critical, warning, info |
| `--kubeconfig` | | string | "" | Path to kubeconfig file |
| `--context` | | string | "" | Kubernetes context to use |
| `--namespace` | `-n` | string | "" | DEPRECATED: Not used for lint command |

**Flag Changes from v1**:
- ❌ **Removed**: `--version` (renamed to `--target-version`)
- ✅ **Added**: `--target-version` (explicit name for clarity)
- ✅ **Unchanged**: All other flags preserve identical behavior

### Examples

#### Basic Usage

```bash
# Validate current cluster state (lint mode)
kubectl odh lint

# Assess upgrade readiness to version 3.0 (upgrade mode)
kubectl odh lint --target-version 3.0

# Output in JSON format
kubectl odh lint -o json

# Run only component checks
kubectl odh lint --checks "components/*"

# Filter by critical severity
kubectl odh lint --severity critical
```

#### Combined Flags

```bash
# Upgrade assessment with JSON output
kubectl odh lint --target-version 3.0 -o json

# Specific checks with severity filter
kubectl odh lint --checks "*dashboard*" --severity warning

# Full example
kubectl odh lint \
  --target-version 3.1 \
  --output yaml \
  --checks "components/*" \
  --severity critical
```

## Behavior Contract

### Lint Mode (No Target Version)

**Trigger**: `kubectl odh lint` (no `--target-version` specified)

**Behavior**:
1. Detect current cluster OpenShift AI version
2. Execute all checks (or filtered by `--checks` pattern)
3. Validate current cluster state against expected configuration
4. Report issues with severity, message, and remediation

**Output Example** (table format):
```
CATEGORY      CHECK               STATUS  SEVERITY  MESSAGE
components    dashboard           PASS    -         Dashboard ready
components    workbenches         PASS    -         Workbenches ready
services      servicemesh         FAIL    WARNING   ServiceMesh deprecated in 3.x
```

### Upgrade Mode (Target Version Specified)

**Trigger**: `kubectl odh lint --target-version X.Y`

**Behavior**:
1. Detect current cluster OpenShift AI version
2. Parse target version
3. Execute all checks (or filtered by `--checks` pattern)
4. Compare current vs target version for breaking changes
5. Report upgrade readiness with blockers and warnings

**Output Example** (table format):
```
CATEGORY      CHECK               STATUS  SEVERITY  MESSAGE
components    kserve              FAIL    CRITICAL  KServe serverless mode removed in 3.0
components    modelmesh           WARN    WARNING   ModelMesh deprecated, migrate to KServe
services      servicemesh         PASS    INFO      ServiceMesh already disabled
```

### Output Formats

#### Table Format (default)

Human-readable table with aligned columns:
```
CATEGORY      CHECK               STATUS  SEVERITY  MESSAGE
components    dashboard           PASS    -         Dashboard ready
```

#### JSON Format (`-o json`)

Machine-parsable JSON array:
```json
{
  "checks": [
    {
      "category": "components",
      "checkID": "dashboard.ready",
      "status": "pass",
      "severity": null,
      "message": "Dashboard ready",
      "remediation": ""
    }
  ],
  "summary": {
    "total": 10,
    "passed": 9,
    "failed": 1,
    "warnings": 0
  }
}
```

#### YAML Format (`-o yaml`)

Machine-parsable YAML:
```yaml
checks:
  - category: components
    checkID: dashboard.ready
    status: pass
    severity: null
    message: Dashboard ready
    remediation: ""
summary:
  total: 10
  passed: 9
  failed: 1
  warnings: 0
```

## Exit Codes

| Code | Meaning | When |
|------|---------|------|
| 0 | Success | All checks passed OR only warnings/info |
| 1 | Validation failed | At least one check failed with CRITICAL/ERROR severity |
| 2 | Command error | Invalid flags, connection errors, or runtime errors |

**Exit Code Examples**:
```bash
# Exit 0 - all checks pass
kubectl odh lint
echo $?  # 0

# Exit 1 - critical issues found
kubectl odh lint --target-version 3.0
echo $?  # 1 (blocking upgrade issues)

# Exit 2 - invalid flag
kubectl odh lint --invalid-flag
echo $?  # 2
```

## Error Handling

### Invalid Target Version

**Input**: `kubectl odh lint --target-version invalid`

**Output**:
```
Error: invalid target version "invalid": must be semantic version (e.g., "2.15.0", "3.0")
```

**Exit Code**: 2

### Removed Flag Usage

**Input**: `kubectl odh lint --version 3.0`

**Output**:
```
Error: unknown flag: --version
Did you mean --target-version?
```

**Exit Code**: 2

### Cluster Connection Failure

**Input**: `kubectl odh lint` (cluster unreachable)

**Output**:
```
Error: failed to connect to cluster: unable to load kubeconfig
```

**Exit Code**: 2

### Permission Denied

**Input**: `kubectl odh lint` (insufficient RBAC)

**Output**:
```
Error: insufficient permissions to list resources
Required: get, list permissions for datascienceclusters.datasciencecluster.opendatahub.io
```

**Exit Code**: 2

## Backward Compatibility

### Breaking Changes

**Command Path**:
- ❌ **Old**: `kubectl odh doctor lint` → **Removed**
- ✅ **New**: `kubectl odh lint` → **Required**

**Flag Names**:
- ❌ **Old**: `--version` → **Removed**
- ✅ **New**: `--target-version` → **Required**

### Migration Guide

Users MUST update:
1. Scripts: Replace `kubectl odh doctor lint` with `kubectl odh lint`
2. Flags: Replace `--version X.Y` with `--target-version X.Y`
3. Documentation: Update all command references

**Before**:
```bash
#!/bin/bash
kubectl odh doctor lint --version 3.0 -o json > results.json
```

**After**:
```bash
#!/bin/bash
kubectl odh lint --target-version 3.0 -o json > results.json
```

## Performance Contract

**Constraints** (from spec SC-006):
- Command execution time MUST remain within 5% of baseline
- Startup time (help/version display) < 500ms
- Full validation execution time dependent on cluster size (unchanged from v1)

**Baseline** (current `kubectl odh doctor lint`):
- Help display: ~100ms
- Lint mode (10 checks): ~2-5 seconds
- Upgrade mode (10 checks): ~2-5 seconds

**Monitoring**:
```bash
# Benchmark command execution
time kubectl odh lint
time kubectl odh lint --target-version 3.0
```

## Testing Contract

### Unit Test Coverage

Required test cases:
- ✅ Command initialization with new package path
- ✅ Flag parsing with `--target-version`
- ✅ Flag validation (invalid version formats)
- ✅ Lint mode execution (no target version)
- ✅ Upgrade mode execution (with target version)
- ✅ Output format rendering (table/JSON/YAML)
- ✅ Error handling (connection failures, invalid flags)

### Integration Test Coverage

Required scenarios:
- ✅ Full command: `kubectl odh lint`
- ✅ Upgrade mode: `kubectl odh lint --target-version X.Y`
- ✅ All output formats work identically
- ✅ All flag combinations preserved
- ✅ Exit codes correct for each scenario

### Regression Test Coverage

Required validations:
- ✅ Existing lint functionality unchanged
- ✅ Check execution logic identical
- ✅ Output format consistency maintained
- ✅ Performance within 5% of baseline

## Version Compatibility

**CLI Version**: Any version with this feature
**Kubernetes Version**: 1.24+ (unchanged)
**OpenShift AI Version**: 2.x, 3.x (detection via cluster API)

**Version Detection**:
```go
// Current version detection unchanged
currentVersion, err := version.Detect(ctx, client)

// Target version parsing
if c.targetVersion != "" {
    targetVersion, err := semver.Parse(c.targetVersion)
}
```

## Security Considerations

**RBAC Requirements** (unchanged):
- `get`, `list` permissions on DataScienceCluster
- `get`, `list` permissions on DSCInitialization
- `get`, `list` permissions on CRDs
- `get`, `list` permissions on ClusterServiceVersions

**No New Security Implications**:
- Command restructuring does not affect security model
- Flag rename does not introduce vulnerabilities
- Same cluster access patterns as v1

## Deprecation Notice

**Deprecated in**: This release
**Removed in**: This release

| Item | Old | New | Status |
|------|-----|-----|--------|
| Command path | `kubectl odh doctor lint` | `kubectl odh lint` | ❌ Removed |
| Parent command | `kubectl odh doctor` | N/A | ❌ Removed |
| Flag | `--version` | `--target-version` | ❌ Removed |

**No Deprecation Period**: This is a breaking change with immediate effect.

## API Stability Guarantee

**Stable** (will not change without major version bump):
- Command name: `lint`
- Flag names: `--target-version`, `--output`, `--checks`, `--severity`
- Output formats: table, json, yaml
- Exit codes: 0, 1, 2 semantics
- Command interface methods: Complete, Validate, Run, AddFlags

**May Change** (minor version updates allowed):
- Check implementations (new checks added)
- Check messages and remediation text
- Help text improvements
- Performance optimizations

## Support

**Issues**: Report to https://github.com/lburgazzoli/odh-cli/issues
**Documentation**: See [quickstart.md](../quickstart.md) for usage guide
**Constitution**: v1.16.0 (Principle XI: Lint Command Architecture)
