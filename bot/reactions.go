package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type ReactionCallback func(*Message)

type Message struct {
	discordgo.Message
}

func (b *JeevesBot) ReactionHandler(session *discordgo.Session, message *discordgo.MessageReactionAdd) {
	fmt.Println(message.Emoji.APIName(), message)
}

func (b *JeevesBot) RegisterMessageReactionCallback(message *Message, cb ReactionCallback) error {
	return nil
}
