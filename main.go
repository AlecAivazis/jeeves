package main

import (
	"fmt"

	"github.com/AlecAivazis/jeeves/bot"
	"github.com/AlecAivazis/jeeves/config"

	_ "github.com/lib/pq"
)

func main() {
	// if we are not running locally
	if !config.LocalMode {
		// load secrets from google
		if err := config.LoadSecrets(); err != nil {
			fmt.Println(err)
			return
		}
	}

	defer func() {
		// if the main goroutine finished with a panic
		if r := recover(); r != nil {
			// keep the bot running
			bot.Start()
		}
	}()

	bot.Start()
}
