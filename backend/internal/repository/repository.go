package repository

import (
	"context"

	"github.com/BeilerOl/HobbyManager/backend/internal/model"
)

// WorkRepository defines data access for works. Implementations can use SQLite, PostgreSQL, etc.
type WorkRepository interface {
	List(ctx context.Context, filter *WorkFilter) ([]*model.Work, error)
	GetByID(ctx context.Context, id int64) (*model.Work, error)
	Create(ctx context.Context, w *model.WorkCreate) (*model.Work, error)
	// CreateMany inserts all works in a single transaction. On any error, no work is persisted.
	CreateMany(ctx context.Context, items []*model.WorkCreate) ([]*model.Work, error)
	Update(ctx context.Context, id int64, w *model.WorkCreate) (*model.Work, error)
	Delete(ctx context.Context, id int64) error
}

// WorkFilter optional query filters for List.
type WorkFilter struct {
	Type *model.WorkType
	Seen *bool
}
