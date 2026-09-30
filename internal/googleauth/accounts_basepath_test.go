package googleauth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNormalizeBasePath(t *testing.T) {
	for in, want := range map[string]string{"": "", "/": "", "gog": "/gog", "/gog/": "/gog", " /a/b/ ": "/a/b"} {
		if got := NormalizeBasePath(in); got != want {
			t.Errorf("NormalizeBasePath(%q)=%q want %q", in, got, want)
		}
	}
}

func TestManagerApplicationBasePath(t *testing.T) {
	app := newTestManagerApplication(t, ManagerOptions{BasePath: "gog/"}, ManagerDependencies{})
	h := app.Handler()

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/gog/", nil))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), `const basePath = "/gog";`) {
		t.Fatalf("index: %d %s", rec.Code, rec.Body.String()[strings.Index(rec.Body.String(), "basePath")-20:][:80])
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/gog", nil))
	if rec.Code != http.StatusMovedPermanently {
		t.Fatalf("redirect: %d", rec.Code)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/gog/accounts", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("accounts: %d", rec.Code)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/accounts", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("unprefixed: %d", rec.Code)
	}
}
