package bot

import (
	"context"

	"github.com/AlecAivazis/jeeves/db"
	"github.com/AlecAivazis/jeeves/db/guild"
)

func (b *JeevesBot) GuildFromContext(ctx *CommandContext) (*db.Guild, error) {
	// look up the database entry associated with this guild
	return b.Database.Guild.Query().
		Where(guild.DiscordID(ctx.GuildID)).
		Only(context.Background())
}
