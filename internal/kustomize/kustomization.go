package kustomize

import "gopkg.in/yaml.v3"

type Kustomization struct {
	GeneratorOptions   *GeneratorOptions     `yaml:"generatorOptions,omitempty"`
	ConfigMapGenerator []*ConfigMapGenerator `yaml:"configMapGenerator,omitempty"`
}

type GeneratorOptions struct {
	DisableNameSuffixHash bool `yaml:"disableNameSuffixHash,omitempty"`
}

type ConfigMapGenerator struct {
	Name     string   `yaml:"name,omitempty"`
	Files    []string `yaml:"files,omitempty"`
	Literals []string `yaml:"literals,omitempty"`
}

func (c *ConfigMapGenerator) Marshal() ([]byte, error) {
	generators := []*ConfigMapGenerator{c}
	k := &Kustomization{
		GeneratorOptions:   &GeneratorOptions{DisableNameSuffixHash: true},
		ConfigMapGenerator: generators,
	}
	return yaml.Marshal(k)
}
