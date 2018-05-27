/*
Copyright 2017 Sergey Nuzhdin.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/package v1alpha1

import (
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: "kubereplay.lwolf.org", Version: "v1alpha1"}

// Kind takes an unqualified kind and returns back a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

// Adds the list of known types to Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Harvester{},
		&HarvesterList{},
		&Refinery{},
		&RefineryList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type HarvesterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Harvester `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type RefineryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Refinery `json:"items"`
}

// CRD Generation
func getFloat(f float64) *float64 {
	return &f
}

func getInt(i int64) *int64 {
	return &i
}

var (
	// Define CRDs for resources
	HarvesterCRD = v1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "harvesters.kubereplay.lwolf.org",
		},
		Spec: v1beta1.CustomResourceDefinitionSpec{
			Group:   "kubereplay.lwolf.org",
			Version: "v1alpha1",
			Names: v1beta1.CustomResourceDefinitionNames{
				Kind:   "Harvester",
				Plural: "harvesters",
			},
			Scope: "Namespaced",
			Validation: &v1beta1.CustomResourceValidation{
				OpenAPIV3Schema: &v1beta1.JSONSchemaProps{
					Type: "object",
					Properties: map[string]v1beta1.JSONSchemaProps{
						"apiVersion": v1beta1.JSONSchemaProps{
							Type: "string",
						},
						"kind": v1beta1.JSONSchemaProps{
							Type: "string",
						},
						"metadata": v1beta1.JSONSchemaProps{
							Type: "object",
						},
						"spec": v1beta1.JSONSchemaProps{
							Type: "object",
							Properties: map[string]v1beta1.JSONSchemaProps{
								"refinery": v1beta1.JSONSchemaProps{
									Type: "string",
								},
								"segment": v1beta1.JSONSchemaProps{
									Type:   "integer",
									Format: "int32",
								},
								"selector": v1beta1.JSONSchemaProps{
									Type: "object",
								},
							},
						},
						"status": v1beta1.JSONSchemaProps{
							Type: "object",
							Properties: map[string]v1beta1.JSONSchemaProps{
								"segment": v1beta1.JSONSchemaProps{
									Type:   "integer",
									Format: "int32",
								},
							},
						},
					},
				},
			},
		},
	}
	// Define CRDs for resources
	RefineryCRD = v1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "refineries.kubereplay.lwolf.org",
		},
		Spec: v1beta1.CustomResourceDefinitionSpec{
			Group:   "kubereplay.lwolf.org",
			Version: "v1alpha1",
			Names: v1beta1.CustomResourceDefinitionNames{
				Kind:   "Refinery",
				Plural: "refineries",
			},
			Scope: "Namespaced",
			Validation: &v1beta1.CustomResourceValidation{
				OpenAPIV3Schema: &v1beta1.JSONSchemaProps{
					Type: "object",
					Properties: map[string]v1beta1.JSONSchemaProps{
						"apiVersion": v1beta1.JSONSchemaProps{
							Type: "string",
						},
						"kind": v1beta1.JSONSchemaProps{
							Type: "string",
						},
						"metadata": v1beta1.JSONSchemaProps{
							Type: "object",
						},
						"spec": v1beta1.JSONSchemaProps{
							Type: "object",
							Properties: map[string]v1beta1.JSONSchemaProps{
								"output": v1beta1.JSONSchemaProps{
									Type: "object",
									Properties: map[string]v1beta1.JSONSchemaProps{
										"elasticsearch": v1beta1.JSONSchemaProps{
											Type: "object",
											Properties: map[string]v1beta1.JSONSchemaProps{
												"enabled": v1beta1.JSONSchemaProps{
													Type: "boolean",
												},
												"uri": v1beta1.JSONSchemaProps{
													Type: "string",
												},
											},
										},
										"file": v1beta1.JSONSchemaProps{
											Type: "object",
											Properties: map[string]v1beta1.JSONSchemaProps{
												"append": v1beta1.JSONSchemaProps{
													Type: "boolean",
												},
												"enabled": v1beta1.JSONSchemaProps{
													Type: "boolean",
												},
												"filelimit": v1beta1.JSONSchemaProps{
													Type: "string",
												},
												"filename": v1beta1.JSONSchemaProps{
													Type: "string",
												},
												"flushinterval": v1beta1.JSONSchemaProps{
													Type: "string",
												},
												"queuesize": v1beta1.JSONSchemaProps{
													Type:   "integer",
													Format: "int32",
												},
											},
										},
										"http": v1beta1.JSONSchemaProps{
											Type: "object",
											Properties: map[string]v1beta1.JSONSchemaProps{
												"debug": v1beta1.JSONSchemaProps{
													Type: "boolean",
												},
												"enabled": v1beta1.JSONSchemaProps{
													Type: "boolean",
												},
												"uri": v1beta1.JSONSchemaProps{
													Type: "string",
												},
											},
										},
										"kafka": v1beta1.JSONSchemaProps{
											Type: "object",
											Properties: map[string]v1beta1.JSONSchemaProps{
												"enabled": v1beta1.JSONSchemaProps{
													Type: "boolean",
												},
												"json": v1beta1.JSONSchemaProps{
													Type: "boolean",
												},
												"topic": v1beta1.JSONSchemaProps{
													Type: "string",
												},
												"uri": v1beta1.JSONSchemaProps{
													Type: "string",
												},
											},
										},
										"stdout": v1beta1.JSONSchemaProps{
											Type: "object",
											Properties: map[string]v1beta1.JSONSchemaProps{
												"enabled": v1beta1.JSONSchemaProps{
													Type: "boolean",
												},
											},
										},
										"tcp": v1beta1.JSONSchemaProps{
											Type: "object",
											Properties: map[string]v1beta1.JSONSchemaProps{
												"enabled": v1beta1.JSONSchemaProps{
													Type: "boolean",
												},
												"uri": v1beta1.JSONSchemaProps{
													Type: "string",
												},
											},
										},
									},
								},
								"timeout": v1beta1.JSONSchemaProps{
									Type: "string",
								},
							},
						},
						"status": v1beta1.JSONSchemaProps{
							Type: "object",
							Properties: map[string]v1beta1.JSONSchemaProps{
								"deployed": v1beta1.JSONSchemaProps{
									Type: "boolean",
								},
							},
						},
					},
				},
			},
		},
	}
)
