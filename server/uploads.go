package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const maxUpload = 100 << 20 // 100 MB

var allowedTypes = map[string]string{
	"image/jpeg":      ".jpg",
	"image/png":       ".png",
	"image/webp":      ".webp",
	"application/pdf": ".pdf",
}

// saveUpload sniffs the content type from magic bytes (ignoring the client's
// filename/header), stores the file under a random name and returns its URL
// path and human-readable size.
func saveUpload(dir string, src io.Reader) (url, size string, err error) {
	head := make([]byte, 512)
	n, err := io.ReadFull(src, head)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return "", "", err
	}
	head = head[:n]
	ext, ok := allowedTypes[sniff(head)]
	if !ok {
		return "", "", fmt.Errorf("unsupported file type")
	}

	id := make([]byte, 12)
	if _, err := rand.Read(id); err != nil {
		return "", "", err
	}
	name := hex.EncodeToString(id) + ext
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", "", err
	}
	f, err := os.OpenFile(filepath.Join(dir, name), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return "", "", err
	}
	written, err := io.Copy(f, io.MultiReader(bytesReader(head), src))
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		os.Remove(filepath.Join(dir, name))
		return "", "", err
	}
	return "/uploads/" + name, humanSize(written), nil
}

func sniff(head []byte) string {
	if len(head) >= 5 && string(head[:5]) == "%PDF-" {
		return "application/pdf"
	}
	if len(head) >= 12 && string(head[:4]) == "RIFF" && string(head[8:12]) == "WEBP" {
		return "image/webp"
	}
	t := http.DetectContentType(head)
	if t == "image/jpeg" || t == "image/png" {
		return t
	}
	return ""
}

func humanSize(n int64) string {
	if n >= 1<<20 {
		return fmt.Sprintf("%.1f МБ", float64(n)/(1<<20))
	}
	kb := n / 1024
	if kb < 1 {
		kb = 1
	}
	return fmt.Sprintf("%d КБ", kb)
}
