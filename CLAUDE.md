# emojistats

Discord bot that tracks emoji usage statistics.

## Commands

```bash
# Build
go build -o emojistats ./cmd

# Run
go run ./cmd

# Test
go test ./...

# Start PostgreSQL
docker compose up -d
```

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| TOKEN | Yes | Discord bot token |
| APPLICATION_ID | Yes | Discord application ID |
| GUILD_ID | No | Restrict bot to specific guild |
| HEALTH_CHECK_ADDR | No | Health check endpoint address |
| LOG_LEVEL | No | Logging level |

## Architecture

Uses [elliotwms/bot](https://github.com/elliotwms/bot) framework which wraps discordgo. The pattern:

1. Configure bot in `internal/emojistats/` with `bot.Configure()`
2. Register event handlers from `internal/eventhandlers/` using `bot.Handle()`
3. Event handlers receive discordgo session and event, return error

## Database

PostgreSQL 17 via Docker Compose. Start with `docker compose up -d`.
