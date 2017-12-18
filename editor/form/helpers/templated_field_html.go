package helpers

import (
	"bytes"
	"html/template"

	"github.com/dpb587/bosh-dross/editor/form"
)

func TemplatedFieldHTML(tmplStr string, field form.Field) ([]byte, error) {
	tmpl, err := template.New("field").Parse(tmplStr)
	if err != nil {
		panic(err)
	}

	var out bytes.Buffer

	err = tmpl.Execute(&out, struct{ Field form.Field }{Field: field})
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
