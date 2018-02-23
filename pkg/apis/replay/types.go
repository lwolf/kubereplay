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

type RefineryStorage struct {
	File          FileSilo
	Tcp           TcpSilo
	Stdout        StdoutSilo
	Http          HttpSilo
	Elasticsearch ElasticsearchSilo
	Kafka         KafkaSilo
}

type RefinerySpec struct {
	Workers int
	Timeout string
	Storage RefineryStorage
}

type RefineryStatus struct {
	Deployed bool
}

type FileSilo struct {
	Enabled       bool
	Filename      string
	Append        bool
	FlushInterval string
	QueueSize     int
	FileLimit     string
}
type TcpSilo struct {
	Enabled bool
	Uri     string
}
type StdoutSilo struct {
	Enabled bool
}
type HttpSilo struct {
	Enabled        bool
	Uri            string
	Debug          bool
	ResponseBuffer int
}
type ElasticsearchSilo struct {
	Enabled bool
	Uri     string
}
type KafkaSilo struct {
	Enabled bool
	Uri     string
	Json    bool
	Topic   string
}

type SiloStatus struct {
	Deployed bool
}

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Harvester struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec HarvesterSpec
}

type HarvesterSpec struct {
	Selector    string
	Refinery    string
	SegmentSize float32
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type HarvesterList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []Harvester
}
