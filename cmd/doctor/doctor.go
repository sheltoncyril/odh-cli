package doctor

import (
	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

const (
	cmdName  = "doctor"
	cmdShort = "Diagnose and validate OpenShift AI installation"
	cmdLong  = `
The doctor command provides diagnostic and validation tools for OpenShift AI clusters.

Available subcommands:
  lint  - Validate current cluster configuration or assess upgrade readiness
          Use without --version to validate current state
          Use with --version to assess upgrade readiness to a target version
`
)

// AddCommand adds the doctor command and its subcommands to the root command.
func AddCommand(root *cobra.Command, flags *genericclioptions.ConfigFlags) {
	cmd := &cobra.Command{
		Use:   cmdName,
		Short: cmdShort,
		Long:  cmdLong,
	}

	// Add subcommands
	AddLintCommand(cmd, flags)

	root.AddCommand(cmd)
}
