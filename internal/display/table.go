package display

import (
	"fmt"
	"strings"
)

type Table struct {
	headers []string
	rows    [][]string
}

func NewTable(headers ...string) *Table {
	return &Table{headers: headers}
}

func (t *Table) AddRow(cols ...string) {
	t.rows = append(t.rows, cols)
}

func (t *Table) Render() string {
	if len(t.headers) == 0 {
		return ""
	}

	// Calculate column widths
	widths := make([]int, len(t.headers))
	for i, h := range t.headers {
		widths[i] = len(h)
	}
	for _, row := range t.rows {
		for i, col := range row {
			if i < len(widths) && len(col) > widths[i] {
				widths[i] = len(col)
			}
		}
	}

	var sb strings.Builder

	// Header row
	sb.WriteString(renderRow(t.headers, widths))
	sb.WriteByte('\n')

	// Separator
	sb.WriteString(renderSeparator(widths))
	sb.WriteByte('\n')

	// Data rows
	for _, row := range t.rows {
		sb.WriteString(renderRow(row, widths))
		sb.WriteByte('\n')
	}

	return sb.String()
}

func renderRow(cols []string, widths []int) string {
	parts := make([]string, len(widths))
	for i := range widths {
		val := ""
		if i < len(cols) {
			val = cols[i]
		}
		parts[i] = fmt.Sprintf("%-*s", widths[i], val)
	}
	return strings.Join(parts, " | ")
}

func renderSeparator(widths []int) string {
	parts := make([]string, len(widths))
	for i, w := range widths {
		parts[i] = strings.Repeat("-", w)
	}
	return strings.Join(parts, "-|-")
}
