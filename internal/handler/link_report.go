package handler

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/spanish-english-discord/api/internal/service"
)

type LinkReportHandler struct {
	service *service.LinkReportService
}

func NewLinkReportHandler(service *service.LinkReportService) *LinkReportHandler {
	return &LinkReportHandler{service: service}
}

func (h *LinkReportHandler) Report(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	podcastID := chi.URLParam(r, "id")

	// Get reporter IP for deduplication
	ip := getClientIP(r)

	report, err := h.service.Report(ctx, podcastID, ip)
	if err != nil {
		if isNotFoundError(err) {
			respondNotFound(w, "podcast")
			return
		}
		respondInternalError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, report)
}

func (h *LinkReportHandler) GetByPodcast(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	podcastID := chi.URLParam(r, "id")

	reports, err := h.service.GetByPodcastID(ctx, podcastID)
	if err != nil {
		respondInternalError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, reports)
}

func (h *LinkReportHandler) GetCount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	podcastID := chi.URLParam(r, "id")

	count, err := h.service.CountByPodcastID(ctx, podcastID)
	if err != nil {
		respondInternalError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]int{"count": count})
}

func (h *LinkReportHandler) GetAllCounts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	counts, err := h.service.GetAllCounts(ctx)
	if err != nil {
		respondInternalError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, counts)
}

func (h *LinkReportHandler) ClearReports(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	podcastID := chi.URLParam(r, "id")

	if err := h.service.ClearReports(ctx, podcastID); err != nil {
		respondInternalError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// getClientIP extracts the client IP from the request
func getClientIP(r *http.Request) *string {
	// Check X-Forwarded-For header first (for proxies/load balancers)
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		ips := strings.Split(xff, ",")
		ip := strings.TrimSpace(ips[0])
		return &ip
	}

	// Check X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return &xri
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	// Remove port if present
	if colonIdx := strings.LastIndex(ip, ":"); colonIdx != -1 {
		// Check if it's IPv6 with brackets
		if bracketIdx := strings.LastIndex(ip, "]"); bracketIdx != -1 && bracketIdx < colonIdx {
			ip = ip[:colonIdx]
		} else if !strings.Contains(ip, "[") {
			ip = ip[:colonIdx]
		}
	}
	return &ip
}
