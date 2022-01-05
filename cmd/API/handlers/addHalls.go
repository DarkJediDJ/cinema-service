package handlers

import (
	"net/http"
	"strconv"
)

func AddHalls(response http.ResponseWriter, request *http.Request) {
	claims, valid, err := validToken(response, request)
	if err != nil {
		response.WriteHeader(http.StatusUnauthorized)
		return
	}
	if !valid {
		return
	}
	response.Write([]byte(`{"AddHalls":"` + strconv.FormatBool(claims.AddHalls) + `"AddMovies":"` + strconv.FormatBool(claims.AddMovies) + `"AddSessions":"` + strconv.FormatBool(claims.AddSessions) + `"}`))
}


