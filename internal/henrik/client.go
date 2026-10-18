// Package henrik talks to HenrikDev's unofficial Valorant MMR API.
package henrik

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/dennn/valorant-twitch-commands/internal/upstream"
)

var ErrNotFound = errors.New("henrik account not found")

type Tier struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Current struct {
	Tier       Tier `json:"tier"`
	RR         int  `json:"rr"`
	LastChange int  `json:"last_change"`
	Elo        int  `json:"elo"`
}

type Peak struct {
	Tier Tier `json:"tier"`
}

type Season struct {
	Wins    int  `json:"wins"`
	Games   int  `json:"games"`
	EndTier Tier `json:"end_tier"`
}

type MMR struct {
	Account struct {
		PUUID string `json:"puuid"`
		Name  string `json:"name"`
		Tag   string `json:"tag"`
	} `json:"account"`
	Peak     Peak     `json:"peak"`
	Current  Current  `json:"current"`
	Seasonal []Season `json:"seasonal"`
}

type Client struct {
	apiKey string
	http   *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		http:   &http.Client{Timeout: 10 * time.Second},
	}
}

// GetMMRv3 loads current rank, RR, peak, and seasonal record.
// platform is typically "pc". region is a Henrik affinity such as na, eu, ap, or kr.
func (c *Client) GetMMRv3(ctx context.Context, region, platform, name, tag string) (MMR, error) {
	if platform == "" {
		platform = "pc"
	}
	endpoint := fmt.Sprintf(
		"https://api.henrikdev.xyz/valorant/v3/mmr/%s/%s/%s/%s",
		url.PathEscape(region),
		url.PathEscape(platform),
		url.PathEscape(name),
		url.PathEscape(tag),
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return MMR{}, err
	}
	if c.apiKey != "" {
		req.Header.Set("Authorization", c.apiKey)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return MMR{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return MMR{}, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return MMR{}, ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return MMR{}, &upstream.Error{
			Source:     "henrik",
			Status:     resp.StatusCode,
			RetryAfter: resp.Header.Get("Retry-After"),
			Detail:     string(body),
		}
	}

	var payload struct {
		Data MMR `json:"data"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return MMR{}, err
	}
	return payload.Data, nil
}
