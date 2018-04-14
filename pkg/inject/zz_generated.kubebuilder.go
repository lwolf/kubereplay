/*
Copyright 2017 Sergey Nuzhdin.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/package inject

import (
	"github.com/kubernetes-sigs/kubebuilder/pkg/inject/run"
	kubereplayv1alpha1 "github.com/lwolf/kubereplay/pkg/apis/kubereplay/v1alpha1"
	rscheme "github.com/lwolf/kubereplay/pkg/client/clientset/versioned/scheme"
	"github.com/lwolf/kubereplay/pkg/controller/harvester"
	"github.com/lwolf/kubereplay/pkg/controller/refinery"
	"github.com/lwolf/kubereplay/pkg/inject/args"
	appsv1 "k8s.io/api/apps/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
)

func init() {
	rscheme.AddToScheme(scheme.Scheme)

	// Inject Informers
	Inject = append(Inject, func(arguments args.InjectArgs) error {
		Injector.ControllerManager = arguments.ControllerManager

		if err := arguments.ControllerManager.AddInformerProvider(&kubereplayv1alpha1.Harvester{}, arguments.Informers.Kubereplay().V1alpha1().Harvesters()); err != nil {
			return err
		}
		if err := arguments.ControllerManager.AddInformerProvider(&kubereplayv1alpha1.Refinery{}, arguments.Informers.Kubereplay().V1alpha1().Refineries()); err != nil {
			return err
		}

		// Add Kubernetes informers
		if err := arguments.ControllerManager.AddInformerProvider(&appsv1.Deployment{}, arguments.KubernetesInformers.Apps().V1().Deployments()); err != nil {
			return err
		}

		if c, err := harvester.ProvideController(arguments); err != nil {
			return err
		} else {
			arguments.ControllerManager.AddController(c)
		}
		if c, err := refinery.ProvideController(arguments); err != nil {
			return err
		} else {
			arguments.ControllerManager.AddController(c)
		}
		return nil
	})

	// Inject CRDs
	Injector.CRDs = append(Injector.CRDs, &kubereplayv1alpha1.HarvesterCRD)
	Injector.CRDs = append(Injector.CRDs, &kubereplayv1alpha1.RefineryCRD)
	// Inject PolicyRules
	Injector.PolicyRules = append(Injector.PolicyRules, rbacv1.PolicyRule{
		APIGroups: []string{"kubereplay.lwolf.org"},
		Resources: []string{"*"},
		Verbs:     []string{"*"},
	})
	Injector.PolicyRules = append(Injector.PolicyRules, rbacv1.PolicyRule{
		APIGroups: []string{
			"apps",
		},
		Resources: []string{
			"deployments",
		},
		Verbs: []string{
			"create", "delete", "get", "list", "patch", "update", "watch",
		},
	})
	// Inject GroupVersions
	Injector.GroupVersions = append(Injector.GroupVersions, schema.GroupVersion{
		Group:   "kubereplay.lwolf.org",
		Version: "v1alpha1",
	})
	Injector.RunFns = append(Injector.RunFns, func(arguments run.RunArguments) error {
		Injector.ControllerManager.RunInformersAndControllers(arguments)
		return nil
	})
}
