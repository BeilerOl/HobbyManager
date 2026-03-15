package importxls

import (
	"strings"

	"github.com/BeilerOl/HobbyManager/backend/internal/model"
)

// PreviewRow is one row in the import preview (row index, optional error, optional mapped work).
type PreviewRow struct {
	RowIndex int              `json:"row_index"`
	Error    string           `json:"error,omitempty"`
	Work     *model.WorkCreate `json:"work,omitempty"`
}

// MapRows uses the first row as headers, maps each data row to WorkCreate, and returns preview rows.
// Column names are matched case-insensitive and trimmed. Known headers: titre/title/nom, auteur(s)/author,
// type/catégorie, origine/origin/source, disponibilité/availability/statut, vu/lu/seen/fait.
func MapRows(rows [][]string) []PreviewRow {
	if len(rows) < 2 {
		return nil
	}
	headers := rows[0]
	col := make(map[string]int)
	for i, h := range headers {
		key := strings.TrimSpace(strings.ToLower(h))
		if key != "" {
			col[key] = i
		}
	}
	// Map common header variants to our field names (exact match first, then contains)
	titleCol := firstColOrContains(col, []string{"titre", "title", "nom"}, nil)
	authorsCol := firstColOrContains(col, []string{"auteur(s)", "auteur", "auteurs", "author", "authors"}, []string{"auteur", "author"})
	typeCol := firstColOrContains(col, []string{"type", "catégorie", "categorie", "category"}, nil)
	originCol := firstColOrContains(col, []string{"origine", "origin", "source"}, nil)
	availCol := firstColOrContains(col, []string{"disponibilité", "disponibilite", "availability", "statut", "status"}, nil)
	seenCol := firstColOrContains(col, []string{"vu", "lu", "seen", "fait", "read", "viewed"}, nil)

	result := make([]PreviewRow, 0, len(rows)-1)
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		pr := PreviewRow{RowIndex: i}
		w, err := mapRow(row, titleCol, authorsCol, typeCol, originCol, availCol, seenCol)
		if err != nil {
			pr.Error = err.Error()
		} else {
			pr.Work = w
		}
		result = append(result, pr)
	}
	return result
}

func firstCol(col map[string]int, keys ...string) int {
	for _, k := range keys {
		if i, ok := col[k]; ok {
			return i
		}
	}
	return -1
}

// firstColOrContains returns the column index: first by exact key match, then by header containing any of the substrings.
func firstColOrContains(col map[string]int, exactKeys []string, containsSubstrs []string) int {
	if i := firstCol(col, exactKeys...); i >= 0 {
		return i
	}
	if len(containsSubstrs) == 0 {
		return -1
	}
	for header, idx := range col {
		h := strings.ToLower(header)
		for _, sub := range containsSubstrs {
			if strings.Contains(h, sub) {
				return idx
			}
		}
	}
	return -1
}

func cell(row []string, col int) string {
	if col < 0 || col >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[col])
}

func mapRow(row []string, titleCol, authorsCol, typeCol, originCol, availCol, seenCol int) (*model.WorkCreate, error) {
	title := cell(row, titleCol)
	if title == "" {
		return nil, nil // skip empty rows without error; caller can treat as "skip"
	}
	typ := normalizeType(cell(row, typeCol))
	if typ == "" {
		typ = string(model.WorkTypeRoman) // default
	}
	if !validWorkType(model.WorkType(typ)) {
		return nil, &mapError{"type invalide: " + cell(row, typeCol)}
	}
	authors := parseAuthors(cell(row, authorsCol))
	if len(authors) == 0 {
		authors = []string{""}
	}
	origin := cell(row, originCol)
	if origin == "" {
		origin = "import"
	}
	availability := cell(row, availCol)
	if availability == "" {
		availability = "inconnu"
	}
	seen := parseSeen(cell(row, seenCol))
	return &model.WorkCreate{
		Type:         model.WorkType(typ),
		Title:        title,
		Authors:      authors,
		Origin:       origin,
		Availability: availability,
		Seen:         seen,
	}, nil
}

type mapError struct{ msg string }

func (e *mapError) Error() string { return e.msg }

var typeMap = map[string]string{
	"roman":                    "roman",
	"livre":                    "livre_culture_generale",
	"livre_culture_generale":   "livre_culture_generale",
	"culture générale":         "livre_culture_generale",
	"culture generale":        "livre_culture_generale",
	"bd":                       "livre_culture_generale",
	"manga":                    "livre_culture_generale",
	"film":                     "film",
	"serie":                    "serie_tv",
	"série":                    "serie_tv",
	"serie tv":                 "serie_tv",
	"série tv":                 "serie_tv",
	"jeu de société":          "jeu_societe",
	"jeu de societe":          "jeu_societe",
	"jeu societe":             "jeu_societe",
	"jeu_societe":             "jeu_societe",
	"jeu vidéo":               "jeu_video",
	"jeu video":               "jeu_video",
	"jeu_video":               "jeu_video",
}

func normalizeType(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	if v, ok := typeMap[s]; ok {
		return v
	}
	// try exact
	switch s {
	case "roman", "film", "serie_tv", "jeu_societe", "jeu_video", "livre_culture_generale":
		return s
	}
	return ""
}

func validWorkType(t model.WorkType) bool {
	_, ok := map[model.WorkType]bool{
		model.WorkTypeRoman:                true,
		model.WorkTypeLivreCultureGenerale: true,
		model.WorkTypeFilm:                 true,
		model.WorkTypeSerieTV:              true,
		model.WorkTypeJeuSociete:           true,
		model.WorkTypeJeuVideo:             true,
	}[t]
	return ok
}

func parseAuthors(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	// Normalize: replace " et ", newlines, " - " by a single separator for consistent split
	s = strings.ReplaceAll(s, "\r\n", ",")
	s = strings.ReplaceAll(s, "\n", ",")
	s = strings.ReplaceAll(s, "\r", ",")
	s = strings.ReplaceAll(s, " et ", ",")
	s = strings.ReplaceAll(s, " - ", ",")
	parts := strings.FieldsFunc(s, func(r rune) bool { return r == ',' || r == ';' || r == '|' })
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		// Single value with no separator, or only separators
		return []string{strings.TrimSpace(s)}
	}
	return out
}

func parseSeen(s string) bool {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return false
	}
	switch s {
	case "1", "oui", "yes", "true", "x", "v", "lu", "vu", "fait", "read", "viewed":
		return true
	}
	return false
}
