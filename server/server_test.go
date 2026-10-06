package main

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newTestAPI(t *testing.T) (*API, http.Handler) {
	t.Helper()
	dir := t.TempDir()
	st, err := OpenStore(filepath.Join(dir, "content.json"))
	if err != nil {
		t.Fatal(err)
	}
	api := &API{store: st, auth: NewAuth("pw", "secret"), uploadDir: filepath.Join(dir, "uploads")}
	mux := http.NewServeMux()
	api.routes(mux)
	return api, mux
}

func do(h http.Handler, method, url string, body any, cookie *http.Cookie) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	r := httptest.NewRequest(method, url, &buf)
	r.RemoteAddr = "1.2.3.4:5"
	if cookie != nil {
		r.AddCookie(cookie)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func login(t *testing.T, h http.Handler) *http.Cookie {
	t.Helper()
	w := do(h, "POST", "/api/admin/login", map[string]string{"password": "pw"}, nil)
	if w.Code != 204 {
		t.Fatalf("login: %d", w.Code)
	}
	return w.Result().Cookies()[0]
}

func TestAdminRequiresAuth(t *testing.T) {
	_, h := newTestAPI(t)
	for _, c := range []struct{ m, u string }{
		{"PUT", "/api/admin/content"}, {"POST", "/api/admin/works"},
		{"DELETE", "/api/admin/works/x"}, {"POST", "/api/admin/upload"},
	} {
		if w := do(h, c.m, c.u, nil, nil); w.Code != 401 {
			t.Errorf("%s %s = %d, want 401", c.m, c.u, w.Code)
		}
	}
	if w := do(h, "POST", "/api/admin/login", map[string]string{"password": "bad"}, nil); w.Code != 401 {
		t.Errorf("bad login = %d", w.Code)
	}
}

func TestLoginRateLimit(t *testing.T) {
	_, h := newTestAPI(t)
	var last int
	for i := 0; i < 7; i++ {
		last = do(h, "POST", "/api/admin/login", map[string]string{"password": "bad"}, nil).Code
	}
	if last != 429 {
		t.Errorf("got %d, want 429", last)
	}
}

func TestWorksCRUD(t *testing.T) {
	_, h := newTestAPI(t)
	ck := login(t, h)

	w := do(h, "POST", "/api/admin/works", Work{Cat: "home", Title: "Новая"}, ck)
	if w.Code != 201 {
		t.Fatalf("create %d %s", w.Code, w.Body)
	}
	var created Work
	_ = json.Unmarshal(w.Body.Bytes(), &created)

	if w := do(h, "POST", "/api/admin/works", Work{Cat: "bad", Title: "x"}, ck); w.Code != 400 {
		t.Errorf("bad cat = %d", w.Code)
	}

	var c Content
	_ = json.Unmarshal(do(h, "GET", "/api/content", nil, nil).Body.Bytes(), &c)
	if c.Works[0].ID != created.ID || len(c.Works) != 7 {
		t.Fatalf("new work should be first of 7, got %d", len(c.Works))
	}

	if w := do(h, "DELETE", "/api/admin/works/"+created.ID, nil, ck); w.Code != 204 {
		t.Errorf("delete %d", w.Code)
	}
	if w := do(h, "DELETE", "/api/admin/works/"+created.ID, nil, ck); w.Code != 404 {
		t.Errorf("second delete %d", w.Code)
	}
}

func TestPutContentKeepsWorks(t *testing.T) {
	_, h := newTestAPI(t)
	ck := login(t, h)
	var c Content
	_ = json.Unmarshal(do(h, "GET", "/api/content", nil, nil).Body.Bytes(), &c)
	c.Works = nil
	c.Tools[0].Level = 99
	w := do(h, "PUT", "/api/admin/content", c, ck)
	var out Content
	_ = json.Unmarshal(w.Body.Bytes(), &out)
	if len(out.Works) != 6 || out.Tools[0].Level != 5 {
		t.Errorf("works=%d level=%d", len(out.Works), out.Tools[0].Level)
	}
}

func multipartReq(t *testing.T, content []byte, name string) *http.Request {
	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	fw, _ := mw.CreateFormFile("file", name)
	fw.Write(content)
	mw.Close()
	r := httptest.NewRequest("POST", "/api/admin/upload", &b)
	r.Header.Set("Content-Type", mw.FormDataContentType())
	return r
}

func TestUploadValidatesMagicBytes(t *testing.T) {
	api, h := newTestAPI(t)
	ck := login(t, h)

	r := multipartReq(t, append([]byte("%PDF-1.4\n"), bytes.Repeat([]byte("x"), 2000)...), "a.pdf")
	r.AddCookie(ck)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 201 {
		t.Fatalf("pdf upload %d %s", w.Code, w.Body)
	}
	var out map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &out)
	if !strings.HasSuffix(out["url"], ".pdf") {
		t.Errorf("url %q", out["url"])
	}
	if _, err := os.Stat(filepath.Join(api.uploadDir, filepath.Base(out["url"]))); err != nil {
		t.Error(err)
	}

	// an executable renamed to .pdf must be rejected
	r = multipartReq(t, []byte("MZ\x90\x00 not a pdf"), "evil.pdf")
	r.AddCookie(ck)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 400 {
		t.Errorf("fake pdf = %d, want 400", w.Code)
	}
}
