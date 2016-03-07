package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"net/url"
)

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
