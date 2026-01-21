package commands

import (
	"database/sql"

	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/bot/interactions/router"
	"github.com/elliotwms/emojistats/internal/stats"
)

var (
	publicOption = &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionBoolean,
		Name:        "public",
		Description: "Make the response visible to everyone (default: private)",
		Required:    false,
	}

	statsCommand = &discordgo.ApplicationCommand{
		Type:        discordgo.ChatApplicationCommand,
		Name:        "stats",
		Description: "View reaction statistics for this server",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "start_date",
				Description: "Start date (YYYY-MM-DD format)",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "end_date",
				Description: "End date (YYYY-MM-DD format)",
				Required:    false,
			},
			publicOption,
		},
	}

	emojiStatsCommand = &discordgo.ApplicationCommand{
		Type:        discordgo.ChatApplicationCommand,
		Name:        "emoji-stats",
		Description: "View statistics for a specific emoji",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "emoji",
				Description: "The emoji to analyze",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "start_date",
				Description: "Start date (YYYY-MM-DD format)",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "end_date",
				Description: "End date (YYYY-MM-DD format)",
				Required:    false,
			},
			publicOption,
		},
	}
)

// Commands returns the application commands and their handlers
func Commands(db *sql.DB) map[*discordgo.ApplicationCommand]router.ApplicationCommandHandler {
	repo := stats.NewRepository(db)

	return map[*discordgo.ApplicationCommand]router.ApplicationCommandHandler{
		statsCommand:      NewStatsHandler(repo),
		emojiStatsCommand: NewEmojiStatsHandler(repo),
	}
}
