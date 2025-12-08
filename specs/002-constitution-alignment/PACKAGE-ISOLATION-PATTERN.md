# Package Isolation Pattern for Diagnostic Checks

**Date**: 2025-12-08
**Constitution Reference**: Principle XI (Package Granularity), v1.14.0

## Pattern Description

All diagnostic checks MUST be in their own dedicated package following the pattern:

```
pkg/doctor/checks/<category>/<check>/
├── <check>.go       # Implementation
└── <check>_test.go  # Tests
```

## Categories

- `components/` - DataScienceCluster component checks (e.g., codeflare, kserve, kueue, modelmesh)
- `dependencies/` - Dependency operator checks (e.g., certmanager, kueueoperator, servicemeshoperator)
- `services/` - Service checks (e.g., servicemesh)
- `workloads/` - High-level custom resource checks (no Pods, Deployments, StatefulSets)

## Verified Isolated Checks

### Components (4 checks)
- ✅ `pkg/doctor/checks/components/codeflare/` (codeflare.go + codeflare_test.go)
- ✅ `pkg/doctor/checks/components/kserve/` (kserve.go + kserve_test.go)
- ✅ `pkg/doctor/checks/components/kueue/` (kueue.go + kueue_test.go)
- ✅ `pkg/doctor/checks/components/modelmesh/` (modelmesh.go + modelmesh_test.go)

### Dependencies (3 checks)
- ✅ `pkg/doctor/checks/dependencies/certmanager/`
- ✅ `pkg/doctor/checks/dependencies/kueueoperator/`
- ✅ `pkg/doctor/checks/dependencies/servicemeshoperator/`

### Services (1 check)
- ✅ `pkg/doctor/checks/services/servicemesh/`

**Total**: 8 isolated check packages

## Rationale

1. **Prevents Package Bloat**: As more checks are added, each has its own isolated namespace
2. **Clear Dependencies**: Check dependencies are explicit via imports, not hidden in shared package
3. **Independent Testing**: Each check can be tested in isolation
4. **Constitutional Compliance**: Aligns with Package Granularity principle from constitution v1.14.0

## Design Decisions

### No Dashboard Check
The dashboard component does NOT have a dedicated check package. This is an intentional design decision - the dashboard does not require dedicated diagnostics.

### Check Registration
Checks are registered in `pkg/doctor/check/registry.go` using the `init()` pattern:

```go
// In pkg/doctor/checks/components/codeflare/codeflare.go
func init() {
    check.Registry().Register(&CodeFlareCheck{})
}
```

### Shared Utilities
Shared check utilities are located in:
- `pkg/doctor/checks/shared/results/` - Common result creation helpers

## Future Checks

When adding new diagnostic checks:

1. Create new directory under appropriate category: `pkg/doctor/checks/<category>/<check>/`
2. Implement `<check>.go` with check struct satisfying `check.Check` interface
3. Add `<check>_test.go` with gomega assertions using struct matchers (HaveField/MatchFields)
4. Register check in `init()` function
5. Follow constitutional standards:
   - Use centralized mocks from `pkg/util/test/mocks/`
   - Use JQ-based field access (no unstructured.NestedField)
   - Message constants at package level
   - godoc comments on exported identifiers
   - Explain WHY not WHAT in code comments

## Validation

Package isolation can be verified by:

```bash
# List all check packages
find pkg/doctor/checks -type d -mindepth 2 -maxdepth 2

# Verify each has implementation + test
find pkg/doctor/checks -name "*.go" -not -name "*_test.go"
find pkg/doctor/checks -name "*_test.go"
```

## References

- Constitution v1.14.0 - Principle XI: Package Granularity
- Constitution v1.14.0 - Diagnostic Check Package Isolation Standard
- `specs/002-constitution-alignment/spec.md` - User Story 1
