package service

import (
	"context"

	"github.com/spanish-english-discord/api/internal/model"
	"github.com/spanish-english-discord/api/internal/repository"
)

type PodcastService struct {
	repo *repository.PodcastRepository
}

func NewPodcastService(repo *repository.PodcastRepository) *PodcastService {
	return &PodcastService{repo: repo}
}

func (s *PodcastService) GetAll(ctx context.Context, filters *model.PodcastFilters) ([]model.Podcast, error) {
	podcasts, err := s.repo.GetAll(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Return empty slice instead of nil for JSON serialization
	if podcasts == nil {
		podcasts = []model.Podcast{}
	}

	return podcasts, nil
}

func (s *PodcastService) GetByID(ctx context.Context, id string) (*model.Podcast, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PodcastService) Create(ctx context.Context, input *model.CreatePodcastInput) (*model.Podcast, error) {
	return s.repo.Create(ctx, input)
}

func (s *PodcastService) Update(ctx context.Context, id string, input *model.UpdatePodcastInput) (*model.Podcast, error) {
	// Verify podcast exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.repo.Update(ctx, id, input)
}

func (s *PodcastService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
