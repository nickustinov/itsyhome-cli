package display

import (
	"strings"
	"testing"
)

func TestTableRender(t *testing.T) {
	tbl := NewTable("Name", "Type", "Room")
	tbl.AddRow("Lamp", "light", "Office")
	tbl.AddRow("AC Unit", "thermostat", "Bedroom")

	out := tbl.Render()
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")

	if len(lines) != 4 {
		t.Fatalf("expected 4 lines, got %d: %q", len(lines), out)
	}

	// Header
	if !strings.Contains(lines[0], "Name") || !strings.Contains(lines[0], "Type") {
		t.Errorf("header missing columns: %s", lines[0])
	}

	// Separator
	if !strings.Contains(lines[1], "-|-") {
		t.Errorf("separator malformed: %s", lines[1])
	}

	// Data alignment
	if !strings.Contains(lines[2], "Lamp") {
		t.Errorf("row 1 missing data: %s", lines[2])
	}
	if !strings.Contains(lines[3], "AC Unit") {
		t.Errorf("row 2 missing data: %s", lines[3])
	}
}

func TestTableColumnWidths(t *testing.T) {
	tbl := NewTable("A", "B")
	tbl.AddRow("Long value here", "x")

	out := tbl.Render()
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")

	// Separator should have dashes matching longest content
	if !strings.Contains(lines[1], "---------------") {
		t.Errorf("separator too short for content: %s", lines[1])
	}
}

func TestTableEmpty(t *testing.T) {
	tbl := NewTable("X", "Y")
	out := tbl.Render()
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")

	if len(lines) != 2 {
		t.Fatalf("expected 2 lines for empty table, got %d", len(lines))
	}
}

func TestTableNoHeaders(t *testing.T) {
	tbl := &Table{}
	if tbl.Render() != "" {
		t.Error("expected empty string for no headers")
	}
}
