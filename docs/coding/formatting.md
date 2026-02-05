# Code Formatting

This document covers code formatting rules for odh-cli development.

For other coding conventions, see [conventions.md](conventions.md) and [patterns.md](patterns.md).

## Code Formatting

**CRITICAL: MUST use `make fmt` to format code. NEVER use `gci` or other formatters directly.**

```bash
# ✓ CORRECT - Format all code
make fmt

# ❌ WRONG - DO NOT use gci directly
gci write ./...
gci write -s standard -s default ./...

# ❌ WRONG - DO NOT use gofmt directly
gofmt -w .

# ❌ WRONG - DO NOT use goimports directly
goimports -w .
```

**Why you MUST use `make fmt`:**
- **Safety**: The Makefile applies correct flags and configuration
- **Consistency**: All developers use identical formatting configuration
- **Completeness**: `make fmt` runs all necessary formatters in the correct order

**What `make fmt` does:**
1. Runs `go fmt` for basic formatting
2. Runs `gci` with project-specific import grouping rules

**Never run formatting tools directly.** Always use `make fmt`.
