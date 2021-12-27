package config

type User struct {
	Login       string `json:"Login"`
	Password    string `json:"Password"`
	AddHalls    bool   `json:"AddHalls"`
	AddMovies   bool   `json:"AddMovies"`
	AddSessions bool   `json:"AddSessions"`
}

