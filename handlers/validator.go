package handlers

import "github.com/xeipuuv/gojsonschema"

// Validate validates byte data against json schema.
func Validate(schema string, data []byte) []string {
	want := gojsonschema.NewStringLoader(schema)
	got := gojsonschema.NewBytesLoader(data)

	results, err := gojsonschema.Validate(want, got)

	if err != nil {
		return []string{err.Error()}
	}

	errors := []string{}
	for _, e := range results.Errors() {
		errors = append(errors, e.String())
	}

	return errors
}
