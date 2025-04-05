package main

import (
	"dialogue/internal/store/sqlstore"
	"flag"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	databaseURL = "postgres://postgres:world555@localhost:5431/rpg?sslmode=disable"
)

var (
	token   string
	buttons = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("List of stories"),
			tgbotapi.NewKeyboardButton("Create a story"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("FAQ"),
			tgbotapi.NewKeyboardButton("Donate"),
		),
	)
	keyboard = tgbotapi.ReplyKeyboardMarkup{
		Keyboard:        buttons.Keyboard,
		OneTimeKeyboard: true,
	}
)

func init() {
	flag.StringVar(&token, "token", "", "input bot token here")
}

func connectDB(databaseURL string) (db *gorm.DB, err error) {
	db, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	flag.Parse()

	db, err := connectDB(databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	store := sqlstore.New(db)
	stories = &StoryStruct{
		db: store,
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		if update.Message == nil {
			continue
		}

		switch update.Message.Text {
		case "List of stories":
			stories, err := stories.db.Story().GetAllStories()
			if err != nil {
				continue
			}

			var storyButtons []tgbotapi.InlineKeyboardButton
			for _, story := range stories {
				storyButtons = append(storyButtons, tgbotapi.NewInlineKeyboardButtonSwitch(story.StoryTitle, fmt.Sprintf("story_%d", story.StoryID)))
			}

			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Select a story:")
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(storyButtons)
		}

		switch update.Message.Command() {
		case "start":
			msg.ReplyMarkup = keyboard
		case "close":
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		}

		if _, err := bot.Send(msg); err != nil {
			continue
		}
	}
}
