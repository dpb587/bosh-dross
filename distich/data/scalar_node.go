package data

type ScalarNode struct {
	node
	data interface{}
}

var _ Node = &ScalarNode{}

func (n *ScalarNode) Store(data Node) error {
	return PathTraversalNotSupported
}

func (n *ScalarNode) Traverse(path string) (Node, error) {
	return nil, PathTraversalNotSupported
}

func (n *ScalarNode) Import(data interface{}) {
	n.data = data
}

func (n *ScalarNode) Export() interface{} {
	return n.data
}

func (n *ScalarNode) Visit(visitor NodeVisitor) error {
	err := visitor.EnterNode(n)
	if err != nil {
		return err
	}

	return visitor.LeaveNode(n)
}
