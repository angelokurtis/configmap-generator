package handler

import (
	"context"

	"github.com/angelokurtis/reconciler"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ConfigMapCreation struct{ reconciler.Funcs }

func (c *ConfigMapCreation) Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	return c.Next(ctx, obj)
}
