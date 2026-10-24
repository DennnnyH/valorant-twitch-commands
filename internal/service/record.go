package service

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dennn/valorant-twitch-commands/internal/henrik"
)

// Record counts competitive wins and losses since the stream session started.
type Record struct {
	henrik *henrik.Client

	mu        sync.Mutex
	startedAt time.Time
}

func NewRecord(henrikClient *henrik.Client, startedAt time.Time) *Record {
	return &Record{henrik: henrikClient, startedAt: startedAt.UTC()}
}

func (r *Record) Start(at time.Time) time.Time {
	if at.IsZero() {
		at = time.Now()
	}
	at = at.UTC()
	r.mu.Lock()
	r.startedAt = at
	r.mu.Unlock()
	return at
}

func (r *Record) StartedAt() time.Time {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.startedAt
}

type RecordResponse struct {
	StartedAt string       `json:"startedAt"`
	Wins      int          `json:"wins"`
	Losses    int          `json:"losses"`
	Games     int          `json:"games"`
	RRChange  int          `json:"rrChange"`
	Matches   []RecordGame `json:"matches"`
}

type RecordGame struct {
	MatchID  string `json:"matchId"`
	Map      string `json:"map"`
	Result   string `json:"result"`
	RRChange int    `json:"rrChange"`
	RR       int    `json:"rr"`
	PlayedAt string `json:"playedAt"`
}

func (r *Record) Lookup(ctx context.Context, q Query) (RecordResponse, error) {
	q.Name = strings.TrimSpace(q.Name)
	q.Tag = strings.TrimSpace(q.Tag)
	q.Region = strings.ToLower(strings.TrimSpace(q.Region))
	q.Platform = strings.ToLower(strings.TrimSpace(q.Platform))
	if q.Platform == "" {
		q.Platform = "pc"
	}
	if q.Name == "" || q.Tag == "" || q.Region == "" {
		return RecordResponse{}, &APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_PARAMS",
			Message: "name, tag, and region are required",
		}
	}

	history, err := r.henrik.GetMMRHistoryV2(ctx, q.Region, q.Platform, q.Name, q.Tag)
	if err != nil {
		return RecordResponse{}, mapUpstream("henrik", err)
	}

	started := r.StartedAt()
	out := RecordResponse{
		StartedAt: started.Format(time.RFC3339),
		Matches:   []RecordGame{},
	}
	for _, entry := range history {
		played, ok := entryTime(entry)
		if !ok || played.Before(started) {
			continue
		}
		result := "even"
		switch {
		case entry.LastChange > 0:
			result = "win"
			out.Wins++
		case entry.LastChange < 0:
			result = "loss"
			out.Losses++
		}
		out.Games++
		out.RRChange += entry.LastChange
		out.Matches = append(out.Matches, RecordGame{
			MatchID:  entry.MatchID,
			Map:      entry.Map.Name,
			Result:   result,
			RRChange: entry.LastChange,
			RR:       entry.RR,
			PlayedAt: played.UTC().Format(time.RFC3339),
		})
	}
	return out, nil
}

func entryTime(entry henrik.HistoryEntry) (time.Time, bool) {
	if entry.DateRaw > 0 {
		return unixTime(entry.DateRaw), true
	}
	if entry.Date == "" {
		return time.Time{}, false
	}
	if n, err := strconv.ParseInt(entry.Date, 10, 64); err == nil {
		return unixTime(n), true
	}
	for _, layout := range []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05.000Z",
		"Monday, January 2, 2006 3:04 PM",
	} {
		if t, err := time.Parse(layout, entry.Date); err == nil {
			return t.UTC(), true
		}
	}
	return time.Time{}, false
}

func unixTime(n int64) time.Time {
	if n > 1_000_000_000_000 {
		return time.UnixMilli(n).UTC()
	}
	return time.Unix(n, 0).UTC()
}
