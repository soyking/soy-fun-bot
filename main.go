package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"gopkg.in/telegram-bot-api.v1"
	"log"
	"net/url"
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
		command := update.Message.Command()

		switch command {
		case "/start":
			text := `本机器人提供查询天气功能~！发送你的地址来吧

感谢 Baidu LBS API
			`
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			msg.ParseMode = "HTML"
			sendAndCheck(msg)
		default:
			latitude := update.Message.Location.Latitude
			longitude := update.Message.Location.Longitude
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
					log.Println(err)
					text = "可能出了什么差错，本机器人也不知道该怎么办 :("
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

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
				msg.ParseMode = "HTML"
				sendAndCheck(msg)
			}
		}
	}
}

func sendAndCheck(msg tgbotapi.Chattable) {
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func baiduWeather(latitude float32, longitude float32, data interface{}) error {
	basePath := fmt.Sprintf("/telematics/v3/weather?location=%f,%f&output=json", longitude, latitude)
	queryURL := baiduLBSAPI(basePath)

	request := httplib.Get(queryURL)
	bytes, err := request.Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, data)
}

var (
	baiduLBSAK = "***"
	baiduLBSSK = "***"
)

func baiduLBSAPI(basePath string) string {
	sn := baiduCalculateSN(basePath)
	return fmt.Sprintf("http://api.map.baidu.com%s&ak=%s&sn=%s", basePath, baiduLBSAK, sn)
}

// 计算百度LBS API 中的 SN 值
func baiduCalculateSN(path string) string {
	path = path + "&ak=" + baiduLBSAK + baiduLBSSK
	encodedStr := url.QueryEscape(path)
	return fmt.Sprintf("%x", md5.Sum([]byte(encodedStr)))
}
