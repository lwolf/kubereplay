package utils

import (
	"fmt"
	"log"
	"strconv"

	"github.com/lwolf/kube-replay/pkg/apis/replay"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func fileOutputToArgs(spec *replay.FileOutput) *[]string {
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

func tcpOutputToArgs(spec *replay.TcpOutput) *[]string {
	var args []string
	if spec.Uri == "" {
		log.Fatalf("Uri is required for tcp output mode")
	}
	args = append(args, "--output-tcp")
	args = append(args, spec.Uri)
	return &args
}

func stdoutOutputToArgs(spec *replay.StdoutOutput) *[]string {
	var args []string
	args = append(args, "--output-stdout")
	return &args
}

func httpOutputToArgs(spec *replay.HttpOutput) *[]string {
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

func elasticsearchOutputToArgs(spec *replay.ElasticsearchOutput) *[]string {
	var args []string
	if spec.Uri == "" {
		log.Fatalf("Uri is required for elasticsearch output mode")
	}
	args = append(args, "--output-http-elasticsearch")
	args = append(args, spec.Uri)

	return &args
}

func kafkaOutputToArgs(spec *replay.KafkaOutput) *[]string {
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

func CreateDeployment(siloName string, spec *replay.RefinerySpec) *appsv1.Deployment {

	var args []string

	if spec.Output.File.Enabled == true {
		fileArgs := fileOutputToArgs(&spec.Output.File)
		args = mergeArgs(*fileArgs, args)
	}
	if spec.Output.Tcp.Enabled == true {
		tcpArgs := tcpOutputToArgs(&spec.Output.Tcp)
		args = mergeArgs(*tcpArgs, args)
	}

	if spec.Output.Stdout.Enabled == true {
		stdoutArgs := stdoutOutputToArgs(&spec.Output.Stdout)
		args = mergeArgs(*stdoutArgs, args)
	}

	if spec.Output.Http.Enabled == true {
		httpArgs := httpOutputToArgs(&spec.Output.Http)
		args = mergeArgs(*httpArgs, args)
	}

	if spec.Output.Elasticsearch.Enabled == true {
		elasticsearchArgs := elasticsearchOutputToArgs(&spec.Output.Elasticsearch)
		args = mergeArgs(*elasticsearchArgs, args)
	}

	if spec.Output.Kafka.Enabled == true {
		kafkaArgs := kafkaOutputToArgs(&spec.Output.Kafka)
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
