package bot

import (
	"errors"

	"strings"
)

// ItemIDGold is the item id that we use to internall represent a gold deposit/withdrawl
var ItemIDGold = "gold"

var itemNames = map[string]string{}
var itemIDs = map[string]string{}

// ItemID returns the WoW item ID for the item with the given name
func ItemID(name string) (string, error) {
	// the id of gold is the special string
	if name == ItemIDGold {
		return name, nil
	}

	// if we have an entry for that item, use it
	if item, ok := itemIDs[strings.ToLower(name)]; ok {
		return item, nil
	}

	// there was no entry for that name
	return "", errors.New("I don't recognize this item: " + name)
}

// ItemName returns the name of an item given its id
func ItemName(id string) (string, error) {
	name, ok := itemNames[id]
	if !ok {
		return "", errors.New("could not find the name for an item with id " + id)
	}

	// return the name
	return name, nil
}

func init() {
	// invert each entry in the map so we can look up names if we have the ID
	for name, id := range itemData {
		itemNames[id] = name
		itemIDs[strings.ToLower(name)] = id
	}
}
