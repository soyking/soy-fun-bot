package main

import (
	"github.com/soyking/telegram-bot-api"
	"log"
)

var (
	bot *tgbotapi.BotAPI
)

func main() {
	var err error
	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		b := NewBotWrapper(&update.Message)
		command := update.Message.Command()
		query := update.Message.CommandArguments()

		switch command {
		case "/start":
			text := `本机器人并不知道存在的意义
- 发送地址可以得到天气预报
- /movie [电影名] 可以收到一张电影海报

感谢 Baidu LBS API、Douban API
			`
			b.SendText(text)
			continue

		case "/movie":
			if query == "" {
				b.SendText("告诉我电影的名字哦:)")
				continue
			}

			posterURL, err := getPoster(query)
			if err != nil {
				if err == errNoSearchResults {
					b.SendText("找不到呀找不到:(")
				}
				b.SendErr(err)
				continue
			}

			b.SendImage(posterURL)
			continue
		}

		// 发送位置
		latitude := update.Message.Location.Latitude
		longitude := update.Message.Location.Longitude
		if latitude != 0.0 && longitude != 0.0 {
			text, err := getWeatherMsg(latitude, longitude)
			if err != nil {
				b.SendErr(err)
				continue
			}

			b.SendText(text)
			continue
		}

		// @SoyFunBot 电影
		if update.InlineQuery.ID != "" {
			config, err := getMovieArticles(update.InlineQuery.Query, update.InlineQuery.ID)
			if err != nil {
				b.SendErr(err)
				continue
			}

			b.SendInlineAnswer(*config)
		}

	}
}
