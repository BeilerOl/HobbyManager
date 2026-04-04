package handler

import (
	"errors"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/BeilerOl/HobbyManager/backend/internal/model"
	"github.com/BeilerOl/HobbyManager/backend/internal/repository"
)

// WorksHandler handles /api/v1/works and /api/v1/works/{id}.
type WorksHandler struct {
	Repo repository.WorkRepository
}

// ServeHTTP dispatches by method and path.
func (h *WorksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/works")
	path = strings.Trim(path, "/")

	if path == "" {
		switch r.Method {
		case http.MethodGet:
			h.List(w, r)
			return
		case http.MethodPost:
			h.Create(w, r)
			return
		}
	} else {
		id, err := strconv.ParseInt(path, 10, 64)
		if err != nil {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "work not found")
			return
		}
		switch r.Method {
		case http.MethodGet:
			h.Get(w, r, id)
			return
		case http.MethodPut:
			h.Update(w, r, id)
			return
		case http.MethodDelete:
			h.Delete(w, r, id)
			return
		}
	}

	writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
}

func (h *WorksHandler) List(w http.ResponseWriter, r *http.Request) {
	var filter repository.WorkFilter
	if t := r.URL.Query().Get("type"); t != "" {
		wt := model.WorkType(t)
		filter.Type = &wt
	}
	if s := r.URL.Query().Get("seen"); s != "" {
		seen := strings.EqualFold(s, "true") || s == "1"
		filter.Seen = &seen
	}

	list, err := h.Repo.List(r.Context(), &filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	if list == nil {
		list = []*model.Work{}
	}
	writeJSON(w, http.StatusOK, list)
}

func (h *WorksHandler) Get(w http.ResponseWriter, r *http.Request, id int64) {
	work, err := h.Repo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	if work == nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "work not found")
		return
	}
	writeJSON(w, http.StatusOK, work)
}

func (h *WorksHandler) Create(w http.ResponseWriter, r *http.Request) {
	var create model.WorkCreate
	if err := json.NewDecoder(r.Body).Decode(&create); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid JSON: "+err.Error())
		return
	}
	if err := validateWorkCreate(&create); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	work, err := h.Repo.Create(r.Context(), &create)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, work)
}

func (h *WorksHandler) Update(w http.ResponseWriter, r *http.Request, id int64) {
	var create model.WorkCreate
	if err := json.NewDecoder(r.Body).Decode(&create); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid JSON: "+err.Error())
		return
	}
	if err := validateWorkCreate(&create); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	work, err := h.Repo.Update(r.Context(), id, &create)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	if work == nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "work not found")
		return
	}
	writeJSON(w, http.StatusOK, work)
}

func (h *WorksHandler) Delete(w http.ResponseWriter, r *http.Request, id int64) {
	existing, err := h.Repo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	if existing == nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "work not found")
		return
	}
	if err := h.Repo.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

var (
	ErrInvalidType     = errors.New("invalid type")
	ErrTitleRequired   = errors.New("title is required")
)

var validTypes = map[model.WorkType]bool{
	model.WorkTypeRoman:                true,
	model.WorkTypeLivreCultureGenerale: true,
	model.WorkTypeFilm:                 true,
	model.WorkTypeSerieTV:              true,
	model.WorkTypeJeuSociete:           true,
	model.WorkTypeJeuVideo:             true,
}

// ValidateWorkCreateErrors returns validation issues for import preview (empty if valid).
func ValidateWorkCreateErrors(w *model.WorkCreate) []string {
	var errs []string
	if !validTypes[w.Type] {
		errs = append(errs, ErrInvalidType.Error())
	}
	if strings.TrimSpace(w.Title) == "" {
		errs = append(errs, ErrTitleRequired.Error())
	}
	return errs
}

func validateWorkCreate(w *model.WorkCreate) error {
	errs := ValidateWorkCreateErrors(w)
	if len(errs) == 0 {
		return nil
	}
	return errors.New(errs[0])
}
