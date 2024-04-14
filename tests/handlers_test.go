package tests

import (
	"banners/internal/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProcessError(t *testing.T) {
	w := httptest.NewRecorder()
	err := handlers.ErrIncorrectBannerContent
	code := http.StatusBadRequest

	handlers.ProcessError(w, err, code)
	got := w.Body.String()

	want := "{\"error\":\"incorrect banner content\"}"
	if got != want {
		t.Errorf("incorrect response body: got %q, want %q", got, want)
	}
	if w.Result().StatusCode != code {
		t.Errorf("incorrect response status: got %d, want %d", w.Result().StatusCode, code)
	}
}
