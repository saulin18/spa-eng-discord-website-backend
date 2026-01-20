package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spanish-english-discord/api/internal/model"
)

var ErrNotFound = errors.New("resource not found")

type PodcastRepository struct {
	pool *pgxpool.Pool
}

func NewPodcastRepository(pool *pgxpool.Pool) *PodcastRepository {
	return &PodcastRepository{pool: pool}
}

func (r *PodcastRepository) GetAll(ctx context.Context, filters *model.PodcastFilters) ([]model.Podcast, error) {
	query := `
		SELECT id, title, description, image_url, language, level, country, topic, url, archived, created_at, updated_at
		FROM podcasts
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	// Filter out archived by default
	includeArchived := filters != nil && filters.IncludeArchived
	if !includeArchived {
		query += fmt.Sprintf(" AND archived = $%d", argIndex)
		args = append(args, false)
		argIndex++
	}

	if filters != nil {
		if filters.Language != nil {
			query += fmt.Sprintf(" AND language = $%d", argIndex)
			args = append(args, *filters.Language)
			argIndex++
		}
		if filters.Level != nil {
			query += fmt.Sprintf(" AND level = $%d", argIndex)
			args = append(args, *filters.Level)
			argIndex++
		}
		if filters.Country != nil {
			query += fmt.Sprintf(" AND country = $%d", argIndex)
			args = append(args, *filters.Country)
			argIndex++
		}
		if filters.Topic != nil {
			query += fmt.Sprintf(" AND topic = $%d", argIndex)
			args = append(args, *filters.Topic)
			argIndex++
		}
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query podcasts: %w", err)
	}
	defer rows.Close()

	var podcasts []model.Podcast
	for rows.Next() {
		var p model.Podcast
		err := rows.Scan(
			&p.ID, &p.Title, &p.Description, &p.ImageURL,
			&p.Language, &p.Level, &p.Country, &p.Topic,
			&p.URL, &p.Archived, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan podcast: %w", err)
		}
		podcasts = append(podcasts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating podcasts: %w", err)
	}

	return podcasts, nil
}

func (r *PodcastRepository) GetByID(ctx context.Context, id string) (*model.Podcast, error) {
	query := `
		SELECT id, title, description, image_url, language, level, country, topic, url, archived, created_at, updated_at
		FROM podcasts
		WHERE id = $1
	`

	var p model.Podcast
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.Title, &p.Description, &p.ImageURL,
		&p.Language, &p.Level, &p.Country, &p.Topic,
		&p.URL, &p.Archived, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get podcast: %w", err)
	}

	return &p, nil
}

func (r *PodcastRepository) Create(ctx context.Context, input *model.CreatePodcastInput) (*model.Podcast, error) {
	query := `
		INSERT INTO podcasts (title, description, image_url, language, level, country, topic, url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, title, description, image_url, language, level, country, topic, url, archived, created_at, updated_at
	`

	var p model.Podcast
	err := r.pool.QueryRow(ctx, query,
		input.Title, input.Description, input.ImageURL,
		input.Language, input.Level, input.Country, input.Topic, input.URL,
	).Scan(
		&p.ID, &p.Title, &p.Description, &p.ImageURL,
		&p.Language, &p.Level, &p.Country, &p.Topic,
		&p.URL, &p.Archived, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create podcast: %w", err)
	}

	return &p, nil
}

func (r *PodcastRepository) Update(ctx context.Context, id string, input *model.UpdatePodcastInput) (*model.Podcast, error) {
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if input.Title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, *input.Title)
		argIndex++
	}
	if input.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *input.Description)
		argIndex++
	}
	if input.ImageURL != nil {
		setParts = append(setParts, fmt.Sprintf("image_url = $%d", argIndex))
		args = append(args, *input.ImageURL)
		argIndex++
	}
	if input.Language != nil {
		setParts = append(setParts, fmt.Sprintf("language = $%d", argIndex))
		args = append(args, *input.Language)
		argIndex++
	}
	if input.Level != nil {
		setParts = append(setParts, fmt.Sprintf("level = $%d", argIndex))
		args = append(args, *input.Level)
		argIndex++
	}
	if input.Country != nil {
		setParts = append(setParts, fmt.Sprintf("country = $%d", argIndex))
		args = append(args, *input.Country)
		argIndex++
	}
	if input.Topic != nil {
		setParts = append(setParts, fmt.Sprintf("topic = $%d", argIndex))
		args = append(args, *input.Topic)
		argIndex++
	}
	if input.URL != nil {
		setParts = append(setParts, fmt.Sprintf("url = $%d", argIndex))
		args = append(args, *input.URL)
		argIndex++
	}

	if len(setParts) == 0 {
		return r.GetByID(ctx, id)
	}

	setParts = append(setParts, "updated_at = NOW()")
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE podcasts
		SET %s
		WHERE id = $%d
		RETURNING id, title, description, image_url, language, level, country, topic, url, archived, created_at, updated_at
	`, strings.Join(setParts, ", "), argIndex)

	var p model.Podcast
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&p.ID, &p.Title, &p.Description, &p.ImageURL,
		&p.Language, &p.Level, &p.Country, &p.Topic,
		&p.URL, &p.Archived, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to update podcast: %w", err)
	}

	return &p, nil
}

func (r *PodcastRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM podcasts WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete podcast: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *PodcastRepository) Archive(ctx context.Context, id string, archived bool) (*model.Podcast, error) {
	query := `
		UPDATE podcasts
		SET archived = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, title, description, image_url, language, level, country, topic, url, archived, created_at, updated_at
	`

	var p model.Podcast
	err := r.pool.QueryRow(ctx, query, archived, id).Scan(
		&p.ID, &p.Title, &p.Description, &p.ImageURL,
		&p.Language, &p.Level, &p.Country, &p.Topic,
		&p.URL, &p.Archived, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to archive podcast: %w", err)
	}

	return &p, nil
}
