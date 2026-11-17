package store

import (
	"path/filepath"

	"example.com/coldchain-lake/internal/engine"
)

type ArtifactRepository struct {
	Root string
}

func NewArtifactRepository(root string) *ArtifactRepository { return &ArtifactRepository{Root: root} }

func (r *ArtifactRepository) Save(result engine.Result) error {
	if err := WriteJSONL(filepath.Join(r.Root, "artifacts", "excursions.jsonl"), result.Excursions); err != nil {
		return err
	}
	return WriteJSONL(filepath.Join(r.Root, "artifacts", "quality.jsonl"), result.Quality)
}
