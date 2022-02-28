package handler

import (
	"context"
	"net/url"
	"os"

	"github.com/angelokurtis/reconciler"
	sourcev1beta1 "github.com/fluxcd/source-controller/api/v1beta1"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/angelokurtis/configmap-generator/api/v1beta1"
)

type SourceFromGitRepository struct {
	reconciler.Funcs
	reader     GitRepositoryReader
	downloader ArtifactDownloader
}

func NewSourceFromGitRepository(reader GitRepositoryReader, downloader ArtifactDownloader) *SourceFromGitRepository {
	return &SourceFromGitRepository{reader: reader, downloader: downloader}
}

type GitRepositoryReader interface {
	FetchGitRepository(ctx context.Context, key client.ObjectKey) (*sourcev1beta1.GitRepository, error)
}

type ArtifactDownloader interface {
	DownloadArtifact(ctx context.Context, artifact *sourcev1beta1.Artifact, dest string) (bool, error)
}

func (s *SourceFromGitRepository) Reconcile(ctx context.Context, obj client.Object) (ctrl.Result, error) {
	gen, ok := obj.(*v1beta1.ConfigMapGenerator)
	if !ok || gen.Spec.SourceRef.Kind != "GitRepository" {
		return s.Next(ctx, obj)
	}
	return s.reconcile(ctx, gen)
}

func (s *SourceFromGitRepository) reconcile(ctx context.Context, gen *v1beta1.ConfigMapGenerator) (ctrl.Result, error) {
	log := logr.FromContextOrDiscard(ctx)

	repo, err := s.reader.FetchGitRepository(ctx, gen.GetSourceRefNamespacedName())
	if err != nil {
		return s.RequeueOnErr(ctx, err)
	}

	if repo == nil {
		log.Info("GitRepository was not found")
		return s.Next(ctx, gen)
	}

	artifact := repo.GetArtifact()
	if artifact == nil {
		log.Info("GitRepository is not ready")
		return s.Next(ctx, gen)
	}

	u, err := url.Parse(artifact.URL)
	if err != nil {
		return s.RequeueOnErr(ctx, err)
	}
	dest := os.TempDir() + u.Path
	ok, err := s.downloader.DownloadArtifact(ctx, artifact, dest)
	if err != nil {
		return s.RequeueOnErr(ctx, err)
	}
	if ok {
		log.Info("Source downloaded", "path", dest, "checksum", artifact.Checksum)
	} else {
		log.Info("Source is already available locally", "path", dest, "checksum", artifact.Checksum)
	}

	return s.Next(contextWithSource(ctx, dest), gen)
}
