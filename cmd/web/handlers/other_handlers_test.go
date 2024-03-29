package handlers

import (
	"net/http"
	"testing"

	"github.com/klshriharsha/snippetbox/internal/assert"
	"github.com/klshriharsha/snippetbox/internal/testutils"
)

func TestPing(t *testing.T) {
	// setup
	app := testutils.NewTestApplication(t)
	ts := testutils.NewTestServer(t, Routes(app))
	defer ts.Close()

	// make a GET request to /ping
	statusCode, _, body := ts.Get(t, "/ping")
	assert.Equal(t, statusCode, http.StatusOK)
	assert.Equal(t, string(body), "OK")
}
