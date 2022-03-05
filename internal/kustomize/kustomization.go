package kustomize

import "gopkg.in/yaml.v3"

type Kustomization struct {
	GeneratorOptions   *GeneratorOptions     `yaml:"generatorOptions,omitempty"`
	ConfigMapGenerator []*ConfigMapGenerator `yaml:"configMapGenerator,omitempty"`
}

type GeneratorOptions struct {
	DisableNameSuffixHash bool              `yaml:"disableNameSuffixHash,omitempty"`
	Labels                map[string]string `json:"labels,omitempty"`
	Annotations           map[string]string `json:"annotations,omitempty"`
}

type ConfigMapGenerator struct {
	Name     string   `yaml:"name,omitempty"`
	Files    []string `yaml:"files,omitempty"`
	Literals []string `yaml:"literals,omitempty"`
}

func (k *Kustomization) Marshal() ([]byte, error) {
	return yaml.Marshal(k)
}
