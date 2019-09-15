package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecretKey = []byte("wordthy_secret_key_1029101")

const (
	sessionTime = 5 * time.Minute
	usernameKey = "username"
	passwordKey = "password"
)

var users = map[string]string{
	"han": "hanpw",
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username, ok := r.Form[usernameKey]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	password, ok := r.Form[passwordKey]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userPassword, ok := users[username[0]]
	if !ok || password[0] != userPassword {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expiresAt := time.Now().Add(sessionTime)
	claims := &Claims{
		Username: username[0],
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//w.WriteHeader(http.StatusOK)
	http.SetCookie(w, &http.Cookie{
		Name:    "session-token",
		Value:   tokenString,
		Expires: expiresAt,
	})
}

func main() {
	http.HandleFunc("/signin", signinHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
