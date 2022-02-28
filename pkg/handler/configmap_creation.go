package handler

import (
	"context"

	mf "github.com/manifestival/manifestival"

	"github.com/angelokurtis/reconciler"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/angelokurtis/configmap-generator/api/v1beta1"
	"github.com/angelokurtis/configmap-generator/internal/kustomize"
)

type ConfigMapCreation struct {
	reconciler.Funcs
	kustomize Kustomize
}

func NewConfigMapCreation(kustomize Kustomize) *ConfigMapCreation {
	return &ConfigMapCreation{kustomize: kustomize}
}

type Kustomize interface {
	GenerateConfigMap(ctx context.Context, source string, generator *kustomize.ConfigMapGenerator) (mf.Manifest, error)
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

func (c *ConfigMapCreation) reconcile(ctx context.Context, source string, gen *v1beta1.ConfigMapGenerator) (ctrl.Result, error) {
	generator := &kustomize.ConfigMapGenerator{Name: gen.GetName(), Files: gen.Spec.Files, Literals: gen.Spec.Literals}
	manifest, err := c.kustomize.GenerateConfigMap(ctx, source, generator)
	if err != nil {
		return c.RequeueOnErr(ctx, err)
	}
	_ = manifest
	return c.Next(ctx, gen)
}
