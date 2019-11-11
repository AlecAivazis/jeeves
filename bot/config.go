package bot

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// BotToken holds the token to use to authenticate the bot
var BotToken string

// DB configuration
var DBHost string
var DBName string
var DBPort string
var DBUser string
var DBPassword string

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found.")
	}

	// load the variables from the environment
	if t := os.Getenv("TOKEN"); t != "" {
		BotToken = t
	}

	// load the db configuration
	if t := os.Getenv("DB_HOST"); t != "" {
		DBHost = t
	}
	if t := os.Getenv("DB_NAME"); t != "" {
		DBName = t
	}
	if t := os.Getenv("DB_PASSWORD"); t != "" {
		DBPassword = t
	}
	if t := os.Getenv("DB_USER"); t != "" {
		DBUser = t
	}
	if t := os.Getenv("DB_PORT"); t != "" {
		DBPort = t
	}
}
