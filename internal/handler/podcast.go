package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/spanish-english-discord/api/internal/model"
	"github.com/spanish-english-discord/api/internal/service"
)

type PodcastHandler struct {
	service  *service.PodcastService
	validate *validator.Validate
}

func NewPodcastHandler(service *service.PodcastService) *PodcastHandler {
	return &PodcastHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *PodcastHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query params for filters and pagination
	var filters *model.PodcastFilters
	query := r.URL.Query()

	hasFilters := query.Has("language") || query.Has("level") || query.Has("country") || query.Has("topic") || query.Has("includeArchived")
	hasPagination := query.Has("page") || query.Has("pageSize")

	if hasFilters || hasPagination {
		filters = &model.PodcastFilters{}

		if lang := query.Get("language"); lang != "" {
			l := model.Language(lang)
			filters.Language = &l
		}
		if level := query.Get("level"); level != "" {
			l := model.Level(level)
			filters.Level = &l
		}
		if country := query.Get("country"); country != "" {
			filters.Country = &country
		}
		if topic := query.Get("topic"); topic != "" {
			filters.Topic = &topic
		}
		if query.Get("includeArchived") == "true" {
			filters.IncludeArchived = true
		}

		if pageStr := query.Get("page"); pageStr != "" {
			if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
				filters.Page = &page
			}
		}
		if pageSizeStr := query.Get("pageSize"); pageSizeStr != "" {
			if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
				filters.PageSize = &pageSize
			}
		}
	}

	podcasts, pagination, err := h.service.GetAll(ctx, filters)
	if err != nil {
		respondInternalError(w, err)
		return
	}

	// Return response with pagination metadata
	response := map[string]interface{}{
		"items":      podcasts,
		"pagination": pagination,
	}

	respondJSON(w, http.StatusOK, response)
}

func (h *PodcastHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	podcast, err := h.service.GetByID(ctx, id)
	if err != nil {
		if isNotFoundError(err) {
			respondNotFound(w, "podcast")
			return
		}
		respondInternalError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, podcast)
}

func (h *PodcastHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var input model.CreatePodcastInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondBadRequest(w, "invalid JSON body")
		return
	}

	if err := h.validate.Struct(input); err != nil {
		respondBadRequest(w, formatValidationError(err))
		return
	}

	podcast, err := h.service.Create(ctx, &input)
	if err != nil {
		respondInternalError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, podcast)
}

func (h *PodcastHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	var input model.UpdatePodcastInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondBadRequest(w, "invalid JSON body")
		return
	}

	if err := h.validate.Struct(input); err != nil {
		respondBadRequest(w, formatValidationError(err))
		return
	}

	podcast, err := h.service.Update(ctx, id, &input)
	if err != nil {
		if isNotFoundError(err) {
			respondNotFound(w, "podcast")
			return
		}
		respondInternalError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, podcast)
}

func (h *PodcastHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if err := h.service.Delete(ctx, id); err != nil {
		if isNotFoundError(err) {
			respondNotFound(w, "podcast")
			return
		}
		respondInternalError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PodcastHandler) Archive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	podcast, err := h.service.Archive(ctx, id, true)
	if err != nil {
		if isNotFoundError(err) {
			respondNotFound(w, "podcast")
			return
		}
		respondInternalError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, podcast)
}

func (h *PodcastHandler) Unarchive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	podcast, err := h.service.Archive(ctx, id, false)
	if err != nil {
		if isNotFoundError(err) {
			respondNotFound(w, "podcast")
			return
		}
		respondInternalError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, podcast)
}

func formatValidationError(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		if len(validationErrors) > 0 {
			e := validationErrors[0]
			return "validation failed for field '" + e.Field() + "': " + e.Tag()
		}
	}
	return "validation failed"
}
