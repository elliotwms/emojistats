package tests

import (
	"context"
	"strings"
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

type CommandStage struct {
	t           *testing.T
	session     *discordgo.Session
	require     *require.Assertions
	assert      *assert.Assertions
	snowflake   *snowflake.Node
	fakediscord *fakediscord.Client

	channel         *discordgo.Channel
	message         *discordgo.Message
	interaction     *discordgo.InteractionCreate
	emoji           string
	emojiForCommand string // Emoji in MessageFormat for command queries
	userID          string
}

func NewCommandStage(t *testing.T) (*CommandStage, *CommandStage, *CommandStage) {
	s := &CommandStage{
		t:           t,
		session:     session,
		require:     require.New(t),
		assert:      assert.New(t),
		snowflake:   node,
		fakediscord: fakediscord.NewClient(session.Token),
	}

	ctx, cancel := context.WithCancel(context.Background())

	c := emojistats.NewConfig(session, session.State.User.ID)
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

func (s *CommandStage) and() *CommandStage {
	return s
}

func (s *CommandStage) a_channel() *CommandStage {
	c, err := s.session.GuildChannelCreate(testGuildID, "test-channel", discordgo.ChannelTypeGuildText)
	s.require.NoError(err)

	s.t.Cleanup(func() {
		_, err = s.session.ChannelDelete(c.ID)
		s.assert.NoError(err)
	})

	s.channel = c

	return s
}

func (s *CommandStage) a_message() *CommandStage {
	m, err := s.session.ChannelMessageSend(s.channel.ID, "Test message")
	s.require.NoError(err)
	s.message = m

	return s
}

func (s *CommandStage) a_user() *CommandStage {
	s.userID = s.session.State.User.ID

	return s
}

func (s *CommandStage) a_custom_emoji(name string) *CommandStage {
	id := s.snowflake.Generate().String()
	s.emoji = name + ":" + id
	// fakediscord stores custom emojis as <:id:name> based on how it parses the name:id input
	s.emojiForCommand = "<:" + id + ":" + name + ">"

	return s
}

func (s *CommandStage) a_default_emoji(emoji string) *CommandStage {
	s.emoji = emoji
	s.emojiForCommand = emoji // Default emojis use the same format

	return s
}

func (s *CommandStage) the_user_adds_a_reaction() *CommandStage {
	err := s.session.MessageReactionAdd(s.channel.ID, s.message.ID, s.emoji)
	s.require.NoError(err)

	// Wait for reaction to be saved
	s.require.Eventually(func() bool {
		var count int
		err := db.QueryRow(`SELECT COUNT(*) FROM reactions WHERE message_id = $1`, s.message.ID).Scan(&count)
		return err == nil && count > 0
	}, 5*time.Second, 100*time.Millisecond)

	return s
}

func (s *CommandStage) the_stats_command_is_invoked() *CommandStage {
	return s.the_stats_command_is_invoked_with_public(false)
}

func (s *CommandStage) the_stats_command_is_invoked_with_public(public bool) *CommandStage {
	options := []*discordgo.ApplicationCommandInteractionDataOption{}
	if public {
		options = append(options, &discordgo.ApplicationCommandInteractionDataOption{
			Name:  "public",
			Type:  discordgo.ApplicationCommandOptionBoolean,
			Value: true,
		})
	}

	i := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			ID:    s.snowflake.Generate().String(),
			AppID: s.session.State.User.ID,
			Type:  discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				ID:          s.snowflake.Generate().String(),
				Name:        "stats",
				CommandType: discordgo.ChatApplicationCommand,
				Options:     options,
			},
			GuildID:   testGuildID,
			ChannelID: s.channel.ID,
			Member: &discordgo.Member{
				User: &discordgo.User{
					ID: s.userID,
				},
			},
			Version: 1,
		},
	}

	var err error
	s.interaction, err = s.fakediscord.Interaction(i)
	s.require.NoError(err)
	s.require.NotEmpty(s.interaction)

	return s
}

func (s *CommandStage) the_emoji_stats_command_is_invoked_with_emoji(emoji string) *CommandStage {
	i := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			ID:    s.snowflake.Generate().String(),
			AppID: s.session.State.User.ID,
			Type:  discordgo.InteractionApplicationCommand,
			Data: discordgo.ApplicationCommandInteractionData{
				ID:          s.snowflake.Generate().String(),
				Name:        "emoji-stats",
				CommandType: discordgo.ChatApplicationCommand,
				Options: []*discordgo.ApplicationCommandInteractionDataOption{
					{
						Name:  "emoji",
						Type:  discordgo.ApplicationCommandOptionString,
						Value: emoji,
					},
				},
			},
			GuildID:   testGuildID,
			ChannelID: s.channel.ID,
			Member: &discordgo.Member{
				User: &discordgo.User{
					ID: s.userID,
				},
			},
			Version: 1,
		},
	}

	var err error
	s.interaction, err = s.fakediscord.Interaction(i)
	s.require.NoError(err)
	s.require.NotEmpty(s.interaction)

	return s
}

func (s *CommandStage) the_response_should_contain(text string) *CommandStage {
	s.require.Eventually(func() bool {
		res, err := s.session.InteractionResponse(s.interaction.Interaction)
		if err != nil {
			return false
		}

		return strings.Contains(res.Content, text)
	}, 5*time.Second, 100*time.Millisecond)

	return s
}

func (s *CommandStage) the_response_should_be_public() *CommandStage {
	s.require.Eventually(func() bool {
		res, err := s.session.InteractionResponse(s.interaction.Interaction)
		if err != nil {
			return false
		}

		return res.Flags&discordgo.MessageFlagsEphemeral == 0
	}, 5*time.Second, 100*time.Millisecond)

	return s
}

func (s *CommandStage) cleanupReactions() {
	if s.message != nil {
		_, _ = db.Exec(`DELETE FROM reactions WHERE message_id = $1`, s.message.ID)
	}
}
