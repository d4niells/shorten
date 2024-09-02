package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/d4niells/shorten/internal/entity"
	"github.com/d4niells/shorten/internal/service"
)

type URLHandler struct {
	urlService service.URLService
}

func NewURLHandler(urlService service.URLService) *URLHandler {
	return &URLHandler{urlService}
}

func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req struct {
		LongURL string `json:"long_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	url, err := h.urlService.Shorten(r.Context(), req.LongURL)
	if err != nil {
		if errors.Is(err, entity.ErrEmptyLongURL) {
			http.Error(w, "long_url cannot be empty", http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(url)
}
