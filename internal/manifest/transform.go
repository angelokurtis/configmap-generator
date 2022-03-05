package manifest

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

type Transformer func(u *unstructured.Unstructured) error

func (f *fullList) Transform(funcs ...Transformer) (List, error) {
	resources := make([]*unstructured.Unstructured, 0, f.Size())
	for _, v := range f.Resources() {
		resource := v.DeepCopy()
		for _, transform := range funcs {
			if err := transform(resource); err != nil {
				return &fullList{}, err
			}
		}
		resources = append(resources, resource)
	}
	return &fullList{resources: resources, client: f.client, mapper: f.mapper}, nil
}
