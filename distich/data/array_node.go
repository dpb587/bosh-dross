package data

import (
	"fmt"
	"strconv"
)

type ArrayNode struct {
	node
	data []Node
}

var _ Node = &ArrayNode{}

func (n *ArrayNode) Store(data Node) error {
	index, err := strconv.Atoi(data.GetRelativePath())
	if err != nil {
		return err
	}

	if index >= len(n.data) {
		if index > len(n.data)+1 {
			// @todo no sparse support
			return fmt.Errorf("cannot store for more than len(n.data); set intermediate values first: %d vs %d", index, len(n.data))
		}

		data.SetParent(n)

		n.data = append(n.data, data)
	} else {
		n.data[index] = data
	}

	return nil
}

func (n *ArrayNode) Traverse(path string) (Node, error) {
	index, err := strconv.Atoi(path)
	if err != nil {
		return nil, err
	} else if index >= len(n.data) {
		return nil, PathNotFound
	}

	return n.data[index], nil
}

func (n *ArrayNode) Export() interface{} {
	r := []interface{}{}

	for _, v := range n.data {
		r = append(r, v.Export())
	}

	return r
}

func (n *ArrayNode) Visit(visitor NodeVisitor) error {
	err := visitor.EnterNode(n)
	if err != nil {
		return err
	}

	for _, v := range n.data {
		err = v.Visit(visitor)
		if err != nil {
			return err
		}
	}

	return visitor.LeaveNode(n)
}
