package eventhandlers

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

func Ready(s *discordgo.Session, _ *discordgo.Ready) {
	slog.Info("emojistats is ready")
}
