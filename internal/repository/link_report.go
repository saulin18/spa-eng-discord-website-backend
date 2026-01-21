package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spanish-english-discord/api/internal/model"
)

type LinkReportRepository struct {
	pool *pgxpool.Pool
}

func NewLinkReportRepository(pool *pgxpool.Pool) *LinkReportRepository {
	return &LinkReportRepository{pool: pool}
}

func (r *LinkReportRepository) Create(ctx context.Context, podcastID string, reporterIP *string) (*model.LinkReport, error) {
	query := `
		INSERT INTO link_reports (podcast_id, reporter_ip)
		VALUES ($1, $2)
		ON CONFLICT (podcast_id, reporter_ip) WHERE reporter_ip IS NOT NULL
		DO UPDATE SET created_at = link_reports.created_at
		RETURNING id, podcast_id, reporter_ip, created_at
	`

	var report model.LinkReport
	err := r.pool.QueryRow(ctx, query, podcastID, reporterIP).Scan(
		&report.ID, &report.PodcastID, &report.ReporterIP, &report.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create link report: %w", err)
	}

	return &report, nil
}

func (r *LinkReportRepository) GetByPodcastID(ctx context.Context, podcastID string) ([]model.LinkReport, error) {
	query := `
		SELECT id, podcast_id, reporter_ip, created_at
		FROM link_reports
		WHERE podcast_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, podcastID)
	if err != nil {
		return nil, fmt.Errorf("failed to query link reports: %w", err)
	}
	defer rows.Close()

	var reports []model.LinkReport
	for rows.Next() {
		var report model.LinkReport
		err := rows.Scan(&report.ID, &report.PodcastID, &report.ReporterIP, &report.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan link report: %w", err)
		}
		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating link reports: %w", err)
	}

	return reports, nil
}

func (r *LinkReportRepository) CountByPodcastID(ctx context.Context, podcastID string) (int, error) {
	query := `SELECT COUNT(*) FROM link_reports WHERE podcast_id = $1`

	var count int
	err := r.pool.QueryRow(ctx, query, podcastID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count link reports: %w", err)
	}

	return count, nil
}

func (r *LinkReportRepository) GetAllCounts(ctx context.Context) ([]model.LinkReportCount, error) {
	query := `
		SELECT podcast_id, COUNT(*) as count
		FROM link_reports
		GROUP BY podcast_id
		ORDER BY count DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query link report counts: %w", err)
	}
	defer rows.Close()

	var counts []model.LinkReportCount
	for rows.Next() {
		var c model.LinkReportCount
		err := rows.Scan(&c.PodcastID, &c.Count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan link report count: %w", err)
		}
		counts = append(counts, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating link report counts: %w", err)
	}

	return counts, nil
}

func (r *LinkReportRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM link_reports WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete link report: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *LinkReportRepository) DeleteByPodcastID(ctx context.Context, podcastID string) error {
	query := `DELETE FROM link_reports WHERE podcast_id = $1`

	_, err := r.pool.Exec(ctx, query, podcastID)
	if err != nil {
		return fmt.Errorf("failed to delete link reports: %w", err)
	}

	return nil
}

func (r *LinkReportRepository) PodcastExists(ctx context.Context, podcastID string) (bool, error) {
	query := `SELECT 1 FROM podcasts WHERE id = $1`

	var exists int
	err := r.pool.QueryRow(ctx, query, podcastID).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check podcast existence: %w", err)
	}

	return true, nil
}
