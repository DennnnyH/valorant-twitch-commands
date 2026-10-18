// Package riot talks to Riot's official account, content, match, and ranked APIs.
package riot

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

var ErrNotFound = errors.New("riot account not found")

type Account struct {
	PUUID    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
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

// GetAccountByRiotID resolves a Riot ID on the routing region that matches affinity
// (na/latam/br → americas, eu → europe, ap/kr → asia).
func (c *Client) GetAccountByRiotID(ctx context.Context, region, name, tag string) (Account, error) {
	routing, err := accountRouting(region)
	if err != nil {
		return Account{}, err
	}

	endpoint := fmt.Sprintf(
		"https://%s.api.riotgames.com/riot/account/v1/accounts/by-riot-id/%s/%s",
		routing,
		url.PathEscape(name),
		url.PathEscape(tag),
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return Account{}, err
	}
	req.Header.Set("X-Riot-Token", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return Account{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return Account{}, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return Account{}, ErrNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return Account{}, &upstream.Error{
			Source:     "riot",
			Status:     resp.StatusCode,
			RetryAfter: resp.Header.Get("Retry-After"),
			Detail:     string(body),
		}
	}

	var account Account
	if err := json.Unmarshal(body, &account); err != nil {
		return Account{}, err
	}
	return account, nil
}

func accountRouting(region string) (string, error) {
	switch region {
	case "na", "latam", "br":
		return "americas", nil
	case "eu":
		return "europe", nil
	case "ap", "kr":
		return "asia", nil
	default:
		return "", fmt.Errorf("unsupported region: %s", region)
	}
}
