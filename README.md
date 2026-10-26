# valorant-twitch-commands

A small Valorant chat command. Viewers type `!rank` or `!record`. Nightbot, Fossabot, or StreamElements fetches one plain-text line: current rank, RR, peak, and the wins and losses since the stream session started.

The Go server resolves `Name#Tag` with the [Riot Account API](https://developer.riotgames.com/), then reads MMR from [HenrikDev](https://docs.henrikdev.xyz/). A React app (Vite) builds the bot command so you can paste it once.

Try it without installing anything: [valorant-twitch-commands on Cloud Run](https://valorant-twitch-commands-443346128252.us-central1.run.app/format).

## Run locally

You need Go 1.22, Node 22, a [Riot API key](https://developer.riotgames.com/), and a [HenrikDev API key](https://docs.henrikdev.xyz/).

```bash
cp .env.example .env
```

Fill in `RIOT_API_KEY` and `HENRIK_API_KEY`. `.env` is gitignored.

```bash
# API on :8080
go run ./cmd/server

# Frontend on :5173, in another terminal
cd frontend && npm ci && npm run dev
```

Open `http://localhost:5173`, enter a Riot ID, and copy the command for your bot. The frontend calls the API with `CORS_ALLOWED_ORIGIN` (default `http://localhost:5173`).

## API

All lookups take `name`, `tag`, and `region`. `platform` defaults to `pc`. Regions: `na`, `latam`, `br`, `eu`, `ap`, `kr`.

| Method | Path | Response |
| --- | --- | --- |
| `GET` | `/api/valorant/line` | Plain text for the bot. Optional `template`, at most 240 characters. |
| `GET` | `/api/valorant/rank` | JSON: tier, RR, last RR change, elo, peak, season wins and games. |
| `GET` | `/api/valorant/record` | JSON: wins, losses, and RR change since the session started. |
| `POST` | `/api/valorant/record/start` | Resets the session clock to now. |
| `GET` | `/healthz` | `ok` |

Default line: `{name}#{tag} → {rank} {rr} RR`.

```bash
curl "http://localhost:8080/api/valorant/line?name=TenZ&tag=SEN&region=na"
```

Tokens: `name`, `tag`, `rank`, `rr`, `lastRrChange`, `elo`, `peak`, `wins`, `losses`, `games`, `rrChange`. `{wins}` prints `3W`; `{wins:num}` prints `3`. `:lower` and `:upper` change case.

The session clock starts when the process starts. `POST /api/valorant/record/start` moves it to the current time. It is stored in memory and resets on restart.

## Docker

The image builds the frontend, compiles a static Go binary, and runs on `gcr.io/distroless/static-debian12`.

```bash
docker build -t valorant-twitch-commands .
docker run --rm -p 8080:8080 --env-file .env -e STATIC_DIR=/app/frontend/dist valorant-twitch-commands
```

The API and the built frontend are both served on port 8080. `STATIC_DIR` defaults to `frontend/dist` when that directory exists.

## Configuration

| Variable | Default | Purpose |
| --- | --- | --- |
| `RIOT_API_KEY` | | Riot `X-Riot-Token` |
| `HENRIK_API_KEY` | | HenrikDev key |
| `PORT` | `8080` | Listen port |
| `CORS_ALLOWED_ORIGIN` | `http://localhost:5173` | Allowed browser origin |
| `STATIC_DIR` | `frontend/dist` | Built frontend, if present |
