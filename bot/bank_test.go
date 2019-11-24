package bot_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AlecAivazis/jeeves/bot"
)

func TestParseTransaction(t *testing.T) {
	lavaCoreID, _ := bot.ItemID("Lava Core")

	// the table
	table := []struct {
		Entry    string
		Expected bot.Transaction
	}{
		{
			Entry: "2xLava Core",
			Expected: bot.Transaction{
				Amount: 2,
				Item:   lavaCoreID,
			},
		},
		{
			Entry: "2x Lava Core",
			Expected: bot.Transaction{
				Amount: 2,
				Item:   lavaCoreID,
			},
		},
		{
			Entry: " 2xLava Core",
			Expected: bot.Transaction{
				Amount: 2,
				Item:   lavaCoreID,
			},
		},
		{
			Entry: " 2x Lava Core",
			Expected: bot.Transaction{
				Amount: 2,
				Item:   lavaCoreID,
			},
		},
		{
			Entry: "2c",
			Expected: bot.Transaction{
				Amount: 2,
				Item:   bot.ItemIDGold,
			},
		},
		{
			Entry: "2s",
			Expected: bot.Transaction{
				Amount: 200,
				Item:   bot.ItemIDGold,
			},
		},
		{
			Entry: "2g",
			Expected: bot.Transaction{
				Amount: 20000,
				Item:   bot.ItemIDGold,
			},
		},
	}

	for _, row := range table {
		t.Run(row.Entry, func(t *testing.T) {
			tx, err := bot.ParseTransaction(row.Entry)
			if !assert.Nil(t, err) {
				return
			}

			// make sure that we parse the results correctly
			assert.Equal(t, row.Expected, tx)
		})
	}
}
