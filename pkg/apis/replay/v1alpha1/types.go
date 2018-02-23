package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Refinery struct {
	metav1.TypeMeta       `json:",inline"`
	metav1.ObjectMeta     `json:"metadata"`
	Spec   RefinerySpec   `json:"spec"`
	Status RefineryStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type RefineryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Refinery `json:"items"`
}

type RefinerySpec struct {
	Workers int             `json:"workers"`
	Timeout string          `json:"timeout"`
	Storage RefineryStorage `json:"output"`
}

type RefineryStorage struct {
	File          FileSilo          `json:"file,omitempty"`
	Tcp           TcpSilo           `json:"tcp,omitempty"`
	Stdout        StdoutSilo        `json:"stdout,omitempty"`
	Http          HttpSilo          `json:"http,omitempty"`
	Elasticsearch ElasticsearchSilo `json:"elasticsearch,omitempty"`
	Kafka         KafkaSilo         `json:"kafka,omitempty"`
}

type RefineryStatus struct {
	Deployed bool `json:"deployed"`
}

type FileSilo struct {
	Enabled       bool   `json:"enabled"`
	Filename      string `json:"filename"`
	Append        bool   `json:"append"`
	FlushInterval string `json:"flushinterval"`
	QueueSize     int    `json:"queuesize"`
	FileLimit     string `json:"filelimit"`
}

type TcpSilo struct {
	Enabled bool   `json:"enabled"`
	Uri     string `json:"uri"`
}

type StdoutSilo struct {
	Enabled bool `json:"enabled"`
}

type HttpSilo struct {
	Enabled        bool   `json:"enabled"`
	Uri            string `json:"uri"`
	Debug          bool   `json:"debug"`
	ResponseBuffer int    `json:"response_buffer"`
}

type ElasticsearchSilo struct {
	Enabled bool   `json:"enabled"`
	Uri     string `json:"uri"`
}

type KafkaSilo struct {
	Enabled bool   `json:"enabled"`
	Uri     string `json:"uri"`
	Json    bool   `json:"json"`
	Topic   string `json:"topic"`
}

type SiloStatus struct {
	Deployed bool `json:"deployed"`
}

// +genclient
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
	Refinery    string  `json:"refinery"`
	SegmentSize float32 `json:"segment"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type HarvesterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Harvester `json:"items"`
}
