package bank

import (
	"bytes"
	"fmt"
	"html/template"
	"math"
	"sort"
	"strings"

	"github.com/AlecAivazis/jeeves/data"
	"github.com/AlecAivazis/jeeves/db"
	"github.com/AlecAivazis/jeeves/db/guildbank"
)

//////////////////////////////////
//
// Guild Bank Display
//
//////////////////////////////////

type bankDisplayData struct {
	Items   []*db.BankItem
	Balance int
}

// UpdateBankListing is called whenever jeeves needs to rerender the bank display
func (b *JeevesBot) UpdateBankListing(ctx *CommandContext) error {
	// find the channel ID for the bank channel for this guild
	bank, err := b.GuildBankFromContext(ctx)
	if err != nil {
		return err
	}

	// get the items in the bank
	items, err := b.Database.GuildBank.Query().
		Where(guildbank.ID(bank.ID)).
		QueryItems().All(ctx)
	if err != nil {
		return err
	}

	// sort the items based on their display name
	sort.SliceStable(items, func(i, j int) bool {
		// figure out the display names of the two items
		nameA, _ := data.ItemName(items[i].ItemID)
		nameB, _ := data.ItemName(items[j].ItemID)

		// i should come before j if i's name is less than j
		return nameA < nameB
	})

	// execute the template
	var contents bytes.Buffer
	err = displayTemplate.Execute(&contents, &bankDisplayData{
		Items:   items,
		Balance: bank.Balance,
	})
	if err != nil {
		return err
	}

	// delete every message in the channel
	messages, err := b.Discord.ChannelMessages(bank.ChannelID, 100, "", "", "")
	if err != nil {
		return err
	}
	messageIDs := []string{}
	for _, msg := range messages {
		messageIDs = append(messageIDs, msg.ID)
	}
	err = b.Discord.ChannelMessagesBulkDelete(bank.ChannelID, messageIDs)
	if err != nil {
		fmt.Println("err", err)
	}

	// we need to break the message we are about to send in 2000 character chunks.
	messagesToSend := []string{}
	lines := strings.Split(contents.String(), "\n")
	currentMessage := ""

	for _, line := range lines {
		// if the line wouldn't push the message over the limit, we can just add it to the list
		if len(currentMessage)+len(line)+1 < 2000 {
			currentMessage += line + "\n"
			continue
		}

		// this message would push the current over the list
		messagesToSend = append(messagesToSend, currentMessage)

		// start with this line
		currentMessage = line + "\n"
	}

	// add the remaining message to the list
	messagesToSend = append(messagesToSend, currentMessage)

	// send each message to the bank
	for _, msg := range messagesToSend {
		_, err := b.Discord.ChannelMessageSend(bank.ChannelID, msg)
		if err != nil {
			return err
		}
	}

	// nothing went wrong
	return nil
}

var displayTemplate *template.Template

// BankDisplayContents is the template used by jeeves to show what's in the bank
const BankDisplayContents = `
Current Balance: {{ format .Balance }}

Bank Contents:
{{- range .Items }}
{{ .Quantity}}x {{ itemName .ItemID }}
{{- end }}
`

func formatGold(balance int) string {
	// the amount of copper will be what's left
	copper := float64(balance)

	// the amount of gold
	gold := math.Floor(float64(balance) / 10000)
	// remove the amount of gold
	copper -= gold * 10000

	// the amount of silver left
	silver := math.Floor(float64(copper) / 100)
	// remove the amount of silver
	copper -= silver * 100

	// return the for
	return fmt.Sprintf("%vg %vs %vc", gold, silver, copper)
}

func init() {
	displayTemplate = template.Must(template.New("bank-display").Funcs(template.FuncMap{
		"itemName": data.ItemName,
		"format":   formatGold,
	}).Parse(BankDisplayContents))
}
