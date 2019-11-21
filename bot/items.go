package bot

import (
	"errors"
)

// ItemIDGold is the item id that we use to internall represent a gold deposit/withdrawl
var ItemIDGold = "gold"

var itemNames = map[string]string{}

// ItemID returns the WoW item ID for the item with the given name
func ItemID(name string) (string, error) {
	// if we are looking up gold (not an item) we have to return a special ID
	if name == "gold" {
		return ItemIDGold, nil
	}

	// if we have an entry for that item, use it
	if item, ok := itemNumbers[properTitle(name)]; ok {
		return item, nil
	}

	// there was no entry for that name
	return "", errors.New("I don't recognize this item: " + name)
}

func init() {
	// invert each entry in the map so we can look up names if we have the ID
	for name := range itemNumbers {
		itemID, _ := ItemID(name)
		itemNames[itemID] = name
	}
}
