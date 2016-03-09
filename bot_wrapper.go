package main

import (
	"github.com/astaxie/beego/httplib"
	"github.com/soyking/telegram-bot-api"
	"log"
)

type botWrapper struct {
	bot *tgbotapi.BotAPI
	m   *tgbotapi.Message
}

func NewBotWrapper(m *tgbotapi.Message) *botWrapper {
	return &botWrapper{bot: bot, m: m}
}

func (b *botWrapper) Send(msg tgbotapi.Chattable) {
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func (b *botWrapper) SendText(text string) {
	msg := tgbotapi.NewMessage(b.m.Chat.ID, text)
	msg.ParseMode = "HTML"
	b.Send(msg)
}

func (b *botWrapper) SendErr(err error) {
	log.Println(err)
	b.SendText("可能出了什么差错，本机器人也不知道该怎么办 :(")
}

func (b *botWrapper) SendImage(imgURL string) {
	request := httplib.Get(imgURL)
	bytes, err := request.Bytes()
	if err != nil {
		log.Println(err)
		return
	}

	msg := tgbotapi.NewPhotoUpload(b.m.Chat.ID, tgbotapi.FileBytes{Bytes: bytes, Name: "image.jpg"})
	b.Send(msg)
}

func (b *botWrapper) SendInlineAnswer(config tgbotapi.InlineConfig) {
	_, err := bot.AnswerInlineQuery(config)
	if err != nil {
		log.Println(err)
	}
}
