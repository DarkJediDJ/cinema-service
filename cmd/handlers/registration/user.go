package handlers

type User struct {
	Login       string `json:"Login"`
	Password    string `json:"password"`
	AddHalls    bool   `json:"AddHalls"`
	AddMovies   bool   `json:"AddMovies"`
	AddSessions bool   `json:"AddSessions"`
}
