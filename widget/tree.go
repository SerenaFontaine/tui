package widget

import "tui"

// TreeNode represents a node in a tree view.
type TreeNode struct {
	Text     string
	Children []*TreeNode
	Data     any
	Expanded bool
}

// NewTreeNode creates a tree node.
func NewTreeNode(text string, children ...*TreeNode) *TreeNode {
	return &TreeNode{
		Text:     text,
		Children: children,
		Expanded: true,
	}
}

// Tree is a navigable tree view widget.
type Tree struct {
	Root          []*TreeNode
	Selected      int // index in the flattened visible list
	Style         tui.Style
	SelectedStyle tui.Style
	Block         *tui.Block
	IndentSize    int

	flat   []flatNode // cached flattened view
	offset int
}

type flatNode struct {
	node  *TreeNode
	depth int
}

// NewTree creates a tree from root nodes.
func NewTree(roots ...*TreeNode) *Tree {
	t := &Tree{
		Root:          roots,
		SelectedStyle: tui.NewStyle().Reverse(true),
		IndentSize:    2,
	}
	t.rebuild()
	return t
}

// SetBlock adds a border block.
func (t *Tree) SetBlock(b tui.Block) *Tree { t.Block = &b; return t }

// SetStyle sets the default node style.
func (t *Tree) SetStyle(s tui.Style) *Tree { t.Style = s; return t }

// SetSelectedStyle sets the selected node style.
func (t *Tree) SetSelectedStyle(s tui.Style) *Tree { t.SelectedStyle = s; return t }

// SelectedNode returns the currently selected tree node.
func (t *Tree) SelectedNode() *TreeNode {
	if t.Selected >= 0 && t.Selected < len(t.flat) {
		return t.flat[t.Selected].node
	}
	return nil
}

// Update handles navigation and expand/collapse.
func (t *Tree) Update(msg tui.Msg) (*Tree, tui.Cmd) {
	switch msg := msg.(type) {
	case tui.KeyMsg:
		switch msg.Type {
		case tui.KeyUp, tui.KeyCtrlP:
			if t.Selected > 0 {
				t.Selected--
			}
		case tui.KeyDown, tui.KeyCtrlN:
			if t.Selected < len(t.flat)-1 {
				t.Selected++
			}
		case tui.KeyRight:
			// Expand
			if node := t.SelectedNode(); node != nil && len(node.Children) > 0 {
				node.Expanded = true
				t.rebuild()
			}
		case tui.KeyLeft:
			// Collapse
			if node := t.SelectedNode(); node != nil && node.Expanded && len(node.Children) > 0 {
				node.Expanded = false
				t.rebuild()
			}
		case tui.KeyEnter, tui.KeySpace:
			// Toggle
			if node := t.SelectedNode(); node != nil && len(node.Children) > 0 {
				node.Expanded = !node.Expanded
				t.rebuild()
			}
		case tui.KeyRune:
			switch msg.Rune {
			case 'j':
				if t.Selected < len(t.flat)-1 {
					t.Selected++
				}
			case 'k':
				if t.Selected > 0 {
					t.Selected--
				}
			case 'l':
				if node := t.SelectedNode(); node != nil && len(node.Children) > 0 {
					node.Expanded = true
					t.rebuild()
				}
			case 'h':
				if node := t.SelectedNode(); node != nil && node.Expanded && len(node.Children) > 0 {
					node.Expanded = false
					t.rebuild()
				}
			}
		}
	}
	return t, nil
}

// Render draws the tree.
func (t *Tree) Render(buf *tui.Buffer, area tui.Rect) {
	if area.IsEmpty() {
		return
	}

	inner := area
	if t.Block != nil {
		inner = t.Block.Render(buf, area)
	}

	if inner.IsEmpty() || len(t.flat) == 0 {
		return
	}

	// Adjust scroll
	visibleHeight := inner.Height
	if t.Selected < t.offset {
		t.offset = t.Selected
	}
	if t.Selected >= t.offset+visibleHeight {
		t.offset = t.Selected - visibleHeight + 1
	}

	for i := 0; i < visibleHeight; i++ {
		idx := t.offset + i
		if idx >= len(t.flat) {
			break
		}

		fn := t.flat[idx]
		style := t.Style
		if idx == t.Selected {
			style = t.SelectedStyle
			for x := inner.X; x < inner.Right(); x++ {
				buf.SetChar(x, inner.Y+i, ' ', style)
			}
		}

		// Indentation
		indent := fn.depth * t.IndentSize
		x := inner.X + indent

		// Tree prefix
		prefix := ""
		if len(fn.node.Children) > 0 {
			if fn.node.Expanded {
				prefix = "▼ "
			} else {
				prefix = "▶ "
			}
		} else {
			prefix = "  "
		}

		n := buf.SetString(x, inner.Y+i, prefix, style)
		text := fn.node.Text
		maxLen := inner.Width - indent - n
		if maxLen > 0 {
			if len(text) > maxLen {
				text = text[:maxLen]
			}
			buf.SetString(x+n, inner.Y+i, text, style)
		}
	}
}

func (t *Tree) rebuild() {
	t.flat = t.flat[:0]
	for _, root := range t.Root {
		t.flatten(root, 0)
	}
	if t.Selected >= len(t.flat) {
		t.Selected = len(t.flat) - 1
	}
	if t.Selected < 0 {
		t.Selected = 0
	}
}

func (t *Tree) flatten(node *TreeNode, depth int) {
	t.flat = append(t.flat, flatNode{node: node, depth: depth})
	if node.Expanded {
		for _, child := range node.Children {
			t.flatten(child, depth+1)
		}
	}
}
