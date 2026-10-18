// Package service maps Riot account data and Henrik MMR into the public rank API.
package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dennn/valorant-twitch-commands/internal/henrik"
	"github.com/dennn/valorant-twitch-commands/internal/riot"
	"github.com/dennn/valorant-twitch-commands/internal/upstream"
)

type Rank struct {
	riot   *riot.Client
	henrik *henrik.Client
}

func NewRank(riotClient *riot.Client, henrikClient *henrik.Client) *Rank {
	return &Rank{riot: riotClient, henrik: henrikClient}
}

type Query struct {
	Name     string
	Tag      string
	Region   string
	Platform string
}

type Response struct {
	RiotID  RiotID  `json:"riotId"`
	PUUID   string  `json:"puuid"`
	Current Current `json:"current"`
	Peak    Peak    `json:"peak"`
	Season  Season  `json:"season"`
	Sources Sources `json:"sources"`
}

type RiotID struct {
	Name string `json:"name"`
	Tag  string `json:"tag"`
}

type Current struct {
	TierID       int    `json:"tierId"`
	RankName     string `json:"rankName"`
	RR           int    `json:"rr"`
	LastRRChange int    `json:"lastRrChange"`
	Elo          int    `json:"elo"`
}

type Peak struct {
	TierID   int    `json:"tierId"`
	RankName string `json:"rankName"`
}

type Season struct {
	Wins        int    `json:"wins"`
	Games       int    `json:"games"`
	EndRankName string `json:"endRankName"`
}

type Sources struct {
	Account string `json:"account"`
	MMR     string `json:"mmr"`
}

// APIError is returned to the HTTP layer with a stable code and status.
type APIError struct {
	Status     int
	Code       string
	Message    string
	RetryAfter string
}

func (e *APIError) Error() string {
	return e.Message
}

func (r *Rank) Lookup(ctx context.Context, q Query) (Response, error) {
	q.Name = strings.TrimSpace(q.Name)
	q.Tag = strings.TrimSpace(q.Tag)
	q.Region = strings.ToLower(strings.TrimSpace(q.Region))
	q.Platform = strings.ToLower(strings.TrimSpace(q.Platform))
	if q.Platform == "" {
		q.Platform = "pc"
	}
	if q.Name == "" || q.Tag == "" || q.Region == "" {
		return Response{}, &APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_PARAMS",
			Message: "name, tag, and region are required",
		}
	}

	account, err := r.riot.GetAccountByRiotID(ctx, q.Region, q.Name, q.Tag)
	if err != nil {
		return Response{}, mapUpstream("riot", err)
	}

	mmr, err := r.henrik.GetMMRv3(ctx, q.Region, q.Platform, q.Name, q.Tag)
	if err != nil {
		return Response{}, mapUpstream("henrik", err)
	}

	resp := Response{
		RiotID: RiotID{Name: account.GameName, Tag: account.TagLine},
		PUUID:  account.PUUID,
		Current: Current{
			TierID:       mmr.Current.Tier.ID,
			RankName:     mmr.Current.Tier.Name,
			RR:           mmr.Current.RR,
			LastRRChange: mmr.Current.LastChange,
			Elo:          mmr.Current.Elo,
		},
		Peak: Peak{
			TierID:   mmr.Peak.Tier.ID,
			RankName: mmr.Peak.Tier.Name,
		},
		Sources: Sources{Account: "riot", MMR: "henrikdev"},
	}
	if len(mmr.Seasonal) > 0 {
		resp.Season = Season{
			Wins:        mmr.Seasonal[0].Wins,
			Games:       mmr.Seasonal[0].Games,
			EndRankName: mmr.Seasonal[0].EndTier.Name,
		}
	}
	return resp, nil
}

func mapUpstream(source string, err error) error {
	if errors.Is(err, riot.ErrNotFound) || errors.Is(err, henrik.ErrNotFound) {
		return &APIError{
			Status:  http.StatusNotFound,
			Code:    "NOT_FOUND",
			Message: "account not found",
		}
	}
	var up *upstream.Error
	if errors.As(err, &up) {
		if up.Status == http.StatusTooManyRequests {
			return &APIError{
				Status:     http.StatusTooManyRequests,
				Code:       "UPSTREAM",
				Message:    source + " rate limit exceeded",
				RetryAfter: up.RetryAfter,
			}
		}
		if up.Status == http.StatusBadRequest {
			return &APIError{
				Status:  http.StatusBadRequest,
				Code:    "INVALID_PARAMS",
				Message: source + " rejected the request",
			}
		}
		return &APIError{
			Status:  http.StatusBadGateway,
			Code:    "UPSTREAM",
			Message: fmt.Sprintf("%s returned status %d", source, up.Status),
		}
	}
	if strings.HasPrefix(err.Error(), "unsupported region:") {
		return &APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_PARAMS",
			Message: err.Error(),
		}
	}
	return &APIError{
		Status:  http.StatusBadGateway,
		Code:    "UPSTREAM",
		Message: source + " request failed",
	}
}
