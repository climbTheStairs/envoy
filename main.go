package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

type sessionInfo struct {
	Messages []string
	File string
	Username string
}

func main() {
	http.HandleFunc("/", mux)
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func mux(w http.ResponseWriter, r *http.Request) {
	info := &sessionInfo{
		File: r.RequestURI[1:], // Remove leading "/"
		Username: verifyUserAndGetUsername(r),
	}
	if info.File == "" {
		info.File = "index"
	}
	if info.File == "index" && info.Username == "" {
		info.File = "guest"
	}
	if info.File == "register" && (r.Method == "GET" || r.Method == "HEAD") && info.Username != "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if info.File == "register" && r.Method == "POST" {
		register(w, r, info)
		return
	}
	if info.File == "login" && (r.Method == "GET" || r.Method == "HEAD") && info.Username != "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	if info.File == "login" && r.Method == "POST" {
		login(w, r, info)
		return
	}
	if info.File == "logout" && r.Method == "POST" {
		logout(w, r)
		return
	}
	if filepath.Ext(info.File) == "" && (r.Method == "GET" || r.Method == "HEAD") {
		servePage(w, r, info)
		return
	}
	if r.Method == "GET" || r.Method == "HEAD" {
		serveFile(w, r, info)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func sendError(w http.ResponseWriter, status int, error string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	fmt.Fprintf(w, "%d %s %s\n",
		status, http.StatusText(status), error)
}

