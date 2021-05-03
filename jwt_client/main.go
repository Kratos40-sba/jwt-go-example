package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var mySecretKey = []byte(os.Getenv("SECRET_KEY"))

// export SECRET_KEY=verysecret
func GetJwt() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["client"] = "Krissanawat"
	claims["aud"] = "billing.jwtgo.io"
	claims["iss"] = "jwtgo.io"
	claims["exp"] = time.Now().Add(time.Minute * 2).Unix()
	tokenString, err := token.SignedString(mySecretKey)
	if err != nil {
		_ = fmt.Errorf("something went wrong : %s", err.Error())
		return "", err
	}
	return tokenString, nil

}
func Index(w http.ResponseWriter, r *http.Request) {
	validToken, err := GetJwt()
	if err != nil {
		log.Println("Failed to generate token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = fmt.Fprintf(w, string(validToken))

}
func handleRequests() {
	http.HandleFunc("/", Index)
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
func main() {
	handleRequests()
}
