package helpers

import (
	"testing"
	"time"

	"github.com/go-test/deep"
	appsv1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func refToIntstr(i intstr.IntOrString) *intstr.IntOrString {
	return &i
}

var mockContainer = apiv1.Container{
	Name:  "test",
	Image: "docker.io/container",
}

var mockDeployment = &appsv1.Deployment{
	ObjectMeta: metav1.ObjectMeta{
		Name:                       "echoheaders",
		GenerateName:               "",
		Namespace:                  "default",
		SelfLink:                   "/apis/apps/v1/namespaces/default/deployments/echoheaders",
		UID:                        "e1d8b5b3-6956-11e8-9871-080027cf8168",
		ResourceVersion:            "64735",
		Generation:                 2,
		CreationTimestamp:          metav1.NewTime(time.Now().UTC()),
		DeletionTimestamp:          nil,
		DeletionGracePeriodSeconds: nil,
		Labels: map[string]string{
			"app":    "kubereplay",
			"module": "test",
		},
		Annotations: map[string]string{
			"deployment.kubernetes.io/revision": "1",
			"my-custom-annotations":             "annotation-value",
			"kubereplay.lwolf.org/mode":         "green",
			"kubereplay.lwolf.org/replicas":     "3",
			"kubereplay.lwolf.org/shadow":       "echoheaders-gor",
		},
		OwnerReferences: []metav1.OwnerReference{},
		Finalizers:      []string{},
		ClusterName:     "",
		Initializers:    nil,
	},
	Spec: appsv1.DeploymentSpec{
		Replicas: int32Ptr(3),
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app":    "kubereplay",
				"module": "test",
			},
			MatchExpressions: []metav1.LabelSelectorRequirement{},
		},
		Template: apiv1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Name:                       "Name",
				GenerateName:               "GenerateName",
				Namespace:                  "Namespace",
				SelfLink:                   "SelfLink",
				UID:                        "UID",
				ResourceVersion:            "ResourceVersion",
				Generation:                 0,
				CreationTimestamp:          metav1.NewTime(time.Now().UTC()),
				DeletionTimestamp:          nil,
				DeletionGracePeriodSeconds: nil,
				Labels: map[string]string{
					"app":    "kubereplay",
					"module": "test",
				},
				Annotations:     map[string]string{},
				OwnerReferences: []metav1.OwnerReference{},
				Finalizers:      []string{},
				ClusterName:     "",
				Initializers:    nil,
			},
			Spec: apiv1.PodSpec{
				Volumes:                       []apiv1.Volume{},
				Containers:                    []apiv1.Container{mockContainer},
				RestartPolicy:                 apiv1.RestartPolicyAlways,
				TerminationGracePeriodSeconds: int64Ptr(30),
				ActiveDeadlineSeconds:         nil,
				DNSPolicy:                     apiv1.DNSClusterFirst,
				NodeSelector:                  map[string]string{},
				ServiceAccountName:            "",
				DeprecatedServiceAccount:      "",
				NodeName:                      "",
				HostNetwork:                   false,
				HostPID:                       false,
				HostIPC:                       false,
				SecurityContext: &apiv1.PodSecurityContext{
					SELinuxOptions:     nil,
					RunAsUser:          nil,
					RunAsNonRoot:       nil,
					SupplementalGroups: []int64{},
					FSGroup:            nil,
					RunAsGroup:         nil,
				},
				ImagePullSecrets:             []apiv1.LocalObjectReference{},
				Hostname:                     "",
				Subdomain:                    "",
				Affinity:                     nil,
				SchedulerName:                "default-scheduler",
				InitContainers:               []apiv1.Container{},
				AutomountServiceAccountToken: nil,
				Tolerations:                  []apiv1.Toleration{},
				HostAliases:                  []apiv1.HostAlias{},
				PriorityClassName:            "",
				Priority:                     nil,
				DNSConfig:                    nil,
				ShareProcessNamespace:        nil,
			},
		},
		Strategy: appsv1.DeploymentStrategy{
			Type: appsv1.RollingUpdateDeploymentStrategyType,
			RollingUpdate: &appsv1.RollingUpdateDeployment{
				MaxUnavailable: refToIntstr(intstr.FromInt(25)),
				MaxSurge:       refToIntstr(intstr.FromString("25%")),
			},
		},
		MinReadySeconds:         0,
		RevisionHistoryLimit:    int32Ptr(10),
		Paused:                  false,
		ProgressDeadlineSeconds: int32Ptr(600),
	},
	Status: appsv1.DeploymentStatus{
		ObservedGeneration:  2,
		Replicas:            3,
		UpdatedReplicas:     3,
		AvailableReplicas:   0,
		UnavailableReplicas: 3,
		Conditions:          []appsv1.DeploymentCondition{},
		ReadyReplicas:       0,
		CollisionCount:      nil,
	},
}

var cleanDeployment = &appsv1.Deployment{
	ObjectMeta: metav1.ObjectMeta{
		Namespace: "default",
		Labels: map[string]string{
			"app":    "kubereplay",
			"module": "test",
		},
		Annotations: map[string]string{
			"my-custom-annotations": "annotation-value",
		},
	},
	Spec: appsv1.DeploymentSpec{
		Replicas: int32Ptr(3),
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app":    "kubereplay",
				"module": "test",
			},
			MatchExpressions: []metav1.LabelSelectorRequirement{},
		},
		Template: apiv1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Namespace:                  "Namespace",
				DeletionGracePeriodSeconds: nil,
				Labels: map[string]string{
					"app":    "kubereplay",
					"module": "test",
				},
				Annotations: map[string]string{},
			},
			Spec: apiv1.PodSpec{
				Volumes:                       []apiv1.Volume{},
				Containers:                    []apiv1.Container{mockContainer},
				RestartPolicy:                 apiv1.RestartPolicyAlways,
				TerminationGracePeriodSeconds: int64Ptr(30),
				ActiveDeadlineSeconds:         nil,
				DNSPolicy:                     apiv1.DNSClusterFirst,
				NodeSelector:                  map[string]string{},
				ServiceAccountName:            "",
				NodeName:                      "",
				HostNetwork:                   false,
				SecurityContext: &apiv1.PodSecurityContext{
					SELinuxOptions:     nil,
					RunAsUser:          nil,
					RunAsNonRoot:       nil,
					SupplementalGroups: []int64{},
					FSGroup:            nil,
					RunAsGroup:         nil,
				},
				ImagePullSecrets: []apiv1.LocalObjectReference{},
				Affinity:         nil,
				SchedulerName:    "default-scheduler",
				InitContainers:   []apiv1.Container{},
				Tolerations:      []apiv1.Toleration{},
				HostAliases:      []apiv1.HostAlias{},
			},
		},
		Strategy: appsv1.DeploymentStrategy{
			Type: appsv1.RollingUpdateDeploymentStrategyType,
			RollingUpdate: &appsv1.RollingUpdateDeployment{
				MaxUnavailable: refToIntstr(intstr.FromInt(25)),
				MaxSurge:       refToIntstr(intstr.FromString("25%")),
			},
		},
		MinReadySeconds:         0,
		RevisionHistoryLimit:    int32Ptr(10),
		ProgressDeadlineSeconds: int32Ptr(600),
	},
}

func TestCleanupDeployment(t *testing.T) {

	t.Run("isOk", func(t *testing.T) {
		got := CleanupDeployment(mockDeployment)
		want := cleanDeployment
		if diff := deep.Equal(got, want); diff != nil {
			t.Error(diff)
		}
	})
	t.Run("isIgnoringGoreplayContainer", func(t *testing.T) {
		d := mockDeployment.DeepCopy()
		d.Spec.Template.Spec.Containers = append(
			d.Spec.Template.Spec.Containers,
			apiv1.Container{Name: "goreplay"},
		)
		got := CleanupDeployment(d)
		want := cleanDeployment
		if diff := deep.Equal(got, want); diff != nil {
			t.Error(diff)
		}
	})
}
