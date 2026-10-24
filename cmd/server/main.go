package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dennn/valorant-twitch-commands/internal/config"
	"github.com/dennn/valorant-twitch-commands/internal/henrik"
	"github.com/dennn/valorant-twitch-commands/internal/httpapi"
	"github.com/dennn/valorant-twitch-commands/internal/riot"
	"github.com/dennn/valorant-twitch-commands/internal/service"
)

func main() {
	cfg := config.Load()
	riotClient := riot.NewClient(cfg.RiotAPIKey)
	henrikClient := henrik.NewClient(cfg.HenrikAPIKey)
	rank := service.NewRank(riotClient, henrikClient)
	record := service.NewRecord(henrikClient, time.Now())
	handler := httpapi.NewRouter(cfg, rank, record)

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}
