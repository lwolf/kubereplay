package refinery

import (
	"fmt"
	"log"
	"strconv"

	kubereplayv1alpha1 "github.com/lwolf/kubereplay/pkg/apis/kubereplay/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func fileSiloToArgs(spec *kubereplayv1alpha1.FileSilo) *[]string {
	var args []string
	if spec.Filename == "" {
		log.Fatalf("Filename is required for file output mode")
	} else {
		args = append(args, "--output-file")
		args = append(args, spec.Filename)
	}
	if spec.Append == true {
		args = append(args, "--output-file-append")
	}

	if spec.FlushInterval != "" {
		args = append(args, "--output-file-flush-interval")
		args = append(args, spec.FlushInterval)
	}

	if spec.QueueSize > 0 {
		args = append(args, "--output-file-queue-limit")
		args = append(args, fmt.Sprint(spec.QueueSize))
	}

	if spec.FileLimit != "" {
		args = append(args, "--output-file-size-limit")
		args = append(args, spec.FileLimit)

	}

	return &args
}

func tcpSiloToArgs(spec *kubereplayv1alpha1.TcpSilo) *[]string {
	var args []string
	if spec.Uri == "" {
		log.Fatalf("Uri is required for tcp output mode")
	}
	args = append(args, "--output-tcp")
	args = append(args, spec.Uri)
	return &args
}

func stdoutSiloToArgs(spec *kubereplayv1alpha1.StdoutSilo) *[]string {
	var args []string
	args = append(args, "--output-stdout")
	return &args
}

func httpSiloToArgs(spec *kubereplayv1alpha1.HttpSilo) *[]string {
	var args []string
	if spec.Uri == "" {
		log.Fatalf("Uri is required for http output mode")
	}
	args = append(args, "--output-http")
	args = append(args, spec.Uri)

	if spec.Debug == true {
		args = append(args, "--output-http-debug")
	}

	if spec.ResponseBuffer > 0 {
		args = append(args, "--output-http-response-buffer")
		args = append(args, strconv.Itoa(spec.ResponseBuffer))
	}

	return &args
}

func elasticsearchSiloToArgs(spec *kubereplayv1alpha1.ElasticsearchSilo) *[]string {
	var args []string
	if spec.Uri == "" {
		log.Fatalf("Uri is required for elasticsearch output mode")
	}
	args = append(args, "--output-http-elasticsearch")
	args = append(args, spec.Uri)

	return &args
}

func kafkaSiloToArgs(spec *kubereplayv1alpha1.KafkaSilo) *[]string {
	var args []string
	if spec.Uri == "" {
		log.Fatalf("Uri is required for kafka output mode")
	}
	args = append(args, "--output-kafka-host")
	args = append(args, spec.Uri)

	if spec.Json == true {
		args = append(args, "--output-kafka-json-format")
	}

	if spec.Topic != "" {
		args = append(args, "--output-kafka-topic")
		args = append(args, spec.Topic)
	}

	return &args
}

func mergeArgs(newArgs []string, args []string) []string {
	for _, a := range newArgs {
		args = append(args, a)
	}
	return args
}

func GenerateSidecar(refinerySvc string, port uint32) *apiv1.Container {
	return &apiv1.Container{
		Name:  "goreplay",
		Image: "buger/goreplay:latest",
		Args: []string{
			"--input-raw",
			fmt.Sprintf(":%d", port),
			"--output-tcp",
			fmt.Sprintf("%s:28020", refinerySvc),
		},
		//Resources: apiv1.ResourceRequirements{
		//	Limits: apiv1.ResourceList{},
		//	Requests: apiv1.ResourceList{},
		//},
	}
}

func ConfigmapName(harvesterName string) string {
	return fmt.Sprintf("%s-sidecar", harvesterName)
}

func GenerateConfigmap(name string, spec *kubereplayv1alpha1.HarvesterSpec) *apiv1.ConfigMap {
	container := GenerateSidecar(
		fmt.Sprintf("refinery-%s.default", spec.Refinery),
		spec.AppPort,
	)
	return &apiv1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
			Labels: map[string]string{
				"kubereplay-app": name,
			},
		},
		// TODO: pretty print container data
		Data: map[string]string{
			"key": fmt.Sprintf("|\n %s", container.String()),
		},
	}
}

func GenerateService(name string, spec *kubereplayv1alpha1.RefinerySpec) *apiv1.Service {
	return &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("refinery-%s", name),
			Namespace: "default",
			Labels: map[string]string{
				"kubereplay-app": name,
			},
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Name:       "gor",
					Protocol:   apiv1.ProtocolTCP,
					Port:       28020,
					TargetPort: intstr.FromInt(28020),
				},
			},
			Selector: map[string]string{
				"kubereplay-app": name,
			},
			Type: apiv1.ServiceTypeClusterIP,
		},
	}
}

func GenerateDeployment(name string, r *kubereplayv1alpha1.Refinery) *appsv1.Deployment {

	var args []string
	// Confugure input arguments
	args = append(args, "--input-tcp")
	args = append(args, ":28020")

	spec := r.Spec

	if spec.Storage.File.Enabled == true {
		fileArgs := fileSiloToArgs(spec.Storage.File)
		args = mergeArgs(*fileArgs, args)
	}
	if spec.Storage.Tcp.Enabled == true {
		tcpArgs := tcpSiloToArgs(spec.Storage.Tcp)
		args = mergeArgs(*tcpArgs, args)
	}

	if spec.Storage.Stdout.Enabled == true {
		stdoutArgs := stdoutSiloToArgs(spec.Storage.Stdout)
		args = mergeArgs(*stdoutArgs, args)
	}

	if spec.Storage.Http.Enabled == true {
		httpArgs := httpSiloToArgs(spec.Storage.Http)
		args = mergeArgs(*httpArgs, args)
	}

	if spec.Storage.Elasticsearch.Enabled == true {
		elasticsearchArgs := elasticsearchSiloToArgs(spec.Storage.Elasticsearch)
		args = mergeArgs(*elasticsearchArgs, args)
	}

	if spec.Storage.Kafka.Enabled == true {
		kafkaArgs := kafkaSiloToArgs(spec.Storage.Kafka)
		args = mergeArgs(*kafkaArgs, args)
	}

	if spec.Workers > 0 {
		args = append(args, "--output-http-workers")
		args = append(args, fmt.Sprint(spec.Workers))
	}

	if spec.Timeout != "" {
		args = append(args, "-output-http-timeout")
		args = append(args, spec.Timeout)
	}

	ownerReferences := []metav1.OwnerReference{
		{
			Name:       r.Name,
			UID:        r.UID,
			Kind:       "Refinery",
			APIVersion: kubereplayv1alpha1.SchemeGroupVersion.String(),
		},
	}

	deployment := &appsv1.Deployment{

		ObjectMeta: metav1.ObjectMeta{
			OwnerReferences: ownerReferences,
			Name:            fmt.Sprintf("refinery-%s", name),
			Labels: map[string]string{
				"kubereplay-app": name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"kubereplay-app": name,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"kubereplay-app": name,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "goreplay",
							Image: "buger/goreplay:latest",
							Args:  args,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "tcp",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 28020,
								},
							},
						},
					},
				},
			},
		},
	}

	return deployment

}

func int32Ptr(i int32) *int32 { return &i }
