package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hanstar17/wordy/auth"
	"github.com/hanstar17/wordy/db"
	"github.com/hanstar17/wordy/env"
)

func main() {
	env.Init()

	db.Init()
	err := db.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Disconnect()

	auth.Init()
	auth.RegisterHandlers()

	RegisterHandlers()
	log.Fatalln(http.ListenAndServe(fmt.Sprintf("%s:%s", env.GetString("HOST"), env.GetString("PORT")), nil))
}
