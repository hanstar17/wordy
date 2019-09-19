package env

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type errNotFound struct {
	varName string
}

func (e errNotFound) Error() string {
	return e.varName + "Env variable is not found."
}

type errParseFailed struct {
	varName string
}

func (e errParseFailed) Error() string {
	return e.varName + "Env variable couldn't be parsed."
}

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
}

func GetString(name string) string {
	s, ok := os.LookupEnv(name)
	if !ok {
		log.Fatalln(errNotFound{name})
	}
	return s
}

func GetStringFallback(name string, fallback string) string {
	s, ok := os.LookupEnv(name)
	if !ok {
		return fallback
	}
	return s
}

func GetInt(name string) int {
	s, ok := os.LookupEnv(name)
	if !ok {
		log.Fatalln(errNotFound{name})
	}

	v, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalln(errParseFailed{name})
	}
	return v
}

func GetIntFallback(name string, fallback int) int {
	s, ok := os.LookupEnv(name)
	if !ok {
		return fallback
	}

	v, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return v
}
