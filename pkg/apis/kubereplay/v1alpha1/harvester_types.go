package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!
// Created by "kubebuilder create resource" for you to implement the Harvester resource schema definition
// as a go struct

// HarvesterSpec defines the desired state of Harvester
type HarvesterSpec struct {
	// INSERT YOUR CODE HERE - define desired state schema
	Selector    map[string]string `json:"selector,omitempty"`
	AppPort     uint32            `json:"appPort,omitempty"`
	Refinery    string            `json:"refinery,omitempty"`
	SegmentSize uint32            `json:"segment,omitempty"`
}

// HarvesterStatus defines the observed state of Harvester
type HarvesterStatus struct {
	// INSERT YOUR CODE HERE - define observed state schema
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
