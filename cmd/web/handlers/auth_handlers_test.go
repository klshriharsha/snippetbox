package handlers

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/klshriharsha/snippetbox/internal/assert"
	"github.com/klshriharsha/snippetbox/internal/testutils"
)

func TestSignupHandler(t *testing.T) {
	app := testutils.NewTestApplication(t)
	ts := testutils.NewTestServer(t, Routes(app))
	defer ts.Close()

	_, _, body := ts.Get(t, "/user/signup")
	csrfToken := testutils.ExtractCSRFToken(t, body)

	const (
		validName     = "Snippetbox"
		validPassword = "password"
		validEmail    = "snippetbox@example.com"
		formTag       = "<form action='/user/signup' method='POST' novalidate>"
	)

	tests := []struct {
		name        string
		username    string
		email       string
		password    string
		csrfToken   string
		wantCode    int
		wantFormTag string
	}{
		{
			name:      "Valid submission",
			username:  validName,
			email:     validEmail,
			password:  validPassword,
			csrfToken: csrfToken,
			wantCode:  http.StatusSeeOther,
		},
		{
			name:      "Invalid CSRF token",
			username:  validName,
			email:     validEmail,
			password:  validPassword,
			csrfToken: "wrongtoken",
			wantCode:  http.StatusBadRequest,
		},
		{
			name:        "Empty name",
			username:    "",
			email:       validEmail,
			password:    validPassword,
			csrfToken:   csrfToken,
			wantCode:    http.StatusUnprocessableEntity,
			wantFormTag: formTag,
		},
		{
			name:        "Empty email",
			username:    validName,
			email:       "",
			password:    validPassword,
			csrfToken:   csrfToken,
			wantCode:    http.StatusUnprocessableEntity,
			wantFormTag: formTag,
		},
		{
			name:        "Empty password",
			username:    validName,
			email:       validEmail,
			password:    "",
			csrfToken:   csrfToken,
			wantCode:    http.StatusUnprocessableEntity,
			wantFormTag: formTag,
		},
		{
			name:        "Invalid email",
			username:    validName,
			email:       "snippetbox@example.",
			password:    validPassword,
			csrfToken:   csrfToken,
			wantCode:    http.StatusUnprocessableEntity,
			wantFormTag: formTag,
		},
		{
			name:        "Short password",
			username:    validName,
			email:       validEmail,
			password:    "pass",
			csrfToken:   csrfToken,
			wantCode:    http.StatusUnprocessableEntity,
			wantFormTag: formTag,
		},
		{
			name:        "Duplicate email",
			username:    validName,
			email:       "duplicate@example.com",
			password:    "pass",
			csrfToken:   csrfToken,
			wantCode:    http.StatusUnprocessableEntity,
			wantFormTag: formTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.username)
			form.Add("email", tt.email)
			form.Add("password", tt.password)
			form.Add("csrf_token", tt.csrfToken)

			statusCode, _, body := ts.Post(t, "/user/signup", form)
			assert.Equal(t, statusCode, tt.wantCode)
			if tt.wantFormTag != "" {
				assert.StringContains(t, body, tt.wantFormTag)
			}
		})
	}
}
