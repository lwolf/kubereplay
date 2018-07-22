package helpers

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	kubereplayv1alpha1 "github.com/lwolf/kubereplay/pkg/apis/kubereplay/v1alpha1"
	appsv1 "k8s.io/api/apps/v1beta2"
	apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func AppLabels(name string) map[string]string {
	return map[string]string{
		"kubereplay-app": name,
	}
}

func fileSiloToArgs(spec *kubereplayv1alpha1.FileSilo) (*[]string, error) {
	var args []string
	if spec.Filename == "" {
		return nil, errors.New("filename is required for file output mode")
	} else {
		args = append(args, "-output-file")
		args = append(args, spec.Filename)
	}
	if spec.Append == true {
		args = append(args, "-output-file-append")
	}

	if spec.FlushInterval != "" {
		args = append(args, "-output-file-flush-interval")
		args = append(args, spec.FlushInterval)
	}

	if spec.QueueSize > 0 {
		args = append(args, "-output-file-queue-limit")
		args = append(args, fmt.Sprint(spec.QueueSize))
	}

	if spec.FileLimit != "" {
		args = append(args, "-output-file-size-limit")
		args = append(args, spec.FileLimit)
	}

	return &args, nil
}

func tcpSiloToArgs(spec *kubereplayv1alpha1.TcpSilo) (*[]string, error) {
	var args []string
	if spec.Uri == "" {
		return nil, errors.New("uri is required for tcp output mode")
	}
	args = append(args, "-output-tcp")
	args = append(args, spec.Uri)
	return &args, nil
}

func stdoutSiloToArgs(spec *kubereplayv1alpha1.StdoutSilo) (*[]string, error) {
	var args []string
	args = append(args, "-output-stdout")
	return &args, nil
}

func httpSiloToArgs(spec *kubereplayv1alpha1.HttpSilo) (*[]string, error) {
	var args []string
	if spec.Uri == "" {
		return nil, errors.New("uri is required for http output mode")
	}
	args = append(args, "-output-http")
	args = append(args, spec.Uri)

	if spec.Debug == true {
		args = append(args, "-output-http-debug")
	}

	if spec.ResponseBuffer > 0 {
		args = append(args, "-output-http-response-buffer")
		args = append(args, strconv.Itoa(spec.ResponseBuffer))
	}

	return &args, nil
}

func elasticsearchSiloToArgs(spec *kubereplayv1alpha1.ElasticsearchSilo) (*[]string, error) {
	var args []string
	if spec.Uri == "" {
		return nil, errors.New("uri is required for elasticsearch output mode")
	}
	args = append(args, "-output-http-elasticsearch")
	args = append(args, spec.Uri)

	return &args, nil
}

func kafkaSiloToArgs(spec *kubereplayv1alpha1.KafkaSilo) (*[]string, error) {
	var args []string
	if spec.Uri == "" {
		return nil, errors.New("uri is required for kafka output mode")
	}
	args = append(args, "-output-kafka-host")
	args = append(args, spec.Uri)

	if spec.Json == true {
		args = append(args, "-output-kafka-json-format")
	}

	if spec.Topic != "" {
		args = append(args, "-output-kafka-topic")
		args = append(args, spec.Topic)
	}

	return &args, nil
}

func mergeArgs(newArgs []string, args []string) []string {
	for _, a := range newArgs {
		args = append(args, a)
	}
	return args
}

func GenerateService(name string, namespace string, spec *kubereplayv1alpha1.RefinerySpec) *apiv1.Service {
	return &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("refinery-%s", name),
			Namespace: namespace,
			Labels:    AppLabels(name),
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Name:       "goreplay",
					Protocol:   apiv1.ProtocolTCP,
					Port:       28020,
					TargetPort: intstr.FromInt(28020),
				},
			},
			Selector: AppLabels(name),
			Type:     apiv1.ServiceTypeClusterIP,
		},
	}
}

func argsFromSpec(spec *kubereplayv1alpha1.RefinerySpec) *[]string {
	var args []string
	// Confugure input arguments
	args = append(args, "-input-tcp")
	args = append(args, ":28020")

	if spec.Storage.File != nil && spec.Storage.File.Enabled == true {
		fileArgs, err := fileSiloToArgs(spec.Storage.File)
		if err != nil {
			log.Print(err)
		} else {
			args = mergeArgs(*fileArgs, args)
		}
	}
	if spec.Storage.Tcp != nil && spec.Storage.Tcp.Enabled == true {
		tcpArgs, err := tcpSiloToArgs(spec.Storage.Tcp)
		if err != nil {
			log.Print(err)
		} else {
			args = mergeArgs(*tcpArgs, args)
		}
	}

	if spec.Storage.Stdout != nil && spec.Storage.Stdout.Enabled == true {
		stdoutArgs, err := stdoutSiloToArgs(spec.Storage.Stdout)
		if err != nil {
			log.Print(err)
		} else {
			args = mergeArgs(*stdoutArgs, args)
		}
	}

	if spec.Storage.Http != nil && spec.Storage.Http.Enabled == true {
		httpArgs, err := httpSiloToArgs(spec.Storage.Http)
		if err != nil {
			log.Print(err)
		} else {
			args = mergeArgs(*httpArgs, args)
		}
	}

	if spec.Storage.Elasticsearch != nil && spec.Storage.Elasticsearch.Enabled == true {
		elasticsearchArgs, err := elasticsearchSiloToArgs(spec.Storage.Elasticsearch)
		if err != nil {
			log.Print(err)
		} else {
			args = mergeArgs(*elasticsearchArgs, args)
		}
	}

	if spec.Storage.Kafka != nil && spec.Storage.Kafka.Enabled == true {
		kafkaArgs, err := kafkaSiloToArgs(spec.Storage.Kafka)
		if err != nil {
			log.Print(err)
		} else {
			args = mergeArgs(*kafkaArgs, args)
		}
	}

	if spec.Workers > 0 {
		args = append(args, "-output-http-workers")
		args = append(args, fmt.Sprint(spec.Workers))
	}

	if spec.Timeout != "" {
		args = append(args, "-output-http-timeout")
		args = append(args, spec.Timeout)
	}
	return &args
}

func GenerateDeployment(name string, r *kubereplayv1alpha1.Refinery) *appsv1.Deployment {
	if &r.Spec == nil || r.Spec.Storage == nil {
		return nil
	}
	var imagePullSecrets []corev1.LocalObjectReference
	var image string
	var imagePullPolicy apiv1.PullPolicy
	var resources corev1.ResourceRequirements

	args := argsFromSpec(&r.Spec)

	ownerReferences := []metav1.OwnerReference{
		{
			Name:       r.Name,
			UID:        r.UID,
			Kind:       "Refinery",
			APIVersion: kubereplayv1alpha1.SchemeGroupVersion.String(),
		},
	}
	if r.Spec.Goreplay != nil && r.Spec.Goreplay.ImagePullSecrets != nil {
		imagePullSecrets = r.Spec.Goreplay.ImagePullSecrets
	}
	if r.Spec.Goreplay != nil && r.Spec.Goreplay.Image != "" {
		image = r.Spec.Goreplay.Image
	} else {
		image = "buger/goreplay:latest"
	}
	if r.Spec.Goreplay != nil && r.Spec.Goreplay.ImagePullPolicy != "" {
		imagePullPolicy = r.Spec.Goreplay.ImagePullPolicy
	} else {
		imagePullPolicy = apiv1.PullAlways
	}
	if r.Spec.Resources != nil {
		resources = *r.Spec.Resources
	} else {
		resources = corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("10m"),
				corev1.ResourceMemory: resource.MustParse("64Mi"),
			},
			Limits: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("10m"),
				corev1.ResourceMemory: resource.MustParse("64Mi"),
			},
		}
	}

	deployment := &appsv1.Deployment{

		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: ownerReferences,
			Name:            fmt.Sprintf("refinery-%s", name),
			Labels:          AppLabels(name),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: AppLabels(name),
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: AppLabels(name),
				},
				Spec: apiv1.PodSpec{
					ImagePullSecrets: imagePullSecrets,
					Containers: []apiv1.Container{
						{
							Name:            "goreplay",
							Image:           image,
							ImagePullPolicy: imagePullPolicy,
							Args:            *args,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "tcp",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 28020,
								},
							},
							Resources: resources,
						},
					},
				},
			},
		},
	}

	return deployment

}

func int32Ptr(i int32) *int32 { return &i }
func int64Ptr(i int64) *int64 { return &i }
