package data

import (
	"strconv"
)

func CreateNode(d interface{}) (Node, error) {
	return createNode("", d)
}

func createNode(path string, d interface{}) (Node, error) {
	var localnode Node

	switch d.(type) {
	case map[string]interface{}:
		localnode = &HashNode{node: node{path: path}}

		for k, v := range d.(map[string]interface{}) {
			vn, err := createNode(k, v)
			if err != nil {
				return nil, err
			}

			err = localnode.Store(vn)
			if err != nil {
				return nil, err
			}
		}
	case map[interface{}]interface{}:
		localnode = &HashNode{node: node{path: path}}

		for k, v := range d.(map[interface{}]interface{}) {
			vn, err := createNode(k.(string), v) // @todo panic cast
			if err != nil {
				return nil, err
			}

			err = localnode.Store(vn)
			if err != nil {
				return nil, err
			}
		}
	case []interface{}:
		localnode = &ArrayNode{node: node{path: path}}

		for k, v := range d.([]interface{}) {
			vn, err := createNode(strconv.Itoa(k), v)
			if err != nil {
				return nil, err
			}

			err = localnode.Store(vn)
			if err != nil {
				return nil, err
			}
		}
	default:
		localnode = &ScalarNode{node: node{path: path}}

		localnode.(*ScalarNode).Import(d)
	}

	return localnode, nil
}
