package main

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/jeeves/bot"

	_ "github.com/lib/pq"
)

func main() {
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
	bot.Start()
}
