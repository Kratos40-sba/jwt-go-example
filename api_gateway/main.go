package main

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"os"
)

var mySecretKey = []byte(os.Getenv("SECRET_KEY"))

func isAuth(e http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			// jwt parsing logic
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("invalid Signing Method \n")
				}
				if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
					return nil, fmt.Errorf("expired token \n")
				}
				aud := "billing.jwtgo.io"
				checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
				if !checkAud {
					return nil, fmt.Errorf("invalid aud \n")
				}
				iss := "jwtgo.io"
				checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
				if !checkIss {
					return nil, fmt.Errorf("invalid iss \n")
				}
				return mySecretKey, nil

			})
			if err != nil {
				_, _ = fmt.Fprintf(w, err.Error())
			}
			if token != nil {
				if token.Valid {
					e(w, r)
				}
			}
		} else {
			_, _ = fmt.Fprintf(w, "No Authorization Token provided \n")
		}
	})
}
func homePage(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Super Secret Informations \n")
}
func handleRequests() {
	http.Handle("/", isAuth(homePage))
	log.Fatalln(http.ListenAndServe(":9001", nil))
}
func main() {
	fmt.Println("SERVER STARTS")
	handleRequests()
}
