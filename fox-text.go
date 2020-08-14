package main

import (
	"encoding/json"
	"github.com/ChimeraCoder/anaconda"
	"io/ioutil"
	"math/rand"
	"time"
)

// FoxText : array of tweet text
type FoxText struct {
	Body []string `json:"text"`
}

// PostTextTweet : post tweet text
func PostTextTweet(api *anaconda.TwitterApi) {
	// body : array of tweet body
	var body FoxText

	jsonData, err := ioutil.ReadFile("./tweet_body.json")
	if err != nil {
		panic(err.Error())
	}

	json.Unmarshal(jsonData, &body)

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(body.Body))
	text := body.Body[n]
	tweet, err := api.PostTweet(text+"\n(bot)", nil)

	if err != nil {
		panic(err)
	}

	print(tweet.Text)
	return
}
