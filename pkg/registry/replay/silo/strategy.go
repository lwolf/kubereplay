package silo

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

type siloStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

func (siloStrategy) NamespaceScoped() bool {
	return true
}

func (siloStrategy) PrepareForCreate(ctx genericapirequest.Context, obj runtime.Object) {
}

func (siloStrategy) PrepareForUpdate(ctx genericapirequest.Context, obj, old runtime.Object) {
}

func (siloStrategy) Validate(ctx genericapirequest.Context, obj runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

func (siloStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (siloStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (siloStrategy) Canonicalize(obj runtime.Object) {
}

func (siloStrategy) ValidateUpdate(ctx genericapirequest.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

func NewStrategy(typer runtime.ObjectTyper) siloStrategy {
	return siloStrategy{typer, names.SimpleNameGenerator}
}

func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, bool, error) {
	apiserver, ok := obj.(*replay.Silo)
	if !ok {
		return nil, nil, false, fmt.Errorf("given object is not a Silo.")
	}
	return labels.Set(apiserver.ObjectMeta.Labels), SiloToSelectableFields(apiserver), apiserver.Initializers != nil, nil
}

// MatchFlunder is the filter used by the generic etcd backend to watch events
// from etcd to clients of the apiserver only interested in specific labels/fields.
func MatchSilo(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

// FlunderToSelectableFields returns a field set that represents the object.
func SiloToSelectableFields(obj *replay.Silo) fields.Set {
	return generic.ObjectMetaFieldsSet(&obj.ObjectMeta, true)
}