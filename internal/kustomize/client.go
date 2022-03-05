package kustomize

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/hashicorp/go-getter"

	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/kyaml/filesys"

	"github.com/angelokurtis/configmap-generator/internal/manifest"
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
	FromBytes(data []byte) (manifest.List, error)
}

func (c *Client) Build(ctx context.Context, source string, kustomization *Kustomization) (manifest.List, error) {
	dest := source[0 : len(source)-len(".tar.gz")]
	defer os.RemoveAll(dest)
	if err := getter.GetAny(dest, source); err != nil {
		return nil, fmt.Errorf("error extracting Source artifact %s: %w", source, err)
	}
	content, err := c.build(ctx, dest, kustomization)
	if err != nil {
		return nil, err
	}

	return c.manifest.FromBytes(content)
}

func (c *Client) build(ctx context.Context, dir string, kustomization *Kustomization) ([]byte, error) {
	content, err := kustomization.Marshal()
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
