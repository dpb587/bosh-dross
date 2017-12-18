package data

import "fmt"

type node struct {
	parent Node
	path   string
}

func (n *node) GetRelativePath() string {
	return n.path
}

func (n *node) GetPath() string {
	prefix := ""

	if n.parent != nil {
		prefix = fmt.Sprintf("%s/", n.parent.GetPath())
	}

	return fmt.Sprintf("%s%s", prefix, n.GetRelativePath())
}

func (n *node) SetParent(parent Node) {
	n.parent = parent
}
