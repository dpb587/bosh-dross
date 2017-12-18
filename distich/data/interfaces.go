package data

type Node interface {
	GetPath() string
	GetRelativePath() string

	Visit(NodeVisitor) error

	Store(data Node) error
	Traverse(path string) (Node, error)
	Export() interface{}

	SetParent(parent Node)
}

type NodeVisitor interface {
	EnterNode(Node) error
	LeaveNode(Node) error
}
