package commands

import (
	"context"
	"log/slog"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/bot/interactions/router"
	"github.com/elliotwms/emojistats/internal/stats"
)

// NewStatsHandler creates a handler for the /stats command
func NewStatsHandler(repo *stats.Repository) router.ApplicationCommandHandler {
	return func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
		guildID := i.GuildID

		dateRange, err := parseDateRange(data.Options)
		if err != nil {
			return respondWithError(s, i, "Invalid date format. Please use YYYY-MM-DD.")
		}

		guildStats, err := repo.GetGuildStats(ctx, guildID, dateRange)
		if err != nil {
			slog.Error("failed to get guild stats", "error", err, "guild_id", guildID)
			return respondWithError(s, i, "Failed to retrieve statistics.")
		}

		content := stats.FormatGuildStats(guildStats, guildID)
		return respond(s, i, content)
	}
}

func parseDateRange(options []*discordgo.ApplicationCommandInteractionDataOption) (stats.DateRange, error) {
	var dateRange stats.DateRange

	for _, opt := range options {
		switch opt.Name {
		case "start_date":
			t, err := time.Parse("2006-01-02", opt.StringValue())
			if err != nil {
				return dateRange, err
			}
			dateRange.Start = &t
		case "end_date":
			t, err := time.Parse("2006-01-02", opt.StringValue())
			if err != nil {
				return dateRange, err
			}
			// Add 1 day to make end_date inclusive
			t = t.AddDate(0, 0, 1)
			dateRange.End = &t
		}
	}

	return dateRange, nil
}

func respond(s *discordgo.Session, i *discordgo.InteractionCreate, content string) error {
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
	return err
}

func respondWithError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) error {
	return respond(s, i, message)
}
