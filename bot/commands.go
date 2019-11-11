package jeeves

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// CommandContext holds the contextual information for a message that we receive
type CommandContext struct {
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
	words := strings.Split(message.Content[1:], " ")
	command := words[0]
	args := words[1:]

	// construct the context object
	ctx := &CommandContext{
		GuildID:   message.GuildID,
		ChannelID: message.ChannelID,
	}

	var err error
	// check the command against our known strings
	switch command {
	case CommandDeposit:
		err = b.DepositItems(ctx, args)
	case CommandWithdraw:
		err = b.WithdrawItems(ctx, args)
	case CommandAssignChannel:
		err = b.RegisterChannelRole(ctx, args[0])
	}
	// if the command failed
	if err != nil {
		// send the error to the channel we received the message on
		b.ReportError(message.ChannelID, err)
		return
	}
}
