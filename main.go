package main

import (
	"github.com/AlecAivazis/jeeves/bot"
)

func main() {
	defer func() {
		// if the main goroutine finished with a panic
		if r := recover(); r != nil {
			// keep the bot running
			bot.Start()

		}
	}()

	bot.Start()
}
