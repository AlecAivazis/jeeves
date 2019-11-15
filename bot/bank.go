// // Jeeve's bank features are summarized with the following commands:
// !deposit 1xA 2xB C - adds one A, two B, and one C to the guild's bank
// !withdraw 1xA 2xB C (for arcanite reaper) - removes one A, two B, and one C to the guild's bank with the provided note

package bot

import (
	"context"

	"github.com/AlecAivazis/jeeves/db/bankitem"
	"github.com/AlecAivazis/jeeves/db/guild"
	"github.com/AlecAivazis/jeeves/db/predicate"
)

const (
	// CommandDeposit defines the command used to deposit items into the guild bank
	CommandDeposit = "deposit"
	// CommandWithdraw defines the command used to withdraw items from the guild bank
	CommandWithdraw = "withdraw"
)

func (b *JeevesBot) WithdrawItems(ctx *CommandContext, items []string) error {
	return nil
}

func (b *JeevesBot) DepositItems(ctx *CommandContext, items []string) error {
	// we need to add each item to the database
	for _, item := range items {
		// we need a get or update on the same guild channel config
		wherePredicates := []predicate.BankItem{
			bankitem.HasGuildWith(guild.DiscordID(ctx.GuildID)),
			bankitem.ItemID(item),
		}

		// look if we have an existing record for the item
		previousRecord, err := b.Database.BankItem.Query().
			Where(wherePredicates...).
			Only(context.Background())
		if err != nil {
			return err
		}

		// if this is a new item for the bank we need to add it to the database
		if previousRecord == nil {
			// grab the guild from context
			guild, err := b.GuildFromContext(ctx)
			if err != nil {
				return err
			}

			// store the association in the database
			b.Database.BankItem.Create().
				SetGuild(guild).
				SetItemID(item).
				SetQuantity(1).
				Save(context.Background())

		} else {
			// otherwise we need to update the existing channel association
			b.Database.BankItem.Update().
				Where(wherePredicates...).
				SetQuantity(previousRecord.Quantity + 1).
				Save(context.Background())
		}
	}

	// once we are done adding the items we should update the listing
	return b.UpdateBankListing()
}

func (b *JeevesBot) UpdateBankListing() error {

	// nothing went wrong
	return nil
}
