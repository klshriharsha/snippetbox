package testutils

import (
	"html"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form"
	"github.com/klshriharsha/snippetbox/cmd/web/config"
	"github.com/klshriharsha/snippetbox/internal/models/mocks"
)

// NewTestApplication initializes a new `Application` for the test environment with discarded logs
func NewTestApplication(t *testing.T) *config.Application {
	cache, err := config.NewTemplateCache()
	if err != nil {
		t.Fatal(err)
	}
	decoder := form.NewDecoder()
	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	return &config.Application{
		ErrorLog: log.New(io.Discard, "", 0),
		InfoLog:  log.New(io.Discard, "", 0),

		Snippets: &mocks.SnippetModel{},
		Users:    &mocks.UserModel{},

		TemplateCache:  cache,
		FormDecoder:    decoder,
		SessionManager: sessionManager,
	}
}

// TestServer encapsulates the test server and provides functions to make requests
type TestServer struct {
	*httptest.Server
}

// NewTestServer initializes a new HTTPS server for the test environment
func NewTestServer(t *testing.T, h http.Handler) *TestServer {
	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	ts.Client().Jar = jar

	// disable automatically following a redirect response from a server
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &TestServer{ts}
}

// Get method makes a GET request to the given `urlPath`
func (ts *TestServer) Get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, string(body)
}

// Post method makes a POST request to the given `urlPath` with the given `form` payload
func (ts *TestServer) Post(t *testing.T, urlPath string, form url.Values) (int, http.Header, string) {
	rs, err := ts.Client().PostForm(ts.URL+urlPath, form)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, string(body)
}

var csrfTokenRX = regexp.MustCompile(`<input type='hidden' name='csrf_token' value='(.+)'>`)

// ExtractCSRFToken finds and returns the CSRF token in the given HTML `body`
func ExtractCSRFToken(t *testing.T, body string) string {
	matches := csrfTokenRX.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no CSRF token found in body")
	}

	return html.UnescapeString(matches[1])
}
