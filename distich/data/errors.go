package data

import "errors"

var PathNotFound = errors.New("Path does not exist")
var PathTraversalNotSupported = errors.New("Path traversal is not supported")
