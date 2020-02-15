// !deposit 1xA 2xB C - adds one A, two B, and one C to the guild's bank
// !withdraw 1xA 2xB C (for arcanite reaper) - removes one A, two B, and one C to the guild's bank with the provided note

package bot

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/AlecAivazis/jeeves/db"
	"github.com/AlecAivazis/jeeves/db/bankitem"
	"github.com/AlecAivazis/jeeves/db/guild"
	"github.com/AlecAivazis/jeeves/db/guildbank"
	"github.com/AlecAivazis/jeeves/db/predicate"
	"github.com/bwmarrin/discordgo"
)

const (
	// CommandDeposit defines the command used to deposit items into the guild bank
	CommandDeposit = "deposit"
	// CommandWithdraw defines the command used to withdraw items from the guild bank
	CommandWithdraw = "withdraw"
	// CommandAssignBankChannel defines the command used to assign a channel to use to display the bank
	CommandAssignBankChannel = "jeeves-assign-bank"
	// CommandRefreshBank can be used to force the bank to be re-rendered
	CommandRefreshBank = "refresh-bank"
	// CommandRequest can be used to submit a request to be fulfilled by a banker
	CommandRequest = "request"
)

const (
	// RoleBanker defines the public name of the role to give non-admin users permissions to modify the bank
	RoleBanker = "Banker"
)

const (
	// QuantityDelimiter is the character that separates the amount from the item. Ie, "x" in 2xLava Core
	QuantityDelimiter = 'x'
	// CopperDelimiter is the character that designates a copper deposit
	CopperDelimiter = "c"
	// SilverDelimiter is the character that designates a silver deposit
	SilverDelimiter = "s"
	// GoldDelimiter is the character that designates a Gold deposit
	GoldDelimiter = "g"
)

type Transaction struct {
	Item   string
	Amount int
}

// InitializeBankChannel is called when the user intends to assign a channel for use to display the bank
func (b *JeevesBot) InitializeBankChannel(ctx *CommandContext) (bool, error) {
	// confirm the action with the user
	_, err := b.Reply(ctx, "Okay! Please give me a moment to set up your guild bank...")
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
		_, err = b.Reply(ctx, "I am creating the Banker role. Assign this to non-Admin users you want"+
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
			b.Reply(ctx, "Hmmm something happened when I tried to do that. Maybe I don't have permissions to add roles?"+
				" Before you try again, either make sure I have that permission or you can create the \"Banker\" role and I"+
				" won't try to make it next time.")
		}

	}

	// send the display message now so they know what they can delete
	display, err := b.Reply(ctx, "All set! Your guild bank's contents will go here. You are free to"+
		" delete any other message in this channel but please do not delete this message. I will update it as your bankers"+
		" add items to the bank.")
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
		guild, err := b.GuildFromContext(ctx)
		if err != nil {
			return false, err
		}

		// create the entry for the bank
		_, err = b.Database.GuildBank.Create().
			SetGuild(guild).
			SetChannelID(ctx.ChannelID).
			SetDisplayMessageID(display.ID).
			Save(ctx)
		if err != nil {
			return false, err
		}

	} else {
		// update the bank entry to have the new channel
		_, err = b.Database.GuildBank.Update().
			Where(wherePredicates...).
			SetChannelID(ctx.ChannelID).
			SetDisplayMessageID(display.ID).
			Save(ctx)
		if err != nil {
			return false, err
		}
	}

	// nothing went wrong
	return true, nil
}

// WithdrawItems is used when the user wants to withdraw the specified items from the bank. Will update the display message.
func (b *JeevesBot) WithdrawItems(ctx *CommandContext, items []string) (bool, error) {
	// make sure the user has the right permissions
	if err := b.userCanModifyBank(ctx, ctx.Message.Member); err != nil {
		return false, err
	}

	// make sure the user has the right permissions
	if err := b.ValidateWithdraw(ctx, items); err != nil {
		return false, err
	}

	// find the bank for this guild
	guildBank, err := b.GuildBank(ctx)
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
		if item == ItemIDGold {
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
func (b *JeevesBot) DepositItems(ctx *CommandContext, items []string) (bool, error) {
	// make sure the user has the right permissions
	if err := b.userCanModifyBank(ctx, ctx.Message.Member); err != nil {
		return false, err
	}

	// find the bank for this guild
	guildBank, err := b.GuildBank(ctx)
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
		if item == ItemIDGold {
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

func (b *JeevesBot) RequestItems(ctx *CommandContext, items []string) (bool, error) {
	// make sure we would be able to perform the withdraw (ie, the item names are valid)
	if err := b.ValidateWithdraw(ctx, items); err != nil {
		return false, err
	}

	// confirm that we see the withdraw
	b.Discord.MessageReactionAdd(ctx.ChannelID, ctx.Message.ID, "👀")

	// listen for the indication that the banker sent the items
	// that reaction can be any of the following: ☑️, ✔️, or ✅
	return false, ctx.Bot.RegisterMessageReactionCallback(ctx.Message, func(reaction *discordgo.MessageReactionAdd) {
		// get the user object
		member, err := b.Discord.State.Member(reaction.GuildID, reaction.UserID)
		if err != nil {
			fmt.Println(err)
			return
		}

		// make sure the user has the right permissions
		if err := b.userCanModifyBank(ctx, member); err != nil {
			return
		}

		// if the reaction is one of the approved ones
		if strings.Contains("☑️✔️✅", reaction.Emoji.APIName()) {
			// perform the withdraw
			_, err := b.WithdrawItems(ctx, items)
			if err != nil {
				b.ReportError(ctx.ChannelID, err)
				return
			}

			// confirm that we performed the withdraw
			b.Discord.MessageReactionAdd(ctx.ChannelID, ctx.Message.ID, "👍")

			// we don't need to listen for reactions to this message any more
			ctx.Bot.UnregisterMessageReactionCallback(ctx.Message)
		}
	})
}

func (b *JeevesBot) ValidateWithdraw(ctx *CommandContext, items []string) error {
	// find the bank for this guild
	guildBank, err := b.GuildBank(ctx)
	if err != nil {
		return err
	}

	// parse the transactions into something we can process
	transactions, err := ParseTransactions(items)
	if err != nil {
		return err
	}

	// we need to add each item to the database
	for _, transaction := range transactions {
		// pull out the constants of the transaction
		item := transaction.Item
		amount := transaction.Amount

		// if we are depositing gold
		if item == ItemIDGold {
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

// GuildBank returns the build bank object associated with the current context
func (b *JeevesBot) GuildBank(ctx *CommandContext) (*db.GuildBank, error) {
	return b.Database.GuildBank.Query().
		Where(guildbank.HasGuildWith(guild.DiscordID(ctx.GuildID))).
		Only(ctx)
}

type bankDisplayData struct {
	Items   []*db.BankItem
	Balance int
}

const numbers = "1234567890"

func ParseTransactions(before []string) ([]Transaction, error) {
	transactions := []Transaction{}

	for _, entry := range before {
		txn, err := parseTransaction(entry)
		if err != nil {
			return nil, err
		}

		// add the transaction to the list
		transactions = append(transactions, txn)
	}

	return transactions, nil
}

// userCanModifyBank returns an error if the user shouldn't be allowed to touch the bank
func (b *JeevesBot) userCanModifyBank(ctx *CommandContext, member *discordgo.Member) error {
	// look up the roles in the guild so we can compare role IDs
	roles, err := b.Discord.GuildRoles(ctx.GuildID)
	if err != nil {
		return errors.New("I could not look up the roles of this server. Maybe I dont have the right persmissions?")
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
		return errors.New("it doesn't look like you have the Banker role defined")
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

// ParseTransaction takes a string like "2xLava Core" and extracts the quantity and item referenced
func parseTransaction(entry string) (Transaction, error) {
	// get the name ready and normalized
	item := strings.ToLower(strings.Trim(entry, " "))

	// the transaction to return
	transaction := Transaction{
		Amount: 1,
	}

	// we are going to consume until we find something that's not a number
	amount := ""

	// if the first character is a number we want to keep looking down the string
	// and group up all of the numbers to form a single quantity
	if strings.Contains(numbers, string(item[0])) {

		// look at all of the characters in the word
		for i, char := range item {
			// if the character is a number
			if strings.Contains(numbers, string(item[i])) {
				// add it to the running total
				amount += string(char)

				// we found something that's not a number
			} else {
				// try to parse the quantity as a number
				quantity, _ := strconv.Atoi(amount)
				transaction.Amount = quantity

				// we want to "eat up" what we've treated as the number
				item = item[i:]

				// stop consuming text
				break
			}
		}
	}

	// remove any spaces around the item
	item = strings.Trim(item, " ")

	// if the user is depositing gold
	if item == GoldDelimiter {
		transaction.Item = ItemIDGold
		transaction.Amount *= 10000 // 1 gold = 100 silver = 10000 copper

		// the user is depositing silver
	} else if item == SilverDelimiter {
		transaction.Item = ItemIDGold
		transaction.Amount *= 100 // 1 gold = 100 silver

		// the user is depositing copper
	} else if item == CopperDelimiter {
		transaction.Item = ItemIDGold

		// if we started the message with a number and the next character
		// is the quantity delimiter then we are depositing some number of an item
	} else if amount != "" && item[0] == QuantityDelimiter {
		// convert the item name into the normalized ID
		itemID, err := ItemID(strings.Trim(item[1:], " "))
		if err != nil {
			return transaction, err
		}
		transaction.Item = itemID
	} else {
		// convert the item name into the normalized ID
		itemID, err := ItemID(strings.Trim(entry, " "))
		if err != nil {
			return transaction, err
		}
		transaction.Item = itemID

	}

	// we're done
	return transaction, nil
}

//////////////////////////////////
//
// Guild Bank Display
//
//////////////////////////////////

// UpdateBankListing is called whenever jeeves needs to rerender the bank display
func (b *JeevesBot) UpdateBankListing(ctx *CommandContext) error {
	// find the channel ID for the bank channel for this guild
	bank, err := b.GuildBank(ctx)
	if err != nil {
		return err
	}

	// get the items in the bank
	items, err := b.Database.GuildBank.Query().
		Where(guildbank.ID(bank.ID)).
		QueryItems().All(ctx)
	if err != nil {
		return err
	}

	// sort the items based on their display name
	sort.SliceStable(items, func(i, j int) bool {
		// figure out the display names of the two items
		nameA, _ := ItemName(items[i].ItemID)
		nameB, _ := ItemName(items[j].ItemID)

		// i should come before j if i's name is less than j
		return nameA < nameB
	})

	// execute the template
	var contents bytes.Buffer
	err = displayTemplate.Execute(&contents, &bankDisplayData{
		Items:   items,
		Balance: bank.Balance,
	})
	if err != nil {
		return err
	}
	// update the display message with the items
	_, err = b.Discord.ChannelMessageEdit(bank.ChannelID, bank.DisplayMessageID, contents.String())
	if err != nil {
		return err
	}

	// nothing went wrong
	return nil
}

var displayTemplate *template.Template

// BankDisplayContents is the template used by jeeves to show what's in the bank
const BankDisplayContents = `
Current Balance: {{ format .Balance }}

Bank Contents:
{{- range .Items }}
{{ .Quantity}}x {{ itemName .ItemID }}
{{- end }}
`

func init() {
	displayTemplate = template.Must(template.New("bank-display").Funcs(template.FuncMap{
		"itemName": func(id string) string {
			// if the id is something we recognize
			if name, ok := itemNames[id]; ok {
				return name
			}

			// backwards compatability is hard
			return id

		},
		"format": func(balance int) string {
			// the amount of copper will be what's left
			copper := float64(balance)

			// the amount of gold
			gold := math.Floor(float64(balance) / 10000)
			// remove the amount of gold
			copper -= gold * 10000

			// the amount of silver left
			silver := math.Floor(float64(copper) / 100)
			// remove the amount of silver
			copper -= silver * 100

			// return the for
			return fmt.Sprintf("%vg %vs %vc", gold, silver, copper)
		},
	}).Parse(BankDisplayContents))
}
