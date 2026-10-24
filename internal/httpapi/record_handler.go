package httpapi

import (
	"net/http"
	"time"

	"github.com/dennn/valorant-twitch-commands/internal/service"
)

func recordHandler(record *service.Record) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		resp, err := record.Lookup(r.Context(), service.Query{
			Name:     q.Get("name"),
			Tag:      q.Get("tag"),
			Region:   q.Get("region"),
			Platform: q.Get("platform"),
		})
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, resp)
	}
}

func recordStartHandler(record *service.Record) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		started := record.Start(time.Now())
		writeJSON(w, http.StatusOK, map[string]string{
			"startedAt": started.Format(time.RFC3339),
		})
	}
}
