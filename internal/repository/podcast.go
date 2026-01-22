package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spanish-english-discord/api/internal/model"
	"github.com/spanish-english-discord/api/internal/shared"
)

type PodcastRepository struct {
	pool *pgxpool.Pool
}

func NewPodcastRepository(pool *pgxpool.Pool) *PodcastRepository {
	return &PodcastRepository{pool: pool}
}

func (r *PodcastRepository) GetAll(ctx context.Context, filters *model.PodcastFilters) ([]model.Podcast, shared.OffsetPaginationResult, error) {

	query := `
		SELECT id, title, description, image_url, language, level, country, topic, url, archived, created_at, 
		updated_at
		FROM podcasts
		WHERE 1=1
	`

	whereClause := ""
	args := []interface{}{}
	argIndex := 1

	page := 1
	pageSize := 0
	offset := 0
	limit := 0

	if filters != nil {
		if filters.Page != nil && *filters.Page > 0 {
			page = *filters.Page
		}
		if filters.PageSize != nil && *filters.PageSize > 0 {
			pageSize = *filters.PageSize
			limit = pageSize
			offset = (page - 1) * pageSize
		}
	}

	includeArchived := filters != nil && filters.IncludeArchived
	if !includeArchived {
		whereClause += fmt.Sprintf(" AND archived = $%d", argIndex)
		args = append(args, false)
		argIndex++
	}

	filtersMap := map[string]interface{}{}
	if filters != nil {
		if filters.Language != nil {
			filtersMap["language"] = *filters.Language
		}
		if filters.Level != nil {
			filtersMap["level"] = *filters.Level
		}
		if filters.Country != nil {
			filtersMap["country"] = *filters.Country
		}
		if filters.Topic != nil {
			filtersMap["topic"] = *filters.Topic
		}

		qb := shared.NewQueryBuilder()
		whereClause += qb.BuildFilters(filtersMap, &argIndex, &args)
	}

	query += whereClause

	countQuery := fmt.Sprintf(`
	SELECT COUNT(*) FROM podcasts WHERE 1=1 %s
	`, whereClause)


	var totalCount int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, shared.OffsetPaginationResult{}, fmt.Errorf("failed to count podcasts: %w", err)
	}

	query += " ORDER BY created_at DESC"

	if limit > 0 || offset > 0 {
		qb := shared.NewQueryBuilder()
		query += qb.BuildOffsetPagination(offset, limit, &argIndex, &args)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, shared.OffsetPaginationResult{}, fmt.Errorf("failed to query podcasts: %w", err)
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
			return nil, shared.OffsetPaginationResult{}, fmt.Errorf("failed to scan podcast: %w", err)
		}
		podcasts = append(podcasts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, shared.OffsetPaginationResult{}, fmt.Errorf("error iterating podcasts: %w", err)
	}

	totalPages := 0
	if pageSize > 0 && totalCount > 0 {
		totalPages = totalCount / pageSize
		if totalCount%pageSize != 0 {
			totalPages++
		}
	}

	paginationResult := shared.OffsetPaginationResult{
		TotalCount:  totalCount,
		TotalPages:  totalPages,
		CurrentPage: page,
	}

	return podcasts, paginationResult, nil
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
			return nil, shared.PodcastNotFound
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
	args := []interface{}{}
	argIndex := 1

	fieldMap := map[string]interface{}{}

	if input != nil {
		if input.Title != nil {
			fieldMap["title"] = *input.Title
		}
		if input.Description != nil {
			fieldMap["description"] = *input.Description
		}
		if input.ImageURL != nil {
			fieldMap["image_url"] = *input.ImageURL
		}
		if input.Language != nil {
			fieldMap["language"] = *input.Language
		}
		if input.Level != nil {
			fieldMap["level"] = *input.Level
		}
		if input.Country != nil {
			fieldMap["country"] = *input.Country
		}
		if input.Topic != nil {
			fieldMap["topic"] = *input.Topic
		}
		if input.URL != nil {
			fieldMap["url"] = *input.URL
		}
	}

	if len(fieldMap) == 0 {
		return r.GetByID(ctx, id)
	}

	qb := shared.NewQueryBuilder()
	setClause := qb.BuildUpdates(fieldMap, &argIndex, &args)

	setClause += ", updated_at = NOW()"
	args = append(args, id)
	whereArgIndex := argIndex

	query := fmt.Sprintf(`
		UPDATE podcasts
		SET %s
		WHERE id = $%d
		RETURNING id, title, description, image_url, language, level, country, topic, url, archived, created_at, updated_at
		`, setClause, whereArgIndex)

	var p model.Podcast
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&p.ID, &p.Title, &p.Description, &p.ImageURL,
		&p.Language, &p.Level, &p.Country, &p.Topic,
		&p.URL, &p.Archived, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, shared.PodcastNotFound
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
		return shared.PodcastNotFound
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
			return nil, shared.PodcastNotFound
		}
		return nil, fmt.Errorf("failed to archive podcast: %w", err)
	}

	return &p, nil
}
