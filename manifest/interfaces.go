package manifest

import (
	"github.com/dpb587/bosh-dross/distich/schema"
)

type SchemaGuesser interface {
	Guess(Manifest) (*schema.Node, error)
}
