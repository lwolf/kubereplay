package replay

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
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
		File          SiloFileOutput
		Tcp           SiloTcpOutput
		Stdout        SiloStdoutOutput
		Http          SiloHttpOutput
		Elasticsearch SiloElasticsearchOutput
		Kafka         SiloKafkaOutput
	}
}

type SiloFileOutput struct {
	Enabled       bool
	Filename      string
	Append        bool
	FlushInterval string
	QueueSize     int
	FileLimit     string
}
type SiloTcpOutput struct {
	Enabled bool
	Uri     string
}
type SiloStdoutOutput struct {
	Enabled bool
}
type SiloHttpOutput struct {
	Enabled        bool
	Uri            string
	Debug          bool
	ResponseBuffer int
}
type SiloElasticsearchOutput struct {
	Enabled bool
	Uri     string
}
type SiloKafkaOutput struct {
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
