package main

import (
	"context"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	gogpt "github.com/sashabaranov/go-gpt3"
	"log"
	"os"
)

func telegramBot() func(func(string) string) {
	tokenBot := os.Getenv("token_telegram_bot")
	if tokenBot == "" {
		log.Panic("TokenTelegramBot not found.")
	}
	session, err := tgbotapi.NewBotAPI(tokenBot) // используя токен создаем новый инстанс бота
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", session.Self.UserName)
	newUpdate := tgbotapi.NewUpdate(0) // получчаем апдейты
	newUpdate.Timeout = 60
	chanUpdates, err := session.GetUpdatesChan(newUpdate) // используя конфиг u создаем канал в который будут прилетать новые сообщения
	// в канал chanUpdates прилетают структуры типа Update. Читываем их и обрабатываем
	log.Print("Service started")
	return func(sessionChatGPT func(string) string) {
		var respBot string
		for reqBot := range chanUpdates {
			switch {
			case reqBot.Message == nil:
				continue
			case reqBot.Message.IsCommand():
				switch reqBot.Message.Command() { // обрабатываем команды
				case "start":
					respBot = "1 2 3 Go..."
				}
			default:
				respBot = sessionChatGPT(reqBot.Message.Text)
			}
			msgBot := tgbotapi.NewMessage(reqBot.Message.Chat.ID, respBot) // формируем ответ
			log.Printf("[%s]->[%s] %s", reqBot.Message.From.UserName, session.Self.UserName, reqBot.Message.Text)
			log.Printf("[%s]->[%s] %s", session.Self.UserName, reqBot.Message.From.UserName, respBot)
			session.Send(msgBot) // отправляем копию сообщения
		}
	}
}

func chatGPT() func(string) string {
	tokenChat := os.Getenv("token_chatgpt")
	if tokenChat == "" {
		log.Panic("TokenChatGPT not found.")
	}
	c := gogpt.NewClient(tokenChat)
	ctx := context.Background()

	return func(msg string) string {
		req := gogpt.CompletionRequest{
			Model:       gogpt.GPT3TextDavinci003,
			Prompt:      msg,
			MaxTokens:   1000,
			Temperature: 1.0,
		}
		resp, err := c.CreateCompletion(ctx, req)
		if err != nil {
			return "Difficult... Say it again!"
		}
		return resp.Choices[0].Text
	}
}

func main() {
	sessionTelegramBot := telegramBot()
	sessionChatGPT := chatGPT()
	sessionTelegramBot(sessionChatGPT)
}
