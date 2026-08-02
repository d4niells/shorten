package handler

import (
	"html/template"
	"net/http"

	"github.com/d4niells/shorten/internal/service"
)

type WebHandler struct {
	templates  *template.Template
	urlService service.URLService
}

type PageData struct {
	ShortURL string
	LongURL  string
	Error    string
}

func NewWebHandler(tmpl *template.Template, urlService service.URLService) *WebHandler {
	return &WebHandler{
		templates:  tmpl,
		urlService: urlService,
	}
}

func (h *WebHandler) Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := PageData{}

	// Handle POST request (form submission)
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			data.Error = "Failed to parse form"
		} else {
			longURL := r.FormValue("longUrl")
			if longURL == "" {
				data.Error = "Please enter a URL"
			} else {
				// Call the URL shortening service
				url, err := h.urlService.Shorten(r.Context(), longURL)
				if err != nil {
					data.Error = err.Error()
				} else {
					data.ShortURL = url.ShortURL
					data.LongURL = url.LongURL
				}
			}
		}
	}

	if err := h.templates.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
