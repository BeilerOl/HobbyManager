package repository

import (
	"context"

	"github.com/BeilerOl/HobbyManager/backend/internal/model"
)

// MockWorkRepository is an in-memory implementation for tests.
type MockWorkRepository struct {
	Works  []*model.Work
	NextID int64
}

func (m *MockWorkRepository) List(ctx context.Context, filter *WorkFilter) ([]*model.Work, error) {
	list := make([]*model.Work, 0, len(m.Works))
	for _, w := range m.Works {
		if filter != nil {
			if filter.Type != nil && w.Type != *filter.Type {
				continue
			}
			if filter.Seen != nil && w.Seen != *filter.Seen {
				continue
			}
		}
		list = append(list, w)
	}
	return list, nil
}

func (m *MockWorkRepository) GetByID(ctx context.Context, id int64) (*model.Work, error) {
	for _, w := range m.Works {
		if w.ID == id {
			return w, nil
		}
	}
	return nil, nil
}

func (m *MockWorkRepository) Create(ctx context.Context, w *model.WorkCreate) (*model.Work, error) {
	m.NextID++
	work := &model.Work{
		ID:           m.NextID,
		Type:         w.Type,
		Title:        w.Title,
		Authors:      append([]string(nil), w.Authors...),
		Origin:       w.Origin,
		Availability: w.Availability,
		Seen:         w.Seen,
	}
	// added_at set to zero for mock; tests can ignore or set
	m.Works = append(m.Works, work)
	return work, nil
}

func (m *MockWorkRepository) CreateMany(ctx context.Context, items []*model.WorkCreate) ([]*model.Work, error) {
	out := make([]*model.Work, 0, len(items))
	for _, w := range items {
		created, err := m.Create(ctx, w)
		if err != nil {
			return nil, err
		}
		out = append(out, created)
	}
	return out, nil
}

func (m *MockWorkRepository) Update(ctx context.Context, id int64, w *model.WorkCreate) (*model.Work, error) {
	for i, work := range m.Works {
		if work.ID == id {
			m.Works[i] = &model.Work{
				ID:           id,
				Type:         w.Type,
				Title:        w.Title,
				Authors:      append([]string(nil), w.Authors...),
				Origin:       w.Origin,
				Availability: w.Availability,
				Seen:         w.Seen,
				AddedAt:      m.Works[i].AddedAt,
			}
			return m.Works[i], nil
		}
	}
	return nil, nil
}

func (m *MockWorkRepository) Delete(ctx context.Context, id int64) error {
	for i, w := range m.Works {
		if w.ID == id {
			m.Works = append(m.Works[:i], m.Works[i+1:]...)
			return nil
		}
	}
	return nil
}
