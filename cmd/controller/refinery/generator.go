package refinery

import (
	"fmt"
	"log"
	"strconv"

	"github.com/lwolf/kube-replay/pkg/apis/replay"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func fileSiloToArgs(spec *replay.FileSilo) *[]string {
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
		args = append(args, strconv.Itoa(spec.QueueSize))
	}

	if spec.FileLimit != "" {
		args = append(args, "--output-file-size-limit")
		args = append(args, spec.FileLimit)

	}

	return &args
}

func tcpSiloToArgs(spec *replay.TcpSilo) *[]string {
	var args []string
	if spec.Uri == "" {
		log.Fatalf("Uri is required for tcp output mode")
	}
	args = append(args, "--output-tcp")
	args = append(args, spec.Uri)
	return &args
}

func stdoutSiloToArgs(spec *replay.StdoutSilo) *[]string {
	var args []string
	args = append(args, "--output-stdout")
	return &args
}

func httpSiloToArgs(spec *replay.HttpSilo) *[]string {
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

func elasticsearchSiloToArgs(spec *replay.ElasticsearchSilo) *[]string {
	var args []string
	if spec.Uri == "" {
		log.Fatalf("Uri is required for elasticsearch output mode")
	}
	args = append(args, "--output-http-elasticsearch")
	args = append(args, spec.Uri)

	return &args
}

func kafkaSiloToArgs(spec *replay.KafkaSilo) *[]string {
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

func GenerateDeployment(siloName string, spec *replay.RefinerySpec) *appsv1.Deployment {

	var args []string
	// Confugure input arguments
	args = append(args, "--input-tcp")
	args = append(args, ":28020")

	if spec.Storage.File.Enabled == true {
		fileArgs := fileSiloToArgs(&spec.Storage.File)
		args = mergeArgs(*fileArgs, args)
	}
	if spec.Storage.Tcp.Enabled == true {
		tcpArgs := tcpSiloToArgs(&spec.Storage.Tcp)
		args = mergeArgs(*tcpArgs, args)
	}

	if spec.Storage.Stdout.Enabled == true {
		stdoutArgs := stdoutSiloToArgs(&spec.Storage.Stdout)
		args = mergeArgs(*stdoutArgs, args)
	}

	if spec.Storage.Http.Enabled == true {
		httpArgs := httpSiloToArgs(&spec.Storage.Http)
		args = mergeArgs(*httpArgs, args)
	}

	if spec.Storage.Elasticsearch.Enabled == true {
		elasticsearchArgs := elasticsearchSiloToArgs(&spec.Storage.Elasticsearch)
		args = mergeArgs(*elasticsearchArgs, args)
	}

	if spec.Storage.Kafka.Enabled == true {
		kafkaArgs := kafkaSiloToArgs(&spec.Storage.Kafka)
		args = mergeArgs(*kafkaArgs, args)
	}

	if spec.Workers > 0 {
		args = append(args, "--output-http-workers")
		args = append(args, strconv.Itoa(spec.Workers))
	}

	if spec.Timeout != "" {
		args = append(args, "-output-http-timeout")
		args = append(args, spec.Timeout)
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("refinery-%s", siloName),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "kubereplay",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "kubereplay",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "replay-silo",
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
