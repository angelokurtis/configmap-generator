package manifest

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func (f *fullList) Delete(ctx context.Context) error {
	for _, v := range f.Resources() {
		if err := f.delete(ctx, v); err != nil {
			return err
		}
	}
	return nil
}

func (f *fullList) delete(ctx context.Context, obj *unstructured.Unstructured) error {
	log := logr.FromContextOrDiscard(ctx)

	gvk := obj.GroupVersionKind()
	mapper, err := f.mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	resource := f.client.Resource(mapper.Resource).Namespace(obj.GetNamespace())
	if mapper.Scope.Name() == meta.RESTScopeNameRoot {
		resource = f.client.Resource(mapper.Resource)
	}

	err = resource.Delete(ctx, obj.GetName(), metav1.DeleteOptions{})
	if errors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}

	kind := fmt.Sprintf("%s.%s", strings.ToLower(gvk.Kind), gvk.Group)
	if len(gvk.Group) == 0 {
		kind = strings.ToLower(gvk.Kind)
	}
	log.Info(fmt.Sprintf("%s %q deleted", kind, obj.GetName()))
	return nil
}
