package bot

import (
	"errors"
)

var itemNames = map[string]string{}

// ItemID returns the WoW item ID for the item with the given name
func ItemID(name string) (string, error) {
	// do we have an entry for that item
	if item, ok := itemNumbers[properTitle(name)]; ok {
		// stringify the int value
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
