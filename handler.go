package main

import (
	"fmt"
	"net/http"

	"github.com/hanstar17/wordy/auth"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if auth.Authenticate(r) {
		w.Write([]byte(fmt.Sprintf("Authed!")))
	} else {
		w.Write([]byte(fmt.Sprintf("Please Sign In")))
	}
}

func RegisterHandlers() {
	http.HandleFunc("/", indexHandler)
}
