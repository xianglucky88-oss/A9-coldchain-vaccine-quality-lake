package store

import (
	"path/filepath"

	"example.com/coldchain-lake/internal/model"
)

type Repository struct {
	Root string
}

func NewRepository(root string) *Repository { return &Repository{Root: root} }

func (r *Repository) AppendSamples(samples []model.Sample) error {
	return AppendJSONL(filepath.Join(r.Root, "raw", "samples.jsonl"), samples)
}

func (r *Repository) AppendEvents(events []model.DoorEvent) error {
	return AppendJSONL(filepath.Join(r.Root, "raw", "door_events.jsonl"), events)
}

func (r *Repository) ReadSamples() ([]model.Sample, error) {
	return ReadJSONL[model.Sample](filepath.Join(r.Root, "raw", "samples.jsonl"))
}

func (r *Repository) ReadEvents() ([]model.DoorEvent, error) {
	return ReadJSONL[model.DoorEvent](filepath.Join(r.Root, "raw", "door_events.jsonl"))
}
