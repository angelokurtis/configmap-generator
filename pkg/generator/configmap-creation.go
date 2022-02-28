package generator

import (
	"context"

	"github.com/angelokurtis/reconciler"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ConfigmapCreation struct{ reconciler.Funcs }

func (c *ConfigmapCreation) Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	return c.Next(ctx, obj)
}
