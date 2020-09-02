// !deposit 1xA 2xB C - adds one A, two B, and one C to the guild's bank
// !withdraw 1xA 2xB C (for arcanite reaper) - removes one A, two B, and one C to the guild's bank with the provided note

package bank

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/AlecAivazis/jeeves/data"
	"github.com/AlecAivazis/jeeves/db"
	"github.com/AlecAivazis/jeeves/db/bankitem"
	"github.com/AlecAivazis/jeeves/db/guild"
	"github.com/AlecAivazis/jeeves/db/guildbank"
	"github.com/AlecAivazis/jeeves/db/predicate"
	"github.com/bwmarrin/discordgo"
)

const (
	// RoleBanker defines the public name of the role to give non-admin users permissions to modify the bank
	RoleBanker = "Banker"
)

const (
	// CommandDeposit defines the command used to deposit items into the guild bank
	CommandDeposit = "deposit"
	// CommandWithdraw defines the command used to withdraw items from the guild bank
	CommandWithdraw = "withdraw"
	// CommandRequest can be used to submit a request to be fulfilled by a banker
	CommandRequest = "request"
	// CommandAssignBankChannel defines the command used to assign a channel to use to display the bank
	CommandAssignBankChannel = "jeeves-assign-bank"
	// CommandRefreshBank can be used to force the bank to be re-rendered
	CommandRefreshBank = "refresh-bank"
	// CommandResetBank can be  used to reset the banks current inventory
	CommandResetBank = "reset-bank"
)

// CommandHandler handles the parsing and dispatching of commands for Jeeves
func (b *Banker) CommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	// only look at commands
	if message.Content[0] != '!' {
		return
	}

	// since the message is presumably text, we care about words, not letters
	words := strings.SplitN(message.Content[1:], " ", 2)
	command := words[0]

	// construct the context object
	ctx := &CommandContext{
		Banker:    b,
		GuildID:   message.GuildID,
		ChannelID: message.ChannelID,
		Context:   context.Background(),
		Message:   message.Message,
	}

	// split the items
	items := strings.Split(strings.Trim(words[1], ", "), ",")

	var err error
	var done = true
	// check the command against our known strings
	switch command {
	case CommandAssignBankChannel:
		done, err = b.InitializeBankChannel(ctx)
	case CommandDeposit:
		if len(words) == 1 {
			err = errors.New("you did not give me items to deposit")
		} else {
			done, err = b.DepositItems(ctx, items)
		}
	case CommandWithdraw:
		if len(words) == 1 {
			err = errors.New("you did not give me items to withdraw")
		} else {
			done, err = b.WithdrawItems(ctx, items)
		}
	case CommandRequest:
		if len(words) == 1 {
			err = errors.New("you did not give me items to request")
		} else {
			done, err = b.RequestItems(ctx, items)
		}
	case CommandResetBank:
		err = b.ResetBank(ctx)
	case CommandRefreshBank:
		err = b.UpdateBankListing(ctx)
	default:
		return
	}
	// if the command failed
	if err != nil {
		// send the error to the channel we received the message on
		b.Discord.ChannelMessageSend(message.ChannelID, "Sorry, "+err.Error())
		return
	}

	// if we are supposed to confirm
	if done {
		// confirm the action with a reaction
		err = b.Discord.MessageReactionAdd(message.ChannelID, message.ID, "👍")
		if err != nil {
			b.Discord.ChannelMessageSend(message.ChannelID, "Sorry, "+err.Error())
		}
	}
}

// InitializeBankChannel is called when the user intends to assign a channel for use to display the bank
func (b *Banker) InitializeBankChannel(ctx *CommandContext) (bool, error) {
	// confirm the action with the user
	_, err := b.Discord.ChannelMessageSend(ctx.ChannelID, "Okay! Please give me a moment to set up your guild bank...")
	if err != nil {
		return false, err
	}

	// if we haven't defined the banker role yet
	roles, err := b.Discord.GuildRoles(ctx.GuildID)
	if err != nil {
		return false, err
	}
	definedBanker := false
	for _, role := range roles {
		if role.Name == RoleBanker {
			definedBanker = true
			break
		}
	}
	// if we have to define the banker role now
	if !definedBanker {
		// tell the user about it
		_, err = b.Discord.ChannelMessageSend(ctx.ChannelID, "I am creating the Banker role. Assign this to non-Admin users you want"+
			" to give permission to move items in and out of the bank.")
		if err != nil {
			return false, err
		}

		// create the banker role
		role, err := b.Discord.GuildRoleCreate(ctx.GuildID)

		// jeeves might not have the permissions to make roles (that's a high level one) but that's okay!

		// if we succeed
		if err == nil {
			// edit the role we just made (not sure why we couldn't do this when we created it to begin with...)
			_, err = b.Discord.GuildRoleEdit(ctx.GuildID, role.ID, RoleBanker, role.Color, role.Hoist, role.Permissions, role.Mentionable)
			if err != nil {
				return false, err
			}

			// if there was a problem, we still want to continue
		} else {
			// tell them what happened
			b.Discord.ChannelMessageSend(ctx.ChannelID, "Hmmm something happened when I tried to do that. Maybe I don't have permissions to add roles?"+
				" Before you try again, either make sure I have that permission or you can create the \"Banker\" role and I"+
				" won't try to make it next time.")
		}

	}

	// send the display message now so they know what they can delete
	_, err = b.Discord.ChannelMessageSend(ctx.ChannelID, "All set! Your guild bank's contents will go here. You are free to"+
		" delete any message in this channel. I will update it as your bankers add items to the bank.")
	if err != nil {
		return false, err
	}

	// we need to find and update the bank record for this guild
	wherePredicates := []predicate.GuildBank{
		guildbank.HasGuildWith(guild.DiscordID(ctx.GuildID)),
	}

	// look if we have an existing record for the bank
	previousRecord, err := b.Database.GuildBank.Query().
		Where(wherePredicates...).
		All(ctx)
	if err != nil {
		return false, err
	}

	// if we have never recorded a bank for this guild
	if len(previousRecord) == 0 {
		// grab the guild from context
		guild, err := ctx.Guild()
		if err != nil {
			return false, err
		}

		// create the entry for the bank
		_, err = b.Database.GuildBank.Create().
			SetGuild(guild).
			SetChannelID(ctx.ChannelID).
			Save(ctx)
		if err != nil {
			return false, err
		}

	} else {
		// update the bank entry to have the new channel
		_, err = b.Database.GuildBank.Update().
			Where(wherePredicates...).
			SetChannelID(ctx.ChannelID).
			Save(ctx)
		if err != nil {
			return false, err
		}
	}

	// nothing went wrong
	return true, nil
}

// WithdrawItems is used when the user wants to withdraw the specified items from the bank. Will update the display message.
func (b *Banker) WithdrawItems(ctx *CommandContext, items []string) (bool, error) {
	// make sure the user has the right permissions
	if err := b.CheckInventory(ctx, items); err != nil {
		return false, err
	}

	// find the bank for this guild
	guildBank, err := ctx.GuildBank()
	if err != nil {
		return false, err
	}

	transactions, err := ParseTransactions(items)
	if err != nil {
		return false, err
	}

	// we need to add each item to the database
	for _, transaction := range transactions {
		// pull out the constants of the transaction
		item := transaction.Item
		amount := transaction.Amount

		// if we are depositing gold
		if item == data.ItemIDGold {
			// decrement the guild bance
			guildBank.Update().AddBalance(-amount).Exec(ctx)

			// we're done processing it
			continue
		}

		// does this bank have a record for the item
		existingItems, err := guildBank.
			QueryItems().
			Where(bankitem.ItemID(item)).
			All(ctx)
		if err != nil {
			return false, err
		}

		// if withdrawing this item will take its quantity to zero
		if amount == existingItems[0].Quantity {
			// just remove it from the database
			err = b.Database.BankItem.DeleteOneID(existingItems[0].ID).Exec(ctx)
		} else {
			// update the existing record
			err = existingItems[0].Update().
				AddQuantity(-amount).
				Exec(ctx)
		}

		if err != nil {
			return false, err
		}
	}

	// once we are done adding the items we should update the listing
	return true, b.UpdateBankListing(ctx)
}

// DepositItems is used when the user wants to deposit the specified items into the bank. Will update the display message.
func (b *Banker) DepositItems(ctx *CommandContext, items []string) (bool, error) {
	// make sure the user has the right permissions
	if err := b.userCanModifyBank(ctx, ctx.Message.Member); err != nil {
		return false, err
	}

	// find the bank for this guild
	guildBank, err := ctx.GuildBank()
	if err != nil {
		return false, err
	}

	transactions, err := ParseTransactions(items)
	if err != nil {
		return false, err
	}

	// we need to add each item to the database
	for _, transaction := range transactions {
		// pull out the constants of the transaction
		item := transaction.Item
		amount := transaction.Amount

		// if we are depositing gold
		if item == data.ItemIDGold {
			// add the deposit to the guild bank
			guildBank.Update().AddBalance(amount).Exec(ctx)

			// we're done processing it
			continue
		}

		// does this bank have a record for the item
		existingItems, err := guildBank.
			QueryItems().
			Where(bankitem.ItemID(item)).
			All(ctx)
		if err != nil {
			return false, err
		}

		// if we haven't seen the item before
		if len(existingItems) == 0 {
			// create a bank item entry
			_, err := b.Database.BankItem.Create().
				SetItemID(item).
				SetQuantity(amount).
				SetBank(guildBank).
				Save(ctx)
			if err != nil {
				return false, err
			}

			// we're done processing this item
			continue
		}

		// we are adding an item to an existing record in the bank
		err = b.Database.BankItem.Update().
			Where(bankitem.ID(existingItems[0].ID)).
			AddQuantity(amount).
			Exec(ctx)
		if err != nil {
			return false, err
		}
	}

	// once we are done adding the items we should update the listing
	return true, b.UpdateBankListing(ctx)
}

// RequestItems is called when a user wants something from the guild bank
func (b *Banker) RequestItems(ctx *CommandContext, items []string) (bool, error) {
	// make sure we would be able to perform the withdraw (ie, the item names are valid)
	if err := b.CheckInventory(ctx, items); err != nil {
		return false, err
	}

	// confirm that we see the withdraw
	b.Discord.MessageReactionAdd(ctx.ChannelID, ctx.Message.ID, "👀")

	// notify all of the bankers that there is a request

	// the message to send
	message := fmt.Sprintf(
		"Hey! It looks like %s wants to withdraw: %s",
		ctx.MemberName(ctx.Message.Author),
		strings.Join(items, ","),
	)

	// find the bankers in this guild
	bankers, err := b.Bankers(ctx)
	if err != nil {
		return false, err
	}
	for _, banker := range bankers {
		// grab a reference to the channel with the user
		channel, err := b.Discord.UserChannelCreate(banker.User.ID)
		if err != nil {
			return false, err
		}

		// send a notification on the channel
		_, err = b.Discord.ChannelMessageSend(channel.ID, message)
		if err != nil {
			return false, err
		}
	}

	return false, nil
}

// ResetBank should be called when the bank contents for this guild should be reset
func (b *Banker) ResetBank(ctx *CommandContext) error {
	// make  sure that the author of the message can modify the bank
	if err := b.userCanModifyBank(ctx, ctx.Message.Member); err != nil {
		return err
	}

	// delete every bank item associated with the guild bank
	_, err := b.Database.BankItem.Delete().
		Where(
			bankitem.HasBankWith(
				guildbank.HasGuildWith(
					guild.DiscordID(ctx.GuildID),
				),
			),
		).Exec(ctx)
	if err != nil {
		return err
	}

	// reset the currency of the guild bank
	guildBank, err := ctx.GuildBank()
	if err != nil {
		return err
	}
	// zero out the guild balance
	guildBank.Update().AddBalance(-guildBank.Balance).Exec(ctx)

	// nothing went wrong
	return b.UpdateBankListing(ctx)
}

// CommandContext holds the contextual information for a message that we receive
type CommandContext struct {
	context.Context
	Banker    *Banker
	GuildID   string
	ChannelID string
	Message   *discordgo.Message
}

// Guild returns the database entry for the current guild
func (ctx *CommandContext) Guild() (*db.Guild, error) {
	// look up the database entry associated with this guild
	return ctx.Banker.Database.Guild.Query().
		Where(guild.DiscordID(ctx.GuildID)).
		Only(context.Background())
}

// GuildBank returns the build bank object associated with the current context
func (ctx *CommandContext) GuildBank() (*db.GuildBank, error) {
	return ctx.Banker.Database.GuildBank.Query().
		Where(guildbank.HasGuildWith(guild.DiscordID(ctx.GuildID))).
		Only(ctx)
}

// MemberName returns the display name for a member
func (ctx *CommandContext) MemberName(user *discordgo.User) string {
	// look up the membership for this user
	member, err := ctx.Banker.Discord.GuildMember(ctx.GuildID, user.ID)
	if err != nil {
		return ""
	}

	// if there is a nickname use it
	if member.Nick != "" {
		return member.Nick
	}

	// theres no nickname so the username will have to do
	return user.Username
}
