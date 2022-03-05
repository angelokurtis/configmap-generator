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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

const fieldManager = "configmap-generator"

func (f *fullList) Apply(ctx context.Context) error {
	for _, v := range f.Resources() {
		if err := f.apply(ctx, v); err != nil {
			return err
		}
	}
	return nil
}

func (f *fullList) apply(ctx context.Context, obj *unstructured.Unstructured) error {
	log := logr.FromContextOrDiscard(ctx)

	current, err := f.find(ctx, obj)
	if err != nil {
		return err
	}

	gvk := obj.GroupVersionKind()
	kind := fmt.Sprintf("%s.%s", strings.ToLower(gvk.Kind), gvk.Group)
	if len(gvk.Group) == 0 {
		kind = strings.ToLower(gvk.Kind)
	}

	if current == nil { // create
		if _, err = f.patch(ctx, obj); err != nil {
			return err
		}
		log.Info(fmt.Sprintf("%s %q created", kind, obj.GetName()))
		return nil
	}

	// update
	updated, err := f.patch(ctx, obj)
	if err != nil {
		return err
	}
	if current.GetResourceVersion() == updated.GetResourceVersion() {
		log.Info(fmt.Sprintf("%s %q unchanged", kind, obj.GetName()))
		return nil
	}

	log.Info(fmt.Sprintf("%s %q configured", kind, obj.GetName()))
	return nil
}

func (f *fullList) find(ctx context.Context, obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	gvk := obj.GroupVersionKind()
	mapper, err := f.mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, err
	}

	resource := f.client.Resource(mapper.Resource).Namespace(obj.GetNamespace())
	if mapper.Scope.Name() == meta.RESTScopeNameRoot {
		resource = f.client.Resource(mapper.Resource)
	}

	result, err := resource.Get(ctx, obj.GetName(), metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (f *fullList) patch(ctx context.Context, obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	gvk := obj.GroupVersionKind()
	mapper, err := f.mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, err
	}

	data, err := runtime.Encode(unstructured.UnstructuredJSONScheme, obj)
	if err != nil {
		return nil, err
	}

	resource := f.client.Resource(mapper.Resource).Namespace(obj.GetNamespace())
	if mapper.Scope.Name() == meta.RESTScopeNameRoot {
		resource = f.client.Resource(mapper.Resource)
	}

	force := true
	options := metav1.PatchOptions{Force: &force, FieldManager: fieldManager}
	return resource.Patch(ctx, obj.GetName(), types.ApplyPatchType, data, options)
}
