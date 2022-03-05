package manifest

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type List interface {
	Delete(ctx context.Context) error
	Apply(ctx context.Context) error
	Filter(funcs ...Filter) List
	Transform(funcs ...Transformer) (List, error)
	Resources() []*unstructured.Unstructured
	Size() int
	Append(mfs ...List) List
}
