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

package v1beta1

import (
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigMapGeneratorSpec defines the desired state of ConfigMapGenerator
type ConfigMapGeneratorSpec struct {
	Files       []string          `json:"files,omitempty"`
	Literals    []string          `json:"literals,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`

	// Reference of the source where the kustomization file is.
	// +required
	SourceRef kustomizev1.CrossNamespaceSourceReference `json:"sourceRef"`
}

// ConfigMapGeneratorStatus defines the observed state of ConfigMapGenerator
type ConfigMapGeneratorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ConfigMapGenerator is the Schema for the configmapgenerators API
type ConfigMapGenerator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConfigMapGeneratorSpec   `json:"spec,omitempty"`
	Status ConfigMapGeneratorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ConfigMapGeneratorList contains a list of ConfigMapGenerator
type ConfigMapGeneratorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ConfigMapGenerator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ConfigMapGenerator{}, &ConfigMapGeneratorList{})
}
