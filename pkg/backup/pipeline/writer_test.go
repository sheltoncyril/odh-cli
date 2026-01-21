package pipeline_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/lburgazzoli/odh-cli/pkg/backup/dependencies"
	"github.com/lburgazzoli/odh-cli/pkg/backup/pipeline"
	"github.com/lburgazzoli/odh-cli/pkg/util/iostreams"

	. "github.com/onsi/gomega"
)

const (
	testNamespace = "test-namespace"
	notebookName  = "test-notebook"
)

func createTestWorkload() *unstructured.Unstructured {
	workload := &unstructured.Unstructured{}
	workload.SetNamespace(testNamespace)
	workload.SetName(notebookName)

	return workload
}

func TestWriterStage(t *testing.T) {
	g := NewWithT(t)
	ctx := t.Context()

	t.Run("should write workload and dependencies", func(t *testing.T) {
		io := iostreams.NewIOStreams(nil, nil, nil)
		var writtenResources []string

		writeFunc := func(gvr schema.GroupVersionResource, obj *unstructured.Unstructured) error {
			writtenResources = append(writtenResources, fmt.Sprintf("%s/%s", obj.GetNamespace(), obj.GetName()))

			return nil
		}

		writer := &pipeline.WriterStage{
			WriteResource: writeFunc,
			IO:            io,
		}

		input := make(chan pipeline.WorkloadWithDeps, 1)

		// Create test workload with dependencies
		workload := createTestWorkload()
		dep1 := createTestWorkload()
		dep1.SetName("dep-1")
		dep2 := createTestWorkload()
		dep2.SetName("dep-2")

		input <- pipeline.WorkloadWithDeps{
			GVR:      schema.GroupVersionResource{Group: "test", Version: "v1", Resource: "tests"},
			Instance: workload,
			Dependencies: []dependencies.Dependency{
				{
					GVR:      schema.GroupVersionResource{Group: "v1", Version: "", Resource: "configmaps"},
					Resource: dep1,
				},
				{
					GVR:      schema.GroupVersionResource{Group: "v1", Version: "", Resource: "persistentvolumeclaims"},
					Resource: dep2,
				},
			},
		}
		close(input)

		// Run writer
		err := writer.Run(ctx, input)
		g.Expect(err).ToNot(HaveOccurred())

		// Verify workload and dependencies were written
		g.Expect(writtenResources).To(HaveLen(3))
		g.Expect(writtenResources[0]).To(Equal(fmt.Sprintf("%s/%s", testNamespace, notebookName)))
		g.Expect(writtenResources[1]).To(Equal(testNamespace + "/dep-1"))
		g.Expect(writtenResources[2]).To(Equal(testNamespace + "/dep-2"))
	})

	t.Run("should handle write errors gracefully", func(t *testing.T) {
		io := iostreams.NewIOStreams(nil, nil, nil)
		var writeCalls int

		writeFunc := func(gvr schema.GroupVersionResource, obj *unstructured.Unstructured) error {
			writeCalls++

			return errors.New("write error")
		}

		writer := &pipeline.WriterStage{
			WriteResource: writeFunc,
			IO:            io,
		}

		input := make(chan pipeline.WorkloadWithDeps, 1)

		workload := createTestWorkload()
		input <- pipeline.WorkloadWithDeps{
			GVR:          schema.GroupVersionResource{Group: "test", Version: "v1", Resource: "tests"},
			Instance:     workload,
			Dependencies: nil,
		}
		close(input)

		// Run writer
		err := writer.Run(ctx, input)
		g.Expect(err).ToNot(HaveOccurred())

		// Verify write was attempted
		g.Expect(writeCalls).To(Equal(1))
	})

	t.Run("should handle context cancellation", func(t *testing.T) {
		io := iostreams.NewIOStreams(nil, nil, nil)

		writeFunc := func(gvr schema.GroupVersionResource, obj *unstructured.Unstructured) error {
			return nil
		}

		writer := &pipeline.WriterStage{
			WriteResource: writeFunc,
			IO:            io,
		}

		input := make(chan pipeline.WorkloadWithDeps)

		cancelCtx, cancel := context.WithCancel(ctx)
		cancel() // Cancel immediately

		err := writer.Run(cancelCtx, input)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("writer cancelled"))
	})

	t.Run("should handle empty input", func(t *testing.T) {
		io := iostreams.NewIOStreams(nil, nil, nil)
		var writeCalls int

		writeFunc := func(gvr schema.GroupVersionResource, obj *unstructured.Unstructured) error {
			writeCalls++

			return nil
		}

		writer := &pipeline.WriterStage{
			WriteResource: writeFunc,
			IO:            io,
		}

		input := make(chan pipeline.WorkloadWithDeps)
		close(input)

		// Run writer
		err := writer.Run(ctx, input)
		g.Expect(err).ToNot(HaveOccurred())

		// Verify no writes occurred
		g.Expect(writeCalls).To(Equal(0))
	})

	t.Run("should handle timeout", func(t *testing.T) {
		io := iostreams.NewIOStreams(nil, nil, nil)

		writeFunc := func(gvr schema.GroupVersionResource, obj *unstructured.Unstructured) error {
			return nil
		}

		writer := &pipeline.WriterStage{
			WriteResource: writeFunc,
			IO:            io,
		}

		input := make(chan pipeline.WorkloadWithDeps)

		// Create context with very short timeout
		timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Nanosecond)
		defer cancel()

		// Wait for timeout
		time.Sleep(10 * time.Millisecond)

		err := writer.Run(timeoutCtx, input)
		g.Expect(err).To(HaveOccurred())
		g.Expect(err.Error()).To(ContainSubstring("writer cancelled"))
	})
}
