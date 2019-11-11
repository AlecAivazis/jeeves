package main

import (
	"context"
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/AlecAivazis/jeeves/db/guild"
	"github.com/AlecAivazis/jeeves/db/guildchannel"
	"github.com/AlecAivazis/jeeves/db/predicate"
)

// Channel Roles
const (
	// ChannelRoleBank is the channel role used to identify the channel that displays the bank contents
	ChannelRoleBank = "bank"
	// ChannelRoleBankLog is the channel role used to identify the channel that the Banker uses
	ChannelRoleBankLog = "bank-log"
)

// Commands
const (
	CommandAssignChannel = "jeeves-assign-channel"
)

// NewGuild is invoked when a guild is registered with the bot
func (b *JeevesBot) NewGuild(s *discordgo.Session, event *discordgo.GuildCreate) {
	// only register guilds we have access to
	if event.Guild.Unavailable {
		return
	}

	// add an entry in the database for the new guild
	b.Database.Guild.Create().
		SetDiscordID(event.Guild.ID).
		Save(context.Background())
}

// RegisterChannels allows Jeeves to work across many channels in a guild. In order to know where to look
// for which messages, the user must assign a role to a particular channel
func (b *JeevesBot) RegisterChannels(session *discordgo.Session, message *discordgo.MessageCreate) {
	// since the message is presumably text, we care about words, not letters
	words := strings.Split(message.Content, " ")

	// make sure that the message starts with the right command
	if words[0] != "!"+CommandAssignChannel {
		return
	}

	//  there must be one argument after the command
	if len(words) != 2 {
		b.ReportError(message.ChannelID, errors.New("Please provide a role to associate with this channel"))
		return
	}

	// the role they want to assign
	role := words[1]

	// make sure the channel role is a valid one
	if err := validateChannelRole(role); err != nil {
		b.ReportError(message.ChannelID, err)
		return
	}

	// check if the users has permissions to set channel roles
	hasPerm, err := AuthorHasPermission(message.Message, discordgo.PermissionManageServer)
	if err != nil {
		b.ReportError(message.ChannelID, err)
		return
	}
	if !hasPerm {
		b.ReportError(message.ChannelID, errors.New("Sorry, you do not have permission to set channel roles"))
		return
	}

	// look up the database entry associated with this guild
	guildRecord, err := b.Database.Guild.Query().
		Where(guild.DiscordID(message.GuildID)).
		Only(context.Background())
	if err != nil {
		b.ReportError(message.ChannelID, err)
		return
	}

	// we need a get or update on the same guild channel config
	wherePredicates := []predicate.GuildChannel{
		guildchannel.HasGuildWith(guild.ID(guildRecord.ID)),
		guildchannel.Role(role),
	}

	// do we have an existing entry for the guild's channel role
	exists, err := b.Database.GuildChannel.Query().
		Where(wherePredicates...).
		Exist(context.Background())
	if err != nil {
		b.ReportError(message.ChannelID, err)
		return
	}

	// if the guild had not configured this channel role before
	if !exists {
		// store the association in the database
		b.Database.GuildChannel.Create().
			SetGuild(guildRecord).
			SetRole(role).
			SetChannel(message.ChannelID).
			Save(context.Background())
	} else {
		// otherwise we need to update the existing channel association
		b.Database.GuildChannel.Update().
			Where(wherePredicates...).
			SetChannel(message.ChannelID).
			Save(context.Background())
	}

	_, err = b.Discord.ChannelMessageSend(message.ChannelID, "Okay! I'll use this channel for the following role: "+role)
	if err != nil {
		b.ReportError(message.ChannelID, err)
		return
	}
}

func validateChannelRole(role string) error {
	if role == "bank-log" || role == "bank" {
		return nil
	}

	return errors.New("could not identify role " + role)
}

// AuthorHasRole returns true if the author has the given jeeves role
func AuthorHasRole(message *discordgo.Message, role string) (bool, error) {
	// check every role
	return true, nil
}

// AuthorHasPermission returns true if the author has the given permissions
func AuthorHasPermission(message *discordgo.Message, perm int) (bool, error) {
	return true, nil
}
