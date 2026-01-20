package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/spanish-english-discord/api/internal/database"
	"github.com/spanish-english-discord/api/internal/repository"
	"github.com/spanish-english-discord/api/internal/service"
)

type Handler struct {
	router   *chi.Mux
	db       *database.DB
	podcasts *PodcastHandler
}

func New(db *database.DB) *Handler {
	h := &Handler{
		router: chi.NewRouter(),
		db:     db,
	}

	// Initialize repositories
	podcastRepo := repository.NewPodcastRepository(db.Pool)

	// Initialize services
	podcastService := service.NewPodcastService(podcastRepo)

	// Initialize handlers
	h.podcasts = NewPodcastHandler(podcastService)

	h.setupMiddleware()
	h.setupRoutes()

	return h
}

func (h *Handler) setupMiddleware() {
	h.router.Use(middleware.RequestID)
	h.router.Use(middleware.RealIP)
	h.router.Use(middleware.Logger)
	h.router.Use(middleware.Recoverer)
	h.router.Use(middleware.Compress(5))

	h.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Configure for production
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
}

func (h *Handler) setupRoutes() {
	h.router.Get("/health", h.healthCheck)

	h.router.Route("/api", func(r chi.Router) {
		r.Route("/podcasts", func(r chi.Router) {
			r.Get("/", h.podcasts.GetAll)
			r.Post("/", h.podcasts.Create)
			r.Get("/{id}", h.podcasts.GetByID)
			r.Patch("/{id}", h.podcasts.Update)
			r.Delete("/{id}", h.podcasts.Delete)
		})
	})
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	if err := h.db.Health(r.Context()); err != nil {
		slog.Error("health check failed", "error", err)
		respondError(w, http.StatusServiceUnavailable, "database unavailable")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
	})
}

// Response helpers

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			slog.Error("failed to encode response", "error", err)
		}
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	})
}

func respondNotFound(w http.ResponseWriter, resource string) {
	respondError(w, http.StatusNotFound, resource+" not found")
}

func respondBadRequest(w http.ResponseWriter, message string) {
	respondError(w, http.StatusBadRequest, message)
}

func respondInternalError(w http.ResponseWriter, err error) {
	slog.Error("internal error", "error", err)
	respondError(w, http.StatusInternalServerError, "internal server error")
}

func isNotFoundError(err error) bool {
	return errors.Is(err, repository.ErrNotFound)
}
