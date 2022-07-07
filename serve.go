package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const defaultMimeType = "application/octet-stream"
var mimeTypes = map[string]string{
	"html": "text/html",
	"css": "text/css",
	"js": "text/javascript",
}
var templates = template.Must(template.ParseGlob("share/templates/*.html"))

func servePage(w http.ResponseWriter, r *http.Request, info *sessionInfo) {
	t := templates.Lookup(info.File + ".html")
	if t == nil {
		sendError(w, http.StatusNotFound, r.RequestURI)
		return
	}
	w.Header().Set("Content-Type", mimeTypes["html"])
	w.WriteHeader(http.StatusOK)
	if err := t.Execute(w, info); err != nil {
		log.Fatal(err)
	}
}

func serveFile(w http.ResponseWriter, r *http.Request, info *sessionInfo) {
	f, err := os.ReadFile("share/srv/" + info.File)
	if err == os.ErrNotExist {
		sendError(w, http.StatusNotFound, r.RequestURI)
		return
	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", getMimeType(info.File))
	w.WriteHeader(http.StatusOK)
	w.Write(f)
}

func getMimeType(f string) string {
	// TODO: Fix panic when there is no extension
	ext := filepath.Ext(f)[1:] // Get extension without leading "."
	if mimeType, ok := mimeTypes[ext]; ok {
		return mimeType
	}
	return defaultMimeType
}
