package service

import (
	"context"
	"errors"

	"github.com/spanish-english-discord/api/internal/model"
	"github.com/spanish-english-discord/api/internal/repository"
)

var ErrPodcastNotFound = errors.New("podcast not found")

type LinkReportService struct {
	repo *repository.LinkReportRepository
}

func NewLinkReportService(repo *repository.LinkReportRepository) *LinkReportService {
	return &LinkReportService{repo: repo}
}

func (s *LinkReportService) Report(ctx context.Context, podcastID string, reporterIP *string) (*model.LinkReport, error) {
	// Check if podcast exists
	exists, err := s.repo.PodcastExists(ctx, podcastID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrPodcastNotFound
	}

	return s.repo.Create(ctx, podcastID, reporterIP)
}

func (s *LinkReportService) GetByPodcastID(ctx context.Context, podcastID string) ([]model.LinkReport, error) {
	reports, err := s.repo.GetByPodcastID(ctx, podcastID)
	if err != nil {
		return nil, err
	}

	if reports == nil {
		reports = []model.LinkReport{}
	}

	return reports, nil
}

func (s *LinkReportService) CountByPodcastID(ctx context.Context, podcastID string) (int, error) {
	return s.repo.CountByPodcastID(ctx, podcastID)
}

func (s *LinkReportService) GetAllCounts(ctx context.Context) ([]model.LinkReportCount, error) {
	counts, err := s.repo.GetAllCounts(ctx)
	if err != nil {
		return nil, err
	}

	if counts == nil {
		counts = []model.LinkReportCount{}
	}

	return counts, nil
}

func (s *LinkReportService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *LinkReportService) ClearReports(ctx context.Context, podcastID string) error {
	return s.repo.DeleteByPodcastID(ctx, podcastID)
}
