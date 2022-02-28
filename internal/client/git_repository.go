package client

import (
	"context"
	"fmt"

	sourcev1beta1 "github.com/fluxcd/source-controller/api/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type GitRepository struct {
	client client.Client
}

func NewGitRepository(client client.Client) *GitRepository {
	return &GitRepository{client: client}
}

func (g *GitRepository) FetchGitRepository(ctx context.Context, key client.ObjectKey) (*sourcev1beta1.GitRepository, error) {
	m := new(sourcev1beta1.GitRepository)
	err := g.client.Get(ctx, key, m)
	if errors.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to lookup GitRepository: %w", err)
	}
	return m, nil

}
