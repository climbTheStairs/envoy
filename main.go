package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

const webroot = "/home/stairs/share/envoy/srv"
const defaultMimeType = "application/octet-stream"
var mimeTypes = map[string]string{
	"html": "text/html",
	"css": "text/css",
	"js": "text/javascript",
}

func main() {
	http.HandleFunc("/", respond)
	http.ListenAndServe(":8888", nil)
}

func f(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}

func respond(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" || r.Method == "HEAD" {
		servePage(w, r)
		return
	}
	if r.Method == "POST" {
		login(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func servePage(w http.ResponseWriter, r *http.Request) {
	url := webroot + r.RequestURI
	if r.RequestURI == "/" {
		url += "index.html"
	}

	f, err := os.ReadFile(url)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		msg := fmt.Sprintf("%d %s %s\n",
			http.StatusNotFound,
			http.StatusText(http.StatusNotFound),
			r.RequestURI)
		w.Write([]byte(msg))
		return
	}

	w.Header().Set("Content-Type", getMimeType(url))
	w.WriteHeader(http.StatusOK)
	w.Write(f)
}

func getMimeType(url string) string {
	ext := filepath.Ext(url)[1:]
	if mimeType, ok := mimeTypes[ext]; ok {
		return mimeType
	}
	return defaultMimeType
}

