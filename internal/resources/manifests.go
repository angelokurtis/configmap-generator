package resources

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	mf "github.com/manifestival/manifestival"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type ManifestProducer struct{ client mf.Client }

func NewManifestProducer(client mf.Client) *ManifestProducer {
	return &ManifestProducer{client: client}
}

func (s *ManifestProducer) FromBytes(ctx context.Context, manifests []byte) (mf.Manifest, error) {
	log := logr.FromContextOrDiscard(ctx).V(1)
	reader := bytes.NewReader(manifests)
	m, err := mf.ManifestFrom(mf.Reader(reader), mf.UseClient(s.client), mf.UseLogger(log))
	if err != nil {
		return mf.Manifest{}, fmt.Errorf("failed to read manifests from bytes: %w", err)
	}
	return m, nil
}

func (s *ManifestProducer) FromString(ctx context.Context, manifests string) (mf.Manifest, error) {
	log := logr.FromContextOrDiscard(ctx).V(1)
	reader := strings.NewReader(manifests)
	m, err := mf.ManifestFrom(mf.Reader(reader), mf.UseClient(s.client), mf.UseLogger(log))
	if err != nil {
		return mf.Manifest{}, fmt.Errorf("failed to read manifests from string: %w", err)
	}
	return m, nil
}

func (s *ManifestProducer) FromUnstructured(ctx context.Context, manifests []unstructured.Unstructured) (mf.Manifest, error) {
	log := logr.FromContextOrDiscard(ctx).V(1)
	m, err := mf.ManifestFrom(mf.Slice(manifests), mf.UseClient(s.client), mf.UseLogger(log))
	if err != nil {
		return mf.Manifest{}, fmt.Errorf("failed to read manifests from string: %w", err)
	}
	return m, nil
}
