package distich

import (
	"strings"

	"github.com/dpb587/bosh-dross/distich/data"
	"github.com/dpb587/bosh-dross/distich/schema"
)

func Traverse(d data.Node, s schema.Node, path string) (data.Node, schema.Node, error) {
	pathPieces := strings.SplitN(strings.TrimPrefix(path, "/"), "/", 2)

	if pathPieces[0] == "" {
		if len(pathPieces) == 1 {
			return d, s, nil
		}

		return Traverse(d, s, pathPieces[1])
	}

	dc, err := d.Traverse(pathPieces[0])
	if err != nil {
		return nil, schema.Node{}, err
	}

	sc, err := s.Traverse(pathPieces[0])
	if err != nil {
		return nil, schema.Node{}, err
	}

	if len(pathPieces) == 2 {
		return Traverse(dc, sc, pathPieces[1])
	}

	return dc, sc, nil
}
