package service

import (
	"example.com/coldchain-lake/internal/engine"
	"example.com/coldchain-lake/internal/model"
	"example.com/coldchain-lake/internal/store"
)

type Service struct {
	repository *store.Repository
	artifacts  *store.ArtifactRepository
}

func New(root string) *Service {
	return &Service{
		repository: store.NewRepository(root),
		artifacts:  store.NewArtifactRepository(root),
	}
}

// Raw records are retained so the compute stage can report data-quality
// rejects instead of silently losing bad source rows during ingestion.
func (s *Service) IngestSamples(samples []model.Sample) error {
	return s.repository.AppendSamples(samples)
}

func (s *Service) IngestEvents(events []model.DoorEvent) error {
	return s.repository.AppendEvents(events)
}

func (s *Service) Compute(config engine.Config) (engine.Result, error) {
	samples, err := s.repository.ReadSamples()
	if err != nil {
		return engine.Result{}, err
	}
	events, err := s.repository.ReadEvents()
	if err != nil {
		return engine.Result{}, err
	}
	result := engine.Compute(samples, events, config)
	if err := s.artifacts.Save(result); err != nil {
		return engine.Result{}, err
	}
	return result, nil
}
