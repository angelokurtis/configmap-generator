package manifest

import (
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
)

type fullList struct {
	resources []*unstructured.Unstructured
	client    dynamic.Interface
	mapper    meta.RESTMapper
}

func (f *fullList) Resources() []*unstructured.Unstructured {
	return f.resources
}

func (f *fullList) Size() int {
	return len(f.resources)
}

func (f *fullList) Append(mfs ...List) List {
	resources := make([]*unstructured.Unstructured, 0, f.Size())
	for _, v := range f.Resources() {
		resource := v.DeepCopy()
		resources = append(resources, resource)
	}
	for _, mf := range mfs {
		for _, v := range mf.Resources() {
			resource := v.DeepCopy()
			resources = append(resources, resource)
		}
	}
	return &fullList{resources: resources, client: f.client, mapper: f.mapper}
}
