package config

import "os"

type User struct {
	UserID int `json:"UserID"`
	Login       string `json:"Login"`
	Password    string `json:"Password"`
	AddHalls    bool   `json:"AddHalls"`
	AddMovies   bool   `json:"AddMovies"`
	AddSessions bool   `json:"AddSessions"`
}

type RedisConfig struct {
	Host string
	Port string
}

func LoadRedisConfig() RedisConfig {
	return RedisConfig{
		Host: os.Getenv("REDIS_HOST"),
		Port: os.Getenv("REDIS_PORT"),
	}
}
