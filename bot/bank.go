package jeeves

import (
	"github.com/bwmarrin/discordgo"
)

// Jeeve's bank features are summarized with the following commands:
// !deposit 1xA 2xB C - adds one A, two B, and one C to the guild's bank
// !withdraw 1xA 2xB C (for arcanite reaper) - removes one A, two B, and one C to the guild's bank with the provided note

// BankHandler is invoked on every message that is sent to a channel the
// bot is authenticated with
func (b *JeevesBot) BankHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	//

}
