package rhbok

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"sigs.k8s.io/yaml"

	"github.com/lburgazzoli/odh-cli/pkg/migrate/action"
	"github.com/lburgazzoli/odh-cli/pkg/migrate/action/result"
	"github.com/lburgazzoli/odh-cli/pkg/resources"
)

func backupResources(
	ctx context.Context,
	target *action.ActionTarget,
) {
	step := target.Recorder.Child(
		"backup-resources",
		"Backup Kueue resources",
	)

	//nolint:mnd,gosec // 0o755 is standard directory permission for backup directory
	if err := os.MkdirAll(target.BackupPath, 0o755); err != nil {
		step.Complete(result.StepFailed, "Failed to create backup directory: %v", err)

		return
	}

	timestamp := time.Now().Format("20060102-150405")

	clusterQueues, err := target.Client.ListResources(ctx, resources.ClusterQueue.GVR())
	if err != nil {
		step.Complete(result.StepFailed, "Failed to list ClusterQueues: %v", err)

		return
	}

	if len(clusterQueues) > 0 {
		cqPath := filepath.Join(target.BackupPath, fmt.Sprintf("clusterqueues-%s.yaml", timestamp))
		if err := writeYAML(cqPath, clusterQueues); err != nil {
			step.Complete(result.StepFailed, "Failed to backup ClusterQueues: %v", err)

			return
		}
		target.IO.Errorf("✓ Backed up %d ClusterQueues to %s", len(clusterQueues), cqPath)
	}

	localQueues, err := target.Client.ListResources(ctx, resources.LocalQueue.GVR())
	if err != nil {
		step.Complete(result.StepFailed, "Failed to list LocalQueues: %v", err)

		return
	}

	if len(localQueues) > 0 {
		lqPath := filepath.Join(target.BackupPath, fmt.Sprintf("localqueues-%s.yaml", timestamp))
		if err := writeYAML(lqPath, localQueues); err != nil {
			step.Complete(result.StepFailed, "Failed to backup LocalQueues: %v", err)

			return
		}
		target.IO.Errorf("✓ Backed up %d LocalQueues to %s", len(localQueues), lqPath)
	}

	dsc, err := target.Client.GetDataScienceCluster(ctx)
	if err != nil {
		step.Complete(result.StepFailed, "Failed to get DataScienceCluster: %v", err)

		return
	}

	dscPath := filepath.Join(target.BackupPath, fmt.Sprintf("datasciencecluster-%s.yaml", timestamp))
	if err := writeYAML(dscPath, dsc); err != nil {
		step.Complete(result.StepFailed, "Failed to backup DataScienceCluster: %v", err)

		return
	}
	target.IO.Errorf("✓ Backed up DataScienceCluster to %s", dscPath)

	step.Complete(result.StepCompleted,
		"Backed up %d ClusterQueues, %d LocalQueues to %s",
		len(clusterQueues), len(localQueues), target.BackupPath)
}

func writeYAML(path string, data any) error {
	yamlBytes, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshaling YAML: %w", err)
	}

	//nolint:mnd,gosec // 0o644 is standard file permission for backup files
	if err := os.WriteFile(path, yamlBytes, 0o644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}
