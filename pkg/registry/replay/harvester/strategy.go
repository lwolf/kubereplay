package harvester

import (
	"fmt"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"

	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"github.com/lwolf/kube-replay/pkg/apis/replay"
)


type harvesterStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

func (harvesterStrategy) NamespaceScoped() bool {
	return true
}

func (harvesterStrategy) PrepareForCreate(ctx genericapirequest.Context, obj runtime.Object) {
}

func (harvesterStrategy) PrepareForUpdate(ctx genericapirequest.Context, obj, old runtime.Object) {
}

func (harvesterStrategy) Validate(ctx genericapirequest.Context, obj runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

func (harvesterStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (harvesterStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (harvesterStrategy) Canonicalize(obj runtime.Object) {
}

func (harvesterStrategy) ValidateUpdate(ctx genericapirequest.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}


func NewStrategy(typer runtime.ObjectTyper) harvesterStrategy {
	return harvesterStrategy{typer, names.SimpleNameGenerator}
}

func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, bool, error) {
	apiserver, ok := obj.(*replay.Harvester)
	if !ok {
		return nil, nil, false, fmt.Errorf("given object is not a Harvester.")
	}
	return labels.Set(apiserver.ObjectMeta.Labels), HarvesterToSelectableFields(apiserver), apiserver.Initializers != nil, nil
}

// MatchFlunder is the filter used by the generic etcd backend to watch events
// from etcd to clients of the apiserver only interested in specific labels/fields.
func MatchHarvester(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

// FlunderToSelectableFields returns a field set that represents the object.
func HarvesterToSelectableFields(obj *replay.Harvester) fields.Set {
	return generic.ObjectMetaFieldsSet(&obj.ObjectMeta, true)
}