package schema

type Loader interface {
	Load(uri string) ([]byte, error)
	IsSupported(uri string) bool
}
