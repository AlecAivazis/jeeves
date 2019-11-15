package bot

import (
	"context"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// CommandContext holds the contextual information for a message that we receive
type CommandContext struct {
	context.Context
	GuildID   string
	ChannelID string
}

// CommandHandler handles the parsing and dispatching of commands for Jeeves
func (b *JeevesBot) CommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	// only look at commands
	if message.Content[0] != '!' {
		return
	}

	// since the message is presumably text, we care about words, not letters
	words := strings.SplitN(message.Content[1:], " ", 1)
	command := words[0]

	// construct the context object
	ctx := &CommandContext{
		GuildID:   message.GuildID,
		ChannelID: message.ChannelID,
		Context:   context.Background(),
	}

	var err error
	// check the command against our known strings
	switch command {
	case CommandAssignBankChannel:
		err = b.InitializeBankChannel(ctx)
	case CommandDeposit:
		err = b.DepositItems(ctx, strings.Split(words[1], ","))
	case CommandWithdraw:
		err = b.WithdrawItems(ctx, strings.Split(words[1], ","))
	}
	// if the command failed
	if err != nil {
		// send the error to the channel we received the message on
		b.ReportError(message.ChannelID, err)
		return
	}
}
