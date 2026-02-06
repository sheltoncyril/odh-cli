package client

import (
	"context"

	operatorsv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	olmclientset "github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/versioned"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// olmReaderImpl wraps an OLM clientset to provide read-only access.
// Nil-safe: returns empty results when the underlying client is nil.
type olmReaderImpl struct {
	client olmclientset.Interface
}

func newOLMReader(client olmclientset.Interface) OLMReader {
	return &olmReaderImpl{client: client}
}

func (r *olmReaderImpl) Available() bool {
	return r.client != nil
}

func (r *olmReaderImpl) Subscriptions(namespace string) SubscriptionReader {
	if r.client == nil {
		return &nilSubscriptionReader{}
	}

	return r.client.OperatorsV1alpha1().Subscriptions(namespace)
}

func (r *olmReaderImpl) ClusterServiceVersions(namespace string) CSVReader {
	if r.client == nil {
		return &nilCSVReader{}
	}

	return r.client.OperatorsV1alpha1().ClusterServiceVersions(namespace)
}

// nilSubscriptionReader returns empty results when OLM is not available.
type nilSubscriptionReader struct{}

func (n *nilSubscriptionReader) List(_ context.Context, _ metav1.ListOptions) (*operatorsv1alpha1.SubscriptionList, error) {
	return &operatorsv1alpha1.SubscriptionList{}, nil
}

func (n *nilSubscriptionReader) Get(_ context.Context, _ string, _ metav1.GetOptions) (*operatorsv1alpha1.Subscription, error) {
	return nil, nil
}

// nilCSVReader returns empty results when OLM is not available.
type nilCSVReader struct{}

func (n *nilCSVReader) List(_ context.Context, _ metav1.ListOptions) (*operatorsv1alpha1.ClusterServiceVersionList, error) {
	return &operatorsv1alpha1.ClusterServiceVersionList{}, nil
}

func (n *nilCSVReader) Get(_ context.Context, _ string, _ metav1.GetOptions) (*operatorsv1alpha1.ClusterServiceVersion, error) {
	return nil, nil
}
