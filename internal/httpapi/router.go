package httpapi

import (
	"net/http"
	"os"
	"strings"

	"github.com/dennn/valorant-twitch-commands/internal/config"
	"github.com/dennn/valorant-twitch-commands/internal/service"
)

func NewRouter(cfg config.Config, rank *service.Rank, record *service.Record) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/valorant/rank", rankHandler(rank))
	mux.HandleFunc("GET /api/valorant/record", recordHandler(record))
	mux.HandleFunc("GET /api/valorant/line", lineHandler(rank, record))
	mux.HandleFunc("POST /api/valorant/record/start", recordStartHandler(record))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("GET /format/riot.txt", riotVerification)
	if dir := staticDir(); dir != "" {
		mux.Handle("/", spa(dir))
	}
	return withCORS(cfg.CORSAllowedOrigin, mux)
}

func staticDir() string {
	dir := os.Getenv("STATIC_DIR")
	if dir == "" {
		dir = "frontend/dist"
	}
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return ""
	}
	return dir
}

func spa(dir string) http.Handler {
	files := http.Dir(dir)
	fileServer := http.FileServer(files)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			if _, err := files.Open(strings.TrimPrefix(r.URL.Path, "/")); err != nil {
				r.URL.Path = "/"
			}
		}
		fileServer.ServeHTTP(w, r)
	})
}

func withCORS(origin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
