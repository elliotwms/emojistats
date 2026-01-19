package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/emojistats/internal/database"
	"github.com/elliotwms/emojistats/internal/emojistats"
)

func main() {
	logLevel := getLogLevel(os.Getenv("LOG_LEVEL"))
	slog.SetLogLoggerLevel(logLevel)

	db, err := database.Connect(mustGetEnv("DATABASE_URL"))
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			slog.Error("failed to close database", "error", err)
		}
	}()

	if err := database.Migrate(db); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	s := buildSession(logLevel)

	c := emojistats.NewConfig(s, mustGetEnv("APPLICATION_ID"))
	c.HealthCheckAddr = os.Getenv("HEALTH_CHECK_ADDR")
	c.GuildID = os.Getenv("GUILD_ID")
	c.Logger = slog.Default()
	c.DB = db

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	if err := emojistats.Run(c, ctx); err != nil {
		slog.Error("completed with error", "error", err)
		os.Exit(1)
	}
}

func getLogLevel(s string) (l slog.Level) {
	if s == "" {
		return slog.LevelInfo
	}

	if err := l.UnmarshalText([]byte(s)); err != nil {
		panic(err)
	}

	return l
}

func buildSession(l slog.Level) *discordgo.Session {
	s, err := discordgo.New("Bot " + mustGetEnv("TOKEN"))
	if err != nil {
		panic(err)
	}

	if l <= slog.LevelDebug {
		s.LogLevel = discordgo.LogDebug
	}
	return s
}

func mustGetEnv(s string) string {
	token := os.Getenv(s)
	if token == "" {
		panic(fmt.Sprintf("Missing '%s'", s))
	}
	return token
}
