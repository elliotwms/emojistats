package commands

import (
	"context"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/bot/interactions/router"
	"github.com/elliotwms/emojistats/internal/stats"
)

// NewEmojiStatsHandler creates a handler for the /emoji-stats command
func NewEmojiStatsHandler(repo *stats.Repository) router.ApplicationCommandHandler {
	return func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
		public := parsePublicOption(data.Options)
		if err := deferResponse(s, i, public); err != nil {
			return err
		}

		guildID := i.GuildID

		var emojiID string
		for _, opt := range data.Options {
			if opt.Name == "emoji" {
				emojiID = opt.StringValue()
				break
			}
		}

		if emojiID == "" {
			return respondWithError(s, i, "Please provide an emoji.")
		}

		dateRange, err := parseDateRange(data.Options)
		if err != nil {
			return respondWithError(s, i, "Invalid date format. Please use YYYY-MM-DD.")
		}

		emojiStats, err := repo.GetEmojiStats(ctx, guildID, emojiID, dateRange)
		if err != nil {
			slog.Error("failed to get emoji stats", "error", err, "guild_id", guildID, "emoji_id", emojiID)
			return respondWithError(s, i, "Failed to retrieve emoji statistics.")
		}

		if emojiStats.TotalUses == 0 {
			return respondWithError(s, i, "No reactions found for this emoji.")
		}

		content := stats.FormatEmojiStats(emojiStats, guildID)
		return respond(s, i, content)
	}
}
