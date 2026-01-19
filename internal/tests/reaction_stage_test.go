package tests

import (
	"context"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/bwmarrin/snowflake"
	"github.com/elliotwms/emojistats/internal/emojistats"
	"github.com/elliotwms/fakediscord/pkg/fakediscord"
	"github.com/neilotoole/slogt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ReactionStage struct {
	t           *testing.T
	session     *discordgo.Session
	require     *require.Assertions
	assert      *assert.Assertions
	snowflake   *snowflake.Node
	fakediscord *fakediscord.Client

	channel *discordgo.Channel
	message *discordgo.Message
	emoji   string
	userID  string
}

func NewReactionStage(t *testing.T) (*ReactionStage, *ReactionStage, *ReactionStage) {
	s := &ReactionStage{
		t:           t,
		session:     session,
		require:     require.New(t),
		assert:      assert.New(t),
		snowflake:   node,
		fakediscord: fakediscord.NewClient(session.Token),
	}

	ctx, cancel := context.WithCancel(context.Background())

	c := emojistats.NewConfig(session, "appid")
	c.Logger = slogt.New(t)
	c.DB = db
	c.GuildID = testGuildID

	done := make(chan struct{})

	go func() {
		s.require.NoError(emojistats.Run(c, ctx))
		close(done)
	}()

	t.Cleanup(func() {
		cancel()
		<-done
		s.cleanupReactions()
	})

	return s, s, s
}

func (s *ReactionStage) and() *ReactionStage {
	return s
}

func (s *ReactionStage) a_channel() *ReactionStage {
	c, err := s.session.GuildChannelCreate(testGuildID, "test-channel", discordgo.ChannelTypeGuildText)
	s.require.NoError(err)

	s.t.Cleanup(func() {
		_, err = s.session.ChannelDelete(c.ID)
		s.assert.NoError(err)
	})

	s.channel = c

	return s
}

func (s *ReactionStage) a_message() *ReactionStage {
	m, err := s.session.ChannelMessageSend(s.channel.ID, "Test message")
	s.require.NoError(err)
	s.message = m

	return s
}

func (s *ReactionStage) a_custom_emoji(name string) *ReactionStage {
	s.emoji = name + ":" + s.snowflake.Generate().String()

	return s
}

func (s *ReactionStage) a_default_emoji(emoji string) *ReactionStage {
	s.emoji = emoji

	return s
}

func (s *ReactionStage) a_user() *ReactionStage {
	// Use the session's user ID since reactions are added via the session
	s.userID = s.session.State.User.ID

	return s
}

func (s *ReactionStage) the_user_adds_a_reaction() *ReactionStage {
	err := s.session.MessageReactionAdd(s.channel.ID, s.message.ID, s.emoji)
	s.require.NoError(err)

	return s
}

func (s *ReactionStage) the_user_removes_the_reaction() *ReactionStage {
	err := s.session.MessageReactionRemove(s.channel.ID, s.message.ID, s.emoji, s.session.State.User.ID)
	s.require.NoError(err)

	return s
}

func (s *ReactionStage) the_reaction_should_be_saved() *ReactionStage {
	s.require.Eventually(func() bool {
		var count int
		err := db.QueryRow(`
			SELECT COUNT(*) FROM reactions
			WHERE message_id = $1 AND sender_user_id = $2`,
			s.message.ID, s.userID,
		).Scan(&count)

		return err == nil && count == 1
	}, 5*time.Second, 100*time.Millisecond)

	return s
}

func (s *ReactionStage) the_reaction_should_be_removed() *ReactionStage {
	s.require.Eventually(func() bool {
		var count int
		err := db.QueryRow(`
			SELECT COUNT(*) FROM reactions
			WHERE message_id = $1 AND sender_user_id = $2`,
			s.message.ID, s.userID,
		).Scan(&count)

		return err == nil && count == 0
	}, 5*time.Second, 100*time.Millisecond)

	return s
}

func (s *ReactionStage) the_reaction_should_be_marked_as_default() *ReactionStage {
	s.require.Eventually(func() bool {
		var isDefault bool
		err := db.QueryRow(`
			SELECT is_default FROM reactions
			WHERE message_id = $1 AND sender_user_id = $2`,
			s.message.ID, s.userID,
		).Scan(&isDefault)

		return err == nil && isDefault
	}, 5*time.Second, 100*time.Millisecond)

	return s
}

func (s *ReactionStage) the_reaction_should_not_be_marked_as_default() *ReactionStage {
	s.require.Eventually(func() bool {
		var isDefault bool
		err := db.QueryRow(`
			SELECT is_default FROM reactions
			WHERE message_id = $1 AND sender_user_id = $2`,
			s.message.ID, s.userID,
		).Scan(&isDefault)

		return err == nil && !isDefault
	}, 5*time.Second, 100*time.Millisecond)

	return s
}

func (s *ReactionStage) cleanupReactions() {
	if s.message != nil {
		_, _ = db.Exec(`DELETE FROM reactions WHERE message_id = $1`, s.message.ID)
	}
}
