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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kurtisdevbrv1beta1 "github.com/angelokurtis/configmap-generator/api/v1beta1"
	"github.com/angelokurtis/configmap-generator/pkg/handler"
)

// ConfigMapGeneratorReconciler reconciles a ConfigMapGenerator object
type ConfigMapGeneratorReconciler struct {
	reconciler.Result
	client  client.Client
	handler ConfigMapGeneratorHandler
}

func NewConfigMapGeneratorReconciler(client client.Client, handler ConfigMapGeneratorHandler) *ConfigMapGeneratorReconciler {
	return &ConfigMapGeneratorReconciler{client: client, handler: handler}
}

// ConfigMapGeneratorHandler aggregates the handlers of ConfigMapGenerator
type ConfigMapGeneratorHandler struct {
	SourceFromGitRepository *handler.SourceFromGitRepository
	ConfigmapCreation       *handler.ConfigMapCreation
}

//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kurtis.dev.br,resources=configmapgenerators,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kurtis.dev.br,resources=configmapgenerators/finalizers,verbs=update
//+kubebuilder:rbac:groups=kurtis.dev.br,resources=configmapgenerators/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=source.toolkit.fluxcd.io,resources=gitrepositories,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ConfigMapGeneratorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	resource := &kurtisdevbrv1beta1.ConfigMapGenerator{}
	err := r.client.Get(ctx, req.NamespacedName, resource)
	if errors.IsNotFound(err) {
		return r.Finish(ctx) // Ignoring since object must be deleted
	}
	if err != nil {
		return r.RequeueOnErr(ctx, err) // Failed to get ConfigMapGenerator
	}

	return reconciler.Chain(
		r.handler.SourceFromGitRepository,
		r.handler.ConfigmapCreation,
	).Reconcile(ctx, resource)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConfigMapGeneratorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kurtisdevbrv1beta1.ConfigMapGenerator{}).
		Owns(&corev1.ConfigMap{}).
		// Watches(&source.Kind{Type: &v1beta1.GitRepository{}}, &EnqueueRequestForGitRepositoryRef{}). // TODO: watches for GitRepository changes
		Complete(r)
}
