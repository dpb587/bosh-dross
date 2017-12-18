package schema

import (
	"fmt"
	"strconv"
)

type Node struct {
	Ref_  string `json:"$ref"`
	Node_ string `json:"$schema"`

	ID          string `json:"-"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Property    string `json:"-"`

	Definitions map[string]*Node `json:"definitions"`
	Items       *Node            `json:"items"`
	Properties  map[string]*Node `json:"properties"`

	Required []string `json:"required"`

	AdditionalProperties *bool `json:"additionalProperties"`

	Enum []string `json:"enum"`

	OneOf []*Node `json:"oneOf"`
	AnyOf []*Node `json:"anyOf"`
	AllOf []*Node `json:"allOf"`
}

func (n *Node) ApplyID(uri string) {
	n.ID = uri

	for idx := range n.Definitions {
		x := n.Definitions[idx]
		x.ApplyID(fmt.Sprintf("%s/definitions/%s", uri, idx))
	}

	if n.Items != nil {
		n.Items.ApplyID(fmt.Sprintf("%s/items", uri))
	}

	for idx := range n.Properties {
		x := n.Properties[idx]
		x.ApplyID(fmt.Sprintf("%s/properties/%s", uri, idx))
	}

	for idx := range n.OneOf {
		n.OneOf[idx].ApplyID(fmt.Sprintf("%s/oneOf/%d", uri, idx))
	}

	for idx := range n.AnyOf {
		n.AnyOf[idx].ApplyID(fmt.Sprintf("%s/oneOf/%d", uri, idx))
	}

	for idx := range n.AllOf {
		n.AllOf[idx].ApplyID(fmt.Sprintf("%s/oneOf/%d", uri, idx))
	}
}

func (n *Node) Traverse(path string) (Node, error) {
	if path == "-" {
		// @todo not responsibility of this; should mean new node
		return *n, nil
	} else if _, err := strconv.Atoi(path); err == nil {
		if n.Type != "array" {
			return Node{}, fmt.Errorf(
				`Expected array at "%s" of schema "%s", but found "%s"`,
				path,
				n.Ref_,
				n.Type,
			)
		}

		return *n.Items, nil
	} else if n.Type != "object" {
		return Node{}, fmt.Errorf(
			`Expected object at "%s" of schema "%s", but found "%s"`,
			path,
			n.Ref_,
			n.Type,
		)
	}

	property, found := n.Properties[path]
	if !found {
		return Node{}, PathNotFound
	}

	return *property, nil
}
