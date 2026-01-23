package display

import (
	"strings"
	"testing"
)

func TestTreeSingleRoomSingleDevice(t *testing.T) {
	tree := &Tree{
		Root: TreeNode{
			Label: "Home (1 room, 1 device, 0 unreachable)",
			Children: []TreeNode{
				{
					Label: "Office",
					Children: []TreeNode{
						{Label: "Desk Lamp         light       on    80%"},
					},
				},
			},
		},
	}

	result := tree.Render()

	expected := strings.Join([]string{
		"Home (1 room, 1 device, 0 unreachable)",
		"└── Office",
		"    └── Desk Lamp         light       on    80%",
		"",
	}, "\n")

	if result != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestTreeMultipleRoomsMultipleDevices(t *testing.T) {
	tree := &Tree{
		Root: TreeNode{
			Label: "Home (3 rooms, 5 devices, 1 unreachable)",
			Children: []TreeNode{
				{
					Label: "Office",
					Children: []TreeNode{
						{Label: "Desk Lamp         light       on    80%"},
						{Label: "AC Unit           thermostat  on    22.5°"},
					},
				},
				{
					Label: "Bedroom",
					Children: []TreeNode{
						{Label: "Ceiling Light     light       off"},
						{Label: "Fan               fan         on    speed 85%"},
					},
				},
				{
					Label: "Living Room",
					Children: []TreeNode{
						{Label: "Floor Lamp        light       unreachable"},
					},
				},
			},
		},
	}

	result := tree.Render()

	if !strings.Contains(result, "├── Office") {
		t.Error("expected ├── Office")
	}
	if !strings.Contains(result, "├── Bedroom") {
		t.Error("expected ├── Bedroom")
	}
	if !strings.Contains(result, "└── Living Room") {
		t.Error("expected └── Living Room (last child)")
	}
	if !strings.Contains(result, "│   ├── Desk Lamp") {
		t.Error("expected │   ├── for non-last child of non-last parent")
	}
	if !strings.Contains(result, "│   └── AC Unit") {
		t.Error("expected │   └── for last child of non-last parent")
	}
	if !strings.Contains(result, "    └── Floor Lamp") {
		t.Error("expected     └── for last child of last parent")
	}
}

func TestTreeEmptyChildren(t *testing.T) {
	tree := &Tree{
		Root: TreeNode{
			Label:    "Home (0 rooms, 0 devices, 0 unreachable)",
			Children: nil,
		},
	}

	result := tree.Render()
	expected := "Home (0 rooms, 0 devices, 0 unreachable)\n"

	if result != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestTreeRoomWithNoDevices(t *testing.T) {
	tree := &Tree{
		Root: TreeNode{
			Label: "Home (1 room, 0 devices, 0 unreachable)",
			Children: []TreeNode{
				{Label: "Empty Room", Children: nil},
			},
		},
	}

	result := tree.Render()

	expected := strings.Join([]string{
		"Home (1 room, 0 devices, 0 unreachable)",
		"└── Empty Room",
		"",
	}, "\n")

	if result != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestTreeConnectorsContinuation(t *testing.T) {
	tree := &Tree{
		Root: TreeNode{
			Label: "Root",
			Children: []TreeNode{
				{
					Label: "A",
					Children: []TreeNode{
						{Label: "A1"},
						{Label: "A2"},
					},
				},
				{
					Label: "B",
					Children: []TreeNode{
						{Label: "B1"},
					},
				},
			},
		},
	}

	result := tree.Render()

	lines := strings.Split(result, "\n")
	// Root
	if lines[0] != "Root" {
		t.Errorf("line 0: expected Root, got %q", lines[0])
	}
	// A (not last) uses ├──
	if lines[1] != "├── A" {
		t.Errorf("line 1: expected '├── A', got %q", lines[1])
	}
	// A's children use │ prefix since A is not last
	if lines[2] != "│   ├── A1" {
		t.Errorf("line 2: expected '│   ├── A1', got %q", lines[2])
	}
	if lines[3] != "│   └── A2" {
		t.Errorf("line 3: expected '│   └── A2', got %q", lines[3])
	}
	// B (last) uses └──
	if lines[4] != "└── B" {
		t.Errorf("line 4: expected '└── B', got %q", lines[4])
	}
	// B's children use space prefix since B is last
	if lines[5] != "    └── B1" {
		t.Errorf("line 5: expected '    └── B1', got %q", lines[5])
	}
}
