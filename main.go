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

	// start the jeeves bot
	err = jeeves.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// after we are running
	defer func() {
		// if we aren't running because of a panic
		if r := recover(); r != nil {
			// keep the bot running
			if err := jeeves.Start(); err != nil {
				fmt.Println(err)
			}

		} else {
			// we stopped naturally so make sure the server cleans up
			jeeves.Stop()
		}
	}()
}
