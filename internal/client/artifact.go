package client

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	sourcev1beta1 "github.com/fluxcd/source-controller/api/v1beta1"
)

type Artifact struct{}

func (a *Artifact) DownloadArtifact(ctx context.Context, artifact *sourcev1beta1.Artifact, dest string) (bool, error) {
	if _, err := os.Stat(dest); !errors.Is(err, os.ErrNotExist) && checksum(dest, artifact.Checksum) {
		return false, nil
	}

	res, err := http.Get(artifact.URL)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	index := strings.LastIndex(dest, "/")
	if err = os.MkdirAll(dest[:index], os.ModePerm); err != nil {
		return false, err
	}

	out, err := os.Create(dest)
	if err != nil {
		return false, err
	}
	defer out.Close()

	if _, err = io.Copy(out, res.Body); err != nil {
		return false, err
	}

	return true, nil
}

func checksum(filepath, checksum string) bool {
	f, err := os.Open(filepath)
	if err != nil {
		return false
	}
	defer f.Close()

	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return false
	}

	return fmt.Sprintf("%x", h.Sum(nil)) == checksum
}
