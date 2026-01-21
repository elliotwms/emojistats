package emojistats

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/bot"
	"github.com/elliotwms/bot/interactions/router"
	"github.com/elliotwms/emojistats/internal/commands"
	"github.com/elliotwms/emojistats/internal/eventhandlers"
)

const intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessageReactions

type Config struct {
	Session         *discordgo.Session
	ApplicationID   string
	HealthCheckAddr string
	GuildID         string
	Logger          *slog.Logger
	DB              *sql.DB
}

func NewConfig(s *discordgo.Session, appID string) Config {
	return Config{
		Session:       s,
		ApplicationID: appID,
	}
}

func Run(config Config, ctx context.Context) error {
	r := router.New(router.WithDeferredResponse(true))

	b := bot.
		New(config.ApplicationID, config.Session).
		WithLogger(config.Logger).
		WithIntents(intents).
		WithHandler(eventhandlers.Ready).
		WithHandler(eventhandlers.NewReactionAddHandler(config.DB)).
		WithHandler(eventhandlers.NewReactionRemoveHandler(config.DB)).
		WithRouter(r).
		WithApplicationCommands(commands.Commands(config.DB)).
		WithMigrationEnabled(true)

	if config.HealthCheckAddr != "" {
		b.WithHealthCheck(config.HealthCheckAddr)
	}

	if config.GuildID != "" {
		b.WithGuildID(config.GuildID)
	}

	return b.Build().Run(ctx)
}
