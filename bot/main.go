package bot

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/AlecAivazis/jeeves/config"
	"github.com/AlecAivazis/jeeves/db"
	"github.com/bwmarrin/discordgo"
)

func Start() {
	// if there is no token
	if config.BotToken == "" {
		// tell the user what happened
		fmt.Println("Please provide a token via the TOKEN environment variable")
		// don't continue
		os.Exit(1)
		return
	}

	// create a new Discord session using the provided bot token
	dg, err := discordgo.New("Bot " + config.BotToken)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}
	// make sure we close the bot when we're done
	defer dg.Close()

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
	defer client.Close()

	// make sure the schema is up to date
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// instantiate the bot
	bot := &JeevesBot{
		Database:          client,
		Discord:           dg,
		ReactionCallbacks: make(map[string][]ReactionCallback),
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
}
