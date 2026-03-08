package sqlite

import (
	"context"
	"testing"

	"github.com/BeilerOl/HobbyManager/backend/internal/model"
)

func TestDB_CRUD(t *testing.T) {
	db, err := NewDB("file::memory:?cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	ctx := context.Background()

	created, err := db.Create(ctx, &model.WorkCreate{
		Type:         model.WorkTypeRoman,
		Title:        "Test Book",
		Authors:      []string{"Author A", "Author B"},
		Origin:       "friend",
		Availability: "library",
		Seen:         false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.ID == 0 || created.Title != "Test Book" || len(created.Authors) != 2 {
		t.Fatalf("unexpected created: %+v", created)
	}

	got, err := db.GetByID(ctx, created.ID)
	if err != nil || got == nil || got.Title != "Test Book" {
		t.Fatalf("GetByID: %v %+v", err, got)
	}

	list, err := db.List(ctx, nil)
	if err != nil || len(list) != 1 {
		t.Fatalf("List: %v len=%d", err, len(list))
	}

	updated, err := db.Update(ctx, created.ID, &model.WorkCreate{
		Type: model.WorkTypeFilm, Title: "Updated", Authors: []string{"X"}, Origin: "", Availability: "", Seen: true,
	})
	if err != nil || updated == nil || updated.Title != "Updated" || !updated.Seen {
		t.Fatalf("Update: %v %+v", err, updated)
	}

	err = db.Delete(ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	got, _ = db.GetByID(ctx, created.ID)
	if got != nil {
		t.Fatal("expected nil after delete")
	}
}
