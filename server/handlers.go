package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type API struct {
	store     *Store
	auth      *Auth
	uploadDir string
}

func (a *API) routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/content", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, a.store.Get())
	})

	mux.HandleFunc("POST /api/admin/login", a.login)
	mux.HandleFunc("POST /api/admin/logout", func(w http.ResponseWriter, r *http.Request) {
		a.auth.clearCookie(w)
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("GET /api/admin/session", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]bool{"authed": a.auth.authed(r)})
	})

	admin := http.NewServeMux()
	admin.HandleFunc("PUT /api/admin/content", a.putContent)
	admin.HandleFunc("POST /api/admin/works", a.createWork)
	admin.HandleFunc("PUT /api/admin/works/order", a.orderWorks)
	admin.HandleFunc("PUT /api/admin/works/{id}", a.updateWork)
	admin.HandleFunc("DELETE /api/admin/works/{id}", a.deleteWork)
	admin.HandleFunc("POST /api/admin/upload", a.upload)
	mux.Handle("/api/admin/", a.auth.Require(admin))
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		writeErr(w, http.StatusBadRequest, "bad json")
		return false
	}
	return true
}

func (a *API) login(w http.ResponseWriter, r *http.Request) {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	if !a.auth.allow(ip) {
		writeErr(w, http.StatusTooManyRequests, "слишком много попыток, подождите минуту")
		return
	}
	var in struct {
		Password string `json:"password"`
	}
	if !decode(w, r, &in) {
		return
	}
	if !a.auth.checkPassword(in.Password) {
		writeErr(w, http.StatusUnauthorized, "неверный пароль")
		return
	}
	a.auth.setCookie(w, r)
	w.WriteHeader(http.StatusNoContent)
}

// putContent replaces everything except works (managed via /works endpoints).
func (a *API) putContent(w http.ResponseWriter, r *http.Request) {
	var in Content
	if !decode(w, r, &in) {
		return
	}
	for i := range in.Tools {
		in.Tools[i].Level = min(max(in.Tools[i].Level, 1), 5)
	}
	err := a.store.Update(func(c *Content) error {
		in.Works = c.Works
		*c = in
		return nil
	})
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "save failed")
		return
	}
	writeJSON(w, http.StatusOK, a.store.Get())
}

func validWork(w http.ResponseWriter, k *Work) bool {
	if k.Cat != "home" && k.Cat != "house" && k.Cat != "commercial" {
		writeErr(w, http.StatusBadRequest, "неверная категория")
		return false
	}
	if strings.TrimSpace(k.Title) == "" {
		writeErr(w, http.StatusBadRequest, "нужно название")
		return false
	}
	return true
}

func (a *API) createWork(w http.ResponseWriter, r *http.Request) {
	var in Work
	if !decode(w, r, &in) || !validWork(w, &in) {
		return
	}
	b := make([]byte, 6)
	_, _ = rand.Read(b)
	in.ID = hex.EncodeToString(b)
	err := a.store.Update(func(c *Content) error {
		c.Works = append([]Work{in}, c.Works...)
		return nil
	})
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "save failed")
		return
	}
	writeJSON(w, http.StatusCreated, in)
}

func (a *API) updateWork(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var in Work
	if !decode(w, r, &in) || !validWork(w, &in) {
		return
	}
	in.ID = id
	found := false
	err := a.store.Update(func(c *Content) error {
		for i := range c.Works {
			if c.Works[i].ID == id {
				c.Works[i] = in
				found = true
			}
		}
		return nil
	})
	switch {
	case err != nil:
		writeErr(w, http.StatusInternalServerError, "save failed")
	case !found:
		writeErr(w, http.StatusNotFound, "not found")
	default:
		writeJSON(w, http.StatusOK, in)
	}
}

func (a *API) deleteWork(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var gone *Work
	err := a.store.Update(func(c *Content) error {
		out := c.Works[:0]
		for _, k := range c.Works {
			if k.ID == id {
				k := k
				gone = &k
				continue
			}
			out = append(out, k)
		}
		c.Works = out
		return nil
	})
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "save failed")
		return
	}
	if gone == nil {
		writeErr(w, http.StatusNotFound, "not found")
		return
	}
	a.removeUnreferenced(gone.Cover, gone.PDF)
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) orderWorks(w http.ResponseWriter, r *http.Request) {
	var in struct {
		IDs []string `json:"ids"`
	}
	if !decode(w, r, &in) {
		return
	}
	err := a.store.Update(func(c *Content) error {
		byID := map[string]Work{}
		for _, k := range c.Works {
			byID[k.ID] = k
		}
		var out []Work
		for _, id := range in.IDs {
			if k, ok := byID[id]; ok {
				out = append(out, k)
				delete(byID, id)
			}
		}
		for _, k := range c.Works { // anything not listed keeps its place at the end
			if _, ok := byID[k.ID]; ok {
				out = append(out, k)
			}
		}
		c.Works = out
		return nil
	})
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "save failed")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUpload)
	f, _, err := r.FormFile("file")
	if err != nil {
		writeErr(w, http.StatusBadRequest, "нужен файл (поле file), максимум 100 МБ")
		return
	}
	defer f.Close()
	url, size, err := saveUpload(a.uploadDir, f)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"url": url, "size": size})
}

// removeUnreferenced deletes uploaded files that no content field points to any more.
func (a *API) removeUnreferenced(urls ...string) {
	b, _ := json.Marshal(a.store.Get())
	body := string(b)
	for _, u := range urls {
		if !strings.HasPrefix(u, "/uploads/") || strings.Contains(body, `"`+u+`"`) {
			continue
		}
		_ = os.Remove(filepath.Join(a.uploadDir, filepath.Base(u)))
	}
}
