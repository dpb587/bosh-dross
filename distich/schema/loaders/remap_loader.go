package loaders

import (
	"fmt"
	"strings"

	"github.com/dpb587/bosh-dross/distich/schema"
)

type RemappedLoader struct {
	loader   schema.Loader
	old, new string
}

var _ schema.Loader = RemappedLoader{}

func NewRemappedLoader(loader schema.Loader, old, new string) RemappedLoader {
	return RemappedLoader{
		loader: loader,
		old:    old,
		new:    new,
	}
}

func (l RemappedLoader) IsSupported(uri string) bool {
	return l.loader.IsSupported(l.remap(uri))
}

func (l RemappedLoader) Load(uri string) ([]byte, error) {
	return l.loader.Load(l.remap(uri))
}

func (l RemappedLoader) remap(uri string) string {
	if strings.HasPrefix(uri, l.old) {
		return fmt.Sprintf("%s%s", l.new, uri[len(l.old):])
	}

	return uri
}
