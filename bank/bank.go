package bank

import (
	"github.com/AlecAivazis/jeeves/db"
	"github.com/bwmarrin/discordgo"
)

// Banker is the central entity for managing a guild's bank
type Banker struct {
	Database *db.Client
	Discord  *discordgo.Session
}

func AddHandlers(dg *discordgo.Session, database *db.Client) {
	// instantiate a banker we'll use as an execution context
	banker := &Banker{
		Database: database,
		Discord:  dg,
	}

	// add the necessary handlers
	dg.AddHandler(banker.CommandHandler)
}
