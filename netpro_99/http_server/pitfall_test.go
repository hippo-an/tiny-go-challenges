package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHeandlerWriteHeader(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Bad request"))
		w.WriteHeader(http.StatusBadRequest)
	}

	r := httptest.NewRequest(http.MethodGet, "http://test", nil)
	w := httptest.NewRecorder()

	handler(w, r)

	t.Logf("Response status: %q", w.Result().Status)
}
