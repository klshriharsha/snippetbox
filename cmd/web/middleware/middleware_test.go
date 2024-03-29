package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/klshriharsha/snippetbox/internal/assert"
)

func TestSecureHeaders(t *testing.T) {
	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	SecureHeaders(next).ServeHTTP(rr, r)

	rs := rr.Result()
	expected := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, rs.Header.Get("Content-Security-Policy"), expected)
	assert.Equal(t, rs.Header.Get("Referrer-Policy"), "origin-when-cross-origin")
	assert.Equal(t, rs.Header.Get("X-Content-Type-Options"), "nosniff")
	assert.Equal(t, rs.Header.Get("X-Frame-Options"), "deny")
	assert.Equal(t, rs.Header.Get("X-XSS-Protection"), "0")
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(body), "OK")
}
