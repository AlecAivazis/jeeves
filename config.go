package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// BotToken holds the token to use to authenticate the bot
var BotToken string

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found.")
	}

	// load the variables from the environment
	if t := os.Getenv("TOKEN"); t != "" {
		BotToken = t
	}
}
