package eventhandlers

import (
	"database/sql"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func NewReactionAddHandler(db *sql.DB) func(*discordgo.Session, *discordgo.MessageReactionAdd) {
	return func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
		msg, err := s.ChannelMessage(r.ChannelID, r.MessageID)
		if err != nil {
			slog.Error("failed to get message", "error", err)
			return
		}

		id := r.Emoji.ID
		isDefault := false

		// default emojis don't have an id, but a name
		if id == "" {
			id = r.Emoji.Name
			isDefault = true
		}

		_, err = db.Exec(`
			INSERT INTO reactions (emoji_id, sender_user_id, receiver_user_id, channel_id, message_id, guild_id, is_default)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			id,
			r.UserID,
			msg.Author.ID,
			r.ChannelID,
			r.MessageID,
			r.GuildID,
			isDefault,
		)
		if err != nil {
			slog.Error("failed to insert reaction", "error", err)
			return
		}

		slog.Info("reaction saved",
			"emoji_id", r.Emoji.ID,
			"sender", r.UserID,
			"receiver", msg.Author.ID,
		)
	}
}
