package main

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/httplib"
)

var (
	errNoSearchResults = errors.New("no search results")
)

func getPoster(query string) ([]byte, error) {
	posters, err := getPosterURL(query)
	if err != nil {
		return []byte{}, err
	}
	if len(posters) == 0 {
		return []byte{}, errNoSearchResults
	}

	request := httplib.Get(posters[0])
	return request.Bytes()
}

func getPosterURL(query string) ([]string, error) {
	url := "http://api.douban.com/v2/movie/search?q=" + query
	request := httplib.Get(url)
	bytes, err := request.Bytes()
	if err != nil {
		return []string{}, err
	}

	var response struct {
		Subjects []struct {
			Images struct {
				Large string `json:"large" bson:"large"`
			} `json:"images" bson:"images"`
		} `json:"subjects" bson:"subjects"`
	}
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return []string{}, err
	}

	posters := []string{}
	for _, subject := range response.Subjects {
		posters = append(posters, subject.Images.Large)
	}
	return posters, nil
}
