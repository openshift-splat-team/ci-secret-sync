package schema

import "fmt"

type Generic struct {
	SchemaInterface

	Data []byte
}

// GetFields generic wrapper around the source data
// other schema types parse the data in to fields
func (g *Generic) GetFields() ([]string, error) {
	return []string{string(g.Data)}, nil
}

func (g *Generic) GetField(idx int) (string, error) {

	fields, _ := g.GetFields()
	if idx > 0 {
		return "", fmt.Errorf("field index is out of range")
	}
	return string(fields[idx]), nil
}
