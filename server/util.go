package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func bytesReader(b []byte) io.Reader { return bytes.NewReader(b) }

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
