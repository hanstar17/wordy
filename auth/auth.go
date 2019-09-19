package auth

import (
	"net/http"
	"time"

	"github.com/hanstar17/wordy/env"
)

const (
	tokenKey = "token"
)

var refreshWindow time.Duration
var sessionLifetime time.Duration
var jwtKey []byte

func Init() {
	refreshWindow = time.Duration(env.GetInt("AUTH_SESSION_REFRESH_WINDOW_SEC")) * time.Second
	sessionLifetime = time.Duration(env.GetInt("AUTH_SESSION_LIFETIME_SEC")) * time.Second
	jwtKey = []byte(env.GetString("AUTH_JWT_SECRET_KEY"))
}

func RegisterHandlers() {
	http.HandleFunc("/signin", signinHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/refresh", refreshHandler)
}
