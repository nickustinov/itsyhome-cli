package display

import "strings"

type TreeNode struct {
	Label    string
	Children []TreeNode
}

type Tree struct {
	Root TreeNode
}

func (t *Tree) Render() string {
	var sb strings.Builder
	sb.WriteString(t.Root.Label)
	sb.WriteByte('\n')
	renderChildren(&sb, t.Root.Children, "")
	return sb.String()
}

func renderChildren(sb *strings.Builder, children []TreeNode, prefix string) {
	for i, child := range children {
		last := i == len(children)-1

		connector := "├── "
		if last {
			connector = "└── "
		}

		sb.WriteString(prefix)
		sb.WriteString(connector)
		sb.WriteString(child.Label)
		sb.WriteByte('\n')

		if len(child.Children) > 0 {
			childPrefix := prefix + "│   "
			if last {
				childPrefix = prefix + "    "
			}
			renderChildren(sb, child.Children, childPrefix)
		}
	}
}
