package stats

import (
	"context"
	"database/sql"
	"strconv"
)

// Repository handles database queries for stats
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new stats repository
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// GetGuildStats retrieves aggregated stats for a guild
func (r *Repository) GetGuildStats(ctx context.Context, guildID string, dateRange DateRange) (*GuildStats, error) {
	stats := &GuildStats{}

	total, err := r.getTotalReactions(ctx, guildID, dateRange)
	if err != nil {
		return nil, err
	}
	stats.TotalReactions = total

	topEmojis, err := r.getTopEmojis(ctx, guildID, dateRange, 10)
	if err != nil {
		return nil, err
	}
	stats.TopEmojis = topEmojis

	topSenders, err := r.getTopSenders(ctx, guildID, "", dateRange, 3)
	if err != nil {
		return nil, err
	}
	stats.TopSenders = topSenders

	topReceivers, err := r.getTopReceivers(ctx, guildID, "", dateRange, 3)
	if err != nil {
		return nil, err
	}
	stats.TopReceivers = topReceivers

	return stats, nil
}

// GetEmojiStats retrieves detailed stats for a specific emoji
func (r *Repository) GetEmojiStats(ctx context.Context, guildID, emojiID string, dateRange DateRange) (*EmojiStats, error) {
	stats := &EmojiStats{
		EmojiID: emojiID,
	}

	total, isDefault, err := r.getEmojiTotalUses(ctx, guildID, emojiID, dateRange)
	if err != nil {
		return nil, err
	}
	stats.TotalUses = total
	stats.IsDefault = isDefault

	topMessages, err := r.getTopMessages(ctx, guildID, emojiID, dateRange, 10)
	if err != nil {
		return nil, err
	}
	stats.TopMessages = topMessages

	topSenders, err := r.getTopSenders(ctx, guildID, emojiID, dateRange, 10)
	if err != nil {
		return nil, err
	}
	stats.TopSenders = topSenders

	topReceivers, err := r.getTopReceivers(ctx, guildID, emojiID, dateRange, 10)
	if err != nil {
		return nil, err
	}
	stats.TopReceivers = topReceivers

	return stats, nil
}

func (r *Repository) getTotalReactions(ctx context.Context, guildID string, dateRange DateRange) (int, error) {
	query := `SELECT COUNT(*) FROM reactions WHERE guild_id = $1`
	args := []any{guildID}

	query, args = appendDateFilter(query, args, dateRange)

	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *Repository) getTopEmojis(ctx context.Context, guildID string, dateRange DateRange, limit int) ([]EmojiCount, error) {
	query := `
		SELECT emoji_id, is_default, COUNT(*) as count
		FROM reactions
		WHERE guild_id = $1`
	args := []any{guildID}

	query, args = appendDateFilter(query, args, dateRange)
	query += ` GROUP BY emoji_id, is_default ORDER BY count DESC LIMIT $` + argNum(len(args)+1)
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []EmojiCount
	for rows.Next() {
		var ec EmojiCount
		if err := rows.Scan(&ec.EmojiID, &ec.IsDefault, &ec.Count); err != nil {
			return nil, err
		}
		results = append(results, ec)
	}
	return results, rows.Err()
}

func (r *Repository) getTopSenders(ctx context.Context, guildID, emojiID string, dateRange DateRange, limit int) ([]UserCount, error) {
	query := `
		SELECT sender_user_id, COUNT(*) as count
		FROM reactions
		WHERE guild_id = $1`
	args := []any{guildID}

	if emojiID != "" {
		args = append(args, emojiID)
		query += ` AND emoji_id = $` + argNum(len(args))
	}

	query, args = appendDateFilter(query, args, dateRange)
	query += ` GROUP BY sender_user_id ORDER BY count DESC LIMIT $` + argNum(len(args)+1)
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []UserCount
	for rows.Next() {
		var uc UserCount
		if err := rows.Scan(&uc.UserID, &uc.Count); err != nil {
			return nil, err
		}
		results = append(results, uc)
	}
	return results, rows.Err()
}

func (r *Repository) getTopReceivers(ctx context.Context, guildID, emojiID string, dateRange DateRange, limit int) ([]UserCount, error) {
	query := `
		SELECT receiver_user_id, COUNT(*) as count
		FROM reactions
		WHERE guild_id = $1`
	args := []any{guildID}

	if emojiID != "" {
		args = append(args, emojiID)
		query += ` AND emoji_id = $` + argNum(len(args))
	}

	query, args = appendDateFilter(query, args, dateRange)
	query += ` GROUP BY receiver_user_id ORDER BY count DESC LIMIT $` + argNum(len(args)+1)
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []UserCount
	for rows.Next() {
		var uc UserCount
		if err := rows.Scan(&uc.UserID, &uc.Count); err != nil {
			return nil, err
		}
		results = append(results, uc)
	}
	return results, rows.Err()
}

func (r *Repository) getTopMessages(ctx context.Context, guildID, emojiID string, dateRange DateRange, limit int) ([]MessageCount, error) {
	query := `
		SELECT message_id, channel_id, COUNT(*) as count
		FROM reactions
		WHERE guild_id = $1 AND emoji_id = $2`
	args := []any{guildID, emojiID}

	query, args = appendDateFilter(query, args, dateRange)
	query += ` GROUP BY message_id, channel_id ORDER BY count DESC LIMIT $` + argNum(len(args)+1)
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []MessageCount
	for rows.Next() {
		var mc MessageCount
		if err := rows.Scan(&mc.MessageID, &mc.ChannelID, &mc.Count); err != nil {
			return nil, err
		}
		results = append(results, mc)
	}
	return results, rows.Err()
}

func (r *Repository) getEmojiTotalUses(ctx context.Context, guildID, emojiID string, dateRange DateRange) (int, bool, error) {
	query := `SELECT COUNT(*), COALESCE(bool_or(is_default), false) FROM reactions WHERE guild_id = $1 AND emoji_id = $2`
	args := []any{guildID, emojiID}

	query, args = appendDateFilter(query, args, dateRange)

	var count int
	var isDefault bool
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count, &isDefault)
	return count, isDefault, err
}

func appendDateFilter(query string, args []any, dateRange DateRange) (string, []any) {
	if dateRange.Start != nil {
		args = append(args, *dateRange.Start)
		query += ` AND created_at >= $` + argNum(len(args))
	}
	if dateRange.End != nil {
		args = append(args, *dateRange.End)
		query += ` AND created_at < $` + argNum(len(args))
	}
	return query, args
}

func argNum(n int) string {
	return strconv.Itoa(n)
}
