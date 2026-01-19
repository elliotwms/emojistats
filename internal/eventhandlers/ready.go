package eventhandlers

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/emojistats/internal/build"
)

func Ready(s *discordgo.Session, _ *discordgo.Ready) {
	slog.Info("emojistats is ready", "version", build.Version)

	err := s.UpdateGameStatus(0, build.Version)
	if err != nil {
		slog.Error("Could not update game status", "error", err)
		return
	}
}
