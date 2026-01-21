package commands

import (
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDateRange_Empty(t *testing.T) {
	options := []*discordgo.ApplicationCommandInteractionDataOption{}

	dateRange, err := parseDateRange(options)

	require.NoError(t, err)
	assert.Nil(t, dateRange.Start)
	assert.Nil(t, dateRange.End)
}

func TestParseDateRange_StartDateOnly(t *testing.T) {
	options := []*discordgo.ApplicationCommandInteractionDataOption{
		{Name: "start_date", Type: discordgo.ApplicationCommandOptionString, Value: "2024-01-15"},
	}

	dateRange, err := parseDateRange(options)

	require.NoError(t, err)
	require.NotNil(t, dateRange.Start)
	assert.Equal(t, 2024, dateRange.Start.Year())
	assert.Equal(t, time.January, dateRange.Start.Month())
	assert.Equal(t, 15, dateRange.Start.Day())
	assert.Nil(t, dateRange.End)
}

func TestParseDateRange_EndDateOnly(t *testing.T) {
	options := []*discordgo.ApplicationCommandInteractionDataOption{
		{Name: "end_date", Type: discordgo.ApplicationCommandOptionString, Value: "2024-01-20"},
	}

	dateRange, err := parseDateRange(options)

	require.NoError(t, err)
	assert.Nil(t, dateRange.Start)
	require.NotNil(t, dateRange.End)
	// End date should be +1 day to make it inclusive
	assert.Equal(t, 2024, dateRange.End.Year())
	assert.Equal(t, time.January, dateRange.End.Month())
	assert.Equal(t, 21, dateRange.End.Day())
}

func TestParseDateRange_BothDates(t *testing.T) {
	options := []*discordgo.ApplicationCommandInteractionDataOption{
		{Name: "start_date", Type: discordgo.ApplicationCommandOptionString, Value: "2024-01-01"},
		{Name: "end_date", Type: discordgo.ApplicationCommandOptionString, Value: "2024-01-31"},
	}

	dateRange, err := parseDateRange(options)

	require.NoError(t, err)
	require.NotNil(t, dateRange.Start)
	require.NotNil(t, dateRange.End)
	assert.Equal(t, 1, dateRange.Start.Day())
	assert.Equal(t, 1, dateRange.End.Day())   // 31 + 1 = Feb 1
	assert.Equal(t, time.February, dateRange.End.Month())
}

func TestParseDateRange_InvalidStartDate(t *testing.T) {
	options := []*discordgo.ApplicationCommandInteractionDataOption{
		{Name: "start_date", Type: discordgo.ApplicationCommandOptionString, Value: "not-a-date"},
	}

	_, err := parseDateRange(options)

	assert.Error(t, err)
}

func TestParseDateRange_InvalidEndDate(t *testing.T) {
	options := []*discordgo.ApplicationCommandInteractionDataOption{
		{Name: "end_date", Type: discordgo.ApplicationCommandOptionString, Value: "01/20/2024"},
	}

	_, err := parseDateRange(options)

	assert.Error(t, err)
}

func TestParseDateRange_IgnoresUnknownOptions(t *testing.T) {
	options := []*discordgo.ApplicationCommandInteractionDataOption{
		{Name: "emoji", Type: discordgo.ApplicationCommandOptionString, Value: "👍"},
		{Name: "start_date", Type: discordgo.ApplicationCommandOptionString, Value: "2024-01-15"},
	}

	dateRange, err := parseDateRange(options)

	require.NoError(t, err)
	require.NotNil(t, dateRange.Start)
	assert.Equal(t, 15, dateRange.Start.Day())
}
