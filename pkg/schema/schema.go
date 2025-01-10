package schema

// Source secrets may have various formats of data in a specific field
// that may need to be parsed in to fields. An example of this is getting
// the pull credentials from a docker config file.
type SchemaInterface interface {
	GetFields() ([]string, error)
	GetField(idx int) (string, error)
}
