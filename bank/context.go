package bot

import (
	"context"

	"github.com/AlecAivazis/jeeves/db"
	"github.com/AlecAivazis/jeeves/db/guild"
	"github.com/AlecAivazis/jeeves/db/guildbank"
)

// CommandContext holds the contextual information for a message that we receive
type CommandContext struct {
	context.Context
	Bot       *JeevesBot
	GuildID   string
	ChannelID string
	Message   *Message
}

func (b *JeevesBot) GuildFromContext(ctx *CommandContext) (*db.Guild, error) {
	// look up the database entry associated with this guild
	return b.Database.Guild.Query().
		Where(guild.DiscordID(ctx.GuildID)).
		Only(context.Background())
}

// GuildBank returns the build bank object associated with the current context
func (b *JeevesBot) GuildBankFromContext(ctx *CommandContext) (*db.GuildBank, error) {
	return b.Database.GuildBank.Query().
		Where(guildbank.HasGuildWith(guild.DiscordID(ctx.GuildID))).
		Only(ctx)
}
