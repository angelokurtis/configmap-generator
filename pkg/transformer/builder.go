package transformer

import "github.com/angelokurtis/configmap-generator/api/v1beta1"

type Builder struct{ reference ObjectReference }

func NewBuilder(reference ObjectReference) *Builder {
	return &Builder{reference: reference}
}

type ConfigMapGenerator struct {
	resource  *v1beta1.ConfigMapGenerator
	reference ObjectReference
}

func (b *Builder) ForConfigMapGenerator(resource *v1beta1.ConfigMapGenerator) *ConfigMapGenerator {
	return &ConfigMapGenerator{
		resource:  resource,
		reference: b.reference,
	}
}
