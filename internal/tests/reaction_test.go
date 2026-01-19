package tests

import (
	"testing"
)

func TestReactionAdd(t *testing.T) {
	given, when, then := NewReactionStage(t)

	given.
		a_channel().and().
		a_message().and().
		a_default_emoji("👍").and().
		a_user()

	when.
		the_user_adds_a_reaction()

	then.
		the_reaction_should_be_saved().and().
		the_reaction_should_be_marked_as_default()
}

func TestReactionAddCustomEmoji(t *testing.T) {
	given, when, then := NewReactionStage(t)

	given.
		a_channel().and().
		a_message().and().
		a_custom_emoji("good").and().
		a_user()

	when.
		the_user_adds_a_reaction()

	then.
		the_reaction_should_be_saved().and().
		the_reaction_should_not_be_marked_as_default()
}

func TestReactionRemove(t *testing.T) {
	given, when, then := NewReactionStage(t)

	given.
		a_channel().and().
		a_message().and().
		a_default_emoji("👍").and().
		a_user().and().
		the_user_adds_a_reaction().and().
		the_reaction_should_be_saved()

	when.
		the_user_removes_the_reaction()

	then.
		the_reaction_should_be_removed()
}
