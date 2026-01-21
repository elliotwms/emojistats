package stats

import "time"

// DateRange represents an optional date range for filtering queries
type DateRange struct {
	Start *time.Time
	End   *time.Time
}

// EmojiCount represents an emoji and its usage count
type EmojiCount struct {
	EmojiID   string
	IsDefault bool
	Count     int
}

// UserCount represents a user and their reaction count
type UserCount struct {
	UserID string
	Count  int
}

// MessageCount represents a message and its reaction count
type MessageCount struct {
	MessageID string
	ChannelID string
	Count     int
}

// GuildStats contains aggregated stats for a guild
type GuildStats struct {
	TotalReactions int
	TopEmojis      []EmojiCount
	TopSenders     []UserCount
	TopReceivers   []UserCount
}

// EmojiStats contains detailed stats for a specific emoji
type EmojiStats struct {
	EmojiID     string
	IsDefault   bool
	TotalUses   int
	TopMessages []MessageCount
	TopSenders  []UserCount
	TopReceivers []UserCount
}
