package main

import (
	"fmt"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if !verifyUser(username) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "User \"%s\" does not exist\n", username)
		return
	}

	if !verifyPassword(username, password) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Password \"%s\" is incorrect\n", password)
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
	fmt.Fprintf(w, "Successfully logged in as \"%s\"!\n", username)
}

func verifyUser(username string) bool {
	_, ok := accounts[username]
	return ok
}

func verifyPassword(username, password string) bool {
	return password == accounts[username]
}

