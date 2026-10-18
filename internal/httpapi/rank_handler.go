package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dennn/valorant-twitch-commands/internal/service"
)

func rankHandler(rank *service.Rank) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		resp, err := rank.Lookup(r.Context(), service.Query{
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

func writeError(w http.ResponseWriter, err error) {
	var apiErr *service.APIError
	if !errors.As(err, &apiErr) {
		apiErr = &service.APIError{
			Status:  http.StatusBadGateway,
			Code:    "UPSTREAM",
			Message: "rank lookup failed",
		}
	}
	if apiErr.RetryAfter != "" {
		w.Header().Set("Retry-After", apiErr.RetryAfter)
	}
	writeJSON(w, apiErr.Status, map[string]string{
		"error": apiErr.Message,
		"code":  apiErr.Code,
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
