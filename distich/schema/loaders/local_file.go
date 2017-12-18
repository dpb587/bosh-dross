package loaders

import (
	"io/ioutil"
	"strings"

	"github.com/dpb587/bosh-dross/distich/schema"
)

type LocalFile struct{}

var _ schema.Loader = LocalFile{}

func (LocalFile) IsSupported(uri string) bool {
	return strings.HasPrefix(uri, "file://")
}

func (LocalFile) Load(uri string) ([]byte, error) {
	uri = uri[7:]

	return ioutil.ReadFile(uri)
}
