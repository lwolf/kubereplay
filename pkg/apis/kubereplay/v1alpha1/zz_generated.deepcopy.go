// +build !ignore_autogenerated

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
*/
// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ElasticsearchSilo) DeepCopyInto(out *ElasticsearchSilo) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ElasticsearchSilo.
func (in *ElasticsearchSilo) DeepCopy() *ElasticsearchSilo {
	if in == nil {
		return nil
	}
	out := new(ElasticsearchSilo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FileSilo) DeepCopyInto(out *FileSilo) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FileSilo.
func (in *FileSilo) DeepCopy() *FileSilo {
	if in == nil {
		return nil
	}
	out := new(FileSilo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Harvester) DeepCopyInto(out *Harvester) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Harvester.
func (in *Harvester) DeepCopy() *Harvester {
	if in == nil {
		return nil
	}
	out := new(Harvester)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Harvester) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HarvesterList) DeepCopyInto(out *HarvesterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Harvester, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HarvesterList.
func (in *HarvesterList) DeepCopy() *HarvesterList {
	if in == nil {
		return nil
	}
	out := new(HarvesterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HarvesterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HarvesterSpec) DeepCopyInto(out *HarvesterSpec) {
	*out = *in
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HarvesterSpec.
func (in *HarvesterSpec) DeepCopy() *HarvesterSpec {
	if in == nil {
		return nil
	}
	out := new(HarvesterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HarvesterStatus) DeepCopyInto(out *HarvesterStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HarvesterStatus.
func (in *HarvesterStatus) DeepCopy() *HarvesterStatus {
	if in == nil {
		return nil
	}
	out := new(HarvesterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HttpSilo) DeepCopyInto(out *HttpSilo) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HttpSilo.
func (in *HttpSilo) DeepCopy() *HttpSilo {
	if in == nil {
		return nil
	}
	out := new(HttpSilo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KafkaSilo) DeepCopyInto(out *KafkaSilo) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KafkaSilo.
func (in *KafkaSilo) DeepCopy() *KafkaSilo {
	if in == nil {
		return nil
	}
	out := new(KafkaSilo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Refinery) DeepCopyInto(out *Refinery) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Refinery.
func (in *Refinery) DeepCopy() *Refinery {
	if in == nil {
		return nil
	}
	out := new(Refinery)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Refinery) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RefineryList) DeepCopyInto(out *RefineryList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Refinery, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RefineryList.
func (in *RefineryList) DeepCopy() *RefineryList {
	if in == nil {
		return nil
	}
	out := new(RefineryList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RefineryList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RefinerySpec) DeepCopyInto(out *RefinerySpec) {
	*out = *in
	if in.Storage != nil {
		in, out := &in.Storage, &out.Storage
		if *in == nil {
			*out = nil
		} else {
			*out = new(RefineryStorage)
			(*in).DeepCopyInto(*out)
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RefinerySpec.
func (in *RefinerySpec) DeepCopy() *RefinerySpec {
	if in == nil {
		return nil
	}
	out := new(RefinerySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RefineryStatus) DeepCopyInto(out *RefineryStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RefineryStatus.
func (in *RefineryStatus) DeepCopy() *RefineryStatus {
	if in == nil {
		return nil
	}
	out := new(RefineryStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RefineryStorage) DeepCopyInto(out *RefineryStorage) {
	*out = *in
	if in.File != nil {
		in, out := &in.File, &out.File
		if *in == nil {
			*out = nil
		} else {
			*out = new(FileSilo)
			**out = **in
		}
	}
	if in.Tcp != nil {
		in, out := &in.Tcp, &out.Tcp
		if *in == nil {
			*out = nil
		} else {
			*out = new(TcpSilo)
			**out = **in
		}
	}
	if in.Stdout != nil {
		in, out := &in.Stdout, &out.Stdout
		if *in == nil {
			*out = nil
		} else {
			*out = new(StdoutSilo)
			**out = **in
		}
	}
	if in.Http != nil {
		in, out := &in.Http, &out.Http
		if *in == nil {
			*out = nil
		} else {
			*out = new(HttpSilo)
			**out = **in
		}
	}
	if in.Elasticsearch != nil {
		in, out := &in.Elasticsearch, &out.Elasticsearch
		if *in == nil {
			*out = nil
		} else {
			*out = new(ElasticsearchSilo)
			**out = **in
		}
	}
	if in.Kafka != nil {
		in, out := &in.Kafka, &out.Kafka
		if *in == nil {
			*out = nil
		} else {
			*out = new(KafkaSilo)
			**out = **in
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RefineryStorage.
func (in *RefineryStorage) DeepCopy() *RefineryStorage {
	if in == nil {
		return nil
	}
	out := new(RefineryStorage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StdoutSilo) DeepCopyInto(out *StdoutSilo) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StdoutSilo.
func (in *StdoutSilo) DeepCopy() *StdoutSilo {
	if in == nil {
		return nil
	}
	out := new(StdoutSilo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TcpSilo) DeepCopyInto(out *TcpSilo) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TcpSilo.
func (in *TcpSilo) DeepCopy() *TcpSilo {
	if in == nil {
		return nil
	}
	out := new(TcpSilo)
	in.DeepCopyInto(out)
	return out
}
