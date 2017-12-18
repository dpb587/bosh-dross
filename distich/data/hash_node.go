package data

type HashNode struct {
	node
	data map[string]Node
}

var _ Node = &HashNode{}

func (n *HashNode) Store(data Node) error {
	if n.data == nil {
		n.data = map[string]Node{}
	}

	data.SetParent(n)

	n.data[data.GetRelativePath()] = data

	return nil
}

func (n *HashNode) Traverse(path string) (Node, error) {
	data, found := n.data[path]
	if !found {
		return nil, PathNotFound
	}

	return data, nil
}

func (n *HashNode) Export() interface{} {
	r := map[string]interface{}{}

	for k, v := range n.data {
		r[k] = v.Export()
	}

	return r
}

func (n *HashNode) Visit(visitor NodeVisitor) error {
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
