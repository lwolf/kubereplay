package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!
// Created by "kubebuilder create resource" for you to implement the Refinery resource schema definition
// as a go struct

// RefinerySpec defines the desired state of Refinery
type RefinerySpec struct {
	Workers int32            `json:"workers,omitempчty"`
	Timeout string           `json:"timeout,omitempty"`
	Storage *RefineryStorage `json:"output,omitempty"`
}

// RefineryStatus defines various storages available for Refinery
type RefineryStorage struct {
	File          *FileSilo          `json:"file,omitempty"`
	Tcp           *TcpSilo           `json:"tcp,omitempty"`
	Stdout        *StdoutSilo        `json:"stdout,omitempty"`
	Http          *HttpSilo          `json:"http,omitempty"`
	Elasticsearch *ElasticsearchSilo `json:"elasticsearch,omitempty"`
	Kafka         *KafkaSilo         `json:"kafka,omitempty"`
}

type FileSilo struct {
	Enabled       bool   `json:"enabled,omitempty"`
	Filename      string `json:"filename,omitempty"`
	Append        bool   `json:"append,omitempty"`
	FlushInterval string `json:"flushinterval,omitempty"`
	QueueSize     int32  `json:"queuesize,omitempty"`
	FileLimit     string `json:"filelimit,omitempty"`
}

type TcpSilo struct {
	Enabled bool   `json:"enabled,omitempty"`
	Uri     string `json:"uri,omitempty"`
}

type StdoutSilo struct {
	Enabled bool `json:"enabled,omitempty"`
}

type HttpSilo struct {
	Enabled        bool   `json:"enabled,omitempty"`
	Uri            string `json:"uri,omitempty"`
	Debug          bool   `json:"debug,omitempty"`
	ResponseBuffer int    `json:"response_buffer,omitempty"`
}

type ElasticsearchSilo struct {
	Enabled bool   `json:"enabled,omitempty"`
	Uri     string `json:"uri,omitempty"`
}

type KafkaSilo struct {
	Enabled bool   `json:"enabled,omitempty"`
	Uri     string `json:"uri,omitempty"`
	Json    bool   `json:"json,omitempty"`
	Topic   string `json:"topic,omitempty"`
}

// RefineryStatus defines the observed state of Refinery
type RefineryStatus struct {
	// INSERT YOUR CODE HERE - define observed state schema
	Deployed bool `json:"deployed,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Refinery
// +k8s:openapi-gen=true
// +resource:path=refineries
type Refinery struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RefinerySpec   `json:"spec,omitempty"`
	Status RefineryStatus `json:"status,omitempty"`
}
