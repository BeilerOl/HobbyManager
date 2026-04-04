package server

import (
	"net/http"

	"github.com/BeilerOl/HobbyManager/backend/internal/handler"
	"github.com/BeilerOl/HobbyManager/backend/internal/repository"
)

// NewMux returns an http.ServeMux with API routes registered.
func NewMux(repo repository.WorkRepository) *http.ServeMux {
	mux := http.NewServeMux()
	ih := &handler.ImportHandler{Repo: repo}
	mux.HandleFunc("/api/v1/works/import/preview", ih.Preview)
	mux.HandleFunc("/api/v1/works/import", ih.Import)

	wh := &handler.WorksHandler{Repo: repo}
	mux.Handle("/api/v1/works", wh)
	mux.Handle("/api/v1/works/", wh)
	return mux
}
