package apis

import (
	"github.com/kubernetes-sigs/kubebuilder/pkg/builders"
	"github.com/lwolf/kubereplay/pkg/apis/kubereplay"
	kubereplayv1alpha1 "github.com/lwolf/kubereplay/pkg/apis/kubereplay/v1alpha1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type MetaData struct{}

var APIMeta = MetaData{}

// GetAllApiBuilders returns all known APIGroupBuilders
// so they can be registered with the apiserver
func (MetaData) GetAllApiBuilders() []*builders.APIGroupBuilder {
	return []*builders.APIGroupBuilder{
		GetKubereplayAPIBuilder(),
	}
}

// GetCRDs returns all the CRDs for known resource types
func (MetaData) GetCRDs() []v1beta1.CustomResourceDefinition {
	return []v1beta1.CustomResourceDefinition{
		kubereplayv1alpha1.HarvesterCRD,
		kubereplayv1alpha1.RefineryCRD,
	}
}

func (MetaData) GetRules() []rbacv1.PolicyRule {
	return []rbacv1.PolicyRule{
		{
			APIGroups: []string{"kubereplay.lwolf.org"},
			Resources: []string{"*"},
			Verbs:     []string{"*"},
		},
	}
}

func (MetaData) GetGroupVersions() []schema.GroupVersion {
	return []schema.GroupVersion{
		{
			Group:   "kubereplay.lwolf.org",
			Version: "v1alpha1",
		},
	}
}

var kubereplayApiGroup = builders.NewApiGroupBuilder(
	"kubereplay.lwolf.org",
	"github.com/lwolf/kubereplay/pkg/apis/kubereplay").
	WithUnVersionedApi(kubereplay.ApiVersion).
	WithVersionedApis(
		kubereplayv1alpha1.ApiVersion,
	).
	WithRootScopedKinds()

func GetKubereplayAPIBuilder() *builders.APIGroupBuilder {
	return kubereplayApiGroup
}
