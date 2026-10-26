package service

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// DefaultLine is the chat sentence used when a bot omits a template.
const DefaultLine = "{name}#{tag} → {rank} {rr} RR"

var tokenPattern = regexp.MustCompile(`\{([a-zA-Z]+)(?::(lower|upper|num))?\}`)

// LineData is the token set a chat template can print.
type LineData struct {
	Name         string
	Tag          string
	Rank         string
	RR           int
	LastRRChange int
	Elo          int
	Peak         string
	Wins         int
	Losses       int
	Games        int
	RRChange     int
}

// FormatLine fills a chat template. Unknown tokens are left in place.
// {wins} prints 3W, {wins:num} prints 3. :lower and :upper change case.
func FormatLine(template string, data LineData) string {
	if strings.TrimSpace(template) == "" {
		template = DefaultLine
	}
	values := map[string]string{
		"name":         data.Name,
		"tag":          data.Tag,
		"rank":         data.Rank,
		"rr":           strconv.Itoa(data.RR),
		"lastRrChange": signed(data.LastRRChange),
		"elo":          strconv.Itoa(data.Elo),
		"peak":         data.Peak,
		"wins":         strconv.Itoa(data.Wins),
		"losses":       strconv.Itoa(data.Losses),
		"games":        strconv.Itoa(data.Games),
		"rrChange":     signed(data.RRChange),
	}
	return tokenPattern.ReplaceAllStringFunc(template, func(match string) string {
		parts := tokenPattern.FindStringSubmatch(match)
		key, mod := parts[1], parts[2]
		value, ok := values[key]
		if !ok {
			return match
		}
		switch mod {
		case "num":
			return strings.Map(func(r rune) rune {
				if (r >= '0' && r <= '9') || r == '-' {
					return r
				}
				return -1
			}, value)
		case "lower":
			return strings.ToLower(value)
		case "upper":
			return strings.ToUpper(value)
		}
		if key == "wins" {
			return value + "W"
		}
		if key == "losses" {
			return value + "L"
		}
		return value
	})
}

func signed(value int) string {
	if value > 0 {
		return fmt.Sprintf("+%d", value)
	}
	return strconv.Itoa(value)
}
