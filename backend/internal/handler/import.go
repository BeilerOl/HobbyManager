package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/BeilerOl/HobbyManager/backend/internal/importxls"
	"github.com/BeilerOl/HobbyManager/backend/internal/repository"
)

const importMaxMemory = 10 << 20 // 10 MB

// ImportHandler handles POST /api/v1/import/preview and /api/v1/import/execute.
type ImportHandler struct {
	Repo repository.WorkRepository
}

// ServeHTTP dispatches to Preview or Execute by path.
func (h *ImportHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed")
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/import")
	path = strings.Trim(path, "/")
	switch path {
	case "preview":
		h.Preview(w, r)
		return
	case "execute":
		h.Execute(w, r)
		return
	}
	writeError(w, http.StatusNotFound, "NOT_FOUND", "not found")
}

// ImportPreviewResponse matches OpenAPI ImportPreviewResponse.
type ImportPreviewResponse struct {
	Rows []ImportPreviewRow `json:"rows"`
}

// ImportPreviewRow matches OpenAPI ImportPreviewRow.
type ImportPreviewRow struct {
	RowIndex int                  `json:"row_index"`
	Error    string               `json:"error,omitempty"`
	Work     *importWorkCreateJSON `json:"work,omitempty"`
}

type importWorkCreateJSON struct {
	Type         string   `json:"type"`
	Title        string   `json:"title"`
	Authors      []string `json:"authors"`
	Origin       string   `json:"origin"`
	Availability string   `json:"availability"`
	Seen         bool     `json:"seen"`
}

// ImportExecuteResponse matches OpenAPI ImportExecuteResponse.
type ImportExecuteResponse struct {
	Created int                     `json:"created"`
	Errors  []ImportExecuteErrorRow `json:"errors,omitempty"`
}

// ImportExecuteErrorRow matches OpenAPI ImportExecuteErrorRow.
type ImportExecuteErrorRow struct {
	RowIndex int    `json:"row_index"`
	Message  string `json:"message"`
}

func (h *ImportHandler) Preview(w http.ResponseWriter, r *http.Request) {
	rows, err := h.parseFileRows(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_FILE", err.Error())
		return
	}
	preview := importxls.MapRows(rows)
	resp := ImportPreviewResponse{Rows: make([]ImportPreviewRow, len(preview))}
	for i, p := range preview {
		resp.Rows[i] = ImportPreviewRow{RowIndex: p.RowIndex, Error: p.Error}
		if p.Work != nil {
			resp.Rows[i].Work = &importWorkCreateJSON{
				Type:         string(p.Work.Type),
				Title:        p.Work.Title,
				Authors:      p.Work.Authors,
				Origin:       p.Work.Origin,
				Availability: p.Work.Availability,
				Seen:         p.Work.Seen,
			}
		}
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *ImportHandler) Execute(w http.ResponseWriter, r *http.Request) {
	rows, err := h.parseFileRows(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_FILE", err.Error())
		return
	}
	preview := importxls.MapRows(rows)

	var selectedIndices map[int]bool
	if s := r.FormValue("row_indices"); s != "" {
		var indices []int
		if err := json.Unmarshal([]byte(s), &indices); err == nil {
			selectedIndices = make(map[int]bool)
			for _, i := range indices {
				selectedIndices[i] = true
			}
		}
	}

	var created int
	var errRows []ImportExecuteErrorRow
	for _, p := range preview {
		if p.Work == nil || p.Error != "" {
			if p.Error != "" {
				errRows = append(errRows, ImportExecuteErrorRow{RowIndex: p.RowIndex, Message: p.Error})
			}
			continue
		}
		if selectedIndices != nil && !selectedIndices[p.RowIndex] {
			continue
		}
		_, err := h.Repo.Create(r.Context(), p.Work)
		if err != nil {
			errRows = append(errRows, ImportExecuteErrorRow{RowIndex: p.RowIndex, Message: err.Error()})
			continue
		}
		created++
	}
	writeJSON(w, http.StatusOK, ImportExecuteResponse{Created: created, Errors: errRows})
}

// parseFileRows parses the multipart form and returns sheet rows from the "file" part.
func (h *ImportHandler) parseFileRows(r *http.Request) ([][]string, error) {
	if err := r.ParseMultipartForm(importMaxMemory); err != nil {
		return nil, err
	}
	f, _, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	rows, err := importxls.ReadFirstSheet(f)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
