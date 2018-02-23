/*
Copyright 2018 The Kubernetes Authors.

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

// This file was automatically generated by lister-gen

package v1alpha1

import (
	v1alpha1 "github.com/lwolf/kube-replay/pkg/apis/replay/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// RefineryLister helps list Refineries.
type RefineryLister interface {
	// List lists all Refineries in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.Refinery, err error)
	// Refineries returns an object that can list and get Refineries.
	Refineries(namespace string) RefineryNamespaceLister
	RefineryListerExpansion
}

// refineryLister implements the RefineryLister interface.
type refineryLister struct {
	indexer cache.Indexer
}

// NewRefineryLister returns a new RefineryLister.
func NewRefineryLister(indexer cache.Indexer) RefineryLister {
	return &refineryLister{indexer: indexer}
}

// List lists all Refineries in the indexer.
func (s *refineryLister) List(selector labels.Selector) (ret []*v1alpha1.Refinery, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Refinery))
	})
	return ret, err
}

// Refineries returns an object that can list and get Refineries.
func (s *refineryLister) Refineries(namespace string) RefineryNamespaceLister {
	return refineryNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// RefineryNamespaceLister helps list and get Refineries.
type RefineryNamespaceLister interface {
	// List lists all Refineries in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.Refinery, err error)
	// Get retrieves the Refinery from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.Refinery, error)
	RefineryNamespaceListerExpansion
}

// refineryNamespaceLister implements the RefineryNamespaceLister
// interface.
type refineryNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Refineries in the indexer for a given namespace.
func (s refineryNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.Refinery, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Refinery))
	})
	return ret, err
}

// Get retrieves the Refinery from the indexer for a given namespace and name.
func (s refineryNamespaceLister) Get(name string) (*v1alpha1.Refinery, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("refinery"), name)
	}
	return obj.(*v1alpha1.Refinery), nil
}
