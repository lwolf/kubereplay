package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Silo describes a database.
type Silo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SiloSpec   `json:"spec"`
	Status SiloStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SiloList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Silo `json:"items"`
}

// SiloSpec is the spec for a Silo resource
type SiloSpec struct {
	Workers int    `json:"workers"`
	Timeout string `json:"timeout"`
	Output struct {
		File          SiloFileOutput          `json:"file,omitempty"`
		Tcp           SiloTcpOutput           `json:"tcp,omitempty"`
		Stdout        SiloStdoutOutput        `json:"stdout,omitempty"`
		Http          SiloHttpOutput          `json:"http,omitempty"`
		Elasticsearch SiloElasticsearchOutput `json:"elasticsearch,omitempty"`
		Kafka         SiloKafkaOutput         `json:"kafka,omitempty"`
	}
}

type SiloFileOutput struct {
	Enabled       bool   `json:"enabled"`
	Filename      string `json:"filename"`
	Append        bool   `json:"append"`
	FlushInterval string `json:"flushinterval"`
	QueueSize     int    `json:"queuesize"`
	FileLimit     string `json:"filelimit"`
}

type SiloTcpOutput struct {
	Enabled bool   `json:"enabled"`
	Uri     string `json:"uri"`
}

type SiloStdoutOutput struct {
	Enabled bool `json:"enabled"`
}

type SiloHttpOutput struct {
	Enabled        bool   `json:"enabled"`
	Uri            string `json:"uri"`
	Debug          bool   `json:"debug"`
	ResponseBuffer int    `json:"response_buffer"`
}

type SiloElasticsearchOutput struct {
	Enabled bool   `json:"enabled"`
	Uri     string `json:"uri"`
}

type SiloKafkaOutput struct {
	Enabled bool   `json:"enabled"`
	Uri     string `json:"uri"`
	Json    bool   `json:"json"`
	Topic   string `json:"topic"`
}

type SiloStatus struct {
	Deployed bool `json:"deployed"`
}

// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Harvester describes resources
type Harvester struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec HarvesterSpec `json:"items"`
}

// HarvesterSpec is the spec for a Harvester resource
type HarvesterSpec struct {
	Selector    string  `json:"selector"`
	Silo        string  `json:"silo"`
	SegmentSize float32 `json:"segment"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type HarvesterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Harvester `json:"items"`
}
