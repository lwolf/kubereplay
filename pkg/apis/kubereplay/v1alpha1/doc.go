// Api versions allow the api contract for a resource to be changed while keeping
// backward compatibility by support multiple concurrent versions
// of the same resource

// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=package,register
// +k8s:conversion-gen=github.com/lwolf/kubereplay/pkg/apis/kubereplay
// +k8s:defaulter-gen=TypeMeta
// +groupName=kubereplay.lwolf.org
package v1alpha1 // import "github.com/lwolf/kubereplay/pkg/apis/kubereplay/v1alpha1"
