package main

import (
	"encoding/json"
	"fmt"
	"io"
	"myslik_bot/config"
	"net/http"
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

	// var usersData

	bot.Handle(telebot.OnText, func(context telebot.Context) error {
		message := context.Message()
		jsonMessage, err := json.MarshalIndent(message, "", "  ")
		if err != nil {
			return err
		} else {
			fmt.Println(string(jsonMessage))
		}

		chatId := message.Chat.ID
		text := message.Text

		if !strings.Contains(text, "/") {

			response, err := getResponse(text)
			if err != nil {
				return err
			}
			context.Send(response)

		} else {
			switch text {
			case "/start":
				firstMeeting(chatId, context)
				return nil
			case "/restart":
			case "/menu":
				menu(chatId, context)
				return nil
			}

		}
		return nil
	})

	bot.Handle(telebot.OnCallback, func(context telebot.Context) error {
		query := context.Callback()
		jsonQuery, err := json.MarshalIndent(query, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(jsonQuery))

		chatId := query.Message.Chat.ID
		data := strings.TrimPrefix(query.Data, "\f")

		switch data {
		case "menu":
			menu(chatId, context)
			return nil
		}
		return nil
	})

	bot.Start()
}

func getResponse(request string) (string, error) {

	const url string = "https://openrouter.ai/api/v1/chat/completions"
	var headers = map[string]string{
		"Authorization": "Bearer " + config.Config["metaKey"].(string),
		"Content-Type":  "application/json",
	}

	payload := map[string]interface{}{
		"model": "meta-llama/llama-4-maverick",
		"messages": []interface{}{map[string]string{
			"role":    "system",
			"content": "Общайся на ты и используй разговорный стиль общения.",
		}, map[string]string{
			"role":    "user",
			"content": request,
		}},
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return "", err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err)
		return "", err
	}

	var respObj struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	json.Unmarshal(body, &respObj)

	return respObj.Choices[0].Message.Content, nil

}

func firstMeeting(chatId int64, context telebot.Context) {
	context.Send(fmt.Sprintf("Привет <b>%s</b>!\n<b>Это Мыслик!</b>\n \nНапиши /menu", context.Message().Chat.FirstName), &telebot.SendOptions{
		ParseMode: telebot.ModeHTML})
}

func menu(chatId int64, context telebot.Context) {
	context.Send("Меню:\n\n1. <b>Пункт 1</b>\n2. <b>Пункт 2</b>\n3. <b>Пункт 3</b>", &telebot.SendOptions{
		ParseMode: telebot.ModeHTML,
		ReplyMarkup: &telebot.ReplyMarkup{
			InlineKeyboard: [][]telebot.InlineButton{{telebot.InlineButton{Unique: "menu", Text: "Переотправить меню"}}},
		},
	})
}
