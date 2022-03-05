package manifest

import (
	"bytes"
	"errors"
	"io"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
)

type Reader struct {
	client dynamic.Interface
	mapper meta.RESTMapper
}

func NewReader(client dynamic.Interface, mapper meta.RESTMapper) *Reader {
	return &Reader{client: client, mapper: mapper}
}

func (s *Reader) FromBytes(data []byte) (List, error) {
	reader := bytes.NewReader(data)
	decoder := yaml.NewYAMLToJSONDecoder(reader)
	var resources []*unstructured.Unstructured
	var err error
	for {
		out := &unstructured.Unstructured{}
		err = decoder.Decode(out)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil || len(out.Object) == 0 {
			continue
		}
		resources = append(resources, out)
	}
	if !errors.Is(err, io.EOF) {
		return &fullList{}, err
	}
	return &fullList{resources: resources, client: s.client, mapper: s.mapper}, nil
}
