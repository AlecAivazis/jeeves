package main

import (
	"fmt"
	"net/http"
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

	// make sure we have a status check for Cloud Run
	http.HandleFunc("/", StatusCheck)
	err = http.ListenAndServe(":"+config.StatusCheckPort, nil)
	if err != nil {
		fmt.Println(err)
	}
}

// StatusCheck is an http endpoint that we can use to make sure the bot is still running
func StatusCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}
