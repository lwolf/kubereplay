package replay

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Refinery struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec   RefinerySpec
	Status RefineryStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type RefineryList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []Refinery
}

type RefineryOutput struct {
	File          FileOutput
	Tcp           TcpOutput
	Stdout        StdoutOutput
	Http          HttpOutput
	Elasticsearch ElasticsearchOutput
	Kafka         KafkaOutput
}

type RefinerySpec struct {
	Workers int
	Timeout string
	Output  RefineryOutput
}

type RefineryOutputTest struct {
	Enabled bool
}

type RefineryStatus struct {
	Deployed bool
}

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Silo struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	Spec   SiloSpec
	Status SiloStatus
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SiloList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []Silo
}

type SiloSpec struct {
	Workers int
	Timeout string
	Output struct {
		File          FileOutput
		Tcp           TcpOutput
		Stdout        StdoutOutput
		Http          HttpOutput
		Elasticsearch ElasticsearchOutput
		Kafka         KafkaOutput
	}
}

type FileOutput struct {
	Enabled       bool
	Filename      string
	Append        bool
	FlushInterval string
	QueueSize     int
	FileLimit     string
}
type TcpOutput struct {
	Enabled bool
	Uri     string
}
type StdoutOutput struct {
	Enabled bool
}
type HttpOutput struct {
	Enabled        bool
	Uri            string
	Debug          bool
	ResponseBuffer int
}
type ElasticsearchOutput struct {
	Enabled bool
	Uri     string
}
type KafkaOutput struct {
	Enabled bool
	Uri     string
	Json    bool
	Topic   string
}

type SiloStatus struct {
	Deployed bool
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Harvester struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec HarvesterSpec
}

type HarvesterSpec struct {
	Selector    string
	Silo        string
	SegmentSize float32
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type HarvesterList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []Harvester
}
