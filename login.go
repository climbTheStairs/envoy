package main

import (
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if !verifyUser(username) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(f(`User "%s" does not exist`, username)))
		return
	}

	if !verifyPassword(username, password) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(f(`Password "%s" is incorrect`, password)))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: "username",
		Value: username,
	})
	http.SetCookie(w, &http.Cookie{
		Name: "password",
		Value: password,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(f(`Successfully logged in as "%s"!`, username)))
}

func verifyUser(username string) bool {
	_, ok := accounts[username]
	return ok
}

func verifyPassword(username, password string) bool {
	return password == accounts[username]
}

