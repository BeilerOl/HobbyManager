package importxls

import (
	"bytes"
	"io"

	"github.com/extrame/xls"
)

// ReadFirstSheet reads the first sheet of an .xls file and returns rows as [][]string.
// Row i, cell j is at [i][j]; empty cells are empty strings.
// The caller can treat row 0 as headers.
func ReadFirstSheet(r io.Reader) ([][]string, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	seek := bytes.NewReader(buf)
	wb, err := xls.OpenReader(seek, "utf-8")
	if err != nil {
		return nil, err
	}
	sheet := wb.GetSheet(0)
	if sheet == nil {
		return [][]string{}, nil
	}
	maxRow := int(sheet.MaxRow) + 1
	rows := make([][]string, 0, maxRow)
	for i := 0; i < maxRow; i++ {
		row := sheet.Row(i)
		lastCol := row.LastCol()
		if lastCol < 0 {
			lastCol = 0
		}
		cells := make([]string, lastCol+1)
		for j := 0; j <= lastCol; j++ {
			cells[j] = row.Col(j)
		}
		rows = append(rows, cells)
	}
	return rows, nil
}
