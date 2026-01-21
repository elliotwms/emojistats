package stats

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatGuildStats(t *testing.T) {
	stats := &GuildStats{
		TotalReactions: 100,
		TopEmojis: []EmojiCount{
			{EmojiID: "👍", IsDefault: true, Count: 50},
			{EmojiID: "<:pepe:123456789>", IsDefault: false, Count: 30},
		},
		TopSenders: []UserCount{
			{UserID: "111", Count: 40},
			{UserID: "222", Count: 30},
		},
		TopReceivers: []UserCount{
			{UserID: "333", Count: 25},
		},
	}

	result := FormatGuildStats(stats, "guild123")

	assert.Contains(t, result, "## Reaction Statistics")
	assert.Contains(t, result, "**Total Reactions:** 100")
	assert.Contains(t, result, "### Top 10 Reactions")
	assert.Contains(t, result, "1. 👍 - 50")
	assert.Contains(t, result, "2. <:pepe:123456789> - 30")
	assert.Contains(t, result, "### Top 3 Reaction Givers")
	assert.Contains(t, result, "<@111>")
	assert.Contains(t, result, "### Top 3 Reaction Receivers")
	assert.Contains(t, result, "<@333>")
}

func TestFormatGuildStats_Empty(t *testing.T) {
	stats := &GuildStats{
		TotalReactions: 0,
	}

	result := FormatGuildStats(stats, "guild123")

	assert.Contains(t, result, "**Total Reactions:** 0")
	assert.NotContains(t, result, "### Top 10 Reactions")
	assert.NotContains(t, result, "### Top 3 Reaction Givers")
}

func TestFormatEmojiStats(t *testing.T) {
	stats := &EmojiStats{
		EmojiID:   "👍",
		IsDefault: true,
		TotalUses: 50,
		TopMessages: []MessageCount{
			{MessageID: "msg1", ChannelID: "chan1", Count: 10},
			{MessageID: "msg2", ChannelID: "chan2", Count: 5},
		},
		TopSenders: []UserCount{
			{UserID: "111", Count: 20},
		},
		TopReceivers: []UserCount{
			{UserID: "222", Count: 15},
		},
	}

	result := FormatEmojiStats(stats, "guild123")

	assert.Contains(t, result, "## 👍 Statistics")
	assert.Contains(t, result, "**Total Uses:** 50")
	assert.Contains(t, result, "### Top 10 Messages")
	assert.Contains(t, result, "https://discord.com/channels/guild123/chan1/msg1")
	assert.Contains(t, result, "### Top 10 Recipients")
	assert.Contains(t, result, "<@222>")
	assert.Contains(t, result, "### Top 10 Senders")
	assert.Contains(t, result, "<@111>")
}

func TestFormatEmojiStats_CustomEmoji(t *testing.T) {
	stats := &EmojiStats{
		EmojiID:   "<:pepe:123456789>",
		IsDefault: false,
		TotalUses: 25,
	}

	result := FormatEmojiStats(stats, "guild123")

	assert.Contains(t, result, "## <:pepe:123456789> Statistics")
}

func TestFormatEmoji(t *testing.T) {
	tests := []struct {
		name      string
		emojiID   string
		isDefault bool
		expected  string
	}{
		{"default emoji", "👍", true, "👍"},
		{"custom emoji", "<:pepe:123456789>", false, "<:pepe:123456789>"},
		{"animated emoji", "<a:dance:987654321>", false, "<a:dance:987654321>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatEmoji(tt.emojiID, tt.isDefault)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatMessageLink(t *testing.T) {
	link := formatMessageLink("guild1", "channel1", "message1")
	assert.Equal(t, "https://discord.com/channels/guild1/channel1/message1", link)
}

func TestFormatRank(t *testing.T) {
	tests := []struct {
		position int
		expected string
	}{
		{1, "🥇"},
		{2, "🥈"},
		{3, "🥉"},
		{4, "4️⃣"},
		{5, "5️⃣"},
		{6, "6️⃣"},
		{7, "7️⃣"},
		{8, "8️⃣"},
		{9, "9️⃣"},
		{10, "🔟"},
		{11, "11."},
	}

	for _, tt := range tests {
		result := formatRank(tt.position)
		assert.Equal(t, tt.expected, result)
	}
}

func TestFormatGuildStats_RankingOrder(t *testing.T) {
	stats := &GuildStats{
		TotalReactions: 100,
		TopEmojis: []EmojiCount{
			{EmojiID: "first", IsDefault: true, Count: 50},
			{EmojiID: "second", IsDefault: true, Count: 30},
			{EmojiID: "third", IsDefault: true, Count: 10},
		},
	}

	result := FormatGuildStats(stats, "guild123")

	firstIdx := strings.Index(result, "1. first")
	secondIdx := strings.Index(result, "2. second")
	thirdIdx := strings.Index(result, "3. third")

	assert.True(t, firstIdx < secondIdx)
	assert.True(t, secondIdx < thirdIdx)
}
