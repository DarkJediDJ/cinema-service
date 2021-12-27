package handlers

import "net/http"

func AddHalls(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	tokenAuth, err := ExtractTokenMetadata(request)
	if err != nil {
		response.Write([]byte(`{"response":"unauthorized"}`))
		return
	}
	A, err := FetchAuth(tokenAuth)
	if err != nil {
		response.Write([]byte(`{"response":"unauthorized"}`))
		return
	}
	response.Write([]byte(`{"response":"` + A + `"}`))
}
