package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/BeilerOl/HobbyManager/backend/internal/importer"
	"github.com/BeilerOl/HobbyManager/backend/internal/model"
	"github.com/BeilerOl/HobbyManager/backend/internal/repository"
)

const defaultImportMaxBytes = 10 << 20 // 10 MiB

// ImportHandler handles POST /api/v1/works/import/preview and POST /api/v1/works/import.
type ImportHandler struct {
	Repo     repository.WorkRepository
	MaxBytes int64
}

func (h *ImportHandler) maxBytes() int64 {
	if h.MaxBytes <= 0 {
		return defaultImportMaxBytes
	}
	return h.MaxBytes
}

// Preview handles multipart upload and returns parsed rows with validation hints.
func (h *ImportHandler) Preview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}
	if err := r.ParseMultipartForm(h.maxBytes()); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid multipart form: "+err.Error())
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "missing file field")
		return
	}
	defer func() { _ = file.Close() }()

	data, err := io.ReadAll(io.LimitReader(file, h.maxBytes()))
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "cannot read file: "+err.Error())
		return
	}

	res, err := importer.ParseXLSX(data)
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	type rowOut struct {
		RowIndex int              `json:"row_index"`
		Work     model.WorkCreate `json:"work"`
		Errors   []string         `json:"errors"`
	}
	out := struct {
		SheetWarnings []string `json:"sheet_warnings,omitempty"`
		Rows          []rowOut `json:"rows"`
	}{
		SheetWarnings: res.SheetWarnings,
	}

	for _, pr := range res.Rows {
		errs := ValidateWorkCreateErrors(&pr.Work)
		out.Rows = append(out.Rows, rowOut{
			RowIndex: pr.RowIndex,
			Work:     pr.Work,
			Errors:   errs,
		})
	}

	writeJSON(w, http.StatusOK, out)
}

// ImportBody is the JSON body for POST /api/v1/works/import.
type ImportBody struct {
	Items []model.WorkCreate `json:"items"`
}

// ImportResponse is the JSON body for a successful import.
type ImportResponse struct {
	Created []*model.Work       `json:"created"`
	Failed  []ImportFailureItem `json:"failed"`
}

// ImportFailureItem describes a validation failure for one item index.
type ImportFailureItem struct {
	Index   int    `json:"index"`
	Message string `json:"message"`
}

// Import validates and inserts works in a transaction (all valid rows or none on DB error).
func (h *ImportHandler) Import(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}
	var body ImportBody
	if err := json.NewDecoder(io.LimitReader(r.Body, h.maxBytes())).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid JSON: "+err.Error())
		return
	}

	var failed []ImportFailureItem
	var toCreate []*model.WorkCreate
	for i, item := range body.Items {
		wc := item
		if errs := ValidateWorkCreateErrors(&wc); len(errs) > 0 {
			failed = append(failed, ImportFailureItem{
				Index:   i,
				Message: strings.Join(errs, "; "),
			})
			continue
		}
		copyWC := wc
		toCreate = append(toCreate, &copyWC)
	}

	resp := ImportResponse{Failed: failed}
	if len(toCreate) == 0 {
		writeJSON(w, http.StatusOK, resp)
		return
	}

	created, err := h.Repo.CreateMany(r.Context(), toCreate)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	resp.Created = created
	writeJSON(w, http.StatusOK, resp)
}
