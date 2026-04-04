package importer

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/BeilerOl/HobbyManager/backend/internal/model"
	"github.com/xuri/excelize/v2"
)

// PreviewRow is one sheet row mapped to WorkCreate.
type PreviewRow struct {
	RowIndex int              `json:"row_index"`
	Work     model.WorkCreate `json:"work"`
}

// ParseResult is the outcome of parsing the first worksheet.
type ParseResult struct {
	Rows          []PreviewRow `json:"rows"`
	SheetWarnings []string     `json:"sheet_warnings,omitempty"`
}

// ParseXLSX reads the first worksheet and maps columns by header name (row 1).
func ParseXLSX(data []byte) (*ParseResult, error) {
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("open spreadsheet: %w", err)
	}
	defer func() { _ = f.Close() }()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("no worksheet")
	}
	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, fmt.Errorf("read rows: %w", err)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("empty sheet")
	}

	col := mapHeaders(rows[0])
	var warnings []string
	if _, ok := col["type"]; !ok {
		warnings = append(warnings, "No column mapped for type (expected header such as: type, genre, catégorie).")
	}
	if _, ok := col["title"]; !ok {
		warnings = append(warnings, "No column mapped for title (expected header such as: titre, title, nom).")
	}

	var out []PreviewRow
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		excelRow := i + 1
		if rowIsEmpty(row) {
			continue
		}
		w := buildWorkCreate(row, col)
		out = append(out, PreviewRow{
			RowIndex: excelRow,
			Work:     w,
		})
	}
	return &ParseResult{Rows: out, SheetWarnings: warnings}, nil
}

func rowIsEmpty(row []string) bool {
	for _, c := range row {
		if strings.TrimSpace(c) != "" {
			return false
		}
	}
	return true
}

func mapHeaders(header []string) map[string]int {
	normToCol := make(map[string]int)
	for i, h := range header {
		k := normalizeHeader(h)
		if k == "" {
			continue
		}
		if _, ok := normToCol[k]; !ok {
			normToCol[k] = i
		}
	}

	fields := []struct {
		name    string
		aliases []string
	}{
		{"type", []string{"type", "genre", "categorie", "category", "catégorie"}},
		{"title", []string{"title", "titre", "nom", "name"}},
		{"authors", []string{"authors", "auteur", "auteurs", "author", "realisateur", "réalisateur"}},
		{"origin", []string{"origin", "origine"}},
		{"availability", []string{"availability", "disponibilite", "disponibilité", "ou trouver", "où trouver"}},
		{"seen", []string{"seen", "vu", "deja_vu", "déjà_vu", "statut"}},
	}

	col := make(map[string]int)
	for _, f := range fields {
		for _, a := range f.aliases {
			k := normalizeHeader(a)
			if j, ok := normToCol[k]; ok {
				col[f.name] = j
				break
			}
		}
	}
	return col
}

func normalizeHeader(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, "\u00a0", " ")
	repl := strings.NewReplacer(
		"é", "e", "è", "e", "ê", "e", "ë", "e",
		"à", "a", "â", "a",
		"ù", "u", "û", "u", "ü", "u",
		"ô", "o", "ö", "o",
		"î", "i", "ï", "i",
		"ç", "c",
	)
	s = repl.Replace(s)
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	return strings.TrimSpace(s)
}

func cellAt(row []string, col map[string]int, field string) string {
	j, ok := col[field]
	if !ok || j < 0 || j >= len(row) {
		return ""
	}
	return row[j]
}

func buildWorkCreate(row []string, col map[string]int) model.WorkCreate {
	title := strings.TrimSpace(cellAt(row, col, "title"))
	typ := parseWorkType(cellAt(row, col, "type"))
	authors := parseAuthors(cellAt(row, col, "authors"))
	origin := strings.TrimSpace(cellAt(row, col, "origin"))
	availability := strings.TrimSpace(cellAt(row, col, "availability"))
	seen := parseSeen(cellAt(row, col, "seen"))
	return model.WorkCreate{
		Type:         typ,
		Title:        title,
		Authors:      authors,
		Origin:       origin,
		Availability: availability,
		Seen:         seen,
	}
}

func parseAuthors(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return []string{}
	}
	for _, sep := range []string{";", ","} {
		if strings.Contains(s, sep) {
			var parts []string
			for _, p := range strings.Split(s, sep) {
				p = strings.TrimSpace(p)
				if p != "" {
					parts = append(parts, p)
				}
			}
			return parts
		}
	}
	return []string{s}
}

func parseSeen(s string) bool {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "true" || s == "1" || s == "yes" || s == "oui" || s == "v" || s == "x" || s == "vu" {
		return true
	}
	return false
}

func parseWorkType(raw string) model.WorkType {
	s := normalizeHeader(raw)
	if s == "" {
		return ""
	}
	switch s {
	case "roman":
		return model.WorkTypeRoman
	case "livre_culture_generale":
		return model.WorkTypeLivreCultureGenerale
	case "film":
		return model.WorkTypeFilm
	case "serie_tv":
		return model.WorkTypeSerieTV
	case "jeu_societe":
		return model.WorkTypeJeuSociete
	case "jeu_video":
		return model.WorkTypeJeuVideo
	}
	switch s {
	case "romans":
		return model.WorkTypeRoman
	case "livre culture generale", "culture generale", "essai":
		return model.WorkTypeLivreCultureGenerale
	case "cinema", "cinéma":
		return model.WorkTypeFilm
	case "serie", "série", "series", "séries", "tv":
		return model.WorkTypeSerieTV
	case "jeu de societe", "jeu societe", "boardgame":
		return model.WorkTypeJeuSociete
	case "jeu video", "jeu vidéo":
		return model.WorkTypeJeuVideo
	}
	if strings.Contains(s, "serie") || strings.Contains(s, "série") {
		return model.WorkTypeSerieTV
	}
	if strings.Contains(s, "societe") || strings.Contains(s, "société") {
		return model.WorkTypeJeuSociete
	}
	if strings.Contains(s, "video") || strings.Contains(s, "vidéo") {
		return model.WorkTypeJeuVideo
	}
	return model.WorkType("")
}
