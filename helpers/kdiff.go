package helpers

import (
	appsv1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

func cleanupObjectMeta(meta *metav1.ObjectMeta) metav1.ObjectMeta {
	annotationPrefixesToRemove := [2]string{
		"kubereplay.lwolf.org",
		"deployment.kubernetes.io/revision",
	}
	annotations := make(map[string]string)
	for key, value := range meta.Annotations {
		shouldAdd := true
		for _, i := range annotationPrefixesToRemove {
			if strings.HasPrefix(key, i) {
				shouldAdd = false
			}
		}
		if shouldAdd {
			annotations[key] = value
		}

	}

	return metav1.ObjectMeta{
		Namespace:   meta.Namespace,
		Annotations: annotations,
		Labels:      meta.Labels,
	}
}

func filteredContainers(containers []apiv1.Container) []apiv1.Container {
	var result []apiv1.Container
	for _, c := range containers {
		if c.Name != "goreplay" {
			result = append(result, c)
		}
	}
	return result
}

func cleanupPodSpec(in *apiv1.PodSpec) apiv1.PodSpec {
	return apiv1.PodSpec{
		Volumes:                       in.Volumes,
		Containers:                    filteredContainers(in.Containers),
		RestartPolicy:                 in.RestartPolicy,
		TerminationGracePeriodSeconds: in.TerminationGracePeriodSeconds,
		ActiveDeadlineSeconds:         in.ActiveDeadlineSeconds,
		DNSPolicy:                     in.DNSPolicy,
		NodeSelector:                  in.NodeSelector,
		ServiceAccountName:            in.ServiceAccountName,
		NodeName:                      in.NodeName,
		HostNetwork:                   in.HostNetwork,
		SecurityContext:               in.SecurityContext,
		ImagePullSecrets:              in.ImagePullSecrets,
		Affinity:                      in.Affinity,
		SchedulerName:                 in.SchedulerName,
		InitContainers:                in.InitContainers,
		Tolerations:                   in.Tolerations,
		HostAliases:                   in.HostAliases,
	}
}

func CleanupDeployment(d *appsv1.Deployment) *appsv1.Deployment {
	dep := d.DeepCopy()
	return &appsv1.Deployment{
		ObjectMeta: cleanupObjectMeta(&dep.ObjectMeta),
		Spec: appsv1.DeploymentSpec{
			Replicas: dep.Spec.Replicas,
			Selector: dep.Spec.Selector,
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: cleanupObjectMeta(&dep.Spec.Template.ObjectMeta),
				Spec:       cleanupPodSpec(&dep.Spec.Template.Spec),
			},
			Strategy:                dep.Spec.Strategy,
			MinReadySeconds:         dep.Spec.MinReadySeconds,
			RevisionHistoryLimit:    dep.Spec.RevisionHistoryLimit,
			ProgressDeadlineSeconds: dep.Spec.ProgressDeadlineSeconds,
		},
	}
}
