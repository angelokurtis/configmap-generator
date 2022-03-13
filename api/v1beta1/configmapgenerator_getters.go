package v1beta1

import (
	"github.com/fluxcd/kustomize-controller/api/v1beta2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (in *ConfigMapGenerator) GetSourceRefKey() client.ObjectKey {
	ref := in.GetSourceRef()
	return client.ObjectKey{Namespace: ref.Namespace, Name: ref.Name}
}

func (in *ConfigMapGenerator) GetSourceRef() v1beta2.CrossNamespaceSourceReference {
	ref := in.Spec.SourceRef
	if ref.Namespace == "" {
		ref.Namespace = in.Namespace
	}
	return ref
}
