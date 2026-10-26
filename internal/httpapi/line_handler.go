package httpapi

import (
	"net/http"
	"unicode/utf8"

	"github.com/dennn/valorant-twitch-commands/internal/service"
)

func lineHandler(rank *service.Rank, record *service.Record) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		query := service.Query{
			Name:     q.Get("name"),
			Tag:      q.Get("tag"),
			Region:   q.Get("region"),
			Platform: q.Get("platform"),
		}
		template := q.Get("template")
		if utf8.RuneCountInString(template) > 240 {
			writePlain(w, http.StatusBadRequest, "template must be 240 characters or fewer")
			return
		}

		rankResp, err := rank.Lookup(r.Context(), query)
		if err != nil {
			writePlainError(w, err)
			return
		}
		recordResp, err := record.Lookup(r.Context(), query)
		if err != nil {
			writePlainError(w, err)
			return
		}

		line := service.FormatLine(template, service.LineData{
			Name:         rankResp.RiotID.Name,
			Tag:          rankResp.RiotID.Tag,
			Rank:         rankResp.Current.RankName,
			RR:           rankResp.Current.RR,
			LastRRChange: rankResp.Current.LastRRChange,
			Elo:          rankResp.Current.Elo,
			Peak:         rankResp.Peak.RankName,
			Wins:         recordResp.Wins,
			Losses:       recordResp.Losses,
			Games:        recordResp.Games,
			RRChange:     recordResp.RRChange,
		})
		writePlain(w, http.StatusOK, line)
	}
}

func writePlainError(w http.ResponseWriter, err error) {
	message := "rank lookup failed"
	status := http.StatusBadGateway
	if apiErr, ok := err.(*service.APIError); ok {
		message = apiErr.Message
		status = apiErr.Status
		if apiErr.RetryAfter != "" {
			w.Header().Set("Retry-After", apiErr.RetryAfter)
		}
	}
	writePlain(w, status, message)
}

func writePlain(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body))
}
