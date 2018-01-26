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

// This file was automatically generated by informer-gen

package internalversion

import (
	time "time"

	replay "github.com/lwolf/kube-replay/pkg/apis/replay"
	clientset_internalversion "github.com/lwolf/kube-replay/pkg/client/clientset/internalversion"
	internalinterfaces "github.com/lwolf/kube-replay/pkg/client/informers/internalversion/internalinterfaces"
	internalversion "github.com/lwolf/kube-replay/pkg/client/listers/replay/internalversion"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// SiloInformer provides access to a shared informer and lister for
// Silos.
type SiloInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() internalversion.SiloLister
}

type siloInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewSiloInformer constructs a new informer for Silo type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewSiloInformer(client clientset_internalversion.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredSiloInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredSiloInformer constructs a new informer for Silo type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredSiloInformer(client clientset_internalversion.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.Kubereplay().Silos(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.Kubereplay().Silos(namespace).Watch(options)
			},
		},
		&replay.Silo{},
		resyncPeriod,
		indexers,
	)
}

func (f *siloInformer) defaultInformer(client clientset_internalversion.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredSiloInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *siloInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&replay.Silo{}, f.defaultInformer)
}

func (f *siloInformer) Lister() internalversion.SiloLister {
	return internalversion.NewSiloLister(f.Informer().GetIndexer())
}