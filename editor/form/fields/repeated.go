package fields

import (
	"fmt"

	"github.com/dpb587/bosh-dross/editor/form"
	"github.com/dpb587/bosh-dross/editor/form/helpers"
)

type Repeated struct {
	BaseField
	Value []interface{}

	FieldFactory form.FieldFactory
	FieldType    string
	Fields       []form.Field
}

var _ form.Field = &Repeated{}

func (f *Repeated) HTML() []byte {
	html, err := helpers.TemplatedFieldHTML(`
    <div class="grid form-row">
      <div class="col form-col">
        <div class="form-unit">
          <div class="grid grid-nogutter label-row">
            <div class="col">
              <label for="{{ .Field.ID }}">{{ .Field.Title }}</label>
            </div>
            <div class="col col-fixed col-middle post-label"></div>
          </div>
          <div class="field-row">
            <select name="{{ .Field.Name }}" id="{{ .Field.ID }}" class="">
              {{ range k, v := .Field.Options }}
                <option {{ range .Field.Value }}{{ if k == .Value }} selected="selected"{{ end }} value="{{ k }}">{{ v }}</option>
              {{ end }}
            </select>
          </div>
          <div class="help-row type-dark-5">{{ .Field.Help }}</div>
        </div>
      </div>
    </div>
`, f)
	if err != nil {
		panic(err)
	}

	return html
}

func (f *Repeated) Set(data interface{}) error {
	f.Value = data.([]interface{})

	for valueIdx, value := range f.Value {
		field, err := f.FieldFactory.Create(f.FieldType, fmt.Sprintf("%s/%d", f.BaseField.Path, valueIdx), f.BaseField.Options)
		if err != nil {
			return err
		}

		field.Set(value)

		f.Fields = append(f.Fields, field)
	}

	return nil
}
