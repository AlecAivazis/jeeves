package bot

import (
	"context"
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"

	"github.com/AlecAivazis/jeeves/db"
)

// JeevesBot provides context for the discord handlers
type JeevesBot struct {
	Database          *db.Client
	Discord           *discordgo.Session
	ReactionCallbacks map[string][]ReactionCallback
	cbMutex           sync.Mutex
}

type Message struct {
	discordgo.Message
}

type ReactionCallback func(*discordgo.MessageReactionAdd)

// ReportError sends the error to the specified channel
func (b *JeevesBot) ReportError(channel string, errorToReport error) (err error) {
	// send the error message to the channel
	_, err = b.Discord.ChannelMessageSend(channel, "Sorry, "+errorToReport.Error())

	return err
}

// NewGuild is invoked when a guild is registered with the bot
func (b *JeevesBot) NewGuild(s *discordgo.Session, event *discordgo.GuildCreate) {
	// only register guilds we have access to
	if event.Guild.Unavailable {
		return
	}

	// add an entry in the database for the new guild
	b.Database.Guild.Create().
		SetDiscordID(event.Guild.ID).
		Save(context.Background())
}

// Reply sends a message to the channel in the given context
func (b *JeevesBot) Reply(ctx *CommandContext, message string) (*discordgo.Message, error) {
	return b.Discord.ChannelMessageSend(ctx.ChannelID, message)
}

func (b *JeevesBot) ReactionHandler(session *discordgo.Session, message *discordgo.MessageReactionAdd) {
	// if we have a callback for the message
	if b.ReactionCallbacks[message.MessageID] == nil {
		return
	}

	// invoke each of the callbacks with the reaction
	for _, cb := range b.ReactionCallbacks[message.MessageID] {
		cb(message)
	}
}

func (b *JeevesBot) UnregisterMessageReactionCallback(message *Message) error {
	// remove the message id from the dispatch map
	b.cbMutex.Lock()
	delete(b.ReactionCallbacks, message.ID)
	b.cbMutex.Unlock()

	// nothing went wrong
	return nil
}

func (b *JeevesBot) RegisterMessageReactionCallback(message *Message, cb ReactionCallback) error {
	fmt.Println("Registering callback for ", message.ID)

	// ensure atomic access to the list of callbacks
	b.cbMutex.Lock()

	// if we dont have an registered callbacks for the message make sure there is a list to add to
	if b.ReactionCallbacks[message.ID] == nil {
		b.ReactionCallbacks[message.ID] = []ReactionCallback{}
	}

	// add the callback to the list
	b.ReactionCallbacks[message.ID] = append(b.ReactionCallbacks[message.ID], cb)

	// we're done modifying the callback list
	b.cbMutex.Unlock()

	// nothing went wrong
	return nil
}
