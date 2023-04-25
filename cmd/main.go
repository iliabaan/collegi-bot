package main

import (
	"collegi-bot/authentication"
	"collegi-bot/command"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	// Set up the bot
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

	// Set up an update listener
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Fatal(err)
	}

	spotifyClient, err := authentication.Spotify()
	if err != nil {
		return
	}

	yaClient := authentication.YaMusic()

	// Handle incoming updates
	for update := range updates {
		if update.Message == nil { // ignore non-message updates
			continue
		}

		if err != nil {
			return
		}

		// handle commands
		if update.Message.IsCommand() {
			fmt.Println("Command:", update.Message.Command())
			switch update.Message.Command() {
			case "add":
				if update.Message.CommandArguments() != "" {
					command.AddYandex(ctx, update.Message.CommandArguments(), yaClient)
					track := command.AddSpotify(update.Message.CommandArguments(), *spotifyClient)
					if err != nil {
						return
					}
					msgText := fmt.Sprintf("Добавлен трек: %s – %s в плейлист!", track.Name, track.Artists[0].Name)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
					_, err = bot.Send(msg)

				}
			}
		}
	}
}
