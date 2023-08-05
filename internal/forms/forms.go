package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Requied(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// checks if form field is in post and not empty
func (f *Form) Has(field string) bool {
	x := f.Get(field)
	if x == "" {
		f.Errors.Add(field, "This field cannot be blank")
		return false
	}
	return true
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// checks for string minimum length
func (f *Form) MinLength(field string, length int) bool {
	x := f.Get(field)
	if len((x)) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "invalid email address")
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
