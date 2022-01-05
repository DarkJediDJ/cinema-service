package handlers

import (
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

func getTokenData(tknStr string, claims *Claims, response http.ResponseWriter) (*jwt.Token, *Claims, error) {
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			response.WriteHeader(http.StatusUnauthorized)
			return nil, nil, err
		}
		response.WriteHeader(http.StatusBadRequest)
		return nil, nil, err
	}
	if !tkn.Valid {
		response.WriteHeader(http.StatusUnauthorized)
		return nil, nil, errors.New(tkn.Claims.Valid().Error())
	}
	return tkn, claims, nil
}
func getCoockie(name string, response http.ResponseWriter, request *http.Request) (string, error) {
	c, err := request.Cookie(name)
	if err != nil {
		if err == http.ErrNoCookie {
			response.WriteHeader(http.StatusUnauthorized)
			return "", err
		}
		response.WriteHeader(http.StatusBadRequest)
		return "", err
	}
	return c.Value, nil
}

func validToken(response http.ResponseWriter, request *http.Request) (*Claims, bool, error) {
	claims := &Claims{}
	tknStr, err := getCoockie("token", response, request)
	if err != nil {
		refreshTkn, err := getCoockie("refresh_token", response, request)
		if err != nil {
			response.WriteHeader(http.StatusUnauthorized)
			return nil, false, err
		}
		refreshToken, claims, err := getTokenData(refreshTkn, claims, response)
		if err != nil {
			response.WriteHeader(http.StatusUnauthorized)
			return nil, false, err
		}
		if refreshToken.Valid {
			tknStr, err = GenerateJWT(claims.AddHalls, claims.AddMovies, claims.AddSessions, claims.UserID)
			response.Write([]byte("Used refresh token")) // Just for displaying
			if err != nil {

				response.WriteHeader(http.StatusUnauthorized)
				return nil, false, err
			}
		}
	}
	token, claims, err := getTokenData(tknStr, claims, response)
	if err != nil {
		response.WriteHeader(http.StatusUnauthorized)
		return nil, false, err
	}
	return claims, token.Valid, nil
}
