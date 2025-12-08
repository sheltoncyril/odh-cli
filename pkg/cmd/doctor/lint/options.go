package lint

import (
	"context"
	"errors"
	"fmt"

	"github.com/lburgazzoli/odh-cli/pkg/cmd/doctor"
	"github.com/lburgazzoli/odh-cli/pkg/doctor/check"
	"github.com/lburgazzoli/odh-cli/pkg/doctor/discovery"
	"github.com/lburgazzoli/odh-cli/pkg/doctor/version"
)

// Options contains options for the lint command.
type Options struct {
	*doctor.SharedOptions
}

// NewOptions creates a new Options with defaults.
func NewOptions(shared *doctor.SharedOptions) *Options {
	return &Options{
		SharedOptions: shared,
	}
}

// Complete populates Options and performs pre-validation setup.
func (o *Options) Complete() error {
	// Complete shared options (creates client)
	if err := o.SharedOptions.Complete(); err != nil {
		return fmt.Errorf("completing shared options: %w", err)
	}

	return nil
}

// Validate checks that all required options are valid.
func (o *Options) Validate() error {
	// Validate shared options
	if err := o.SharedOptions.Validate(); err != nil {
		return fmt.Errorf("validating shared options: %w", err)
	}

	return nil
}

// Run executes the lint command.
func (o *Options) Run(ctx context.Context) error {
	// Create context with timeout to prevent hanging on slow clusters
	ctx, cancel := context.WithTimeout(ctx, o.Timeout)
	defer cancel()

	// Detect cluster version
	clusterVersion, err := version.Detect(ctx, o.Client)
	if err != nil {
		return fmt.Errorf("detecting cluster version: %w", err)
	}

	_, _ = fmt.Fprintf(o.Out, "Detected OpenShift AI version: %s\n\n", clusterVersion)

	// Discover components and services
	_, _ = fmt.Fprint(o.Out, "Discovering OpenShift AI components and services...\n")
	components, err := discovery.DiscoverComponentsAndServices(ctx, o.Client)
	if err != nil {
		return fmt.Errorf("discovering components and services: %w", err)
	}
	_, _ = fmt.Fprintf(o.Out, "Found %d API groups\n", len(components))
	for _, comp := range components {
		_, _ = fmt.Fprintf(o.Out, "  - %s/%s (%d resources)\n", comp.APIGroup, comp.Version, len(comp.Resources))
	}
	_, _ = fmt.Fprintln(o.Out)

	// Discover workloads
	_, _ = fmt.Fprint(o.Out, "Discovering workload custom resources...\n")
	workloads, err := discovery.DiscoverWorkloads(ctx, o.Client)
	if err != nil {
		return fmt.Errorf("discovering workloads: %w", err)
	}
	_, _ = fmt.Fprintf(o.Out, "Found %d workload types\n", len(workloads))
	for _, gvr := range workloads {
		_, _ = fmt.Fprintf(o.Out, "  - %s/%s %s\n", gvr.Group, gvr.Version, gvr.Resource)
	}
	_, _ = fmt.Fprintln(o.Out)

	// Get the global check registry
	registry := check.GetGlobalRegistry()

	// Execute component and service checks (Resource: nil)
	_, _ = fmt.Fprint(o.Out, "Running component and service checks...\n")
	componentTarget := &check.CheckTarget{
		Client:         o.Client,
		CurrentVersion: clusterVersion, // For lint, current = target
		Version:        clusterVersion,
		Resource:       nil, // No specific resource for component/service checks
	}

	executor := check.NewExecutor(registry)

	// Execute all component checks
	componentResults, err := executor.ExecuteSelective(ctx, componentTarget, o.CheckSelector, check.CategoryComponent)
	if err != nil {
		// Log error but continue with other checks
		_, _ = fmt.Fprintf(o.ErrOut, "Warning: Failed to execute component checks: %v\n", err)
		componentResults = []check.CheckExecution{}
	}

	// Execute all service checks
	serviceResults, err := executor.ExecuteSelective(ctx, componentTarget, o.CheckSelector, check.CategoryService)
	if err != nil {
		// Log error but continue with other checks
		_, _ = fmt.Fprintf(o.ErrOut, "Warning: Failed to execute service checks: %v\n", err)
		serviceResults = []check.CheckExecution{}
	}

	// Execute all dependency checks
	dependencyResults, err := executor.ExecuteSelective(ctx, componentTarget, o.CheckSelector, check.CategoryDependency)
	if err != nil {
		// Log error but continue with other checks
		_, _ = fmt.Fprintf(o.ErrOut, "Warning: Failed to execute dependency checks: %v\n", err)
		dependencyResults = []check.CheckExecution{}
	}

	// Execute workload checks for each discovered workload instance
	_, _ = fmt.Fprint(o.Out, "Running workload checks...\n")
	var workloadResults []check.CheckExecution

	for _, gvr := range workloads {
		// List all instances of this workload type
		instances, err := o.Client.ListResources(ctx, gvr)
		if err != nil {
			// Skip workloads we can't access
			_, _ = fmt.Fprintf(o.ErrOut, "Warning: Failed to list %s: %v\n", gvr.Resource, err)

			continue
		}

		// Run workload checks for each instance
		for i := range instances {
			workloadTarget := &check.CheckTarget{
				Client:         o.Client,
				CurrentVersion: clusterVersion, // For lint, current = target
				Version:        clusterVersion,
				Resource:       &instances[i],
			}

			results, err := executor.ExecuteSelective(ctx, workloadTarget, o.CheckSelector, check.CategoryWorkload)
			if err != nil {
				return fmt.Errorf("executing workload checks: %w", err)
			}

			workloadResults = append(workloadResults, results...)
		}
	}

	// Group results by category
	resultsByCategory := map[check.CheckCategory][]check.CheckExecution{
		check.CategoryComponent:  componentResults,
		check.CategoryService:    serviceResults,
		check.CategoryDependency: dependencyResults,
		check.CategoryWorkload:   workloadResults,
	}

	// Filter results by minimum severity if specified
	filteredResults := doctor.FilterResultsBySeverity(resultsByCategory, o.MinSeverity)

	// Format and output results based on output format
	if err := o.formatAndOutputResults(filteredResults); err != nil {
		return err
	}

	// Determine exit code based on fail-on flags
	return o.determineExitCode(filteredResults)
}

// determineExitCode returns an error if fail-on conditions are met.
func (o *Options) determineExitCode(resultsByCategory map[check.CheckCategory][]check.CheckExecution) error {
	var hasCritical, hasWarning bool

	for _, results := range resultsByCategory {
		for _, result := range results {
			if result.Result.Severity != nil {
				//nolint:revive // exhaustive linter requires explicit SeverityInfo case
				switch *result.Result.Severity {
				case check.SeverityCritical:
					hasCritical = true
				case check.SeverityWarning:
					hasWarning = true
				case check.SeverityInfo:
					// Info doesn't affect exit code
				default:
					// Unknown severities don't affect exit code
				}
			}
		}
	}

	if o.FailOnCritical && hasCritical {
		return errors.New("critical findings detected")
	}

	if o.FailOnWarning && hasWarning {
		return errors.New("warning findings detected")
	}

	return nil
}

// formatAndOutputResults formats and outputs check results based on the output format.
func (o *Options) formatAndOutputResults(resultsByCategory map[check.CheckCategory][]check.CheckExecution) error {
	switch o.OutputFormat {
	case doctor.OutputFormatTable:
		return o.outputTable(resultsByCategory)
	case doctor.OutputFormatJSON:
		if err := doctor.OutputJSON(o.Out, resultsByCategory); err != nil {
			return fmt.Errorf("outputting JSON: %w", err)
		}

		return nil
	case doctor.OutputFormatYAML:
		if err := doctor.OutputYAML(o.Out, resultsByCategory); err != nil {
			return fmt.Errorf("outputting YAML: %w", err)
		}

		return nil
	default:
		return fmt.Errorf("unsupported output format: %s", o.OutputFormat)
	}
}

// outputTable outputs results in table format.
func (o *Options) outputTable(resultsByCategory map[check.CheckCategory][]check.CheckExecution) error {
	_, _ = fmt.Fprintln(o.Out)
	_, _ = fmt.Fprintln(o.Out, "Check Results:")
	_, _ = fmt.Fprintln(o.Out, "==============")

	if err := doctor.OutputTable(o.Out, resultsByCategory); err != nil {
		return fmt.Errorf("outputting table: %w", err)
	}

	return nil
}
