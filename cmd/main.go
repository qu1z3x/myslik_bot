package main

import (
	"encoding/json"
	"fmt"
	"myslik_bot/config"
	"strings"
	"time"

	"gopkg.in/telebot.v4"
)

func main() {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  config.Config["TOKENs"].([]string)[1],
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	// var usersData = map[string]interface{}{}

	bot.Handle(telebot.OnText, func(context telebot.Context) error {
		message := context.Message()

		jsonData, err := json.Marshal(message)
		if err != nil {
			fmt.Println("Ошибка сериализации:", err)
		} else {
			fmt.Println(string(jsonData))
		}

		chatId := message.Chat.ID
		firstName := message.Chat.FirstName

		if strings.Contains(message.Text, "/") {
			switch message.Text {
			case "/start":
			case "/restart":
				firstMeeting(context)
			}
		}

		return context.Send(fmt.Sprintf("<b>Привет, %s</b>\n\nYour Id: %d", firstName, chatId), &telebot.SendOptions{
			ParseMode: telebot.ModeHTML,
		})
	})

	bot.Start()
}

func firstMeeting(context interface{}) {

}

func menu(context interface{}) {

}
