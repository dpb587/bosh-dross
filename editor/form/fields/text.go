package fields

import (
	"github.com/dpb587/bosh-dross/editor/form"
	"github.com/dpb587/bosh-dross/editor/form/helpers"
)

type Text struct {
	BaseField
	Value string
}

var _ form.Field = &Text{}

func (f *Text) HTML() []byte {
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
            <input type="text" name="{{ .Field.Name }}" id="{{ .Field.ID }}" value="{{ .Field.Value }}" class="">
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

func (f *Text) Set(data interface{}) error {
	f.Value = data.(string) // @todo panic

	return nil
}
