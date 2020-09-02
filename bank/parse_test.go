package bank_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AlecAivazis/jeeves/bot"
	"github.com/AlecAivazis/jeeves/data"
)

func TestParseTransaction(t *testing.T) {
	lavaCoreID, _ := data.ItemID("Lava Core")

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
				Item:   data.ItemIDGold,
			},
		},
		{
			Entry: "2s",
			Expected: bot.Transaction{
				Amount: 200,
				Item:   data.ItemIDGold,
			},
		},
		{
			Entry: "2g",
			Expected: bot.Transaction{
				Amount: 20000,
				Item:   data.ItemIDGold,
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
