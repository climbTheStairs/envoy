package main

import (
	"errors"
	"net/http"
	"time"
)

func register(w http.ResponseWriter, r *http.Request, info *sessionInfo) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm-password")
	info.Username = username

	if username == "" {
		info.Messages = append(info.Messages,
			"Please actually enter a username")
	} else if len(username) > 50 {
		info.Messages = append(info.Messages,
			"Username must not exceed 50 characters in length")
	} else if verifyUsername(username) {
		info.Messages = append(info.Messages,
			`User "` + username + `" already exists`)
	}
	if len(password) < 4 {
		info.Messages = append(info.Messages,
			"Your password must be at least 4 characters long")
	}
	if password != confirmPassword {
		info.Messages = append(info.Messages,
			"Passwords do not match")
	}

	if len(info.Messages) != 0 {
		servePage(w, r, info)
		return
	}
		
	accounts[username] = password
	http.Redirect(w, r, "/login", http.StatusFound)
}

func login(w http.ResponseWriter, r *http.Request, info *sessionInfo) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	if !verifyUsername(username) {
		info.Messages = []string{`User "` + username +
			`" does not exist`}
		servePage(w, r, info)
		return
	}
	if !verifyPassword(username, password) {
		info.Username = username
		info.Messages = []string{"Password is incorrect"}
		servePage(w, r, info)
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
	http.Redirect(w, r, "/", http.StatusFound)
}

// logout signs out the user by deleting their cookies
// before redirecting them to the login page.
func logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name: "username",
		Value: "",
		Expires: time.Unix(0, 0),
	})
	http.SetCookie(w, &http.Cookie{
		Name: "password",
		Value: "",
		Expires: time.Unix(0, 0),
	})
	http.Redirect(w, r, "/login", 302)
}

func verifyUserAndGetUsername(r *http.Request) string {
	// The only error that r.Cookie can return is http.ErrNoCookie
	usernameCookie, err := r.Cookie("username")
	if errors.Is(err, http.ErrNoCookie) {
		return ""
	}
	passwordCookie, err := r.Cookie("password")
	if errors.Is(err, http.ErrNoCookie) {
		return ""
	}
	username := usernameCookie.Value
	password := passwordCookie.Value
	if !verifyPassword(username, password) {
		return ""
	}
	return username
}

func verifyUsername(username string) bool {
	_, ok := accounts[username]
	return ok
}

func verifyPassword(username, password string) bool {
	// verifyPassword automatically verifies username as well
	correctPassword, ok := accounts[username]
	return ok && password == correctPassword
}
