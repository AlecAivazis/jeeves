package bot

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// BotToken holds the token to use to authenticate the bot
var BotToken string

// DB configuration
var DBHost = "0.0.0.0"
var DBName = "prisma"
var DBPort = "5432"
var DBUser = "prisma"
var DBPassword = "prisma"

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
