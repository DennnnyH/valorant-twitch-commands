package config

import (
	"os"
	"strings"
)

type Config struct {
	RiotAPIKey        string
	HenrikAPIKey      string
	Port              string
	CORSAllowedOrigin string
}

func Load() Config {
	loadDotEnv(".env")
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

// loadDotEnv fills variables from a local .env file. Existing non-empty
// environment variables win so a shell export can still override the file.
func loadDotEnv(path string) {
	body, err := os.ReadFile(path)
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(body), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		if key == "" {
			continue
		}
		if current, exists := os.LookupEnv(key); exists && current != "" {
			continue
		}
		_ = os.Setenv(key, value)
	}
}
