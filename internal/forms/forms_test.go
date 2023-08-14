package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/something", nil)
	form := New(r.PostForm)
	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should get valid")
	}
}
func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/something", nil)
	form := New(r.PostForm)
	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("shows valid when requred fileds missing")
	}
	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")
	r = httptest.NewRequest("POST", "/something", nil)
	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows that doesnt have required filed but it does")
	}
}

func TestForm_Has(t *testing.T) {

	postedData := url.Values{}
	form := New(url.Values{})
	form.Has("qwerty")
	if form.Has("qwerty") {
		t.Error("has field when shouldn't")
	}

	postedData = url.Values{}
	postedData.Add("a", "a")

	form = New(postedData)

	has := form.Has("a")
	if !has {
		t.Error("shows form does not have field when it should")
	}
}

func TestForm_MinLenghth(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.MinLength("a", 10)
	if form.Valid() {
		t.Error("form shows min length for non-existent field")
	}

	isError := form.Errors.Get("a")
	if isError == "" {
		t.Error("should have error but did not get one")
	}

	postedData = url.Values{}
	postedData.Add("afield", "avalue")

	form = New(postedData)

	form.MinLength("afield", 100)
	if form.Valid() {
		t.Error("shows min length of 100 met when data is shorter")
	}

	postedData = url.Values{}
	postedData.Add("b_field", "abc123")
	form = New(postedData)

	form.MinLength("b_field", 1)
	if !form.Valid() {
		t.Error("shows min length if 1 is not met when it is")
	}

	isError = form.Errors.Get("b_field")
	if isError != "" {
		t.Error("should not have error but got one")
	}

}

func TestForm_IsEmail(t *testing.T) {
	r := httptest.NewRequest("POST", "/something", nil)

	postedData := url.Values{}
	postedData.Add("q", "w")
	form := New(r.PostForm)
	form.IsEmail("q")
	if form.Valid() {
		t.Error("hasnt error when should")
	}

	postedData.Add("qw", "w@a.com")

	r = httptest.NewRequest("POST", "/something", nil)
	r.PostForm = postedData
	form = New(r.PostForm)
	form.IsEmail("qw")
	if !form.Valid() {
		t.Error("have error when shouldt")
	}
}
