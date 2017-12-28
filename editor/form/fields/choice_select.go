package fields

import (
	"github.com/dpb587/bosh-dross/editor/form"
	"github.com/dpb587/bosh-dross/editor/form/helpers"
)

type ChoiceSelect struct {
	BaseField
	Value   []string
	Options map[interface{}]string
}

var _ form.Field = &ChoiceSelect{}

func (f *ChoiceSelect) HTML() []byte {
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

func (f *ChoiceSelect) Set(data interface{}) error {
	f.Value = data.([]string) // @todo panic

	return nil
}
