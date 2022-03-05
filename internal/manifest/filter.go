package manifest

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

type Filter func(u *unstructured.Unstructured) bool

func (f *fullList) Filter(funcs ...Filter) List {
	resources := make([]*unstructured.Unstructured, 0, f.Size())
	for _, v := range f.Resources() {
		resource := v.DeepCopy()
		for _, filter := range funcs {
			if filter(resource) {
				resources = append(resources, resource)
			}
		}
	}
	return &fullList{resources: resources, client: f.client, mapper: f.mapper}
}
