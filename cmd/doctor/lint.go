package doctor

import (
	"fmt"

	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"

	"github.com/lburgazzoli/odh-cli/pkg/cmd/doctor/lint"
)

const (
	lintCmdName  = "lint"
	lintCmdShort = "Validate current OpenShift AI installation or assess upgrade readiness"
	lintCmdLong  = `
Validates the current OpenShift AI installation or assesses upgrade readiness.

LINT MODE (without --version):
  Validates the current cluster state and reports configuration issues.

UPGRADE MODE (with --version):
  Assesses upgrade readiness by comparing current version against target version.

The lint command performs comprehensive validation across four categories:
  - Components: Core OpenShift AI components (Dashboard, Workbenches, etc.)
  - Services: Platform services (OAuth, monitoring, etc.)
  - Dependencies: External dependencies (CertManager, Kueue, etc.)
  - Workloads: User-created custom resources (Notebooks, InferenceServices, etc.)

Each issue is reported with:
  - Severity level (Critical, Warning, Info)
  - Detailed description of the problem
  - Remediation guidance for fixing the issue

Examples:
  # Validate current cluster state
  kubectl odh doctor lint

  # Assess upgrade readiness for version 3.0
  kubectl odh doctor lint --version 3.0

  # Validate with JSON output
  kubectl odh doctor lint -o json

  # Validate only component checks
  kubectl odh doctor lint --checks "components/*"
`
	lintCmdExample = `
  # Validate current cluster state
  kubectl odh doctor lint

  # Assess upgrade readiness for version 3.0
  kubectl odh doctor lint --version 3.0

  # Output results in JSON format
  kubectl odh doctor lint -o json

  # Run only dashboard-related checks
  kubectl odh doctor lint --checks "*dashboard*"

  # Check upgrade to version 3.1 with critical issues only
  kubectl odh doctor lint --version 3.1 --severity critical
`
)

// AddLintCommand adds the lint subcommand to the doctor command.
func AddLintCommand(parent *cobra.Command, flags *genericclioptions.ConfigFlags) {
	streams := genericiooptions.IOStreams{
		In:     parent.InOrStdin(),
		Out:    parent.OutOrStdout(),
		ErrOut: parent.ErrOrStderr(),
	}

	// Create command using new pattern (FR-014: SharedOptions initialized internally)
	command := lint.NewCommand(streams)

	// Use the ConfigFlags from parent instead of creating new ones
	command.ConfigFlags = flags

	cmd := &cobra.Command{
		Use:     lintCmdName,
		Short:   lintCmdShort,
		Long:    lintCmdLong,
		Example: lintCmdExample,
		RunE: func(cmd *cobra.Command, _ []string) error {
			// Complete phase
			if err := command.Complete(); err != nil {
				return fmt.Errorf("completing lint command: %w", err)
			}

			// Validate phase
			if err := command.Validate(); err != nil {
				return fmt.Errorf("validating lint command: %w", err)
			}

			// Run phase
			return command.Run(cmd.Context())
		},
	}

	// Register flags using AddFlags method
	command.AddFlags(cmd.Flags())

	parent.AddCommand(cmd)
}
