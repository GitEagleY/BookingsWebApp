package forms

// errors is a custom type that represents a map of form field names to error messages.
type errors map[string][]string

// Add adds an error message for a given form field.
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message) // Append the error message to the list of messages for the specified field.
}

// Get returns the first error message for a given form field.
func (e errors) Get(field string) string {
	es := e[field] // Retrieve the list of error messages for the specified field.
	if len(es) == 0 {
		return "" // Return an empty string if there are no error messages for the field.
	}
	return es[0] // Return the first error message for the field.
}
