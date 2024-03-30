package handlers

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/klshriharsha/snippetbox/internal/assert"
	"github.com/klshriharsha/snippetbox/internal/testutils"
)

func TestSnippetViewHandler(t *testing.T) {
	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "An old silent pond...",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/snippet/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/snippet/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snippet/view/1.56",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/snippet/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}

	app := testutils.NewTestApplication(t)
	ts := testutils.NewTestServer(t, Routes(app))
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, _, body := ts.Get(t, tt.urlPath)
			assert.Equal(t, statusCode, tt.wantCode)
			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}

func TestSnippetCreateHandler(t *testing.T) {
	t.Run("Unauthenticated", func(t *testing.T) {
		app := testutils.NewTestApplication(t)
		ts := testutils.NewTestServer(t, Routes(app))
		defer ts.Close()

		statusCode, header, _ := ts.Get(t, "/snippet/create")
		assert.Equal(t, statusCode, http.StatusSeeOther)
		assert.Equal(t, header.Get("Location"), "/user/login")

	})

	t.Run("Authenticated", func(t *testing.T) {
		app := testutils.NewTestApplication(t)
		ts := testutils.NewTestServer(t, Routes(app))
		defer ts.Close()

		_, _, body := ts.Get(t, "/user/login")
		csrfToken := testutils.ExtractCSRFToken(t, body)

		form := url.Values{}
		form.Add("email", "snippetbox@example.com")
		form.Add("password", "password")
		form.Add("csrf_token", csrfToken)
		ts.Post(t, "/user/login", form)

		statusCode, _, body := ts.Get(t, "/snippet/create")
		assert.Equal(t, statusCode, http.StatusOK)
		assert.StringContains(t, body, `<form action='/snippet/create' method='POST'>`)
	})
}
