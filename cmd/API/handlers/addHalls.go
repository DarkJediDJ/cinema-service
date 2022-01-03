package handlers

import (
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
)

func AddHalls(response http.ResponseWriter, request *http.Request) {
	c, err := request.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			response.WriteHeader(http.StatusUnauthorized)
			return
		}
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	tknStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			response.WriteHeader(http.StatusUnauthorized)
			return
		}
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		response.WriteHeader(http.StatusUnauthorized)
		return
	}
	response.Write([]byte(strconv.FormatBool(claims.AddHalls)))
}
