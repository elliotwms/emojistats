package tests

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/elliotwms/emojistats/internal/database"
	"github.com/elliotwms/fakediscord/pkg/fakediscord"
)

const testGuildName = "Emojistats Integration Testing"

var (
	session     *discordgo.Session
	testGuildID string
	db          *sql.DB
)

var node *snowflake.Node

func TestMain(m *testing.M) {
	fakediscord.Configure("http://localhost:8080/")

	if os.Getenv("TEST_DEBUG") != "" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	node, _ = snowflake.NewNode(0)

	openSession("bot")
	connectDB()

	code := m.Run()

	// todo should guild be deleted?
	//_ = session.GuildDelete(testGuildID)
	closeSession()
	closeDB()

	os.Exit(code)
}

func openSession(token string) {
	var err error
	session, err = discordgo.New(fmt.Sprintf("Bot %s", token))
	if err != nil {
		panic(err)
	}

	if os.Getenv("TEST_DEBUG") != "" {
		session.LogLevel = discordgo.LogDebug
		session.Debug = true
	}

	if err := session.Open(); err != nil {
		panic(err)
	}

	createGuild()
}

func createGuild() {
	guild, err := session.GuildCreate(testGuildName)
	if err != nil {
		panic(err)
	}
	testGuildID = guild.ID
}

func closeSession() {
	if err := session.Close(); err != nil {
		panic(err)
	}
}

func connectDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://emojistats:emojistats@localhost:5432/emojistats?sslmode=disable"
	}

	var err error
	db, err = database.Connect(dsn)
	if err != nil {
		panic(err)
	}

	if err := database.Migrate(db); err != nil {
		panic(err)
	}
}

func closeDB() {
	if err := db.Close(); err != nil {
		panic(err)
	}
}
