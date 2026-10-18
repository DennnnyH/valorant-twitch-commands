package config

import "os"

type Config struct {
	RiotAPIKey        string
	HenrikAPIKey      string
	Port              string
	CORSAllowedOrigin string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	origin := os.Getenv("CORS_ALLOWED_ORIGIN")
	if origin == "" {
		origin = "http://localhost:5173"
	}
	return Config{
		RiotAPIKey:        os.Getenv("RIOT_API_KEY"),
		HenrikAPIKey:      os.Getenv("HENRIK_API_KEY"),
		Port:              port,
		CORSAllowedOrigin: origin,
	}
}
