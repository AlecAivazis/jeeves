package bank

import (
	"errors"

	"github.com/AlecAivazis/jeeves/data"
	"github.com/AlecAivazis/jeeves/db/bankitem"
	"github.com/bwmarrin/discordgo"
)

func (b *Banker) CheckInventory(ctx *CommandContext, items []string) error {
	// before we can check item quantity, we should check for spelling mistakes
	transactions, err := ParseTransactions(items)
	if err != nil {
		return err
	}

	// find the bank for this guild
	guildBank, err := ctx.GuildBank()
	if err != nil {
		return err
	}

	// we need to add each item to the database
	for _, transaction := range transactions {
		// pull out the constants of the transaction
		item := transaction.Item
		amount := transaction.Amount

		// if we are depositing gold
		if item == data.ItemIDGold {
			// if the guild does not have enough balance
			if guildBank.Balance < amount {
				return errors.New("we don't have that much money in the bank")
			}

			// we're done processing it
			continue
		}

		// does this bank have a record for the item
		existingItems, err := guildBank.
			QueryItems().
			Where(bankitem.ItemID(item)).
			All(ctx)
		if err != nil {
			return err
		}

		// if we haven't seen the item before
		if len(existingItems) == 0 {
			// we can't withdraw it!
			return errors.New("it does not look like we have that item in the bank")
		}

		// make sure there are enough items in the bank
		if amount > existingItems[0].Quantity {
			return errors.New("there is not enough of that item in the bank")
		}
	}

	// the withdraw is valid
	return nil
}

// userCanModifyBank returns an error if the user shouldn't be allowed to touch the bank
func (b *Banker) userCanModifyBank(ctx *CommandContext, member *discordgo.Member) error {
	// grab the id of the banker role
	bankRoleID, err := b.BankerRoleID(ctx)
	if err != nil {
		return err
	}

	// see if the author of th emessage has the bank role
	for _, roleID := range member.Roles {
		if roleID == bankRoleID {
			return nil
		}
	}

	// the user does not have the right permissions
	return errors.New("only Bankers can modify the bank")
}

// BankerRoleID returns the id corresponding to the banker role in this guild.
func (b *Banker) BankerRoleID(ctx *CommandContext) (string, error) {
	// look up the roles in the guild so we can compare role IDs
	roles, err := b.Discord.GuildRoles(ctx.GuildID)
	if err != nil {
		return "", errors.New("I could not look up the roles of this server. Maybe I dont have the right persmissions?")
	}

	// find the id of the bank role
	var bankRoleID string
	for _, role := range roles {
		if role.Name == RoleBanker {
			bankRoleID = role.ID
		}
	}

	// if we couldn't find the role
	if bankRoleID == "" {
		return "", errors.New("it doesn't look like you have the Banker role defined")
	}

	// we found the role
	return bankRoleID, nil
}

// Bankers returns a list of all of the bankers in the guild
func (b *Banker) Bankers(ctx *CommandContext) ([]*discordgo.Member, error) {
	// get the first list of members
	members, err := b.Discord.GuildMembers(ctx.GuildID, "", 1000)
	if err != nil {
		return nil, err
	}

	// we should keep looking
	for {
		// get the next 1000
		nextPage, err := b.Discord.GuildMembers(ctx.GuildID, members[len(members)-1].User.ID, 1000)
		if err != nil {
			return nil, err
		}

		// if there are no more
		if len(nextPage) == 0 {
			// stop looking
			break
		}

		// add what we found to the list
		members = append(members, nextPage...)
	}

	// the list of bankers
	bankers := []*discordgo.Member{}

	for _, member := range members {
		if b.userCanModifyBank(ctx, member) == nil {
			bankers = append(bankers, member)
		}
	}

	return bankers, nil
}
