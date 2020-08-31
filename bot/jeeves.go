package bot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"github.com/AlecAivazis/jeeves/config"
	"github.com/AlecAivazis/jeeves/db"
)

// JeevesBot provides context for the discord handlers
type JeevesBot struct {
	Database *db.Client
	Discord  *discordgo.Session
}

type Message struct {
	discordgo.Message
}

type ReactionCallback func(*discordgo.MessageReactionAdd)

func New() (*JeevesBot, error) {
	return &JeevesBot{}, nil
}

func (b *JeevesBot) Start() error {
	// if there is no token
	if config.BotToken == "" {
		// don't continue
		return errors.New("Please provide a token via the TOKEN environment variable")
	}

	// create a new Discord session using the provided bot token
	dg, err := discordgo.New("Bot " + config.BotToken)
	if err != nil {
		return errors.New("Error creating Discord session: " + err.Error())
	}

	// add the various handlers
	dg.AddHandler(b.NewGuild)
	dg.AddHandler(b.CommandHandler)

	// open up a client with the configured values
	client, err := db.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	))
	if err != nil {
		panic(err)
	}

	// hold onto the references to our clients
	b.Discord = dg
	b.Database = client

	// make sure the schema is up to date
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// Open a websocket connection to Discord and begin listening.
	if err := b.Discord.Open(); err != nil {
		return errors.New("error opening connection: " + err.Error())
	}

	// after we are running
	defer func() {
		// if we aren't running because of a panic
		if r := recover(); r != nil {
			// keep the bot running
			if err := b.Start(); err != nil {
				fmt.Println(err)
				b.Stop()
			}

		} else {
			// we stopped naturally so make sure the server cleans up
			b.Stop()
		}

	}()

	// wait for some kind of signal to stop
	fmt.Println("Jeeves is now running. Press ctrl+c to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	return nil
}

func (b *JeevesBot) Stop() error {
	if err := b.Database.Close(); err != nil {
		fmt.Println(err)
	}
	// make sure we close the bot when we're done
	if err := b.Discord.Close(); err != nil {
		fmt.Println(err)
	}

	return nil
}

// ReportError sends the error to the specified channel
func (b *JeevesBot) ReportError(channel string, errorToReport error) (err error) {
	// send the error message to the channel
	_, err = b.Discord.ChannelMessageSend(channel, "Sorry, "+errorToReport.Error())

	return err
}

// NewGuild is invoked when a guild is registered with the bot
func (b *JeevesBot) NewGuild(s *discordgo.Session, event *discordgo.GuildCreate) error {
	// only register guilds we have access to
	if event.Guild.Unavailable {
		return errors.New("that guild is unavailable")
	}

	// add an entry in the database for the new guild
	_, err := b.Database.Guild.Create().
		SetDiscordID(event.Guild.ID).
		Save(context.Background())

	return err
}

// Reply sends a message to the channel in the given context
func (b *JeevesBot) Reply(ctx *CommandContext, message string) (*discordgo.Message, error) {
	return b.Discord.ChannelMessageSend(ctx.ChannelID, message)
}

func (b *JeevesBot) MemberName(ctx *CommandContext, user *discordgo.User) string {
	// look up the membership for this user
	member, err := b.Discord.GuildMember(ctx.GuildID, user.ID)
	if err != nil {
		return ""
	}

	// if there is a nickname use it
	if member.Nick != "" {
		return member.Nick
	}

	// theres no nickname so the username will have to do
	return user.Username
}
