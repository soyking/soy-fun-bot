package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"net/url"
	"strings"
)

func getWeatherMsg(latitude, longitude float32) (string, error) {
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
		return "", err
	}

	if len(data.Results) == 0 || len(data.Results[0].WeatherData) == 0 {
		return "这个地方我可能没去过 :(", nil
	}

	text := "<b>城市</b>：%s\n<b>实时天气</b>：%s\n<b>今日天气</b>：%s\n<b>PM2.5</b>：%s"
	currentTemperature := strings.Split(data.Results[0].WeatherData[0].Date, "：")[1]
	currentTemperature = currentTemperature[:len(currentTemperature)-1]
	text = fmt.Sprintf(text,
		data.Results[0].CurrentCity,
		currentTemperature,
		data.Results[0].WeatherData[0].Temperature,
		data.Results[0].PM25)
	return text, nil
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
