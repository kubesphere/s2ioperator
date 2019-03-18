/*
Copyright 2019 The Kubesphere Authors.

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

package fake

import (
	v1alpha1 "github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeS2iBuilderTemplates implements S2iBuilderTemplateInterface
type FakeS2iBuilderTemplates struct {
	Fake *FakeDevopsV1alpha1
}

var s2ibuildertemplatesResource = schema.GroupVersionResource{Group: "devops.kubesphere.io", Version: "v1alpha1", Resource: "s2ibuildertemplates"}

var s2ibuildertemplatesKind = schema.GroupVersionKind{Group: "devops.kubesphere.io", Version: "v1alpha1", Kind: "S2iBuilderTemplate"}

// Get takes name of the s2iBuilderTemplate, and returns the corresponding s2iBuilderTemplate object, and an error if there is any.
func (c *FakeS2iBuilderTemplates) Get(name string, options v1.GetOptions) (result *v1alpha1.S2iBuilderTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(s2ibuildertemplatesResource, name), &v1alpha1.S2iBuilderTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.S2iBuilderTemplate), err
}

// List takes label and field selectors, and returns the list of S2iBuilderTemplates that match those selectors.
func (c *FakeS2iBuilderTemplates) List(opts v1.ListOptions) (result *v1alpha1.S2iBuilderTemplateList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(s2ibuildertemplatesResource, s2ibuildertemplatesKind, opts), &v1alpha1.S2iBuilderTemplateList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.S2iBuilderTemplateList{ListMeta: obj.(*v1alpha1.S2iBuilderTemplateList).ListMeta}
	for _, item := range obj.(*v1alpha1.S2iBuilderTemplateList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested s2iBuilderTemplates.
func (c *FakeS2iBuilderTemplates) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(s2ibuildertemplatesResource, opts))
}

// Create takes the representation of a s2iBuilderTemplate and creates it.  Returns the server's representation of the s2iBuilderTemplate, and an error, if there is any.
func (c *FakeS2iBuilderTemplates) Create(s2iBuilderTemplate *v1alpha1.S2iBuilderTemplate) (result *v1alpha1.S2iBuilderTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(s2ibuildertemplatesResource, s2iBuilderTemplate), &v1alpha1.S2iBuilderTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.S2iBuilderTemplate), err
}

// Update takes the representation of a s2iBuilderTemplate and updates it. Returns the server's representation of the s2iBuilderTemplate, and an error, if there is any.
func (c *FakeS2iBuilderTemplates) Update(s2iBuilderTemplate *v1alpha1.S2iBuilderTemplate) (result *v1alpha1.S2iBuilderTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(s2ibuildertemplatesResource, s2iBuilderTemplate), &v1alpha1.S2iBuilderTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.S2iBuilderTemplate), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeS2iBuilderTemplates) UpdateStatus(s2iBuilderTemplate *v1alpha1.S2iBuilderTemplate) (*v1alpha1.S2iBuilderTemplate, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(s2ibuildertemplatesResource, "status", s2iBuilderTemplate), &v1alpha1.S2iBuilderTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.S2iBuilderTemplate), err
}

// Delete takes name of the s2iBuilderTemplate and deletes it. Returns an error if one occurs.
func (c *FakeS2iBuilderTemplates) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(s2ibuildertemplatesResource, name), &v1alpha1.S2iBuilderTemplate{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeS2iBuilderTemplates) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(s2ibuildertemplatesResource, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.S2iBuilderTemplateList{})
	return err
}

// Patch applies the patch and returns the patched s2iBuilderTemplate.
func (c *FakeS2iBuilderTemplates) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.S2iBuilderTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(s2ibuildertemplatesResource, name, pt, data, subresources...), &v1alpha1.S2iBuilderTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.S2iBuilderTemplate), err
}
