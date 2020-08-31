package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// BotToken holds the token to use to authenticate the bot
var BotToken string

// LocalMode prevents Jeeves from reaching out to external services wherever possible
var LocalMode = false

// GoogleCloudProject points the google secrets manager to the right home
var GoogleCloudProject = "aivazis"

// DB configuration
var DBHost = "0.0.0.0"
var DBName = "jeeves"
var DBPort = "5432"
var DBUser = "jeeves"
var DBPassword = "password"

// load the variables from the environment
func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found.")
	}

	if t := os.Getenv("LOCAL"); t == "true" {
		LocalMode = true
	}

	if t := os.Getenv("TOKEN"); t != "" {
		BotToken = t
	}

	if t := os.Getenv("DB_HOST"); t != "" {
		DBHost = t
	}
	if t := os.Getenv("DB_NAME"); t != "" {
		DBName = t
	}
	if t := os.Getenv("DB_PORT"); t != "" {
		DBPort = t
	}
	if t := os.Getenv("DB_PASSWORD"); t != "" {
		DBPassword = t
	}
	if t := os.Getenv("DB_USER"); t != "" {
		DBUser = t
	}
}
