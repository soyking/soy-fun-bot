package main

import (
	"fmt"
	"github.com/soyking/telegram-bot-api"
	"log"
	"strings"
)

var (
	botToken = "***"
	bot      *tgbotapi.BotAPI
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
		m := &update.Message
		command := m.Command()

		switch command {
		case "/start":
			text := `本机器人提供查询天气功能~！发送你的地址来吧
- /movie [电影名] 可以收到一张电影海报

感谢 Baidu LBS API、Douban API
			`
			msg := tgbotapi.NewMessage(m.Chat.ID, text)
			sendAndCheck(msg)

		case "/movie":
			query := m.CommandArguments()
			if query == "" {
				msg := tgbotapi.NewMessage(m.Chat.ID, "告诉我电影的名字哦:)")
				sendAndCheck(msg)
			} else {
				poster, err := getPoster(query)
				if err != nil {
					var text string
					if err == errNoSearchResults {
						text = "找不到呀找不到:("
					} else {
						sendErr(m, err)
					}
					msg := tgbotapi.NewMessage(m.Chat.ID, text)
					sendAndCheck(msg)
				} else {
					msg := tgbotapi.NewPhotoUpload(m.Chat.ID, tgbotapi.FileBytes{Bytes: poster, Name: "poster.jpg"})
					sendAndCheck(msg)
				}
			}

		default:
			latitude := m.Location.Latitude
			longitude := m.Location.Longitude
			if latitude != 0.0 && longitude != 0.0 {
				var text string
				var data struct {
					Results []struct {
						CurrentCity string `json:"currentCity" bson:"currentCity"`
						PM25        string `json:"pm25" bson:"pm25"`
						WeatherData []struct {
							Date        string `json:"date" bson:"date"`
							Temperature string `json:"temperature" bson:"temperature"`
							Weather     string `json:"weather" bson:"weather"`
						} `json:"weather_data" bson:"weather_data"`
					} `json:"results" bson:"results"`
				}
				err := baiduWeather(latitude, longitude, &data)
				if err != nil {
					sendErr(m, err)
					continue
				} else if len(data.Results) == 0 || len(data.Results[0].WeatherData) == 0 {
					text = "这个地方我可能没去过 :("
				} else {
					text = "<b>城市</b>：%s\n<b>实时天气</b>：%s\n<b>今日天气</b>：%s\n<b>PM2.5</b>：%s"
					currentTemperature := strings.Split(data.Results[0].WeatherData[0].Date, "：")[1]
					currentTemperature = currentTemperature[:len(currentTemperature)-1]
					text = fmt.Sprintf(text,
						data.Results[0].CurrentCity,
						currentTemperature,
						data.Results[0].WeatherData[0].Temperature,
						data.Results[0].PM25)
				}

				msg := tgbotapi.NewMessage(m.Chat.ID, text)
				msg.ParseMode = "HTML"
				sendAndCheck(msg)
			} else {
				if update.InlineQuery.ID != "" {
					config := tgbotapi.InlineConfig{InlineQueryID: update.InlineQuery.ID}
					subjects, err := getPosterURL(update.InlineQuery.Query)
					if err != nil {
						sendErr(m, err)
						continue
					}
					for i, subject := range subjects {
						article := tgbotapi.InlineQueryResultArticle{}
						article.ID = fmt.Sprint("%d", i)
						if len(subject.Directors) > 0 {
							directors := "导演："
							for _, cast := range subject.Directors {
								directors = directors + cast.Name + " "
							}
							article.Description = article.Description + directors + "\n"
						}
						if len(subject.Casts) > 0 {
							casts := "主演："
							for _, cast := range subject.Casts {
								casts = casts + cast.Name + " "
							}
							article.Description = article.Description + casts + "\n"
						}
						if len(subject.Genres) > 0 {
							topics := "主题："
							for _, genre := range subject.Genres {
								topics = topics + genre + " "
							}
							article.Description = article.Description + topics + "\n"
						}
						article.Description += fmt.Sprintf("评分：%.1f", subject.Rating.Average)
						article.Title = subject.Title + " (" + subject.Year + ")"
						article.URL = subject.URL
						article.ThumbURL = subject.Images.Large
						article.MessageText = "<a href=\"" + subject.URL + "\">" + subject.Title + "</a>\n<pre>" + article.Description + "</pre>"
						article.ParseMode = "HTML"
						config.Results = append(config.Results, &article)
					}
					_, err = bot.AnswerInlineQuery(config)
					if err != nil {
						sendErr(m, err)
						continue
					}
				}
			}
		}
	}
}

func sendErr(m *tgbotapi.Message, err error) {
	log.Println(err)
	msg := tgbotapi.NewMessage(m.Chat.ID, "可能出了什么差错，本机器人也不知道该怎么办 :(")
	sendAndCheck(msg)
}

func sendAndCheck(msg tgbotapi.Chattable) {
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}
