package tests

import (
	"testing"
)

func TestStatsCommand(t *testing.T) {
	given, when, then := NewCommandStage(t)

	given.
		a_channel().and().
		a_message().and().
		a_default_emoji("👍").and().
		a_user().and().
		the_user_adds_a_reaction()

	when.
		the_stats_command_is_invoked()

	then.
		the_response_should_contain("## Reaction Statistics").and().
		the_response_should_contain("**Total Reactions:**").and().
		the_response_should_contain("### Top 10 Reactions").and().
		the_response_should_contain("👍")
}

func TestStatsCommandWithNoReactions(t *testing.T) {
	given, when, then := NewCommandStage(t)

	given.
		a_channel().and().
		a_user()

	when.
		the_stats_command_is_invoked()

	then.
		the_response_should_contain("## Reaction Statistics").and().
		the_response_should_contain("**Total Reactions:** 0")
}

func TestEmojiStatsCommand(t *testing.T) {
	given, when, then := NewCommandStage(t)

	given.
		a_channel().and().
		a_message().and().
		a_default_emoji("👍").and().
		a_user().and().
		the_user_adds_a_reaction()

	when.
		the_emoji_stats_command_is_invoked_with_emoji("👍")

	then.
		the_response_should_contain("## 👍 Statistics").and().
		the_response_should_contain("**Total Uses:**")
}

func TestEmojiStatsCommandWithCustomEmoji(t *testing.T) {
	given, when, then := NewCommandStage(t)

	given.
		a_channel().and().
		a_message().and().
		a_custom_emoji("custom").and().
		a_user().and().
		the_user_adds_a_reaction()

	when.
		the_emoji_stats_command_is_invoked_with_emoji(given.emojiForCommand)

	then.
		the_response_should_contain("Statistics").and().
		the_response_should_contain("**Total Uses:**")
}

func TestEmojiStatsCommandWithNoReactions(t *testing.T) {
	given, when, then := NewCommandStage(t)

	given.
		a_channel().and().
		a_user()

	when.
		the_emoji_stats_command_is_invoked_with_emoji("🎉")

	then.
		the_response_should_contain("No reactions found for this emoji.")
}
