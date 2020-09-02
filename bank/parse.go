package bank

import (
	"strconv"
	"strings"

	"github.com/AlecAivazis/jeeves/data"
)

type Transaction struct {
	Item   string
	Amount int
}

func ParseTransactions(before []string) ([]Transaction, error) {
	transactions := []Transaction{}

	for _, entry := range before {
		txn, err := ParseTransaction(entry)
		if err != nil {
			return nil, err
		}

		// add the transaction to the list
		transactions = append(transactions, txn)
	}

	return transactions, nil
}

const numbers = "1234567890"

// ParseTransaction takes a string like "2xLava Core" and extracts the quantity and item referenced
func ParseTransaction(entry string) (Transaction, error) {
	// get the name ready and normalized
	item := strings.ToLower(strings.Trim(entry, " "))

	// the transaction to return
	transaction := Transaction{
		Amount: 1,
	}

	// we are going to consume until we find something that's not a number
	amount := ""

	// if the first character is a number we want to keep looking down the string
	// and group up all of the numbers to form a single quantity
	if strings.Contains(numbers, string(item[0])) {

		// look at all of the characters in the word
		for i, char := range item {
			// if the character is a number
			if strings.Contains(numbers, string(item[i])) {
				// add it to the running total
				amount += string(char)

				// we found something that's not a number
			} else {
				// try to parse the quantity as a number
				quantity, _ := strconv.Atoi(amount)
				transaction.Amount = quantity

				// we want to "eat up" what we've treated as the number
				item = item[i:]

				// stop consuming text
				break
			}
		}
	}

	// remove any spaces around the item
	item = strings.Trim(item, " ")

	// if the user is depositing gold
	if item == GoldDelimiter {
		transaction.Item = data.ItemIDGold
		transaction.Amount *= 10000 // 1 gold = 100 silver = 10000 copper

		// the user is depositing silver
	} else if item == SilverDelimiter {
		transaction.Item = data.ItemIDGold
		transaction.Amount *= 100 // 1 gold = 100 silver

		// the user is depositing copper
	} else if item == CopperDelimiter {
		transaction.Item = data.ItemIDGold

		// if we started the message with a number and the next character
		// is the quantity delimiter then we are depositing some number of an item
	} else if amount != "" && item[0] == QuantityDelimiter {
		// convert the item name into the normalized ID
		itemID, err := data.ItemID(strings.Trim(item[1:], " "))
		if err != nil {
			return transaction, err
		}
		transaction.Item = itemID
	} else {
		// convert the item name into the normalized ID
		itemID, err := data.ItemID(strings.Trim(entry, " "))
		if err != nil {
			return transaction, err
		}
		transaction.Item = itemID

	}

	// we're done
	return transaction, nil
}
