package handlers

import "net/http"

func Refresh(response http.ResponseWriter, request *http.Request) {
	c, err := request.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			response.WriteHeader(http.StatusUnauthorized)
			return
		}
		response.WriteHeader(http.StatusBadRequest)
		return
	}
	refreshToken, err := RefreshJWT(response, c.Value)
	if err != nil {
		response.WriteHeader(http.StatusUnauthorized)
		return
	}
	http.SetCookie(response, &http.Cookie{
		Name:    "session_token",
		Value:   refreshToken,
		Expires: ExpirationTime,
	})
}
