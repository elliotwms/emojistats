-- +goose Up
-- Indexes for efficient stats queries
CREATE INDEX idx_reactions_guild_emoji ON reactions (guild_id, emoji_id);
CREATE INDEX idx_reactions_guild_sender ON reactions (guild_id, sender_user_id);
CREATE INDEX idx_reactions_guild_receiver ON reactions (guild_id, receiver_user_id);
CREATE INDEX idx_reactions_guild_created ON reactions (guild_id, created_at);
CREATE INDEX idx_reactions_emoji_message ON reactions (emoji_id, message_id);
CREATE INDEX idx_reactions_guild_channel ON reactions (guild_id, channel_id);

-- +goose Down
DROP INDEX idx_reactions_guild_emoji;
DROP INDEX idx_reactions_guild_sender;
DROP INDEX idx_reactions_guild_receiver;
DROP INDEX idx_reactions_guild_created;
DROP INDEX idx_reactions_emoji_message;
DROP INDEX idx_reactions_guild_channel;
