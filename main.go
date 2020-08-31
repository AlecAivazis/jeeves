package main

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/jeeves/bot"
	"github.com/AlecAivazis/jeeves/config"

	_ "github.com/lib/pq"
)

func main() {
	// if we aren't running locally
	if !config.LocalMode {
		// load config values from google secrets
		if err := config.LoadSecrets(); err != nil {
			fmt.Println(err)
			return
		}

	}

	// create an instance of jeeves we can run
	jeeves, err := bot.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer func() {
		// if the main goroutine finished with a panic
		if r := recover(); r != nil {
			// keep the bot running
			jeeves.Start()
		}
	}()

	// start the jeeves bot
	jeeves.Start()
}
