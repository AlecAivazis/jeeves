package bot

import (
	"context"
	"errors"

	"github.com/bwmarrin/discordgo"

	"github.com/AlecAivazis/jeeves/db/guild"
	"github.com/AlecAivazis/jeeves/db/guildchannel"
	"github.com/AlecAivazis/jeeves/db/predicate"
)

// Channel Roles
const (
	// ChannelRoleBank is the channel role used to identify the channel that displays the bank contents
	ChannelRoleBank = "bank"
)

// Commands
const (
	CommandAssignChannel = "jeeves-assign-channel"
)

// RegisterChannelRole allows Jeeves to work across many channels in a guild. In order to know where to look
// for which messages, the user must assign a role to a particular channel
func (b *JeevesBot) RegisterChannelRole(ctx *CommandContext, role string) error {
	// make sure the channel role is a valid one
	if err := validateChannelRole(role); err != nil {
		return err
	}

	// check if the users has permissions to set channel roles
	hasPerm, err := AuthorHasPermission(ctx, discordgo.PermissionManageServer)
	if err != nil {
		return err
	}
	if !hasPerm {
		return errors.New("Sorry, you do not have permission to set channel roles")
	}

	// look up the database entry associated with this guild
	guildRecord, err := b.Database.Guild.Query().
		Where(guild.DiscordID(ctx.GuildID)).
		Only(context.Background())
	if err != nil {
		return err
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
		return err
	}

	// if the guild had not configured this channel role before
	if !exists {
		// store the association in the database
		b.Database.GuildChannel.Create().
			SetGuild(guildRecord).
			SetRole(role).
			SetChannel(ctx.ChannelID).
			Save(context.Background())
	} else {
		// otherwise we need to update the existing channel association
		b.Database.GuildChannel.Update().
			Where(wherePredicates...).
			SetChannel(ctx.ChannelID).
			Save(context.Background())
	}

	_, err = b.Discord.ChannelMessageSend(ctx.ChannelID, "Okay! I'll use this channel for the following role: "+role)
	if err != nil {
		return err
	}

	return nil
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
func AuthorHasPermission(message *CommandContext, perm int) (bool, error) {
	return true, nil
}
