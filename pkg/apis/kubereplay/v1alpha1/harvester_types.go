package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Created by "kubebuilder create resource" for you to implement the Harvester resource schema definition
// as a go struct

// HarvesterSpec defines the desired state of Harvester
type HarvesterSpec struct {
	Selector    map[string]string `json:"selector,omitempty"`
	AppPort     uint32            `json:"appPort,omitempty"`
	Refinery    string            `json:"refinery,omitempty"`
	SegmentSize uint32            `json:"segment,omitempty"`
}

// HarvesterStatus defines the observed state of Harvester
type HarvesterStatus struct {
	ControlledObjects []ControlledObject `json:"controlled,omitempty"`
}

type ControlledObject struct {
	Kind          string `json:"kind"`
	BlueName      string `json:"blue_name"`
	BlueReplicas  int32  `json:"blue_replicas"`
	GreenName     string `json:"green_name"`
	GreenReplicas int32  `json:"green_replicas"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Harvester
// +k8s:openapi-gen=true
// +resource:path=harvesters
type Harvester struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HarvesterSpec   `json:"spec,omitempty"`
	Status HarvesterStatus `json:"status,omitempty"`
}
