package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/BeilerOl/HobbyManager/backend/internal/model"
	"github.com/BeilerOl/HobbyManager/backend/internal/repository"
	_ "modernc.org/sqlite"
)

// DB wraps *sql.DB for SQLite-backed WorkRepository.
type DB struct {
	db *sql.DB
}

// NewDB opens a SQLite database at dataSource (e.g. "file:hobby.db") and creates the works table if needed.
func NewDB(dataSource string) (*DB, error) {
	db, err := sql.Open("sqlite", dataSource)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	d := &DB{db: db}
	if err := d.migrate(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return d, nil
}

func (d *DB) migrate() error {
	_, err := d.db.Exec(`
		CREATE TABLE IF NOT EXISTS works (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			type TEXT NOT NULL,
			title TEXT NOT NULL,
			authors TEXT NOT NULL,
			added_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			origin TEXT NOT NULL,
			availability TEXT NOT NULL,
			seen INTEGER NOT NULL DEFAULT 0
		)
	`)
	return err
}

// Close closes the database connection.
func (d *DB) Close() error {
	return d.db.Close()
}

// List returns works optionally filtered by type and seen.
func (d *DB) List(ctx context.Context, filter *repository.WorkFilter) ([]*model.Work, error) {
	query := `SELECT id, type, title, authors, added_at, origin, availability, seen FROM works WHERE 1=1`
	args := []interface{}{}
	if filter != nil {
		if filter.Type != nil {
			query += ` AND type = ?`
			args = append(args, string(*filter.Type))
		}
		if filter.Seen != nil {
			v := 0
			if *filter.Seen {
				v = 1
			}
			query += ` AND seen = ?`
			args = append(args, v)
		}
	}
	query += ` ORDER BY added_at DESC`

	rows, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*model.Work
	for rows.Next() {
		w, err := scanWork(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, w)
	}
	return list, rows.Err()
}

// GetByID returns one work by id, or nil if not found.
func (d *DB) GetByID(ctx context.Context, id int64) (*model.Work, error) {
	row := d.db.QueryRowContext(ctx,
		`SELECT id, type, title, authors, added_at, origin, availability, seen FROM works WHERE id = ?`, id)
	w, err := scanWorkRow(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return w, err
}

// Create inserts a new work and returns it with id and added_at set.
func (d *DB) Create(ctx context.Context, w *model.WorkCreate) (*model.Work, error) {
	authorsJSON, err := json.Marshal(w.Authors)
	if err != nil {
		return nil, err
	}
	res, err := d.db.ExecContext(ctx,
		`INSERT INTO works (type, title, authors, origin, availability, seen) VALUES (?, ?, ?, ?, ?, ?)`,
		string(w.Type), w.Title, string(authorsJSON), w.Origin, w.Availability, boolToInt(w.Seen))
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return d.GetByID(ctx, id)
}

// CreateMany inserts all works in one transaction.
func (d *DB) CreateMany(ctx context.Context, items []*model.WorkCreate) ([]*model.Work, error) {
	if len(items) == 0 {
		return []*model.Work{}, nil
	}
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	ids := make([]int64, len(items))
	for i, w := range items {
		authorsJSON, err := json.Marshal(w.Authors)
		if err != nil {
			return nil, err
		}
		res, err := tx.ExecContext(ctx,
			`INSERT INTO works (type, title, authors, origin, availability, seen) VALUES (?, ?, ?, ?, ?, ?)`,
			string(w.Type), w.Title, string(authorsJSON), w.Origin, w.Availability, boolToInt(w.Seen))
		if err != nil {
			return nil, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		ids[i] = id
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	out := make([]*model.Work, len(ids))
	for i, id := range ids {
		w, err := d.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}
		out[i] = w
	}
	return out, nil
}

// Update replaces a work by id and returns the updated work, or nil if not found.
func (d *DB) Update(ctx context.Context, id int64, w *model.WorkCreate) (*model.Work, error) {
	authorsJSON, err := json.Marshal(w.Authors)
	if err != nil {
		return nil, err
	}
	res, err := d.db.ExecContext(ctx,
		`UPDATE works SET type=?, title=?, authors=?, origin=?, availability=?, seen=? WHERE id=?`,
		string(w.Type), w.Title, string(authorsJSON), w.Origin, w.Availability, boolToInt(w.Seen), id)
	if err != nil {
		return nil, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return nil, nil
	}
	return d.GetByID(ctx, id)
}

// Delete removes a work by id. Returns nil if not found (no error).
func (d *DB) Delete(ctx context.Context, id int64) error {
	res, err := d.db.ExecContext(ctx, `DELETE FROM works WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return nil
	}
	return nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func scanWork(rows *sql.Rows) (*model.Work, error) {
	var id int64
	var typ, title, authorsJSON, origin, availability string
	var addedAt string
	var seen int
	if err := rows.Scan(&id, &typ, &title, &authorsJSON, &addedAt, &origin, &availability, &seen); err != nil {
		return nil, err
	}
	var authors []string
	if err := json.Unmarshal([]byte(authorsJSON), &authors); err != nil {
		return nil, err
	}
	// SQLite returns datetime as string; parse for JSON date-time
	t, err := parseSQLiteTime(addedAt)
	if err != nil {
		t = time.Now()
	}
	return &model.Work{
		ID:           id,
		Type:         model.WorkType(typ),
		Title:        title,
		Authors:      authors,
		AddedAt:      t,
		Origin:       origin,
		Availability: availability,
		Seen:         seen != 0,
	}, nil
}

func scanWorkRow(row *sql.Row) (*model.Work, error) {
	var id int64
	var typ, title, authorsJSON, origin, availability string
	var addedAt string
	var seen int
	if err := row.Scan(&id, &typ, &title, &authorsJSON, &addedAt, &origin, &availability, &seen); err != nil {
		return nil, err
	}
	var authors []string
	if err := json.Unmarshal([]byte(authorsJSON), &authors); err != nil {
		return nil, err
	}
	t, err := parseSQLiteTime(addedAt)
	if err != nil {
		t = time.Now()
	}
	return &model.Work{
		ID:           id,
		Type:         model.WorkType(typ),
		Title:        title,
		Authors:      authors,
		AddedAt:      t,
		Origin:       origin,
		Availability: availability,
		Seen:         seen != 0,
	}, nil
}

// parseSQLiteTime parses SQLite datetime string to time.Time.
func parseSQLiteTime(s string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", s)
}
