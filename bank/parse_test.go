package bank_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/AlecAivazis/jeeves/bank"
	"github.com/AlecAivazis/jeeves/data"
)

func TestParseTransaction(t *testing.T) {
	lavaCoreID, _ := data.ItemID("Lava Core")

	// the table
	table := []struct {
		Entry    string
		Expected bank.Transaction
	}{
		{
			Entry: "2xLava Core",
			Expected: bank.Transaction{
				Amount: 2,
				Item:   lavaCoreID,
			},
		},
		{
			Entry: "2x Lava Core",
			Expected: bank.Transaction{
				Amount: 2,
				Item:   lavaCoreID,
			},
		},
		{
			Entry: " 2xLava Core",
			Expected: bank.Transaction{
				Amount: 2,
				Item:   lavaCoreID,
			},
		},
		{
			Entry: " 2x Lava Core",
			Expected: bank.Transaction{
				Amount: 2,
				Item:   lavaCoreID,
			},
		},
		{
			Entry: "2c",
			Expected: bank.Transaction{
				Amount: 2,
				Item:   data.ItemIDGold,
			},
		},
		{
			Entry: "2s",
			Expected: bank.Transaction{
				Amount: 200,
				Item:   data.ItemIDGold,
			},
		},
		{
			Entry: "2g",
			Expected: bank.Transaction{
				Amount: 20000,
				Item:   data.ItemIDGold,
			},
		},
	}

	for _, row := range table {
		t.Run(row.Entry, func(t *testing.T) {
			tx, err := bank.ParseTransaction(row.Entry)
			if !assert.Nil(t, err) {
				return
			}

			// make sure that we parse the results correctly
			assert.Equal(t, row.Expected, tx)
		})
	}
}
