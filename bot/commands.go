package bot

import (
	"context"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// CommandContext holds the contextual information for a message that we receive
type CommandContext struct {
	context.Context
	Bot       *JeevesBot
	GuildID   string
	ChannelID string
	Message   *Message
}

// CommandHandler handles the parsing and dispatching of commands for Jeeves
func (b *JeevesBot) CommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	// only look at commands
	if message.Content[0] != '!' {
		return
	}

	// since the message is presumably text, we care about words, not letters
	words := strings.SplitN(message.Content[1:], " ", 2)
	command := words[0]

	// construct the context object
	ctx := &CommandContext{
		Bot:       b,
		GuildID:   message.GuildID,
		ChannelID: message.ChannelID,
		Context:   context.Background(),
		Message:   &Message{Message: *message.Message},
	}

	var err error
	// check the command against our known strings
	switch command {
	case CommandAssignBankChannel:
		err = b.InitializeBankChannel(ctx)
	case CommandDeposit:
		err = b.DepositItems(ctx, ParseItems(words[1]))
	case CommandWithdraw:
		err = b.WithdrawItems(ctx, ParseItems(words[1]))
	case CommandRequest:
		err = b.RequestItems(ctx, ParseItems(words[1]))
	case CommandRefreshBank:
		err = b.UpdateBankListing(ctx)
	}
	// if the command failed
	if err != nil {
		// send the error to the channel we received the message on
		b.ReportError(message.ChannelID, err)
		return
	}

	// confirm the action with a reaction
	err = b.Discord.MessageReactionAdd(message.ChannelID, message.ID, "👍")
	if err != nil {
		b.ReportError(message.ChannelID, err)
	}
}

func ParseItems(input string) []string {
	return strings.Split(strings.Trim(input, ", "), ",")
}
