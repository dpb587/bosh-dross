package fields

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/dpb587/bosh-dross/editor/form"
)

type BaseField struct {
	id string

	Path    string
	Title   string
	Name    string
	Help    string
	Options form.FieldOptions
}

func (f BaseField) ID() string {
	if f.id == "" {
		hasher := md5.New()
		hasher.Write([]byte(f.Path))

		f.id = hex.EncodeToString(hasher.Sum(nil))
	}

	return f.id
}
