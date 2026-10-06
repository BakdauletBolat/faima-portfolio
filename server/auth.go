package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	cookieName = "admin_session"
	sessionTTL = 24 * time.Hour
)

type Auth struct {
	password string
	secret   []byte

	mu       sync.Mutex
	attempts map[string][]time.Time
}

func NewAuth(password, secret string) *Auth {
	return &Auth{password: password, secret: []byte(secret), attempts: map[string][]time.Time{}}
}

func (a *Auth) sign(msg string) string {
	m := hmac.New(sha256.New, a.secret)
	m.Write([]byte(msg))
	return base64.RawURLEncoding.EncodeToString(m.Sum(nil))
}

func (a *Auth) newToken() string {
	exp := strconv.FormatInt(time.Now().Add(sessionTTL).Unix(), 10)
	return exp + "." + a.sign(exp)
}

func (a *Auth) valid(token string) bool {
	exp, sig, ok := strings.Cut(token, ".")
	if !ok {
		return false
	}
	if subtle.ConstantTimeCompare([]byte(sig), []byte(a.sign(exp))) != 1 {
		return false
	}
	ts, err := strconv.ParseInt(exp, 10, 64)
	return err == nil && time.Now().Unix() < ts
}

func (a *Auth) checkPassword(p string) bool {
	return subtle.ConstantTimeCompare([]byte(p), []byte(a.password)) == 1
}

// allow permits at most 5 login attempts per ip per minute.
func (a *Auth) allow(ip string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	cut := time.Now().Add(-time.Minute)
	kept := a.attempts[ip][:0]
	for _, t := range a.attempts[ip] {
		if t.After(cut) {
			kept = append(kept, t)
		}
	}
	if len(kept) >= 5 {
		a.attempts[ip] = kept
		return false
	}
	a.attempts[ip] = append(kept, time.Now())
	return true
}

func (a *Auth) setCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name: cookieName, Value: a.newToken(), Path: "/",
		MaxAge: int(sessionTTL.Seconds()), HttpOnly: true,
		SameSite: http.SameSiteStrictMode, Secure: r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https",
	})
}

func (a *Auth) clearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{Name: cookieName, Value: "", Path: "/", MaxAge: -1, HttpOnly: true})
}

func (a *Auth) authed(r *http.Request) bool {
	c, err := r.Cookie(cookieName)
	return err == nil && a.valid(c.Value)
}

func (a *Auth) Require(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !a.authed(r) {
			writeErr(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	})
}
