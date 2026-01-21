package stats

import (
	"fmt"
	"strings"
)

// FormatGuildStats formats guild stats as Discord markdown
func FormatGuildStats(stats *GuildStats, guildID string) string {
	var sb strings.Builder

	sb.WriteString("## Reaction Statistics\n\n")
	sb.WriteString(fmt.Sprintf("**Total Reactions:** %d\n\n", stats.TotalReactions))

	if len(stats.TopEmojis) > 0 {
		sb.WriteString("### Top 10 Reactions\n")
		for i, e := range stats.TopEmojis {
			sb.WriteString(fmt.Sprintf("%d. %s - %d\n", i+1, formatEmoji(e.EmojiID, e.IsDefault), e.Count))
		}
		sb.WriteString("\n")
	}

	if len(stats.TopSenders) > 0 {
		sb.WriteString("### Top 3 Reaction Givers\n")
		for i, u := range stats.TopSenders {
			sb.WriteString(fmt.Sprintf("%s <@%s> - %d\n", formatRank(i+1), u.UserID, u.Count))
		}
		sb.WriteString("\n")
	}

	if len(stats.TopReceivers) > 0 {
		sb.WriteString("### Top 3 Reaction Receivers\n")
		for i, u := range stats.TopReceivers {
			sb.WriteString(fmt.Sprintf("%s <@%s> - %d\n", formatRank(i+1), u.UserID, u.Count))
		}
	}

	return sb.String()
}

// FormatEmojiStats formats emoji-specific stats as Discord markdown
func FormatEmojiStats(stats *EmojiStats, guildID string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("## %s Statistics\n\n", formatEmoji(stats.EmojiID, stats.IsDefault)))
	sb.WriteString(fmt.Sprintf("**Total Uses:** %d\n\n", stats.TotalUses))

	if len(stats.TopMessages) > 0 {
		sb.WriteString("### Top 10 Messages\n")
		for i, m := range stats.TopMessages {
			link := formatMessageLink(guildID, m.ChannelID, m.MessageID)
			sb.WriteString(fmt.Sprintf("%s [Jump to message](%s) - %d\n", formatRank(i+1), link, m.Count))
		}
		sb.WriteString("\n")
	}

	if len(stats.TopReceivers) > 0 {
		sb.WriteString("### Top 10 Recipients\n")
		for i, u := range stats.TopReceivers {
			sb.WriteString(fmt.Sprintf("%s <@%s> - %d\n", formatRank(i+1), u.UserID, u.Count))
		}
		sb.WriteString("\n")
	}

	if len(stats.TopSenders) > 0 {
		sb.WriteString("### Top 10 Senders\n")
		for i, u := range stats.TopSenders {
			sb.WriteString(fmt.Sprintf("%s <@%s> - %d\n", formatRank(i+1), u.UserID, u.Count))
		}
	}

	return sb.String()
}

func formatEmoji(emojiID string, _ bool) string {
	return emojiID
}

func formatMessageLink(guildID, channelID, messageID string) string {
	return fmt.Sprintf("https://discord.com/channels/%s/%s/%s", guildID, channelID, messageID)
}

func formatRank(position int) string {
	switch position {
	case 1:
		return "🥇"
	case 2:
		return "🥈"
	case 3:
		return "🥉"
	case 4:
		return "4️⃣"
	case 5:
		return "5️⃣"
	case 6:
		return "6️⃣"
	case 7:
		return "7️⃣"
	case 8:
		return "8️⃣"
	case 9:
		return "9️⃣"
	case 10:
		return "🔟"
	default:
		return fmt.Sprintf("%d.", position)
	}
}
