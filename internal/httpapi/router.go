package httpapi

import (
	"net/http"

	"github.com/dennn/valorant-twitch-commands/internal/config"
	"github.com/dennn/valorant-twitch-commands/internal/service"
)

func NewRouter(cfg config.Config, rank *service.Rank, record *service.Record) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/valorant/rank", rankHandler(rank))
	mux.HandleFunc("GET /api/valorant/record", recordHandler(record))
	mux.HandleFunc("POST /api/valorant/record/start", recordStartHandler(record))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	return withCORS(cfg.CORSAllowedOrigin, mux)
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
