package bot

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"

	_ "github.com/lib/pq"

	"github.com/AlecAivazis/jeeves/db"
)

func Start() {
	// if there is no token
	if BotToken == "" {
		// tell the user what happened
		fmt.Println("Please provide a token via the TOKEN environment variable")
		// don't continue
		os.Exit(1)
		return
	}

	// create a new Discord session using the provided bot token
	dg, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}
	// make sure we close the bot when we're done
	defer dg.Close()

	// open up a client with the configured values
	client, err := db.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", DBHost, DBPort, DBUser, DBPassword, DBName))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// make sure the schema is up to date
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// instantiate the bot
	bot := &JeevesBot{
		Database: client,
		Discord:  dg,
	}

	// add the various handlers
	dg.AddHandler(bot.NewGuild)
	dg.AddHandler(bot.CommandHandler)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// wait for some kind of signal to stop
	fmt.Println("Jeeves is now running. Press ctrl+c to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

// JeevesBot provides context for the discord handlers
type JeevesBot struct {
	Database *db.Client
	Discord  *discordgo.Session
}

// ReportError sends the error to the specified channel
func (b *JeevesBot) ReportError(channel string, errorToReport error) (err error) {
	// send the error message to the channel
	_, err = b.Discord.ChannelMessageSend(channel, errorToReport.Error())

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
