package stats

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/elliotwms/emojistats/internal/database"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://emojistats:emojistats@localhost:5432/emojistats?sslmode=disable"
	}

	var err error
	testDB, err = database.Connect(dsn)
	if err != nil {
		panic(err)
	}

	if err := database.Migrate(testDB); err != nil {
		panic(err)
	}

	code := m.Run()

	_ = testDB.Close()
	os.Exit(code)
}

func setupTest(t *testing.T) (*Repository, string, func()) {
	t.Helper()

	guildID := "test-guild-" + time.Now().Format("20060102150405.000000000")

	cleanup := func() {
		_, _ = testDB.Exec("DELETE FROM reactions WHERE guild_id = $1", guildID)
	}

	return NewRepository(testDB), guildID, cleanup
}

func insertReaction(t *testing.T, guildID, emojiID, senderID, receiverID, channelID, messageID string, isDefault bool, createdAt time.Time) {
	t.Helper()
	_, err := testDB.Exec(`
		INSERT INTO reactions (guild_id, emoji_id, sender_user_id, receiver_user_id, channel_id, message_id, is_default, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		guildID, emojiID, senderID, receiverID, channelID, messageID, isDefault, createdAt)
	require.NoError(t, err)
}

func TestGetGuildStats_Empty(t *testing.T) {
	repo, guildID, cleanup := setupTest(t)
	defer cleanup()

	stats, err := repo.GetGuildStats(context.Background(), guildID, DateRange{})

	require.NoError(t, err)
	assert.Equal(t, 0, stats.TotalReactions)
	assert.Empty(t, stats.TopEmojis)
	assert.Empty(t, stats.TopSenders)
	assert.Empty(t, stats.TopReceivers)
}

func TestGetGuildStats_WithReactions(t *testing.T) {
	repo, guildID, cleanup := setupTest(t)
	defer cleanup()

	now := time.Now()
	insertReaction(t, guildID, "👍", "sender1", "receiver1", "chan1", "msg1", true, now)
	insertReaction(t, guildID, "👍", "sender1", "receiver2", "chan1", "msg2", true, now)
	insertReaction(t, guildID, "👍", "sender2", "receiver1", "chan1", "msg3", true, now)
	insertReaction(t, guildID, "❤️", "sender1", "receiver1", "chan1", "msg4", true, now)

	stats, err := repo.GetGuildStats(context.Background(), guildID, DateRange{})

	require.NoError(t, err)
	assert.Equal(t, 4, stats.TotalReactions)
	require.Len(t, stats.TopEmojis, 2)
	assert.Equal(t, "👍", stats.TopEmojis[0].EmojiID)
	assert.Equal(t, 3, stats.TopEmojis[0].Count)
	assert.Equal(t, "❤️", stats.TopEmojis[1].EmojiID)
	assert.Equal(t, 1, stats.TopEmojis[1].Count)

	require.Len(t, stats.TopSenders, 2)
	assert.Equal(t, "sender1", stats.TopSenders[0].UserID)
	assert.Equal(t, 3, stats.TopSenders[0].Count)

	require.Len(t, stats.TopReceivers, 2)
	assert.Equal(t, "receiver1", stats.TopReceivers[0].UserID)
	assert.Equal(t, 3, stats.TopReceivers[0].Count)
}

func TestGetGuildStats_DateRangeFilter(t *testing.T) {
	repo, guildID, cleanup := setupTest(t)
	defer cleanup()

	oldDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	newDate := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)

	insertReaction(t, guildID, "👍", "sender1", "receiver1", "chan1", "msg1", true, oldDate)
	insertReaction(t, guildID, "👍", "sender1", "receiver1", "chan1", "msg2", true, newDate)
	insertReaction(t, guildID, "👍", "sender1", "receiver1", "chan1", "msg3", true, newDate)

	startDate := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
	dateRange := DateRange{Start: &startDate, End: &endDate}

	stats, err := repo.GetGuildStats(context.Background(), guildID, dateRange)

	require.NoError(t, err)
	assert.Equal(t, 2, stats.TotalReactions)
}

func TestGetEmojiStats_Empty(t *testing.T) {
	repo, guildID, cleanup := setupTest(t)
	defer cleanup()

	stats, err := repo.GetEmojiStats(context.Background(), guildID, "👍", DateRange{})

	require.NoError(t, err)
	assert.Equal(t, 0, stats.TotalUses)
	assert.Empty(t, stats.TopMessages)
	assert.Empty(t, stats.TopSenders)
	assert.Empty(t, stats.TopReceivers)
}

func TestGetEmojiStats_WithReactions(t *testing.T) {
	repo, guildID, cleanup := setupTest(t)
	defer cleanup()

	now := time.Now()
	insertReaction(t, guildID, "👍", "sender1", "receiver1", "chan1", "msg1", true, now)
	insertReaction(t, guildID, "👍", "sender2", "receiver1", "chan1", "msg1", true, now)
	insertReaction(t, guildID, "👍", "sender1", "receiver2", "chan1", "msg2", true, now)
	insertReaction(t, guildID, "❤️", "sender1", "receiver1", "chan1", "msg3", true, now)

	stats, err := repo.GetEmojiStats(context.Background(), guildID, "👍", DateRange{})

	require.NoError(t, err)
	assert.Equal(t, 3, stats.TotalUses)
	assert.True(t, stats.IsDefault)

	require.Len(t, stats.TopMessages, 2)
	assert.Equal(t, "msg1", stats.TopMessages[0].MessageID)
	assert.Equal(t, 2, stats.TopMessages[0].Count)

	require.Len(t, stats.TopSenders, 2)
	assert.Equal(t, "sender1", stats.TopSenders[0].UserID)
	assert.Equal(t, 2, stats.TopSenders[0].Count)

	require.Len(t, stats.TopReceivers, 2)
	assert.Equal(t, "receiver1", stats.TopReceivers[0].UserID)
	assert.Equal(t, 2, stats.TopReceivers[0].Count)
}

func TestGetEmojiStats_CustomEmoji(t *testing.T) {
	repo, guildID, cleanup := setupTest(t)
	defer cleanup()

	now := time.Now()
	insertReaction(t, guildID, "pepe:123456789", "sender1", "receiver1", "chan1", "msg1", false, now)

	stats, err := repo.GetEmojiStats(context.Background(), guildID, "pepe:123456789", DateRange{})

	require.NoError(t, err)
	assert.Equal(t, 1, stats.TotalUses)
	assert.False(t, stats.IsDefault)
}

func TestGetEmojiStats_DateRangeFilter(t *testing.T) {
	repo, guildID, cleanup := setupTest(t)
	defer cleanup()

	oldDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	newDate := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)

	insertReaction(t, guildID, "👍", "sender1", "receiver1", "chan1", "msg1", true, oldDate)
	insertReaction(t, guildID, "👍", "sender1", "receiver1", "chan1", "msg2", true, newDate)

	startDate := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	dateRange := DateRange{Start: &startDate}

	stats, err := repo.GetEmojiStats(context.Background(), guildID, "👍", dateRange)

	require.NoError(t, err)
	assert.Equal(t, 1, stats.TotalUses)
}

func TestGuildIsolation(t *testing.T) {
	repo, guildID1, cleanup1 := setupTest(t)
	defer cleanup1()

	_, guildID2, cleanup2 := setupTest(t)
	defer cleanup2()

	now := time.Now()
	insertReaction(t, guildID1, "👍", "sender1", "receiver1", "chan1", "msg1", true, now)
	insertReaction(t, guildID2, "👍", "sender1", "receiver1", "chan1", "msg2", true, now)
	insertReaction(t, guildID2, "👍", "sender1", "receiver1", "chan1", "msg3", true, now)

	stats1, err := repo.GetGuildStats(context.Background(), guildID1, DateRange{})
	require.NoError(t, err)
	assert.Equal(t, 1, stats1.TotalReactions)

	stats2, err := repo.GetGuildStats(context.Background(), guildID2, DateRange{})
	require.NoError(t, err)
	assert.Equal(t, 2, stats2.TotalReactions)
}
