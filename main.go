package main

import (
	"fmt"
	"net/http"

	"github.com/AlecAivazis/jeeves/bot"
	"github.com/AlecAivazis/jeeves/config"

	_ "github.com/lib/pq"
)

func main() {
	// if we are not running locally
	if !config.LocalMode {
		fmt.Println("loading secrets")
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

	// start the jeeves bot
	bot.Start()

	// make sure we have a status check for Cloud Run
	http.HandleFunc("/", StatusCheck)
	err := http.ListenAndServe(":"+config.StatusCheckPort, nil)
	if err != nil {
		fmt.Println(err)
	}
}

// StatusCheck is an http endpoint that we can use to make sure the bot is still running
func StatusCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}
