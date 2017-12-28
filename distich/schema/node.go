package schema

import (
	"fmt"
	"strconv"
	"strings"
)

type Node struct {
	Ref_  string `json:"$ref,omitempty"`
	Node_ string `json:"$schema,omitempty"`

	ID          string      `json:"-"`
	Type        string      `json:"type,omitempty"`
	Title       string      `json:"title,omitempty"`
	Description string      `json:"description,omitempty"`
	Property    string      `json:"_debug2,omitempty"`
	Default     interface{} `json:"default,omitempty"`

	Definitions map[string]*Node `json:"definitions,omitempty"`
	Items       *Node            `json:"items,omitempty"`
	Properties  map[string]*Node `json:"properties,omitempty"`

	Required []string `json:"required,omitempty"`

	AdditionalProperties *bool `json:"additionalProperties,omitempty"`

	Enum []string `json:"enum,omitempty"`

	OneOf []*Node `json:"oneOf,omitempty"`
	AnyOf []*Node `json:"anyOf,omitempty"`
	AllOf []*Node `json:"allOf,omitempty"`
}

func (n *Node) ApplyID(uri string) {
	n.ID = uri

	if n.Ref_ != "" && n.Ref_[0] == '#' {
		n.Ref_ = fmt.Sprintf("%s%s", strings.SplitN(n.ID, "#", 2)[0], n.Ref_)
	}

	for idx := range n.Definitions {
		x := n.Definitions[idx]
		x.ApplyID(fmt.Sprintf("%s/definitions/%s", n.ID, idx))
	}

	if n.Items != nil {
		n.Items.ApplyID(fmt.Sprintf("%s/items", n.ID))
	}

	for idx := range n.Properties {
		x := n.Properties[idx]
		x.ApplyID(fmt.Sprintf("%s/properties/%s", n.ID, idx))
	}

	for idx := range n.OneOf {
		n.OneOf[idx].ApplyID(fmt.Sprintf("%s/oneOf/%d", n.ID, idx))
	}

	for idx := range n.AnyOf {
		n.AnyOf[idx].ApplyID(fmt.Sprintf("%s/oneOf/%d", n.ID, idx))
	}

	for idx := range n.AllOf {
		n.AllOf[idx].ApplyID(fmt.Sprintf("%s/oneOf/%d", n.ID, idx))
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
