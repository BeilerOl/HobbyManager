package model

import "time"

// WorkType represents the type of cultural work (matches OpenAPI enum).
type WorkType string

const (
	WorkTypeRoman               WorkType = "roman"
	WorkTypeLivreCultureGenerale WorkType = "livre_culture_generale"
	WorkTypeFilm                WorkType = "film"
	WorkTypeSerieTV             WorkType = "serie_tv"
	WorkTypeJeuSociete          WorkType = "jeu_societe"
	WorkTypeJeuVideo            WorkType = "jeu_video"
)

// Work is the full entity as returned by the API (id and added_at set by server).
type Work struct {
	ID           int64     `json:"id"`
	Type         WorkType  `json:"type"`
	Title        string    `json:"title"`
	Authors      []string  `json:"authors"`
	AddedAt      time.Time `json:"added_at"`
	Origin       string    `json:"origin"`
	Availability string    `json:"availability"`
	Seen         bool      `json:"seen"`
}

// WorkCreate is the payload for creating or replacing a work (no id, no added_at).
type WorkCreate struct {
	Type         WorkType `json:"type"`
	Title        string   `json:"title"`
	Authors      []string `json:"authors"`
	Origin       string   `json:"origin"`
	Availability string   `json:"availability"`
	Seen         bool     `json:"seen"`
}
