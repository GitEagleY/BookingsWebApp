package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates a custom form struct that embeds a url.Values object and holds validation errors.
type Form struct {
	url.Values
	Errors errors // Custom errors map to store validation errors.
}

// New initializes a new Form struct with provided data and an empty errors map.
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required checks if specified form fields are not empty.
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Has checks if a form field exists in the post data and is not empty.
func (f *Form) Has(field string) bool {
	x := f.Get(field)
	if x == "" {
		f.Errors.Add(field, "This field cannot be blank")
		return false
	}
	return true
}

// Valid checks if the form has no validation errors.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// MinLength checks if a string form field meets the minimum length requirement.
func (f *Form) MinLength(field string, length int) bool {
	x := f.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

// IsEmail checks if a form field contains a valid email address using the govalidator package.
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}

//my email validator
/*
func (f *Form) IsValidEmailFormat(field string, r *http.Request) bool {
	x := r.Form.Get(field)

	if strings.Contains("@", x) == false {
		return false
	}

	splitted := strings.Split(x, "@")
	if len(splitted) != 2 {
		return false
	}

	for _, v := range splitted {
		if len(v) < 3 {
			return false
		}
	}
	domainSuffixes := []string{".com", ".org", ".me", ".edu"}
	for _, suffix := range domainSuffixes {
		if !strings.HasSuffix(splitted[len(splitted)-1], suffix) {
			return false
		}
	}
	return true

}
*/
