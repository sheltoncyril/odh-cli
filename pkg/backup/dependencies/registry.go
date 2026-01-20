package dependencies

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Registry holds all registered dependency resolvers.
type Registry struct {
	resolvers []Resolver
}

// NewRegistry creates a new resolver registry.
func NewRegistry() *Registry {
	return &Registry{
		resolvers: make([]Resolver, 0),
	}
}

// Register adds a resolver to the registry.
func (r *Registry) Register(resolver Resolver) {
	r.resolvers = append(r.resolvers, resolver)
}

// MustRegister registers a resolver and panics if the resolver is nil.
// Use this for explicit registration in command construction.
func (r *Registry) MustRegister(resolver Resolver) {
	if resolver == nil {
		panic("cannot register nil resolver")
	}
	r.Register(resolver)
}

// GetResolver finds the appropriate resolver for the given GVR.
func (r *Registry) GetResolver(gvr schema.GroupVersionResource) (Resolver, error) {
	for _, resolver := range r.resolvers {
		if resolver.CanHandle(gvr) {
			return resolver, nil
		}
	}

	return nil, fmt.Errorf("no dependency resolver registered for %s", gvr.String())
}
