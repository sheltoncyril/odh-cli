# Quickstart: kubectl odh lint

**Feature**: 004-promote-lint-command
**Date**: 2025-12-09

## Overview

The `kubectl odh lint` command validates your OpenShift AI cluster and assesses upgrade readiness. This guide provides quick-reference examples for common use cases.

## Installation

```bash
# Ensure kubectl-odh is in your PATH
kubectl odh version

# Verify lint command is available
kubectl odh lint --help
```

## Basic Usage

### Validate Current Cluster

Validate your current OpenShift AI installation:

```bash
kubectl odh lint
```

**Output**:
```
CATEGORY      CHECK               STATUS  SEVERITY  MESSAGE
components    dashboard           PASS    -         Dashboard ready
components    workbenches         PASS    -         Workbenches ready
services      oauth               PASS    -         OAuth configured
```

### Assess Upgrade Readiness

Check if your cluster is ready to upgrade to a specific version:

```bash
kubectl odh lint --target-version 3.0
```

**Output**:
```
CATEGORY      CHECK               STATUS  SEVERITY  MESSAGE
components    kserve              FAIL    CRITICAL  KServe serverless mode removed in 3.0
components    modelmesh           WARN    WARNING   ModelMesh deprecated, migrate to KServe
services      servicemesh         PASS    INFO      ServiceMesh already disabled
```

## Output Formats

### Human-Readable Table (default)

```bash
kubectl odh lint
```

### Machine-Parsable JSON

```bash
kubectl odh lint -o json
```

**Output**:
```json
{
  "checks": [
    {
      "category": "components",
      "checkID": "dashboard.ready",
      "status": "pass",
      "severity": null,
      "message": "Dashboard ready"
    }
  ],
  "summary": {
    "total": 10,
    "passed": 9,
    "failed": 1
  }
}
```

### Machine-Parsable YAML

```bash
kubectl odh lint -o yaml
```

## Filtering Checks

### Run Specific Check Categories

```bash
# Only component checks
kubectl odh lint --checks "components/*"

# Only service checks
kubectl odh lint --checks "services/*"

# Only dashboard-related checks
kubectl odh lint --checks "*dashboard*"
```

### Filter by Severity

```bash
# Show only critical issues
kubectl odh lint --severity critical

# Show warnings and above
kubectl odh lint --severity warning

# Show all including info
kubectl odh lint --severity info
```

## Common Scenarios

### Pre-Upgrade Validation

Before upgrading OpenShift AI:

```bash
# Check upgrade readiness
kubectl odh lint --target-version 3.0

# Get detailed results in JSON
kubectl odh lint --target-version 3.0 -o json > upgrade-assessment.json

# Check only critical blockers
kubectl odh lint --target-version 3.0 --severity critical
```

### Continuous Integration

Use in CI/CD pipelines:

```bash
#!/bin/bash
# Validate cluster health
kubectl odh lint -o json > lint-results.json

# Exit code 0 = passed, 1 = failed, 2 = error
if [ $? -eq 0 ]; then
  echo "Cluster validation passed"
else
  echo "Cluster validation failed"
  cat lint-results.json
  exit 1
fi
```

### Troubleshooting

Debug specific component issues:

```bash
# Check specific component
kubectl odh lint --checks "components/dashboard"

# Get detailed output with all severities
kubectl odh lint --checks "*workbench*" -o yaml
```

## Exit Codes

| Code | Meaning | Action |
|------|---------|--------|
| 0 | All checks passed | Proceed with operations |
| 1 | Validation failed | Review and fix critical issues |
| 2 | Command error | Check flags and cluster connection |

**Example**:
```bash
kubectl odh lint --target-version 3.0
EXIT_CODE=$?

if [ $EXIT_CODE -eq 0 ]; then
  echo "Ready to upgrade"
elif [ $EXIT_CODE -eq 1 ]; then
  echo "Blocking issues found, do not upgrade"
else
  echo "Command failed, check cluster connection"
fi
```

## Flag Reference

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--target-version` | | string | "" | Target version for upgrade assessment |
| `--output` | `-o` | string | table | Output format: table, json, yaml |
| `--checks` | | string | * | Glob pattern to filter checks |
| `--severity` | | string | "" | Filter by severity: critical, warning, info |

## Examples by Use Case

### Daily Health Checks

```bash
# Quick validation
kubectl odh lint

# Save results for tracking
kubectl odh lint -o json | jq '.' > daily-health-$(date +%Y%m%d).json
```

### Pre-Upgrade Planning

```bash
# Assess upgrade to 3.0
kubectl odh lint --target-version 3.0

# Get only blockers
kubectl odh lint --target-version 3.0 --severity critical

# Full assessment report
kubectl odh lint --target-version 3.0 -o yaml > upgrade-3.0-assessment.yaml
```

### Component-Specific Validation

```bash
# Validate dashboard
kubectl odh lint --checks "components/dashboard"

# Validate all workload components
kubectl odh lint --checks "components/*" --severity warning

# Check KServe configuration
kubectl odh lint --checks "*kserve*"
```

### Automated Monitoring

```bash
#!/bin/bash
# Run hourly validation
*/60 * * * * kubectl odh lint -o json > /var/log/odh-lint.json

# Alert on failures
kubectl odh lint || send-alert "OpenShift AI validation failed"
```

## Migration from Old Command

If you were using `kubectl odh doctor lint`:

**Old Command**:
```bash
kubectl odh doctor lint
kubectl odh doctor lint --version 3.0
```

**New Command** (use this):
```bash
kubectl odh lint
kubectl odh lint --target-version 3.0
```

**Changes**:
- ❌ Removed `doctor` parent command
- ❌ Renamed `--version` to `--target-version`
- ✅ All other functionality identical

## Troubleshooting

### Command Not Found

```bash
$ kubectl odh lint
Error: unknown command "lint" for "kubectl-odh"
```

**Solution**: Update kubectl-odh to latest version with lint command support.

### Unknown Flag Error

```bash
$ kubectl odh lint --version 3.0
Error: unknown flag: --version
```

**Solution**: Use `--target-version` instead of `--version`:
```bash
kubectl odh lint --target-version 3.0
```

### Permission Denied

```bash
$ kubectl odh lint
Error: insufficient permissions to list resources
```

**Solution**: Ensure your user/service account has required RBAC permissions:
- `get`, `list` on DataScienceCluster
- `get`, `list` on DSCInitialization
- `get`, `list` on CRDs

### Connection Refused

```bash
$ kubectl odh lint
Error: failed to connect to cluster: connection refused
```

**Solution**: Verify kubeconfig and cluster accessibility:
```bash
kubectl cluster-info
kubectl get nodes
```

## Advanced Usage

### Custom Kubeconfig

```bash
kubectl odh lint --kubeconfig=/path/to/kubeconfig
```

### Specific Context

```bash
kubectl odh lint --context=production-cluster
```

### Combined Flags

```bash
kubectl odh lint \
  --target-version 3.1 \
  --output json \
  --checks "components/*" \
  --severity critical \
  --kubeconfig=/path/to/kubeconfig
```

## Integration Examples

### Ansible Playbook

```yaml
- name: Validate OpenShift AI cluster
  command: kubectl odh lint -o json
  register: lint_results
  failed_when: lint_results.rc == 1

- name: Parse results
  set_fact:
    lint_summary: "{{ (lint_results.stdout | from_json).summary }}"

- name: Report
  debug:
    msg: "Validation: {{ lint_summary.passed }}/{{ lint_summary.total }} checks passed"
```

### GitHub Actions

```yaml
steps:
  - name: Run cluster validation
    run: |
      kubectl odh lint -o json > lint-results.json
    continue-on-error: true

  - name: Upload results
    uses: actions/upload-artifact@v3
    with:
      name: lint-results
      path: lint-results.json
```

### Prometheus Alerting

```bash
# Export metrics
kubectl odh lint -o json | jq '.summary.failed' > /var/lib/prometheus/odh_lint_failures.prom

# Alert rule
- alert: ODHValidationFailed
  expr: odh_lint_failures > 0
  annotations:
    summary: "OpenShift AI validation failed"
```

## Tips & Best Practices

1. **Run regularly**: Schedule daily validations to catch configuration drift
2. **Use JSON output**: Easier to parse and integrate with tools
3. **Filter checks**: Use `--checks` to focus on specific areas
4. **Save results**: Track validation history for trend analysis
5. **CI/CD integration**: Validate after deployments automatically

## Next Steps

- See [command-api.md](./contracts/command-api.md) for complete API documentation
- See [plan.md](./plan.md) for implementation details
- Run `/speckit.tasks` to see implementation task list

## Support

**Issues**: https://github.com/lburgazzoli/odh-cli/issues
**Documentation**: See full command API in contracts/command-api.md
