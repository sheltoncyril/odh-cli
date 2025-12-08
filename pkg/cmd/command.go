package cmd

import (
	"context"

	"github.com/spf13/pflag"
)

// Command represents a standard command with lifecycle methods.
//
// All commands should implement this interface to ensure consistent
// behavior across the CLI. The interface enforces a three-phase execution:
//  1. Complete: Populate derived fields and perform pre-validation setup
//  2. Validate: Check that all required fields are valid
//  3. Run: Execute the command logic
//
// Commands should also implement AddFlags to register their flags with
// a FlagSet, enabling testability and decoupling from Cobra.
type Command interface {
	// Complete populates derived fields and performs pre-validation setup.
	// This is called after flags are parsed but before validation.
	// Example: parsing version strings, creating clients, etc.
	Complete() error

	// Validate checks that all required fields are valid.
	// This is called after Complete and before Run.
	// Example: checking required flags, validating patterns, etc.
	Validate() error

	// Run executes the command logic.
	// This is called after Complete and Validate have succeeded.
	Run(ctx context.Context) error

	// AddFlags registers command-specific flags with the provided FlagSet.
	// This enables testing flag registration independently of Cobra.
	AddFlags(fs *pflag.FlagSet)
}
