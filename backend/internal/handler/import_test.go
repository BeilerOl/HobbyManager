package handler

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BeilerOl/HobbyManager/backend/internal/model"
	"github.com/BeilerOl/HobbyManager/backend/internal/repository"
	"github.com/xuri/excelize/v2"
)

func testXLSX(t *testing.T) []byte {
	t.Helper()
	f := excelize.NewFile()
	sheet := f.GetSheetName(0)
	_ = f.SetSheetRow(sheet, "A1", &[]interface{}{"Titre", "Type", "Auteurs", "Origine", "Disponibilité", "Vu"})
	_ = f.SetSheetRow(sheet, "A2", &[]interface{}{"Le Test", "film", "A, B", "ami", "librairie", "non"})
	_ = f.SetSheetRow(sheet, "A3", &[]interface{}{"", "roman", "", "", "", ""})
	buf, err := f.WriteToBuffer()
	if err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestImportHandler_Preview(t *testing.T) {
	repo := &repository.MockWorkRepository{NextID: 0}
	h := &ImportHandler{Repo: repo}

	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	part, err := w.CreateFormFile("file", "test.xlsx")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write(testXLSX(t)); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/works/import/preview", &body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	rec := httptest.NewRecorder()
	h.Preview(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}
	var out struct {
		Rows []struct {
			RowIndex int              `json:"row_index"`
			Work     model.WorkCreate `json:"work"`
			Errors   []string         `json:"errors"`
		} `json:"rows"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	if len(out.Rows) < 1 {
		t.Fatalf("expected at least one row, got %d", len(out.Rows))
	}
	if out.Rows[0].Work.Title != "Le Test" || out.Rows[0].Work.Type != model.WorkTypeFilm {
		t.Fatalf("unexpected first row: %+v", out.Rows[0])
	}
}

func TestImportHandler_Import(t *testing.T) {
	repo := &repository.MockWorkRepository{NextID: 0}
	h := &ImportHandler{Repo: repo}

	body := `{"items":[{"type":"film","title":"X","authors":[],"origin":"","availability":"","seen":false},{"type":"bad","title":"Y","authors":[],"origin":"","availability":"","seen":false}]}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/works/import", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Import(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}
	var out ImportResponse
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	if len(out.Created) != 1 || out.Created[0].Title != "X" {
		t.Fatalf("created: %+v", out.Created)
	}
	if len(out.Failed) != 1 || out.Failed[0].Index != 1 {
		t.Fatalf("failed: %+v", out.Failed)
	}
}

func TestImporter_ParseXLSX_EmptyTypeRowSkippedValidation(t *testing.T) {
	// second row has empty title — should still appear with errors
	repo := &repository.MockWorkRepository{NextID: 0}
	h := &ImportHandler{Repo: repo}

	f := excelize.NewFile()
	sheet := f.GetSheetName(0)
	_ = f.SetSheetRow(sheet, "A1", &[]interface{}{"titre", "type"})
	_ = f.SetSheetRow(sheet, "A2", &[]interface{}{"Ok", "jeu_video"})
	_ = f.SetSheetRow(sheet, "A3", &[]interface{}{"", "film"})
	buf, err := f.WriteToBuffer()
	if err != nil {
		t.Fatal(err)
	}

	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	part, _ := w.CreateFormFile("file", "t.xlsx")
	_, _ = part.Write(buf.Bytes())
	_ = w.Close()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/works/import/preview", &body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	rec := httptest.NewRecorder()
	h.Preview(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatal(rec.Body.String())
	}
	var out struct {
		Rows []struct {
			Errors []string `json:"errors"`
		} `json:"rows"`
	}
	_ = json.NewDecoder(rec.Body).Decode(&out)
	if len(out.Rows) != 2 {
		t.Fatalf("rows: %d", len(out.Rows))
	}
	if len(out.Rows[1].Errors) == 0 {
		t.Fatal("expected validation error on empty title row")
	}
}

func TestImporter_ParseXLSX_NoFile(t *testing.T) {
	h := &ImportHandler{Repo: &repository.MockWorkRepository{}}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/works/import/preview", nil)
	rec := httptest.NewRecorder()
	h.Preview(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got %d", rec.Code)
	}
}

func TestImportHandler_Import_EmptyBody(t *testing.T) {
	h := &ImportHandler{Repo: &repository.MockWorkRepository{}}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/works/import", bytes.NewReader([]byte(`{"items":[]}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Import(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatal(rec.Code)
	}
}

