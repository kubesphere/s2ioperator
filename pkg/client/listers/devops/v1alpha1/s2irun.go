/*
Copyright 2019 The KubeSphere Authors.

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
// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// S2iRunLister helps list S2iRuns.
type S2iRunLister interface {
	// List lists all S2iRuns in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.S2iRun, err error)
	// S2iRuns returns an object that can list and get S2iRuns.
	S2iRuns(namespace string) S2iRunNamespaceLister
	S2iRunListerExpansion
}

// s2iRunLister implements the S2iRunLister interface.
type s2iRunLister struct {
	indexer cache.Indexer
}

// NewS2iRunLister returns a new S2iRunLister.
func NewS2iRunLister(indexer cache.Indexer) S2iRunLister {
	return &s2iRunLister{indexer: indexer}
}

// List lists all S2iRuns in the indexer.
func (s *s2iRunLister) List(selector labels.Selector) (ret []*v1alpha1.S2iRun, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.S2iRun))
	})
	return ret, err
}

// S2iRuns returns an object that can list and get S2iRuns.
func (s *s2iRunLister) S2iRuns(namespace string) S2iRunNamespaceLister {
	return s2iRunNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// S2iRunNamespaceLister helps list and get S2iRuns.
type S2iRunNamespaceLister interface {
	// List lists all S2iRuns in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.S2iRun, err error)
	// Get retrieves the S2iRun from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.S2iRun, error)
	S2iRunNamespaceListerExpansion
}

// s2iRunNamespaceLister implements the S2iRunNamespaceLister
// interface.
type s2iRunNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all S2iRuns in the indexer for a given namespace.
func (s s2iRunNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.S2iRun, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.S2iRun))
	})
	return ret, err
}

// Get retrieves the S2iRun from the indexer for a given namespace and name.
func (s s2iRunNamespaceLister) Get(name string) (*v1alpha1.S2iRun, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("s2irun"), name)
	}
	return obj.(*v1alpha1.S2iRun), nil
}
