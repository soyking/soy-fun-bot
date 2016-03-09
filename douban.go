package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/soyking/telegram-bot-api"
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

func getSubjects(query string) ([]DoubanSubject, error) {
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

func getPoster(query string) (string, error) {
	subjects, err := getSubjects(query)
	if err != nil {
		return "", err
	}

	if len(subjects) == 0 || subjects[0].Images.Large == "" {
		return "", errNoSearchResults
	}

	return subjects[0].Images.Large, nil
}

func getMovieArticles(query, inlineQueryID string) (*tgbotapi.InlineConfig, error) {
	config := tgbotapi.InlineConfig{InlineQueryID: inlineQueryID}
	subjects, err := getSubjects(query)
	if err != nil {
		return nil, err
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
			article.Description = article.Description + topics
		}

		article.Title = subject.Title + " (" + subject.Year + ")  " + fmt.Sprintf("评分:%.1f", subject.Rating.Average)
		article.URL = subject.URL
		article.ThumbURL = subject.Images.Large
		article.MessageText = "<a href=\"" + subject.URL + "\">" + subject.Title + "</a>\n<pre>" + article.Description + "</pre>"
		article.ParseMode = "HTML"
		config.Results = append(config.Results, &article)
	}

	return &config, nil
}
