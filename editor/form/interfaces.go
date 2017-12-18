package form

type FieldOptions map[string]interface{}

type FieldFactory interface {
	Create(uri, path string, options FieldOptions) (Field, error)
	IsSupported(fieldType string) bool
}

type Field interface {
	// Bind()
	HTML() []byte
	Set(interface{}) error
}
