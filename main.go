package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("wordthy_secret_key_1029101")

const (
	sessionTime   = 5 * time.Minute
	refreshWindow = 30 * time.Second
	usernameKey   = "username"
	passwordKey   = "password"
	tokenKey      = "token"
)

var users = map[string]string{
	"han": "hanpw",
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	claims := authenticate(r)
	if claims != nil {
		w.Write([]byte(fmt.Sprintf("Authed! Valid until %d", claims.ExpiresAt)))
	} else {
		w.Write([]byte(fmt.Sprintf("Please Sign In")))
	}
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

	if !issueToken(w, username[0]) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func refreshHandler(w http.ResponseWriter, r *http.Request) {
	claims := authenticate(r)
	if claims == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > refreshWindow {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !issueToken(w, claims.Username) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func authenticate(r *http.Request) *Claims {
	cookie, err := r.Cookie(tokenKey)
	if err != nil {
		return nil
	}

	token := cookie.Value
	claims := &Claims{}
	jwtToken, err := jwt.ParseWithClaims(
		token,
		claims,
		func(_ *jwt.Token) (interface{}, error) { return jwtKey, nil },
	)

	if err != nil || !jwtToken.Valid {
		return nil
	}

	return claims
}

func issueToken(w http.ResponseWriter, username string) bool {
	expiresAt := time.Now().Add(sessionTime)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return false
	}

	http.SetCookie(w, &http.Cookie{
		Name:    tokenKey,
		Value:   tokenString,
		Expires: expiresAt,
	})
	return true
}

func makeAuthed(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if authenticate(r) != nil {
			handler(w, r)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/signin", signinHandler)
	http.HandleFunc("/refresh", refreshHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
