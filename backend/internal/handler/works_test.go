package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BeilerOl/HobbyManager/backend/internal/model"
	"github.com/BeilerOl/HobbyManager/backend/internal/repository"
)

func TestWorksHandler_List(t *testing.T) {
	repo := &repository.MockWorkRepository{NextID: 0}
	_, _ = repo.Create(nil, &model.WorkCreate{
		Type: model.WorkTypeRoman, Title: "Test", Authors: []string{"A"}, Origin: "", Availability: "", Seen: false,
	})
	h := &WorksHandler{Repo: repo}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/works", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d", rec.Code)
	}
	var list []*model.Work
	if err := json.NewDecoder(rec.Body).Decode(&list); err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].Title != "Test" {
		t.Fatalf("unexpected list: %+v", list)
	}
}

func TestWorksHandler_Create(t *testing.T) {
	repo := &repository.MockWorkRepository{NextID: 0}
	h := &WorksHandler{Repo: repo}

	body := `{"type":"film","title":"Un film","authors":["X"],"origin":"ami","availability":"librairie","seen":false}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/works", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("got status %d body %s", rec.Code, rec.Body.String())
	}
	var work model.Work
	if err := json.NewDecoder(rec.Body).Decode(&work); err != nil {
		t.Fatal(err)
	}
	if work.ID != 1 || work.Title != "Un film" {
		t.Fatalf("unexpected work: %+v", work)
	}
}

func TestWorksHandler_Get_NotFound(t *testing.T) {
	repo := &repository.MockWorkRepository{}
	h := &WorksHandler{Repo: repo}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/works/999", nil)
	rec := httptest.NewRecorder()
	// Simulate path already trimmed to "999"
	req.URL.Path = "/api/v1/works/999"
	h.Get(rec, req, 999)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("got status %d", rec.Code)
	}
}
