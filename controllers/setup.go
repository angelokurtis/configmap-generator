package controllers

import (
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/kustomize/api/krusty"

	"sigs.k8s.io/kustomize/kyaml/filesys"

	"github.com/angelokurtis/configmap-generator/internal/client"
	"github.com/angelokurtis/configmap-generator/internal/kustomize"
	"github.com/angelokurtis/configmap-generator/internal/manifest"
	"github.com/angelokurtis/configmap-generator/internal/object"
	"github.com/angelokurtis/configmap-generator/pkg/handler"
	"github.com/angelokurtis/configmap-generator/pkg/transformer"
)

func SetupWithManager(mgr ctrl.Manager) error {
	dynamicRESTMapper, err := apiutil.NewDynamicRESTMapper(mgr.GetConfig())
	if err != nil {
		return err
	}
	configMapGeneratorHandler := ConfigMapGeneratorHandler{
		SourceFromGitRepository: handler.NewSourceFromGitRepository(
			client.NewGitRepository(mgr.GetClient()),
			&client.Artifact{},
		),
		ConfigmapCreation: handler.NewConfigMapCreation(
			kustomize.NewClient(
				filesys.MakeFsOnDisk(),
				krusty.MakeKustomizer(krusty.MakeDefaultOptions()),
				manifest.NewReader(
					dynamic.NewForConfigOrDie(mgr.GetConfig()),
					dynamicRESTMapper,
				),
			),
			transformer.NewBuilder(object.NewReference(mgr.GetScheme())),
		),
	}
	return NewConfigMapGeneratorReconciler(mgr.GetClient(), configMapGeneratorHandler).SetupWithManager(mgr)
}
