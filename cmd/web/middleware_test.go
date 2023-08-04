package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNosurf(t *testing.T) {
	var myH myHandler
	h := NoSurf(&myH)

	switch v := h.(type) {
	case http.Handler:
		//ok
	default:
		t.Error(fmt.Sprintf("type is not http.Handler %T", v))
	}
}
