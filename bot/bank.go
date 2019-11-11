package bot

// Jeeve's bank features are summarized with the following commands:
// !deposit 1xA 2xB C - adds one A, two B, and one C to the guild's bank
// !withdraw 1xA 2xB C (for arcanite reaper) - removes one A, two B, and one C to the guild's bank with the provided note

const (
	// CommandDeposit defines the command used to deposit items into the guild bank
	CommandDeposit = "deposit"
	// CommandWithdraw defines the command used to withdraw items from the guild bank
	CommandWithdraw = "withdraw"
)

func (b *JeevesBot) DepositItems(ctx *CommandContext, items []string) error {
	return nil
}

func (b *JeevesBot) WithdrawItems(ctx *CommandContext, items []string) error {
	return nil
}
