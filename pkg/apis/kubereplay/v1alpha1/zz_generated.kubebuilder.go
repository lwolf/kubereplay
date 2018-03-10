package v1alpha1

import (
	"github.com/kubernetes-sigs/kubebuilder/pkg/builders"
	"github.com/lwolf/kubereplay/pkg/apis/kubereplay"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

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
								"appPort": v1beta1.JSONSchemaProps{
									Type:   "integer",
									Format: "int32",
								},
								"refinery": v1beta1.JSONSchemaProps{
									Type: "string",
								},
								"segment": v1beta1.JSONSchemaProps{
									Type:   "integer",
									Format: "int32",
								},
								"selector": v1beta1.JSONSchemaProps{
									Type: "object",
									AdditionalProperties: &v1beta1.JSONSchemaPropsOrBool{
										Allows: true,
										//Schema: &,
									},
								},
							},
						},
						"status": v1beta1.JSONSchemaProps{
							Type:       "object",
							Properties: map[string]v1beta1.JSONSchemaProps{},
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

	// Collect CRDs in this group version
	CRDs = []v1beta1.CustomResourceDefinition{
		HarvesterCRD,
		RefineryCRD,
	}

	// Register APIs
	kubereplayHarvesterStorage = builders.NewApiResource( // Resource status endpoint
		kubereplay.InternalHarvester,
		HarvesterSchemeFns{},
		func() runtime.Object { return &Harvester{} },     // Register versioned resource
		func() runtime.Object { return &HarvesterList{} }, // Register versioned resource list
		&HarvesterStrategy{builders.StorageStrategySingleton},
	)
	kubereplayRefineryStorage = builders.NewApiResource( // Resource status endpoint
		kubereplay.InternalRefinery,
		RefinerySchemeFns{},
		func() runtime.Object { return &Refinery{} },     // Register versioned resource
		func() runtime.Object { return &RefineryList{} }, // Register versioned resource list
		&RefineryStrategy{builders.StorageStrategySingleton},
	)
	ApiVersion = builders.NewApiVersion("kubereplay.lwolf.org", "v1alpha1").WithResources(
		kubereplayHarvesterStorage,
		builders.NewApiResource( // Resource status endpoint
			kubereplay.InternalHarvesterStatus,
			HarvesterSchemeFns{},
			func() runtime.Object { return &Harvester{} },     // Register versioned resource
			func() runtime.Object { return &HarvesterList{} }, // Register versioned resource list
			&HarvesterStatusStrategy{builders.StatusStorageStrategySingleton},
		), kubereplayRefineryStorage,
		builders.NewApiResource( // Resource status endpoint
			kubereplay.InternalRefineryStatus,
			RefinerySchemeFns{},
			func() runtime.Object { return &Refinery{} },     // Register versioned resource
			func() runtime.Object { return &RefineryList{} }, // Register versioned resource list
			&RefineryStatusStrategy{builders.StatusStorageStrategySingleton},
		))

	// Required by code generated by go2idl
	AddToScheme        = ApiVersion.SchemaBuilder.AddToScheme
	SchemeBuilder      = ApiVersion.SchemaBuilder
	localSchemeBuilder = &SchemeBuilder
	SchemeGroupVersion = ApiVersion.GroupVersion
)

func getFloat(f float64) *float64 {
	return &f
}

// Required by code generated by go2idl
// Kind takes an unqualified kind and returns a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Required by code generated by go2idl
// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

//
// Harvester Functions and Structs
//
// +k8s:deepcopy-gen=false
type HarvesterSchemeFns struct {
	builders.DefaultSchemeFns
}

// +k8s:deepcopy-gen=false
type HarvesterStrategy struct {
	builders.DefaultStorageStrategy
}

// +k8s:deepcopy-gen=false
type HarvesterStatusStrategy struct {
	builders.DefaultStatusStorageStrategy
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type HarvesterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Harvester `json:"items"`
}

//
// Refinery Functions and Structs
//
// +k8s:deepcopy-gen=false
type RefinerySchemeFns struct {
	builders.DefaultSchemeFns
}

// +k8s:deepcopy-gen=false
type RefineryStrategy struct {
	builders.DefaultStorageStrategy
}

// +k8s:deepcopy-gen=false
type RefineryStatusStrategy struct {
	builders.DefaultStatusStorageStrategy
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type RefineryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Refinery `json:"items"`
}
