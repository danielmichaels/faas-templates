package handlers

import (
	"fmt"
	jwt "github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
	"time"
)

func JwtExample() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		domain := os.Getenv("DOMAIN")

		if r.Body != nil {
			defer r.Body.Close()
		}

		t := jwt.New(jwt.GetSigningMethod("HS256"))
		t.Claims = jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 48)),
			Audience:  jwt.ClaimStrings{fmt.Sprintf("%s/function/go-example", domain)},
		}

		signingKey := []byte("secret")
		res, err := t.SignedString(signingKey)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(res))
	}
}
