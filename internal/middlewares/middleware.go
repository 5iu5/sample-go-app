package middlewares

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Retrieve cookie from r
		log.Println(r)
		tokenString, err := r.Cookie("Authorization")
		log.Println("tokenString: ", tokenString)
		if err != nil {
			if err == http.ErrNoCookie {
				log.Println("Error 1 here: ", err)
				http.Error(w, "cookie not found", http.StatusBadRequest)
				return
			}
			http.Error(w, "Error retrieving cookie", http.StatusInternalServerError)
			fmt.Println("Error retrieving cookie: err")
			return
		}

		//validate cookie
		token, err := jwt.Parse(tokenString.Value, func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("SECRET")), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
		if err != nil || !token.Valid {
			http.Error(w, "invalid cookie, unauthorized", http.StatusBadRequest)
			log.Println("Error 2 here")
			log.Println("invalid cookie: ", err)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "claims not in expected format", http.StatusUnauthorized)
			return
		}

		//check expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			http.Error(w, "token expired", http.StatusUnauthorized)
			return
		}

		//cookie is valid and user is authorized
		//extract userID from token's claims
		userID := int(claims["sub"].(float64))
		if userID == 0 {
			http.Error(w, "Error extracting user id from claim", http.StatusUnauthorized)
			return
		}
		// store userID in the request context

		key := CtxUserKey("user_id") //avoid collision by using struct instead of string as key
		ctx := context.WithValue(r.Context(), key, userID)
		// call the next handler with the new context
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

type CtxUserKey string
