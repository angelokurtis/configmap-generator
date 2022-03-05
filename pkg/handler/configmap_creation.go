package handler

import (
	"context"

	"github.com/angelokurtis/reconciler"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/angelokurtis/configmap-generator/api/v1beta1"
	"github.com/angelokurtis/configmap-generator/internal/kustomize"
	"github.com/angelokurtis/configmap-generator/internal/manifest"
	"github.com/angelokurtis/configmap-generator/pkg/transformer"
)

type ConfigMapCreation struct {
	reconciler.Funcs
	kustomize    Kustomize
	transformers *transformer.Builder
}

func NewConfigMapCreation(kustomize Kustomize, transformers *transformer.Builder) *ConfigMapCreation {
	return &ConfigMapCreation{kustomize: kustomize, transformers: transformers}
}

type Kustomize interface {
	Build(ctx context.Context, source string, kustomization *kustomize.Kustomization) (manifest.List, error)
}

func (c *ConfigMapCreation) Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	gen, ok := obj.(*v1beta1.ConfigMapGenerator)
	if !ok {
		return c.Next(ctx, obj)
	}
	source := sourceFromContext(ctx)
	if source == "" {
		return c.Next(ctx, gen)
	}
	return c.reconcile(ctx, source, gen)
}

func (c *ConfigMapCreation) reconcile(ctx context.Context, source string, resource *v1beta1.ConfigMapGenerator) (ctrl.Result, error) {
	mf, err := c.kustomize.Build(ctx, source, kustomizationFor(resource))
	if err != nil {
		return c.RequeueOnErr(ctx, err)
	}

	transform := c.transformers.ForConfigMapGenerator(resource)
	mf, err = mf.Transform(transform.Metadata)
	if err != nil {
		return c.RequeueOnErr(ctx, err)
	}

	if err = mf.Apply(ctx); err != nil {
		return c.RequeueOnErr(ctx, err)
	}
	return c.Next(ctx, resource)
}

func kustomizationFor(resource *v1beta1.ConfigMapGenerator) *kustomize.Kustomization {
	return &kustomize.Kustomization{
		GeneratorOptions: &kustomize.GeneratorOptions{
			DisableNameSuffixHash: true,
			Labels:                resource.Spec.Labels,
			Annotations:           resource.Spec.Annotations,
		},
		ConfigMapGenerator: []*kustomize.ConfigMapGenerator{{
			Name:     resource.GetName(),
			Files:    resource.Spec.Files,
			Literals: resource.Spec.Literals,
		}},
	}
}
