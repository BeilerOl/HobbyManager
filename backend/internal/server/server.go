package server

import (
	"net/http"

	"github.com/BeilerOl/HobbyManager/backend/internal/handler"
	"github.com/BeilerOl/HobbyManager/backend/internal/repository"
)

// NewMux returns an http.ServeMux with API routes registered.
func NewMux(repo repository.WorkRepository) *http.ServeMux {
	mux := http.NewServeMux()
	wh := &handler.WorksHandler{Repo: repo}
	mux.Handle("/api/v1/works", wh)
	mux.Handle("/api/v1/works/", wh)
	return mux
}
