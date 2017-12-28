package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Resolver struct {
	loader Loader
	cache  map[string]*Node
}

func NewResolver(loader Loader) Resolver {
	return Resolver{
		loader: loader,
		cache:  map[string]*Node{},
	}
}

func (r *Resolver) Load(uri string) (*Node, error) {
	uriSplit := strings.SplitN(uri, "#", 2)

	if _, found := r.cache[uriSplit[0]]; !found {
		sc, err := r.parse(uriSplit[0])
		if err != nil {
			return nil, err
		}

		r.cache[uriSplit[0]] = sc
	}

	schema := r.cache[uriSplit[0]]

	if len(uriSplit) == 2 && uriSplit[1] != "" {
		schema2, err := r.traverse(schema, strings.TrimPrefix(uriSplit[1], "/"))
		if err != nil {
			return nil, err
		}

		schema = schema2
	}

	return schema, nil
}

func (r *Resolver) parse(uri string) (*Node, error) {
	sbytes, err := r.loader.Load(uri)
	if err != nil {
		return nil, err
	}

	var s Node

	err = json.Unmarshal(sbytes, &s)
	if err != nil {
		return nil, err
	}

	s.ApplyID(fmt.Sprintf("%s#", uri))

	// @todo hack to support self-referencing $ref
	r.cache[uri] = &s

	err = r.recurseParse(&s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// @todo naive; can panic
func (r Resolver) traverse(node *Node, fragment string) (*Node, error) {
	if fragment == "" {
		return node, nil
	}

	fragmentSplit := strings.SplitN(fragment, "/", 2)

	switch fragmentSplit[0] {
	case "definitions":
		fragmentSplit = strings.SplitN(fragmentSplit[1], "/", 2)
		node = node.Definitions[fragmentSplit[0]]
	case "items":
		node = node.Items
	case "properties":
		fragmentSplit = strings.SplitN(fragmentSplit[1], "/", 2)
		node = node.Properties[fragmentSplit[0]]
	case "oneOf":
		fragmentSplit = strings.SplitN(fragmentSplit[1], "/", 2)
		idx, _ := strconv.Atoi(fragmentSplit[0])
		node = node.OneOf[idx]
	case "anyOf":
		fragmentSplit = strings.SplitN(fragmentSplit[1], "/", 2)
		idx, _ := strconv.Atoi(fragmentSplit[0])
		node = node.AnyOf[idx]
	case "allOf":
		fragmentSplit = strings.SplitN(fragmentSplit[1], "/", 2)
		idx, _ := strconv.Atoi(fragmentSplit[0])
		node = node.AllOf[idx]
	default:
		return nil, errors.New("invalid fragment")
	}

	if len(fragmentSplit) == 2 {
		return r.traverse(node, fragmentSplit[1])
	}

	return node, nil
}

// @todo self-references will recurse
// @todo merge $ref'd node properties?
func (r Resolver) recurseParse(node *Node) error {
	var err error

	for subnodeIdx, subnode := range node.Definitions {
		if subnode.Ref_ != "" {
			node.Definitions[subnodeIdx], err = r.Load(subnode.Ref_)
			if err != nil {
				return err
			}
		} else {
			err = r.recurseParse(node.Definitions[subnodeIdx])
			if err != nil {
				return err
			}
		}
	}

	if node.Items != nil {
		if node.Items.Ref_ != "" {
			node.Items, err = r.Load(node.Ref_)
			if err != nil {
				return err
			}
		} else {
			err = r.recurseParse(node.Items)
			if err != nil {
				return err
			}
		}
	}

	for subnodeIdx, subnode := range node.Properties {
		if subnode.Ref_ != "" {
			node.Properties[subnodeIdx], err = r.Load(subnode.Ref_)
			if err != nil {
				return err
			}
		} else {
			err = r.recurseParse(node.Properties[subnodeIdx])
			if err != nil {
				return err
			}
		}
	}

	for subnodeIdx, subnode := range node.OneOf {
		if subnode.Ref_ != "" {
			node.OneOf[subnodeIdx], err = r.Load(subnode.Ref_)
			if err != nil {
				return err
			}
		} else {
			err = r.recurseParse(node.OneOf[subnodeIdx])
			if err != nil {
				return err
			}
		}
	}

	for subnodeIdx, subnode := range node.AnyOf {
		if subnode.Ref_ != "" {
			node.AnyOf[subnodeIdx], err = r.Load(subnode.Ref_)
			if err != nil {
				return err
			}
		} else {
			err = r.recurseParse(node.AnyOf[subnodeIdx])
			if err != nil {
				return err
			}
		}
	}

	for subnodeIdx, subnode := range node.AllOf {
		if subnode.Ref_ != "" {
			node.AllOf[subnodeIdx], err = r.Load(subnode.Ref_)
			if err != nil {
				return err
			}
		} else {
			err = r.recurseParse(node.AllOf[subnodeIdx])
			if err != nil {
				return err
			}
		}
	}

	return nil
}
