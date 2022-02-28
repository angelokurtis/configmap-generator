package v1beta1

import "k8s.io/apimachinery/pkg/types"

func (in *ConfigMapGenerator) GetSourceRefNamespacedName() types.NamespacedName {
	name := in.Spec.SourceRef.Name
	namespace := in.Spec.SourceRef.Namespace
	if namespace == "" {
		namespace = in.Namespace
	}
	return types.NamespacedName{Namespace: namespace, Name: name}
}
