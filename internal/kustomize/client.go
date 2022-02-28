package kustomize

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/hashicorp/go-getter"

	mf "github.com/manifestival/manifestival"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

type Client struct {
	fs         filesys.FileSystem
	kustomizer *krusty.Kustomizer
	manifest   Manifest
}

func NewClient(fs filesys.FileSystem, kustomizer *krusty.Kustomizer, manifest Manifest) *Client {
	return &Client{fs: fs, kustomizer: kustomizer, manifest: manifest}
}

type Manifest interface {
	FromBytes(ctx context.Context, manifests []byte) (mf.Manifest, error)
}

func (c *Client) GenerateConfigMap(ctx context.Context, source string, generator *ConfigMapGenerator) (mf.Manifest, error) {
	destination := source[0 : len(source)-len(".tar.gz")]
	defer os.RemoveAll(destination)
	if err := getter.GetAny(destination, source); err != nil {
		return mf.Manifest{}, fmt.Errorf("error extracting Source artifact %s: %w", source, err)
	}
	content, err := c.generateConfigMap(ctx, destination, generator)
	if err != nil {
		return mf.Manifest{}, err
	}

	return c.manifest.FromBytes(ctx, content)
}

func (c *Client) generateConfigMap(ctx context.Context, dir string, generator *ConfigMapGenerator) ([]byte, error) {
	content, err := generator.Marshal()
	if err != nil {
		return nil, err
	}

	file := path.Join(dir, "kustomization.yaml")
	err = os.WriteFile(file, content, 0600)
	if err != nil {
		return nil, err
	}
	defer os.Remove(file)

	m, err := c.kustomizer.Run(c.fs, dir)
	if err != nil {
		return nil, err
	}

	return m.AsYaml()
}
