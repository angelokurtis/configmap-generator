package transformer

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func (c *ConfigMapGenerator) Metadata(u *unstructured.Unstructured) error {
	u.SetName(c.resource.GetName())
	u.SetNamespace(c.resource.GetNamespace())
	if err := c.reference.SetController(c.resource, u); err != nil {
		return err
	}
	labels := u.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}
	labels["app.kubernetes.io/managed-by"] = "configmap-generator"
	labels["app.kubernetes.io/name"] = c.resource.GetName()
	u.SetLabels(labels)
	return nil
}
