/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	json "encoding/json"
	"fmt"
	"time"

	v1alpha1 "github.com/puller-io/puller/pkg/apis/puller/v1alpha1"
	pullerv1alpha1 "github.com/puller-io/puller/pkg/generated/applyconfiguration/puller/v1alpha1"
	scheme "github.com/puller-io/puller/pkg/generated/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// PullersGetter has a method to return a PullerInterface.
// A group's client should implement this interface.
type PullersGetter interface {
	Pullers() PullerInterface
}

// PullerInterface has methods to work with Puller resources.
type PullerInterface interface {
	Create(ctx context.Context, puller *v1alpha1.Puller, opts v1.CreateOptions) (*v1alpha1.Puller, error)
	Update(ctx context.Context, puller *v1alpha1.Puller, opts v1.UpdateOptions) (*v1alpha1.Puller, error)
	UpdateStatus(ctx context.Context, puller *v1alpha1.Puller, opts v1.UpdateOptions) (*v1alpha1.Puller, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.Puller, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.PullerList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Puller, err error)
	Apply(ctx context.Context, puller *pullerv1alpha1.PullerApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Puller, err error)
	ApplyStatus(ctx context.Context, puller *pullerv1alpha1.PullerApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Puller, err error)
	PullerExpansion
}

// pullers implements PullerInterface
type pullers struct {
	client rest.Interface
}

// newPullers returns a Pullers
func newPullers(c *PullerV1alpha1Client) *pullers {
	return &pullers{
		client: c.RESTClient(),
	}
}

// Get takes name of the puller, and returns the corresponding puller object, and an error if there is any.
func (c *pullers) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Puller, err error) {
	result = &v1alpha1.Puller{}
	err = c.client.Get().
		Resource("pullers").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Pullers that match those selectors.
func (c *pullers) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.PullerList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.PullerList{}
	err = c.client.Get().
		Resource("pullers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested pullers.
func (c *pullers) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("pullers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a puller and creates it.  Returns the server's representation of the puller, and an error, if there is any.
func (c *pullers) Create(ctx context.Context, puller *v1alpha1.Puller, opts v1.CreateOptions) (result *v1alpha1.Puller, err error) {
	result = &v1alpha1.Puller{}
	err = c.client.Post().
		Resource("pullers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(puller).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a puller and updates it. Returns the server's representation of the puller, and an error, if there is any.
func (c *pullers) Update(ctx context.Context, puller *v1alpha1.Puller, opts v1.UpdateOptions) (result *v1alpha1.Puller, err error) {
	result = &v1alpha1.Puller{}
	err = c.client.Put().
		Resource("pullers").
		Name(puller.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(puller).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *pullers) UpdateStatus(ctx context.Context, puller *v1alpha1.Puller, opts v1.UpdateOptions) (result *v1alpha1.Puller, err error) {
	result = &v1alpha1.Puller{}
	err = c.client.Put().
		Resource("pullers").
		Name(puller.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(puller).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the puller and deletes it. Returns an error if one occurs.
func (c *pullers) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("pullers").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *pullers) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("pullers").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched puller.
func (c *pullers) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Puller, err error) {
	result = &v1alpha1.Puller{}
	err = c.client.Patch(pt).
		Resource("pullers").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied puller.
func (c *pullers) Apply(ctx context.Context, puller *pullerv1alpha1.PullerApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Puller, err error) {
	if puller == nil {
		return nil, fmt.Errorf("puller provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(puller)
	if err != nil {
		return nil, err
	}
	name := puller.Name
	if name == nil {
		return nil, fmt.Errorf("puller.Name must be provided to Apply")
	}
	result = &v1alpha1.Puller{}
	err = c.client.Patch(types.ApplyPatchType).
		Resource("pullers").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *pullers) ApplyStatus(ctx context.Context, puller *pullerv1alpha1.PullerApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Puller, err error) {
	if puller == nil {
		return nil, fmt.Errorf("puller provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(puller)
	if err != nil {
		return nil, err
	}

	name := puller.Name
	if name == nil {
		return nil, fmt.Errorf("puller.Name must be provided to Apply")
	}

	result = &v1alpha1.Puller{}
	err = c.client.Patch(types.ApplyPatchType).
		Resource("pullers").
		Name(*name).
		SubResource("status").
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
