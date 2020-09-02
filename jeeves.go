package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"

	"github.com/AlecAivazis/jeeves/bank"
	"github.com/AlecAivazis/jeeves/config"
	"github.com/AlecAivazis/jeeves/db"
	"github.com/AlecAivazis/jeeves/db/guild"
)

type JeevesBot struct {
	Database *db.Client
	Discord  *discordgo.Session
}

// Start opens Jeeve's connection with the database and discord clients
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

	// add the various handlers
	dg.AddHandler(b.NewGuild)
	bank.AddHandlers(dg, client)

	// make sure the schema is up to date
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// Open a websocket connection to Discord and begin listening.
	if err := b.Discord.Open(); err != nil {
		return errors.New("error opening connection: " + err.Error())
	}

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

// NewGuild is invoked when a guild is registered with the bot
func (b *JeevesBot) NewGuild(s *discordgo.Session, event *discordgo.GuildCreate) {
	// only register guilds we have access to
	if event.Guild.Unavailable {
		fmt.Println("that guild is unavailable")
		return
	}

	// have we recorded this guild before?
	exists, err := b.Database.Guild.Query().
		Where(guild.DiscordID(event.Guild.ID)).
		Exist(context.Background())
	if err != nil {
		fmt.Println("error checking if we have seen guild before: " + err.Error())
		return
	}

	// if we're never seen the guild before
	if !exists {
		// add an entry in the database for the new guild
		_, err := b.Database.Guild.Create().
			SetDiscordID(event.Guild.ID).
			Save(context.Background())
		if err != nil {
			fmt.Println(err)
		}
	}
}
