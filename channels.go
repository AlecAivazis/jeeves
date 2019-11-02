package main

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

// NewGuild is invoked when a guild is registered with the bot
func (b *JeevesBot) NewGuild(s *discordgo.Session, event *discordgo.GuildCreate) {
	// only register guilds we have access to
	if event.Guild.Unavailable {
		return
	}

	// add an entry in the database for the new guild
	b.Client.Guild.Create().
		SetID(event.Guild.ID).
		Save(context.Background())
}

// RegisterChannels allows Jeeves to work across many channels in a guild. In order to know where to look
// for which messages, the user must assign a role to a particular channel
func (b *JeevesBot) RegisterChannels(session *discordgo.Session, message *discordgo.MessageCreate) {
	//
}
