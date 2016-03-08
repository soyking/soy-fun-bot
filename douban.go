package main

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/httplib"
)

var (
	errNoSearchResults = errors.New("no search results")
)

type DoubanSubject struct {
	Images struct {
		Large string `json:"large" bson:"large"`
	} `json:"images" bson:"images"`
	URL    string `json:"alt" bson:"alt"`
	Title  string `json:"title" bson:"title"`
	Year   string `json:"year" bson:"year"`
	Rating struct {
		Average float32 `json:"average" bson:"average"`
	} `json:"rating" bson:"rating"`
	Genres []string `json:"genres" bson:"genres"`
	Casts  []struct {
		Name string `json:"name" bson:"name"`
	} `json:"casts" bson:"casts"`
	Directors []struct {
		Name string `json:"name" bson:"name"`
	} `json:"directors" bson:"directors"`
}

func getPoster(query string) ([]byte, error) {
	subjects, err := getPosterURL(query)
	if err != nil {
		return []byte{}, err
	}
	if len(subjects) == 0 || subjects[0].Images.Large == "" {
		return []byte{}, errNoSearchResults
	}

	request := httplib.Get(subjects[0].Images.Large)
	return request.Bytes()
}

func getPosterURL(query string) ([]DoubanSubject, error) {
	url := "http://api.douban.com/v2/movie/search?start=0&count=3&q=" + query
	request := httplib.Get(url)
	bytes, err := request.Bytes()
	if err != nil {
		return nil, err
	}

	var response struct {
		Subjects []DoubanSubject `json:"subjects" bson:"subjects"`
	}

	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	return response.Subjects, nil
}
