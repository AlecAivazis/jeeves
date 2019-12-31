package bot

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"

	_ "github.com/lib/pq"

	"github.com/AlecAivazis/jeeves/db"
)

// JeevesBot provides context for the discord handlers
type JeevesBot struct {
	Database *db.Client
	Discord  *discordgo.Session
}

type ReactionCallback func(*Message)

type Message struct {
	discordgo.Message
}

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
	fmt.Println(message.Emoji.APIName(), message.Emoji)
}

func (b *JeevesBot) RegisterMessageReactionCallback(message *Message, cb ReactionCallback) error {
	return nil
}
