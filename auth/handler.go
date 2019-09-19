package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/hanstar17/wordy/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Authenticate(r *http.Request) bool {
	return authenticate(r) != nil
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
	expiresAt := time.Now().Add(sessionLifetime)
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

func MakeAuthed(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if Authenticate(r) {
			handler(w, r)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	var creds credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	filter := bson.D{
		{"username", bson.D{{"$eq", creds.Username}}},
		{"password", bson.D{{"$eq", creds.Password}}},
	}
	user := db.FindUser(filter)
	if user == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !issueToken(w, creds.Username) {
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

func signupHandler(w http.ResponseWriter, r *http.Request) {
	var signup db.User
	err := json.NewDecoder(r.Body).Decode(&signup)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	filter := bson.D{{
		"username", bson.D{{"$eq", signup.Username}},
	}}

	user := db.FindUser(filter)
	if user != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	insertedId := db.InsertUser(signup)
	if insertedId == primitive.NilObjectID {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user = db.FindUser(bson.D{{"_id", bson.D{{"$eq", insertedId}}}})
	if user == nil {
		log.Fatalln("signup failed in db silently.")
	}
}
