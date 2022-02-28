/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/angelokurtis/reconciler"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kurtisdevbrv1beta1 "github.com/angelokurtis/configmap-generator/api/v1beta1"
	"github.com/angelokurtis/configmap-generator/pkg/generator"
)

// ConfigMapGeneratorReconciler reconciles a ConfigMapGenerator object
type ConfigMapGeneratorReconciler struct {
	reconciler.Result
	client.Client
}

//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kurtis.dev.br,resources=configmapgenerators,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kurtis.dev.br,resources=configmapgenerators/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kurtis.dev.br,resources=configmapgenerators/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ConfigMapGeneratorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	gen := &kurtisdevbrv1beta1.ConfigMapGenerator{}
	err := r.Get(ctx, req.NamespacedName, gen)
	if errors.IsNotFound(err) {
		return r.Finish(ctx) // Ignoring since object must be deleted
	}
	if err != nil {
		return r.RequeueOnErr(ctx, err) // Failed to get ConfigMapGenerator
	}

	return reconciler.Chain(
		&generator.ConfigmapCreation{},
	).Reconcile(ctx, gen)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConfigMapGeneratorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kurtisdevbrv1beta1.ConfigMapGenerator{}).
		Complete(r)
}
