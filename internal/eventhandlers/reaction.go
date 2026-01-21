package eventhandlers

import (
	"database/sql"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func NewReactionAddHandler(db *sql.DB) func(*discordgo.Session, *discordgo.MessageReactionAdd) {
	return func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
		id := r.Emoji.MessageFormat()

		slog.Debug("reaction add event received",
			"emoji_id", id,
			"emoji_name", r.Emoji.Name,
			"user_id", r.UserID,
			"channel_id", r.ChannelID,
			"message_id", r.MessageID,
			"guild_id", r.GuildID,
		)

		msg, err := s.ChannelMessage(r.ChannelID, r.MessageID)
		if err != nil {
			slog.Error("failed to get message", "error", err)
			return
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
			r.Emoji.ID == "", // isDefault is reserved for future query use
		)
		if err != nil {
			slog.Error("failed to insert reaction", "error", err)
			return
		}

		slog.Info("reaction saved",
			"emoji_id", id,
			"sender", r.UserID,
			"receiver", msg.Author.ID,
		)
	}
}

func NewReactionRemoveHandler(db *sql.DB) func(*discordgo.Session, *discordgo.MessageReactionRemove) {
	return func(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
		id := r.Emoji.MessageFormat()

		slog.Debug("reaction remove event received",
			"emoji_id", id,
			"emoji_name", r.Emoji.Name,
			"user_id", r.UserID,
			"channel_id", r.ChannelID,
			"message_id", r.MessageID,
			"guild_id", r.GuildID,
		)

		res, err := db.Exec(`
			DELETE FROM reactions
			WHERE emoji_id = $1 AND sender_user_id = $2 AND message_id = $3`,
			id,
			r.UserID,
			r.MessageID,
		)
		if err != nil {
			slog.Error("failed to delete reaction", "error", err)
			return
		}

		rows, err := res.RowsAffected()
		if err != nil {
			slog.Error("failed to get rows affected", "error", err)
			return
		}

		if rows != 1 {
			slog.Warn("unexpected number of reactions deleted",
				"expected", 1,
				"actual", rows,
				"emoji_id", id,
				"sender", r.UserID,
				"message_id", r.MessageID,
			)
			return
		}

		slog.Info("reaction removed",
			"emoji_id", id,
			"sender", r.UserID,
		)
	}
}
